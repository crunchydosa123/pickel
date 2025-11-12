"use client";

import { Button } from "@/components/ui/button";
import { Card } from "@/components/ui/card";
import { Settings, LogOut } from "lucide-react";
import { useRouter } from "next/navigation";

export default function UserCard() {
  const router = useRouter();

  const handleLogout = async () => {
    try {
      const res = await fetch("/api/auth/logout", {
        method: "POST",
        credentials: "include",
      });

      if (res.ok) {
        localStorage.removeItem("token"); // optional: clear any local tokens
        router.push("/login"); // redirect to login page
      } else {
        console.error("Logout failed");
      }
    } catch (err) {
      console.error("Error logging out:", err);
    }
  };

  return (
    <Card className="flex justify-between px-2 py-2 items-center">
      <div className="flex items-center gap-2">
        <div className="font-medium">User</div>
        <Settings className="cursor-pointer" />
      </div>
      <Button
        variant="ghost"
        size="icon"
        onClick={handleLogout}
        title="Logout"
      >
        <LogOut />
      </Button>
    </Card>
  );
}
