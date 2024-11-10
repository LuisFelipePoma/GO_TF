// src/pages/MovieInfo.tsx
import React from 'react'
import { useLocation } from 'react-router-dom'
import type { Movie } from '../types/movies'
import Card from '../components/Card'

const MovieInfo: React.FC = () => {
  const location = useLocation()
  const { movie } = location.state as { movie: Movie }
  console.log(movie)

  return (
    <div>
      PEPEPE
      <h1>{movie.title}</h1>
      <p>{movie.overview}</p>
      <Card movie={movie}/>
    </div>
  )
}

export default MovieInfo
