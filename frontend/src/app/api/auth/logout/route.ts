import { NextResponse } from "next/server";

export async function POST() {
  // Delete the auth cookie
  const res = NextResponse.json({ message: "Logged out successfully" });
  res.cookies.set("token", "", { path: "/", maxAge: 0 });
  return res;
}
