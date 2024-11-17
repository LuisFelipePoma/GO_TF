import React from 'react'

interface Props {
  children: React.ReactNode
  handleDirection: () => void
}

export const NavigationBtn: React.FC<Props> = ({
  children,
  handleDirection
}) => {
  return (
    <button
      onClick={handleDirection}
      className='px-4 py-1 w-fit text-light rounded-md 
                 bg-gray-200 bg-secondary/95 hover:text-white 
                 transition-transform duration-600  backdrop-blur-md
                 shadow-md hover:shadow-lg hover:bg-primary/95'
    >
      {children}
    </button>
  )
}
