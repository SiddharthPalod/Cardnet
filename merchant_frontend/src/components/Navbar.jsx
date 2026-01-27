"use client";
import { motion } from "framer-motion";
import Image from "next/image";
import circle from "../../public/circle.svg";
import Link from "next/link";
import { useState } from "react"
import Newsletter from "./Newsletter";

export default function Navbar() {
  const [open, setOpen] = useState(false)

  return (
    <motion.nav 
      initial={{ y: -40, opacity: 0 }}
      animate={{ y: 0, opacity: 1 }}
      transition={{ duration: 0.6 }}
      className="flex items-center py-8"
    >
        <motion.div
          initial={{ scaleY: 0 }}
          animate={{ scaleY: 1 }}
          transition={{ duration: 0.6 }}
          className="origin-top"
        >
          <Image
            src={circle}
            alt="Real-Time Indicator"
            width={56}
            height={56}
            className="w-14 mr-16"
          />
        </motion.div>
      <div className="hidden md:flex items-center gap-28">
        <span className="text-black text-xl font-['Poppins'] uppercase">
        <Link href="/">Home</Link>
        </span>
        <span className="text-black text-xl font-['Poppins'] uppercase">
          <Link href="/auth">Simulate</Link>
        </span>

        <span className="text-black text-5xl font-['Pacifico']">
          CardNet
        </span>

      <span
        onClick={() => setOpen(true)}
        className="text-black text-xl font-poppins uppercase cursor-pointer hover:opacity-70 transition"
      >
        Connect
      </span>
      {open && (
        <div
          className="fixed inset-0 z-50 flex items-center justify-center bg-black/50 backdrop-blur-sm px-4"
          onClick={() => setOpen(false)}
        >
          {/* Modal Content */}
          <div
            className="relative w-full max-w-lg animate-fade-in"
            onClick={(e) => e.stopPropagation()}
          >
            <Newsletter />
          </div>
        </div>
      )}


        <span className="text-black text-xl font-['Poppins'] uppercase">
          <Link href="/history">History</Link>
        </span>
      </div>

      <div className="md:hidden text-black text-2xl font-bold">☰</div>
    </motion.nav>
  );
}
