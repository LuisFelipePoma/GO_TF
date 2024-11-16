// src/components/Layout.tsx
import React, { useEffect, useRef, useState } from 'react'
import { Outlet } from 'react-router-dom'
import { Header } from '../components/Header'
import { Response } from '../types/movies'
import { Popup, PopupData } from '../components/modals/PopUp'
import { useStore } from '../services/store'
import { URL_IMG, URL_WS } from '../consts/api'
import { AnimatePresence, motion } from 'framer-motion'
import Config from '../components/modals/Config'
import { IconConfig } from '../assets/IconConfig'

const Layout: React.FC = () => {
  const [popup, setPopup] = useState<PopupData | null>(null)
  // eslint-disable-next-line @typescript-eslint/no-explicit-any
  const backgroundPath = useStore((state: any) => state.backgroundPath)
  // eslint-disable-next-line @typescript-eslint/no-explicit-any
  const nMoviesUser = useStore((state: any) => state.nMoviesUser)
  const [isSidebarOpen, setIsSidebarOpen] = useState(false)
  const connectionWS = useRef<WebSocket | null>(null)
  const [isPopupVisible, setIsPopupVisible] = useState(false)

  const toggleSidebar = () => {
    setIsSidebarOpen(!isSidebarOpen)
  }

  const togglePopupVisibility = () => {
    setIsPopupVisible(!isPopupVisible)
  }

  useEffect(() => {
    const socket = new WebSocket(URL_WS)
    connectionWS.current = socket

    socket.onopen = () => {
      console.log('Conexión WebSocket establecida')
      socket.send(JSON.stringify({ n: nMoviesUser }))
    }

    socket.onmessage = event => {
      const data: Response = JSON.parse(event.data)
      const showPopupTimer = setTimeout(() => {
        setPopup({
          message: data.movie_response || [],
          id: data.target_movie ?? '0'
        })
      }, 3000)

      return () => clearTimeout(showPopupTimer)
    }

    socket.onclose = () => {
      console.log('Conexión WebSocket cerrada')
    }

    socket.onerror = error => {
      console.error('Error en WebSocket:', error)
    }

    return () => {
      socket.close()
    }
  }, [nMoviesUser])

  useEffect(() => {
    if (
      connectionWS.current &&
      connectionWS.current.readyState === WebSocket.OPEN
    ) {
      connectionWS.current.send(JSON.stringify({ n: nMoviesUser }))
      console.log('Enviado nuevo valor para nMoviesUser:', nMoviesUser)
    }
  }, [nMoviesUser])

  return (
    <div className='grid place-items-center w-[100vw] h-[100vh] px-[10vw] relative'>
      {!isSidebarOpen && (
        <button
          onClick={toggleSidebar}
          className='fixed top-4 left-4 text-white transform transition-transform duration-300 hover:rotate-180 hover:scale-125'
        >
          <IconConfig />
        </button>
      )}

      <Config isOpen={isSidebarOpen} onClose={toggleSidebar} />

      <div className='absolute inset-0 -z-10 overflow-hidden'>
        <AnimatePresence>
          <motion.img
            initial={{ opacity: 0.75, scale: 1.15 }}
            animate={{
              opacity: [0.6, 0.75, 0.6],
              scale: [1.1, 1, 1.1]
            }}
            transition={{
              duration: 25,
              repeat: Infinity,
              repeatType: 'mirror',
              ease: 'easeInOut'
            }}
            className='w-full h-full object-cover filter blur-[0.85px] absolute top-0 left-0'
            src={
              backgroundPath ? URL_IMG(backgroundPath, 'original') : '/bg.webp'
            }
            alt='Background'
          />
        </AnimatePresence>
        <div className='absolute inset-0 bg-gradient-to-b from-[#0B0000]/75 via-transparent to-transparent' />
        <div className='absolute inset-0 bg-gradient-to-t from-dark to-transparent' />
      </div>

      <Header />
      <Outlet />

      <AnimatePresence>
        {popup && (
          <>
            <Popup
              data={popup}
              isVisible={isPopupVisible}
              setIsVisible={setIsPopupVisible}
            />
            {!isPopupVisible && (
              <motion.button
                onClick={togglePopupVisibility}
                className='fixed top-1/2 right-0 transform -translate-y-1/2 bg-primary text-white py-4 px-3 rounded-l'
                aria-label='Toggle Popup'
                initial={{ opacity: 0, x: 20 }}
                animate={{ opacity: 1, x: 0 }}
                exit={{ opacity: 0, x: -20 }}
                transition={{ duration: 1.5 }}
              >
                <svg
                  xmlns='http://www.w3.org/2000/svg'
                  viewBox='0 0 24 24'
                  fill='none'
                  stroke='currentColor'
                  strokeLinecap='round'
                  strokeLinejoin='round'
                  width={24}
                  height={24}
                  strokeWidth={2}
                >
                  <path d='M13 17l-5-5 5-5' />
                </svg>
              </motion.button>
            )}
          </>
        )}
      </AnimatePresence>
    </div>
  )
}

export default Layout
