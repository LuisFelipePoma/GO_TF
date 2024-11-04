package types

// Movie representa la estructura de una película.
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
}

// SimilarMovie representa una película similar con su ID y similaridad.
type SimilarMovie struct {
	ID         int     `json:"id"`
	Similarity float64 `json:"similarity"`
}

type MovieResponse struct {
	ID          int     `json:"id"`
	Title       string  `json:"title"`
	Characters  string  `json:"characters"`
	Actors      string  `json:"actors"`
	Director    string  `json:"director"`
	Genres      string  `json:"genres"`
	ImdbId      string  `json:"imdb_id"`
	VoteAverage float64 `json:"vote_average"`
}

type Response struct {
	Error         string          `json:"error"`
	MovieResponse []MovieResponse `json:"movie_response"`
	TargetMovie   string          `json:"target_movie"`
}

type Request struct {
	Option int    `json:"option"`
	Data   string `json:"data"`
}