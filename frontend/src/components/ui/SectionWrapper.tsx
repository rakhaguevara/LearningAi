"use client";

import { motion } from "framer-motion";
import { cn } from "@/lib/utils";

interface SectionWrapperProps {
  children: React.ReactNode;
  id?: string;
  className?: string;
  withContainer?: boolean;
}

export function SectionWrapper({
  children,
  id,
  className,
  withContainer = true,
}: SectionWrapperProps) {
  return (
    <section id={id} className={cn("relative py-24 lg:py-32", className)}>
      {withContainer ? (
        <div className="mx-auto max-w-7xl px-6 lg:px-8">
          <motion.div
            initial={{ opacity: 0 }}
            whileInView={{ opacity: 1 }}
            viewport={{ once: true, margin: "-100px" }}
            transition={{ duration: 0.6 }}
          >
            {children}
          </motion.div>
        </div>
      ) : (
        children
      )}
    </section>
  );
}
