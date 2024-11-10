// src/components/Header.tsx
import React from 'react'
import { Link } from 'react-router-dom'

const Header: React.FC = () => {
  return (
    <header>
      <h1 className='p-0 m-0 w-full'>
        <Link to='/' className='w-fit h-fit hover:text-primary'>
          Movies Recomender
        </Link>
      </h1>
    </header>
  )
}

export default Header
