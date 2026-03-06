import type { Metadata } from 'next';
import { OnboardingPage } from './OnboardingPage';

export const metadata: Metadata = {
    title: 'Set Up Your Profile — Learny AI',
    description: 'Personalize your AI learning experience with a quick setup.',
};

export default function Page() {
    return <OnboardingPage />;
}
