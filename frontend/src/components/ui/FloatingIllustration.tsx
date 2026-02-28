"use client";

import { motion } from "framer-motion";
import { cn } from "@/lib/utils";

interface FloatingIllustrationProps {
  className?: string;
  variant?: "orb" | "ring" | "dots";
  color?: "purple" | "cyan" | "pink";
  size?: "sm" | "md" | "lg";
}

const colorMap = {
  purple: {
    bg: "bg-brand-500/20",
    border: "border-brand-400/30",
    shadow: "shadow-brand-500/20",
  },
  cyan: {
    bg: "bg-accent-cyan/20",
    border: "border-accent-cyan/30",
    shadow: "shadow-accent-cyan/20",
  },
  pink: {
    bg: "bg-accent-pink/20",
    border: "border-accent-pink/30",
    shadow: "shadow-accent-pink/20",
  },
};

const sizeMap = {
  sm: "w-32 h-32",
  md: "w-48 h-48",
  lg: "w-72 h-72",
};

export function FloatingIllustration({
  className,
  variant = "orb",
  color = "purple",
  size = "md",
}: FloatingIllustrationProps) {
  const colors = colorMap[color];
  const sizeClass = sizeMap[size];

  if (variant === "ring") {
    return (
      <motion.div
        animate={{ rotate: 360 }}
        transition={{ duration: 20, repeat: Infinity, ease: "linear" }}
        className={cn(
          sizeClass,
          "rounded-full border-2 border-dashed",
          colors.border,
          "opacity-30",
          className
        )}
      />
    );
  }

  if (variant === "dots") {
    return (
      <motion.div
        animate={{ y: [-10, 10, -10] }}
        transition={{ duration: 5, repeat: Infinity, ease: "easeInOut" }}
        className={cn("flex gap-2", className)}
      >
        {[0, 1, 2].map((i) => (
          <motion.div
            key={i}
            animate={{ scale: [1, 1.3, 1] }}
            transition={{
              duration: 2,
              repeat: Infinity,
              delay: i * 0.3,
              ease: "easeInOut",
            }}
            className={cn("w-3 h-3 rounded-full", colors.bg)}
          />
        ))}
      </motion.div>
    );
  }

  // Default: orb
  return (
    <motion.div
      animate={{
        y: [-15, 15, -15],
        scale: [1, 1.05, 1],
      }}
      transition={{ duration: 6, repeat: Infinity, ease: "easeInOut" }}
      className={cn(
        sizeClass,
        "rounded-full blur-3xl",
        colors.bg,
        "opacity-50",
        className
      )}
    />
  );
}
