// src/components/Layout.tsx
import React from 'react'
import { Outlet } from 'react-router-dom'
import Header from './Header'

const Layout: React.FC = () => {
  return (
    <div className='grid place-items-center w-[100vw] h-[100vh] px-[10vw]'>
      <Header />
      <Outlet />
    </div>
  )
}

export default Layout
