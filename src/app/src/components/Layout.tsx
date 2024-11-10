// src/components/Layout.tsx
import React from 'react'
import { Outlet } from 'react-router-dom'
import Header from './Header'

const Layout: React.FC = () => {
  return (
    <div className='grid place-content-center place-items-center'>
      <Header />
      <Outlet />
    </div>
  )
}

export default Layout
