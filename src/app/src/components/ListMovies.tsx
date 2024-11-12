import React from 'react'
import { MovieResponse } from '../types/movies'
import Card from './Card'
import { averageStyles } from '../consts/styles'

interface Props {
  movies: MovieResponse[]
}

export const ListMovies: React.FC<Props> = ({ movies }) => {
  return (
    <div
      className='grid grid-cols-[repeat(auto-fill,minmax(200px,1fr))] 
    gap-y-10 gap-x-4 w-full h-[85%]'
    >
      {movies.map(movie => (
        <Card key={movie.id} movie={movie}>
          <section className='grid py-2 px-2 gap-2 h-[100px] w-full'>
            <article className='flex w-full justify-between min-h-[25px]'>
              <h5 className='h-full text-body-16 group-hover:underline line-clamp-2'>
                {movie.title}
              </h5>
              <p
                className={`hover:brightness-110 transition-all duration-1000 ease-in-out 
                  text-body-14 rounded-full w-fit h-fit font-bold text-black px-2 py-1 ${averageStyles(
                    movie.vote_average!
                  )}`}
              >
                {movie.vote_average?.toPrecision(2)}
              </p>
            </article>
            <div className='flex flex-wrap gap-1 items-start'>
              {movie.genres?.split(',').map(genre => (
                <span
                  className='w-fit h-fit text-body-12 upper px-[0.3rem] py-[0.1rem] bg-secondary rounded-md hover:bg-tertiary transition-colors duration-500 ease-in-out'
                  key={genre}
                >
                  {genre}
                </span>
              ))}
            </div>
          </section>
        </Card>
      ))}
    </div>
  )
}
