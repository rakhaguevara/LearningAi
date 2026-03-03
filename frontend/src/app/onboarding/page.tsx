'use client';

import { useState, useEffect } from 'react';
import { motion, AnimatePresence } from 'framer-motion';
import { useRouter } from 'next/navigation';
import Link from 'next/link';

interface Interest {
    id: string;
    name: string;
    icon: string;
    category: string;
}

interface LearningStyle {
    id: string;
    name: string;
    description: string;
    icon: string;
}

const interests: Interest[] = [
    { id: 'anime', name: 'Anime', icon: '🎌', category: 'Entertainment' },
    { id: 'gaming', name: 'Gaming', icon: '🎮', category: 'Entertainment' },
    { id: 'sports', name: 'Sports', icon: '⚽', category: 'Lifestyle' },
    { id: 'music', name: 'Music', icon: '🎵', category: 'Arts' },
    { id: 'technology', name: 'Technology', icon: '💻', category: 'Science' },
    { id: 'cooking', name: 'Cooking', icon: '🍳', category: 'Lifestyle' },
    { id: 'movies', name: 'Movies', icon: '🎬', category: 'Entertainment' },
    { id: 'travel', name: 'Travel', icon: '✈️', category: 'Lifestyle' },
    { id: 'fashion', name: 'Fashion', icon: '👗', category: 'Lifestyle' },
    { id: 'fitness', name: 'Fitness', icon: '💪', category: 'Health' },
    { id: 'reading', name: 'Reading', icon: '📚', category: 'Education' },
    { id: 'art', name: 'Art & Design', icon: '🎨', category: 'Arts' },
];

const learningStyles: LearningStyle[] = [
    {
        id: 'visual',
        name: 'Visual Learner',
        description: 'I learn best through images, diagrams, and videos',
        icon: '👁️',
    },
    {
        id: 'interactive',
        name: 'Interactive',
        description: 'I prefer hands-on practice and interactive exercises',
        icon: '🤲',
    },
    {
        id: 'story',
        name: 'Story-based',
        description: 'I learn through narratives and real-world examples',
        icon: '📖',
    },
    {
        id: 'analytical',
        name: 'Analytical',
        description: 'I prefer structured, logical explanations',
        icon: '🧠',
    },
];

export default function OnboardingPage() {
    const router = useRouter();
    const [step, setStep] = useState(1);
    const [selectedInterests, setSelectedInterests] = useState<string[]>([]);
    const [selectedStyle, setSelectedStyle] = useState<string>('');
    const [isLoading, setIsLoading] = useState(false);

    // Check if user is authenticated
    useEffect(() => {
        const checkAuth = async () => {
            const token = localStorage.getItem('access_token');
            if (!token) {
                router.push('/login');
                return;
            }
            
            try {
                const res = await fetch('http://localhost:8080/user/profile', {
                    headers: {
                        'Authorization': `Bearer ${token}`,
                    },
                    credentials: 'include',
                });
                if (!res.ok) {
                    router.push('/login');
                }
            } catch (error) {
                router.push('/login');
            }
        };
        checkAuth();
    }, [router]);

    const toggleInterest = (id: string) => {
        setSelectedInterests((prev) =>
            prev.includes(id)
                ? prev.filter((i) => i !== id)
                : [...prev, id]
        );
    };

    const handleContinue = () => {
        if (step === 1 && selectedInterests.length > 0) {
            setStep(2);
        } else if (step === 2 && selectedStyle) {
            savePreferences();
        }
    };

    const savePreferences = async () => {
        setIsLoading(true);
        const token = localStorage.getItem('access_token');
        
        try {
            // Save interests
            for (const interestId of selectedInterests) {
                await fetch('http://localhost:8080/personalization/interest', {
                    method: 'POST',
                    headers: { 
                        'Content-Type': 'application/json',
                        'Authorization': `Bearer ${token}`,
                    },
                    credentials: 'include',
                    body: JSON.stringify({
                        interest: interests.find((i) => i.id === interestId)?.name,
                        category: interests.find((i) => i.id === interestId)?.category,
                    }),
                });
            }

            // Save learning style
            await fetch('http://localhost:8080/personalization/signal', {
                method: 'POST',
                headers: { 
                    'Content-Type': 'application/json',
                    'Authorization': `Bearer ${token}`,
                },
                credentials: 'include',
                body: JSON.stringify({
                    signal_type: 'learning_style_selected',
                    value: selectedStyle,
                }),
            });

            router.push('/dashboard');
        } catch (error) {
            console.error('Failed to save preferences:', error);
            // Still redirect to dashboard even if save fails
            router.push('/dashboard');
        } finally {
            setIsLoading(false);
        }
    };

    const handleSkip = () => {
        router.push('/dashboard');
    };

    return (
        <div className="min-h-screen bg-gradient-to-br from-slate-950 via-indigo-950 to-slate-950 flex items-center justify-center p-4">
            {/* Background Elements */}
            <div className="absolute inset-0 overflow-hidden pointer-events-none">
                <div className="absolute top-20 left-10 w-72 h-72 bg-indigo-500/20 rounded-full blur-3xl" />
                <div className="absolute bottom-20 right-10 w-96 h-96 bg-purple-500/20 rounded-full blur-3xl" />
            </div>

            <motion.div
                initial={{ opacity: 0, y: 30 }}
                animate={{ opacity: 1, y: 0 }}
                transition={{ duration: 0.6, ease: 'easeOut' }}
                className="w-full max-w-3xl relative z-10"
            >
                {/* Progress Bar */}
                <div className="mb-8">
                    <div className="flex items-center justify-between mb-4">
                        <span className="text-white/60 text-sm">Step {step} of 2</span>
                        <button
                            onClick={handleSkip}
                            className="text-white/40 hover:text-white/60 text-sm transition-colors"
                        >
                            Skip for now →
                        </button>
                    </div>
                    <div className="h-2 bg-white/10 rounded-full overflow-hidden">
                        <motion.div
                            className="h-full bg-gradient-to-r from-indigo-500 to-purple-500"
                            initial={{ width: '0%' }}
                            animate={{ width: step === 1 ? '50%' : '100%' }}
                            transition={{ duration: 0.3 }}
                        />
                    </div>
                </div>

                {/* Card */}
                <div className="bg-white/5 backdrop-blur-xl rounded-3xl border border-white/10 p-8 md:p-12">
                    <AnimatePresence mode="wait">
                        {step === 1 ? (
                            <motion.div
                                key="step1"
                                initial={{ opacity: 0, x: -20 }}
                                animate={{ opacity: 1, x: 0 }}
                                exit={{ opacity: 0, x: 20 }}
                                transition={{ duration: 0.3 }}
                            >
                                <div className="text-center mb-8">
                                    <div className="w-16 h-16 mx-auto mb-4 rounded-2xl bg-gradient-to-br from-indigo-500 to-purple-500 flex items-center justify-center text-3xl">
                                        🎯
                                    </div>
                                    <h1 className="text-3xl font-bold text-white mb-3">
                                        What are you interested in?
                                    </h1>
                                    <p className="text-white/60">
                                        Select topics you love. We'll use these to personalize your learning experience.
                                    </p>
                                </div>

                                <div className="grid grid-cols-2 md:grid-cols-3 lg:grid-cols-4 gap-3 mb-8">
                                    {interests.map((interest) => (
                                        <button
                                            key={interest.id}
                                            onClick={() => toggleInterest(interest.id)}
                                            className={`p-4 rounded-xl border-2 transition-all text-center ${
                                                selectedInterests.includes(interest.id)
                                                    ? 'border-indigo-500 bg-indigo-500/20'
                                                    : 'border-white/10 bg-white/5 hover:border-white/30 hover:bg-white/10'
                                            }`}
                                        >
                                            <span className="text-3xl mb-2 block">{interest.icon}</span>
                                            <span className={`text-sm font-medium ${
                                                selectedInterests.includes(interest.id) ? 'text-indigo-400' : 'text-white'
                                            }`}>
                                                {interest.name}
                                            </span>
                                        </button>
                                    ))}
                                </div>

                                <div className="flex items-center justify-between">
                                    <span className="text-white/50 text-sm">
                                        {selectedInterests.length} selected
                                    </span>
                                    <button
                                        onClick={handleContinue}
                                        disabled={selectedInterests.length === 0}
                                        className="px-8 py-3 bg-gradient-to-r from-indigo-500 to-purple-500 hover:from-indigo-600 hover:to-purple-600 disabled:opacity-50 disabled:cursor-not-allowed text-white font-medium rounded-xl transition-all flex items-center gap-2"
                                    >
                                        Continue
                                        <svg className="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                                            <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M9 5l7 7-7 7" />
                                        </svg>
                                    </button>
                                </div>
                            </motion.div>
                        ) : (
                            <motion.div
                                key="step2"
                                initial={{ opacity: 0, x: 20 }}
                                animate={{ opacity: 1, x: 0 }}
                                exit={{ opacity: 0, x: -20 }}
                                transition={{ duration: 0.3 }}
                            >
                                <div className="text-center mb-8">
                                    <div className="w-16 h-16 mx-auto mb-4 rounded-2xl bg-gradient-to-br from-purple-500 to-pink-500 flex items-center justify-center text-3xl">
                                        🧠
                                    </div>
                                    <h1 className="text-3xl font-bold text-white mb-3">
                                        How do you learn best?
                                    </h1>
                                    <p className="text-white/60">
                                        Choose your preferred learning style. We'll adapt our teaching method to match.
                                    </p>
                                </div>

                                <div className="grid grid-cols-1 md:grid-cols-2 gap-4 mb-8">
                                    {learningStyles.map((style) => (
                                        <button
                                            key={style.id}
                                            onClick={() => setSelectedStyle(style.id)}
                                            className={`p-6 rounded-xl border-2 transition-all text-left ${
                                                selectedStyle === style.id
                                                    ? 'border-purple-500 bg-purple-500/10'
                                                    : 'border-white/10 bg-white/5 hover:border-white/30 hover:bg-white/10'
                                            }`}
                                        >
                                            <div className="flex items-start gap-4">
                                                <span className="text-4xl">{style.icon}</span>
                                                <div>
                                                    <h3 className={`font-semibold text-lg mb-1 ${
                                                        selectedStyle === style.id ? 'text-purple-400' : 'text-white'
                                                    }`}>
                                                        {style.name}
                                                    </h3>
                                                    <p className="text-white/60 text-sm">{style.description}</p>
                                                </div>
                                            </div>
                                        </button>
                                    ))}
                                </div>

                                <div className="flex items-center justify-between">
                                    <button
                                        onClick={() => setStep(1)}
                                        className="px-6 py-3 text-white/60 hover:text-white font-medium rounded-xl transition-colors flex items-center gap-2"
                                    >
                                        <svg className="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                                            <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M15 19l-7-7 7-7" />
                                        </svg>
                                        Back
                                    </button>
                                    <button
                                        onClick={handleContinue}
                                        disabled={!selectedStyle || isLoading}
                                        className="px-8 py-3 bg-gradient-to-r from-indigo-500 to-purple-500 hover:from-indigo-600 hover:to-purple-600 disabled:opacity-50 disabled:cursor-not-allowed text-white font-medium rounded-xl transition-all flex items-center gap-2"
                                    >
                                        {isLoading ? (
                                            <>
                                                <svg className="animate-spin h-5 w-5" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24">
                                                    <circle className="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" strokeWidth="4"></circle>
                                                    <path className="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
                                                </svg>
                                                Saving...
                                            </>
                                        ) : (
                                            <>
                                                Get Started
                                                <svg className="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                                                    <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M5 13l4 4L19 7" />
                                                </svg>
                                            </>
                                        )}
                                    </button>
                                </div>
                            </motion.div>
                        )}
                    </AnimatePresence>
                </div>

                {/* Footer */}
                <div className="mt-8 text-center">
                    <Link href="/" className="text-white/40 hover:text-white/60 text-sm transition-colors">
                        ← Back to Home
                    </Link>
                </div>
            </motion.div>
        </div>
    );
}
