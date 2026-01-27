"use client";
import { motion } from "framer-motion";
import Image from "next/image";

export default function Hero() {
  return (
    <motion.section 
      initial={{ opacity: 0 }}
      animate={{ opacity: 1 }}
      transition={{ delay: 0.2 }}
      className="mt-20 flex flex-col md:flex-row items-center justify-between gap-10"
    >
      <motion.h1 
        initial={{ x: -60, opacity: 0 }}
        animate={{ x: 0, opacity: 1 }}
        transition={{ duration: 0.6 }}
        className="text-black text-5xl md:text-7xl font-extrabold font-['Poppins'] uppercase"
      >
        Analytics
      </motion.h1>

      <Image
        src="/card.svg"
        alt="Card Graphic"
        width={224}
        height={180}
        className="w-auto h-auto"
      />

      <motion.h1 
        initial={{ x: 60, opacity: 0 }}
        animate={{ x: 0, opacity: 1 }}
        transition={{ duration: 0.6 }}
        className="text-black text-5xl md:text-7xl font-extrabold font-['Poppins'] uppercase"
      >
        Dashboard
      </motion.h1>
    </motion.section>
  );
}
