import { NextResponse } from "next/server";

export async function GET(req: Request) {
  const backendURL = process.env.BACKEND_URL || "http://localhost:8080";
  const cookieHeader = req.headers.get("cookie") || "";

  try {
    console.log("Fetching models for user...");

    const res = await fetch(`${backendURL}/model/`, {
      method: "GET",
      headers: {
        "Content-Type": "application/json",
        Cookie: cookieHeader,
      },
      credentials: "include",
    });

    if (!res.ok) {
      console.error("Backend returned error:", res.status, res.statusText);
      return NextResponse.json(
        { error: `Backend returned ${res.status}` },
        { status: res.status }
      );
    }

    // ✅ Parse the actual JSON body
    const data = await res.json();

    console.log("Fetched models:", data);

    // ✅ Return that JSON to frontend
    return NextResponse.json(data);
  } catch (error: any) {
    console.error("Error fetching models:", error);
    return NextResponse.json(
      { error: "Failed to fetch models", details: error.message },
      { status: 500 }
    );
  }
}
