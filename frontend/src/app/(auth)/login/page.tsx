'use client'

import React, { useState } from 'react'
import { useRouter } from 'next/navigation'
import { Button } from '@/components/ui/button'
import { Card, CardContent, CardDescription, CardFooter, CardTitle } from '@/components/ui/card'
import { Input } from '@/components/ui/input'
import { Label } from '@/components/ui/label'
import { ArrowRight } from 'lucide-react'
import { toast } from 'sonner'

const Page = () => {
  const router = useRouter()
  const [email, setEmail] = useState('')
  const [password, setPassword] = useState('')
  const [loading, setLoading] = useState(false)

  const onLogin = async () => {
    if (!email || !password) {
      toast.warning("Please enter both email and password")
      return
    }

    setLoading(true)
    try {
      const res = await fetch('/api/auth/login', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ email, password }),
        credentials: 'include'
      })

      const data = await res.json()

      if (!res.ok) {
        toast.error(data?.error || "Invalid credentials")
      } else {
        toast.success("Login successful")
        //localStorage.setItem('token', data.token)
        router.push('/dashboard')
      }
    } catch (err) {
      console.error('Login failed:', err)
      toast.error('Something went wrong. Please try again.')
    } finally {
      setLoading(false)
    }
  }

  return (
    <div className="flex justify-center mt-10">
      <Card className="w-300 max-w-md p-6">
        <CardTitle>
          <div className="text-3xl">Login</div>
        </CardTitle>
        <CardDescription>Welcome back</CardDescription>

        <CardContent className="mt-4 px-0">
          <div className="flex flex-col gap-2 mb-4">
            <Label htmlFor="email">Email</Label>
            <Input id="email" type="email" value={email} onChange={(e) => setEmail(e.target.value)} placeholder="Enter your email" />
          </div>

          <div className="flex flex-col gap-2 mb-4">
            <Label htmlFor="password">Password</Label>
            <Input id="password" type="password" value={password} onChange={(e) => setPassword(e.target.value)} placeholder="Enter your password" />
          </div>
        </CardContent>

        <CardFooter className="px-0">
          <Button className="w-full justify-center" onClick={onLogin} disabled={loading}>
            {loading ? "Logging in..." : <>Login <ArrowRight className="ml-2" /></>}
          </Button>
        </CardFooter>
      </Card>
    </div>
  )
}

export default Page
