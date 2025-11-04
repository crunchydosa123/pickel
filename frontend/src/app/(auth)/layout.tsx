import React from 'react'

type Props = {
  children: React.ReactNode
}

const layout = (props: Props) => {
  return (
    <div className='flex-col'>
      <div>hello</div>
      <div>{props.children}</div>
    </div>
    
  )
}

export default layout