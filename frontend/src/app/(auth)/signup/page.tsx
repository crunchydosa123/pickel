import React from 'react'
import { Button } from '@/components/ui/button'
import { Card, CardContent, CardDescription, CardFooter, CardTitle } from '@/components/ui/card'
import { Input } from '@/components/ui/input'
import { Label } from '@/components/ui/label'
import { ArrowRight } from 'lucide-react'

type Props = {}

const Page = (props: Props) => {
  return (
    <div className="flex justify-center mt-10">
      <Card className="w-300 max-w-md p-6">
        <CardTitle>
          <div className='text-3xl'>Signup</div>
        </CardTitle>
        <CardDescription>Create a new account</CardDescription>

        <CardContent className="mt-4 px-0">
          <div className="flex flex-col gap-2 mb-4">
            <Label htmlFor="email">Email</Label>
            <Input id="email" type="email" placeholder="Enter your email" />
          </div>

          <div className="flex flex-col gap-2 mb-4">
            <Label htmlFor="password">Password</Label>
            <Input id="password" type="password" placeholder="Enter your password" />
          </div>

          <div className="flex flex-col gap-2 mb-4">
            <Label htmlFor="confirmPassword">Confirm Password</Label>
            <Input id="confirmPassword" type="password" placeholder="Confirm your password" />
          </div>
        </CardContent>

        <CardFooter className='px-0'>
          <Button className="w-full justify-center">
            Create Account <ArrowRight className="ml-2" />
          </Button>
        </CardFooter>
      </Card>
    </div>
  )
}

export default Page
