"use client";

import React, { useState } from "react";
import { Button } from "./ui/button";
import { Card, CardFooter } from "./ui/card";
import {
  Popover,
  PopoverContent,
  PopoverTrigger,
} from "@/components/ui/popover";

const AddModelCode = () => {
  const [file, setFile] = useState<File | null>(null);

  const baseUrl = process.env.NEXT_PUBLIC_BASE_URL || "http://localhost:3000";

  const handleUpload = async () => {
    if (!file) return;

    const formData = new FormData();
    formData.append("modelFile", file);

    const res = await fetch(`${baseUrl}/api/model/deploy`, {
      method: "POST",
      body: formData,
      credentials: "include",
    });

    const data = await res.json();
    console.log("Upload result:", data);
  };

  return (
    <Popover>
      <PopoverTrigger asChild>
        <Button>Add Code</Button>
      </PopoverTrigger>

      <PopoverContent className="w-64 p-4">
        <Card className="p-3 space-y-3">
          <input
            type="file"
            className="w-full text-sm"
            onChange={(e) => setFile(e.target.files?.[0] ?? null)}
          />

          {file && (
            <p className="text-xs text-muted-foreground">
              Selected: {file.name}
            </p>
          )}

          <CardFooter className="flex justify-end">
            <Button onClick={handleUpload}>Upload Code</Button>
          </CardFooter>
        </Card>
      </PopoverContent>
    </Popover>
  );
};

export default AddModelCode;
