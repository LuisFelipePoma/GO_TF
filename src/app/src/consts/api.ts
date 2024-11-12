// export const URL_API = 'http://127.0.0.1:3000/api'
export const URL_API = 'https://b3b5h0z3-3000.brs.devtunnels.ms/api'
// export const URL_TMDB = ''
export const URL_IMG = (tmdb_path_img: string, size: number = 200) =>
  `https://image.tmdb.org/t/p/w${size}${tmdb_path_img}`

export const URL_TMDB = (id: number) =>
  `https://api.themoviedb.org/3/movie/${id}`

// export const API_TMDB = import.meta.env.VITE_API_TMDB
export const API_TMDB = "eyJhbGciOiJIUzI1NiJ9.eyJhdWQiOiI2MWI5NzEyMWIzZTY1MTI3N2JiMTkzOWEyZGI2M2VkMyIsIm5iZiI6MTczMDYwOTI2OS43OTA3MTcsInN1YiI6IjY0YTcwNzE5OTU3ZTZkMDEzOWNmMDc2ZCIsInNjb3BlcyI6WyJhcGlfcmVhZCJdLCJ2ZXJzaW9uIjoxfQ.AvK5oxP2qIFQ4BbNpOsIR-Sj1bAZoqDx_VZPTY1Lx0w"
