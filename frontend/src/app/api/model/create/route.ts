import { NextResponse } from "next/server";

export async function POST(request: Request) {
  try {
    const backendURL = process.env.BACKEND_URL;
    const body = await request.json();
    const { name } = body;

    const cookieHeader = request.headers.get("cookie") || "";

    const res = await fetch(`${backendURL}/model/create`, {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
        "Cookie": cookieHeader,
      },
      credentials: "include",
      body: JSON.stringify({ name }),
    });

    return NextResponse.json({ status: res.status });
  } catch (error) {
    console.error("Error connecting to Go Backend: ", error);
    return NextResponse.json({ error: "Internal Server Error" }, { status: 500 });
  }
}
