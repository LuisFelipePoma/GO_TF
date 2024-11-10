export interface Movie {
  id?: number
  keywords?: string
  characters?: string
  actors?: string
  director?: string
  crew?: string
  genres?: string
  overview?: string
  title?: string
  imdb_id?: string
  vote_average?: number
	poster_path?: string
}

export interface Movies {
  error?: string
  movie_response?: MovieResponse[]
}

export interface MovieResponse {
  id?: number
  title?: string
  characters?: string
  actors?: string
  director?: string
  genres?: string
  imdb_id?: string
  vote_average?: number
}
