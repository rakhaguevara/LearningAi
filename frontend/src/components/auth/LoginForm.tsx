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
                window.location.href = '/dashboard';
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
            
            // 3. Tangani jika status bukan 2xx (termasuk 401 Unauthorized)
            if (!res.ok) {
                throw new Error(data.message || data.error || "Google login gagal di verifikasi server");
            }

            // 4. Pastikan token tersedia (Go backend mengembalikan data.data.access_token)
            const tokenToStore = data.token || data?.data?.access_token;
            if (!tokenToStore) {
                 throw new Error("Sistem tidak mengembalikan token autentikasi.");
            }

            // 5. Simpan token ke localStorage (atau cookie) via utility
            storeAuthToken(tokenToStore); // Fungsi utilitas dari codebase
            
            console.log("Login Berhasil!", data.user || data.data);
            
            // 6. Jalankan callback dan Pindahkan User ke /dashboard
            onSuccess?.();
            
            // Menggunakan window.location.href untuk memastikan force reload 
            // sehingga state Auth di root layout membaca localStorage yang baru
            window.location.href = '/dashboard';
            
        } catch (err: any) {
            console.error("Error verifikasi login:", err);
            // Menampilkan error message ke layar pengguna
            setError(err.message || "Terjadi kesalahan saat verifikasi Google Token.");
        } finally {
            setLoading(false);
        }
    };

    return (
        <div className="w-full max-w-md p-6 bg-gray-900 rounded-xl shadow-md">
            {/* Form Area */}
            <form onSubmit={handleSubmit} className="space-y-4">
                {!isLogin && (
                    <input
                        type="text"
                        placeholder="Full Name"
                        value={name}
                        onChange={(e) => setName(e.target.value)}
                        className="w-full p-3 rounded-xl bg-gray-800 border border-gray-700 text-white outline-none"
                        required={!isLogin}
                    />
                )}
                <input
                    type="email"
                    placeholder="Email"
                    value={email}
                    onChange={(e) => setEmail(e.target.value)}
                    className="w-full p-3 rounded-xl bg-gray-800 border border-gray-700 text-white outline-none"
                    required
                />
                <input
                    type="password"
                    placeholder="Password"
                    value={password}
                    onChange={(e) => setPassword(e.target.value)}
                    className="w-full p-3 rounded-xl bg-gray-800 border border-gray-700 text-white outline-none"
                    required
                />
                {error && <p className="text-red-400 text-sm">{error}</p>}
                <button
                    type="submit"
                    disabled={loading}
                    className="w-full py-3 rounded-xl bg-gradient-to-r from-indigo-500 to-purple-500 text-white font-semibold hover:opacity-90 transition"
                >
                    {loading ? "Processing..." : (isLogin ? "Sign In" : "Register")}
                </button>
            </form>

            {/* Toggle Link */}
            <div className="mt-4 text-center">
                <button
                    type="button"
                    onClick={() => setIsLogin(!isLogin)}
                    className="text-sm text-blue-400 hover:text-purple-400 transition-colors"
                >
                    {isLogin ? "Need an account? Sign up" : "Already have an account? Sign in"}
                </button>
            </div>

            {/* Divider */}
            <div className="flex items-center gap-4 my-6">
                <div className="flex-1 h-px bg-gray-700"></div>
                <span className="text-xs text-gray-400">OR</span>
                <div className="flex-1 h-px bg-gray-700"></div>
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