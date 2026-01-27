"use client";
import { motion } from "framer-motion";
import Link from "next/link";

export default function CTA() {
  return (
    <motion.section 
      initial={{ y: 40, opacity: 0 }}
      animate={{ y: 0, opacity: 1 }}
      transition={{ delay: 0.3, duration: 0.6 }}
      className="mt-24 flex flex-col justify-center md:flex-row items-center gap-10"
    >
      <div className="hidden md:block w-12 border-t-2 border-black" />
      <p className="text-black text-xl font-['Poppins'] uppercase max-w-md mr-16">
        Submit authorization requests with merchant, card, and transaction details
      </p>
          <motion.button 
            whileHover={{ scale: 1.06 }}
            whileTap={{ scale: 0.96 }}
            className="bg-black px-10 py-5 text-white text-xl font-bold font-['Poppins'] uppercase"
          >
            <Link href="/auth">       
                Send Authorization
            </Link>
          </motion.button>
    </motion.section>
  );
}
