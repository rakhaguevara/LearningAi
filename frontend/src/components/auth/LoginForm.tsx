'use client';

import { useState } from "react";
import { useRouter } from "next/navigation";
import { GoogleLogin } from "@react-oauth/google";

type Props = {
  onSuccess?: () => void;
};

export default function LoginForm({ onSuccess }: Props) {
  const router = useRouter();
  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState("");

  // LOGIN EMAIL & PASSWORD
  const handleLogin = async (e: React.FormEvent) => {
    e.preventDefault();
    setLoading(true);
    setError("");

    try {
      const res = await fetch("http://localhost:8080/auth/login", {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({ email, password })
      });

      const data = await res.json();

      if (!res.ok) throw new Error(data.message || "Login gagal");

      localStorage.setItem("token", data.token);
      onSuccess?.();
      router.push("/dashboard");
    } catch (err: any) {
      setError(err.message);
    }

    setLoading(false);
  };

  // LOGIN GOOGLE
  const handleGoogleSuccess = async (credentialResponse: any) => {
    try {
      const res = await fetch("http://localhost:8080/auth/google", {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({ token: credentialResponse.credential })
      });

      const data = await res.json();
      if (!res.ok) throw new Error("Google login gagal");

      localStorage.setItem("token", data.token);
      router.push("/dashboard");
    } catch (err) {
      console.error(err);
      setError("Google login gagal");
    }
  };

  return (
    <div className="w-full max-w-md p-6 bg-gray-900 rounded-xl shadow-md">
      {/* Form Email & Password */}
      <form onSubmit={handleLogin} className="space-y-4">
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
          {loading ? "Signing in..." : "Sign In"}
        </button>
      </form>

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