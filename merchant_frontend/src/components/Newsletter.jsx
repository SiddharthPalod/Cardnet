"use client"

import Image from "next/image"
import { useState } from "react"

export default function Newsletter() {
  const [email, setEmail] = useState("")

  const mailtoLink = `mailto:siddharthpalod@gmail.com?subject=Newsletter Subscription&body=Subscribe: ${email}`

  return (
    <section className="relative max-w-xl w-full bg-white outline outline-2 outline-black overflow-hidden rounded-md">

        <div className="w-full h-48 bg-lime-400 flex items-center justify-center">
            <div className="relative w-28 h-20 sm:w-32 sm:h-24">
                <Image
                src="/card.svg"
                alt="Card UI"
                fill
                className="object-contain"
                priority
                />
            </div>
        </div>

      {/* Content */}
      <div className="relative px-10 pt-10 pb-6 flex flex-col gap-5">

        <div className="flex flex-row gap-10">
        <h3 className="text-black uppercase text-sm sm:text-lg font-medium font-poppins ">
          Subscribe to our newsletter
        </h3>
        <p className="text-black uppercase text-xs sm:text-sm font-poppins leading-relaxed opacity-80">
          Stay connected by joining our newsletter getting weekly updates on our latest innovations
        </p>
        </div>

        {/* Input Row */}
        <div className="flex items-center w-full border border-stone-300 bg-white">

          <input
            type="email"
            placeholder="EMAIL ADDRESS"
            value={email}
            onChange={(e) => setEmail(e.target.value)}
            className="flex-1 px-4 py-3 text-black font-poppins text-xs sm:text-sm outline-none placeholder:text-stone-400"
          />

          <a
            href={mailtoLink}
            className="w-12 h-12 bg-black flex items-center justify-center hover:bg-neutral-900 transition"
          >
            <span className="block w-3 h-4 bg-white clip-arrow" />
          </a>

        </div>

      </div>

      {/* Arrow Shape */}
      <style jsx>{`
        .clip-arrow {
          clip-path: polygon(0 0, 100% 50%, 0 100%, 20% 50%);
        }
      `}</style>

    </section>
  )
}
