"use client";
import { motion } from "framer-motion";
import Image from "next/image";
import { useRouter, usePathname } from "next/navigation";
import { useState } from "react";
import { Search } from "lucide-react";
import Link from "next/link";

export default function Navbar() {
  const router = useRouter();
  const pathname = usePathname();

  const [merchantId, setMerchantId] = useState("");

  const handleSearch = () => {
    if (!merchantId.trim()) return;

    if (!pathname.startsWith("/merchants")) {
      router.push(`/merchants/${merchantId}`);
      return;
    }

    window.dispatchEvent(
      new CustomEvent("merchant-search", {
        detail: merchantId,
      })
    );
  };

  return (
    <motion.nav 
      initial={{ y: -40, opacity: 0 }}
      animate={{ y: 0, opacity: 1 }}
      transition={{ duration: 0.6 }}
      className="flex items-center justify-between px-6 md:px-16 py-6 bg-[#B5FF37]"
    >
      {/* Brand */}
      <div className="flex items-center gap-6">
        <motion.div
          initial={{ scaleY: 0 }}
          animate={{ scaleY: 1 }}
          transition={{ duration: 0.6 }}
          className="origin-top"
        >
          <Image
            src="/circle.svg"
            alt="Real-Time Indicator"
            width={56}
            height={56}
            className="w-14"
          />
        </motion.div>

        <span className="text-black text-4xl font-['Pacifico']">
          <Link href="/">CardNet</Link>
        </span>
      </div>

      {/* Search */}
      <form onSubmit={(e) => { e.preventDefault(); handleSearch();}}
        className="hidden md:flex items-center gap-2 bg-white border-2 rounded-xl px-2 focus-within:border-b-4 focus-within:border-r-4 focus-within:border-black transition"
      >
        <input
          type="text"
          className="px-3 py-2 font-['Poppins'] text-sm w-md bg-transparent outline-none placeholder:text-gray-500"
          placeholder="Search Merchant ID (e.g., amazon, walmart, target) "
          value={merchantId}
          onChange={(e) => setMerchantId(e.target.value.toLowerCase())}
        />

        <button
          type="submit"
          disabled={!merchantId.trim()}
          className="p-2 rounded-xl transition active:scale-95 disabled:opacity-40 disabled:cursor-not-allowed"
          aria-label="Search"
        >
          <Search size={18} />
        </button>
      </form>

      <span className="text-black text-xl font-['Poppins'] uppercase">
        <Link href="/history">History</Link>
      </span>

      <Link
        href="http://localhost:3000"
        className="inline-flex items-center justify-center bg-black px-8 py-3 text-white text-sm md:text-base font-bold font-['Poppins'] uppercase hover:scale-[1.02] active:scale-95 transition-transform"
      >
        Merchant Portal
      </Link>
    

    </motion.nav>
  );
}
