"use client"

import CardNav from "@/components/CardNav";
import Particles from "@/components/Particles";
import Squares from "@/components/Squares";


export default function Home() {
  return (
    <div className="h-screen w-full bg-black">
    <Particles
    particleColors={['rgba(188, 188, 188, 0)', 'rgba(181, 95, 95, 0)']}
    particleCount={500}
    particleSpread={6}
    speed={0.1}
    particleBaseSize={100}
    moveParticlesOnHover={true}
    alphaParticles={false}
    disableRotation={false}
  />
  </div>
  );
}
