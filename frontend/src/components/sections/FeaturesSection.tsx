"use client";

import { motion } from "framer-motion";
import { SectionWrapper } from "@/components/ui/SectionWrapper";
import { GradientCard } from "@/components/ui/GradientCard";
import { FloatingIllustration } from "@/components/ui/FloatingIllustration";

const features = [
  {
    icon: (
      <svg width="28" height="28" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="1.5" strokeLinecap="round" strokeLinejoin="round">
        <path d="M12 2a5 5 0 015 5c0 3-2 5.5-5 8.5C9 12.5 7 10 7 7a5 5 0 015-5z" />
        <circle cx="12" cy="7" r="1.5" />
        <path d="M5 20h14" />
      </svg>
    ),
    title: "Interest-Mapped Learning",
    description:
      "Tell us you love anime, and we'll explain quantum physics through Naruto's chakra system. Your passions become the bridge to new knowledge.",
    glow: "purple" as const,
  },
  {
    icon: (
      <svg width="28" height="28" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="1.5" strokeLinecap="round" strokeLinejoin="round">
        <path d="M21 15a2 2 0 01-2 2H7l-4 4V5a2 2 0 012-2h14a2 2 0 012 2z" />
        <path d="M8 10h.01M12 10h.01M16 10h.01" />
      </svg>
    ),
    title: "Adaptive Conversations",
    description:
      "Our AI doesn't just answer questions — it reads how you're progressing and shifts its teaching approach in real time, like a master tutor.",
    glow: "cyan" as const,
  },
  {
    icon: (
      <svg width="28" height="28" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="1.5" strokeLinecap="round" strokeLinejoin="round">
        <rect x="3" y="3" width="18" height="18" rx="2" />
        <circle cx="8.5" cy="8.5" r="1.5" />
        <path d="M21 15l-5-5L5 21" />
      </svg>
    ),
    title: "AI-Generated Visuals",
    description:
      "Complex concepts get custom illustrations generated on the fly — styled to match your interests, making abstract ideas tangible and memorable.",
    glow: "pink" as const,
  },
  {
    icon: (
      <svg width="28" height="28" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="1.5" strokeLinecap="round" strokeLinejoin="round">
        <path d="M2 3h6a4 4 0 014 4v14a3 3 0 00-3-3H2z" />
        <path d="M22 3h-6a4 4 0 00-4 4v14a3 3 0 013-3h7z" />
      </svg>
    ),
    title: "Any Subject, Any Level",
    description:
      "From calculus to world history, beginner to advanced — the platform scales its depth and complexity to exactly where you are in your journey.",
    glow: "purple" as const,
  },
  {
    icon: (
      <svg width="28" height="28" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="1.5" strokeLinecap="round" strokeLinejoin="round">
        <polyline points="22 12 18 12 15 21 9 3 6 12 2 12" />
      </svg>
    ),
    title: "Progress Intelligence",
    description:
      "Track your learning trajectory with AI-powered analytics that highlight strengths, pinpoint gaps, and suggest your optimal next step.",
    glow: "cyan" as const,
  },
  {
    icon: (
      <svg width="28" height="28" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="1.5" strokeLinecap="round" strokeLinejoin="round">
        <path d="M12 22s8-4 8-10V5l-8-3-8 3v7c0 6 8 10 8 10z" />
        <path d="M9 12l2 2 4-4" />
      </svg>
    ),
    title: "Private & Secure",
    description:
      "Your learning data stays yours. Enterprise-grade encryption, GDPR-ready architecture, and full data sovereignty powered by Alibaba Cloud.",
    glow: "pink" as const,
  },
];

const containerVariants = {
  hidden: {},
  visible: {
    transition: { staggerChildren: 0.1 },
  },
};

export function FeaturesSection() {
  return (
    <SectionWrapper id="features" className="relative overflow-visible">
      {/* Decorative background element */}
      <FloatingIllustration
        variant="ring"
        color="purple"
        size="lg"
        className="absolute -left-32 top-10"
      />
      <div className="text-center mb-16 relative z-10">
        <motion.p
          initial={{ opacity: 0 }}
          whileInView={{ opacity: 1 }}
          viewport={{ once: true }}
          className="text-accent-cyan text-sm font-semibold uppercase tracking-widest mb-4"
        >
          Why Learny?
        </motion.p>
        <motion.h2
          initial={{ opacity: 0, y: 20 }}
          whileInView={{ opacity: 1, y: 0 }}
          viewport={{ once: true }}
          transition={{ delay: 0.1 }}
          className="text-4xl md:text-5xl font-bold"
        >
          Learning That{" "}
          <span className="text-gradient">Make You Enjoy</span>
        </motion.h2>
        <motion.p
          initial={{ opacity: 0, y: 20 }}
          whileInView={{ opacity: 1, y: 0 }}
          viewport={{ once: true }}
          transition={{ delay: 0.2 }}
          className="mt-4 text-brand-200/70 max-w-2xl mx-auto text-lg"
        >
          Every feature is designed around one principle: education should
          adapt to the learner, never the other way around.
        </motion.p>
      </div>

      <motion.div
        variants={containerVariants}
        initial="hidden"
        whileInView="visible"
        viewport={{ once: true, margin: "-50px" }}
        className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6"
      >
        {features.map((feature) => (
          <GradientCard key={feature.title} glowColor={feature.glow}>
            <div className="h-12 w-12 rounded-xl bg-white/5 flex items-center justify-center text-brand-300 mb-4">
              {feature.icon}
            </div>
            <h3 className="text-lg font-semibold text-white mb-2">
              {feature.title}
            </h3>
            <p className="text-sm text-brand-200/70 leading-relaxed">
              {feature.description}
            </p>
          </GradientCard>
        ))}
      </motion.div>
    </SectionWrapper>
  );
}
