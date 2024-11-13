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
	Similarity  float64 `json:"similarity"`
}

// Response represents the structure of a response.
type Response struct {
	Error         string          `json:"error"`
	MovieResponse []MovieResponse `json:"movie_response"`
	TargetMovie   string          `json:"target_movie"`
	UserID        int             `json:"user_id"`
}

// DATA TASK
type TaskType string

const (
	TaskNone         TaskType = ""
	TaskRecomend     TaskType = "Recommend"
	TaskSearch       TaskType = "SearchQuery"
	TaskGet          TaskType = "GetMovies"
	TaskUserRecomend TaskType = "UserRecomend"
)

type TaskDistributed struct {
	Type TaskType `json:"type"`
	Data TaskData `json:"data"`
}

type TaskData struct {
	TaskRecomendations     *TaskRecomendations     `json:"recomendations,omitempty"`
	TaskSearch             *TaskSearchQuery        `json:"search,omitempty"`
	TaskUserRecomendations *TaskUserRecomendations `json:"user_recomendations,omitempty"`
	Quantity               int                     `json:"quantity"`
	Movies                 []Movie                 `json:"movies"`
}

type TaskRecomendations struct {
	MovieId     int   `json:"movie_id"`
	TargetMovie Movie `json:"movie"`
}

type TaskSearchQuery struct {
	Query string `json:"query"`
}

// TaskUserRecomendations representa la tarea de recomendaciones para un usuario
type TaskUserRecomendations struct {
	UserID      int             `json:"user_id"`
	User        map[int]float64 `json:"user"`
	UserRatings map[int]User    `json:"user2"`
}

type User struct {
	ID      int
	Ratings map[int]float64
}
