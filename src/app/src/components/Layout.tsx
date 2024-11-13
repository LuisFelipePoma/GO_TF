// src/components/Layout.tsx
import React, { useEffect, useState } from 'react'
import { Outlet } from 'react-router-dom'
import Header from './Header'
import { Response } from '../types/movies'
import { Popup, PopupData } from './PopUp'
import { useStore } from '../services/store'
import { URL_IMG, URL_WS } from '../consts/api'

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

        // Ocultar el popup después de 10 segundos
        const hidePopupTimer = setTimeout(() => {
          setPopup(null)
        }, 10000)

        // Limpieza del hidePopupTimer si el componente se desmonta antes
        return () => clearTimeout(hidePopupTimer)
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
      {/* Background image with gradient overlay */}
      {backgroundPath && (
        <div className='absolute inset-0 -z-10'>
          <img
            className='w-full h-full object-cover opacity-40 filter blur-[0.5px]'
            src={URL_IMG(backgroundPath, 'original')}
            alt=''
          />
          {/* Gradient overlay */}
          <div className='absolute inset-0 bg-gradient-to-t from-dark via-dark/50 via-dark/50 to-transparent' />
        </div>
      )}

      <Header />
      <Outlet />
      {popup && <Popup data={popup} onClose={closePopup} />}
    </div>
  )
}

export default Layout
