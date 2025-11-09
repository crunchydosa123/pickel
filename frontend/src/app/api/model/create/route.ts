import { NextResponse } from "next/server";

export async function POST(request: Request) {
  try {
    const body = await request.json()

    const { name } = body;

    const res = await fetch("http://localhost:8080/models/create", {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
      },
      body: JSON.stringify({ name })
    });

    const data = await res.json();
    return NextResponse.json(data, { status: res.status })
  }catch(error){
    console.error("Error connecting to Go Backend: ", error);
    return NextResponse.json({error: "Internal Server Error: "}, {status: 500});
  }
  
}