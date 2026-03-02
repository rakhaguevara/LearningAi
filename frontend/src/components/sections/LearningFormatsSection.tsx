"use client";

import { motion } from "framer-motion";

const formats = [
  {
    id: 1,
    icon: "🎥",
    title: "Animated Video",
    description: "Visual storytelling with motion graphics.",
  },
  {
    id: 2,
    icon: "🎧",
    title: "Podcast Mode",
    description: "Audio-first learning for multitasking.",
  },
  {
    id: 3,
    icon: "📖",
    title: "Storytelling",
    description: "Narrative-driven explanations.",
  },
  {
    id: 4,
    icon: "📝",
    title: "Smart Summary",
    description: "Concise key points for quick review.",
  },
];

export function LearningFormatsSection() {
  return (
<<<<<<< HEAD
    <section className="relative py-12 md:py-16 overflow-hidden">
=======
    <section className="relative py-12 md:py-16 bg-slate-900/30 overflow-hidden">
>>>>>>> b0df5b7113a5b69f8fed22603b16d2d21641818d
      {/* Subtle background gradient */}
      <div className="absolute inset-0 bg-gradient-to-b from-transparent via-slate-800/20 to-transparent pointer-events-none" />

      <div className="relative max-w-5xl mx-auto px-6">
        {/* Header */}
        <div className="text-center mb-8 md:mb-10">
          <motion.p
            initial={{ opacity: 0 }}
            whileInView={{ opacity: 1 }}
            viewport={{ once: true }}
            className="text-accent-cyan text-xs font-semibold uppercase tracking-widest mb-3"
          >
            Adaptive Formats
          </motion.p>
          <motion.h2
            initial={{ opacity: 0, y: 20 }}
            whileInView={{ opacity: 1, y: 0 }}
            viewport={{ once: true }}
            transition={{ delay: 0.1 }}
            className="text-2xl md:text-3xl font-bold"
          >
            One Material{" "}
            <span className="text-gradient">Multiple Ways.</span>
          </motion.h2>
          <motion.p
            initial={{ opacity: 0, y: 20 }}
            whileInView={{ opacity: 1, y: 0 }}
            viewport={{ once: true }}
            transition={{ delay: 0.2 }}
            className="mt-2 text-brand-200/70 max-w-xl mx-auto text-sm"
          >
            Upload once. Experience it your way.
          </motion.p>
        </div>

        {/* Format Cards Grid */}
        <div className="grid grid-cols-2 md:grid-cols-4 gap-4">
          {formats.map((format, index) => (
            <motion.div
              key={format.id}
              initial={{ opacity: 0, y: 20 }}
              whileInView={{ opacity: 1, y: 0 }}
              viewport={{ once: true, margin: "-50px" }}
              transition={{ duration: 0.4, delay: index * 0.08 }}
              className="group"
            >
              <div className="h-full bg-slate-800/40 backdrop-blur-md border border-white/10 rounded-xl shadow-sm p-4 transition-all duration-300 hover:scale-105 hover:border-white/20 hover:shadow-md hover:bg-slate-800/50">
                {/* Icon */}
                <div className="text-2xl mb-2 opacity-60 group-hover:opacity-80 transition-opacity duration-300">
                  {format.icon}
                </div>

                {/* Title */}
                <h3 className="text-sm font-semibold text-white mb-1">
                  {format.title}
                </h3>

                {/* Description */}
                <p className="text-xs text-brand-200/70 leading-relaxed">
                  {format.description}
                </p>
              </div>
            </motion.div>
          ))}
        </div>
      </div>
    </section>
  );
}
