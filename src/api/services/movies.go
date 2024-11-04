package services

import (
	"github.com/LuisFelipePoma/Movies_Recomender_With_Golang/tree/main/src/api/types"
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"github.com/LuisFelipePoma/Movies_Recomender_With_Golang/tree/main/src/api/master"
)

// Movies representa la estructura de las películas.
type Movies struct {
	Movies            []types.Movie
	Recommendations   []types.MovieResponse
	LastRecomendation string
}

// NewMovies crea una nueva instancia de Movies.
func NewMovies() *Movies {
	return &Movies{
		Movies: []types.Movie{},
	}
}

// LoadMovies carga las películas desde un archivo JSON.
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

// SetRecomendations establece las recomendaciones de películas.
func (m *Movies) SetRecomendations(Recommendations []types.MovieResponse) {
	m.Recommendations = Recommendations
}

// SetLastRecommendedMovieTitle establece el título de la última película recomendada.
func (m *Movies) SetLastRecommendedMovieTitle(title string) {
	m.LastRecomendation = title
}

// GetMovies retorna todas las n películas.
func (m *Movies) GetMovies(n int) []types.Movie {
	return m.Movies[:n]
}

// GetMovieByTitle retorna una película por su título.
func (m *Movies) GetMovieByID(movieID int) *types.Movie {
	for _, movie := range m.Movies {
		if movie.ID == movieID {
			return &movie
		}
	}
	return nil
}

// IsEmptyRecommendationsverifica si no hay recomendaciones.
func (m *Movies) IsEmptyRecommendations() bool {
	if len(m.Recommendations) == 0 {
		fmt.Println("No hay recomendaciones recientes.")
		return true
	}
	return false
}

func (m *Movies) PrintRecomendationsDetails() {
	for _, movie := range m.Recommendations {
		fmt.Printf("Title: %s\n", movie.Title)
		fmt.Printf("Vote Average: %.2f\n", movie.VoteAverage)
		fmt.Printf("Genres: %s\n", movie.Genres)
		fmt.Println("-----------------------------")
	}
}

func (m *Movies) PrintMoviesDetails(movies []types.MovieResponse)  {
	for _, movie := range movies {
		fmt.Printf("Title: %s\n", movie.Title)
		fmt.Printf("Vote Average: %.2f\n", movie.VoteAverage)
		fmt.Printf("Genres: %s\n", movie.Genres)
		fmt.Println("-----------------------------")
	}
}


func (m *Movies) GetRecomendationsByGenre(genre string) []types.MovieResponse {
	var filteredMovies []types.MovieResponse
	for _, movie := range m.Recommendations {
		if strings.Contains(strings.ToLower(movie.Genres), strings.ToLower(genre)) {
			filteredMovies = append(filteredMovies, movie)
		}
	}
	return filteredMovies
}

func (m *Movies) GetMoviesByVoteAverage(minVoteAverage float64) []types.MovieResponse {
	var filteredMovies []types.MovieResponse
	for _, movie := range m.Recommendations {
		if movie.VoteAverage >= minVoteAverage {
			filteredMovies = append(filteredMovies, movie)
		}
	}
	return filteredMovies
}

func (m *Movies) RecomendMoviesByTitle(title string) ([]types.MovieResponse, string) {
	
	var response []types.MovieResponse
	for _, movie := range m.Recommendations {
		if strings.EqualFold(strings.ToLower(movie.Title), strings.ToLower(title)) {
			// response = similarMoviesHandler(movie.ID)
			response = nil
			return response, movie.Title
		}
	}
	return response, ""
}