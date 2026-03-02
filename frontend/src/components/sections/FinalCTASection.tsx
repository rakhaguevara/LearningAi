"use client";

import { useEffect, useRef, useState } from "react";
import { motion } from "framer-motion";
import { AnimatedButton } from "@/components/ui/AnimatedButton";
import { useRouter } from "next/navigation";

export function FinalCTASection() {
  const router = useRouter();
  const sectionRef = useRef<HTMLElement>(null);
  const [isVisible, setIsVisible] = useState(false);
  const [isFullyVisible, setIsFullyVisible] = useState(false);

  useEffect(() => {
    const section = sectionRef.current;
    if (!section) return;

    // Observer for early visibility (icon appears before reaching section)
    const earlyObserver = new IntersectionObserver(
      (entries) => {
        entries.forEach((entry) => {
          if (entry.isIntersecting) {
            setIsVisible(true);
          } else {
            setIsVisible(false);
            setIsFullyVisible(false);
          }
        });
      },
      {
        threshold: 0,
        rootMargin: "200px 0px 0px 0px", // Trigger 200px before entering viewport
      }
    );

    // Observer for full visibility (floating animation)
    const fullObserver = new IntersectionObserver(
      (entries) => {
        entries.forEach((entry) => {
          if (entry.isIntersecting && entry.intersectionRatio > 0.5) {
            setIsFullyVisible(true);
          } else {
            setIsFullyVisible(false);
          }
        });
      },
      {
        threshold: [0, 0.5, 1],
        rootMargin: "0px",
      }
    );

    earlyObserver.observe(section);
    fullObserver.observe(section);
    
    return () => {
      earlyObserver.disconnect();
      fullObserver.disconnect();
    };
  }, []);

  return (
    <section
      ref={sectionRef}
      className="relative py-24 md:py-32 lg:py-40 overflow-hidden"
    >
      {/* Large 3D Background Icon - Appears Early */}
      <motion.div
        className="absolute inset-0 flex items-center justify-center pointer-events-none"
        initial={{ opacity: 0, y: 120, scale: 0.9 }}
        animate={{
          opacity: isVisible ? (isFullyVisible ? 0.35 : 0.2) : 0,
          y: isVisible ? 0 : 120,
          scale: isVisible ? 1 : 0.9,
        }}
        transition={{
          duration: 1,
          ease: [0.25, 0.46, 0.45, 0.94],
        }}
        style={{ willChange: "transform, opacity" }}
      >
        <motion.div
          animate={
            isFullyVisible
              ? {
                  rotate: [0, 3, 0, -3, 0],
                  y: [0, -8, 0, 8, 0],
                }
              : isVisible
              ? {
                  rotate: [0, 1, 0, -1, 0],
                  y: [0, -4, 0, 4, 0],
                }
              : {}
          }
          transition={{
            duration: isFullyVisible ? 8 : 12,
            repeat: Infinity,
            ease: "easeInOut",
          }}
          className="blur-[1px]"
          style={{ willChange: "transform" }}
        >
          <svg
            viewBox="0 0 200 200"
            className="w-[350px] h-[350px] md:w-[500px] md:h-[500px] lg:w-[600px] lg:h-[600px] text-slate-400"
            fill="none"
            stroke="currentColor"
            strokeWidth="0.8"
          >
            {/* Neural Network Abstract Icon */}
            {/* Central Node */}
            <circle cx="100" cy="100" r="12" strokeWidth="1.5" />
            
            {/* Outer Ring */}
            <circle cx="100" cy="100" r="45" opacity="0.6" />
            <circle cx="100" cy="100" r="70" opacity="0.4" />
            <circle cx="100" cy="100" r="90" opacity="0.2" />
            
            {/* Connected Nodes - Inner Ring */}
            <circle cx="100" cy="55" r="6" />
            <circle cx="145" cy="100" r="6" />
            <circle cx="100" cy="145" r="6" />
            <circle cx="55" cy="100" r="6" />
            
            {/* Connected Nodes - Outer Ring */}
            <circle cx="100" cy="30" r="4" opacity="0.7" />
            <circle cx="170" cy="100" r="4" opacity="0.7" />
            <circle cx="100" cy="170" r="4" opacity="0.7" />
            <circle cx="30" cy="100" r="4" opacity="0.7" />
            <circle cx="150" cy="50" r="4" opacity="0.5" />
            <circle cx="150" cy="150" r="4" opacity="0.5" />
            <circle cx="50" cy="150" r="4" opacity="0.5" />
            <circle cx="50" cy="50" r="4" opacity="0.5" />
            
            {/* Connection Lines - Inner */}
            <line x1="100" y1="100" x2="100" y2="55" />
            <line x1="100" y1="100" x2="145" y2="100" />
            <line x1="100" y1="100" x2="100" y2="145" />
            <line x1="100" y1="100" x2="55" y2="100" />
            
            {/* Connection Lines - Cross */}
            <line x1="100" y1="55" x2="145" y2="100" opacity="0.6" />
            <line x1="145" y1="100" x2="100" y2="145" opacity="0.6" />
            <line x1="100" y1="145" x2="55" y2="100" opacity="0.6" />
            <line x1="55" y1="100" x2="100" y2="55" opacity="0.6" />
            
            {/* Connection Lines - Outer */}
            <line x1="100" y1="55" x2="100" y2="30" opacity="0.5" />
            <line x1="145" y1="100" x2="170" y2="100" opacity="0.5" />
            <line x1="100" y1="145" x2="100" y2="170" opacity="0.5" />
            <line x1="55" y1="100" x2="30" y2="100" opacity="0.5" />
            
            {/* Diagonal Connections */}
            <line x1="145" y1="100" x2="150" y2="50" opacity="0.4" />
            <line x1="145" y1="100" x2="150" y2="150" opacity="0.4" />
            <line x1="55" y1="100" x2="50" y2="150" opacity="0.4" />
            <line x1="55" y1="100" x2="50" y2="50" opacity="0.4" />
            <line x1="100" y1="55" x2="150" y2="50" opacity="0.4" />
            <line x1="100" y1="55" x2="50" y2="50" opacity="0.4" />
            <line x1="100" y1="145" x2="150" y2="150" opacity="0.4" />
            <line x1="100" y1="145" x2="50" y2="150" opacity="0.4" />
          </svg>
        </motion.div>
      </motion.div>

      {/* Content */}
      <div className="relative z-10 max-w-7xl mx-auto px-6 text-center">
        <motion.h2
          initial={{ opacity: 0, y: 30 }}
          whileInView={{ opacity: 1, y: 0 }}
          viewport={{ once: true }}
          transition={{ duration: 0.6, delay: 0.1 }}
          className="text-4xl md:text-5xl lg:text-6xl font-bold"
        >
          Ready to{" "}
          <span className="text-gradient">Understand Differently?</span>
        </motion.h2>

        <motion.p
          initial={{ opacity: 0, y: 30 }}
          whileInView={{ opacity: 1, y: 0 }}
          viewport={{ once: true }}
          transition={{ duration: 0.6, delay: 0.2 }}
          className="mt-6 text-lg md:text-xl text-brand-200/80 max-w-2xl mx-auto"
        >
          Stop forcing yourself to learn the hard way. Let your materials adapt to you.
        </motion.p>

        <motion.div
          initial={{ opacity: 0, y: 30 }}
          whileInView={{ opacity: 1, y: 0 }}
          viewport={{ once: true }}
          transition={{ duration: 0.6, delay: 0.3 }}
          className="mt-10"
        >
          <AnimatedButton
            onClick={() => router.push("/login")}
            variant="primary"
            size="lg"
          >
            Start Learning Your Way
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
        </motion.div>
      </div>
    </section>
  );
}
