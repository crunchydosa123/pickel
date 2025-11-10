import { NextResponse } from "next/server";

export async function POST(req: Request) {
  try {
    const { email, password } = await req.json();
    const backendURL = process.env.BACKEND_URL;

    if (!email || !password) {
      return NextResponse.json(
        { error: "Email and password are required" },
        { status: 400 }
      );
    }

    const res = await fetch(`${backendURL}/auth/login`, {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
      },
      body: JSON.stringify({ email, password }),
    });

    const text = await res.text(); // Read the raw response text (in case it's not valid JSON)
    let data: any;
    try {
      data = JSON.parse(text);
    } catch {
      data = { message: text }; // fallback if backend sends plain text (like "Invalid credentials")
    }

    // Handle invalid credentials
    if (res.status === 401) {
      return NextResponse.json(
        { error: "Invalid email or password" },
        { status: 401 }
      );
    }

    // Handle token success
    if (res.status === 200 && data.token) {
      return NextResponse.json(
        { message: "Login successful", token: data.token },
        { status: 200 }
      );
    }

    // Handle unexpected responses
    console.error("Unexpected backend response:", res.status, data);
    return NextResponse.json(
      { error: "Unexpected backend response", data },
      { status: res.status || 500 }
    );

  } catch (error) {
    console.error("Error connecting to Go Backend:", error);
    return NextResponse.json(
      { error: "Internal Server Error" },
      { status: 500 }
    );
  }
}
