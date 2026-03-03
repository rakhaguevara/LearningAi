package ai

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"regexp"
	"strings"
	"time"
	"unicode/utf8"

	"go.uber.org/zap"
)

const maxFileSizeBytes = 10 * 1024 * 1024 // 10 MB

// SupportedMIMETypes maps MIME type to a parser tag.
var SupportedMIMETypes = map[string]string{
	"application/pdf": "pdf",
	"application/vnd.openxmlformats-officedocument.wordprocessingml.document": "docx",
	"text/plain": "txt",
	"image/jpeg": "image",
	"image/png":  "image",
	"image/webp": "image",
	"image/gif":  "image",
}

// ParsedDocument holds the extracted, sanitised text and metadata.
type ParsedDocument struct {
	Text      string `json:"text"`
	WordCount int    `json:"word_count"`
	Source    string `json:"source"`
	FileType  string `json:"file_type"`
}

// FileParser orchestrates extraction from different file formats.
type FileParser struct {
	qwenClient *QwenClient
	log        *zap.Logger
}

func NewFileParser(qwenClient *QwenClient, log *zap.Logger) *FileParser {
	return &FileParser{qwenClient: qwenClient, log: log}
}

// ParseUpload validates and extracts text from an uploaded file.
func (p *FileParser) ParseUpload(ctx context.Context, header *multipart.FileHeader, file multipart.File) (*ParsedDocument, error) {
	if header.Size > maxFileSizeBytes {
		return nil, fmt.Errorf("file exceeds maximum size of 10 MB (got %.2f MB)", float64(header.Size)/1024/1024)
	}

	data, err := io.ReadAll(io.LimitReader(file, maxFileSizeBytes+1))
	if err != nil {
		return nil, fmt.Errorf("reading file: %w", err)
	}
	if len(data) > maxFileSizeBytes {
		return nil, fmt.Errorf("file exceeds maximum size of 10 MB")
	}

	contentType := header.Header.Get("Content-Type")
	if contentType == "" {
		contentType = http.DetectContentType(data)
	}
	contentType = strings.TrimSpace(strings.SplitN(contentType, ";", 2)[0])

	fileTag, supported := SupportedMIMETypes[contentType]
	if !supported {
		fileTag = extensionToTag(header.Filename)
		if fileTag == "" {
			return nil, fmt.Errorf("unsupported file type: %s", contentType)
		}
	}

	p.log.Info("parsing uploaded file",
		zap.String("filename", header.Filename),
		zap.String("type", fileTag),
		zap.Int64("size_bytes", header.Size),
	)

	var rawText string

	switch fileTag {
	case "txt":
		rawText = extractTXT(data)
	case "pdf":
		rawText, err = extractPDF(data)
	case "docx":
		rawText, err = extractDOCX(data)
	case "image":
		rawText, err = p.extractImageOCR(ctx, data, contentType, header.Filename)
	default:
		return nil, fmt.Errorf("no extractor for file type: %s", fileTag)
	}
	if err != nil {
		return nil, fmt.Errorf("extracting text from %s: %w", fileTag, err)
	}

	sanitised := SanitiseText(rawText)
	if strings.TrimSpace(sanitised) == "" {
		return nil, fmt.Errorf("could not extract meaningful text from file")
	}

	return &ParsedDocument{
		Text:      sanitised,
		WordCount: len(strings.Fields(sanitised)),
		Source:    header.Filename,
		FileType:  fileTag,
	}, nil
}

// ──────────────────────────────────────────────────────────────────────────────
// Format extractors
// ──────────────────────────────────────────────────────────────────────────────

func extractTXT(data []byte) string {
	if !utf8.Valid(data) {
		return strings.ToValidUTF8(string(data), "")
	}
	return string(data)
}

func extractPDF(data []byte) (string, error) {
	content := string(data)

	// Try to extract BT...ET text blocks (uncompressed streams)
	btRe := regexp.MustCompile(`BT([\s\S]*?)ET`)
	tjRe := regexp.MustCompile(`\(([^)]*)\)\s*T[jJ]`)

	var sb strings.Builder
	for _, block := range btRe.FindAllStringSubmatch(content, -1) {
		if len(block) < 2 {
			continue
		}
		for _, tm := range tjRe.FindAllStringSubmatch(block[1], -1) {
			if len(tm) > 1 {
				sb.WriteString(decodePDFString(tm[1]))
				sb.WriteByte(' ')
			}
		}
	}

	result := sb.String()
	if strings.TrimSpace(result) == "" {
		// Fallback: plain printable runs (handles simple/unencoded PDFs)
		printableRe := regexp.MustCompile(`[\x20-\x7E]{5,}`)
		matched := printableRe.FindAllString(content, -1)
		// Filter out PDF header noise
		var filtered []string
		for _, m := range matched {
			if !strings.Contains(m, "PDF") && !strings.Contains(m, "obj") {
				filtered = append(filtered, m)
			}
		}
		result = strings.Join(filtered, " ")
	}
	return result, nil
}

func decodePDFString(s string) string {
	octalRe := regexp.MustCompile(`\\([0-7]{1,3})`)
	s = octalRe.ReplaceAllStringFunc(s, func(m string) string {
		var n int
		fmt.Sscanf(m[1:], "%o", &n)
		return string(rune(n))
	})
	s = strings.ReplaceAll(s, `\n`, "\n")
	s = strings.ReplaceAll(s, `\r`, "")
	s = strings.ReplaceAll(s, `\t`, "\t")
	return s
}

func extractDOCX(data []byte) (string, error) {
	zipMagic := []byte{0x50, 0x4B, 0x03, 0x04}
	if !bytes.HasPrefix(data, zipMagic) {
		return "", fmt.Errorf("invalid DOCX: not a ZIP archive")
	}
	content := string(data)
	wtRe := regexp.MustCompile(`<w:t[^>]*>([^<]*)</w:t>`)
	matches := wtRe.FindAllStringSubmatch(content, -1)
	var sb strings.Builder
	for _, m := range matches {
		if len(m) > 1 {
			sb.WriteString(m[1])
			sb.WriteByte(' ')
		}
	}
	return sb.String(), nil
}

// extractImageOCR uses Qwen-VL multimodal API to OCR an image.
func (p *FileParser) extractImageOCR(ctx context.Context, data []byte, mimeType, filename string) (string, error) {
	if p.qwenClient == nil || p.qwenClient.apiKey == "" {
		return "", fmt.Errorf("Qwen API key not configured for OCR")
	}

	b64 := base64.StdEncoding.EncodeToString(data)
	dataURI := fmt.Sprintf("data:%s;base64,%s", mimeType, b64)

	prompt := "Extract ALL text from this image exactly as written. Return only the extracted text, preserving line breaks. If no text is present, return an empty string."

	type vlContent struct {
		Image string `json:"image,omitempty"`
		Text  string `json:"text,omitempty"`
	}
	type vlMessage struct {
		Role    string      `json:"role"`
		Content []vlContent `json:"content"`
	}
	type vlInput struct {
		Messages []vlMessage `json:"messages"`
	}
	type vlParameters struct {
		ResultFormat string `json:"result_format"`
	}
	type vlRequest struct {
		Model      string       `json:"model"`
		Input      vlInput      `json:"input"`
		Parameters vlParameters `json:"parameters"`
	}

	reqBody := vlRequest{
		Model: "qwen-vl-plus",
		Input: vlInput{
			Messages: []vlMessage{
				{
					Role: "user",
					Content: []vlContent{
						{Image: dataURI},
						{Text: prompt},
					},
				},
			},
		},
		Parameters: vlParameters{ResultFormat: "message"},
	}

	bodyBytes, err := json.Marshal(reqBody)
	if err != nil {
		return "", fmt.Errorf("marshalling VL request: %w", err)
	}

	ctx2, cancel := context.WithTimeout(ctx, 60*time.Second)
	defer cancel()

	text, err := p.qwenClient.doVLRequest(ctx2, bodyBytes)
	if err != nil {
		p.log.Warn("qwen OCR failed, returning empty text",
			zap.String("filename", filename),
			zap.Error(err),
		)
		return "", nil // non-fatal
	}

	p.log.Info("OCR completed",
		zap.String("filename", filename),
		zap.Int("text_length", len(text)),
	)
	return text, nil
}

// ──────────────────────────────────────────────────────────────────────────────
// Text Sanitisation — must be applied before any prompt injection
// ──────────────────────────────────────────────────────────────────────────────

var (
	controlCharRe = regexp.MustCompile(`[\x00-\x08\x0B\x0C\x0E-\x1F\x7F]`)
	whitespaceRe  = regexp.MustCompile(`[ \t]{2,}`)
	blankLineRe   = regexp.MustCompile(`\n{3,}`)
)

// SanitiseText removes dangerous control characters and collapses whitespace.
// Hard cap at 50,000 chars to limit prompt injection surface.
func SanitiseText(text string) string {
	text = controlCharRe.ReplaceAllString(text, " ")
	text = whitespaceRe.ReplaceAllString(text, " ")
	text = blankLineRe.ReplaceAllString(text, "\n\n")
	text = strings.TrimSpace(text)
	if len(text) > 50000 {
		text = text[:50000] + "\n\n[Content truncated at 50,000 characters]"
	}
	return text
}

// ──────────────────────────────────────────────────────────────────────────────
// Helpers
// ──────────────────────────────────────────────────────────────────────────────

func extensionToTag(filename string) string {
	lower := strings.ToLower(filename)
	switch {
	case strings.HasSuffix(lower, ".pdf"):
		return "pdf"
	case strings.HasSuffix(lower, ".docx"):
		return "docx"
	case strings.HasSuffix(lower, ".txt"), strings.HasSuffix(lower, ".md"):
		return "txt"
	case strings.HasSuffix(lower, ".jpg"), strings.HasSuffix(lower, ".jpeg"),
		strings.HasSuffix(lower, ".png"), strings.HasSuffix(lower, ".webp"),
		strings.HasSuffix(lower, ".gif"):
		return "image"
	}
	return ""
}

// doVLRequest sends a Qwen-VL multimodal request and returns the text output.
func (c *QwenClient) doVLRequest(ctx context.Context, body []byte) (string, error) {
	url := c.baseURL + "/services/aigc/multimodal-generation/generation"

	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(body))
	if err != nil {
		return "", fmt.Errorf("creating VL request: %w", err)
	}
	httpReq.Header.Set("Authorization", "Bearer "+c.apiKey)
	httpReq.Header.Set("Content-Type", "application/json")

	httpResp, err := c.httpClient.Do(httpReq)
	if err != nil {
		return "", fmt.Errorf("executing VL request: %w", err)
	}
	defer httpResp.Body.Close()

	respBytes, err := io.ReadAll(io.LimitReader(httpResp.Body, 2*1024*1024))
	if err != nil {
		return "", fmt.Errorf("reading VL response: %w", err)
	}
	if httpResp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("qwen VL http %d: %s", httpResp.StatusCode, string(respBytes))
	}

	var result struct {
		Output struct {
			Choices []struct {
				Message struct {
					Content []struct {
						Text string `json:"text"`
					} `json:"content"`
				} `json:"message"`
			} `json:"choices"`
		} `json:"output"`
		Code    string `json:"code"`
		Message string `json:"message"`
	}
	if err := json.Unmarshal(respBytes, &result); err != nil {
		return "", fmt.Errorf("decoding VL response: %w", err)
	}
	if result.Code != "" && result.Code != "200" {
		return "", fmt.Errorf("qwen VL error %s: %s", result.Code, result.Message)
	}
	if len(result.Output.Choices) == 0 {
		return "", fmt.Errorf("qwen VL returned no choices")
	}

	var texts []string
	for _, item := range result.Output.Choices[0].Message.Content {
		if item.Text != "" {
			texts = append(texts, item.Text)
		}
	}
	return strings.Join(texts, "\n"), nil
}
