"use client"

import Image from "next/image"
import { useState } from "react"
import Link from "next/link"

const FEATURES = [
  {
    title: "High-Performance Go Backend",
    desc: [
      "LOW LATENCY",
      "HIGH CONCURRENCY (GOROUTINES)",
      "EFFICIENT MEMORY USAGE",
    ],
  },
  {
    title: "Modular Network Architecture",
    desc: [
      "BANK LOGIC",
      "MERCHANT PROCESSING",
      "NETWORK ROUTER",
      "LEDGER / SETTLEMENT MODULE",
    ],
  },
  {
    title: "Built-in Fraud Detection Layer",
    desc: [
      "VELOCITY CHECKS",
      "SPEND PATTERN ANALYSIS",
      "GEO MISMATCH DETECTION",
      "RISK THRESHOLDS",
    ],
  },
]

function FeatureMap() {
  const [activeIndex, setActiveIndex] = useState(null)

  return (
    <div className="space-y-10 max-w-xs">
      {FEATURES.map((item, index) => {
        const active = activeIndex === index

        return (
          <div
            key={index}
            onMouseEnter={() => setActiveIndex(index)}
            onMouseLeave={() => setActiveIndex(null)}
            onClick={() => setActiveIndex(active ? null : index)}
            className="cursor-pointer space-y-3 transition-all"
          >
            {/* Title Row */}
            <div className="flex items-start gap-4">

              {/* Dash slot (always takes space) */}
              <span
                className={`mt-2 block w-8 h-0.5 transition-all duration-300 ${
                  active ? "bg-black opacity-100" : "bg-transparent opacity-0"
                }`}
              />

              {/* Title */}
              <p
                className={`uppercase text-lg sm:text-xl lg:text-2xl font-medium font-poppins transition-opacity duration-200 ${
                  active ? "opacity-100 text-black" : "opacity-80 text-black"
                }`}
              >
                {item.title}
              </p>
            </div>

            {/* Expandable Description */}
            <div
              className={`ml-9 overflow-hidden transition-all duration-300 ease-in-out ${
                active ? "max-h-40 opacity-100" : "max-h-0 opacity-0"
              }`}
            >
              <ul className="space-y-2 text-sm sm:text-base text-black font-poppins uppercase tracking-wide">
                {item.desc.map((line, i) => (
                  <li key={i} className="opacity-80">
                    {line}
                  </li>
                ))}
              </ul>
            </div>

          </div>
        )
      })}
    </div>
  )
}

export default function HeroTransactions() {
return (
    <section className="relative bg-white overflow-hidden px-6 py-20 lg:px-20">

        {/* Main Heading */}
        <h1 className="text-center text-black font-extrabold uppercase font-poppins text-4xl sm:text-5xl lg:text-7xl leading-tight max-w-5xl mx-auto">
            Best Partner For Your Transactions
        </h1>

        {/* Content Grid */}
    <div className="my-16 flex flex-row items-start gap-10">

    {/* Feature List — Wider */}
    <div className="flex-[1.2]">
            <FeatureMap />
    </div>

    {/* Image Block — Narrower */}
    <div className="flex-[1.4] bg-lime-400 p-3 min-h-[420px] flex items-center justify-center">
            <div className="relative w-full max-w-[560px] h-[400px] bg-white border-8 border-white overflow-hidden mx-auto">
            <Image
                    src="/features.webp"
                    alt="Transaction Visual"
                    fill
                    className="object-contain"
                    priority
            />
            </div>
    </div>

    </div>

        {/* Bottom Section */}
        <div className="mt-24 grid grid-cols-1 lg:grid-cols-2 gap-12 items-center">

            {/* History Card Stack (reduced size) */}
            <div className="w-[420px] h-[360px] relative">
                <div className="w-72 h-72 left-[11%] top-[44%] absolute origin-top-left rotate-[-33.48deg] bg-black" />
                <div className="w-72 h-72 left-0 top-0 absolute">
                    <Link href="/history">
                        <div className="w-72 h-72 left-0 top-0 absolute bg-lime-400" />
                        <div className="absolute inset-0 flex flex-col justify-center items-center gap-3 px-6">
                            <div className="self-stretch text-black text-3xl font-medium uppercase font-poppins">
                                    VIEW HISTORY
                            </div>
                            <div className="w-full text-black text-xl font-normal uppercase font-poppins">BROWSE PAST AUTH <br/> REQUESTS AND THEIR<br/> OUTCOMES</div>
                        </div>
                    </Link>
                </div>
            </div>

            <h2 className="text-black uppercase font-extrabold font-poppins text-4xl sm:text-5xl lg:text-7xl leading-tight text-center lg:text-left">
                    Real-Time <br /> Updates
            </h2>

        </div>
    </section>
)
}
