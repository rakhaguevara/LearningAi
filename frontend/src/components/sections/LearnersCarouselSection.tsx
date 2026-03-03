"use client";

import { useState, useEffect, useCallback } from "react";
import { motion } from "framer-motion";

const learners = [
  {
    id: 1,
    icon: "🎓",
    title: "Students",
    description:
      "Master complex subjects with personalized explanations that match your academic needs and learning pace.",
  },
  {
    id: 2,
    icon: "📚",
    title: "High School Learners",
    description:
      "Transform textbooks into engaging content that makes studying feel less like a chore and more like discovery.",
  },
  {
    id: 3,
    icon: "💼",
    title: "Professionals",
    description:
      "Upskill efficiently with bite-sized learning that fits your busy schedule and career goals.",
  },
  {
    id: 4,
    icon: "🧭",
    title: "Self Learners",
    description:
      "Explore any topic at your own pace with adaptive content that evolves with your curiosity.",
  },
  {
    id: 5,
    icon: "🎨",
    title: "Creators",
    description:
      "Fuel your creative projects with deep knowledge presented in inspiring, story-driven formats.",
  },
];

export function LearnersCarouselSection() {
  const [activeIndex, setActiveIndex] = useState(0);
  const [isAutoPlaying, setIsAutoPlaying] = useState(true);

  const nextSlide = useCallback(() => {
    setActiveIndex((prev) => (prev + 1) % learners.length);
  }, []);

  const prevSlide = useCallback(() => {
    setActiveIndex((prev) => (prev - 1 + learners.length) % learners.length);
  }, []);

  const goToSlide = useCallback((index: number) => {
    setActiveIndex(index);
    setIsAutoPlaying(false);
    setTimeout(() => setIsAutoPlaying(true), 10000);
  }, []);

  useEffect(() => {
    if (!isAutoPlaying) return;
    const interval = setInterval(nextSlide, 5000);
    return () => clearInterval(interval);
  }, [isAutoPlaying, nextSlide]);

  return (
    <section className="relative py-20 overflow-hidden">
      <div className="max-w-6xl mx-auto px-6">
        <div className="text-center mb-12">
          <p className="text-accent-pink text-sm font-semibold uppercase tracking-widest mb-4">
            Who is it For?
          </p>
          <h2 className="text-4xl md:text-5xl font-bold">
            Built for <span className="text-gradient">Every Kind of Learner.</span>
          </h2>
          <p className="mt-4 text-brand-200/70 max-w-2xl mx-auto text-lg">
            Whether you're studying, working, building, or creating — NeuraLearn adapts to your learning style.
          </p>
        </div>

        <div className="relative h-[420px] flex items-center justify-center">
          {learners.map((learner, index) => (
            <motion.div
              key={learner.id}
              animate={{
                scale: index === activeIndex ? 1 : 0.8,
                opacity: index === activeIndex ? 1 : 0.4,
              }}
              transition={{ duration: 0.4 }}
              className="absolute w-80 bg-slate-800 p-8 rounded-2xl shadow-xl cursor-pointer"
              onClick={() => goToSlide(index)}
            >
              <div className="text-4xl mb-4">{learner.icon}</div>
              <h3 className="text-xl font-bold mb-3">{learner.title}</h3>
              <p className="text-sm text-white/70">{learner.description}</p>
            </motion.div>
          ))}
        </div>

        <div className="flex justify-center gap-2 mt-10">
          {learners.map((_, index) => (
            <button
              key={index}
              onClick={() => goToSlide(index)}
              className={`w-3 h-3 rounded-full ${
                index === activeIndex ? "bg-white" : "bg-white/30"
              }`}
            />
          ))}
        </div>
      </div>
    </section>
  );
}