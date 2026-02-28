"use client";

import { motion } from "framer-motion";
import { cn } from "@/lib/utils";

interface GradientCardProps {
  children: React.ReactNode;
  className?: string;
  glowColor?: "purple" | "cyan" | "pink";
  hover?: boolean;
}

const glowColors = {
  purple: "hover:shadow-brand-500/20",
  cyan: "hover:shadow-accent-cyan/20",
  pink: "hover:shadow-accent-pink/20",
};

export function GradientCard({
  children,
  className,
  glowColor = "purple",
  hover = true,
}: GradientCardProps) {
  return (
    <motion.div
      initial={{ opacity: 0, y: 20 }}
      whileInView={{ opacity: 1, y: 0 }}
      viewport={{ once: true, margin: "-50px" }}
      transition={{ duration: 0.5, ease: "easeOut" }}
      whileHover={hover ? { y: -4, transition: { duration: 0.2 } } : undefined}
      className={cn(
        "glass rounded-2xl p-6 transition-shadow duration-300",
        hover && `hover:shadow-2xl ${glowColors[glowColor]}`,
        className
      )}
    >
      {children}
    </motion.div>
  );
}
