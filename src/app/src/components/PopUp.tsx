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
  return (
    <div className='fixed top-2 right-2 bg-black bg-opacity-50 w-[220px]'>
      <article className='bg-transparent backdrop-blur-lg p-4 rounded relative w-full'>
        <p className='text-body-20 font-bold text-left h-[75px] w-full'>
          Peliculas para el Usuario {data.id}
        </p>
        <section className='max-h-[1050px] overflow-y-hidden hover:overflow-y-auto no-scrollbar flex flex-col gap-5 items-center '>
          {data.message.map(movie => (
            <Card
              key={'PopUp-' + movie.id}
              movie={{ id: movie.id! } as MovieResponse}
              width={100}
              height={350}
            />
          ))}
        </section>
        <button
          className='fixed top-4 right-1 w-10 h-10 bg-primary text-light rounded-full flex justify-center items-center shadow-lg transition transform hover:bg-secondary hover:rotate-90 hover:scale-110 duration-300 ease-out'
          onClick={onClose}
          aria-label='Close popup'
        >
          <span className='text-xl font-bold'>×</span>
        </button>
      </article>
    </div>
  )
}
