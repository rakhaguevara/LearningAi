"use client";

import { motion } from "framer-motion";
import { AnimatedButton } from "@/components/ui/AnimatedButton";
import { FloatingIllustration } from "@/components/ui/FloatingIllustration";
import { FloatingLearningIcons } from "@/components/ui/FloatingLearningIcons";
import { useRouter } from "next/navigation";

export function HeroSection() {
  const router = useRouter();

  return (
    <section className="relative min-h-screen flex items-center justify-center overflow-hidden pt-16">
      {/* Background Orbs */}
      <FloatingIllustration
        variant="orb"
        color="purple"
        size="lg"
        className="absolute top-20 left-10"
      />
      <FloatingIllustration
        variant="orb"
        color="cyan"
        size="lg"
        className="absolute bottom-20 right-10"
      />
      <FloatingIllustration
        variant="orb"
        color="pink"
        size="md"
        className="absolute top-40 right-1/4"
      />
      <FloatingIllustration
        variant="ring"
        color="purple"
        size="lg"
        className="absolute top-1/3 left-1/4 hidden lg:block"
      />

      {/* Floating Learning Icons */}
      <FloatingLearningIcons />

      <div className="relative z-10 mx-auto max-w-7xl px-6 lg:px-8 text-center">
        {/* Badge */}
        <motion.div
          initial={{ opacity: 0, y: 20 }}
          animate={{ opacity: 1, y: 0 }}
          transition={{ duration: 0.5 }}
          className="mt-8 md:mt-12"
        >
          <span className="inline-flex items-center gap-2 glass rounded-full px-4 py-1.5 text-sm text-brand-200">
            <span className="h-2 w-2 rounded-full bg-accent-cyan animate-pulse" />
            Powered by Qwen AI &amp; Alibaba Cloud
          </span>
        </motion.div>

        {/* Headline */}
        <motion.h1
          initial={{ opacity: 0, y: 30 }}
          animate={{ opacity: 1, y: 0 }}
          transition={{ duration: 0.6, delay: 0.1 }}
          className="mt-8 text-5xl md:text-7xl lg:text-8xl font-bold tracking-tight leading-[1.08]"
        >
          Learn Anything
          <br />
          <span className="text-gradient">Through What You Love</span>
        </motion.h1>

        {/* Subheadline */}
        <motion.p
          initial={{ opacity: 0, y: 30 }}
          animate={{ opacity: 1, y: 0 }}
          transition={{ duration: 0.6, delay: 0.2 }}
          className="mt-6 text-lg md:text-xl text-brand-200/80 max-w-2xl mx-auto leading-relaxed"
        >
          NeuraLearn transforms complex knowledge into adaptive learning experiences.
          intelligently tailored to how you process, understand, and remember.
        </motion.p>

        {/* CTA Buttons */}
        <motion.div
          initial={{ opacity: 0, y: 30 }}
          animate={{ opacity: 1, y: 0 }}
          transition={{ duration: 0.6, delay: 0.3 }}
          className="mt-10 flex flex-col sm:flex-row gap-4 justify-center"
        >
          <AnimatedButton onClick={() => router.push('/login')} variant="primary" size="lg">
            Start Learning Free
            <svg
              width="20"
              height="20"
              viewBox="0 0 24 24"
              fill="none"
              stroke="currentColor"
              strokeWidth="2"
              strokeLinecap="round"
              strokeLinejoin="round"
            >
              <path d="M5 12h14M12 5l7 7-7 7" />
            </svg>
          </AnimatedButton>
          <AnimatedButton variant="secondary" size="lg">
            Watch Demo
          </AnimatedButton>
        </motion.div>

        {/* Social proof */}
        <motion.div
          initial={{ opacity: 0 }}
          animate={{ opacity: 1 }}
          transition={{ duration: 0.8, delay: 0.6 }}
          className="mt-16 flex flex-col items-center gap-3"
        >
          <div className="flex -space-x-3">
            {[...Array(5)].map((_, i) => (
              <div
                key={i}
                className="w-10 h-10 rounded-full border-2 border-[var(--color-bg-primary)] bg-gradient-to-br from-brand-400 to-accent-cyan"
              />
            ))}
          </div>
          <p className="text-sm text-brand-200/60">
            Used by <span className="text-white font-medium">12,000+</span>{" "}
            learners worldwide
          </p>
        </motion.div>
      </div>
    </section>
  );
}
