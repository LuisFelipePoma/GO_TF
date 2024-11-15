import React, { useEffect, useState } from 'react'
import Skeleton from 'react-loading-skeleton'
import 'react-loading-skeleton/dist/skeleton.css'
import { getTmdbInfo } from '../services/movies'
import { useNavigate } from 'react-router-dom'
import { PLACEHOLDER_URL, URL_IMG } from '../consts/api'
import { MovieResponse } from '../types/movies'
import { TmdbResponse } from '../types/tmdb'
import { MovieInfoDummie } from '../dum/movie'
import { Popularity } from './Popularity'
import { StarRating } from './StartRating'

interface CardProps {
  movie: MovieResponse
}

const Card: React.FC<CardProps> = ({ movie }) => {
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
        // setLoading(false)
      })
    // .finally(() => setLoading(false))
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
      className={`h-[400px] w-[200px] group rounded-md 
        transition-all duration-300 ease-in-out hover:drop-shadow-2xl 
        hover:bg-dark/70 cursor-pointer relative`}
      onClick={() => handleMovieClick(movie)}
    >
      <div
        className={`${
          loading ? 'opacity-100' : 'opacity-0'
        } overflow-hidden w-fit h-[300px] rounded-md relative select-none top-0`}
      >
        <Skeleton
          width={200}
          height={300}
          baseColor='#0B0000'
          highlightColor='#ba0c0c2f'
          borderRadius='10px'
          enableAnimation
          duration={4}
        />
      </div>
      <div
        className={`${
          loading ? 'opacity-0' : 'opacity-100'
        } overflow-hidden absolute top-0 w-fit h-[300px] rounded-md select-none`}
      >
        <main className='relative'>
          <img
            className='w-[200px] h-[300px] object-cover transform group-hover:scale-100 scale-110
            filter grayscale-[55%] group-hover:grayscale-[15%] group-hover:blur-[1px] transition-all duration-1000 ease-in-out'
            src={
              posterPath && posterPath.startsWith('https')
                ? PLACEHOLDER_URL
                : URL_IMG(posterPath)
            }
            alt={movieInfo.title}
            onLoad={() => {
              // timeout
              setTimeout(() => {
                setLoading(false)
              }, 1500)
            }}
          />
          <div
            className='opacity-0 group-hover:opacity-100 absolute top-0 left-0 w-full h-full bg-[#0B0000]/35 bg-opacity-50 
          transition-all duration-1000 ease-in-out p-4'
          >
            <article className='flex flex-wrap gap-2 items-start '>
              {movieInfo.genres?.map(genre => (
                <span
                  className='px-2 py-1 bg-secondary text-body-12 uppercase rounded-md hover:bg-tertiary transition-colors duration-300 ease-in-out cursor-pointer'
                  key={genre.id}
                >
                  {genre.name}
                </span>
              ))}
            </article>
            <p className='mt-5  text-body-14 line-clamp-4 font-semibold'>
              {movieInfo.tagline}
            </p>
          </div>
        </main>
      </div>

      <section className='flex flex-col justify-start gap-3 py-3 px-1 h-[110px] '>
        <div className='flex justify-between'>
          <StarRating voteAverage={movieInfo.vote_average!} />
          <Popularity
            popularity={movieInfo.popularity!}
            className='text-body-12 w-fit h-fit'
          />
        </div>

        <h5 className='h-fit text-body-16 group-hover:underline line-clamp-2'>
          {movieInfo.title} (
          {movieInfo.release_date
            ? new Date(movieInfo.release_date).getFullYear()
            : '20XX'}
          )
        </h5>
      </section>
    </a>
  )
}

export default Card
