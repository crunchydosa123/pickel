import { NextResponse } from "next/server";

export async function POST(req: Request) {
  try {
    const { email, password } = await req.json();
    const backendURL = process.env.BACKEND_URL!;

    const backendRes = await fetch(`${backendURL}/auth/login`, {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({ email, password }),
      credentials: "include",
    });

    const cookie = backendRes.headers.get("set-cookie");
    const text = await backendRes.text();
    let data: any;

    try {
      data = JSON.parse(text);
    } catch {
      data = { message: text };
    }

    if (!backendRes.ok) {
      return NextResponse.json(data, { status: backendRes.status });
    }

    const response = NextResponse.json(data, { status: 200 });

    if (cookie) {
      // 🔥 Forward cookie exactly as received from backend
      response.headers.set("Set-Cookie", cookie);
    }

    return response;
  } catch (err) {
    console.error("Error connecting to backend:", err);
    return NextResponse.json(
      { error: "Internal Server Error" },
      { status: 500 }
    );
  }
}
