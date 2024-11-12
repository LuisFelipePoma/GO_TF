package services

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"github.com/LuisFelipePoma/Movies_Recomender_With_Golang/src/backend/types"
	"golang.org/x/exp/rand"
	"os"
	"strconv"
	"strings"
	"time"
)

// Movies represents the structure of the movies service.
type Movies struct {
	Movies        []types.Movie
	RatingsLength int
	UserRatings   map[int]types.User
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

// Implementa el método LoadRatings
func (m *Movies) LoadRatings(filePath string) error {
	file, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		return err
	}

	userMap := make(map[int]types.User)
	for _, record := range records[1:] { // Saltar el encabezado
		userID, _ := strconv.Atoi(record[0])
		itemID, _ := strconv.Atoi(record[1])
		score, _ := strconv.ParseFloat(record[2], 64)

		if _, exists := userMap[userID]; !exists {
			userMap[userID] = types.User{ID: userID, Ratings: make(map[int]float64)}
		}
		userMap[userID].Ratings[itemID] = score
	}

	fmt.Println("Users:", len(userMap))
	fmt.Println("Total reviews:", len(records))

	m.UserRatings = userMap
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

func (m *Movies) GetRandomUserID() int {
	if len(m.UserRatings) == 0 {
		return 0 // or handle the empty case as needed
	}

	rand.Seed(uint64(time.Now().UnixNano()))
	index := rand.Intn(len(m.UserRatings))
	return m.UserRatings[index].ID
}
