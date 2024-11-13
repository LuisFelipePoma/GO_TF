import { TmdbResponse } from '../types/tmdb'

export const MovieInfoDummie: TmdbResponse = {
  adult: false,
  backdrop_path: '/path/to/backdrop.jpg',
  belongs_to_collection: null,
  budget: 100000000,
  genres: [
    { id: 28, name: 'Action' },
    { id: 12, name: 'Adventure' }
  ],
  homepage: 'https://www.example.com',
  id: -1,
  imdb_id: 'tt1234567',
  original_language: 'en',
  original_title: '--------',
  overview: 'This is a sample overview of the movie.',
  popularity: 150.5,
  poster_path: '/path/to/poster.jpg',
  production_companies: [
    {
      id: 1,
      logo_path: '/logo.png',
      name: 'Example Productions',
      origin_country: 'US'
    }
  ],
  production_countries: [
    {
      iso_3166_1: 'US',
      name: 'United States of America'
    }
  ],
  release_date: new Date('2023-01-01'),
  revenue: 500000000,
  runtime: 120,
  spoken_languages: [
    {
      iso_639_1: 'en',
      name: 'English'
    }
  ],
  status: 'Released',
  tagline: 'An example tagline.',
  title: 'Movie Not Found',
  video: false,
  vote_average: 0.5,
  vote_count: 2500
}
