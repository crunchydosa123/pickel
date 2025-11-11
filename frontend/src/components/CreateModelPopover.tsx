"use client"

import { Button } from "@/components/ui/button";
import { Card, CardContent, CardDescription, CardFooter, CardTitle } from "@/components/ui/card";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { Popover, PopoverContent, PopoverTrigger } from "@/components/ui/popover";
import { ArrowRight } from "lucide-react";
import { useState } from "react";

type Props = {}

const CreateModelPopover = (props: Props) => {
  const [modelName, setModelName] = useState('');

  const baseUrl = process.env.NEXT_PUBLIC_BASE_URL || "http://localhost:3000";


  const createModel = async () => {
    const res = await fetch(`${baseUrl}/api/model/create`, {
      method: 'POST',
      headers: {
        'Content-type': 'application/json'
      },
      body: JSON.stringify({name: modelName}),
      credentials: "include"
    });
  }

  return (
    <Popover>
      <PopoverTrigger><Button>Create Model</Button></PopoverTrigger>
      <PopoverContent>
        <Card>
          <CardTitle>Create a new model</CardTitle>
          <CardDescription>Create a new model</CardDescription>
          <CardContent>
            <Label>Name</Label>
            <Input value={modelName} onChange={(e) => setModelName(e.target.value)}></Input>
          </CardContent>

          <CardFooter>
            <Button onClick={createModel}>Create Model <ArrowRight /></Button>
          </CardFooter>
        </Card>
      </PopoverContent>
    </Popover>
  )
}

export default CreateModelPopover