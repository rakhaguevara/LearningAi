'use client';

import { motion } from 'framer-motion';
import { useTheme } from '@/lib/ThemeContext';

export function SettingsView() {
    const { theme, toggleTheme } = useTheme();

    return (
        <div className="h-full overflow-y-auto p-6">
            <div className="max-w-2xl mx-auto space-y-5">

                {/* Appearance */}
                <motion.section
                    initial={{ opacity: 0, y: 14 }}
                    animate={{ opacity: 1, y: 0 }}
                    className="p-5 rounded-2xl border border-[var(--border)] bg-[var(--bg-surface)]"
                >
                    <h3 className="text-xs font-semibold tracking-widest uppercase text-[var(--text-muted)] mb-4">Appearance</h3>

                    {/* Theme toggle */}
                    <div className="flex items-center justify-between py-3 border-b border-[var(--border)]">
                        <div>
                            <p className="text-sm font-medium text-[var(--text-primary)]">Theme</p>
                            <p className="text-xs text-[var(--text-muted)] mt-0.5">Switch between dark and light mode</p>
                        </div>
                        <div className="flex gap-2 p-1 rounded-xl bg-[var(--bg-overlay)] border border-[var(--border)]">
                            {(['dark', 'light'] as const).map((t) => (
                                <button
                                    key={t}
                                    onClick={() => { if (theme !== t) toggleTheme(); }}
                                    className={`px-4 py-1.5 rounded-lg text-xs font-medium transition-all capitalize ${theme === t
                                            ? 'bg-gradient-to-r from-violet-600 to-fuchsia-600 text-white shadow-md'
                                            : 'text-[var(--text-muted)] hover:text-[var(--text-primary)]'
                                        }`}
                                >
                                    {t === 'dark' ? '🌙 Dark' : '☀️ Light'}
                                </button>
                            ))}
                        </div>
                    </div>

                    {/* Language */}
                    <div className="flex items-center justify-between py-3">
                        <div>
                            <p className="text-sm font-medium text-[var(--text-primary)]">Language</p>
                            <p className="text-xs text-[var(--text-muted)] mt-0.5">Interface language</p>
                        </div>
                        <select className="text-sm text-[var(--text-primary)] bg-[var(--bg-overlay)] border border-[var(--border)] rounded-xl px-3 py-1.5 focus:outline-none focus:border-violet-500/50 transition-all">
                            <option value="en">🇺🇸 English</option>
                            <option value="id">🇮🇩 Bahasa Indonesia</option>
                        </select>
                    </div>
                </motion.section>

                {/* AI Preferences */}
                <motion.section
                    initial={{ opacity: 0, y: 14 }}
                    animate={{ opacity: 1, y: 0 }}
                    transition={{ delay: 0.06 }}
                    className="p-5 rounded-2xl border border-[var(--border)] bg-[var(--bg-surface)]"
                >
                    <h3 className="text-xs font-semibold tracking-widest uppercase text-[var(--text-muted)] mb-4">AI Preferences</h3>
                    {[
                        { label: 'Auto-summarize sources', description: 'Generate a summary when sources are added' },
                        { label: 'Adaptive tone', description: 'Adjust explanation style based on behavior signals' },
                        { label: 'Stream responses', description: 'Show AI responses as they\'re generated' },
                    ].map((item) => (
                        <div key={item.label} className="flex items-center justify-between py-3 border-b border-[var(--border)] last:border-none">
                            <div>
                                <p className="text-sm font-medium text-[var(--text-primary)]">{item.label}</p>
                                <p className="text-xs text-[var(--text-muted)] mt-0.5">{item.description}</p>
                            </div>
                            <Toggle defaultOn />
                        </div>
                    ))}
                </motion.section>

                {/* Pomodoro Preferences */}
                <motion.section
                    initial={{ opacity: 0, y: 14 }}
                    animate={{ opacity: 1, y: 0 }}
                    transition={{ delay: 0.1 }}
                    className="p-5 rounded-2xl border border-[var(--border)] bg-[var(--bg-surface)]"
                >
                    <h3 className="text-xs font-semibold tracking-widest uppercase text-[var(--text-muted)] mb-4">Podomoro</h3>
                    {[
                        { label: 'Sound alerts', description: 'Play a sound when timer ends', defaultOn: true },
                        { label: 'Auto-start breaks', description: 'Automatically start break timer', defaultOn: false },
                        { label: 'Notifications', description: 'Browser notifications for session end', defaultOn: true },
                    ].map((item) => (
                        <div key={item.label} className="flex items-center justify-between py-3 border-b border-[var(--border)] last:border-none">
                            <div>
                                <p className="text-sm font-medium text-[var(--text-primary)]">{item.label}</p>
                                <p className="text-xs text-[var(--text-muted)] mt-0.5">{item.description}</p>
                            </div>
                            <Toggle defaultOn={item.defaultOn} />
                        </div>
                    ))}
                </motion.section>

                {/* Danger zone */}
                <motion.section
                    initial={{ opacity: 0, y: 14 }}
                    animate={{ opacity: 1, y: 0 }}
                    transition={{ delay: 0.14 }}
                    className="p-5 rounded-2xl border border-red-500/20 bg-red-500/5"
                >
                    <h3 className="text-xs font-semibold tracking-widest uppercase text-red-400/70 mb-4">Danger Zone</h3>
                    <div className="flex items-center justify-between">
                        <div>
                            <p className="text-sm font-medium text-[var(--text-primary)]">Delete Account</p>
                            <p className="text-xs text-[var(--text-muted)] mt-0.5">Permanently remove all data</p>
                        </div>
                        <button className="px-4 py-2 rounded-xl border border-red-500/30 text-red-400 text-sm hover:bg-red-500/10 transition-all font-medium">
                            Delete
                        </button>
                    </div>
                </motion.section>
            </div>
        </div>
    );
}

function Toggle({ defaultOn = false }: { defaultOn?: boolean }) {
    const [on, setOn] = useState(defaultOn);

    return (
        <motion.button
            onClick={() => setOn((v) => !v)}
            className={`relative w-11 h-6 rounded-full transition-all duration-200 flex-shrink-0 ${on ? 'bg-gradient-to-r from-violet-600 to-fuchsia-600 shadow-md shadow-violet-500/25' : 'bg-[var(--border-strong)]'
                }`}
        >
            <motion.div
                className="absolute top-0.5 left-0.5 w-5 h-5 rounded-full bg-white shadow-sm"
                animate={{ x: on ? 20 : 0 }}
                transition={{ type: 'spring', stiffness: 500, damping: 30 }}
            />
        </motion.button>
    );
}

import { useState } from 'react';
