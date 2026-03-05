"use client";

import { useState, useEffect, useCallback } from "react";
import { motion, AnimatePresence } from "framer-motion";
import { FloatingIllustration } from "@/components/ui/FloatingIllustration";

const learners = [
  {
    id: 1,
    icon: "🎓",
    title: "Students",
    description: "Master complex subjects with personalized explanations that match your academic needs and learning pace.",
  },
  {
    id: 2,
    icon: "📚",
    title: "High School Learners",
    description: "Transform textbooks into engaging content that makes studying feel less like a chore and more like discovery.",
  },
  {
    id: 3,
    icon: "💼",
    title: "Professionals",
    description: "Upskill efficiently with bite-sized learning that fits your busy schedule and career goals.",
  },
  {
    id: 4,
    icon: "🧭",
    title: "Self Learners",
    description: "Explore any topic at your own pace with adaptive content that evolves with your curiosity.",
  },
  {
    id: 5,
    icon: "🎨",
    title: "Creators",
    description: "Fuel your creative projects with deep knowledge presented in inspiring, story-driven formats.",
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

    const interval = setInterval(() => {
      nextSlide();
    }, 5000);

    return () => clearInterval(interval);
  }, [isAutoPlaying, nextSlide]);

  const getCardStyle = (index: number) => {
    const diff = index - activeIndex;
    const normalizedDiff = ((diff + learners.length + Math.floor(learners.length / 2)) % learners.length) - Math.floor(learners.length / 2);

    if (normalizedDiff === 0) {
      return {
        transform: "translateX(0) scale(1) rotateY(0deg)",
        opacity: 1,
        zIndex: 20,
      };
    } else if (normalizedDiff === 1 || normalizedDiff === -4) {
      return {
        transform: "translateX(60%) scale(0.85) rotateY(-25deg)",
        opacity: 0.5,
        zIndex: 10,
      };
    } else if (normalizedDiff === -1 || normalizedDiff === 4) {
      return {
        transform: "translateX(-60%) scale(0.85) rotateY(25deg)",
        opacity: 0.5,
        zIndex: 10,
      };
    } else {
      return {
        transform: "translateX(0) scale(0.7)",
        opacity: 0,
        zIndex: 0,
      };
    }
  };

  return (
    <section className="relative py-16 md:py-24 overflow-hidden">
      {/* Decorative Orbs */}
      <FloatingIllustration
        variant="orb"
        color="pink"
        size="lg"
        className="absolute -left-20 top-1/4 opacity-20"
      />
      <FloatingIllustration
        variant="dots"
        color="cyan"
        className="absolute right-10 top-1/2 opacity-30 hidden md:flex"
      />

      <div className="relative z-10 max-w-7xl mx-auto px-6">
        {/* Header */}
        <div className="text-center mb-2 md:mb-4">
          <motion.p
            initial={{ opacity: 0 }}
            whileInView={{ opacity: 1 }}
            viewport={{ once: true }}
            className="text-accent-pink text-sm font-semibold uppercase tracking-widest mb-4"
          >
            Who is it For?
          </motion.p>
          <motion.h2
            initial={{ opacity: 0, y: 20 }}
            whileInView={{ opacity: 1, y: 0 }}
            viewport={{ once: true }}
            transition={{ delay: 0.1 }}
            className="text-4xl md:text-5xl font-bold"
          >
            Built for{" "}
            <span className="text-gradient">Every Kind of Learner.</span>
          </motion.h2>
          <motion.p
            initial={{ opacity: 0, y: 20 }}
            whileInView={{ opacity: 1, y: 0 }}
            viewport={{ once: true }}
            transition={{ delay: 0.2 }}
            className="mt-4 text-brand-200/70 max-w-2xl mx-auto text-lg"
          >
            Whether you&apos;re studying, working, building, or creating — NeuraLearn adapts to your learning style.
          </motion.p>
        </div>

        {/* 3D Carousel */}
        <div className="relative h-[400px] md:h-[450px] perspective-1000">
          <div className="relative w-full h-full flex items-center justify-center">
            {learners.map((learner, index) => (
              <motion.div
                key={learner.id}
                className="absolute w-[280px] md:w-[320px] cursor-pointer"
                style={{
                  ...getCardStyle(index),
                  transformStyle: "preserve-3d",
                  transition: "all 0.5s ease-in-out",
                }}
                onClick={() => goToSlide(index)}
              >
                <div className="h-[280px] md:h-[320px] bg-slate-800/40 backdrop-blur-md border border-white/10 rounded-2xl shadow-md p-6 md:p-8 flex flex-col items-center justify-center text-center transition-all duration-300 hover:border-white/20 hover:shadow-lg">
                  {/* Icon */}
                  <div className="text-4xl md:text-5xl mb-4 opacity-60">
                    {learner.icon}
                  </div>

                  {/* Title */}
                  <h3 className="text-lg md:text-xl font-semibold text-white mb-3">
                    {learner.title}
                  </h3>

                  {/* Description */}
                  <p className="text-xs md:text-sm text-brand-200/70 leading-relaxed">
                    {learner.description}
                  </p>
                </div>
              </motion.div>
            ))}
          </div>

          {/* Navigation Arrows */}
          <button
            onClick={() => {
              prevSlide();
              setIsAutoPlaying(false);
              setTimeout(() => setIsAutoPlaying(true), 10000);
            }}
            className="absolute left-2 md:left-8 top-1/2 -translate-y-1/2 z-30 w-10 h-10 md:w-12 md:h-12 rounded-full bg-slate-800/60 backdrop-blur-sm border border-white/10 flex items-center justify-center text-white/70 hover:text-white hover:bg-slate-700/60 transition-all duration-200"
            aria-label="Previous slide"
          >
            <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2">
              <path d="M15 18l-6-6 6-6" />
            </svg>
          </button>

          <button
            onClick={() => {
              nextSlide();
              setIsAutoPlaying(false);
              setTimeout(() => setIsAutoPlaying(true), 10000);
            }}
            className="absolute right-2 md:right-8 top-1/2 -translate-y-1/2 z-30 w-10 h-10 md:w-12 md:h-12 rounded-full bg-slate-800/60 backdrop-blur-sm border border-white/10 flex items-center justify-center text-white/70 hover:text-white hover:bg-slate-700/60 transition-all duration-200"
            aria-label="Next slide"
          >
            <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2">
              <path d="M9 18l6-6-6-6" />
            </svg>
          </button>
        </div>

        {/* Dots Indicator */}
        <div className="flex justify-center gap-2 mt-8">
          {learners.map((_, index) => (
            <button
              key={index}
              onClick={() => goToSlide(index)}
              className={`w-2 h-2 rounded-full transition-all duration-300 ${index === activeIndex
                ? "bg-accent-cyan w-6"
                : "bg-white/20 hover:bg-white/40"
                }`}
              aria-label={`Go to slide ${index + 1}`}
            />
          ))}
        </div>
      </div>
    </section>
  );
}
