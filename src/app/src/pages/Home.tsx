// src/pages/Home.tsx
import React, { useEffect } from 'react'
import { Movie } from '../types/movies'
import { getAllMovies } from '../services/movies'
import { ListMovies } from '../components/ListMovies'
import { Loading } from '../components/Loading'

const Home: React.FC = () => {
  const [movies, setMovies] = React.useState<Movie[]>([])
  const [loading, setLoading] = React.useState<boolean>()

  useEffect(() => {
    setLoading(true)
    getAllMovies(100).then(data => {
      setMovies(data)
      setLoading(false)
    })
  }, [])

  return <>{loading ? <Loading /> : <ListMovies movies={movies} />}</>
}

export default Home
