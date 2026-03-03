'use client';

import { motion } from 'framer-motion';
import { useTheme } from '@/lib/ThemeContext';

export type DashboardPage = 'learn' | 'archive' | 'pomodoro' | 'settings' | 'profile';

interface SidebarProps {
    active: DashboardPage;
    onChange: (page: DashboardPage) => void;
    userName?: string;
    userAvatar?: string;
}

const NAV_ITEMS: { id: DashboardPage; label: string; icon: React.ReactNode }[] = [
    {
        id: 'learn',
        label: 'LearnNow',
        icon: (
            <svg className="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={1.8}
                    d="M12 6.253v13m0-13C10.832 5.477 9.246 5 7.5 5S4.168 5.477 3 6.253v13C4.168 18.477 5.754 18 7.5 18s3.332.477 4.5 1.253m0-13C13.168 5.477 14.754 5 16.5 5c1.747 0 3.332.477 4.5 1.253v13C19.832 18.477 18.247 18 16.5 18c-1.746 0-3.332.477-4.5 1.253" />
            </svg>
        ),
    },
    {
        id: 'archive',
        label: 'Arsip Dokumen',
        icon: (
            <svg className="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={1.8}
                    d="M5 8h14M5 8a2 2 0 110-4h14a2 2 0 110 4M5 8v10a2 2 0 002 2h10a2 2 0 002-2V8m-9 4h4" />
            </svg>
        ),
    },
    {
        id: 'pomodoro',
        label: 'Podomoro Time',
        icon: (
            <svg className="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={1.8}
                    d="M12 8v4l3 3m6-3a9 9 0 11-18 0 9 9 0 0118 0z" />
            </svg>
        ),
    },
];

const BOTTOM_ITEMS: { id: DashboardPage; label: string; icon: React.ReactNode }[] = [
    {
        id: 'settings',
        label: 'Settings',
        icon: (
            <svg className="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={1.8}
                    d="M10.325 4.317c.426-1.756 2.924-1.756 3.35 0a1.724 1.724 0 002.573 1.066c1.543-.94 3.31.826 2.37 2.37a1.724 1.724 0 001.065 2.572c1.756.426 1.756 2.924 0 3.35a1.724 1.724 0 00-1.066 2.573c.94 1.543-.826 3.31-2.37 2.37a1.724 1.724 0 00-2.572 1.065c-.426 1.756-2.924 1.756-3.35 0a1.724 1.724 0 00-2.573-1.066c-1.543.94-3.31-.826-2.37-2.37a1.724 1.724 0 00-1.065-2.572c-1.756-.426-1.756-2.924 0-3.35a1.724 1.724 0 001.066-2.573c-.94-1.543.826-3.31 2.37-2.37.996.608 2.296.07 2.572-1.065z" />
                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={1.8} d="M15 12a3 3 0 11-6 0 3 3 0 016 0z" />
            </svg>
        ),
    },
    {
        id: 'profile',
        label: 'Profile',
        icon: (
            <svg className="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={1.8}
                    d="M16 7a4 4 0 11-8 0 4 4 0 018 0zM12 14a7 7 0 00-7 7h14a7 7 0 00-7-7z" />
            </svg>
        ),
    },
];

function NavItem({
    item,
    active,
    onClick,
}: {
    item: (typeof NAV_ITEMS)[0];
    active: boolean;
    onClick: () => void;
}) {
    return (
        <motion.button
            onClick={onClick}
            whileHover={{ x: 3 }}
            whileTap={{ scale: 0.97 }}
            className={`
                w-full flex items-center gap-3 px-3 py-2.5 rounded-xl text-sm font-medium
                transition-all duration-200 relative group
                ${active
                    ? 'text-white'
                    : 'text-[var(--text-secondary)] hover:text-[var(--text-primary)] hover:bg-[var(--bg-hover)]'
                }
            `}
        >
            {active && (
                <motion.div
                    layoutId="sidebar-active-pill"
                    className="absolute inset-0 rounded-xl bg-[var(--brand-soft)] border border-[var(--brand)]/30"
                    transition={{ type: 'spring', stiffness: 380, damping: 30 }}
                />
            )}
            <span className={`relative z-10 transition-colors ${active ? 'text-violet-400' : ''}`}>
                {item.icon}
            </span>
            <span className="relative z-10">{item.label}</span>
            {active && (
                <span className="relative z-10 ml-auto w-1.5 h-1.5 rounded-full bg-violet-400" />
            )}
        </motion.button>
    );
}

export function Sidebar({ active, onChange }: SidebarProps) {
    const { theme } = useTheme();

    return (
        <aside className="sidebar-surface flex flex-col h-screen w-[240px] flex-shrink-0 sticky top-0">
            {/* Logo */}
            <div className="px-5 py-5 border-b border-[var(--border)]">
                <div className="flex items-center gap-2.5">
                    <div className="w-8 h-8 rounded-xl bg-gradient-to-br from-violet-500 to-fuchsia-500 flex items-center justify-center shadow-lg shadow-violet-500/30">
                        <svg className="w-4 h-4 text-white" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                            <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2.2} d="M13 10V3L4 14h7v7l9-11h-7z" />
                        </svg>
                    </div>
                    <div>
                        <span className="text-sm font-bold text-[var(--text-primary)]">NeuraLearn</span>
                        <div className="text-[10px] text-[var(--text-muted)] font-medium tracking-wider uppercase">AI</div>
                    </div>
                </div>
            </div>

            {/* Main nav */}
            <nav className="flex-1 p-3 space-y-1 overflow-y-auto">
                <p className="text-[10px] font-semibold tracking-widest uppercase text-[var(--text-muted)] px-3 pt-2 pb-1">
                    Workspace
                </p>
                {NAV_ITEMS.map((item) => (
                    <NavItem
                        key={item.id}
                        item={item}
                        active={active === item.id}
                        onClick={() => onChange(item.id)}
                    />
                ))}
            </nav>

            {/* Bottom nav */}
            <div className="p-3 border-t border-[var(--border)] space-y-1">
                {BOTTOM_ITEMS.map((item) => (
                    <NavItem
                        key={item.id}
                        item={item}
                        active={active === item.id}
                        onClick={() => onChange(item.id)}
                    />
                ))}
            </div>
        </aside>
    );
}
