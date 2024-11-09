package model

import (
	"fmt"
	"github.com/LuisFelipePoma/Movies_Recomender_With_Golang/src/backend/types"
	"math"
	"sort"
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
func (r *Recommender) GetSimilarMovies(movies []types.Movie, targetMovie types.Movie) []types.SimilarMovie {
	start := time.Now()
	// Get the feature vector of the target movie
	targetFeatures := r.GetFeatureVector(targetMovie.Keywords + ", " + targetMovie.Characters + ", " + targetMovie.Actors + ", " + targetMovie.Director + ", " + targetMovie.Crew + ", " + targetMovie.Genres + ", " + targetMovie.Overview)
	similarities := make(map[int]float64) // Create a map to store the similarities

	var mu sync.Mutex
	var wg sync.WaitGroup

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
				similarities[movie.ID] = similarity // Store the similarity in the map
				mu.Unlock()
			}(movie)
		}
	}
	// Wait for all goroutines to finish
	wg.Wait()

	var sortedMovies []types.SimilarMovie // Create a slice
	// Store the movieID and similarity in
	for movieID, similarity := range similarities {
		sortedMovies = append(sortedMovies, types.SimilarMovie{ID: movieID, Similarity: similarity})
	}

	// Sort the movies by similarity
	sort.Slice(sortedMovies, func(i, j int) bool {
		return sortedMovies[i].Similarity > sortedMovies[j].Similarity
	})

	// Map the sortedMovies to a slice of SimilarMovie for the response
	var result []types.SimilarMovie
	for _, movie := range sortedMovies {
		for _, m := range movies {
			if m.ID == movie.ID {
				result = append(result, types.SimilarMovie{ID: m.ID, Similarity: movie.Similarity})
				break
			}
		}
	}
	fmt.Printf("Processed in %s\n", time.Since(start))
	return result
}
