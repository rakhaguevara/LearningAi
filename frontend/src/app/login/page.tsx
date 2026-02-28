'use client';

import { motion } from 'framer-motion';
import { useState } from 'react';
import Link from 'next/link';
import LoginForm from '@/components/auth/LoginForm';
import RegisterForm from '@/components/auth/RegisterForm';
import { FloatingIllustration } from '@/components/ui/FloatingIllustration';

export default function LoginPage() {
    const [isLogin, setIsLogin] = useState(true);

    return (
        <div className="min-h-screen flex items-center justify-center p-4 relative overflow-hidden">
            {/* Background Orbs to match homepage theme seamlessly */}
            <FloatingIllustration
                variant="orb"
                color="purple"
                size="lg"
                className="absolute top-20 left-10 opacity-70"
            />
            <FloatingIllustration
                variant="orb"
                color="cyan"
                size="lg"
                className="absolute bottom-20 right-10 opacity-70"
            />
            <FloatingIllustration
                variant="ring"
                color="pink"
                size="md"
                className="absolute top-1/3 left-1/4 hidden lg:block opacity-40"
            />

            {/* Back to Home Navigation */}
            <Link
                href="/"
                className="absolute top-8 left-8 text-white/60 hover:text-white transition-colors flex items-center gap-2 z-20 font-medium"
            >
                <svg xmlns="http://www.w3.org/2000/svg" width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round"><path d="m15 18-6-6 6-6" /></svg>
                Back to Home
            </Link>

            <motion.div
                initial={{ opacity: 0, y: 20 }}
                animate={{ opacity: 1, y: 0 }}
                transition={{ duration: 0.5, ease: 'easeOut' }}
                className="w-full max-w-md p-8 rounded-3xl bg-white/5 backdrop-blur-md border border-white/10 shadow-2xl relative z-10"
            >
                <div className="text-center mb-8">
                    <h2 className="text-3xl font-bold bg-clip-text text-transparent bg-gradient-to-r from-indigo-500 to-purple-500">
                        {isLogin ? 'Welcome Back' : 'Join AI Learning'}
                    </h2>
                    <p className="text-white/70 mt-2 text-sm">
                        {isLogin
                            ? 'Sign in to jump straight back into your tailored learning journey.'
                            : 'Create an account to start your personalized learning experience.'}
                    </p>
                </div>

                <div className="relative">
                    {isLogin ? (
                        <LoginForm onSuccess={() => { }} />
                    ) : (
                        <RegisterForm onSuccess={() => setIsLogin(true)} />
                    )}
                </div>

                <div className="mt-6 text-center">
                    <button
                        onClick={() => setIsLogin(!isLogin)}
                        className="text-sm text-indigo-400 hover:text-purple-400 transition-colors"
                    >
                        {isLogin ? "Don't have an account? Sign up" : "Already have an account? Sign in"}
                    </button>
                </div>
            </motion.div>
        </div>
    );
}
