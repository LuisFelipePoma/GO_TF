import React, { useEffect, useState } from 'react'
import Skeleton from 'react-loading-skeleton'
import 'react-loading-skeleton/dist/skeleton.css'
import { getTmdbInfo } from '../services/movies'
import { useNavigate } from 'react-router-dom'
import { PLACEHOLDER_URL, URL_IMG } from '../consts/api'
import { MovieResponse } from '../types/movies'
import { TmdbResponse } from '../types/tmdb'
import { averageStyles } from '../consts/styles'
import { MovieInfoDummie } from '../dum/movie'

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
    navigate(`/movie/${movie.id}`, { state: { movie, movieInfo } })
    // also take to the top of the page
    window.scrollTo(0, 0)
  }

  return (
    <a
      className={`w-fit max-h-[450px] min-h-[${height}px] h-auto min-w-[${width}px] group rounded-md hover:bg-slate-50/5
       transition-all duration-300 ease-in-out hover:drop-shadow-2xl cursor-pointer`}
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
            alt={movieInfo.title}
          />
        </div>
      )}
      <section className='grid py-2 px-2 gap-2 h-[100px] w-full'>
        <article className='flex w-full justify-between min-h-[25px]'>
          <h5 className='h-full text-body-16 group-hover:underline line-clamp-2'>
            {movieInfo.title} (
            {movieInfo?.release_date
              ? new Date(movieInfo.release_date).getFullYear()
              : '20XX'}
            )
          </h5>
          <p
            className={`hover:brightness-110 transition-all duration-1000 ease-in-out 
                  text-body-14 rounded-full w-fit h-fit font-bold text-black px-2 py-1 ${averageStyles(
                    movieInfo.vote_average!
                  )}`}
          >
            {movieInfo.vote_average?.toPrecision(2)}
          </p>
        </article>
        <div className='flex flex-wrap gap-1 items-start'>
          {movieInfo.genres?.map(genre => (
            <span
              className='w-fit h-fit text-body-12 upper px-[0.3rem] py-[0.1rem] bg-secondary rounded-md hover:bg-tertiary transition-colors duration-500 ease-in-out'
              key={genre.id}
            >
              {genre.name}
            </span>
          ))}
        </div>
      </section>
    </a>
  )
}

export default Card
