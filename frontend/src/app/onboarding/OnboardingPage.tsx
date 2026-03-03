'use client';

import { useEffect, useState } from 'react';
import { motion } from 'framer-motion';
import { useRouter } from 'next/navigation';
import { OnboardingWizard } from '@/components/onboarding/OnboardingWizard';
import { getAuthToken, authHeaders } from '@/lib/auth';
import { siteConfig } from '@/lib/constants';

export function OnboardingPage() {
    const router = useRouter();
    const [checking, setChecking] = useState(true);

    // Guard: if already completed, skip to dashboard
    useEffect(() => {
        async function checkStatus() {
            const token = getAuthToken();

            if (!token) {
                router.replace('/login');
                return;
            }
            try {
                const res = await fetch(`${siteConfig.api}/onboarding/status`, {
                    headers: { ...authHeaders() },
                    credentials: 'include',
                });
                if (res.ok) {
                    const data = await res.json();
                    if (data?.data?.profile_completed) {
                        router.replace('/dashboard');
                        return;
                    }
                }
            } catch {
                // If status check fails, proceed to show the onboarding (better than blocking)
            } finally {
                setChecking(false);
            }
        }
        checkStatus();
    }, [router]);

    if (checking) {
        return (
            <div className="min-h-screen flex items-center justify-center">
                <div className="w-8 h-8 border-2 border-violet-500/30 border-t-violet-500 rounded-full animate-spin" />
            </div>
        );
    }

    return (
        <div className="min-h-screen flex flex-col items-center justify-center relative overflow-hidden p-4">
            {/* Animated gradient background orbs */}
            <div className="absolute inset-0 pointer-events-none">
                <div className="absolute top-[-10%] left-[-5%] w-96 h-96 rounded-full bg-violet-600/20 blur-3xl animate-pulse" />
                <div
                    className="absolute bottom-[-10%] right-[-5%] w-96 h-96 rounded-full bg-fuchsia-600/20 blur-3xl animate-pulse"
                    style={{ animationDelay: '1.5s' }}
                />
                <div
                    className="absolute top-1/2 left-1/2 -translate-x-1/2 -translate-y-1/2 w-64 h-64 rounded-full bg-purple-500/10 blur-2xl animate-pulse"
                    style={{ animationDelay: '0.75s' }}
                />
            </div>

            {/* Header */}
            <motion.div
                initial={{ opacity: 0, y: -20 }}
                animate={{ opacity: 1, y: 0 }}
                transition={{ duration: 0.5 }}
                className="text-center mb-10 relative z-10"
            >
                <div className="inline-flex items-center gap-2 px-4 py-2 rounded-full bg-white/5 border border-white/10 text-white/60 text-xs font-medium mb-4">
                    <span className="w-1.5 h-1.5 rounded-full bg-violet-400 animate-pulse" />
                    Setting up your learning profile
                </div>
                <h1 className="text-3xl md:text-4xl font-bold text-white">
                    Let&rsquo;s personalize your{' '}
                    <span className="bg-clip-text text-transparent bg-gradient-to-r from-violet-400 to-fuchsia-400">
                        AI tutor
                    </span>
                </h1>
                <p className="mt-3 text-white/50 text-sm max-w-sm mx-auto">
                    10 quick questions. Takes about 2 minutes. You can always change these later.
                </p>
            </motion.div>

            {/* Wizard */}
            <motion.div
                initial={{ opacity: 0, y: 20 }}
                animate={{ opacity: 1, y: 0 }}
                transition={{ duration: 0.5, delay: 0.15 }}
                className="w-full relative z-10"
            >
                <OnboardingWizard />
            </motion.div>
        </div>
    );
}
