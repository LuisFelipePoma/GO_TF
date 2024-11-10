import React, { useEffect, useState } from 'react'
import Skeleton from 'react-loading-skeleton'
import 'react-loading-skeleton/dist/skeleton.css'
import { Movie } from '../types/movies'
import { getTmdbInfo } from '../services/movies'
import { useNavigate } from 'react-router-dom'
import { URL_IMG } from '../consts/api'

interface CardProps {
  movie: Movie
  children?: React.ReactNode
}

const PLACEHOLDER_URL = 'https://via.placeholder.com/200x300'

const Card: React.FC<CardProps> = ({ movie, children }) => {
  const [posterPath, setPosterPath] = useState(
    'https://via.placeholder.com/200x300'
  )
  const [loading, setLoading] = useState(true)
  const navigate = useNavigate()

  useEffect(() => {
    const loadPoster = async () => {
      if (!movie.poster_path || !movie.id) {
        setPosterPath(PLACEHOLDER_URL)
        setLoading(false)
        return
      }

      try {
        const response = await fetch(URL_IMG(movie.poster_path))
        if (!response.ok) {
          throw new Error('Network response was not ok')
        }
        setPosterPath(URL_IMG(movie.poster_path))
      } catch (error) {
        console.error('Error fetching poster:', error)
        try {
          const res = await getTmdbInfo(movie.id)
          if (res.poster_path === null) {
            throw new Error('No poster path found')
          }
          setPosterPath(URL_IMG(res.poster_path!))
        } catch (error) {
          console.error('Error fetching TMDB info:', error)
          setPosterPath(PLACEHOLDER_URL)
        }
      } finally {
        setLoading(false)
      }
    }

    loadPoster()
  }, [movie.id, movie.poster_path])

  const handleMovieClick = (movie: Movie) => {
    navigate(`/movie/${movie.id}`, { state: { movie } })
  }

  return (
    <a
      className='w-fit h-min-[375px] group rounded-md  transition-all duration-300 ease-in-out hover:drop-shadow-2xl hover:bg-dark cursor-pointer'
      onClick={() => handleMovieClick(movie)}
    >
      {loading ? (
        <Skeleton
          width={200}
          height={300}
          className='bg-primary'
          baseColor='#202020'
          highlightColor='#444'
          borderRadius='10px'
        />
      ) : (
        <div className='overflow-hidden w-[200px] h-[300px] rounded-md'>
          <img
            className='w-full h-full object-cover transition-transform duration-500 ease-in-out transform group-hover:scale-100 scale-110'
            src={posterPath}
            alt={movie.title}
          />
        </div>
      )}
      {children}
    </a>
  )
}

export default Card
