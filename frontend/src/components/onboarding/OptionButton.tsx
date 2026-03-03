'use client';

import { motion } from 'framer-motion';

interface OptionButtonProps {
    label: string;
    emoji?: string;
    selected: boolean;
    onClick: () => void;
    id: string;
}

export function OptionButton({ label, emoji, selected, onClick, id }: OptionButtonProps) {
    return (
        <motion.button
            id={id}
            onClick={onClick}
            whileHover={{ scale: 1.02, y: -2 }}
            whileTap={{ scale: 0.98 }}
            transition={{ type: 'spring', stiffness: 400, damping: 20 }}
            className={`
                w-full flex items-center gap-3 px-5 py-4 rounded-xl text-left
                border transition-all duration-200 relative overflow-hidden group
                ${selected
                    ? 'border-violet-500/80 bg-violet-500/20 text-white shadow-lg shadow-violet-500/20'
                    : 'border-white/10 bg-white/5 text-white/80 hover:border-white/25 hover:bg-white/10 hover:text-white'
                }
            `}
        >
            {/* Subtle glow pulse when selected */}
            {selected && (
                <motion.div
                    initial={{ opacity: 0 }}
                    animate={{ opacity: 1 }}
                    className="absolute inset-0 bg-gradient-to-r from-violet-500/10 to-fuchsia-500/10 pointer-events-none"
                />
            )}

            {emoji && (
                <span className="text-xl flex-shrink-0">{emoji}</span>
            )}

            <span className="text-sm font-medium flex-1">{label}</span>

            {selected && (
                <motion.div
                    initial={{ scale: 0 }}
                    animate={{ scale: 1 }}
                    className="w-5 h-5 rounded-full bg-violet-500 flex items-center justify-center flex-shrink-0"
                >
                    <svg className="w-3 h-3 text-white" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                        <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={3} d="M5 13l4 4L19 7" />
                    </svg>
                </motion.div>
            )}
        </motion.button>
    );
}
