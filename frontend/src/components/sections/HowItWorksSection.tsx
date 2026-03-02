"use client";

import { motion } from "framer-motion";
import { SectionWrapper } from "@/components/ui/SectionWrapper";
import { FloatingIllustration } from "@/components/ui/FloatingIllustration";

const steps = [
  {
    number: "01",
    title: "Share Your Learning Style",
    description:
      "Tell us how you prefer to learn — visual animations, podcasts, storytelling, or concise summaries. The more we understand your preferences, the more personalized your learning experience becomes.",
    color: "from-brand-400 to-brand-600",
  },
  {
    number: "02",
    title: "Upload Your Material",
    description:
      "Upload any PDF, PPT, or article you want to study. Our system will analyze the structure, key concepts, and relationships within the content to prepare it for transformation.",
    color: "from-accent-cyan to-brand-400",
  },
  {
    number: "03",
    title: "Choose Your Learning Mode",
    description:
      "Select how you want the material to be delivered — animated video, audio explanation, interactive summary, or story-based learning. You control the format and depth.",
    color: "from-accent-pink to-brand-400",
  },
  {
    number: "04",
    title: "Transform & Learn",
    description:
      "The system converts your material into your chosen format, adapting explanations to match your preferences and desired level of detail.",
    color: "from-accent-amber to-accent-pink",
  },
  {
    number: "05",
    title: "Interact & Reinforce",
    description:
      "Ask questions, revisit sections, or test your understanding with quick quizzes. Learning becomes dynamic, not static.",
    color: "from-brand-500 to-accent-cyan",
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
          From Confused to{" "}
          <span className="text-gradient">Confident</span>
        </motion.h2>
        <motion.p
          initial={{ opacity: 0, y: 20 }}
          whileInView={{ opacity: 1, y: 0 }}
          viewport={{ once: true }}
          transition={{ delay: 0.2 }}
          className="mt-4 text-brand-200/70 max-w-2xl mx-auto text-lg"
        >
          Five steps to a learning experience that finally feels natural.
        </motion.p>
      </div>

      {/* Fixed width container for desktop */}
      <div className="relative max-w-6xl mx-auto">
        {/* Connector Line - Centered with absolute positioning */}
        <div
          className="absolute left-8 md:left-1/2 w-px bg-gradient-to-b from-brand-500/50 via-accent-cyan/50 to-accent-pink/50 hidden md:block -translate-x-1/2"
          style={{
            top: 0,
            bottom: 0,
            height: "100%",
          }}
        />

        {/* Mobile vertical line */}
        <div
          className="absolute left-8 w-px bg-gradient-to-b from-brand-500/50 via-accent-cyan/50 to-accent-pink/50 md:hidden"
          style={{
            top: 0,
            bottom: 0,
            height: "100%",
          }}
        />

        <div className="space-y-12 md:space-y-24">
          {steps.map((step, i) => (
            <motion.div
              key={step.number}
              initial={{ opacity: 0, x: i % 2 === 0 ? -40 : 40 }}
              whileInView={{ opacity: 1, x: 0 }}
              viewport={{ once: true, margin: "-80px" }}
              transition={{ duration: 0.6, delay: 0.1 }}
              className={`relative flex flex-col md:flex-row items-center ${
                i % 2 === 1 ? "md:flex-row-reverse" : ""
              }`}
            >
              {/* Left Content - Equal width with adjusted padding */}
              <div
                className={`w-full md:w-1/2 pl-16 md:pl-0 ${
                  i % 2 === 0
                    ? "md:pr-16 md:text-right"
                    : "md:pl-16 md:text-left"
                }`}
              >
                <div
                  className={`inline-flex flex-col space-y-3 ${
                    i % 2 === 0 ? "md:items-end" : "md:items-start"
                  }`}
                >
                  <span
                    className={`text-6xl font-black bg-gradient-to-r ${step.color} bg-clip-text text-transparent opacity-30`}
                  >
                    {step.number}
                  </span>
                  <h3 className="text-2xl font-bold text-white">
                    {step.title}
                  </h3>
                  <p className="text-brand-200/70 max-w-md leading-relaxed">
                    {step.description}
                  </p>
                </div>
              </div>

              {/* Center Dot - Absolutely centered */}
              <div className="absolute left-8 md:left-1/2 md:-translate-x-1/2 z-10 flex items-center justify-center">
                <div
                  className={`w-3 h-3 md:w-4 md:h-4 rounded-full bg-gradient-to-r ${step.color}`}
                />
                <div
                  className={`absolute w-6 h-6 md:w-8 md:h-8 rounded-full bg-gradient-to-r ${step.color} opacity-20 animate-pulse-glow`}
                />
              </div>

              {/* Right Content - Equal width, no padding modification */}
              <div
                className={`hidden md:block w-1/2 ${
                  i % 2 === 1 ? "md:pr-16" : "md:pl-16"
                }`}
              />
            </motion.div>
          ))}
        </div>
      </div>
    </SectionWrapper>
  );
}
