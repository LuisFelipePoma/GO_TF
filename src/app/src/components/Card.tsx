import React, { useState } from 'react'
import Skeleton from 'react-loading-skeleton'
import 'react-loading-skeleton/dist/skeleton.css'
// import { URL_IMG } from '../consts/api'
import { Movie } from '../types/movies'
// import { TmdbResponse } from '../types/tmdb'
// import { getTmdbInfo } from '../services/movies'
import { useNavigate } from 'react-router-dom'

interface CardProps {
  movie: Movie
  children?: React.ReactNode
}

const Card: React.FC<CardProps> = ({ movie, children }) => {
  // const [movieInfo, setMovieInfo] = useState<TmdbResponse>()
  const [loading] = useState(true)
  const navigate = useNavigate()

  // useEffect(() => {
  //   // Simula la carga de la información de la película
  //   getTmdbInfo(movie.id!).then(data => {
  //     setMovieInfo(data)
  //     setLoading(false)
  //   })
  // }, [movie.id])

  function handleMovieClick (movie: Movie) {
    navigate(`/movie/${movie.id}`, { state: { movie } })
  }

  return (
    <a className='w-fit' onClick={() => handleMovieClick(movie)}>
      {loading ? (
        <Skeleton width={200} height={300} />
      ) : (
        <img
          className='w-[200px] h-auto'
          // src={URL_IMG(movie.poster_path)}
          alt={movie.title}
        />
      )}
      {children}
    </a>
  )
}

export default Card
