'use client';

import { useState } from 'react';
import { AnimatePresence, motion } from 'framer-motion';
import { Sidebar, type DashboardPage } from './Sidebar';
import { Topbar } from './Topbar';
import { LearnNowView } from './views/LearnNowView';
import { ArchiveView } from './views/ArchiveView';
import { PomodoroView } from './views/PomodoroView';
import { ProfileView } from './views/ProfileView';
import { SettingsView } from './views/SettingsView';

const PAGE_VIEWS: Record<DashboardPage, React.ComponentType> = {
    learn: LearnNowView,
    archive: ArchiveView,
    pomodoro: PomodoroView,
    profile: ProfileView,
    settings: SettingsView,
};

const pageVariants = {
    initial: { opacity: 0, y: 10 },
    enter: { opacity: 1, y: 0 },
    exit: { opacity: 0, y: -6 },
};

export function DashboardLayout() {
    const [page, setPage] = useState<DashboardPage>('learn');
    const [mobileSidebarOpen, setMobileSidebarOpen] = useState(false);

    const ActiveView = PAGE_VIEWS[page];

    return (
        <div className="flex h-screen overflow-hidden bg-[var(--bg-base)] relative">
            {/* Background Orbs - Matching landing page theme */}
            <div className="absolute inset-0 pointer-events-none overflow-hidden">
                <div className="absolute top-[-10%] left-[-5%] w-96 h-96 rounded-full bg-violet-600/20 blur-3xl animate-pulse" />
                <div className="absolute bottom-[-10%] right-[-5%] w-96 h-96 rounded-full bg-fuchsia-600/20 blur-3xl animate-pulse" style={{ animationDelay: '1.5s' }} />
                <div className="absolute top-[20%] right-[10%] w-64 h-64 rounded-full bg-cyan-600/15 blur-3xl animate-pulse" style={{ animationDelay: '3s' }} />
            </div>
            {/* Desktop Sidebar */}
            <div className="hidden md:block">
                <Sidebar active={page} onChange={setPage} />
            </div>

            {/* Mobile Sidebar Overlay */}
            {mobileSidebarOpen && (
                <div 
                    className="fixed inset-0 bg-black/50 z-40 md:hidden"
                    onClick={() => setMobileSidebarOpen(false)}
                />
            )}

            {/* Mobile Sidebar */}
            <div className={`fixed inset-y-0 left-0 z-50 transform transition-transform duration-300 md:hidden ${mobileSidebarOpen ? 'translate-x-0' : '-translate-x-full'}`}>
                <Sidebar active={page} onChange={(p) => { setPage(p); setMobileSidebarOpen(false); }} />
            </div>

            {/* Main area */}
            <div className="flex-1 flex flex-col min-w-0 overflow-hidden">
                {/* Topbar */}
                <Topbar page={page} onPageChange={setPage} onMenuClick={() => setMobileSidebarOpen(true)} />

                {/* Content */}
                <main className="flex-1 overflow-hidden relative">
                    <AnimatePresence mode="wait" initial={false}>
                        <motion.div
                            key={page}
                            variants={pageVariants}
                            initial="initial"
                            animate="enter"
                            exit="exit"
                            transition={{ duration: 0.22, ease: 'easeOut' }}
                            className="absolute inset-0"
                        >
                            <ActiveView />
                        </motion.div>
                    </AnimatePresence>
                </main>
            </div>
        </div>
    );
}
