'use client';

import { motion } from 'framer-motion';
import { useTheme } from '@/lib/ThemeContext';
import { useLanguage } from '@/lib/LanguageContext';

export function SettingsView() {
    const { theme, toggleTheme } = useTheme();
    const { language, setLanguage, t } = useLanguage();

    return (
        <div className="h-full overflow-y-auto p-6">
            <div className="max-w-2xl mx-auto space-y-5">

                {/* Appearance */}
                <motion.section
                    initial={{ opacity: 0, y: 14 }}
                    animate={{ opacity: 1, y: 0 }}
                    className="p-5 rounded-2xl border border-[var(--border)] bg-[var(--bg-surface)]"
                >
                    <h3 className="text-xs font-semibold tracking-widest uppercase text-[var(--text-muted)] mb-4">{t('appearance.title')}</h3>

                    {/* Theme toggle */}
                    <div className="flex flex-col sm:flex-row sm:items-center justify-between py-3 border-b border-[var(--border)] gap-3">
                        <div>
                            <p className="text-sm font-medium text-[var(--text-primary)]">{t('theme.label')}</p>
                            <p className="text-xs text-[var(--text-muted)] mt-0.5">{t('theme.description')}</p>
                        </div>
                        <div className="flex gap-2 p-1 rounded-xl bg-[var(--bg-overlay)] border border-[var(--border)] self-start sm:self-auto">
                            {(['dark', 'light'] as const).map((tTheme) => (
                                <button
                                    key={tTheme}
                                    onClick={() => { if (theme !== tTheme) toggleTheme(); }}
                                    className={`px-3 sm:px-4 py-1.5 rounded-lg text-xs font-medium transition-all capitalize ${theme === tTheme
                                            ? 'bg-gradient-to-r from-violet-600 to-fuchsia-600 text-white shadow-md'
                                            : 'text-[var(--text-muted)] hover:text-[var(--text-primary)]'
                                        }`}
                                >
                                    {tTheme === 'dark' ? t('theme.dark') : t('theme.light')}
                                </button>
                            ))}
                        </div>
                    </div>

                    {/* Language */}
                    <div className="flex flex-col sm:flex-row sm:items-center justify-between py-3 gap-3">
                        <div>
                            <p className="text-sm font-medium text-[var(--text-primary)]">{t('language.label')}</p>
                            <p className="text-xs text-[var(--text-muted)] mt-0.5">{t('language.description')}</p>
                        </div>
                        <select 
                            value={language}
                            onChange={(e) => setLanguage(e.target.value as 'en' | 'id')}
                            className="text-sm text-[var(--text-primary)] bg-[var(--bg-overlay)] border border-[var(--border)] rounded-xl px-3 py-1.5 focus:outline-none focus:border-violet-500/50 transition-all self-start sm:self-auto"
                        >
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
                    <h3 className="text-xs font-semibold tracking-widest uppercase text-[var(--text-muted)] mb-4">{t('aiPreferences.title')}</h3>
                    {[
                        { labelKey: 'ai.autoSummarize', descKey: 'ai.autoSummarize.desc' },
                        { labelKey: 'ai.adaptiveTone', descKey: 'ai.adaptiveTone.desc' },
                        { labelKey: 'ai.streamResponses', descKey: 'ai.streamResponses.desc' },
                    ].map((item) => (
                        <div key={item.labelKey} className="flex flex-col sm:flex-row sm:items-center justify-between py-3 border-b border-[var(--border)] last:border-none gap-3">
                            <div>
                                <p className="text-sm font-medium text-[var(--text-primary)]">{t(item.labelKey)}</p>
                                <p className="text-xs text-[var(--text-muted)] mt-0.5">{t(item.descKey)}</p>
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
                    <h3 className="text-xs font-semibold tracking-widest uppercase text-[var(--text-muted)] mb-4">{t('pomodoro.title')}</h3>
                    {[
                        { labelKey: 'pomodoro.soundAlerts', descKey: 'pomodoro.soundAlerts.desc', defaultOn: true },
                        { labelKey: 'pomodoro.autoStart', descKey: 'pomodoro.autoStart.desc', defaultOn: false },
                        { labelKey: 'pomodoro.notifications', descKey: 'pomodoro.notifications.desc', defaultOn: true },
                    ].map((item) => (
                        <div key={item.labelKey} className="flex flex-col sm:flex-row sm:items-center justify-between py-3 border-b border-[var(--border)] last:border-none gap-3">
                            <div>
                                <p className="text-sm font-medium text-[var(--text-primary)]">{t(item.labelKey)}</p>
                                <p className="text-xs text-[var(--text-muted)] mt-0.5">{t(item.descKey)}</p>
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
                    <h3 className="text-xs font-semibold tracking-widest uppercase text-red-400/70 mb-4">{t('dangerZone.title')}</h3>
                    <div className="flex flex-col sm:flex-row sm:items-center justify-between gap-3">
                        <div>
                            <p className="text-sm font-medium text-[var(--text-primary)]">{t('dangerZone.deleteAccount')}</p>
                            <p className="text-xs text-[var(--text-muted)] mt-0.5">{t('dangerZone.deleteAccount.desc')}</p>
                        </div>
                        <button className="px-4 py-2 rounded-xl border border-red-500/30 text-red-400 text-sm hover:bg-red-500/10 transition-all font-medium self-start sm:self-auto">
                            {t('dangerZone.delete')}
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
