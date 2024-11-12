// src/pages/Home.tsx
import React, { useEffect, useState } from 'react'
import { getAllMovies, getMoviesByQuery } from '../services/movies'
import { ListMovies } from '../components/ListMovies'
import { Loading } from '../components/Loading'
import { MovieResponse } from '../types/movies'
import { useDebounce } from 'use-debounce' // Import useDebounce

const N_MOVIES = 21
const DEBOUNCE_DELAY = 500 // milliseconds

const Home: React.FC = () => {
  const [movies, setMovies] = useState<MovieResponse[]>([])
  const [loading, setLoading] = useState<boolean>(false)
  const [query, setQuery] = useState<string>('')

  // Debounce the query input
  const [debouncedQuery] = useDebounce(query, DEBOUNCE_DELAY)

  useEffect(() => {
    setLoading(true)
    getAllMovies(N_MOVIES).then(data => {
      setMovies(data.movie_response!)
      setLoading(false)
    })
  }, [])

  useEffect(() => {
    if (debouncedQuery === '') {
      return
    }
    if (debouncedQuery) {
      setLoading(true)
      getMoviesByQuery(debouncedQuery, N_MOVIES).then(data => {
        if (!data.movie_response) {
          setLoading(false)
          return
        }
        setMovies(data.movie_response!)
        setLoading(false)
      })
    } else {
      setMovies([])
      setLoading(false)
    }
  }, [debouncedQuery])

  function handleQuery (e: React.ChangeEvent<HTMLInputElement>) {
    setQuery(e.target.value)
  }

  return (
    <section className='flex flex-col gap-10 h-[85vh] w-full'>
      <input
        type='text'
        placeholder='Search by movie title, genre, keywords or actors.'
        className='w-[35%] p-2 text-body-20 rounded-sm drop-shadow-sm
        outline-none self-end bg-dark border-b border-secondary text-primary placeholder:text-gray '
        value={query}
        onChange={handleQuery}
      />
      {loading ? <Loading /> : <ListMovies movies={movies} />}
    </section>
  )
}

export default Home
