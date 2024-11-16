// src/pages/Home.tsx
import React, { useEffect, useState } from 'react'
import { getAllMovies, getMoviesByQuery } from '../services/movies'
import { ListMovies } from '../components/ListMovies'
import { MovieResponse } from '../types/movies'
import { useDebounce } from 'use-debounce' // Import useDebounce
import { Loader } from '../components/Items/Loader'
import { motion } from 'framer-motion'
import { useStore } from '../services/store'
import { genres } from '../consts/tag'
import { GenreTag } from '../components/Items/GenreTag'
import { DEBOUNCE_DELAY } from '../consts/debounce'

const Home: React.FC = () => {
  const [movies, setMovies] = useState<MovieResponse[]>([])
  const [loading, setLoading] = useState<boolean>(false)
  const [query, setQuery] = useState<string>('')
  const [selectedGenres, setSelectedGenres] = useState<string[]>(['All'])
  // eslint-disable-next-line @typescript-eslint/no-explicit-any
  const nMovieSearch = useStore((state: any) => state.nMoviesSearch)
  // eslint-disable-next-line @typescript-eslint/no-explicit-any
  const nMoviesHome = useStore((state: any) => state.nMoviesHome)

  // Debounce the query input
  const [debouncedQuery] = useDebounce(query, DEBOUNCE_DELAY)

  useEffect(() => {
    setLoading(true)
    getAllMovies(selectedGenres, nMoviesHome).then(data => {
      setMovies(data.movie_response!)
      //timeout
      setTimeout(() => {
        setLoading(false)
      }, 1500)
    })
  }, [nMoviesHome, selectedGenres])

  useEffect(() => {
    if (debouncedQuery === '') {
      return
    }
    if (debouncedQuery) {
      setLoading(true)
      getMoviesByQuery(debouncedQuery, nMovieSearch).then(data => {
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
  }, [debouncedQuery, nMovieSearch])

  function handleQuery (e: React.ChangeEvent<HTMLInputElement>) {
    setQuery(e.target.value)
  }

  function handleGenreSelected (genre: string) {
    // if select All
    if (genre === 'All') {
      setSelectedGenres(['All'])
      return
    }

    if (selectedGenres.includes(genre)) {
      setSelectedGenres(selectedGenres.filter(g => g !== genre))
    } else {
      // avoid 'All'
      if (selectedGenres.includes('All')) {
        setSelectedGenres([genre])
      } else {
        setSelectedGenres([...selectedGenres, genre])
      }
    }
  }

  return (
    <motion.div
      initial={{ opacity: 0 }}
      animate={{ opacity: 1 }}
      exit={{ opacity: 0 }}
      transition={{ duration: 0.75 }}
      className='flex flex-col gap-10 h-[85vh] w-full'
    >
      <div className='flex justify-between items-center'>
        <ul className='flex flex-wrap gap-3 select-none w-[750px]'>
          {genres.map(genre => (
            <GenreTag
              key={genre.id}
              genre={genre}
              onClick={handleGenreSelected}
              className={
                selectedGenres.includes(genre.name) ? 'bg-tertiary' : ''
              }
            />
          ))}
        </ul>
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
      </div>
      {loading ? <Loader /> : <ListMovies movies={movies} />}
    </motion.div>
  )
}

export default Home
