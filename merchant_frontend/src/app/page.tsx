"use client";

import Navbar from "../components/Navbar";
import Hero from "../components/Hero";
import CTA from "../components/CTA";
import AdBanner from "../components/AdBanner";
import HeroTransactions from "../components/Hero2";
import Link from 'next/link';

export default function Home() {
  return (
      <main className="">
        <div className="bg-[#B5FF37] overflow-hidden px-6 md:px-16 min-h-screen">
          <Navbar />
          <Hero />
          <CTA />
          <AdBanner />
        </div>
        <div className="bg-white relative overflow-hidden py-20 px-6 md:px-16 min-h-screen">
        <HeroTransactions />
        </div>
      </main>
  );
}
