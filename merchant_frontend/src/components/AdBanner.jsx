"use client";

import { motion } from "framer-motion";
import Image from "next/image";
import circle from "../../public/circle.svg";

export default function AdBanner() {
  return (
    <motion.section 
      initial={{ opacity: 0 }}
      animate={{ opacity: 1 }}
      transition={{ delay: 0.4 }}
      className="my-28 grid md:grid-cols-2 gap-20 items-center"
    >

      {/* LEFT BLOCK */}
        <div className="relative w-full max-w-xl">

        {/* Image */}
        <motion.div
            initial={{ scale: 0.9, opacity: 0 }}
            animate={{ scale: 1, opacity: 1 }}
            transition={{ duration: 0.6 }}
            className="ml-auto shadow-xl rounded-lg overflow-hidden w-72 md:w-80"
        >
            <Image
            src="https://images.unsplash.com/photo-1712331393873-f671bec94e21?w=600&auto=format&fit=crop&q=60&ixlib=rb-4.1.0&ixid=M3wxMjA3fDB8MHxzZWFyY2h8Nnx8bWVyY2hhbnR8ZW58MHx8MHx8fDA%3D"
            alt="Merchant"
            width={320}
            height={420}
            className="object-cover w-full h-auto"
            />
        </motion.div>

        {/* OVERLAPPING TEXT */}
        <motion.h2 
            initial={{ x: -40, opacity: 0 }}
            animate={{ x: 0, opacity: 1 }}
            transition={{ duration: 0.6 }}
            className="absolute top-[30%] left-0 text-black text-4xl font-medium uppercase max-w-md leading-tight"
        >
            Simulate Card Authorization Requests
        </motion.h2>

        </div>


      {/* RIGHT LABEL */}
    <div className="relative h-56 flex items-center justify-center">

      {/* Center Circle Image */}
      <motion.div
        initial={{ rotate: 0, opacity: 0 }}
        animate={{ rotate: 90, opacity: 1 }}
        transition={{ duration: 0.6 }}
        className="z-10"
      >
        <Image
          src={circle}
          alt="Indicator"
          width={120}
          height={120}
          className="w-28 h-28"
        />
      </motion.div>

      {/* Circular Text Ring */}
      <motion.svg
        initial={{ rotate: -90, opacity: 0 }}
        animate={{ rotate: 0, opacity: 1 }}
        transition={{ duration: 0.8 }}
        viewBox="0 0 240 240"
        className="absolute w-full h-full"
      >
        <defs>
            <path
            id="bottomArc"
            d="
                M 40,120
                A 80,80 0 0,0 200,120
            "
            />
        </defs>

        <text 
          fill="black" 
          fontSize="16" 
          letterSpacing="3" 
          fontWeight="600"
        >
          <textPath href="#bottomArc" startOffset="50%" textAnchor="middle">
            VIEW REAL-TIME RESULTS
          </textPath>
        </text>
      </motion.svg>

    </div>

    </motion.section>
  );
}
