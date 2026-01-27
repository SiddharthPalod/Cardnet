import type { Metadata } from "next";
import { Poppins, Pacifico } from "next/font/google";
import "./globals.css";

// Poppins (Primary UI Font)
const poppins = Poppins({
  subsets: ["latin"],
  weight: ["400", "500", "600", "700", "800"],
  variable: "--font-poppins",
});

// Pacifico (Brand / Logo Font)
const pacifico = Pacifico({
  subsets: ["latin"],
  weight: "400",
  variable: "--font-pacifico",
});

export const metadata: Metadata = {
  title: "CardNet Merchant Simulator",
  description: "Simulate merchant card authorization requests in real-time",
};

export default function RootLayout({
  children,
}: {
  children: React.ReactNode;
}) {
  return (
    <html lang="en">
      <body
        className={`${poppins.variable} ${pacifico.variable} antialiased`}
      >
        {children}
      </body>
    </html>
  );
}
