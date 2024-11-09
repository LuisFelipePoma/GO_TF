package services

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/LuisFelipePoma/Movies_Recomender_With_Golang/src/backend/types"
	"golang.org/x/exp/rand"
)

// Movies represents the structure of the movies service.
type Movies struct {
	Movies            []types.Movie
	Recommendations   []types.MovieResponse
	LastRecomendation string
}

// NewMovies creates a new Movies service.
func NewMovies() *Movies {
	return &Movies{
		Movies: []types.Movie{},
	}
}

// LoadMovies load movies from a JSON file.
func (m *Movies) LoadMovies(filePath string) error {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("error reading JSON file: %w", err)
	}
	if err := json.Unmarshal(data, &m.Movies); err != nil {
		return fmt.Errorf("error deserializing JSON: %w", err)
	}

	return nil
}

// Get All Movies (n: number)
func (m *Movies) GetAllMovies(n int) []types.Movie {
	rand.Seed(uint64(time.Now().UnixNano())) // Seed the random number generator
	
	if n >= len(m.Movies) {
		return m.Movies
	}

	randomMovies := make([]types.Movie, 0, n)
	indices := rand.Perm(len(m.Movies))[:n]

	for _, i := range indices {
		randomMovies = append(randomMovies, m.Movies[i])
	}

	return randomMovies
}

// GetMovieByTitle returns a movie by its id.
func (m *Movies) GetMovieByID(movieID int) *types.Movie {
	// string to int
	for _, movie := range m.Movies {
		if movie.ID == movieID {
			return &movie
		}
	}
	return nil
}

// GetMovieByTitle returns a movie by its title.
func (m *Movies) GetMovieByTitle(title string) *types.Movie {
	for _, movie := range m.Movies {
		if strings.EqualFold(movie.Title, title) {
			return &movie
		}
	}
	return nil
}

// GetMoviesByGenre returns a list of movies by genre.
func (m *Movies) GetRecomendationsByGenre(genre string) []types.MovieResponse {
	var filteredMovies []types.MovieResponse
	for _, movie := range m.Recommendations {
		if strings.Contains(strings.ToLower(movie.Genres), strings.ToLower(genre)) {
			filteredMovies = append(filteredMovies, movie)
		}
	}
	return filteredMovies
}

// GetMoviesByVoteAverage returns a list of movies by vote average.
func (m *Movies) GetMoviesByVoteAverage(voteAverageStr string) []types.MovieResponse {
	minVoteAverage := 0.0
	fmt.Sscanf(voteAverageStr, "%f", &minVoteAverage)
	var filteredMovies []types.MovieResponse
	for _, movie := range m.Recommendations {
		if movie.VoteAverage >= minVoteAverage {
			filteredMovies = append(filteredMovies, movie)
		}
	}
	return filteredMovies
}
