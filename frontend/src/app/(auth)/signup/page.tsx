"use client"

import React, { useState } from 'react'
import { Button } from '@/components/ui/button'
import { Card, CardContent, CardDescription, CardFooter, CardTitle } from '@/components/ui/card'
import { Input } from '@/components/ui/input'
import { Label } from '@/components/ui/label'
import { ArrowRight } from 'lucide-react'
import { toast } from 'sonner'

const Page = () => {
  const [name, setName] = useState('')
  const [email, setEmail] = useState('')
  const [password1, setPassword1] = useState('')
  const [password2, setPassword2] = useState('')

  const Signup = async () => {
    try {
      if (password1 !== password2) {
        toast.error("Passwords do not match")
        return
      }

      const res = await fetch('/api/auth/signup', {
        method: "POST",
        headers: { "Content-type": "application/json" },
        body: JSON.stringify({ name, email, password: password1 })
      })

      console.log("signup page res: ", res)

      if (res.ok) {
        toast.success("User created successfully 🎉")
        setName('')
        setEmail('')
        setPassword1('')
        setPassword2('')
      } else if (res.status === 409) {
        toast.warning("User with this email already exists ⚠️")
      } else {
        toast.error("Something went wrong 😞")
      }

    } catch (error) {
      console.error(error)
      toast.error("Network error. Please try again.")
    }
  }

  return (
    <div className="flex justify-center mt-10">
      <Card className="w-300 max-w-md p-6">
        <CardTitle>
          <div className='text-3xl'>Signup</div>
        </CardTitle>
        <CardDescription>Create a new account</CardDescription>

        <CardContent className="mt-4 px-0">
          <div className="flex flex-col gap-2 mb-4">
            <Label htmlFor="name">Full Name</Label>
            <Input id="name" type="text" placeholder="Enter your name"
              value={name} onChange={(e) => setName(e.target.value)} />
          </div>

          <div className="flex flex-col gap-2 mb-4">
            <Label htmlFor="email">Email</Label>
            <Input id="email" type="email" placeholder="Enter your email"
              value={email} onChange={(e) => setEmail(e.target.value)} />
          </div>

          <div className="flex flex-col gap-2 mb-4">
            <Label htmlFor="password">Password</Label>
            <Input id="password" type="password" placeholder="Enter your password"
              value={password1} onChange={(e) => setPassword1(e.target.value)} />
          </div>

          <div className="flex flex-col gap-2 mb-4">
            <Label htmlFor="confirmPassword">Confirm Password</Label>
            <Input id="confirmPassword" type="password" placeholder="Confirm your password"
              value={password2} onChange={(e) => setPassword2(e.target.value)} />
          </div>
        </CardContent>

        <CardFooter className='px-0'>
          <Button className="w-full justify-center" onClick={Signup}>
            Create Account <ArrowRight className="ml-2" />
          </Button>
        </CardFooter>
      </Card>
    </div>
  )
}

export default Page
