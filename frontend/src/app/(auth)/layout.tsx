import React from 'react'
import { Toaster } from 'sonner'

type Props = {
  children: React.ReactNode
}

const layout = (props: Props) => {
  return (
    <div className='flex flex-col w-full bg-gray-200 h-screen justify-center items-center'>
      <div>{props.children}</div>
      <Toaster position="top-center" richColors />
    </div>
    
  )
}

export default layout