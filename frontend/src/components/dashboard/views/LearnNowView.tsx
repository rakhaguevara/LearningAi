'use client';

import { useState, useRef, useEffect, useCallback } from 'react';
import { motion, AnimatePresence } from 'framer-motion';
import {
    askAI, uploadFile, getSources, generatePPT, generateAudio, translateText,
    buildDownloadURL, OUTPUT_FORMATS,
    type AskResponse, type OutputFormat, type OutputFormatOption,
} from '@/lib/ai';

// ─────────────────────────────────────────────────────────────────────────────
// Types
// ─────────────────────────────────────────────────────────────────────────────

type MessageRole = 'user' | 'ai' | 'system';

interface Message {
    id: string;
    role: MessageRole;
    text: string;
    outputFormat?: OutputFormat;
    contextUsed?: boolean;
    tokensUsed?: number;
    latencyMs?: number;
    downloadURL?: string;
    downloadLabel?: string;
    audioURL?: string;
    illustrationUrl?: string;
    isStructuredJson?: boolean;
    isError?: boolean;
    timestamp: Date;

    // Typed struct optional fields
    summary?: AskResponse['summary'];
    detailed?: AskResponse['detailed'];
    anime?: AskResponse['anime'];
    sports?: AskResponse['sports'];
    academic?: AskResponse['academic'];
}

interface Source {
    name: string;
    addedAt: Date;
}

// ─────────────────────────────────────────────────────────────────────────────
// Sub-components
// ─────────────────────────────────────────────────────────────────────────────

function TypingIndicator() {
    return (
        <div className="flex gap-3">
            <div className="w-7 h-7 rounded-full bg-gradient-to-br from-violet-500 to-fuchsia-500 flex items-center justify-center flex-shrink-0 shadow-md shadow-violet-500/30">
                <svg className="w-3.5 h-3.5 text-white" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                    <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2.2} d="M13 10V3L4 14h7v7l9-11h-7z" />
                </svg>
            </div>
            <div className="px-4 py-3 rounded-2xl rounded-tl-sm bg-[var(--bg-elevated)] border border-[var(--border)] flex items-center gap-1.5">
                {[0, 1, 2].map((i) => (
                    <motion.div
                        key={i}
                        className="w-2 h-2 rounded-full bg-violet-400"
                        animate={{ scale: [1, 1.4, 1], opacity: [0.5, 1, 0.5] }}
                        transition={{ duration: 1.2, repeat: Infinity, delay: i * 0.2 }}
                    />
                ))}
            </div>
        </div>
    );
}

// ─────────────────────────────────────────────────────────────────────────────
// Structured JSON renderer
// ─────────────────────────────────────────────────────────────────────────────

function StructuredAnswer({ msg }: { msg: Message }) {
    const format = msg.outputFormat;

    const Field = ({ label, value }: { label: string; value?: unknown }) => {
        if (value === undefined || value === null || value === '') return null;
        if (Array.isArray(value)) {
            return (
                <div className="mt-2">
                    <p className="text-[10px] font-semibold text-violet-400 uppercase tracking-widest mb-1">{label}</p>
                    <ul className="space-y-1">
                        {(value as string[]).map((v, i) => (
                            <li key={i} className="flex gap-2 text-sm text-[var(--text-secondary)]">
                                <span className="text-violet-400 mt-0.5">▸</span><span>{v}</span>
                            </li>
                        ))}
                    </ul>
                </div>
            );
        }
        return (
            <div className="mt-2">
                <p className="text-[10px] font-semibold text-violet-400 uppercase tracking-widest mb-0.5">{label}</p>
                <p className="text-sm text-[var(--text-secondary)] leading-relaxed">{String(value)}</p>
            </div>
        );
    };

    if (format === 'quick_summary' && msg.summary) {
        return (
            <div className="space-y-1">
                <p className="font-bold text-base text-[var(--text-primary)]">{msg.summary.title}</p>
                <Field label="Core Concept" value={msg.summary.core_concept_explanation} />
                <Field label="Key Points" value={msg.summary.key_points} />
                <Field label="Real-World Example" value={msg.summary.real_world_example} />
                <Field label="Conclusion" value={msg.summary.short_conclusion} />
            </div>
        );
    }
    if (format === 'detailed_explanation' && msg.detailed) {
        return (
            <div className="space-y-1">
                <p className="font-bold text-base text-[var(--text-primary)]">{msg.detailed.title}</p>
                <Field label="Explanation" value={msg.detailed.concept_explanation} />
                <Field label="Step-by-Step" value={msg.detailed.step_by_step_breakdown} />
                <Field label="Example" value={msg.detailed.example} />
                <Field label="Mini Quiz" value={msg.detailed.mini_quiz} />
            </div>
        );
    }
    if (format === 'anime_style' && msg.anime) {
        return (
            <div className="space-y-1">
                <p className="font-bold text-base text-fuchsia-300">🎌 {msg.anime.episode_title}</p>
                <Field label="Protagonist" value={msg.anime.main_character} />
                <Field label="Story Arc" value={msg.anime.story_arc} />
                <Field label="Physics Explained" value={msg.anime.physics_explanation} />
            </div>
        );
    }
    if (format === 'sports_analogy' && msg.sports) {
        return (
            <div className="space-y-1">
                <p className="font-bold text-base text-emerald-300">⚽ {msg.sports.game_title}</p>
                <Field label="Sport" value={msg.sports.sport_used} />
                <Field label="Play Breakdown" value={msg.sports.play_breakdown} />
                <Field label="Coaching Tip" value={msg.sports.coaching_tip} />
                <Field label="Scoreboard Stat" value={msg.sports.scoreboard_summary} />
            </div>
        );
    }
    if (format === 'academic_formal' && msg.academic) {
        return (
            <div className="space-y-1">
                <p className="font-bold text-base text-amber-300">🎓 {msg.academic.title}</p>
                <Field label="Abstract" value={msg.academic.abstract} />
                <Field label="Theoretical Background" value={msg.academic.theoretical_background} />
                <Field label="Methodology" value={msg.academic.methodology} />
                <Field label="Conclusion" value={msg.academic.conclusion} />
            </div>
        );
    }

    // Fallback: raw format if the struct is somehow not available
    return <span className="whitespace-pre-wrap">{msg.text}</span>;
}

// Illustration with fade-in and zoom-on-hover
function IllustrationImage({ url }: { url: string }) {
    return (
        <motion.div
            initial={{ opacity: 0, y: 8, scale: 0.98 }}
            animate={{ opacity: 1, y: 0, scale: 1 }}
            transition={{ duration: 0.5, delay: 0.2 }}
            className="mt-3 overflow-hidden rounded-2xl border border-violet-500/25"
            style={{ maxWidth: '400px' }}
        >
            <motion.img
                src={url}
                alt="AI generated illustration"
                loading="lazy"
                whileHover={{ scale: 1.04 }}
                transition={{ duration: 0.3 }}
                className="w-full object-cover block"
                onError={(e) => { (e.target as HTMLImageElement).style.display = 'none'; }}
            />
            <div className="px-3 py-1.5 bg-violet-500/10 flex items-center gap-1.5">
                <span className="text-[10px] text-violet-400">✨ AI-generated illustration</span>
            </div>
        </motion.div>
    );
}

function ChatBubble({ msg }: { msg: Message }) {
    const isAi = msg.role === 'ai';
    const isSystem = msg.role === 'system';

    if (isSystem) {
        return (
            <motion.div
                initial={{ opacity: 0, scale: 0.95 }}
                animate={{ opacity: 1, scale: 1 }}
                className="flex justify-center"
            >
                <div className="max-w-[85%] px-4 py-3 rounded-2xl bg-gradient-to-br from-violet-500/10 to-fuchsia-500/10 border border-violet-500/25 text-sm text-[var(--text-secondary)] text-center whitespace-pre-wrap leading-relaxed">
                    {msg.text}
                </div>
            </motion.div>
        );
    }

    const showStructured = isAi && msg.isStructuredJson && msg.outputFormat && !msg.isError;

    return (
        <div className={`flex gap-3 ${!isAi ? 'flex-row-reverse' : ''}`}>
            {isAi && (
                <div className="w-7 h-7 rounded-full bg-gradient-to-br from-violet-500 to-fuchsia-500 flex items-center justify-center flex-shrink-0 mt-0.5 shadow-md shadow-violet-500/30">
                    <svg className="w-3.5 h-3.5 text-white" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                        <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2.2} d="M13 10V3L4 14h7v7l9-11h-7z" />
                    </svg>
                </div>
            )}

            <div className="max-w-[82%] space-y-2">
                <div
                    className={`px-4 py-3 rounded-2xl text-sm leading-relaxed ${msg.isError
                        ? 'bg-red-500/10 border border-red-500/30 text-red-300 rounded-tl-sm'
                        : isAi
                            ? 'bg-[var(--bg-elevated)] border border-[var(--border)] text-[var(--text-primary)] rounded-tl-sm'
                            : 'bg-gradient-to-br from-violet-600 to-fuchsia-600 text-white rounded-tr-sm shadow-lg shadow-violet-500/20'
                        }`}
                >
                    {showStructured
                        ? <StructuredAnswer msg={msg} />
                        : <span className="whitespace-pre-wrap">{msg.text}</span>
                    }
                </div>

                {/* Illustration image */}
                {msg.illustrationUrl && <IllustrationImage url={msg.illustrationUrl} />}

                {/* Metadata badges */}
                {isAi && !msg.isError && (
                    <div className="flex items-center gap-2 flex-wrap">
                        {msg.contextUsed && (
                            <span className="text-[10px] px-2 py-0.5 rounded-full bg-cyan-500/10 border border-cyan-500/25 text-cyan-400">
                                RAG Context Used
                            </span>
                        )}
                        {msg.outputFormat && (
                            <span className="text-[10px] px-2 py-0.5 rounded-full bg-violet-500/10 border border-violet-500/25 text-violet-400">
                                {OUTPUT_FORMATS.find(f => f.id === msg.outputFormat)?.emoji}{' '}
                                {OUTPUT_FORMATS.find(f => f.id === msg.outputFormat)?.label}
                            </span>
                        )}
                        {msg.isStructuredJson && (
                            <span className="text-[10px] px-2 py-0.5 rounded-full bg-emerald-500/10 border border-emerald-500/25 text-emerald-400">
                                ✓ Structured
                            </span>
                        )}
                        {msg.tokensUsed && (
                            <span className="text-[10px] text-[var(--text-muted)]">{msg.tokensUsed} tokens · {msg.latencyMs}ms</span>
                        )}
                    </div>
                )}

                {/* Download button */}
                {msg.downloadURL && (
                    <motion.a
                        href={buildDownloadURL(msg.downloadURL)}
                        target="_blank"
                        rel="noopener noreferrer"
                        whileHover={{ scale: 1.02 }}
                        whileTap={{ scale: 0.97 }}
                        className="inline-flex items-center gap-2 px-3 py-1.5 rounded-lg bg-violet-500/15 border border-violet-500/30 text-xs text-violet-300 hover:bg-violet-500/25 transition-all"
                    >
                        <svg className="w-3 h-3" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                            <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M4 16v1a3 3 0 003 3h10a3 3 0 003-3v-1m-4-4l-4 4m0 0l-4-4m4 4V4" />
                        </svg>
                        {msg.downloadLabel || 'Download'}
                    </motion.a>
                )}

                {/* Audio player */}
                {msg.audioURL && (
                    <audio controls className="mt-2 w-full max-w-xs rounded-lg" style={{ height: '36px' }}>
                        <source src={buildDownloadURL(msg.audioURL)} type="audio/mpeg" />
                    </audio>
                )}
            </div>
        </div>
    );
}

function OutputFormatPicker({
    selected,
    onSelect,
}: {
    selected: OutputFormat | null;
    onSelect: (f: OutputFormat) => void;
}) {
    const groups: { label: string; items: OutputFormatOption[] }[] = [
        {
            label: 'Text',
            items: OUTPUT_FORMATS.filter(f =>
                ['quick_summary', 'detailed_explanation', 'academic_formal'].includes(f.id)
            ),
        },
        {
            label: 'Style',
            items: OUTPUT_FORMATS.filter(f =>
                ['anime_style', 'sports_analogy'].includes(f.id)
            ),
        },
        {
            label: 'Media',
            items: OUTPUT_FORMATS.filter(f =>
                ['presentation_slides', 'audio_explanation', 'translation'].includes(f.id)
            ),
        },
    ];

    return (
        <div className="space-y-3">
            {groups.map(group => (
                <div key={group.label}>
                    <p className="text-[10px] font-semibold tracking-widest uppercase text-[var(--text-muted)] mb-1.5">
                        {group.label}
                    </p>
                    <div className="space-y-1">
                        {group.items.map(fmt => (
                            <motion.button
                                key={fmt.id}
                                id={`format-${fmt.id}`}
                                whileHover={{ x: 3 }}
                                whileTap={{ scale: 0.97 }}
                                onClick={() => onSelect(fmt.id)}
                                className={`w-full flex items-center gap-2.5 px-3 py-2 rounded-xl border text-left transition-all text-xs ${selected === fmt.id
                                    ? 'bg-violet-500/20 border-violet-500/50 text-violet-300'
                                    : 'bg-[var(--bg-overlay)] border-[var(--border)] text-[var(--text-secondary)] hover:border-violet-500/30 hover:bg-[var(--bg-hover)]'
                                    }`}
                            >
                                <span className="text-base flex-shrink-0">{fmt.emoji}</span>
                                <div>
                                    <p className="font-medium">{fmt.label}</p>
                                    <p className="text-[10px] text-[var(--text-muted)]">{fmt.description}</p>
                                </div>
                                {selected === fmt.id && (
                                    <motion.div
                                        layoutId="format-check"
                                        className="ml-auto w-4 h-4 rounded-full bg-violet-500 flex items-center justify-center"
                                    >
                                        <svg className="w-2.5 h-2.5 text-white" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                                            <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={3} d="M5 13l4 4L19 7" />
                                        </svg>
                                    </motion.div>
                                )}
                            </motion.button>
                        ))}
                    </div>
                </div>
            ))}
        </div>
    );
}

function SourceCard({
    name,
    onRemove,
}: {
    name: string;
    onRemove?: () => void;
}) {
    const ext = name.split('.').pop()?.toUpperCase() ?? 'FILE';
    const colorMap: Record<string, string> = {
        PDF: 'bg-red-500/15 text-red-400',
        DOCX: 'bg-blue-500/15 text-blue-400',
        DOC: 'bg-blue-500/15 text-blue-400',
        TXT: 'bg-green-500/15 text-green-400',
        MD: 'bg-green-500/15 text-green-400',
        JPG: 'bg-amber-500/15 text-amber-400',
        JPEG: 'bg-amber-500/15 text-amber-400',
        PNG: 'bg-amber-500/15 text-amber-400',
        WEBP: 'bg-amber-500/15 text-amber-400',
    };
    const color = colorMap[ext] ?? 'bg-violet-500/15 text-violet-400';

    return (
        <motion.div
            layout
            initial={{ opacity: 0, x: -10 }}
            animate={{ opacity: 1, x: 0 }}
            exit={{ opacity: 0, x: -10 }}
            whileHover={{ y: -1 }}
            className="p-2.5 rounded-xl border border-[var(--border)] bg-[var(--bg-overlay)] group"
        >
            <div className="flex items-center gap-2.5">
                <div className={`w-8 h-8 rounded-lg flex items-center justify-center text-[10px] font-bold flex-shrink-0 ${color}`}>
                    {ext.slice(0, 4)}
                </div>
                <p className="text-xs font-medium text-[var(--text-primary)] truncate flex-1 min-w-0">
                    {name}
                </p>
                <div className="w-1.5 h-1.5 rounded-full bg-emerald-400 flex-shrink-0" title="Indexed" />
            </div>
        </motion.div>
    );
}

// ─────────────────────────────────────────────────────────────────────────────
// Main Component
// ─────────────────────────────────────────────────────────────────────────────

export function LearnNowView() {
    const [messages, setMessages] = useState<Message[]>([
        {
            id: 'welcome',
            role: 'ai',
            text: "Hello! I'm your adaptive AI tutor. Upload documents on the left or ask me anything. Want me to explain something or generate slides, audio, or a translation?",
            timestamp: new Date(),
        },
    ]);

    const [input, setInput] = useState('');
    const [isLoading, setIsLoading] = useState(false);
    const [selectedFormat, setSelectedFormat] = useState<OutputFormat | null>(null);
    const [sources, setSources] = useState<string[]>([]);
    const [isUploading, setIsUploading] = useState(false);
    const [dragOver, setDragOver] = useState(false);
    const [translationLang, setTranslationLang] = useState('Indonesian');
    const [showLangInput, setShowLangInput] = useState(false);

    const messagesEndRef = useRef<HTMLDivElement>(null);
    const fileInputRef = useRef<HTMLInputElement>(null);
    const textareaRef = useRef<HTMLTextAreaElement>(null);

    // Auto-scroll to bottom
    useEffect(() => {
        messagesEndRef.current?.scrollIntoView({ behavior: 'smooth' });
    }, [messages, isLoading]);

    // Load sources on mount
    useEffect(() => {
        getSources()
            .then(setSources)
            .catch(() => {/* not fatal */ });
    }, []);

    const addMessage = useCallback((msg: Omit<Message, 'id' | 'timestamp'>) => {
        setMessages(prev => [
            ...prev,
            { ...msg, id: crypto.randomUUID(), timestamp: new Date() },
        ]);
    }, []);

    // ── Send message ────────────────────────────────────────────────────────────
    const handleSend = useCallback(async () => {
        const question = input.trim();
        if (!question || isLoading) return;

        setInput('');
        addMessage({ role: 'user', text: question });
        setIsLoading(true);

        try {
            // Always send with a format — default to quick_summary if none selected
            // This prevents the backend "needs_format" loop from ever triggering
            const format = selectedFormat ?? 'quick_summary';

            const req = {
                question,
                output_format: format,
                target_language: format === 'translation' ? translationLang : undefined,
            };

            const resp: AskResponse = await askAI(req);

            // needs_format guard kept as safety net only (should never trigger now)
            if (resp.needs_format) {
                addMessage({
                    role: 'system',
                    text: '⚠️ Please select an output format from the right panel first.',
                });
                setIsLoading(false);
                return;
            }

            // If RAG was active but returned 0 chunks, notify user in the UI
            let prefix = '';
            if (sources.length > 0 && resp.context_found === 0) {
                prefix = '*(⚠️ RAG Active but no matching context found. Proceeding with general knowledge.)*\n\n';
            }

            addMessage({
                role: 'ai',
                text: prefix + resp.answer,
                outputFormat: resp.output_format,
                tokensUsed: resp.tokens_used,
                latencyMs: resp.latency_ms,
                contextUsed: resp.context_used,
                isStructuredJson: resp.is_structured_json,
                illustrationUrl: resp.illustration_url,
                summary: resp.summary,
                detailed: resp.detailed,
                anime: resp.anime,
                sports: resp.sports,
                academic: resp.academic,
                downloadURL: resp.output_format === 'presentation_slides' ? resp.download_url : undefined,
                downloadLabel: resp.output_format === 'presentation_slides' ? 'Download Presentation Slides' : undefined,
                audioURL: resp.output_format === 'audio_explanation' ? resp.download_url : undefined,
            });

        } catch (err: unknown) {
            const msg = err instanceof Error ? err.message : 'Something went wrong. Please try again.';
            addMessage({ role: 'ai', text: `❌ ${msg}`, isError: true });
        } finally {
            setIsLoading(false);
        }
    }, [input, isLoading, selectedFormat, translationLang, addMessage, sources.length]);


    const handleKeyDown = (e: React.KeyboardEvent) => {
        if (e.key === 'Enter' && !e.shiftKey) {
            e.preventDefault();
            handleSend();
        }
    };

    // ── File upload ─────────────────────────────────────────────────────────────
    const handleFileUpload = useCallback(async (files: FileList | null) => {
        if (!files || files.length === 0) return;
        const file = files[0];

        setIsUploading(true);
        addMessage({ role: 'system', text: `⏳ Uploading & indexing "${file.name}"…` });

        try {
            const result = await uploadFile(file);
            addMessage({
                role: 'system',
                text: `✅ "${result.source}" indexed!\n${result.word_count.toLocaleString()} words → ${result.chunks} chunks ready for context.`,
            });
            setSources(prev => Array.from(new Set([...prev, result.source])));
        } catch (err: unknown) {
            const msg = err instanceof Error ? err.message : 'Upload failed';
            addMessage({ role: 'system', text: `❌ Upload failed: ${msg}` });
        } finally {
            setIsUploading(false);
        }
    }, [addMessage]);

    const handleDrop = (e: React.DragEvent) => {
        e.preventDefault();
        setDragOver(false);
        handleFileUpload(e.dataTransfer.files);
    };

    // ─────────────────────────────────────────────────────────────────────────
    // Render
    // ─────────────────────────────────────────────────────────────────────────

    return (
        <div className="flex h-full overflow-hidden flex-col md:flex-row">

            {/* ── LEFT PANEL: Sources ──────────────────────────────────────────────── */}
            <div className="hidden md:flex w-60 flex-shrink-0 border-r border-[var(--border)] flex-col">
                <div className="p-4 border-b border-[var(--border)]">
                    <div className="flex items-center justify-between mb-3">
                        <h3 className="text-xs font-semibold tracking-widest uppercase text-[var(--text-muted)]">Sources</h3>
                        <span className="text-[10px] text-[var(--text-muted)] bg-[var(--bg-overlay)] px-2 py-0.5 rounded-full border border-[var(--border)]">
                            {sources.length}
                        </span>
                    </div>

                    {/* Drop zone */}
                    <motion.div
                        onDragOver={(e) => { e.preventDefault(); setDragOver(true); }}
                        onDragLeave={() => setDragOver(false)}
                        onDrop={handleDrop}
                        animate={{ borderColor: dragOver ? 'rgba(124,58,237,0.7)' : 'rgba(124,58,237,0.2)' }}
                        className="relative w-full py-4 px-2 rounded-xl border-2 border-dashed border-violet-500/20 text-center cursor-pointer hover:border-violet-500/50 transition-all"
                        onClick={() => fileInputRef.current?.click()}
                    >
                        {isUploading ? (
                            <div className="flex flex-col items-center gap-1">
                                <motion.div
                                    animate={{ rotate: 360 }}
                                    transition={{ duration: 1, repeat: Infinity, ease: 'linear' }}
                                    className="w-5 h-5 rounded-full border-2 border-violet-500 border-t-transparent"
                                />
                                <span className="text-[10px] text-violet-400">Indexing…</span>
                            </div>
                        ) : (
                            <div className="flex flex-col items-center gap-1">
                                <svg className="w-5 h-5 text-violet-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                                    <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={1.5} d="M7 16a4 4 0 01-.88-7.903A5 5 0 1115.9 6L16 6a5 5 0 011 9.9M15 13l-3-3m0 0l-3 3m3-3v12" />
                                </svg>
                                <span className="text-[10px] text-[var(--text-muted)]">PDF · DOCX · TXT · Image</span>
                                <span className="text-[10px] text-violet-400 font-medium">Click or drop</span>
                            </div>
                        )}
                    </motion.div>

                    <input
                        ref={fileInputRef}
                        type="file"
                        accept=".pdf,.docx,.txt,.md,.jpg,.jpeg,.png,.webp"
                        className="hidden"
                        onChange={(e) => handleFileUpload(e.target.files)}
                        id="file-upload-input"
                    />
                </div>

                {/* Source list */}
                <div className="flex-1 overflow-y-auto p-3 space-y-2">
                    <AnimatePresence>
                        {sources.length === 0 ? (
                            <motion.p
                                initial={{ opacity: 0 }}
                                animate={{ opacity: 1 }}
                                className="text-[11px] text-[var(--text-muted)] text-center py-4"
                            >
                                No sources yet. Upload a document to start RAG-augmented learning.
                            </motion.p>
                        ) : (
                            sources.map((src) => (
                                <SourceCard key={src} name={src} />
                            ))
                        )}
                    </AnimatePresence>
                </div>

                {/* RAG Status */}
                {sources.length > 0 && (
                    <div className="p-3 border-t border-[var(--border)]">
                        <div className="p-3 rounded-xl bg-gradient-to-br from-violet-500/10 to-fuchsia-500/10 border border-violet-500/20">
                            <div className="flex items-center gap-2 mb-1">
                                <span className="w-1.5 h-1.5 rounded-full bg-emerald-400 animate-pulse" />
                                <span className="text-[10px] font-semibold tracking-widest uppercase text-emerald-400">RAG Active</span>
                            </div>
                            <p className="text-[10px] text-[var(--text-muted)] leading-relaxed">
                                {sources.length} source{sources.length > 1 ? 's' : ''} loaded. AI will prioritise your uploaded context.
                            </p>
                        </div>
                    </div>
                )}
            </div>

            {/* ── CENTER: Chat ─────────────────────────────────────────────────────── */}
            <div className="flex-1 flex flex-col min-w-0 h-full">
                {/* Selected format badge */}
                <div className="px-5 py-2.5 border-b border-[var(--border)] flex items-center gap-2 flex-wrap">
                    <span className="text-[10px] text-[var(--text-muted)] font-medium uppercase tracking-wider">Mode:</span>
                    {selectedFormat ? (
                        <div className="flex items-center gap-1.5">
                            <span className="text-[10px] font-semibold px-2.5 py-1 rounded-full bg-violet-500/20 text-violet-400 border border-violet-500/40">
                                {OUTPUT_FORMATS.find(f => f.id === selectedFormat)?.emoji}{' '}
                                {OUTPUT_FORMATS.find(f => f.id === selectedFormat)?.label}
                            </span>
                            <button
                                onClick={() => setSelectedFormat(null)}
                                className="text-[10px] text-[var(--text-muted)] hover:text-red-400 transition-colors"
                            >
                                ✕
                            </button>
                        </div>
                    ) : (
                        <span className="text-[10px] text-[var(--text-muted)] italic">
                            Not set — AI will ask you each time
                        </span>
                    )}

                    {selectedFormat === 'translation' && (
                        <div className="flex items-center gap-1.5 ml-auto">
                            {showLangInput ? (
                                <input
                                    autoFocus
                                    value={translationLang}
                                    onChange={(e) => setTranslationLang(e.target.value)}
                                    onBlur={() => setShowLangInput(false)}
                                    onKeyDown={(e) => { if (e.key === 'Enter') setShowLangInput(false); }}
                                    className="text-[11px] px-2 py-0.5 rounded bg-[var(--bg-overlay)] border border-violet-500/40 text-violet-300 outline-none w-28"
                                    placeholder="e.g. French"
                                    id="translation-lang-input"
                                />
                            ) : (
                                <button
                                    onClick={() => setShowLangInput(true)}
                                    className="text-[10px] px-2 py-0.5 rounded-full border border-cyan-500/30 bg-cyan-500/10 text-cyan-400 hover:bg-cyan-500/20 transition-all"
                                >
                                    → {translationLang}
                                </button>
                            )}
                        </div>
                    )}
                </div>

                {/* Messages area */}
                <div className="flex-1 overflow-y-auto p-5 space-y-5">
                    <AnimatePresence initial={false}>
                        {messages.map((msg) => (
                            <motion.div
                                key={msg.id}
                                initial={{ opacity: 0, y: 12 }}
                                animate={{ opacity: 1, y: 0 }}
                                exit={{ opacity: 0, y: -8 }}
                                transition={{ duration: 0.25 }}
                            >
                                <ChatBubble msg={msg} />
                            </motion.div>
                        ))}
                    </AnimatePresence>

                    {isLoading && (
                        <motion.div initial={{ opacity: 0 }} animate={{ opacity: 1 }}>
                            <TypingIndicator />
                        </motion.div>
                    )}

                    <div ref={messagesEndRef} />
                </div>

                {/* Input area */}
                <div className="p-4 border-t border-[var(--border)]">
                    <div className={`flex items-end gap-2.5 p-3 rounded-2xl border bg-[var(--bg-overlay)] transition-all ${isLoading ? 'border-[var(--border)] opacity-75' : 'border-[var(--border)] focus-within:border-violet-500/60'
                        }`}>
                        {/* Attachment button */}
                        <button
                            onClick={() => fileInputRef.current?.click()}
                            disabled={isUploading}
                            className="w-8 h-8 flex-shrink-0 rounded-lg flex items-center justify-center text-[var(--text-muted)] hover:text-violet-400 hover:bg-violet-500/10 transition-all disabled:opacity-40"
                            title="Upload document"
                            id="chat-attach-btn"
                        >
                            <svg className="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={1.8} d="M15.172 7l-6.586 6.586a2 2 0 102.828 2.828l6.414-6.586a4 4 0 00-5.656-5.656l-6.415 6.585a6 6 0 108.486 8.486L20.5 13" />
                            </svg>
                        </button>

                        <textarea
                            ref={textareaRef}
                            id="chat-input"
                            value={input}
                            onChange={(e) => {
                                setInput(e.target.value);
                                // Auto-resize
                                e.target.style.height = 'auto';
                                e.target.style.height = Math.min(e.target.scrollHeight, 128) + 'px';
                            }}
                            onKeyDown={handleKeyDown}
                            placeholder={isLoading ? 'AI is thinking…' : 'Ask anything about your sources…'}
                            disabled={isLoading}
                            rows={1}
                            className="flex-1 bg-transparent text-sm text-[var(--text-primary)] placeholder-[var(--text-muted)] resize-none focus:outline-none min-h-[24px] max-h-32 disabled:opacity-50"
                        />

                        <motion.button
                            id="chat-send-btn"
                            whileHover={!isLoading && input.trim() ? { scale: 1.06 } : {}}
                            whileTap={!isLoading && input.trim() ? { scale: 0.94 } : {}}
                            disabled={!input.trim() || isLoading}
                            onClick={handleSend}
                            className="w-9 h-9 rounded-xl bg-gradient-to-br from-violet-600 to-fuchsia-600 flex items-center justify-center text-white disabled:opacity-35 shadow-lg shadow-violet-500/25 flex-shrink-0 transition-all"
                        >
                            {isLoading ? (
                                <motion.div
                                    animate={{ rotate: 360 }}
                                    transition={{ duration: 1, repeat: Infinity, ease: 'linear' }}
                                    className="w-4 h-4 rounded-full border-2 border-white border-t-transparent"
                                />
                            ) : (
                                <svg className="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                                    <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M12 19l9 2-9-18-9 18 9-2zm0 0v-8" />
                                </svg>
                            )}
                        </motion.button>
                    </div>

                    <p className="text-[10px] text-[var(--text-muted)] mt-2 text-center">
                        Shift+Enter for new line · Enter to send · Ctrl+/ to focus
                    </p>
                </div>
            </div>

            {/* ── RIGHT PANEL: AI Studio ────────────────────────────────────────────── */}
            <div className="hidden lg:flex w-60 flex-shrink-0 border-l border-[var(--border)] flex-col">
                <div className="p-4 border-b border-[var(--border)]">
                    <div className="flex items-center gap-2">
                        <div className="w-2 h-2 rounded-full bg-violet-400 animate-pulse" />
                        <h3 className="text-xs font-semibold tracking-widest uppercase text-[var(--text-muted)]">Output Format</h3>
                    </div>
                </div>

                <div className="flex-1 overflow-y-auto p-3">
                    <OutputFormatPicker
                        selected={selectedFormat}
                        onSelect={(f) => {
                            setSelectedFormat(prev => prev === f ? null : f);
                            if (f === 'translation') setShowLangInput(true);
                        }}
                    />
                </div>

                {/* AI Badge */}
                <div className="p-3 border-t border-[var(--border)]">
                    <div className="p-3 rounded-xl relative overflow-hidden" style={{
                        background: 'linear-gradient(135deg, rgba(124,58,237,0.12), rgba(236,72,153,0.08))',
                        border: '1px solid rgba(124,58,237,0.25)',
                    }}>
                        <p className="text-[10px] font-semibold uppercase tracking-widest text-violet-400 mb-1">AI Engine</p>
                        <p className="text-xs font-bold text-[var(--text-primary)]">Qwen Max</p>
                        <p className="text-[10px] text-[var(--text-muted)] mt-0.5">Alibaba DashScope</p>
                        <div className="flex items-center gap-1.5 mt-2">
                            <span className="w-1.5 h-1.5 rounded-full bg-emerald-400 animate-pulse" />
                            <span className="text-[10px] text-emerald-400">Connected</span>
                        </div>
                    </div>
                </div>
            </div>
        </div>
    );
}
