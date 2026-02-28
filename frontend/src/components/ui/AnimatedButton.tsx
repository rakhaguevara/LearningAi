"use client";

import { motion } from "framer-motion";
import { cn } from "@/lib/utils";

interface AnimatedButtonProps {
  children: React.ReactNode;
  variant?: "primary" | "secondary" | "ghost";
  size?: "sm" | "md" | "lg";
  className?: string;
  onClick?: () => void;
  type?: "button" | "submit";
}

const variants = {
  primary:
    "bg-gradient-to-r from-brand-500 to-accent-cyan text-white glow-purple hover:shadow-lg hover:shadow-brand-500/25",
  secondary:
    "glass text-white hover:bg-white/10 border border-white/10",
  ghost:
    "text-brand-200 hover:text-white hover:bg-white/5",
};

const sizes = {
  sm: "px-4 py-2 text-sm rounded-lg",
  md: "px-6 py-3 text-base rounded-xl",
  lg: "px-8 py-4 text-lg rounded-xl",
};

export function AnimatedButton({
  children,
  variant = "primary",
  size = "md",
  className,
  onClick,
  type = "button",
}: AnimatedButtonProps) {
  return (
    <motion.button
      type={type}
      whileHover={{ scale: 1.03, y: -1 }}
      whileTap={{ scale: 0.97 }}
      transition={{ type: "spring", stiffness: 400, damping: 17 }}
      className={cn(
        "font-semibold transition-all duration-300 cursor-pointer inline-flex items-center justify-center gap-2",
        variants[variant],
        sizes[size],
        className
      )}
      onClick={onClick}
    >
      {children}
    </motion.button>
  );
}
