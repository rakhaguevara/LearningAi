'use client';

import { useState } from "react";
import { useRouter } from "next/navigation";
import { GoogleLogin } from "@react-oauth/google";
import { siteConfig } from "@/lib/constants";
import { storeAuthToken } from "@/lib/auth";

type Props = {
    onSuccess?: () => void;
};

export default function LoginForm({ onSuccess }: Props) {
    const router = useRouter();
    const [isLogin, setIsLogin] = useState(true);
    const [name, setName] = useState("");
    const [email, setEmail] = useState("");
    const [password, setPassword] = useState("");
    const [loading, setLoading] = useState(false);
    const [error, setError] = useState("");

    const handleSubmit = async (e: React.FormEvent) => {
        e.preventDefault();
        setLoading(true);
        setError("");

        try {
            const endpoint = isLogin ? "/auth/login" : "/auth/register";
            const body = isLogin
                ? JSON.stringify({ email, password })
                : JSON.stringify({ name, email, password });

            const res = await fetch(`${siteConfig.api}${endpoint}`, {
                method: "POST",
                headers: { "Content-Type": "application/json" },
                body
            });

            const data = await res.json();

            if (!res.ok) throw new Error(data.message || (isLogin ? "Login gagal" : "Register gagal"));

            const validToken = data.token || data?.data?.access_token;
            if (isLogin) {
                storeAuthToken(validToken);
                onSuccess?.();
                router.push("/dashboard");
            } else {
                storeAuthToken(validToken);
                window.location.href = '/onboarding';
            }
        } catch (err: any) {
            setError(err.message);
        }

        setLoading(false);
    };

    const handleGoogleSuccess = async (credentialResponse: any) => {
        setError(""); // Reset error message
        setLoading(true); // Tampilkan status loading

        try {
            // 1. Mengirim token ke backend GoLang API Route
            const res = await fetch(`${siteConfig.api}/auth/google`, {
                method: "POST",
                headers: { "Content-Type": "application/json" },
                body: JSON.stringify({ token: credentialResponse.credential })
            });

            // 2. Parse hasil JSON dari API Route
            const data = await res.json();
            if (!res.ok) throw new Error("Google login gagal");

            storeAuthToken(data.token || data?.data?.access_token);
            router.push("/dashboard");
        } catch (err) {
            console.error(err);
            setError("Google login gagal");
        }
    };

    return (
        <div className="w-full max-w-md glass glass-strong rounded-2xl p-8 backdrop-blur-2xl border border-white/10">
            {/* Header */}
            <div className="mb-8">
                <h2 className="text-2xl font-bold text-white mb-2">
                    {isLogin ? "Welcome Back" : "Join Learny"}
                </h2>
                <p className="text-sm text-white/60">
                    {isLogin 
                        ? "Sign in to continue your learning journey" 
                        : "Start your personalized learning experience"}
                </p>
            </div>

            {/* Form Area */}
            <form onSubmit={handleSubmit} className="space-y-4">
                {!isLogin && (
                    <div className="space-y-2">
                        <label className="block text-xs font-medium text-white/70 px-1">Full Name</label>
                        <input
                            type="text"
                            placeholder="John Doe"
                            value={name}
                            onChange={(e) => setName(e.target.value)}
                            className="w-full px-4 py-3 rounded-lg bg-white/5 border border-white/10 text-white placeholder-white/40 outline-none focus:border-accent-cyan/50 focus:ring-1 focus:ring-accent-cyan/20 transition-all"
                            required={!isLogin}
                        />
                    </div>
                )}
                <div className="space-y-2">
                    <label className="block text-xs font-medium text-white/70 px-1">Email Address</label>
                    <input
                        type="email"
                        placeholder="you@example.com"
                        value={email}
                        onChange={(e) => setEmail(e.target.value)}
                        className="w-full px-4 py-3 rounded-lg bg-white/5 border border-white/10 text-white placeholder-white/40 outline-none focus:border-accent-cyan/50 focus:ring-1 focus:ring-accent-cyan/20 transition-all"
                        required
                    />
                </div>
                <div className="space-y-2">
                    <label className="block text-xs font-medium text-white/70 px-1">Password</label>
                    <input
                        type="password"
                        placeholder="••••••••"
                        value={password}
                        onChange={(e) => setPassword(e.target.value)}
                        className="w-full px-4 py-3 rounded-lg bg-white/5 border border-white/10 text-white placeholder-white/40 outline-none focus:border-accent-cyan/50 focus:ring-1 focus:ring-accent-cyan/20 transition-all"
                        required
                    />
                </div>
                
                {error && (
                    <div className="flex gap-2 p-3 rounded-lg bg-red-500/10 border border-red-500/20">
                        <span className="text-red-400 text-sm">{error}</span>
                    </div>
                )}
                
                <button
                    type="submit"
                    disabled={loading}
                    className="w-full py-3 rounded-lg bg-gradient-to-r from-brand-500 to-accent-cyan text-white font-semibold hover:shadow-lg hover:shadow-brand-500/25 transition-all disabled:opacity-50 disabled:cursor-not-allowed mt-6"
                >
                    {loading ? (
                        <span className="flex items-center justify-center gap-2">
                            <span className="w-4 h-4 border-2 border-white/20 border-t-white rounded-full animate-spin"></span>
                            Processing...
                        </span>
                    ) : (isLogin ? "Sign In" : "Create Account")}
                </button>
            </form>

            {/* Toggle Link */}
            <div className="mt-6 text-center">
                <button
                    type="button"
                    onClick={() => setIsLogin(!isLogin)}
                    className="text-sm text-white/60 hover:text-white transition-colors"
                >
                    {isLogin ? (
                        <>Don't have an account? <span className="text-accent-cyan font-semibold">Sign up</span></>
                    ) : (
                        <>Already have an account? <span className="text-accent-cyan font-semibold">Sign in</span></>
                    )}
                </button>
            </div>

            {/* Divider */}
            <div className="flex items-center gap-4 my-6">
                <div className="flex-1 h-px bg-gradient-to-r from-transparent via-white/10 to-transparent"></div>
                <span className="text-xs text-white/40 font-medium">OR</span>
                <div className="flex-1 h-px bg-gradient-to-r from-transparent via-white/10 to-transparent"></div>
            </div>

            {/* Google Login */}
            <div className="flex justify-center">
                <GoogleLogin
                    onSuccess={handleGoogleSuccess}
                    onError={() => setError("Google login gagal")}
                />
            </div>
        </div>
    );
}