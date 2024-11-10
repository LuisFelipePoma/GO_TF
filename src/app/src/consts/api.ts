export const URL_API = 'http://127.0.0.1:3000/api'
// export const URL_TMDB = ''
export const URL_IMG = (tmdb_path_img: string) =>
  `https://image.tmdb.org/t/p/w500/${tmdb_path_img}`

export const URL_TMDB = (id: number) =>
  `https://api.themoviedb.org/3/movie/${id}`


export const API_TMDB = import.meta.env.VITE_API_TMDB