package ai

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/pgvector/pgvector-go"
	"go.uber.org/zap"
)

// DocumentChunk is a segment of user-uploaded content stored for RAG retrieval.
type DocumentChunk struct {
	ID        uuid.UUID `json:"id"`
	UserID    uuid.UUID `json:"user_id"`
	Content   string    `json:"content"`
	Source    string    `json:"source"` // original file name
	ChunkIdx  int       `json:"chunk_idx"`
	CreatedAt time.Time `json:"created_at"`
}

const (
	maxChunkTokens   = 800 // approximate token target per chunk
	chunksPerQuery   = 5   // max context chunks injected into prompt
	avgCharsPerToken = 4   // rough char:token ratio
)

// RAGEngine handles document storage and pgvector similarity retrieval.
type RAGEngine struct {
	db        *sql.DB
	apiClient *QwenClient
	log       *zap.Logger
}

func NewRAGEngine(db *sql.DB, apiClient *QwenClient, log *zap.Logger) *RAGEngine {
	return &RAGEngine{db: db, apiClient: apiClient, log: log}
}

// StoreChunks splits text into chunks and persists them for a user.
func (r *RAGEngine) StoreChunks(ctx context.Context, userID uuid.UUID, text, source string) (int, error) {
	chunks := chunkText(text, maxChunkTokens*avgCharsPerToken)
	if len(chunks) == 0 {
		return 0, nil
	}

	embeddings, err := r.apiClient.GenerateEmbeddings(ctx, chunks)
	if err != nil {
		r.log.Warn("failed to generate embeddings, chunks will be stored without vectors", zap.Error(err))
	}

	stored := 0
	for idx, chunk := range chunks {
		var emb interface{}
		if idx < len(embeddings) && embeddings[idx] != nil {
			emb = pgvector.NewVector(embeddings[idx])
		}

		const q = `
			INSERT INTO document_chunks (user_id, content, source, chunk_idx, embedding, created_at)
			VALUES ($1, $2, $3, $4, $5, NOW())`

		if _, err := r.db.ExecContext(ctx, q, userID, chunk, source, idx, emb); err != nil {
			return stored, fmt.Errorf("storing chunk %d: %w", idx, err)
		}
		stored++
	}

	r.log.Info("stored document chunks with vectors",
		zap.String("user_id", userID.String()),
		zap.String("source", source),
		zap.Int("chunks", stored),
	)

	return stored, nil
}

// RetrieveContext fetches the most relevant chunks using exact Cosine Similarity.
func (r *RAGEngine) RetrieveContext(ctx context.Context, userID uuid.UUID, query string) (string, int, error) {
	if strings.TrimSpace(query) == "" {
		return "", 0, nil
	}

	// 1. Generate query embedding
	qEmb, err := r.apiClient.GenerateEmbeddings(ctx, []string{query})
	if err != nil || len(qEmb) == 0 {
		return "", 0, fmt.Errorf("failed to generate query embedding: %w", err)
	}

	vec := pgvector.NewVector(qEmb[0])

	// 2. Vector search (Cosine distance: <-> returns distance, we want lowest distance)
	// We only want chunks where distance < 0.6 (or similarity > 0.4)
	q := `
		SELECT content, source
		FROM document_chunks
		WHERE user_id = $1 AND embedding IS NOT NULL
		ORDER BY embedding <=> $2
		LIMIT $3`

	rows, err := r.db.QueryContext(ctx, q, userID, vec, chunksPerQuery)
	if err != nil {
		return "", 0, fmt.Errorf("retrieving chunks: %w", err)
	}
	defer rows.Close()

	var parts []string
	chunksFound := 0
	for rows.Next() {
		var content, source string
		if err := rows.Scan(&content, &source); err != nil {
			continue
		}
		parts = append(parts, fmt.Sprintf("[Source: %s]\n%s", source, content))
		chunksFound++
	}
	if err := rows.Err(); err != nil {
		return "", chunksFound, fmt.Errorf("reading chunk rows: %w", err)
	}

	return strings.Join(parts, "\n\n"), chunksFound, nil
}

// DeleteUserChunks removes all stored chunks for a user (e.g. on account delete).
func (r *RAGEngine) DeleteUserChunks(ctx context.Context, userID uuid.UUID) error {
	_, err := r.db.ExecContext(ctx, `DELETE FROM document_chunks WHERE user_id = $1`, userID)
	return err
}

// ListUserSources returns distinct source filenames for a user.
func (r *RAGEngine) ListUserSources(ctx context.Context, userID uuid.UUID) ([]string, error) {
	rows, err := r.db.QueryContext(ctx,
		`SELECT DISTINCT source FROM document_chunks WHERE user_id = $1 ORDER BY source`,
		userID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var sources []string
	for rows.Next() {
		var s string
		if err := rows.Scan(&s); err != nil {
			continue
		}
		sources = append(sources, s)
	}
	return sources, rows.Err()
}

// ──────────────────────────────────────────────────────────────────────────────
// Helpers
// ──────────────────────────────────────────────────────────────────────────────

// chunkText splits text into chunks of approximately maxChars characters,
// respecting sentence boundaries where possible.
func chunkText(text string, maxChars int) []string {
	text = strings.TrimSpace(text)
	if len(text) == 0 {
		return nil
	}

	// Split on double-newline blocks first (paragraphs).
	paragraphs := strings.Split(text, "\n\n")
	var chunks []string
	var current strings.Builder

	for _, para := range paragraphs {
		para = strings.TrimSpace(para)
		if para == "" {
			continue
		}

		if current.Len()+len(para)+2 > maxChars && current.Len() > 0 {
			chunks = append(chunks, strings.TrimSpace(current.String()))
			current.Reset()
		}

		if len(para) > maxChars {
			// Large paragraph: split by sentences.
			sentences := splitSentences(para)
			for _, sent := range sentences {
				if current.Len()+len(sent)+1 > maxChars && current.Len() > 0 {
					chunks = append(chunks, strings.TrimSpace(current.String()))
					current.Reset()
				}
				current.WriteString(sent)
				current.WriteByte(' ')
			}
		} else {
			if current.Len() > 0 {
				current.WriteString("\n\n")
			}
			current.WriteString(para)
		}
	}

	if current.Len() > 0 {
		chunks = append(chunks, strings.TrimSpace(current.String()))
	}

	return chunks
}

func splitSentences(text string) []string {
	// Simple sentence splitter on ". ", "! ", "? "
	var sents []string
	start := 0
	for i := 0; i+1 < len(text); i++ {
		if (text[i] == '.' || text[i] == '!' || text[i] == '?') && text[i+1] == ' ' {
			sents = append(sents, text[start:i+1])
			start = i + 2
		}
	}
	if start < len(text) {
		sents = append(sents, text[start:])
	}
	return sents
}

// (End of RAG Engine helpers)
