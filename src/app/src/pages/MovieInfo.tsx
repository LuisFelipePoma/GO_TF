import React, { useEffect } from 'react'
import { useLocation } from 'react-router-dom'
import type { MovieResponse } from '../types/movies'
import { getRecommendations } from '../services/movies'
import { ListMovies } from '../components/ListMovies'
import Skeleton from 'react-loading-skeleton'
import { URL_IMG } from '../consts/api'
import { averageStyles } from '../consts/styles'

const MovieInfo: React.FC = () => {
  const location = useLocation()
  const { movie } = location.state as { movie: MovieResponse }
  const { posterPath } = location.state as { posterPath: string }
  const [recomendations, setRecomendations] = React.useState<MovieResponse[]>(
    []
  )
  const [loading, setLoading] = React.useState<boolean>(true)

  console.log(posterPath)

  useEffect(() => {
    let isMounted = true
    setLoading(true)

    getRecommendations(movie.title!, 20).then(res => {
      if (isMounted) {
        setRecomendations(res.movie_response!)
        setLoading(false)
      }
    })

    return () => {
      isMounted = false
    }
  }, [movie])

  return (
    <div className='grid grid-cols-1 gap-5 place-content-center w-full'>
      <section className='flex gap-x-10 overflow-y-hidden'>
        <div className='relative'>
          <img
            src={URL_IMG(posterPath, 300)}
            alt={movie.poster_path}
            className='min-w-[300px] max-w-[300px] h-auto rounded-md filter drop-shadow-md brightness-90 hover:brightness-100
             transition-all duration-1000 ease-in-out'
          />
          <span
            className={`drop-shadow-md select-none absolute top-2 right-2 rounded-full w-fit h-fit font-bold text-black px-2 py-1 ${averageStyles(
              movie.vote_average!
            )}`}
          >
            {movie.vote_average?.toPrecision(2)}
          </span>
        </div>

        <article className='flex flex-col gap-3 justify-start gap-y-10'>
          <h3 className='underline leading-none font-semibold'>
            {movie.title}
          </h3>
          <p className='text-body-20'>{movie.overview}</p>
          <div className='flex gap-x-4 select-none'>
            {movie.genres?.split(',').map((genre, i) => (
              <span
                className='px-2 py-1 bg-secondary rounded-md hover:bg-tertiary transition-colors duration-500 ease-in-out'
                key={i}
              >
                {genre}
              </span>
            ))}
          </div>
          <div className='flex flex-col gap-2 text-balance'>
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
          </div>
        </article>
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
                className='bg-primary'
                baseColor='#202020'
                highlightColor='#444'
                borderRadius='10px'
              />
            ))
          ) : (
            <ListMovies movies={recomendations} />
          )}
        </div>
      </section>
    </div>
  )
}

export default MovieInfo
