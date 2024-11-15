// src/pages/Home.tsx
import React, { useEffect, useState } from 'react'
import { getAllMovies, getMoviesByQuery } from '../services/movies'
import { ListMovies } from '../components/ListMovies'
import { MovieResponse } from '../types/movies'
import { useDebounce } from 'use-debounce' // Import useDebounce
import { Loader } from '../components/Loader'
import { motion } from 'framer-motion'

const N_MOVIES = 21
const DEBOUNCE_DELAY = 750 // milliseconds

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
    <motion.div
      initial={{ opacity: 0 }}
      animate={{ opacity: 1 }}
      exit={{ opacity: 0 }}
      transition={{ duration: 1.65}}
      className='flex flex-col gap-10 h-[85vh] w-full'
    >
      <input
        type='text'
        placeholder='Search by movie title, genre, keywords or actors.'
        className='w-[35%] p-2 text-body-16 rounded-sm 
         outline-none self-end bg-transparent border-b-2 border-secondary text-light
         placeholder:text-gray transition-all duration-300 ease-in-out 
         hover:bg-gray-800 hover:border-primary focus:bg-gray-700 focus:border-primary 
         focus:shadow-lg transform hover:scale-105 focus:scale-105 backdrop-blur-sm'
        value={query}
        onChange={handleQuery}
      />
      {loading ? <Loader /> : <ListMovies movies={movies} />}
    </motion.div>
  )
}

export default Home
