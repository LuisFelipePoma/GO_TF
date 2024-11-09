import { URL_API } from '../consts/api'
import { type Movie } from '../types/movies'

// Service to get all movies
export const getAllMovies = async (n: number): Promise<Movie[]> => {
  return await fetch(`${URL_API}/movies?n=${n}`)
    .then(response => response.json())
    .catch(error => {
      console.error('Error:', error)
      return null
    })
}

// Service to get by id Movie
export const getMovieById = async (id: number): Promise<Movie> => {
  return await fetch(`${URL_API}/movies?id=${id}`)
    .then(response => response.json())
    .catch(error => {
      console.error('Error:', error)
      return null
    })
}

// Service to get the recomendations by title
export const getRecommendations = async (title: string): Promise<Movie[]> => {
  return await fetch(`${URL_API}/movies/similar?title=${title}`).then(
    response => response.json()
  )
}
