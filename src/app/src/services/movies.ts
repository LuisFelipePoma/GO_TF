import { API_TMDB, URL_API, URL_TMDB } from '../consts/api'
import { type Movie, type Movies } from '../types/movies'
import type { TmdbResponse } from '../types/tmdb'

// Service to get all movies
export const getAllMovies = async (n: number): Promise<Movie[]> => {
  try {
    const response = await fetch(`${URL_API}/movies?n=${n}`)
    if (!response.ok) {
      throw new Error(`Error fetching movies: ${response.statusText}`)
    }
    return await response.json()
  } catch (error) {
    console.error('Error fetching all movies:', error)
    throw error
  }
}

// Service to get by id Movie
export const getMovieById = async (id: number): Promise<Movie> => {
  try {
    const response = await fetch(`${URL_API}/movies?id=${id}`)
    if (!response.ok) {
      throw new Error(`Error fetching movie by id: ${response.statusText}`)
    }
    return await response.json()
  } catch (error) {
    console.error('Error fetching movie by id:', error)
    throw error
  }
}

// Service to get the recommendations by title
export const getRecommendations = async (title: string): Promise<Movies> => {
  try {
    const response = await fetch(`${URL_API}/movies/similar?title=${title}`)
    if (!response.ok) {
      throw new Error(`Error fetching recommendations: ${response.statusText}`)
    }
    return await response.json()
  } catch (error) {
    console.error('Error fetching recommendations:', error)
    throw error
  }
}

// Service to get the image
export const getTmdbInfo = async (id: number): Promise<TmdbResponse> => {
  try {
    const response = await fetch(URL_TMDB(id), {
      headers: {
        Authorization: `Bearer ${API_TMDB}`
      }
    })
    if (!response.ok) {
      throw new Error(`Error fetching TMDB info: ${response.statusText}`)
    }
    return await response.json()
  } catch (error) {
    console.error('Error fetching TMDB info:', error)
    throw error
  }
}
