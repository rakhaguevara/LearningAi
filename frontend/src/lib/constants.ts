export const siteConfig = {
  name: "NeuraLearn AI",
  description:
    "AI-powered learning that adapts to your world. Learn any subject through the lens of what you love.",
  url: "https://neuralearn.ai",
  api: process.env.NEXT_PUBLIC_API_URL || "http://localhost:8080",
  googleClientId: process.env.NEXT_PUBLIC_GOOGLE_CLIENT_ID || "",
} as const;

export const navLinks = [
  { label: "Features", href: "#features" },
  { label: "How It Works", href: "#how-it-works" },
  { label: "Demo", href: "#demo" },
] as const;
