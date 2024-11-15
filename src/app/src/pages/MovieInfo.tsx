import React, { useEffect } from 'react'
import { useLocation } from 'react-router-dom'
import type { MovieResponse } from '../types/movies'
import { getRecommendations } from '../services/movies'
import { ListMovies } from '../components/ListMovies'
import Skeleton from 'react-loading-skeleton'
import { PLACEHOLDER_URL, URL_IMG } from '../consts/api'
import { TmdbResponse } from '../types/tmdb'
import { useStore } from '../services/store'
import { VoteAvg } from '../components/VoteAvg'
import { motion } from 'framer-motion'

const MovieInfo: React.FC = () => {
  const location = useLocation()
  const { movie } = location.state as { movie: MovieResponse }
  const { movieInfo } = location.state as { movieInfo: TmdbResponse }
  const [recomendations, setRecomendations] = React.useState<MovieResponse[]>(
    []
  )
  const [loading, setLoading] = React.useState<boolean>(true)
  // eslint-disable-next-line @typescript-eslint/no-explicit-any
  const setBackgroundPath = useStore((state: any) => state.setBackgroundPath)

  useEffect(() => {
    let isMounted = true

    setLoading(true)
    getRecommendations(movieInfo.id!, 21).then(res => {
      if (isMounted) {
        setRecomendations(res.movie_response!)
        setLoading(false)
      }
    })

    setBackgroundPath(movieInfo.backdrop_path)
    return () => {
      isMounted = false
      setBackgroundPath(null)
    }
  }, [movieInfo, setBackgroundPath])

  return (
    <motion.div
      initial={{ opacity: 0, x: 100 }}
      animate={{ opacity: 1, x: 0 }}
      exit={{ opacity: 0, x: -100 }}
      transition={{ duration: 0.5 }}
      className='grid grid-cols-1 gap-5 place-content-center w-[100%]'
    >
      <section className='flex gap-5 h-[500px] w-full items-center'>
        <div className='flex-shrink-0 '>
          <img
            src={
              movieInfo.poster_path
                ? URL_IMG(movieInfo.poster_path, 'w300')
                : PLACEHOLDER_URL
            }
            alt={movieInfo.title}
            className='w-[300px] h-[450px] rounded-md self-center filter drop-shadow-md  brightness-95 hover:brightness-105 transition-all duration-1000 ease-in-out'
          />
        </div>
        <div className='flex flex-auto flex-col h-full gap-3 justify-between gap-y-10 '>
          <article className='flex justify-between items-center'>
            <section className='flex flex-col gap-3'>
              <h3 className='underline font-semibold text-balance pr-1'>
                {movieInfo.title} (
                {movieInfo?.release_date
                  ? new Date(movieInfo.release_date).getFullYear()
                  : '20XX'}
                )
              </h3>
              <p>{movieInfo.tagline}</p>
            </section>
            <VoteAvg vote_average={movieInfo.vote_average} />
          </article>
          <article className='flex flex-col gap-5'>
            <p className='text-body-16 text-balance'>{movieInfo.overview}</p>
            <p className='flex flex-wrap gap-3 select-none'>
              {movieInfo.genres!.map((genre, i) => (
                <span
                  className='w-fit px-3 py-1 bg-secondary rounded-md hover:bg-primary transition-colors duration-500 ease-in-out'
                  key={`${i}-genre${genre.id}`}
                >
                  {genre.name}
                </span>
              ))}
            </p>
          </article>
          <article className='flex flex-col gap-2 text-balance'>
            <p>
              <span className='font-bold text-primary'>Release Date: </span>
              {/* Format to 19 November, 2015 */}
              {new Date(movieInfo.release_date!).toLocaleDateString('en-US', {
                year: 'numeric',
                month: 'long',
                day: 'numeric'
              })}
            </p>
            <p>
              <span className='font-bold text-primary'>Director: </span>
              {movie.director}
            </p>
            <section className='flex gap-5 justify-around'>
              <p className='w-[50%] line-clamp-5'>
                <span className='font-bold text-primary'>Actors: </span>
                {movie.actors?.split(',').join(', ')}
              </p>
              <p className='w-[50%] line-clamp-5'>
                <span className='font-bold text-primary'>Characters: </span>
                {movie.characters?.split(',').join(', ')}
              </p>
            </section>
          </article>
        </div>
      </section>
      <section className='flex flex-col gap-5'>
        <h4 className=''>Peliculas Similares</h4>
        <div className='flex flex-wrap gap-3'>
          {loading ? (
            Array.from({ length: 7 }).map((_, i) => (
              <Skeleton
                key={i}
                width={200}
                height={300}
                baseColor='#0B0000'
                highlightColor='#ba0c0c2f'
                borderRadius='10px'
                direction='rtl'
                enableAnimation={true}
                duration={4}
              />
            ))
          ) : (
            <ListMovies movies={recomendations} />
          )}
        </div>
      </section>
    </motion.div>
  )
}

export default MovieInfo
