import React, { useEffect, useState } from 'react'
import Skeleton from 'react-loading-skeleton'
import 'react-loading-skeleton/dist/skeleton.css'
import { getTmdbInfo } from '../services/movies'
import { useNavigate } from 'react-router-dom'
import { URL_IMG } from '../consts/api'
import { MovieResponse } from '../types/movies'

interface CardProps {
  movie: MovieResponse
  children?: React.ReactNode
}

const PLACEHOLDER_URL = 'https://via.placeholder.com/200x300'

const Card: React.FC<CardProps> = ({ movie, children }) => {
  const [posterPath, setPosterPath] = useState(
    movie.poster_path ?? PLACEHOLDER_URL
  )
  const [loading, setLoading] = useState(true)
  const navigate = useNavigate()

  useEffect(() => {
    const controller = new AbortController()
    const { signal } = controller

    const loadPoster = async () => {
      if (!movie.poster_path || !movie.id) {
        setLoading(false)
        return
      }

      try {
        const response = await fetch(URL_IMG(movie.poster_path), { signal })
        if ((await response.text()).includes('File Not Found')) {
          throw new Error('Network response was not ok')
        }
        setLoading(false)
      } catch (error) {
        if (error instanceof Error && error.name === 'AbortError') {
          console.log('returning')
          return
        }
        console.error('Error fetching poster:', error)
        try {
          const res = await getTmdbInfo(movie.id)
          if (res.poster_path === null) {
            console.log('null')
            throw new Error('No poster path found')
          }
          setPosterPath(res.poster_path ?? PLACEHOLDER_URL)
          setLoading(false)
        } catch (error) {
          console.error('Error fetching TMDB info:', error)
          setPosterPath(PLACEHOLDER_URL)
          setLoading(false)
        }
      }
    }
    loadPoster()

    return () => {
      console.log('Aborting fetch')
      controller.abort()
    }
  }, [movie.id, movie.poster_path])

  const handleMovieClick = (movie: MovieResponse) => {
    navigate(`/movie/${movie.id}`, { state: { movie, posterPath } })
    // also take to the top of the page
    window.scrollTo(0, 0)
  }

  return (
    <a
      className='w-fit h-fit max-h-[450px] min-h-[400px] min-w-[200px] group rounded-md hover:bg-slate-50/5
       transition-all duration-300 ease-in-out hover:drop-shadow-2xl cursor-pointer'
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
          duration={2}
        />
      ) : (
        <div className='overflow-hidden w-full h-[300px] rounded-md'>
          <img
            className='w-full h-full object-cover transition-all duration-1000 ease-in-out transform group-hover:scale-100 scale-110
            filter grayscale-[60%] group-hover:grayscale-[25%]
            '
            src={
              posterPath && posterPath.startsWith('https')
                ? PLACEHOLDER_URL
                : URL_IMG(posterPath)
            }
            alt={movie.title}
          />
        </div>
      )}
      {children}
    </a>
  )
}

export default Card
