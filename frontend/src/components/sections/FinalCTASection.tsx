"use client";

import { useRef } from "react";
import { motion, useInView } from "framer-motion";
import { SectionWrapper } from "@/components/ui/SectionWrapper";
import { AnimatedButton } from "@/components/ui/AnimatedButton";
import { useRouter } from "next/navigation";

export function FinalCTASection() {
    const router = useRouter();
    const containerRef = useRef<HTMLDivElement>(null);
    const isInView = useInView(containerRef, { margin: "-20% 0px -20% 0px", amount: 0.3 });

    // Custom bezier for cinematic feel
    const cinematicEase = [0.16, 1, 0.3, 1];

    return (
        <section
            ref={containerRef}
            className="relative overflow-hidden w-full py-20 md:py-32 lg:py-40 flex flex-col items-center justify-center border-t border-white/5"
        >
            {/* --- ARC BACKGROUND ELEMENT --- */}
            <motion.div
                className="absolute bottom-0 left-1/2 w-[600px] h-[300px] md:w-[1100px] md:h-[550px] lg:w-[1600px] lg:h-[800px] pointer-events-none z-0"
                style={{
                    borderRadius: "50% 50% 0 0 / 100% 100% 0 0",
                    borderTop: "2px solid rgba(196, 181, 253, 0.15)", // brand text color low opacity
                    background: "linear-gradient(to top, rgba(124, 58, 237, 0.05) 0%, transparent 100%)",
                    filter: "blur(2px)",
                    x: "-50%", // Keep it centered horizontally
                }}
                initial={{ opacity: 0, y: 120 }}
                animate={{
                    opacity: isInView ? 0.3 : 0,
                    y: isInView ? 0 : 120
                }}
                transition={{
                    duration: 1.2,
                    ease: cinematicEase
                }}
            >
                {/* Arc Idle Float */}
                <motion.div
                    className="w-full h-full"
                    animate={isInView ? { y: [0, -6, 0] } : {}}
                    transition={{ duration: 8, repeat: Infinity, ease: "easeInOut" }}
                />
            </motion.div>

            {/* --- 3D ABSTRACT ICON BACKGROUND --- */}
            <motion.div
                className="absolute top-1/2 left-1/2 w-[200px] h-[200px] md:w-[320px] md:h-[320px] lg:w-[450px] lg:h-[450px] pointer-events-none z-0"
                style={{ x: "-50%", y: "-50%", filter: "blur(2px)" }}
                initial={{ opacity: 0, y: "-35%" }} // -50% + 70px approx
                animate={{
                    opacity: isInView ? 0.35 : 0,
                    y: isInView ? "-50%" : "-35%"
                }}
                transition={{
                    duration: 1,
                    ease: cinematicEase,
                    delay: 0.1
                }}
            >
                <motion.div
                    className="w-full h-full relative"
                    animate={isInView ? {
                        y: [0, -8, 0],
                        rotate: [0, 3, 0]
                    } : {}}
                    transition={{
                        duration: 7,
                        repeat: Infinity,
                        ease: "easeInOut"
                    }}
                >
                    {/* Abstract Knowledge/Neural Node SVG */}
                    <svg viewBox="0 0 200 200" fill="none" xmlns="http://www.w3.org/2000/svg" className="w-full h-full text-brand-300">
                        {/* Core Orb */}
                        <circle cx="100" cy="100" r="40" fill="currentColor" opacity="0.8" />
                        <circle cx="100" cy="100" r="50" stroke="currentColor" strokeWidth="1" strokeDasharray="4 4" opacity="0.5" />
                        <circle cx="100" cy="100" r="70" stroke="currentColor" strokeWidth="2" opacity="0.3" />

                        {/* Orbiting Nodes */}
                        <circle cx="35" cy="100" r="8" fill="currentColor" opacity="0.6" />
                        <circle cx="165" cy="100" r="10" fill="currentColor" opacity="0.7" />
                        <circle cx="100" cy="25" r="6" fill="currentColor" opacity="0.5" />
                        <circle cx="100" cy="175" r="12" fill="currentColor" opacity="0.8" />
                        <circle cx="50" cy="50" r="5" fill="currentColor" opacity="0.4" />
                        <circle cx="150" cy="150" r="7" fill="currentColor" opacity="0.6" />

                        {/* Connecting Lines */}
                        <path d="M100 60 L100 25 M100 140 L100 175 M60 100 L35 100 M140 100 L165 100 M72 72 L50 50 M128 128 L150 150" stroke="currentColor" strokeWidth="1.5" opacity="0.4" />
                        <path d="M35 100 Q 50 150 100 175 Q 150 150 165 100 Q 150 50 100 25 Q 50 50 35 100" stroke="currentColor" strokeWidth="0.5" fill="none" opacity="0.2" />
                    </svg>
                </motion.div>
            </motion.div>

            {/* --- CONTENT --- */}
            <div className="relative z-10 mx-auto max-w-7xl px-6 lg:px-8 text-center flex flex-col items-center">
                <motion.h2
                    initial={{ opacity: 0, y: 30 }}
                    animate={{ opacity: isInView ? 1 : 0, y: isInView ? 0 : 30 }}
                    transition={{ duration: 0.8, ease: cinematicEase, delay: 0.2 }}
                    className="text-5xl md:text-6xl lg:text-7xl font-bold tracking-tight text-white mb-6"
                >
                    Ready to Understand <br className="hidden md:block" />
                    <span className="text-gradient">Differently?</span>
                </motion.h2>

                <motion.p
                    initial={{ opacity: 0, y: 20 }}
                    animate={{ opacity: isInView ? 1 : 0, y: isInView ? 0 : 20 }}
                    transition={{ duration: 0.8, ease: cinematicEase, delay: 0.3 }}
                    className="text-lg md:text-xl text-brand-100/80 max-w-2xl mx-auto mb-12 font-light leading-relaxed tracking-wide"
                >
                    Stop forcing yourself to learn the hard way. <br className="hidden sm:block" />
                    Let your materials adapt to you.
                </motion.p>

                <motion.div
                    initial={{ opacity: 0, y: 20 }}
                    animate={{ opacity: isInView ? 1 : 0, y: isInView ? 0 : 20 }}
                    transition={{ duration: 0.8, ease: cinematicEase, delay: 0.4 }}
                >
                    <AnimatedButton onClick={() => router.push('/login')} variant="primary" size="lg" className="shadow-[0_0_40px_-10px_rgba(124,58,237,0.5)] hover:shadow-[0_0_60px_-10px_rgba(124,58,237,0.7)] transition-shadow duration-500">
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
                            className="ml-2"
                        >
                            <path d="M5 12h14M12 5l7 7-7 7" />
                        </svg>
                    </AnimatedButton>
                </motion.div>
            </div>
        </section>
    );
}
