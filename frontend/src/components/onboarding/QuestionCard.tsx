'use client';

import { motion, AnimatePresence } from 'framer-motion';
import { OptionButton } from './OptionButton';

interface Option {
    value: string;
    label: string;
    emoji?: string;
}

interface QuestionCardProps {
    stepNumber: number;
    question: string;
    subtitle?: string;
    options: Option[];
    selectedValue: string | string[];
    isMulti?: boolean;
    onSelect: (value: string) => void;
    isTextInput?: boolean;
    textValue?: string;
    onTextChange?: (value: string) => void;
    direction: number; // 1 = forward, -1 = backward
}

const variants = {
    enter: (direction: number) => ({
        x: direction > 0 ? 60 : -60,
        opacity: 0,
    }),
    center: {
        x: 0,
        opacity: 1,
    },
    exit: (direction: number) => ({
        x: direction > 0 ? -60 : 60,
        opacity: 0,
    }),
};

export function QuestionCard({
    stepNumber,
    question,
    subtitle,
    options,
    selectedValue,
    isMulti = false,
    onSelect,
    isTextInput = false,
    textValue = '',
    onTextChange,
    direction,
}: QuestionCardProps) {
    const isSelected = (value: string) => {
        if (Array.isArray(selectedValue)) return selectedValue.includes(value);
        return selectedValue === value;
    };

    return (
        <motion.div
            key={stepNumber}
            custom={direction}
            variants={variants}
            initial="enter"
            animate="center"
            exit="exit"
            transition={{ duration: 0.35, ease: [0.25, 0.46, 0.45, 0.94] }}
            className="w-full"
        >
            {/* Step label */}
            <div className="mb-6">
                <span className="text-xs font-semibold tracking-widest uppercase text-violet-400/80">
                    Question {stepNumber}
                </span>
                <h2 className="mt-2 text-2xl md:text-3xl font-bold text-white leading-snug">
                    {question}
                </h2>
                {subtitle && (
                    <p className="mt-2 text-white/50 text-sm">{subtitle}</p>
                )}
            </div>

            {/* Options or text input */}
            {isTextInput ? (
                <div className="space-y-3">
                    <textarea
                        id={`onboarding-text-q${stepNumber}`}
                        value={textValue}
                        onChange={(e) => onTextChange?.(e.target.value)}
                        maxLength={200}
                        rows={3}
                        placeholder="e.g. Machine learning, data structures, web dev..."
                        className="w-full px-5 py-4 bg-white/5 border border-white/10 rounded-xl text-white placeholder-white/30 focus:outline-none focus:ring-2 focus:ring-violet-500 resize-none transition-all"
                    />
                    <div className="text-right text-xs text-white/30">
                        {textValue.length}/200
                    </div>
                </div>
            ) : (
                <div className="grid gap-3">
                    {isMulti && (
                        <p className="text-xs text-white/40 mb-1">Select all that apply</p>
                    )}
                    {options.map((opt) => (
                        <OptionButton
                            key={opt.value}
                            id={`option-q${stepNumber}-${opt.value}`}
                            label={opt.label}
                            emoji={opt.emoji}
                            selected={isSelected(opt.value)}
                            onClick={() => onSelect(opt.value)}
                        />
                    ))}
                </div>
            )}
        </motion.div>
    );
}
