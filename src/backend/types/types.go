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

// TaskType represents the task type
const (
	TaskNone         TaskType = ""
	TaskRecomend     TaskType = "Recommend"
	TaskSearch       TaskType = "SearchQuery"
	TaskGet          TaskType = "GetMovies"
	TaskUserRecomend TaskType = "UserRecomend"
)

// TaskDistributed represents the task
type TaskDistributed struct {
	Type TaskType `json:"type"`
	Data TaskData `json:"data"`
}

// TaskData represents the task data structure
type TaskData struct {
	TaskRecomendations     *TaskRecomendations     `json:"recomendations,omitempty"`
	TaskSearch             *TaskSearchQuery        `json:"search,omitempty"`
	TaskUserRecomendations *TaskUserRecomendations `json:"user_recomendations,omitempty"`
	Quantity               int                     `json:"quantity"`
	Movies                 []Movie                 `json:"movies"`
}

// TaskRecomendations represents the recommendation task structure
type TaskRecomendations struct {
	MovieId     int   `json:"movie_id"`
	TargetMovie Movie `json:"movie"`
}

// TaskSearchQuery represents the search query task structure
type TaskSearchQuery struct {
	Query string `json:"query"`
}

// TaskUserRecomendations represents the user recommendation task structure
type TaskUserRecomendations struct {
	UserID      int             `json:"user_id"`
	User        map[int]float64 `json:"user"`
	UserRatings map[int]User    `json:"user2"`
}

// User represents the user structure for the user recommendation task
type User struct {
	ID      int
	Ratings map[int]float64
}

// Message represents the recommendation message structure
type Message struct {
	User            string   `json:"user"`
	Recommendations []string `json:"recommendations"`
}

// Task represents a task to be processed by a slave node
type Task struct {
	Slave    string
	TaskData TaskDistributed
}

// Result represents the result from a slave node
type Result struct {
	Response Response
	Error    error
}
