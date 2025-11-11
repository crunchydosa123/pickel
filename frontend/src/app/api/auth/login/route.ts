import { NextResponse } from "next/server";

export async function POST(req: Request) {
  try {
    const { email, password } = await req.json();
    const backendURL = process.env.BACKEND_URL;

    const res = await fetch(`${backendURL}/auth/login`, {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({ email, password }),
      credentials: "include",
    });

    const setCookieHeader = res.headers.get("set-cookie"); // <-- extract backend cookie
    const text = await res.text();
    let data: any;
    try {
      data = JSON.parse(text);
    } catch {
      data = { message: text };
    }

    const response = NextResponse.json(
      { message: "Login successful", data },
      { status: 200 }
    );

    if (setCookieHeader) {
      response.headers.set("Set-Cookie", setCookieHeader); // <-- forward it to browser
    }

    return response;
  } catch (error) {
    console.error("Error connecting to Go Backend:", error);
    return NextResponse.json(
      { error: "Internal Server Error" },
      { status: 500 }
    );
  }
}
