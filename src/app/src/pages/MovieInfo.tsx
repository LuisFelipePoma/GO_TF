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
    <div className='grid grid-cols-1 gap-5 place-content-center w-full'>
      <section className='flex gap-x-10 overflow-y-hidden relative items-center'>
        <img
          src={
            movieInfo.poster_path ?? movieInfo.poster_path
              ? URL_IMG(movieInfo.poster_path, 'w300')
              : PLACEHOLDER_URL
          }
          alt={movieInfo.title}
          className='min-w-[300px] max-w-[300px] h-auto rounded-md filter drop-shadow-md brightness-90 hover:brightness-100
             transition-all duration-1000 ease-in-out select-none pointer-events-none'
        />

        <article className='flex flex-col h-full gap-3 justify-between gap-y-10 '>
          <div className='flex justify-between'>
            <h3 className='underline leading-none font-semibold '>
              {movieInfo.title} (
              {movieInfo?.release_date
                ? new Date(movieInfo.release_date).getFullYear()
                : '20XX'}
              )
            </h3>
            <VoteAvg
              vote_average={movieInfo.vote_average}
              className='text-body-20 h-fit py-2'
            />
          </div>
          <div className='flex flex-col gap-5'>
            <p className='text-body-20'>{movieInfo.overview}</p>
            <p className='flex flex-wrap gap-3'>
              {movieInfo.genres!.map((genre, i) => (
                <span
                  className='w-fit px-2 py-1 bg-secondary rounded-md hover:bg-tertiary transition-colors duration-500 ease-in-out'
                  key={`${i}-genre${genre.id}`}
                >
                  {genre.name}
                </span>
              ))}
            </p>
          </div>
          <div className='flex flex-col gap-2 text-balance'>
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
