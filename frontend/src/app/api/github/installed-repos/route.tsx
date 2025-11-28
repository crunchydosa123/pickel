import { NextRequest, NextResponse } from "next/server";

// route.tsx
export async function GET(req: NextRequest) {
  try {
    const installationId = req.nextUrl.searchParams.get("installation_id"); //TODO: get it from the frontend or keep it hardcoded for how
    const backendURL = process.env.BACKEND_URL;

    const res = await fetch(`${backendURL}/github/fetch-repos?installation_id=${installationId}`, {
      method: "GET",
      headers: { Cookie: req.headers.get("cookie") || "" },
    });

    const text = await res.text();
    let data;
    try {
      data = JSON.parse(text);
    } catch {
      console.warn("Backend returned invalid JSON:", text);
      data = { repositories: [] };
    }

    return NextResponse.json(data);
  } catch (err) {
    console.error("Error proxying request:", err);
    return NextResponse.json({ repositories: [] }, { status: 500 });
  }
}
