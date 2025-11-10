import { NextResponse } from "next/server";

export async function POST(req: Request) {
  try {
    const { name, email, password } = await req.json();
    const backendURL = process.env.BACKEND_URL;

    const res = await fetch(`${backendURL}/auth/signup`, {
      method: "POST",
      headers: {
        "Content-type": "application/json"
      },
      body: JSON.stringify({ name, email, password })
    });

    const data = await res.json().catch(() => ({}));
    const resStatus = res.status;

    if (resStatus === 409) {
      console.log("Conflict (email exists)");
      return NextResponse.json(
        { message: "Email already registered" },
        { status: 409 }
      );
    }

    if (resStatus === 201 || resStatus === 200) {
      console.log("User created successfully");
      return NextResponse.json(
        { message: "User created successfully", data },
        { status: 200 }
      );
    }

    console.error("Unexpected backend response:", resStatus, data);
    return NextResponse.json(
      { error: "Unexpected response from backend", data },
      { status: resStatus || 500 }
    );

  } catch (error) {
    console.error("Error connecting to Go Backend:", error);
    return NextResponse.json(
      { error: "Internal Server Error" },
      { status: 500 }
    );
  }
}
