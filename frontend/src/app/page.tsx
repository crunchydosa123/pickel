"use client";

import CardNav from "@/components/CardNav";
import Particles from "@/components/Particles";

export default function Home() {
  return (
    <div className="relative w-full min-h-screen overflow-hidden">

      {/* PARTICLES BACKGROUND */}
      <div className="absolute inset-0 -z-10">
        <Particles
          particleColors={["#F35D40", "#F35D40"]}
          particleCount={500}
          particleSpread={6}
          speed={0.1}
          particleBaseSize={100}
          moveParticlesOnHover={true}
          alphaParticles={false}
          disableRotation={false}
        />
      </div>

      {/* NAVBAR */}
      <div className="relative z-50">
        <CardNav
          logo="icon.png"
          logoAlt="Company Logo"
          items={[]}
          baseColor="#fff"
          menuColor="#000"
          buttonBgColor="#111"
          buttonTextColor="#fff"
          ease="power3.out"
        />
      </div>

      {/* HERO CONTENT */}
      <div className="
        relative z-40 
        flex flex-col items-start
        px-6 sm:px-12 md:ml-70
        pt-40 md:pt-56
        max-w-5xl
      ">
        <div className="rounded-md bg-white/30 backdrop-blur-md border border-gray-500 py-2 px-4 text-lg">
          Welcome to Pickel
        </div>

        <div className="mt-6 rounded-full border-gray-500 p-2 text-4xl sm:text-6xl md:text-7xl font-semibold">
          Vercel for your ML Models
        </div>
        <div className="mt-2 rounded-full border-gray-500 p-2 text-2xl sm:text-md md:text-md font-semibold">
          Part of a portfolio. For demonstration purposes only.
        </div>
      </div>
    </div>
  );
}
