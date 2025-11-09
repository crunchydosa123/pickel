import { NextResponse } from "next/server";
import { writeFile } from "fs/promises";
import path from "path";


export async function POST(req: Request){
  const formData = await req.formData();
  const file = formData.get("file") as File;

  if(!file) return NextResponse.json({error: "No file uploaded"}, { status: 400});
  const res = await fetch('https://localhost:8080/')
  
  //proxy to go backend which stores the pickle file in s3
}