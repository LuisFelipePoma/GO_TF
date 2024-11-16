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
      className='px-3 py-1 w-fit text-light rounded-md 
                 bg-gray-200 bg-secondary hover:text-white 
                 transition-transform duration-300 
                 transform hover:scale-105 
                 shadow-md hover:shadow-lg'
    >
      {children}
    </button>
  )
}
