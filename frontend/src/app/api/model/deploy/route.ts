import { NextResponse } from "next/server";
import { writeFile } from "fs/promises";
import path from "path";


export async function POST(req: Request){
  const backendUrl = process.env.BACKEND_URL
  const cookieHeader = req.headers.get("cookie") || "";

  try{
    const formData = await req.formData();
    const file = formData.get('modelFile') as File;

    if(!file){
      return NextResponse.json({error: "No file uploaded"}, {status: 400});
    }

    const body = new FormData();
    body.append("modelFile", file);

    const res = await fetch(`${backendUrl}/model/deploy`, {
      method: "POST",
      headers: {
        Cookie: cookieHeader,
      },
      body, 
      credentials: "include"
    });
    
    const data = await res.json();
    return NextResponse.json(data);
  }catch(err: any){
    console.error("Error deploying model: ", err);
    return NextResponse.json(
    { error: "Failed to deploy model", details: err.message },  
    { status: 500 }
    );
  }
}