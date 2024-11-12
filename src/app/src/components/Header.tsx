// src/components/Header.tsx
import React from 'react'
import { Link } from 'react-router-dom'
import { Icon } from '../assets/Icon'

const Header: React.FC = () => {
  return (
    <header className='h-[20vh] flex items-center place-self-start'>
      <Link
        to='/'
        className=' p-0 m-0 flex items-center h-[15px] justify-center group'
      >
        <Icon className='w-[50px] text-secondary stroke-1 h-[50px]' />
        <span className='text-[54px] group-hover:text-primary  transition-all duration-700 ease-in-out'>
          Movies Recomender
        </span>
      </Link>
    </header>
  )
}

export default Header
