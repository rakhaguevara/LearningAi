'use client';

import { useState, useRef, useEffect } from 'react';
import { motion, AnimatePresence } from 'framer-motion';
import { useTheme } from '@/lib/ThemeContext';
import type { DashboardPage } from './Sidebar';

const PAGE_TITLES: Record<DashboardPage, string> = {
    learn: 'LearnNow',
    archive: 'Arsip Dokumen',
    pomodoro: 'Podomoro Time',
    settings: 'Settings',
    profile: 'My Profile',
};

interface TopbarProps {
    page: DashboardPage;
    onPageChange: (p: DashboardPage) => void;
}

export function Topbar({ page, onPageChange }: TopbarProps) {
    const { theme, toggleTheme } = useTheme();
    const [dropdownOpen, setDropdownOpen] = useState(false);
    const dropRef = useRef<HTMLDivElement>(null);

    // Close dropdown on outside click
    useEffect(() => {
        function handler(e: MouseEvent) {
            if (dropRef.current && !dropRef.current.contains(e.target as Node)) {
                setDropdownOpen(false);
            }
        }
        document.addEventListener('mousedown', handler);
        return () => document.removeEventListener('mousedown', handler);
    }, []);

    return (
        <header className="topbar-surface h-14 px-5 flex items-center justify-between flex-shrink-0 sticky top-0 z-30">
            {/* Page title */}
            <div className="flex items-center gap-3">
                <h1 className="text-base font-semibold text-[var(--text-primary)]">
                    {PAGE_TITLES[page]}
                </h1>
            </div>

            {/* Right controls */}
            <div className="flex items-center gap-2">
                {/* Theme toggle */}
                <motion.button
                    id="theme-toggle"
                    onClick={toggleTheme}
                    whileTap={{ scale: 0.9 }}
                    className="w-9 h-9 flex items-center justify-center rounded-xl border border-[var(--border)] bg-[var(--bg-overlay)] hover:bg-[var(--bg-hover)] text-[var(--text-secondary)] hover:text-[var(--text-primary)] transition-all"
                    aria-label="Toggle theme"
                >
                    <AnimatePresence mode="wait" initial={false}>
                        {theme === 'dark' ? (
                            <motion.svg
                                key="sun"
                                initial={{ rotate: -90, opacity: 0 }}
                                animate={{ rotate: 0, opacity: 1 }}
                                exit={{ rotate: 90, opacity: 0 }}
                                transition={{ duration: 0.2 }}
                                className="w-4 h-4"
                                fill="none" stroke="currentColor" viewBox="0 0 24 24"
                            >
                                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={1.8}
                                    d="M12 3v1m0 16v1m9-9h-1M4 12H3m15.364-6.364l-.707.707M6.343 17.657l-.707.707M17.657 17.657l-.707-.707M6.343 6.343l-.707-.707M12 9a3 3 0 100 6 3 3 0 000-6z" />
                            </motion.svg>
                        ) : (
                            <motion.svg
                                key="moon"
                                initial={{ rotate: 90, opacity: 0 }}
                                animate={{ rotate: 0, opacity: 1 }}
                                exit={{ rotate: -90, opacity: 0 }}
                                transition={{ duration: 0.2 }}
                                className="w-4 h-4"
                                fill="none" stroke="currentColor" viewBox="0 0 24 24"
                            >
                                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={1.8}
                                    d="M20.354 15.354A9 9 0 018.646 3.646 9.003 9.003 0 0012 21a9.003 9.003 0 008.354-5.646z" />
                            </motion.svg>
                        )}
                    </AnimatePresence>
                </motion.button>

                {/* Notification bell */}
                <motion.button
                    whileTap={{ scale: 0.9 }}
                    className="w-9 h-9 flex items-center justify-center rounded-xl border border-[var(--border)] bg-[var(--bg-overlay)] hover:bg-[var(--bg-hover)] text-[var(--text-secondary)] hover:text-[var(--text-primary)] transition-all relative"
                >
                    <svg className="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                        <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={1.8}
                            d="M15 17h5l-1.405-1.405A2.032 2.032 0 0118 14.158V11a6.002 6.002 0 00-4-5.659V5a2 2 0 10-4 0v.341C7.67 6.165 6 8.388 6 11v3.159c0 .538-.214 1.055-.595 1.436L4 17h5m6 0v1a3 3 0 11-6 0v-1m6 0H9" />
                    </svg>
                    <span className="absolute top-2 right-2 w-1.5 h-1.5 rounded-full bg-violet-400" />
                </motion.button>

                {/* Avatar dropdown */}
                <div className="relative" ref={dropRef}>
                    <motion.button
                        id="user-avatar-btn"
                        onClick={() => setDropdownOpen((v) => !v)}
                        whileTap={{ scale: 0.95 }}
                        className="w-9 h-9 rounded-xl bg-gradient-to-br from-violet-500 to-fuchsia-500 flex items-center justify-center text-white text-xs font-bold shadow-lg shadow-violet-500/25 hover:shadow-violet-500/40 transition-all"
                    >
                        N
                    </motion.button>

                    <AnimatePresence>
                        {dropdownOpen && (
                            <motion.div
                                initial={{ opacity: 0, scale: 0.95, y: -6 }}
                                animate={{ opacity: 1, scale: 1, y: 0 }}
                                exit={{ opacity: 0, scale: 0.95, y: -6 }}
                                transition={{ duration: 0.15 }}
                                className="absolute right-0 top-11 w-56 rounded-2xl shadow-2xl border border-[var(--border)] overflow-hidden z-50"
                                style={{ background: 'var(--bg-elevated)' }}
                            >
                                <div className="px-4 py-3 border-b border-[var(--border)]">
                                    <p className="text-sm font-semibold text-[var(--text-primary)]">Neura User</p>
                                    <p className="text-xs text-[var(--text-muted)]">user@example.com</p>
                                </div>
                                {[
                                    { label: 'Profile', page: 'profile' as DashboardPage },
                                    { label: 'Settings', page: 'settings' as DashboardPage },
                                ].map((item) => (
                                    <button
                                        key={item.page}
                                        onClick={() => { onPageChange(item.page); setDropdownOpen(false); }}
                                        className="w-full text-left px-4 py-2.5 text-sm text-[var(--text-secondary)] hover:text-[var(--text-primary)] hover:bg-[var(--bg-hover)] transition-colors"
                                    >
                                        {item.label}
                                    </button>
                                ))}
                                <div className="border-t border-[var(--border)]">
                                    <button
                                        onClick={() => { localStorage.removeItem('access_token'); window.location.href = '/login'; }}
                                        className="w-full text-left px-4 py-2.5 text-sm text-red-400 hover:bg-red-400/10 transition-colors"
                                    >
                                        Sign out
                                    </button>
                                </div>
                            </motion.div>
                        )}
                    </AnimatePresence>
                </div>
            </div>
        </header>
    );
}
