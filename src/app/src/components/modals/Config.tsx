/* eslint-disable @typescript-eslint/no-explicit-any */
// Config.tsx
import React, { useState, useEffect } from 'react'
import { motion } from 'framer-motion'
import { useStore } from '../../services/store'
import { IconClose } from '../../assets/IconClose'
import { InputRange } from '../Items/InputRange'

interface SidebarProps {
  isOpen: boolean
  onClose: () => void
}

const Config: React.FC<SidebarProps> = ({ isOpen, onClose }) => {
  const nMoviesSearch = useStore((state: any) => state.nMoviesSearch)
  const nMoviesHome = useStore((state: any) => state.nMoviesHome)
  const setNMoviesHome = useStore((state: any) => state.setNMoviesHome)
  const nMoviesUser = useStore((state: any) => state.nMoviesUser)
  const setNMoviesSearch = useStore((state: any) => state.setNMoviesSearch)
  const setNMoviesUser = useStore((state: any) => state.setNMoviesUser)

  const [formSearch, setFormSearch] = useState(nMoviesSearch)
  const [formHome, setFormHome] = useState(nMoviesHome)
  const [formUser, setFormUser] = useState(nMoviesUser)

  useEffect(() => {
    setFormSearch(nMoviesSearch)
    setFormHome(nMoviesHome)
    setFormUser(nMoviesUser)
  }, [nMoviesSearch, nMoviesHome, nMoviesUser])

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault()
    setNMoviesSearch(formSearch)
    setNMoviesHome(formHome)
    setNMoviesUser(formUser)
  }

  return (
    <>
      {/* Overlay */}
      {isOpen && (
        <div
          className='fixed inset-0 bg-dark opacity-50 z-10'
          onClick={onClose}
        />
      )}
      {/* Sidebar */}
      <motion.div
        className='fixed top-0 left-0 h-full w-80 bg-dark/85 backdrop-blur-md z-20 shadow-2xl drop-shadow-xl'
        initial={{ x: '-100%' }}
        animate={{ x: isOpen ? '0%' : '-100%' }}
        transition={{ type: 'spring', stiffness: 300, damping: 30 }}
      >
        {/* Contenido del sidebar */}
        <div className='flex items-center gap-5 py-2'>
          <button
            onClick={onClose}
            className='p-4 transform transition-transform duration-300 hover:rotate-180 hover:scale-125'
          >
            <IconClose />
          </button>
          <h2 className='text-3xl text-white underline'>Configuración</h2>
        </div>
        {/* Panel de configuracion de los Numeros */}
        <form
          className='flex flex-col gap-4 p-4 h-full'
          onSubmit={handleSubmit}
        >
          <div className='flex flex-col gap-2'>
            <label htmlFor='nMoviesSearch' className='text-white'>
              Número de películas en la búsqueda
            </label>
            <InputRange formSearch={formSearch} setFormSearch={setFormSearch} />
          </div>
          <div className='flex flex-col gap-2'>
            <label htmlFor='nMoviesHome' className='text-white'>
              Número de películas en la página principal
            </label>
            <InputRange formSearch={formHome} setFormSearch={setFormHome} />
          </div>
          <div className='flex flex-col gap-2'>
            <label htmlFor='nMoviesUser' className='text-white'>
              Número de películas en la página de usuario
            </label>
            <InputRange formSearch={formUser} setFormSearch={setFormUser} />
          </div>
          <button
            type='submit'
            className='bg-primary text-white px-4 py-3 rounded mt-5 disabled:bg-tertiary
              disabled:cursor-not-allowed 
              disabled:text-gray
              transition-all duration-300 ease-in-out hover:bg-secondary
            '
            disabled={
              formSearch === nMoviesSearch &&
              formHome === nMoviesHome &&
              formUser === nMoviesUser
            }
          >
            Apply changes
          </button>
        </form>
      </motion.div>
    </>
  )
}

export default Config
