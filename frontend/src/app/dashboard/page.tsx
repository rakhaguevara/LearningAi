import type { Metadata } from 'next';
import { DashboardLayout } from '@/components/dashboard/DashboardLayout';
import { OnboardingGate } from '@/components/onboarding/OnboardingGate';

export const metadata: Metadata = {
    title: 'Dashboard — NeuraLearn AI',
    description: 'Your AI-powered learning workspace.',
};

export default function DashboardPage() {
    return (
        <OnboardingGate>
            <DashboardLayout />
        </OnboardingGate>
    );
}
