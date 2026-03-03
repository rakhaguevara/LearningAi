'use client';

import { useState } from 'react';
import { AnimatePresence, motion } from 'framer-motion';
import { QuestionCard } from './QuestionCard';
import { ProgressBar } from './ProgressBar';
import { getAuthToken, authHeaders } from '@/lib/auth';
import { siteConfig } from '@/lib/constants';

// ── Question definitions ────────────────────────────────────────────
const QUESTIONS = [
    {
        key: 'learning_style',
        question: 'How do you prefer to learn new things?',
        subtitle: 'This shapes how we deliver every explanation.',
        options: [
            { value: 'visual', label: 'Visuals & diagrams — show me, don\'t tell me', emoji: '🎨' },
            { value: 'step_by_step', label: 'Step-by-step — one thing at a time, please', emoji: '🪜' },
            { value: 'stories', label: 'Stories & analogies — make it relatable', emoji: '📖' },
            { value: 'hands_on', label: 'Hands-on — just give me something to try', emoji: '🛠️' },
        ],
    },
    {
        key: 'long_content_behavior',
        question: 'When content feels too long, what do you do?',
        subtitle: 'Helps us calibrate response length for you.',
        options: [
            { value: 'summarize', label: 'Ask for a shorter summary', emoji: '⚡' },
            { value: 'keep_going', label: 'Power through — I like depth', emoji: '💪' },
            { value: 'break_up', label: 'Break it into smaller pieces', emoji: '🧩' },
            { value: 'give_examples', label: 'Just give me a quick example instead', emoji: '💡' },
        ],
    },
    {
        key: 'explanation_format',
        question: 'What\'s your preferred explanation format?',
        options: [
            { value: 'bullet_points', label: 'Bullet points — scannable & concise', emoji: '📌' },
            { value: 'paragraphs', label: 'Flowing paragraphs — full context', emoji: '📝' },
            { value: 'code_examples', label: 'Code examples — show me the logic', emoji: '💻' },
            { value: 'diagrams', label: 'Diagrams & flowcharts — visual mapping', emoji: '📊' },
        ],
    },
    {
        key: 'interest_themes',
        question: 'What topics excite you the most?',
        subtitle: 'Select all that apply — we\'ll weave them into examples.',
        isMulti: true,
        options: [
            { value: 'tech', label: 'Technology & software', emoji: '🤖' },
            { value: 'science', label: 'Science & research', emoji: '🔬' },
            { value: 'arts', label: 'Arts & design', emoji: '🎭' },
            { value: 'sports', label: 'Sports & fitness', emoji: '⚽' },
            { value: 'gaming', label: 'Gaming & esports', emoji: '🎮' },
            { value: 'business', label: 'Business & finance', emoji: '📈' },
            { value: 'music', label: 'Music & audio', emoji: '🎵' },
        ],
    },
    {
        key: 'analogy_theme',
        question: 'Where do your best "aha!" moments come from?',
        subtitle: 'We\'ll pull analogies from the world that resonates with you.',
        options: [
            { value: 'sports', label: 'Sports plays & strategy', emoji: '🏆' },
            { value: 'gaming', label: 'Game mechanics & quests', emoji: '🎯' },
            { value: 'cooking', label: 'Recipes & cooking', emoji: '👨‍🍳' },
            { value: 'movies', label: 'Movie plot & characters', emoji: '🎬' },
            { value: 'nature', label: 'Nature & science', emoji: '🌿' },
        ],
    },
    {
        key: 'depth_preference',
        question: 'How deep do you want to go?',
        subtitle: 'Sets the default complexity of your explanations.',
        options: [
            { value: 'beginner_overview', label: 'High-level overview — what is this thing?', emoji: '🌤️' },
            { value: 'deep_dive', label: 'Deep dive — I want all the details', emoji: '🌊' },
            { value: 'concept_plus_examples', label: 'Concept + examples — best of both worlds', emoji: '⚖️' },
            { value: 'expert', label: 'Expert mode — skip the basics', emoji: '🚀' },
        ],
    },
    {
        key: 'ai_retry_preference',
        question: 'When an AI explanation misses the mark, you…',
        options: [
            { value: 'rephrase', label: 'Rephrase my question differently', emoji: '🔄' },
            { value: 'try_simpler', label: 'Ask for a simpler version', emoji: '🧸' },
            { value: 'ask_examples', label: 'Request a concrete example', emoji: '🔍' },
            { value: 'give_up', label: 'Give up and search elsewhere', emoji: '😤' },
        ],
    },
    {
        key: 'study_focus',
        question: 'What are you studying or want to learn right now?',
        subtitle: 'A few words is enough — topic, tech, skill...',
        isTextInput: true,
        options: [],
    },
    {
        key: 'file_upload_habit',
        question: 'Would you upload your own notes or PDFs to assist the AI?',
        subtitle: 'Helps us understand how much personal context you\'d share.',
        options: [
            { value: 'yes_always', label: 'Yes — I\'d love my notes to be used', emoji: '📂' },
            { value: 'sometimes', label: 'Sometimes, for specific topics', emoji: '📎' },
            { value: 'rarely', label: 'Rarely — usually not worth it', emoji: '🤷' },
            { value: 'no', label: 'No — I prefer keeping things separate', emoji: '🔒' },
        ],
    },
    {
        key: 'learning_goal',
        question: 'What\'s your primary learning goal?',
        options: [
            { value: 'pass_exams', label: 'Pass exams & ace assessments', emoji: '🎓' },
            { value: 'learn_for_fun', label: 'Learn for the joy of it', emoji: '✨' },
            { value: 'career_change', label: 'Pivot my career or skill up professionally', emoji: '💼' },
            { value: 'build_projects', label: 'Build real projects & ship things', emoji: '🏗️' },
        ],
    },
] as const;

type Answers = {
    learning_style: string;
    long_content_behavior: string;
    explanation_format: string;
    interest_themes: string[];
    analogy_theme: string;
    depth_preference: string;
    ai_retry_preference: string;
    study_focus: string;
    file_upload_habit: string;
    learning_goal: string;
};

const INITIAL_ANSWERS: Answers = {
    learning_style: '',
    long_content_behavior: '',
    explanation_format: '',
    interest_themes: [],
    analogy_theme: '',
    depth_preference: '',
    ai_retry_preference: '',
    study_focus: '',
    file_upload_habit: '',
    learning_goal: '',
};

export function OnboardingWizard() {
    const [step, setStep] = useState(0);
    const [direction, setDirection] = useState(1);
    const [answers, setAnswers] = useState<Answers>(INITIAL_ANSWERS);
    const [isSubmitting, setIsSubmitting] = useState(false);
    const [error, setError] = useState('');

    const currentQ = QUESTIONS[step];
    const total = QUESTIONS.length;

    const currentValue = answers[currentQ.key as keyof Answers];

    // ── Answer handlers ──────────────────────────────────────────────
    function handleSelect(value: string) {
        const key = currentQ.key as keyof Answers;
        if ('isMulti' in currentQ && currentQ.isMulti) {
            const prev = answers.interest_themes;
            const updated = prev.includes(value)
                ? prev.filter((v) => v !== value)
                : [...prev, value];
            setAnswers((a) => ({ ...a, interest_themes: updated }));
        } else {
            setAnswers((a) => ({ ...a, [key]: value }));
        }
    }

    function handleTextChange(value: string) {
        setAnswers((a) => ({ ...a, study_focus: value }));
    }

    // ── Navigation ───────────────────────────────────────────────────
    function canProceed() {
        if ('isTextInput' in currentQ && currentQ.isTextInput) {
            return answers.study_focus.trim().length > 0;
        }
        if ('isMulti' in currentQ && currentQ.isMulti) {
            return answers.interest_themes.length > 0;
        }
        return typeof currentValue === 'string' && currentValue !== '';
    }

    function goNext() {
        if (step < total - 1) {
            setDirection(1);
            setStep((s) => s + 1);
        } else {
            handleSubmit();
        }
    }

    function goBack() {
        if (step > 0) {
            setDirection(-1);
            setStep((s) => s - 1);
        }
    }

    // ── Submission ───────────────────────────────────────────────────
    async function handleSubmit() {
        setIsSubmitting(true);
        setError('');

        try {
            const token = getAuthToken();

            // No token → can't authenticate. Send user back to login.
            if (!token) {
                setError('Session expired. Please log in again.');
                setIsSubmitting(false);
                setTimeout(() => { window.location.href = '/login'; }, 1500);
                return;
            }

            const res = await fetch(`${siteConfig.api}/onboarding/submit`, {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json',
                    ...authHeaders(),
                },
                credentials: 'include',
                body: JSON.stringify(answers),
            });

            if (!res.ok) {
                const data = await res.json().catch(() => ({}));
                const message =
                    data?.error?.detail ||
                    data?.error?.message ||
                    data?.message ||
                    (typeof data?.error === 'string' ? data.error : null) ||
                    `Request failed (${res.status})`;

                // Token invalid/expired → back to login
                if (res.status === 401) {
                    localStorage.removeItem('access_token');
                    setError('Session expired. Redirecting to login…');
                    setTimeout(() => { window.location.href = '/login'; }, 1500);
                    return;
                }

                throw new Error(message);
            }

            // Success → go to dashboard
            window.location.href = '/dashboard';
        } catch (err: any) {
            setError(err?.message || 'Something went wrong. Please try again.');
            setIsSubmitting(false);
        }
    }

    const isLastStep = step === total - 1;

    return (
        <div className="w-full max-w-lg mx-auto px-4">
            {/* Progress */}
            <div className="mb-8">
                <ProgressBar current={step + 1} total={total} />
            </div>

            {/* Card */}
            <div className="relative bg-white/5 backdrop-blur-xl border border-white/10 rounded-3xl p-8 shadow-2xl min-h-[420px] flex flex-col justify-between">
                <div className="flex-1">
                    <AnimatePresence mode="wait" custom={direction}>
                        <QuestionCard
                            key={step}
                            stepNumber={step + 1}
                            question={currentQ.question}
                            subtitle={'subtitle' in currentQ ? currentQ.subtitle : undefined}
                            options={currentQ.options as any}
                            selectedValue={currentValue}
                            isMulti={'isMulti' in currentQ && currentQ.isMulti}
                            onSelect={handleSelect}
                            isTextInput={'isTextInput' in currentQ && currentQ.isTextInput}
                            textValue={answers.study_focus}
                            onTextChange={handleTextChange}
                            direction={direction}
                        />
                    </AnimatePresence>
                </div>

                {/* Error banner */}
                {error && (
                    <div className="mt-4 p-3 text-sm text-red-400 bg-red-400/10 border border-red-400/20 rounded-lg">
                        {error}
                    </div>
                )}

                {/* Navigation */}
                <div className="flex items-center justify-between mt-8 gap-4">
                    <motion.button
                        onClick={goBack}
                        disabled={step === 0}
                        whileTap={{ scale: 0.95 }}
                        className="px-5 py-3 text-sm text-white/50 hover:text-white/80 disabled:opacity-0 disabled:pointer-events-none transition-colors flex items-center gap-2"
                    >
                        <svg className="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                            <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M15 19l-7-7 7-7" />
                        </svg>
                        Back
                    </motion.button>

                    <motion.button
                        id="onboarding-next-btn"
                        onClick={goNext}
                        disabled={!canProceed() || isSubmitting}
                        whileHover={canProceed() ? { scale: 1.03 } : {}}
                        whileTap={canProceed() ? { scale: 0.97 } : {}}
                        className={`
                            flex-1 py-3 px-6 rounded-xl font-semibold text-sm transition-all
                            flex items-center justify-center gap-2
                            ${canProceed() && !isSubmitting
                                ? 'bg-gradient-to-r from-violet-600 to-fuchsia-600 text-white shadow-lg shadow-violet-500/30 hover:from-violet-500 hover:to-fuchsia-500'
                                : 'bg-white/10 text-white/30 cursor-not-allowed'
                            }
                        `}
                    >
                        {isSubmitting ? (
                            <>
                                <svg className="animate-spin h-4 w-4 text-white" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24">
                                    <circle className="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" strokeWidth="4" />
                                    <path className="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z" />
                                </svg>
                                Setting up your profile&hellip;
                            </>
                        ) : isLastStep ? (
                            <>
                                Finish & Start Learning
                                <svg className="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                                    <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M5 13l4 4L19 7" />
                                </svg>
                            </>
                        ) : (
                            <>
                                Continue
                                <svg className="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                                    <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M9 5l7 7-7 7" />
                                </svg>
                            </>
                        )}
                    </motion.button>
                </div>
            </div>
        </div>
    );
}
