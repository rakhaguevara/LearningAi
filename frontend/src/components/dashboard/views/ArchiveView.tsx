'use client';

import { useState } from 'react';
import { motion } from 'framer-motion';

const SAMPLE_FILES = [
    { id: 1, name: 'Neural Networks Basics.pdf', type: 'PDF', size: '2.4 MB', date: '2h ago', tags: ['AI', 'ML'], color: 'red' },
    { id: 2, name: 'Python ML Course Notes.txt', type: 'TXT', size: '48 KB', date: '1d ago', tags: ['Python', 'ML'], color: 'blue' },
    { id: 3, name: 'Deep Learning Paper.pdf', type: 'PDF', size: '5.1 MB', date: '3d ago', tags: ['Research', 'DL'], color: 'red' },
    { id: 4, name: 'Data Structures Handbook.pdf', type: 'PDF', size: '3.8 MB', date: '5d ago', tags: ['CS', 'Algorithms'], color: 'red' },
    { id: 5, name: 'Interview Prep Notes.md', type: 'MD', size: '120 KB', date: '1w ago', tags: ['Career'], color: 'green' },
    { id: 6, name: 'React Patterns Cheatsheet.txt', type: 'TXT', size: '22 KB', date: '2w ago', tags: ['Frontend', 'React'], color: 'blue' },
];

const TYPE_COLORS: Record<string, string> = {
    PDF: 'bg-red-500/15 text-red-400',
    TXT: 'bg-blue-500/15 text-blue-400',
    MD: 'bg-green-500/15 text-green-400',
};

const TAG_FILTERS = ['All', 'AI', 'ML', 'Python', 'CS', 'Career', 'Frontend', 'Research'];

export function ArchiveView() {
    const [search, setSearch] = useState('');
    const [activeTag, setActiveTag] = useState('All');

    const filtered = SAMPLE_FILES.filter((f) => {
        const matchSearch = f.name.toLowerCase().includes(search.toLowerCase());
        const matchTag = activeTag === 'All' || f.tags.includes(activeTag);
        return matchSearch && matchTag;
    });

    return (
        <div className="p-6 h-full overflow-y-auto">
            {/* Header row */}
            <div className="flex flex-col sm:flex-row items-stretch sm:items-center gap-3 mb-6">
                {/* Search */}
                <div className="relative flex-1 max-w-full sm:max-w-sm">
                    <svg className="absolute left-3 top-1/2 -translate-y-1/2 w-4 h-4 text-[var(--text-muted)]"
                        fill="none" stroke="currentColor" viewBox="0 0 24 24">
                        <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z" />
                    </svg>
                    <input
                        value={search}
                        onChange={(e) => setSearch(e.target.value)}
                        placeholder="Search documents..."
                        className="w-full pl-9 pr-4 py-2.5 text-sm rounded-xl border border-[var(--border)] bg-[var(--bg-overlay)] text-[var(--text-primary)] placeholder-[var(--text-muted)] focus:outline-none focus:border-violet-500/50 transition-all"
                    />
                </div>

                <motion.button
                    whileHover={{ scale: 1.02 }}
                    whileTap={{ scale: 0.97 }}
                    className="flex items-center gap-2 px-4 py-2.5 rounded-xl bg-gradient-to-r from-violet-600 to-fuchsia-600 text-white text-sm font-medium shadow-lg shadow-violet-500/20 hover:shadow-violet-500/35 transition-all"
                >
                    <svg className="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                        <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M12 4v16m8-8H4" />
                    </svg>
                    Upload
                </motion.button>
            </div>

            {/* Tag filters */}
            <div className="flex items-center gap-2 mb-6 flex-wrap">
                {TAG_FILTERS.map((tag) => (
                    <button
                        key={tag}
                        onClick={() => setActiveTag(tag)}
                        className={`text-xs px-3 py-1.5 rounded-full font-medium transition-all ${activeTag === tag
                                ? 'bg-violet-600 text-white shadow-md shadow-violet-500/25'
                                : 'border border-[var(--border)] text-[var(--text-muted)] hover:text-[var(--text-primary)] hover:border-violet-500/40 bg-[var(--bg-overlay)]'
                            }`}
                    >
                        {tag}
                    </button>
                ))}
            </div>

            {/* Stats row */}
            <div className="grid grid-cols-1 sm:grid-cols-3 gap-3 mb-6">
                {[
                    { label: 'Total Files', value: SAMPLE_FILES.length, icon: '📁' },
                    { label: 'Total Size', value: '11.5 MB', icon: '💾' },
                    { label: 'This Week', value: 2, icon: '📅' },
                ].map((stat) => (
                    <div key={stat.label} className="p-3 rounded-xl border border-[var(--border)] bg-[var(--bg-overlay)] flex items-center gap-3">
                        <span className="text-xl">{stat.icon}</span>
                        <div>
                            <p className="text-lg font-bold text-[var(--text-primary)]">{stat.value}</p>
                            <p className="text-[10px] text-[var(--text-muted)]">{stat.label}</p>
                        </div>
                    </div>
                ))}
            </div>

            {/* File grid */}
            <div className="grid grid-cols-1 md:grid-cols-2 xl:grid-cols-3 gap-3">
                {filtered.map((file, i) => (
                    <motion.div
                        key={file.id}
                        initial={{ opacity: 0, y: 12 }}
                        animate={{ opacity: 1, y: 0 }}
                        transition={{ delay: i * 0.05 }}
                        whileHover={{ y: -3, boxShadow: '0 12px 30px rgba(124,58,237,0.15)' }}
                        className="p-4 rounded-2xl border border-[var(--border)] bg-[var(--bg-surface)] cursor-pointer group transition-all"
                    >
                        {/* File header */}
                        <div className="flex items-start gap-3 mb-3">
                            <div className={`w-10 h-10 rounded-xl flex items-center justify-center text-xs font-bold flex-shrink-0 ${TYPE_COLORS[file.type]}`}>
                                {file.type}
                            </div>
                            <div className="min-w-0 flex-1">
                                <p className="text-sm font-medium text-[var(--text-primary)] truncate group-hover:text-violet-400 transition-colors">
                                    {file.name}
                                </p>
                                <p className="text-xs text-[var(--text-muted)] mt-0.5">{file.size} · {file.date}</p>
                            </div>
                            <button className="opacity-0 group-hover:opacity-100 transition-opacity text-[var(--text-muted)] hover:text-[var(--text-primary)]">
                                <svg className="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                                    <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M12 5v.01M12 12v.01M12 19v.01M12 6a1 1 0 110-2 1 1 0 010 2zm0 7a1 1 0 110-2 1 1 0 010 2zm0 7a1 1 0 110-2 1 1 0 010 2z" />
                                </svg>
                            </button>
                        </div>

                        {/* Tags */}
                        <div className="flex gap-1.5 flex-wrap">
                            {file.tags.map((tag) => (
                                <span key={tag} className="text-[10px] px-2 py-0.5 rounded-full bg-[var(--brand-soft)] text-violet-400 border border-violet-500/20 font-medium">
                                    {tag}
                                </span>
                            ))}
                        </div>

                        {/* Action row */}
                        <div className="mt-3 pt-3 border-t border-[var(--border)] flex gap-2 opacity-0 group-hover:opacity-100 transition-opacity">
                            <button className="flex-1 text-[10px] py-1.5 rounded-lg bg-[var(--bg-hover)] text-[var(--text-secondary)] hover:text-[var(--text-primary)] transition-colors font-medium">
                                View
                            </button>
                            <button className="flex-1 text-[10px] py-1.5 rounded-lg bg-violet-500/15 text-violet-400 hover:bg-violet-500/25 transition-colors font-medium">
                                Use in Chat
                            </button>
                        </div>
                    </motion.div>
                ))}
            </div>

            {filtered.length === 0 && (
                <div className="text-center py-16 text-[var(--text-muted)]">
                    <p className="text-4xl mb-3">📭</p>
                    <p className="text-sm">No documents found</p>
                </div>
            )}
        </div>
    );
}
