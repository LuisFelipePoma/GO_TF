import React, { useEffect, useState } from 'react'
import Skeleton from 'react-loading-skeleton'
import 'react-loading-skeleton/dist/skeleton.css'
import { getTmdbInfo } from '../services/movies'
import { useNavigate } from 'react-router-dom'
import { PLACEHOLDER_URL, URL_IMG } from '../consts/api'
import { MovieResponse } from '../types/movies'
import { TmdbResponse } from '../types/tmdb'
import { MovieInfoDummie } from '../dum/movie'
import { VoteAvg } from './VoteAvg'

interface CardProps {
  movie: MovieResponse
  width?: number
  height?: number
}

const Card: React.FC<CardProps> = ({ movie, width = 200, height = 400 }) => {
  const [posterPath, setPosterPath] = useState(
    movie.poster_path ?? PLACEHOLDER_URL
  )
  const [movieInfo, setMovieInfo] = useState<TmdbResponse>(MovieInfoDummie)
  const [loading, setLoading] = useState(true)
  const navigate = useNavigate()

  useEffect(() => {
    getTmdbInfo(movie.id!)
      .then(res => {
        setMovieInfo(res)
        if (
          res.poster_path === null ||
          res.poster_path === '' ||
          res.poster_path === undefined
        ) {
          throw new Error('No poster path found')
        }
        setPosterPath(res.poster_path!)
      })
      .catch(() => {
        setPosterPath(PLACEHOLDER_URL)
      })
      .finally(() => setLoading(false))
  }, [movie.id, movie.poster_path])

  const handleMovieClick = (movie: MovieResponse) => {
    if (!movieInfo) return
    console.log(movieInfo)
    navigate(`/movie/${movie.id}`, { state: { movie, movieInfo } })
    // also take to the top of the page
    window.scrollTo(0, 0)
  }

  return (
    <a
      className={`h-[${height}px] w-[${width}px] group rounded-md 
        transition-all duration-300 ease-in-out hover:drop-shadow-2xl cursor-pointer relative`}
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
        <div className='overflow-hidden w-full h-[300px] rounded-md relative select-none'>
          <img
            className='w-full h-full object-cover transition-all duration-1000 ease-in-out transform group-hover:scale-100 scale-110
            filter grayscale-[55%] group-hover:grayscale-[15%]'
            src={
              posterPath && posterPath.startsWith('https')
                ? PLACEHOLDER_URL
                : URL_IMG(posterPath)
            }
            alt={movieInfo.title}
          />
          <div
            className='opacity-0 group-hover:opacity-100 absolute top-0 left-0 w-full h-full bg-[#0B0000]/20 bg-opacity-50 
          transition-all duration-1000 ease-in-out'
          >
            <article className='flex flex-wrap gap-2 items-start p-4'>
              {movieInfo.genres?.map(genre => (
                <span
                  className='px-2 py-1 bg-secondary text-body-12 uppercase rounded-md hover:bg-tertiary transition-colors duration-300 ease-in-out cursor-pointer'
                  key={genre.id}
                >
                  {genre.name}
                </span>
              ))}
            </article>
          </div>
        </div>
      )}
      <section className='flex flex-col justify-between py-2 h-[110px] '>
        <h5 className='h-fit text-body-16 group-hover:underline line-clamp-2'>
          {movieInfo.title} (
          {movieInfo.release_date
            ? new Date(movieInfo.release_date).getFullYear()
            : '20XX'}
          )
        </h5>
        <VoteAvg
          vote_average={movieInfo.vote_average}
          className='text-body-12 w-fit h-fit transition-transform duration-300 ease-in-out hover:scale-105 p-4'
        />
      </section>
    </a>
  )
}

export default Card
