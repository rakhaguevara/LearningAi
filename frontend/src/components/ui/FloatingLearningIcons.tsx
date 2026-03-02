"use client";

import { useEffect, useState } from "react";
import { motion } from "framer-motion";

interface IconConfig {
  id: number;
  icon: string;
  position: {
    top?: string;
    bottom?: string;
    left?: string;
    right?: string;
  };
  size: {
    default: string;
    sm?: string;
    md?: string;
    lg?: string;
    xl?: string;
  };
  delay: number;
  duration: number;
  mobile?: boolean;
}

const learningIcons: IconConfig[] = [
  {
    id: 1,
    icon: "BookOpen",
    position: { top: "15%", left: "5%" },
    size: { default: "w-8 h-8", sm: "sm:w-10 sm:h-10", md: "md:w-12 md:h-12", lg: "lg:w-14 lg:h-14" },
    delay: 0,
    duration: 3,
    mobile: true,
  },
  {
    id: 2,
    icon: "Brain",
    position: { top: "25%", right: "8%" },
    size: { default: "w-10 h-10", sm: "sm:w-12 sm:h-12", md: "md:w-14 md:h-14", lg: "lg:w-16 lg:h-16" },
    delay: 0.5,
    duration: 4,
    mobile: true,
  },
  {
    id: 3,
    icon: "Atom",
    position: { top: "60%", left: "3%" },
    size: { default: "w-8 h-8", sm: "sm:w-10 sm:h-10", md: "md:w-12 md:h-12", lg: "lg:w-14 lg:h-14" },
    delay: 1,
    duration: 3.5,
    mobile: false,
  },
  {
    id: 4,
    icon: "Pencil",
    position: { bottom: "20%", right: "5%" },
    size: { default: "w-6 h-6", sm: "sm:w-8 sm:h-8", md: "md:w-10 md:h-10", lg: "lg:w-12 lg:h-12" },
    delay: 0.3,
    duration: 2.5,
    mobile: true,
  },
  {
    id: 5,
    icon: "Video",
    position: { top: "45%", right: "12%" },
    size: { default: "w-7 h-7", sm: "sm:w-9 sm:h-9", md: "md:w-11 md:h-11", lg: "lg:w-12 lg:h-12" },
    delay: 0.8,
    duration: 3.2,
    mobile: false,
  },
  {
    id: 6,
    icon: "Podcast",
    position: { bottom: "30%", left: "8%" },
    size: { default: "w-6 h-6", sm: "sm:w-8 sm:h-8", md: "md:w-10 md:h-10", lg: "lg:w-12 lg:h-12" },
    delay: 1.2,
    duration: 2.8,
    mobile: false,
  },
  {
    id: 7,
    icon: "Lightbulb",
    position: { top: "70%", right: "15%" },
    size: { default: "w-5 h-5", sm: "sm:w-7 sm:h-7", md: "md:w-8 md:h-8", lg: "lg:w-10 lg:h-10" },
    delay: 0.6,
    duration: 3.8,
    mobile: false,
  },
  {
    id: 8,
    icon: "GraduationCap",
    position: { top: "10%", right: "20%" },
    size: { default: "w-7 h-7", sm: "sm:w-9 sm:h-9", md: "md:w-11 md:h-11", lg: "lg:w-12 lg:h-12" },
    delay: 0.2,
    duration: 3.5,
    mobile: true,
  },
];

const iconPaths: Record<string, string> = {
  BookOpen: "M2 3h6a4 4 0 0 1 4 4v14a3 3 0 0 0-3-3H2z M22 3h-6a4 4 0 0 0-4 4v14a3 3 0 0 1 3-3h7z",
  Brain: "M9.5 2A2.5 2.5 0 0 1 12 4.5v15a2.5 2.5 0 0 1-4.96.44 2.5 2.5 0 0 1-2.96-3.08 3 3 0 0 1-.34-5.58 2.5 2.5 0 0 1 1.32-4.24 2.5 2.5 0 0 1 1.98-3A2.5 2.5 0 0 1 9.5 2z M15.5 2a2.5 2.5 0 0 1 2.5 2.5 2.5 2.5 0 0 1 1.98 3 2.5 2.5 0 0 1 1.32 4.24 3 3 0 0 1-.34 5.58 2.5 2.5 0 0 1-2.96 3.08 2.5 2.5 0 0 1-4.96-.44v-15A2.5 2.5 0 0 1 15.5 2z",
  Atom: "M12 2a1 1 0 0 1 1 1v2a1 1 0 0 1-2 0V3a1 1 0 0 1 1-1z M4.929 4.929a1 1 0 0 1 1.414 0l1.414 1.414a1 1 0 0 1-1.414 1.414L4.929 6.343a1 1 0 0 1 0-1.414z M19.071 4.929a1 1 0 0 1 0 1.414l-1.414 1.414a1 1 0 1 1-1.414-1.414l1.414-1.414a1 1 0 0 1 1.414 0z M12 8a4 4 0 1 0 0 8 4 4 0 0 0 0-8z M2 12a1 1 0 0 1 1-1h2a1 1 0 1 1 0 2H3a1 1 0 0 1-1-1z M18 12a1 1 0 0 1 1-1h2a1 1 0 1 1 0 2h-2a1 1 0 0 1-1-1z M4.929 19.071a1 1 0 0 1 0-1.414l1.414-1.414a1 1 0 0 1 1.414 1.414l-1.414 1.414a1 1 0 0 1-1.414 0z M19.071 19.071a1 1 0 0 1-1.414 0l-1.414-1.414a1 1 0 0 1 1.414-1.414l1.414 1.414a1 1 0 0 1 0 1.414z M12 20a1 1 0 0 1 1 1v2a1 1 0 1 1-2 0v-2a1 1 0 0 1 1-1z",
  Pencil: "M17 3a2.828 2.828 0 1 1 4 4L7.5 20.5 2 22l1.5-5.5L17 3z",
  Video: "M23 7l-7 5 7 5V7z M1 7v10h15V7H1z",
  Podcast: "M12 2a3 3 0 0 0-3 3v7a3 3 0 0 0 6 0V5a3 3 0 0 0-3-3z M19 10v2a7 7 0 0 1-14 0v-2 M12 18v3",
  Lightbulb: "M9 18h6 M10 22h4 M12 2a7 7 0 0 0-7 7c0 2.38 1.19 4.47 3 5.74V17a1 1 0 0 0 1 1h6a1 1 0 0 0 1-1v-2.26c1.81-1.27 3-3.36 3-5.74a7 7 0 0 0-7-7z",
  GraduationCap: "M22 10v6M2 10l10-5 10 5-10 5z M6 12v5c0 2 2 3 6 3s6-1 6-3v-5",
};

export function FloatingLearningIcons() {
  const [scrollProgress, setScrollProgress] = useState(0);

  useEffect(() => {
    const handleScroll = () => {
      const heroHeight = window.innerHeight;
      const scrollY = window.scrollY;
      const progress = Math.min(scrollY / (heroHeight * 0.5), 1);
      setScrollProgress(progress);
    };

    window.addEventListener("scroll", handleScroll, { passive: true });
    return () => window.removeEventListener("scroll", handleScroll);
  }, []);

  const getIconStyle = (progress: number) => {
    const opacity = 1 - progress * 0.6;
    const scale = 1 - progress * 0.1;
    const translateY = progress * 30;

    return {
      opacity: Math.max(opacity, 0.2),
      transform: `scale(${scale}) translateY(${translateY}px)`,
      willChange: "transform, opacity",
    };
  };

  return (
    <div className="absolute inset-0 overflow-hidden pointer-events-none">
      {learningIcons.map((icon) => (
        <motion.div
          key={icon.id}
          className={`absolute ${icon.mobile ? "" : "hidden sm:block"}`}
          style={{
            ...icon.position,
            ...getIconStyle(scrollProgress),
            transition: "all 300ms ease-out",
          }}
          initial={{ opacity: 0, scale: 0.8 }}
          animate={{ opacity: 1, scale: 1 }}
          transition={{ duration: 0.5, delay: icon.delay }}
        >
          <motion.div
            animate={{
              y: [0, -8, 0],
            }}
            transition={{
              duration: icon.duration,
              repeat: Infinity,
              ease: "easeInOut",
            }}
            className="blur-[1px]"
          >
            <svg
              viewBox="0 0 24 24"
              fill="none"
              stroke="currentColor"
              strokeWidth="1.5"
              strokeLinecap="round"
              strokeLinejoin="round"
              className={`${icon.size.default} ${icon.size.sm || ""} ${icon.size.md || ""} ${icon.size.lg || ""} text-slate-400/40`}
            >
              <path d={iconPaths[icon.icon]} />
            </svg>
          </motion.div>
        </motion.div>
      ))}
    </div>
  );
}
