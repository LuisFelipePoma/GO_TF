package model

import (
	"fmt"
	"github.com/LuisFelipePoma/Movies_Recomender_With_Golang/src/backend/types"
	"math"
	"strings"
	"sync"
	"time"
)

// Recommender is a struct that contains the methods to calculate the similarity between movies.
type Recommender struct{}

func NewRecommender() *Recommender {
	return &Recommender{}
}

// CosineSimilarity calculates the cosine similarity between two vectors.
func (r *Recommender) CosineSimilarity(vec1, vec2 map[string]float64) float64 {
	var dotProduct, normA, normB float64
	// Calculate the dot product
	for key, val := range vec1 {
		// If the key exists in both vectors, calculate the dot product
		if val2, ok := vec2[key]; ok {
			dotProduct += val * val2
		}
		normA += val * val // Sum of squares of vector A
	}
	// Calculate the sum of squares of vector B
	for _, val := range vec2 {
		normB += val * val
	}
	// If one of the vectors is empty, return
	if normA == 0 || normB == 0 {
		return 0
	}
	// Calculate the cosine similarity
	return dotProduct / (math.Sqrt(normA) * math.Sqrt(normB))
}

// GetFeatureVector returns a map with the features of a movie.
func (r *Recommender) GetFeatureVector(features string) map[string]float64 {
	featureVector := make(map[string]float64)    // Create a map to store the features
	featureList := strings.Split(features, ", ") // Split the features
	// Count the number of times each feature appears
	for _, feature := range featureList {
		featureVector[feature]++
	}
	// Also called term frequency or tf
	return featureVector
}

// GetSimilarMovies returns a list of similar movies to a target movie.
func (r *Recommender) GetSimilarMovies(movies []types.Movie, targetMovie types.Movie) []types.MovieResponse {
	start := time.Now()
	// Get the feature vector of the target movie
	targetFeatures := r.GetFeatureVector(targetMovie.Keywords + ", " + targetMovie.Characters + ", " + targetMovie.Actors + ", " + targetMovie.Director + ", " + targetMovie.Crew + ", " + targetMovie.Genres + ", " + targetMovie.Overview)
	similarities := make(map[int]types.MovieResponse) // Create a map to store the similarities

	var mu sync.Mutex
	var wg sync.WaitGroup

	upperSimilarity := 0.0

	for _, movie := range movies {
		// If the movie is not the target movie, calculate the similarity
		if movie.ID != targetMovie.ID {
			wg.Add(1)
			// Calculate the similarity in a goroutine
			go func(movie types.Movie) {
				defer wg.Done()
				features := r.GetFeatureVector(movie.Keywords + ", " + movie.Characters + ", " + movie.Actors + ", " + movie.Director + ", " + movie.Crew + ", " + movie.Genres + ", " + movie.Overview)
				similarity := r.CosineSimilarity(targetFeatures, features)
				mu.Lock()
				similarities[movie.ID] = types.MovieResponse{
					ID:          movie.ID,
					Title:       movie.Title,
					Characters:  movie.Characters,
					Actors:      movie.Actors,
					Director:    movie.Director,
					Genres:      movie.Genres,
					ImdbId:      movie.ImdbId,
					VoteAverage: movie.VoteAverage,
					PosterPath:  movie.PosterPath,
					Overview:    movie.Overview,
					Similarity:  similarity,
				}
				if similarity > upperSimilarity {
					upperSimilarity = similarity
				}
				mu.Unlock()
			}(movie)
		}
	}
	// Wait for all goroutines to finish
	wg.Wait()

	// filter movies with more than 0.5 of similarity
	var filterMovies []types.MovieResponse // Create a slice
	for _, movie := range similarities {
		// normarlize similiarity
		movie.Similarity = movie.Similarity / upperSimilarity
		if movie.Similarity > 0.35 {
			filterMovies = append(filterMovies, movie)
		}
	}

	fmt.Printf("Se obtuvieron %d peliculas similares\n", len(filterMovies))
	fmt.Printf("Procesado en %s\n", time.Since(start))
	return filterMovies
}
