// src/components/PopUp.tsx
import React from 'react'
import { motion } from 'framer-motion'
import { MovieResponse } from '../../types/movies'
import Card from '../Card'
import { IconClose } from '../../assets/IconClose'

export interface PopupData {
  message: MovieResponse[]
  id: string
}

interface PopupProps {
  data: PopupData
  isVisible: boolean
  setIsVisible: (isVisible: boolean) => void
}

export const Popup: React.FC<PopupProps> = ({
  data,
  isVisible,
  setIsVisible
}) => {
  function handleClose () {
    setIsVisible(false)
  }

  return (
    <motion.div
      initial={{ x: '100%', opacity: 0 }}
      animate={isVisible ? { x: 0, opacity: 1 } : { x: '100%', opacity: 0 }}
      transition={{ duration: 0.5, ease: 'easeIn' }}
      className={`fixed top-0 right-0 bg-dark/90 w-[300px] h-full z-20 filter backdrop-blur-md shadow-2xl `}
    >
      <header className='flex justify-between p-4'>
        <p className='text-body-16 text-left text-balance '>
          Movies that
          <span className='font-bold text-primary text-body-20'>
            {' '}
            User {data.id}{' '}
          </span>
          might like
        </p>
        <button
          onClick={handleClose}
          className='w-10 h-10 bg-primary text-light rounded-full flex justify-center items-center shadow-lg transition 
          transform hover:bg-secondary hover:rotate-90 hover:scale-110 duration-300 ease-out'
        >
          <IconClose />
        </button>
      </header>
      <section className='flex-1 h-full py-5 overflow-y-scroll no-scrollbar flex flex-col gap-5 items-center w-full justify-start'>
        {data.message.map(movie => (
          <Card
            key={'PopUp-' + movie.id}
            movie={{ id: movie.id! } as MovieResponse}
          />
        ))}
      </section>
    </motion.div>
  )
}
