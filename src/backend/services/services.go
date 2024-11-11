package services

import (
	"encoding/json"
	"fmt"
	"github.com/LuisFelipePoma/Movies_Recomender_With_Golang/src/backend/types"
	"os"
	"strings"
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
	fmt.Printf("Se cargaron %d peliculas.\n", len(m.Movies))

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
