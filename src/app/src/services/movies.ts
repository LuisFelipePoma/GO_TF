import { URL_API, URL_TMDB } from '../consts/api'
import { type Movie, type Movies } from '../types/movies'
import type { TmdbResponse } from '../types/tmdb'
import { API_TMDB } from 'astro:env/server'

// Service to get all movies
export const getAllMovies = async (n: number): Promise<Movie[]> => {
  return await fetch(`${URL_API}/movies?n=${n}`)
    .then(response => response.json())
    .catch(error => error)
}

// Service to get by id Movie
export const getMovieById = async (id: number): Promise<Movie> => {
  return await fetch(`${URL_API}/movies?id=${id}`)
    .then(response => response.json())
    .catch(error => error)
}

// Service to get the recomendations by title
export const getRecommendations = async (title: string): Promise<Movies> => {
  return await fetch(`${URL_API}/movies/similar?title=${title}`).then(
    response => response.json().catch(error => error)
  )
}

// Service to get the image
export const getTmdbInfo = async (id: number): Promise<TmdbResponse> => {
  return await fetch(URL_TMDB(id), {
    headers: {
      Authorization: `Bearer ${API_TMDB}`
    }
  })
    .then(response => response.json())
    .catch(error => error)
}
