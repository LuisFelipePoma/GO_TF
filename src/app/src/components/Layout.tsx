// src/components/Layout.tsx
import React, { useEffect, useState } from 'react'
import { Outlet } from 'react-router-dom'
import { Header } from './Header'
import { Response } from '../types/movies'
import { Popup, PopupData } from './PopUp'
import { useStore } from '../services/store'
import { URL_IMG, URL_WS } from '../consts/api'
import { AnimatePresence, motion } from 'framer-motion'

const Layout: React.FC = () => {
  const [popup, setPopup] = useState<PopupData | null>(null)
  // eslint-disable-next-line @typescript-eslint/no-explicit-any
  const backgroundPath = useStore((state: any) => state.backgroundPath)

  useEffect(() => {
    const socket = new WebSocket(URL_WS)

    socket.onopen = () => {
      console.log('Conexión WebSocket establecida')
    }

    socket.onmessage = event => {
      const data: Response = JSON.parse(event.data)
      // Mostrar el popup después de 3 segundos
      const showPopupTimer = setTimeout(() => {
        setPopup({
          message: data.movie_response || [],
          id: data.target_movie ?? '0'
        })
      }, 3000)

      // Limpieza del showPopupTimer si el componente se desmonta antes
      return () => clearTimeout(showPopupTimer)
    }

    socket.onclose = () => {
      console.log('Conexión WebSocket cerrada')
    }

    socket.onerror = error => {
      console.error('Error en WebSocket:', error)
    }

    // Limpieza al desmontar el componente
    return () => {
      socket.close()
    }
  }, [])

  const closePopup = () => {
    setPopup(null)
  }

  return (
    <div className='grid place-items-center w-[100vw] h-[100vh] px-[10vw] relative'>
      <div className='absolute inset-0 -z-10 overflow-hidden'>
        <AnimatePresence>
          <motion.img
            initial={{ opacity: 0.6, scale: 1.15 }}
            animate={{
              opacity: [0.6, 0.65, 0.6],
              scale: [1.1, 1, 1.1]
            }}
            transition={{
              duration: 15, // Aumenta la duración para un bombeo más lento
              repeat: Infinity, // Repetir indefinidamente
              repeatType: 'mirror', // Alternar entre adelante y atrás
              ease: 'easeInOut'
            }}
            className='w-full h-full object-cover filter blur-[0.85px] absolute top-0 left-0'
            src={
              backgroundPath ? URL_IMG(backgroundPath, 'original') : '/bg.webp'
            }
            alt='Background'
          />
        </AnimatePresence>
        {/* Gradient overlay */}
        <div className='absolute inset-0 bg-gradient-to-b from-[#0B0000]/75 via-transparent to-transparent' />
        <div className='absolute inset-0 bg-gradient-to-t from-dark to-transparent' />
      </div>
      <Header />
      <Outlet />
      {popup && <Popup data={popup} onClose={closePopup} />}
    </div>
  )
}

export default Layout
