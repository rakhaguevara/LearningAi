'use client';

import { useEffect, useRef, useState } from 'react';
import { motion, AnimatePresence } from 'framer-motion';

type TimerMode = 'focus' | 'short' | 'long';

const TIMER_DURATIONS: Record<TimerMode, number> = {
    focus: 25 * 60,
    short: 5 * 60,
    long: 15 * 60,
};

const MODE_LABELS: Record<TimerMode, string> = {
    focus: 'Focus',
    short: 'Short Break',
    long: 'Long Break',
};

const HISTORY = [
    { label: 'Neural Networks Study', duration: '25 min', date: 'Today, 09:00' },
    { label: 'Python ML Review', duration: '25 min', date: 'Today, 08:00' },
    { label: 'Short Break', duration: '5 min', date: 'Today, 08:25' },
    { label: 'Data Structures', duration: '25 min', date: 'Yesterday, 21:00' },
];

const RADIUS = 90;
const CIRCUMFERENCE = 2 * Math.PI * RADIUS;

function formatTime(seconds: number) {
    const m = Math.floor(seconds / 60);
    const s = seconds % 60;
    return `${String(m).padStart(2, '0')}:${String(s).padStart(2, '0')}`;
}

export function PomodoroView() {
    const [mode, setMode] = useState<TimerMode>('focus');
    const [timeLeft, setTimeLeft] = useState(TIMER_DURATIONS['focus']);
    const [running, setRunning] = useState(false);
    const [completedSessions, setCompletedSessions] = useState(3);
    const intervalRef = useRef<NodeJS.Timeout | null>(null);

    const total = TIMER_DURATIONS[mode];
    const progress = timeLeft / total;
    const strokeDashoffset = CIRCUMFERENCE * (1 - progress);

    // Switch mode
    function switchMode(m: TimerMode) {
        setRunning(false);
        setMode(m);
        setTimeLeft(TIMER_DURATIONS[m]);
        if (intervalRef.current) clearInterval(intervalRef.current);
    }

    // Tick
    useEffect(() => {
        if (running) {
            intervalRef.current = setInterval(() => {
                setTimeLeft((t) => {
                    if (t <= 1) {
                        clearInterval(intervalRef.current!);
                        setRunning(false);
                        if (mode === 'focus') setCompletedSessions((n) => n + 1);
                        return 0;
                    }
                    return t - 1;
                });
            }, 1000);
        } else {
            if (intervalRef.current) clearInterval(intervalRef.current);
        }
        return () => { if (intervalRef.current) clearInterval(intervalRef.current); };
    }, [running, mode]);

    function reset() {
        setRunning(false);
        setTimeLeft(TIMER_DURATIONS[mode]);
    }

    const ringGradientId = `ring-grad-${mode}`;

    return (
        <div className="h-full overflow-y-auto p-6 flex flex-col items-center">
            <div className="w-full max-w-2xl">

                {/* Mode selector */}
                <div className="flex gap-2 p-1.5 rounded-2xl bg-[var(--bg-overlay)] border border-[var(--border)] mb-10 w-fit mx-auto">
                    {(Object.keys(TIMER_DURATIONS) as TimerMode[]).map((m) => (
                        <button
                            key={m}
                            onClick={() => switchMode(m)}
                            className={`px-5 py-2 rounded-xl text-sm font-medium transition-all ${mode === m
                                    ? 'bg-gradient-to-r from-violet-600 to-fuchsia-600 text-white shadow-lg shadow-violet-500/25'
                                    : 'text-[var(--text-muted)] hover:text-[var(--text-primary)]'
                                }`}
                        >
                            {MODE_LABELS[m]}
                        </button>
                    ))}
                </div>

                {/* Ring + Timer */}
                <div className="flex flex-col items-center mb-10">
                    <div className="relative w-56 h-56">
                        <svg
                            width="224" height="224"
                            viewBox="0 0 224 224"
                            className="-rotate-90"
                        >
                            <defs>
                                <linearGradient id={ringGradientId} x1="0%" y1="0%" x2="100%" y2="0%">
                                    <stop offset="0%" stopColor="#7c3aed" />
                                    <stop offset="100%" stopColor="#ec4899" />
                                </linearGradient>
                            </defs>
                            {/* Track */}
                            <circle
                                cx="112" cy="112" r={RADIUS}
                                fill="none"
                                stroke="var(--border-strong)"
                                strokeWidth="6"
                            />
                            {/* Progress */}
                            <motion.circle
                                cx="112" cy="112" r={RADIUS}
                                fill="none"
                                stroke={`url(#${ringGradientId})`}
                                strokeWidth="6"
                                strokeLinecap="round"
                                strokeDasharray={CIRCUMFERENCE}
                                strokeDashoffset={strokeDashoffset}
                                style={{ transition: 'stroke-dashoffset 1s linear' }}
                            />
                        </svg>

                        {/* Center text */}
                        <div className="absolute inset-0 flex flex-col items-center justify-center">
                            <AnimatePresence mode="wait">
                                <motion.span
                                    key={timeLeft}
                                    className="text-5xl font-bold tabular-nums text-[var(--text-primary)]"
                                    style={{ fontVariantNumeric: 'tabular-nums' }}
                                >
                                    {formatTime(timeLeft)}
                                </motion.span>
                            </AnimatePresence>
                            <span className="text-xs text-[var(--text-muted)] mt-1.5 font-medium">{MODE_LABELS[mode]}</span>
                        </div>
                    </div>

                    {/* Controls */}
                    <div className="flex items-center gap-4 mt-8">
                        <motion.button
                            whileHover={{ scale: 1.05 }}
                            whileTap={{ scale: 0.95 }}
                            onClick={reset}
                            className="w-11 h-11 rounded-2xl border border-[var(--border)] bg-[var(--bg-overlay)] flex items-center justify-center text-[var(--text-muted)] hover:text-[var(--text-primary)] hover:bg-[var(--bg-hover)] transition-all"
                        >
                            <svg className="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M4 4v5h.582m15.356 2A8.001 8.001 0 004.582 9m0 0H9m11 11v-5h-.581m0 0a8.003 8.003 0 01-15.357-2m15.357 2H15" />
                            </svg>
                        </motion.button>

                        <motion.button
                            whileHover={{ scale: 1.04 }}
                            whileTap={{ scale: 0.96 }}
                            onClick={() => setRunning((r) => !r)}
                            className="w-16 h-16 rounded-2xl bg-gradient-to-br from-violet-600 to-fuchsia-600 flex items-center justify-center text-white shadow-xl shadow-violet-500/30 hover:shadow-violet-500/50 transition-all"
                        >
                            {running ? (
                                <svg className="w-6 h-6" fill="currentColor" viewBox="0 0 24 24">
                                    <path d="M6 4h4v16H6V4zm8 0h4v16h-4V4z" />
                                </svg>
                            ) : (
                                <svg className="w-6 h-6 ml-0.5" fill="currentColor" viewBox="0 0 24 24">
                                    <path d="M8 5v14l11-7z" />
                                </svg>
                            )}
                        </motion.button>

                        <motion.button
                            whileHover={{ scale: 1.05 }}
                            whileTap={{ scale: 0.95 }}
                            onClick={() => switchMode(mode === 'focus' ? 'short' : 'focus')}
                            className="w-11 h-11 rounded-2xl border border-[var(--border)] bg-[var(--bg-overlay)] flex items-center justify-center text-[var(--text-muted)] hover:text-[var(--text-primary)] hover:bg-[var(--bg-hover)] transition-all"
                        >
                            <svg className="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M13 5l7 7-7 7M5 5l7 7-7 7" />
                            </svg>
                        </motion.button>
                    </div>

                    {/* Session dots */}
                    <div className="flex gap-2 mt-6 items-center">
                        <span className="text-xs text-[var(--text-muted)] mr-1">Sessions:</span>
                        {Array.from({ length: 4 }).map((_, i) => (
                            <div
                                key={i}
                                className={`w-2.5 h-2.5 rounded-full transition-all ${i < completedSessions
                                        ? 'bg-gradient-to-r from-violet-500 to-fuchsia-500'
                                        : 'bg-[var(--border-strong)]'
                                    }`}
                            />
                        ))}
                        <span className="text-xs text-[var(--text-muted)] ml-1">{completedSessions}/4</span>
                    </div>
                </div>

                {/* Session history */}
                <div className="w-full">
                    <h3 className="text-xs font-semibold tracking-widest uppercase text-[var(--text-muted)] mb-3">
                        Session History
                    </h3>
                    <div className="space-y-2">
                        {HISTORY.map((item, i) => (
                            <motion.div
                                key={i}
                                initial={{ opacity: 0, x: -10 }}
                                animate={{ opacity: 1, x: 0 }}
                                transition={{ delay: i * 0.05 }}
                                className="flex items-center gap-4 p-3 rounded-xl border border-[var(--border)] bg-[var(--bg-overlay)]"
                            >
                                <div className="w-2 h-8 rounded-full bg-gradient-to-b from-violet-500 to-fuchsia-500 flex-shrink-0" />
                                <div className="flex-1 min-w-0">
                                    <p className="text-sm font-medium text-[var(--text-primary)] truncate">{item.label}</p>
                                    <p className="text-xs text-[var(--text-muted)]">{item.date}</p>
                                </div>
                                <span className="text-xs font-semibold text-violet-400 bg-violet-500/10 px-2.5 py-1 rounded-full border border-violet-500/20">
                                    {item.duration}
                                </span>
                            </motion.div>
                        ))}
                    </div>
                </div>
            </div>
        </div>
    );
}
