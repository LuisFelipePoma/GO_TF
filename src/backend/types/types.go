package types

// Movie represents the structure of a movie.
type Movie struct {
	ID          int     `json:"id"`
	Keywords    string  `json:"keywords"`
	Characters  string  `json:"characters"`
	Actors      string  `json:"actors"`
	Director    string  `json:"director"`
	Crew        string  `json:"crew"`
	Genres      string  `json:"genres"`
	Overview    string  `json:"overview"`
	Title       string  `json:"title"`
	ImdbId      string  `json:"imdb_id"`
	VoteAverage float64 `json:"vote_average"`
	PosterPath  string  `json:"poster_path"`
}

// SimilarMovie represents the structure of a similar movie.
type SimilarMovie struct {
	ID         int     `json:"id"`
	Similarity float64 `json:"similarity"`
}

// MovieResponse represents the structure of a movie response.
type MovieResponse struct {
	ID          int     `json:"id"`
	Title       string  `json:"title"`
	Characters  string  `json:"characters"`
	Actors      string  `json:"actors"`
	Director    string  `json:"director"`
	Genres      string  `json:"genres"`
	ImdbId      string  `json:"imdb_id"`
	VoteAverage float64 `json:"vote_average"`
	PosterPath  string  `json:"poster_path"`
	Overview    string  `json:"overview"`
}

// Response represents the structure of a response.
type Response struct {
	Error         string          `json:"error"`
	MovieResponse []MovieResponse `json:"movie_response"`
	TargetMovie   string          `json:"target_movie"`
}

// Request represents the structure of a request.
type Request struct {
	TargetMovie Movie   `json:"movie"`
	Movies      []Movie `json:"movies"`
}
