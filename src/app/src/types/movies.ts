export interface Response {
  error?: string
  movie_response?: MovieResponse[]
  target_movie?: string
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
  poster_path?: string
  overview?: string
  similarity?: number
}
