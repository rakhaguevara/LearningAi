"use client";

import { motion } from "framer-motion";
import { SectionWrapper } from "@/components/ui/SectionWrapper";
import { AnimatedButton } from "@/components/ui/AnimatedButton";
import { FloatingIllustration } from "@/components/ui/FloatingIllustration";
import { useRouter } from "next/navigation";

export function CTASection() {
  const router = useRouter();

  return (
    <SectionWrapper className="relative overflow-hidden">
      <FloatingIllustration
        variant="orb"
        color="purple"
        size="lg"
        className="absolute -left-20 top-0"
      />
      <FloatingIllustration
        variant="orb"
        color="pink"
        size="md"
        className="absolute -right-10 bottom-0"
      />
      <FloatingIllustration
        variant="ring"
        color="cyan"
        size="md"
        className="absolute right-1/4 top-10 hidden lg:block"
      />

      <div className="relative z-10">
        <motion.div
          initial={{ opacity: 0, scale: 0.95 }}
          whileInView={{ opacity: 1, scale: 1 }}
          viewport={{ once: true }}
          transition={{ duration: 0.6 }}
          className="glass-strong rounded-3xl p-10 md:p-16 text-center max-w-4xl mx-auto"
        >
          <motion.h2
            initial={{ opacity: 0, y: 20 }}
            whileInView={{ opacity: 1, y: 0 }}
            viewport={{ once: true }}
            transition={{ delay: 0.1 }}
            className="text-4xl md:text-5xl font-bold"
          >
            The Future of Learning
            <br />
            <span className="text-gradient">Is Personal</span>
          </motion.h2>

          <motion.p
            initial={{ opacity: 0, y: 20 }}
            whileInView={{ opacity: 1, y: 0 }}
            viewport={{ once: true }}
            transition={{ delay: 0.2 }}
            className="mt-6 text-lg text-brand-200/80 max-w-xl mx-auto leading-relaxed"
          >
            Stop forcing yourself into one-size-fits-all education. Join
            thousands of learners who discovered that the fastest path to
            understanding runs straight through their passions.
          </motion.p>

          <motion.div
            initial={{ opacity: 0, y: 20 }}
            whileInView={{ opacity: 1, y: 0 }}
            viewport={{ once: true }}
            transition={{ delay: 0.3 }}
            className="mt-10 flex flex-col sm:flex-row gap-4 justify-center"
          >
            <AnimatedButton onClick={() => router.push('/login')} variant="primary" size="lg">
              Create Free Account
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
              Talk to Our Team
            </AnimatedButton>
          </motion.div>

          <motion.p
            initial={{ opacity: 0 }}
            whileInView={{ opacity: 1 }}
            viewport={{ once: true }}
            transition={{ delay: 0.5 }}
            className="mt-6 text-xs text-brand-200/40"
          >
            No credit card required. Free tier includes 50 AI-powered lessons
            per month.
          </motion.p>
        </motion.div>
      </div>
    </SectionWrapper>
  );
}
