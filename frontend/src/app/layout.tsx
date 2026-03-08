import type { Metadata } from "next";
import { ThemeProvider } from "@/lib/ThemeContext";
import { LanguageProvider } from "@/lib/LanguageContext";
import { GoogleOAuthProvider } from "@react-oauth/google";
import "@/styles/globals.css";

export const metadata: Metadata = {
  title: "Learny AI — Learn Any Subject Through What You Love",
  description:
    "AI-powered learning platform that personalizes teaching using your interests. Anime fan? Sports enthusiast? Gamer? We adapt every explanation to your world.",
  keywords: [
    "AI learning",
    "personalized education",
    "adaptive learning",
    "Qwen AI",
    "online tutor",
  ],
};

export default function RootLayout({
  children,
}: {
  children: React.ReactNode;
}) {
  return (
    <html lang="en" className="antialiased">
      <body className="min-h-screen">
        <GoogleOAuthProvider clientId={process.env.NEXT_PUBLIC_GOOGLE_CLIENT_ID!}>
          <ThemeProvider>
            <LanguageProvider>
              {children}
            </LanguageProvider>
          </ThemeProvider>
        </GoogleOAuthProvider>
      </body>
    </html>
  );
}