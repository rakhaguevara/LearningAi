'use client';

import { motion } from 'framer-motion';

const PROFILE_DATA = {
    name: 'Neura User',
    email: 'user@example.com',
    learningStyle: 'Stories & Analogies',
    depthPreference: 'Concept + Examples',
    interestThemes: ['Technology', 'Science', 'Gaming'],
    analogyTheme: 'Gaming',
    aiRetryPreference: 'Ask for Examples',
    studyFocus: 'Machine learning fundamentals and neural networks',
    fileUploadHabit: 'Sometimes',
    learningGoal: 'Build real projects & ship things',
    longContentBehavior: 'Break it into pieces',
};

const AI_SUMMARY = `Based on your learning profile, you're a hands-on learner who loves to see concepts applied in the real world — especially through the lens of gaming and technology. You prefer bite-sized content with clear examples, and you're most motivated when working toward a concrete project outcome. Your AI tutor is calibrated to deliver concise, analogy-rich explanations at an intermediate-to-advanced depth.`;

interface LabeledRowProps {
    label: string;
    value: string;
    icon?: string;
}

function LabeledRow({ label, value, icon }: LabeledRowProps) {
    return (
        <div className="flex items-start justify-between py-3 border-b border-[var(--border)] last:border-none gap-4">
            <div className="flex items-center gap-2">
                {icon && <span className="text-base">{icon}</span>}
                <span className="text-sm text-[var(--text-muted)]">{label}</span>
            </div>
            <span className="text-sm font-medium text-[var(--text-primary)] text-right max-w-[55%]">{value}</span>
        </div>
    );
}

export function ProfileView() {
    return (
        <div className="h-full overflow-y-auto p-6">
            <div className="max-w-3xl mx-auto space-y-5">

                {/* Avatar + Brief */}
                <motion.div
                    initial={{ opacity: 0, y: 16 }}
                    animate={{ opacity: 1, y: 0 }}
                    className="flex items-center gap-5 p-5 rounded-2xl border border-[var(--border)] bg-[var(--bg-surface)]"
                >
                    <div className="w-16 h-16 rounded-2xl bg-gradient-to-br from-violet-500 to-fuchsia-500 flex items-center justify-center text-white text-2xl font-bold shadow-xl shadow-violet-500/30 flex-shrink-0">
                        N
                    </div>
                    <div className="flex-1 min-w-0">
                        <h2 className="text-xl font-bold text-[var(--text-primary)]">{PROFILE_DATA.name}</h2>
                        <p className="text-sm text-[var(--text-muted)]">{PROFILE_DATA.email}</p>
                        <div className="flex gap-2 mt-2 flex-wrap">
                            {PROFILE_DATA.interestThemes.map((t) => (
                                <span key={t} className="text-[10px] font-semibold px-2 py-0.5 rounded-full bg-[var(--brand-soft)] text-violet-400 border border-violet-500/20">
                                    {t}
                                </span>
                            ))}
                        </div>
                    </div>
                    <motion.button
                        whileHover={{ scale: 1.03 }}
                        whileTap={{ scale: 0.97 }}
                        className="px-4 py-2 rounded-xl border border-[var(--border)] text-sm font-medium text-[var(--text-secondary)] hover:text-[var(--text-primary)] hover:border-violet-500/40 transition-all flex-shrink-0"
                    >
                        Edit Profile
                    </motion.button>
                </motion.div>

                {/* AI Generated Summary */}
                <motion.div
                    initial={{ opacity: 0, y: 16 }}
                    animate={{ opacity: 1, y: 0 }}
                    transition={{ delay: 0.08 }}
                    className="relative p-5 rounded-2xl overflow-hidden"
                    style={{
                        background: 'linear-gradient(135deg, rgba(124,58,237,0.12) 0%, rgba(236,72,153,0.08) 100%)',
                        border: '1px solid rgba(124,58,237,0.3)',
                    }}
                >
                    {/* Glow blob */}
                    <div className="absolute -top-8 -right-8 w-32 h-32 rounded-full bg-violet-500/20 blur-2xl pointer-events-none" />
                    <div className="relative z-10">
                        <div className="flex items-center gap-2 mb-3">
                            <div className="w-6 h-6 rounded-lg bg-gradient-to-br from-violet-500 to-fuchsia-500 flex items-center justify-center shadow-md">
                                <svg className="w-3 h-3 text-white" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                                    <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2.5} d="M13 10V3L4 14h7v7l9-11h-7z" />
                                </svg>
                            </div>
                            <span className="text-xs font-semibold tracking-widest uppercase text-violet-400">AI Learning Summary</span>
                            <span className="text-[10px] text-violet-400/50 ml-auto">Generated by NeuraLearn AI</span>
                        </div>
                        <p className="text-sm text-[var(--text-secondary)] leading-relaxed">{AI_SUMMARY}</p>
                    </div>
                </motion.div>

                {/* Learning Profile Details */}
                <motion.div
                    initial={{ opacity: 0, y: 16 }}
                    animate={{ opacity: 1, y: 0 }}
                    transition={{ delay: 0.14 }}
                    className="p-5 rounded-2xl border border-[var(--border)] bg-[var(--bg-surface)]"
                >
                    <div className="flex items-center justify-between mb-4">
                        <h3 className="text-xs font-semibold tracking-widest uppercase text-[var(--text-muted)]">Learning Profile</h3>
                        <span className="text-[10px] text-violet-400 bg-violet-500/10 px-2 py-0.5 rounded-full border border-violet-500/20">
                            From user_learning_profiles
                        </span>
                    </div>
                    <LabeledRow icon="🎨" label="Learning Style" value={PROFILE_DATA.learningStyle} />
                    <LabeledRow icon="🔬" label="Explanation Depth" value={PROFILE_DATA.depthPreference} />
                    <LabeledRow icon="🎯" label="Analogy Theme" value={PROFILE_DATA.analogyTheme} />
                    <LabeledRow icon="🔄" label="AI Retry Preference" value={PROFILE_DATA.aiRetryPreference} />
                    <LabeledRow icon="📚" label="Long Content Behavior" value={PROFILE_DATA.longContentBehavior} />
                    <LabeledRow icon="📁" label="File Upload Habit" value={PROFILE_DATA.fileUploadHabit} />
                </motion.div>

                {/* Study Focus + Goal */}
                <div className="grid grid-cols-2 gap-3">
                    {[
                        {
                            label: 'Current Study Focus',
                            value: PROFILE_DATA.studyFocus,
                            icon: '🎓',
                            source: 'user_behavior_signals',
                            delay: 0.18,
                        },
                        {
                            label: 'Learning Goal',
                            value: PROFILE_DATA.learningGoal,
                            icon: '🏗️',
                            source: 'user_learning_profiles',
                            delay: 0.22,
                        },
                    ].map((card) => (
                        <motion.div
                            key={card.label}
                            initial={{ opacity: 0, y: 16 }}
                            animate={{ opacity: 1, y: 0 }}
                            transition={{ delay: card.delay }}
                            className="p-4 rounded-2xl border border-[var(--border)] bg-[var(--bg-surface)]"
                        >
                            <div className="flex items-center gap-2 mb-2">
                                <span>{card.icon}</span>
                                <span className="text-xs text-[var(--text-muted)]">{card.label}</span>
                            </div>
                            <p className="text-sm font-medium text-[var(--text-primary)]">{card.value}</p>
                            <span className="text-[9px] text-[var(--text-muted)] mt-2 inline-block">← {card.source}</span>
                        </motion.div>
                    ))}
                </div>
            </div>
        </div>
    );
}
