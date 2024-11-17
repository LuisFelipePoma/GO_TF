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
  const nMoviesRecomendations = useStore(
    (state: any) => state.nMoviesRecomendations
  )
  const setNMoviesRecomendations = useStore(
    (state: any) => state.setNMoviesRecomendations
  )

  const [formSearch, setFormSearch] = useState(nMoviesSearch)
  const [formHome, setFormHome] = useState(nMoviesHome)
  const [formUser, setFormUser] = useState(nMoviesUser)
  const [formRecomendations, setFormRecomendations] = useState(
    nMoviesRecomendations
  )

  useEffect(() => {
    setFormSearch(nMoviesSearch)
    setFormHome(nMoviesHome)
    setFormUser(nMoviesUser)
    setFormRecomendations(nMoviesRecomendations)
  }, [nMoviesSearch, nMoviesHome, nMoviesUser, nMoviesRecomendations])

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault()
    setNMoviesSearch(formSearch)
    setNMoviesHome(formHome)
    setNMoviesUser(formUser)
    setNMoviesRecomendations(formRecomendations)
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
          <InputRange
            formSearch={formHome}
            setFormSearch={setFormHome}
            description='Number of movies on the home page'
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
              <path d='M19.072 21h-14.144a1.928 1.928 0 0 1 -1.928 -1.928v-6.857c0 -.512 .203 -1 .566 -1.365l7.07 -7.063a1.928 1.928 0 0 1 2.727 0l7.071 7.063c.363 .362 .566 .853 .566 1.365v6.857a1.928 1.928 0 0 1 -1.928 1.928z'></path>{' '}
              <path d='M7 13v4h10v-4l-5 -5'></path>{' '}
              <path d='M14.8 5.2l-11.8 11.8'></path> <path d='M7 17v4'></path>{' '}
              <path d='M17 17v4'></path>{' '}
            </svg>
          </InputRange>

          <InputRange
            formSearch={formUser}
            setFormSearch={setFormUser}
            description='Number of recomendations per user'
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
              {' '}
              <path d='M8 7a4 4 0 1 0 8 0a4 4 0 0 0 -8 0'></path>{' '}
              <path d='M6 21v-2a4 4 0 0 1 4 -4h.5'></path>{' '}
              <path d='M18 22l3.35 -3.284a2.143 2.143 0 0 0 .005 -3.071a2.242 2.242 0 0 0 -3.129 -.006l-.224 .22l-.223 -.22a2.242 2.242 0 0 0 -3.128 -.006a2.143 2.143 0 0 0 -.006 3.071l3.355 3.296z'></path>{' '}
            </svg>{' '}
          </InputRange>
          <InputRange
            formSearch={formSearch}
            setFormSearch={setFormSearch}
            description='Number of movies on the search'
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
              {' '}
              <path d='M10 10m-7 0a7 7 0 1 0 14 0a7 7 0 1 0 -14 0'></path>{' '}
              <path d='M21 21l-6 -6'></path>{' '}
            </svg>
          </InputRange>
          <InputRange
            formSearch={formRecomendations}
            setFormSearch={setFormRecomendations}
            description='Number of recomendations per movie'
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
              {' '}
              <path d='M3 7m0 2a2 2 0 0 1 2 -2h14a2 2 0 0 1 2 2v9a2 2 0 0 1 -2 2h-14a2 2 0 0 1 -2 -2z'></path>{' '}
              <path d='M16 3l-4 4l-4 -4'></path>{' '}
            </svg>
          </InputRange>
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
              formUser === nMoviesUser &&
              formRecomendations === nMoviesRecomendations
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
