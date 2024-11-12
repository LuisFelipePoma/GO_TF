import { API_TMDB, URL_API, URL_TMDB } from '../consts/api'
import { type Response } from '../types/movies'
import type { TmdbResponse } from '../types/tmdb'

// Service to get all movies
export const getAllMovies = async (n: number): Promise<Response> => {
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

// Service to get by Query
export const getMoviesByQuery = async (
  query: string,
  n: number
): Promise<Response> => {
  try {
    const response = await fetch(
      `${URL_API}/movies/search?query=${query}&n=${n}`
    )
    if (!response.ok) {
      throw new Error(`Error fetching movies by query: ${response.statusText}`)
    }
    return await response.json()
  } catch (error) {
    console.error('Error fetching movies by query:', error)
    throw error
  }
}

// Service to get the recommendations by title
export const getRecommendations = async (
  title: string,
  n: number
): Promise<Response> => {
  try {
    const response = await fetch(`${URL_API}/movies/similar`, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json'
      },
      body: JSON.stringify({ title, n })
    })
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
