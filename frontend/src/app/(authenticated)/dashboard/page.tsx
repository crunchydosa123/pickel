import CreateModelPopover from "@/components/CreateModelPopover";
import { Button } from "@/components/ui/button";
import { Card, CardContent, CardDescription, CardFooter, CardTitle } from "@/components/ui/card";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { Popover, PopoverContent, PopoverTrigger } from "@/components/ui/popover";
import { ArrowRight } from "lucide-react";

export default async function Page() {
  const baseUrl = process.env.NEXT_PUBLIC_BASE_URL || "http://localhost:3000";
  

  const res = await fetch(`${baseUrl}/api/hello`, {
    cache: "no-store",
  });
  const data = await res.json();

  return (
    <div className="w-full bg-blue-400">
      <h1>Message from API:</h1>
      <p>{data.message}</p>
      <CreateModelPopover />
    </div>
  );
}
