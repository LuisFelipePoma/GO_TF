import React, { useEffect } from 'react'
import { useLocation } from 'react-router-dom'
import type { Movie } from '../types/movies'
import { getRecommendations } from '../services/movies'
import { ListMovies } from '../components/ListMovies'
import Card from '../components/Card'

const MovieInfo: React.FC = () => {
  const location = useLocation()
  const { movie } = location.state as { movie: Movie }
  const [recomendations, setRecomendations] = React.useState<Movie[]>([])

  useEffect(() => {
    let isMounted = true

    getRecommendations(movie.title!).then(res => {
      if (isMounted) {
        console.log(res)
        setRecomendations(res.movie_response!)
      }
    })

    return () => {
      isMounted = false
    }
  }, [movie.title])

  return (
    <div>
      <section>
        <h1>{movie.title}</h1>
        <p>{movie.overview}</p>
        <Card movie={movie} />
      </section>
      <ListMovies movies={recomendations} />
    </div>
  )
}

export default MovieInfo
