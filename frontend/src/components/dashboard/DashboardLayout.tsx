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

    const ActiveView = PAGE_VIEWS[page];

    return (
        <div className="flex h-screen overflow-hidden bg-[var(--bg-base)]">
            {/* Sidebar */}
            <Sidebar active={page} onChange={setPage} />

            {/* Main area */}
            <div className="flex-1 flex flex-col min-w-0 overflow-hidden">
                {/* Topbar */}
                <Topbar page={page} onPageChange={setPage} />

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
