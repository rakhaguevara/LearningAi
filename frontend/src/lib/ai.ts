/**
 * ai.ts — LearnNow AI Workspace API client
 *
 * All AI calls proxy through the Go backend — the API key is NEVER
 * exposed to the frontend. Every method attaches the JWT Bearer token.
 */

import { authHeaders } from './auth';
import { siteConfig } from './constants';

const BASE = siteConfig.api;

// ─────────────────────────────────────────────────────────────────────────────
// Types
// ─────────────────────────────────────────────────────────────────────────────

export type OutputFormat =
    | 'quick_summary'
    | 'detailed_explanation'
    | 'anime_style'
    | 'sports_analogy'
    | 'academic_formal'
    | 'presentation_slides'
    | 'audio_explanation'
    | 'translation';

export interface OutputFormatOption {
    id: OutputFormat;
    label: string;
    emoji: string;
    description: string;
}

export const OUTPUT_FORMATS: OutputFormatOption[] = [
    { id: 'quick_summary', label: 'Quick Summary', emoji: '📝', description: 'Concise bullet points' },
    { id: 'detailed_explanation', label: 'Detailed Explanation', emoji: '📖', description: 'Full deep-dive' },
    { id: 'anime_style', label: 'Anime Style', emoji: '🎌', description: 'Vivid story analogies' },
    { id: 'sports_analogy', label: 'Sports Analogy', emoji: '⚽', description: 'Sport metaphors' },
    { id: 'academic_formal', label: 'Academic Formal', emoji: '🎓', description: 'Formal scholarly tone' },
    { id: 'presentation_slides', label: 'Presentation Slides', emoji: '📊', description: 'Downloadable slides' },
    { id: 'audio_explanation', label: 'Audio Script', emoji: '🔊', description: 'Spoken-word format' },
    { id: 'translation', label: 'Translation', emoji: '🌍', description: 'Translate to any language' },
];

export interface AskRequest {
    question: string;
    output_format?: OutputFormat;
    session_id?: string;
    target_language?: string;
}

// ── Typed sub-struct interfaces (mirror Go orchestrator.go types) ─────────────

export interface AISummary {
    title: string;
    core_concept_explanation: string;
    key_points: string[];
    real_world_example: string;
    short_conclusion: string;
}

export interface AIDetailed {
    title: string;
    concept_explanation: string;
    step_by_step_breakdown: string[];
    example: string;
    mini_quiz: string[];
}

export interface AIAnime {
    episode_title: string;
    main_character: string;
    story_arc: string;
    physics_explanation: string;
    visual_scene_prompt: string;
}

export interface AISports {
    game_title: string;
    sport_used: string;
    play_breakdown: string;
    coaching_tip: string;
    scoreboard_summary: string;
    visual_scene_prompt: string;
}

export interface AIAcademic {
    title: string;
    abstract: string;
    theoretical_background: string;
    methodology: string;
    conclusion: string;
}

export interface AskResponse {
    answer: string;
    output_format: OutputFormat;
    tokens_used: number;
    latency_ms: number;
    needs_format?: boolean;
    format_prompt?: string;
    context_used: boolean;
    context_found?: number;
    download_url?: string;
    illustration_url?: string;
    is_structured_json?: boolean;
    // Typed sub-structs — populated by the orchestrator
    summary?: AISummary;
    detailed?: AIDetailed;
    anime?: AIAnime;
    sports?: AISports;
    academic?: AIAcademic;
}

export interface UploadResponse {
    source: string;
    file_type: string;
    word_count: number;
    chunks: number;
    message: string;
}

export interface PPTResponse {
    file_name: string;
    slide_count: number;
    download_url: string;
}

export interface AudioResponse {
    file_name: string;
    format: string;
    duration_sec: number;
    download_url: string;
}

export interface TranslateRequest {
    text: string;
    target_lang: string;
}

export interface TranslateResponse {
    translated: string;
    source_lang: string;
    tokens_used: number;
}

export interface SourcesResponse {
    sources: string[];
}

// ─────────────────────────────────────────────────────────────────────────────
// API helpers
// ─────────────────────────────────────────────────────────────────────────────

async function apiFetch<T>(
    path: string,
    init: RequestInit = {}
): Promise<T> {
    let res: Response;
    try {
        res = await fetch(`${BASE}${path}`, {
            ...init,
            headers: {
                ...authHeaders(),
                ...(init.headers as Record<string, string> || {}),
            },
        });
    } catch {
        throw new Error('Cannot reach server. Check your connection or Docker containers.');
    }

    // Safely parse JSON — backend may return HTML on fatal gateway errors
    let json: any;
    try {
        json = await res.json();
    } catch {
        throw new Error(`Server returned a non-JSON response (HTTP ${res.status}). Make sure the backend is running.`);
    }

    if (!res.ok || !json.success) {
        const msg: string = json?.error?.message || `Request failed (HTTP ${res.status})`;
        throw new Error(msg);
    }

    return json.data as T;
}

// ─────────────────────────────────────────────────────────────────────────────
// AI API
// ─────────────────────────────────────────────────────────────────────────────

/**
 * Send a question to the AI. If `output_format` is omitted, the backend
 * returns `needs_format: true` and a prompt asking the user to pick a format.
 * Both cases are wrapped in {success: true, data: AskResponse} by the backend.
 */
export async function askAI(req: AskRequest): Promise<AskResponse> {
    return apiFetch<AskResponse>('/ai/ask', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify(req),
    });
}


/**
 * Upload a file for RAG indexing.
 * Supports PDF, DOCX, TXT, JPG/PNG/WEBP images.
 */
export async function uploadFile(file: File): Promise<UploadResponse> {
    const form = new FormData();
    form.append('file', file);

    const res = await fetch(`${BASE}/ai/upload`, {
        method: 'POST',
        headers: authHeaders(),
        body: form,
    });

    const json = await res.json();
    if (!res.ok || !json.success) {
        throw new Error(json?.error?.message || `Upload failed (${res.status})`);
    }
    return json.data as UploadResponse;
}

/**
 * Fetch the list of uploaded sources for the authenticated user.
 */
export async function getSources(): Promise<string[]> {
    const data = await apiFetch<SourcesResponse>('/ai/sources');
    return data.sources ?? [];
}

/**
 * Generate a downloadable presentation (HTML slideshow) from a topic.
 */
export async function generatePPT(topic: string, content: string): Promise<PPTResponse> {
    return apiFetch<PPTResponse>('/ai/generate-ppt', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ topic, content }),
    });
}

/**
 * Convert text to speech. Returns download URL for the audio file.
 */
export async function generateAudio(text: string, voice?: string): Promise<AudioResponse> {
    return apiFetch<AudioResponse>('/ai/generate-audio', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ text, voice }),
    });
}

/**
 * Translate text using Qwen.
 */
export async function translateText(req: TranslateRequest): Promise<TranslateResponse> {
    return apiFetch<TranslateResponse>('/ai/translate', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify(req),
    });
}

/**
 * Build a full download URL for PPT/audio files.
 */
export function buildDownloadURL(path: string): string {
    return `${BASE}${path}`;
}
