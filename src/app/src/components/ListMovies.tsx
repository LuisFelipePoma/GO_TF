import React from 'react'
import { MovieResponse } from '../types/movies'
import Card from './Card'

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
        <Card key={movie.id} movie={movie} />
      ))}
    </div>
  )
}
