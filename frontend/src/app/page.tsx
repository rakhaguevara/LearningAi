import { Navbar } from "@/components/layout/Navbar";
import { Footer } from "@/components/layout/Footer";
import { HeroSection } from "@/components/sections/HeroSection";
import { FeaturesSection } from "@/components/sections/FeaturesSection";
import { HowItWorksSection } from "@/components/sections/HowItWorksSection";
import { LearningFormatsSection } from "@/components/sections/LearningFormatsSection";
import { LearnersCarouselSection } from "@/components/sections/LearnersCarouselSection";
import { DemoSection } from "@/components/sections/DemoSection";
import { FinalCTASection } from "@/components/sections/FinalCTASection";

export default function Home() {
  return (
    <>
      <Navbar />
      <main>
        <HeroSection />
        <FeaturesSection />
        <HowItWorksSection />
        <LearningFormatsSection />
        <LearnersCarouselSection />
        <DemoSection />
        <FinalCTASection />
      </main>
      <Footer />
    </>
  );
}
