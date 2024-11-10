import React from 'react'
import { Movie } from '../types/movies'
import Card from './Card'

interface Props {
  movies: Movie[]
}

const averageStyles = (average: number) => {
  if (average < 5) {
    return 'bg-red-400'
  } else if (average < 7) {
    return 'bg-yellow-400'
  } else {
    return 'bg-green-400'
  }
}

export const ListMovies: React.FC<Props> = ({ movies }) => {
  return (
    <div className='flex flex-wrap gap-x-5 gap-y-[25px] w-full justify-center'>
      {movies.map(movie => {
        return (
          <Card key={movie.id} movie={movie}>
            <div className='flex justify-center py-2'>
              <h5 className='w-[150px] text-body-20 h-[75px]'>{movie.title}</h5>
              <p
                className={`rounded-[50%] w-fit h-fit font-bold text-black px-2 py-1 
                ${averageStyles(movie.vote_average!)}`}
              >
                {movie.vote_average?.toPrecision(2)}
              </p>
            </div>
          </Card>
        )
      })}
    </div>
  )
}
