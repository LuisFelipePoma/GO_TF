// src/components/Layout.tsx
import React, { useEffect, useState } from 'react'
import { Outlet } from 'react-router-dom'
import Header from './Header'
import { MovieResponse, Response } from '../types/movies'

interface PopupData {
  title: string
  message: MovieResponse[]
}

const Popup: React.FC<{ data: PopupData; onClose: () => void }> = ({
  data,
  onClose
}) => (
  <div className='fixed top-2 right-2 bg-black bg-opacity-50 flex justify-center items-center'>
    <div className='bg-transparent backdrop-blur-lg p-4 rounded'>
      <h2 className='text-xl font-bold'>{data.title}</h2>
      {data.message.map(movie => (
        <p className='mt-2'>{movie.id}</p>
      ))}
      <button
        className='mt-4 px-4 py-2 bg-blue-500 text-white rounded'
        onClick={onClose}
      >
        Cerrar
      </button>
    </div>
  </div>
)

const Layout: React.FC = () => {
  const [popup, setPopup] = useState<PopupData | null>(null)

  useEffect(() => {
    const socket = new WebSocket('ws://localhost:3000/ws')

    socket.onopen = () => {
      console.log('Conexión WebSocket establecida')
    }

    socket.onmessage = event => {
      const data: Response = JSON.parse(event.data)
      console.log('Datos recibidos:', data)
      // Mostrar el popup después de 5 segundos
      const showPopupTimer = setTimeout(() => {
        setPopup({
          title: 'Nueva Información',
          message: data.movie_response || []
        })

        // Ocultar el popup después de 5 segundos
        const hidePopupTimer = setTimeout(() => {
          setPopup(null)
        }, 5000)

        // Limpieza del hidePopupTimer si el componente se desmonta antes
        return () => clearTimeout(hidePopupTimer)
      }, 5000)

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
    <div className='grid place-items-center w-[100vw] h-[100vh] px-[10vw]'>
      <Header />
      <Outlet />
      {popup && <Popup data={popup} onClose={closePopup} />}
    </div>
  )
}

export default Layout
