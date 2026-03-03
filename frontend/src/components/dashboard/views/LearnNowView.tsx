'use client';

import { useState } from 'react';
import { motion } from 'framer-motion';

const SAMPLE_SOURCES = [
    { id: 1, title: 'Neural Networks Basics', type: 'PDF', pages: 24, added: '2h ago' },
    { id: 2, title: 'Python ML Course Notes', type: 'TXT', pages: 8, added: '1d ago' },
    { id: 3, title: 'Deep Learning Paper', type: 'PDF', pages: 42, added: '3d ago' },
];

const SAMPLE_MESSAGES = [
    { role: 'ai', text: 'Hello! I\'ve analyzed your sources. What would you like to understand today?' },
    { role: 'user', text: 'Explain backpropagation like I\'m new to ML' },
    { role: 'ai', text: 'Great question! Think of backpropagation like a coach reviewing a game film and giving each player feedback on what went wrong, working backwards from the final score. In neural networks, we calculate the error at the output, then distribute "blame" backwards through each layer, adjusting weights so the network improves next time.\n\nWant me to go deeper into the math, or would a visual example help more?' },
];

function SourceCard({ source }: { source: typeof SAMPLE_SOURCES[0] }) {
    return (
        <motion.div
            whileHover={{ y: -2, borderColor: 'rgba(124,58,237,0.4)' }}
            className="p-3 rounded-xl border border-[var(--border)] bg-[var(--bg-overlay)] cursor-pointer transition-all group"
        >
            <div className="flex items-start gap-3">
                <div className={`w-9 h-9 rounded-lg flex items-center justify-center text-xs font-bold flex-shrink-0
                    ${source.type === 'PDF' ? 'bg-red-500/15 text-red-400' : 'bg-blue-500/15 text-blue-400'}`}>
                    {source.type}
                </div>
                <div className="min-w-0">
                    <p className="text-xs font-medium text-[var(--text-primary)] truncate">{source.title}</p>
                    <p className="text-[10px] text-[var(--text-muted)] mt-0.5">{source.pages} pages · {source.added}</p>
                </div>
            </div>
        </motion.div>
    );
}

function ChatBubble({ msg }: { msg: typeof SAMPLE_MESSAGES[0] }) {
    const isAi = msg.role === 'ai';
    return (
        <div className={`flex gap-3 ${isAi ? '' : 'flex-row-reverse'}`}>
            {isAi && (
                <div className="w-7 h-7 rounded-full bg-gradient-to-br from-violet-500 to-fuchsia-500 flex items-center justify-center flex-shrink-0 mt-0.5 shadow-md shadow-violet-500/30">
                    <svg className="w-3.5 h-3.5 text-white" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                        <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2.2} d="M13 10V3L4 14h7v7l9-11h-7z" />
                    </svg>
                </div>
            )}
            <div className={`max-w-[80%] px-4 py-3 rounded-2xl text-sm leading-relaxed whitespace-pre-wrap
                ${isAi
                    ? 'bg-[var(--bg-elevated)] border border-[var(--border)] text-[var(--text-primary)] rounded-tl-sm'
                    : 'bg-gradient-to-br from-violet-600 to-fuchsia-600 text-white rounded-tr-sm shadow-lg shadow-violet-500/20'
                }`}>
                {msg.text}
            </div>
        </div>
    );
}

export function LearnNowView() {
    const [message, setMessage] = useState('');
    const [tone, setTone] = useState<'adaptive' | 'simple' | 'expert'>('adaptive');

    return (
        <div className="flex h-full overflow-hidden">

            {/* ── Left Panel: Sources ─────────────────────────── */}
            <div className="w-64 flex-shrink-0 border-r border-[var(--border)] flex flex-col">
                <div className="p-4 border-b border-[var(--border)]">
                    <div className="flex items-center justify-between mb-3">
                        <h3 className="text-xs font-semibold tracking-widest uppercase text-[var(--text-muted)]">Sources</h3>
                        <span className="text-[10px] text-[var(--text-muted)] bg-[var(--bg-overlay)] px-2 py-0.5 rounded-full border border-[var(--border)]">
                            {SAMPLE_SOURCES.length}
                        </span>
                    </div>
                    <motion.button
                        whileHover={{ scale: 1.02 }}
                        whileTap={{ scale: 0.97 }}
                        className="w-full py-2 px-3 rounded-xl border-2 border-dashed border-[var(--border)] text-xs text-[var(--text-muted)] hover:border-violet-500/50 hover:text-violet-400 transition-all flex items-center justify-center gap-2"
                    >
                        <svg className="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                            <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M12 4v16m8-8H4" />
                        </svg>
                        Add Source
                    </motion.button>
                </div>
                <div className="flex-1 overflow-y-auto p-3 space-y-2">
                    {SAMPLE_SOURCES.map((s) => <SourceCard key={s.id} source={s} />)}
                </div>

                {/* AI Summary Card */}
                <div className="p-3 border-t border-[var(--border)]">
                    <div className="p-3 rounded-xl bg-gradient-to-br from-violet-500/10 to-fuchsia-500/10 border border-violet-500/20">
                        <div className="flex items-center gap-2 mb-2">
                            <span className="w-1.5 h-1.5 rounded-full bg-violet-400 animate-pulse" />
                            <span className="text-[10px] font-semibold tracking-widest uppercase text-violet-400">AI Summary</span>
                        </div>
                        <p className="text-xs text-[var(--text-secondary)] leading-relaxed">
                            3 sources loaded covering ML fundamentals, neural nets, and deep learning. Ready to explain at any depth.
                        </p>
                    </div>
                </div>
            </div>

            {/* ── Center: Chat Workspace ──────────────────────── */}
            <div className="flex-1 flex flex-col min-w-0">
                {/* Tone badge */}
                <div className="px-6 py-3 border-b border-[var(--border)] flex items-center gap-2">
                    <span className="text-xs text-[var(--text-muted)]">Tone:</span>
                    {(['adaptive', 'simple', 'expert'] as const).map((t) => (
                        <button
                            key={t}
                            onClick={() => setTone(t)}
                            className={`text-[10px] font-semibold tracking-wider uppercase px-2.5 py-1 rounded-full transition-all ${tone === t
                                    ? 'bg-violet-500/20 text-violet-400 border border-violet-500/40'
                                    : 'text-[var(--text-muted)] hover:text-[var(--text-secondary)] border border-transparent hover:border-[var(--border)]'
                                }`}
                        >
                            {t}
                        </button>
                    ))}
                </div>

                {/* Messages */}
                <div className="flex-1 overflow-y-auto p-6 space-y-5">
                    {SAMPLE_MESSAGES.map((msg, i) => (
                        <motion.div
                            key={i}
                            initial={{ opacity: 0, y: 10 }}
                            animate={{ opacity: 1, y: 0 }}
                            transition={{ delay: i * 0.1 }}
                        >
                            <ChatBubble msg={msg} />
                        </motion.div>
                    ))}
                </div>

                {/* Input */}
                <div className="p-4 border-t border-[var(--border)]">
                    <div className="flex items-end gap-3 p-3 rounded-2xl border border-[var(--border)] bg-[var(--bg-overlay)] focus-within:border-violet-500/50 transition-all">
                        <textarea
                            value={message}
                            onChange={(e) => setMessage(e.target.value)}
                            placeholder="Ask anything about your sources..."
                            rows={1}
                            className="flex-1 bg-transparent text-sm text-[var(--text-primary)] placeholder-[var(--text-muted)] resize-none focus:outline-none min-h-[24px] max-h-32"
                            onKeyDown={(e) => { if (e.key === 'Enter' && !e.shiftKey) { e.preventDefault(); setMessage(''); } }}
                        />
                        <motion.button
                            whileHover={{ scale: 1.05 }}
                            whileTap={{ scale: 0.95 }}
                            disabled={!message.trim()}
                            className="w-9 h-9 rounded-xl bg-gradient-to-br from-violet-600 to-fuchsia-600 flex items-center justify-center text-white disabled:opacity-40 shadow-lg shadow-violet-500/25 flex-shrink-0"
                        >
                            <svg className="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M12 19l9 2-9-18-9 18 9-2zm0 0v-8" />
                            </svg>
                        </motion.button>
                    </div>
                    <p className="text-[10px] text-[var(--text-muted)] mt-2 text-center">
                        Shift+Enter for new line · Enter to send
                    </p>
                </div>
            </div>

            {/* ── Right Panel: Studio ─────────────────────────── */}
            <div className="w-64 flex-shrink-0 border-l border-[var(--border)] flex flex-col">
                <div className="p-4 border-b border-[var(--border)]">
                    <h3 className="text-xs font-semibold tracking-widest uppercase text-[var(--text-muted)]">Studio Tools</h3>
                </div>
                <div className="flex-1 overflow-y-auto p-3 space-y-2">
                    {[
                        { label: 'Generate Summary', icon: '✨', color: 'violet' },
                        { label: 'Mind Map', icon: '🗺️', color: 'cyan' },
                        { label: 'Flashcards', icon: '🃏', color: 'pink' },
                        { label: 'Quiz Me', icon: '🧠', color: 'amber' },
                        { label: 'Explain Differently', icon: '🔄', color: 'green' },
                        { label: 'Deep Dive', icon: '🔬', color: 'blue' },
                    ].map((tool) => (
                        <motion.button
                            key={tool.label}
                            whileHover={{ x: 3 }}
                            whileTap={{ scale: 0.97 }}
                            className="w-full flex items-center gap-3 px-3 py-2.5 rounded-xl border border-[var(--border)] bg-[var(--bg-overlay)] hover:bg-[var(--bg-hover)] hover:border-violet-500/30 transition-all text-left"
                        >
                            <span className="text-base">{tool.icon}</span>
                            <span className="text-xs font-medium text-[var(--text-secondary)] group-hover:text-[var(--text-primary)]">{tool.label}</span>
                        </motion.button>
                    ))}
                </div>

                {/* Gradient highlight border accent */}
                <div className="p-3 border-t border-[var(--border)]">
                    <div className="p-3 rounded-xl relative overflow-hidden" style={{
                        background: 'linear-gradient(135deg, rgba(124,58,237,0.12) 0%, rgba(236,72,153,0.08) 100%)',
                        border: '1px solid rgba(124,58,237,0.25)',
                    }}>
                        <p className="text-[10px] font-semibold uppercase tracking-widest text-violet-400 mb-1">Learning Streak</p>
                        <p className="text-2xl font-bold text-[var(--text-primary)]">7 <span className="text-sm font-normal text-[var(--text-muted)]">days</span></p>
                        <div className="flex gap-1 mt-2">
                            {Array.from({ length: 7 }).map((_, i) => (
                                <div key={i} className="flex-1 h-1.5 rounded-full bg-gradient-to-r from-violet-500 to-fuchsia-500 opacity-80" />
                            ))}
                        </div>
                    </div>
                </div>
            </div>
        </div>
    );
}
