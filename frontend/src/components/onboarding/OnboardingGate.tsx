'use client';

import { useEffect, useState } from 'react';
import { useRouter } from 'next/navigation';
import { getAuthToken, authHeaders, clearAuthToken } from '@/lib/auth';
import { siteConfig } from '@/lib/constants';

interface OnboardingGateProps {
    children: React.ReactNode;
}

/**
 * Wrap any protected page with <OnboardingGate> to enforce that the user
 * must complete onboarding before accessing it.
 *
 * If the user has no token → redirects to /login
 * If profile_completed is false → redirects to /onboarding
 * Otherwise → renders children normally
 */
export function OnboardingGate({ children }: OnboardingGateProps) {
    const router = useRouter();
    const [ready, setReady] = useState(false);

    useEffect(() => {
        async function check() {
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

                if (!res.ok) {
                    clearAuthToken();
                    router.replace('/login');
                    return;
                }

                const data = await res.json();
                const completed = data?.data?.profile_completed;

                if (!completed) {
                    // Still needs onboarding
                    router.replace('/onboarding');
                    return;
                }

                // All checks passed
                setReady(true);
            } catch {
                // Network error — allow through rather than blocking entirely
                setReady(true);
            }
        }

        check();
    }, [router]);

    if (!ready) {
        return (
            <div className="min-h-screen flex items-center justify-center">
                <div className="w-8 h-8 border-2 border-violet-500/30 border-t-violet-500 rounded-full animate-spin" />
            </div>
        );
    }

    return <>{children}</>;
}
