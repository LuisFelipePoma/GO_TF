import React from 'react'
import { Movie } from '../types/movies'
import Card from './Card'

interface Props {
  movies: Movie[]
}

export const ListMovies: React.FC<Props> = ({ movies }) => {
  return (
    <div className='flex flex-wrap gap-x-5 gap-y-2 w-full justify-center'>
      {movies.map(movie => {
        return (
          <Card key={movie.id} movie={movie}>
            <p className='w-[150px] text-body-14'>{movie.title}</p>
          </Card>
        )
      })}
    </div>
  )
}
