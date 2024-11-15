// src/components/PopUp.tsx
import { useState, useEffect } from 'react'
import { MovieResponse } from '../types/movies'
import Card from './Card'

export interface PopupData {
  message: MovieResponse[]
  id: string
}

export const Popup: React.FC<{ data: PopupData; onClose: () => void }> = ({
  data,
  onClose
}) => {
  const [isVisible, setIsVisible] = useState(false)

  useEffect(() => {
    setIsVisible(true)
  }, [data.id])

  const handleClose = () => {
    setIsVisible(false)
    setTimeout(onClose, 300) // Duration matches the CSS transition
  }

  return (
    <div
      className={`transition-all duration-300 ease-in fixed top-0 right-0 bg-[#0B0000]/60 w-fit${
        isVisible ? 'translate-x-0 opacity-100' : 'translate-x-full opacity-0'
      } shadow-2xl`}
    >
      <article className='h-full backdrop-blur-md p-3 rounded relative w-[275px]'>
        <p className='text-body-20 font-bold text-left h-[75px] w-[200px]'>
          Peliculas para el Usuario {data.id}
        </p>
        <section
          className='pt-[300px] max-h-[1000px] overflow-y-hidden hover:overflow-y-auto no-scrollbar 
        flex flex-col gap-5 items-center w-full justify-center'
        >
          {data.message.map(movie => (
            <Card
              key={'PopUp-' + movie.id}
              movie={{ id: movie.id! } as MovieResponse}
            />
          ))}
        </section>
        <button
          className='fixed top-4 right-2 w-10 h-10 bg-primary text-light rounded-full flex justify-center items-center shadow-lg transition transform hover:bg-secondary hover:rotate-90 hover:scale-110 duration-300 ease-out'
          onClick={handleClose}
          aria-label='Close popup'
        >
          <span className='text-xl font-bold'>×</span>
        </button>
      </article>
    </div>
  )
}
