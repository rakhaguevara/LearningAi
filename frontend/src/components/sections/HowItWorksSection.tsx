"use client";

import { motion } from "framer-motion";
import { SectionWrapper } from "@/components/ui/SectionWrapper";
import { FloatingIllustration } from "@/components/ui/FloatingIllustration";

const steps = [
  {
    number: "01",
    title: "Share Your Interests",
    description:
      "Sign up and tell us what you're passionate about — anime, gaming, hip-hop, football, cooking, anything. The more we know, the better we teach.",
    color: "from-brand-400 to-brand-600",
  },
  {
    number: "02",
    title: "Pick a Topic to Learn",
    description:
      "Choose any subject from our library or type in your own question. From organic chemistry to macroeconomics, no topic is off limits.",
    color: "from-accent-cyan to-brand-400",
  },
  {
    number: "03",
    title: "AI Crafts Your Lesson",
    description:
      "Qwen AI analyzes your interests and builds a custom explanation with analogies, visuals, and examples drawn from the things you love most.",
    color: "from-accent-pink to-brand-400",
  },
  {
    number: "04",
    title: "Learn, Interact, Grow",
    description:
      "Ask follow-up questions, request illustrations, dive deeper. The AI adapts continuously — every interaction sharpens your personalized learning path.",
    color: "from-accent-amber to-accent-pink",
  },
];

export function HowItWorksSection() {
  return (
    <SectionWrapper id="how-it-works" className="relative overflow-hidden">
      <FloatingIllustration
        variant="orb"
        color="cyan"
        size="lg"
        className="absolute -right-20 top-20"
      />

      <div className="text-center mb-16">
        <motion.p
          initial={{ opacity: 0 }}
          whileInView={{ opacity: 1 }}
          viewport={{ once: true }}
          className="text-accent-pink text-sm font-semibold uppercase tracking-widest mb-4"
        >
          How It Works
        </motion.p>
        <motion.h2
          initial={{ opacity: 0, y: 20 }}
          whileInView={{ opacity: 1, y: 0 }}
          viewport={{ once: true }}
          transition={{ delay: 0.1 }}
          className="text-4xl md:text-5xl font-bold"
        >
          From Curious to{" "}
          <span className="text-gradient">Confident</span>
        </motion.h2>
        <motion.p
          initial={{ opacity: 0, y: 20 }}
          whileInView={{ opacity: 1, y: 0 }}
          viewport={{ once: true }}
          transition={{ delay: 0.2 }}
          className="mt-4 text-brand-200/70 max-w-2xl mx-auto text-lg"
        >
          Four steps to a learning experience that finally feels natural.
        </motion.p>
      </div>

      <div className="relative">
        {/* Connector Line */}
        <div className="absolute left-8 md:left-1/2 top-0 bottom-0 w-px bg-gradient-to-b from-brand-500/50 via-accent-cyan/50 to-accent-pink/50 hidden md:block" />

        <div className="space-y-12 md:space-y-24">
          {steps.map((step, i) => (
            <motion.div
              key={step.number}
              initial={{ opacity: 0, x: i % 2 === 0 ? -40 : 40 }}
              whileInView={{ opacity: 1, x: 0 }}
              viewport={{ once: true, margin: "-80px" }}
              transition={{ duration: 0.6, delay: 0.1 }}
              className={`relative flex flex-col md:flex-row items-center gap-8 ${
                i % 2 === 1 ? "md:flex-row-reverse" : ""
              }`}
            >
              {/* Content */}
              <div className="flex-1 md:text-right">
                <div
                  className={`inline-block ${
                    i % 2 === 1 ? "md:text-left" : "md:text-right"
                  }`}
                >
                  <span
                    className={`text-6xl font-black bg-gradient-to-r ${step.color} bg-clip-text text-transparent opacity-30`}
                  >
                    {step.number}
                  </span>
                  <h3 className="text-2xl font-bold text-white mt-2">
                    {step.title}
                  </h3>
                  <p className="mt-3 text-brand-200/70 max-w-md leading-relaxed">
                    {step.description}
                  </p>
                </div>
              </div>

              {/* Center Dot */}
              <div className="relative z-10 hidden md:flex items-center justify-center">
                <div
                  className={`w-4 h-4 rounded-full bg-gradient-to-r ${step.color}`}
                />
                <div
                  className={`absolute w-8 h-8 rounded-full bg-gradient-to-r ${step.color} opacity-20 animate-pulse-glow`}
                />
              </div>

              {/* Spacer for layout */}
              <div className="flex-1 hidden md:block" />
            </motion.div>
          ))}
        </div>
      </div>
    </SectionWrapper>
  );
}
