"use client";

import { useState } from "react";
import { motion, AnimatePresence } from "framer-motion";
import { SectionWrapper } from "@/components/ui/SectionWrapper";
import { GradientCard } from "@/components/ui/GradientCard";

const demos = [
  {
    interest: "Anime Fan",
    topic: "Quantum Superposition",
    explanation:
      "Think of Schrodinger's cat like Naruto's Shadow Clone Jutsu. Before you observe which clone is real, every single clone exists simultaneously as both the real Naruto and a copy. The moment someone lands a hit and checks — only then does reality 'collapse' into one definite Naruto. In quantum physics, particles exist in multiple states at once (superposition) until measured. Just like you can't tell which clone is the original until the jutsu breaks, you can't know a particle's state until you observe it. The act of observation itself forces reality to pick one outcome.",
    tags: ["Physics", "Naruto", "Adaptive"],
  },
  {
    interest: "Basketball Fan",
    topic: "Supply & Demand",
    explanation:
      "Imagine limited-edition LeBron sneakers just dropped. There are only 1,000 pairs (supply) but 50,000 people want them (demand). The price skyrockets — resellers charge 5x retail because scarcity creates value. Now imagine Nike releases 500,000 pairs. Suddenly everyone can get a pair and prices drop to retail. That's supply and demand in action: when something is scarce and desired, its value rises; when it's abundant, value falls. The same principle drives everything from concert ticket prices to housing markets.",
    tags: ["Economics", "NBA", "Adaptive"],
  },
  {
    interest: "Gamer",
    topic: "Neural Networks",
    explanation:
      "A neural network is like the skill tree in an RPG. Each node (neuron) in the tree connects to others, and as you invest experience points (training data), certain paths get stronger. When you train your character for combat, the strength and sword nodes get heavily invested — just like how a neural network strengthens connections between neurons that fire together. Backpropagation is like a game review after a boss fight: you look at what went wrong, trace back through your decisions, and redistribute points to perform better next time. After thousands of iterations, the network becomes an expert.",
    tags: ["Computer Science", "RPG", "Adaptive"],
  },
];

export function DemoSection() {
  const [active, setActive] = useState(0);

  return (
    <SectionWrapper id="demo">
      <div className="text-center mb-16">
        <motion.p
          initial={{ opacity: 0 }}
          whileInView={{ opacity: 1 }}
          viewport={{ once: true }}
          className="text-accent-amber text-sm font-semibold uppercase tracking-widest mb-4"
        >
          See It In Action
        </motion.p>
        <motion.h2
          initial={{ opacity: 0, y: 20 }}
          whileInView={{ opacity: 1, y: 0 }}
          viewport={{ once: true }}
          transition={{ delay: 0.1 }}
          className="text-4xl md:text-5xl font-bold"
        >
          One Topic,{" "}
          <span className="text-gradient">Three Worlds</span>
        </motion.h2>
        <motion.p
          initial={{ opacity: 0, y: 20 }}
          whileInView={{ opacity: 1, y: 0 }}
          viewport={{ once: true }}
          transition={{ delay: 0.2 }}
          className="mt-4 text-brand-200/70 max-w-2xl mx-auto text-lg"
        >
          The same complex idea, explained through completely different lenses.
          This is what personalized AI teaching looks like.
        </motion.p>
      </div>

      {/* Interest Tabs */}
      <div className="flex flex-wrap justify-center gap-3 mb-10">
        {demos.map((demo, i) => (
          <motion.button
            key={demo.interest}
            whileHover={{ scale: 1.05 }}
            whileTap={{ scale: 0.95 }}
            onClick={() => setActive(i)}
            className={`px-5 py-2.5 rounded-xl text-sm font-medium transition-all duration-300 cursor-pointer ${
              active === i
                ? "bg-gradient-to-r from-brand-500 to-accent-cyan text-white glow-purple"
                : "glass text-brand-200 hover:text-white"
            }`}
          >
            {demo.interest}
          </motion.button>
        ))}
      </div>

      {/* Demo Card */}
      <AnimatePresence mode="wait">
        <motion.div
          key={active}
          initial={{ opacity: 0, y: 20 }}
          animate={{ opacity: 1, y: 0 }}
          exit={{ opacity: 0, y: -20 }}
          transition={{ duration: 0.4 }}
        >
          <GradientCard
            className="max-w-4xl mx-auto p-8 lg:p-10"
            hover={false}
          >
            {/* Header */}
            <div className="flex flex-wrap items-center gap-3 mb-6">
              <span className="px-3 py-1 rounded-full bg-brand-500/20 text-brand-300 text-xs font-medium">
                {demos[active].interest}
              </span>
              <span className="text-brand-200/40">|</span>
              <span className="text-white font-semibold">
                Topic: {demos[active].topic}
              </span>
            </div>

            {/* AI Response Simulation */}
            <div className="glass rounded-xl p-6 mb-6">
              <div className="flex items-center gap-2 mb-4">
                <div className="h-6 w-6 rounded-full bg-gradient-to-br from-brand-400 to-accent-cyan flex items-center justify-center">
                  <span className="text-white text-xs font-bold">N</span>
                </div>
                <span className="text-sm font-medium text-brand-200">
                  Learny AI
                </span>
              </div>
              <p className="text-brand-200/90 leading-relaxed text-sm md:text-base">
                {demos[active].explanation}
              </p>
            </div>

            {/* Tags */}
            <div className="flex flex-wrap gap-2">
              {demos[active].tags.map((tag) => (
                <span
                  key={tag}
                  className="px-3 py-1 rounded-full bg-white/5 text-brand-200/60 text-xs"
                >
                  {tag}
                </span>
              ))}
            </div>
          </GradientCard>
        </motion.div>
      </AnimatePresence>
    </SectionWrapper>
  );
}
