package model

import (
	"fmt"
	"math"
	"sort"
	"strings"
	"sync"
	"time"
	"github.com/LuisFelipePoma/Movies_Recomender_With_Golang/src/api/types"
)

type Recommender struct{}

func NewRecommender() *Recommender {
	return &Recommender{}
}

// CosineSimilarity calcula la similitud del coseno entre dos vectores.
func (r *Recommender) CosineSimilarity(vec1, vec2 map[string]float64) float64 {
	var dotProduct, normA, normB float64
	for key, val := range vec1 {
		if val2, ok := vec2[key]; ok {
			dotProduct += val * val2
		}
		normA += val * val
	}
	for _, val := range vec2 {
		normB += val * val
	}
	if normA == 0 || normB == 0 {
		return 0
	}
	return dotProduct / (math.Sqrt(normA) * math.Sqrt(normB))
}

// GetFeatureVector convierte una cadena de características en un vector de características.
func (r *Recommender) GetFeatureVector(features string) map[string]float64 {
	featureVector := make(map[string]float64)
	featureList := strings.Split(features, ", ")
	for _, feature := range featureList {
		featureVector[feature]++
	}
	return featureVector
}

// GetSimilarMovies retorna una lista de películas similares a la película dada.
func (r *Recommender) GetSimilarMovies(movies []types.Movie, targetMovie types.Movie) []types.SimilarMovie {
	start := time.Now()

	targetFeatures := r.GetFeatureVector(targetMovie.Keywords + ", " + targetMovie.Characters + ", " + targetMovie.Actors + ", " + targetMovie.Director + ", " + targetMovie.Crew + ", " + targetMovie.Genres + ", " + targetMovie.Overview)

	similarities := make(map[int]float64)
	var mu sync.Mutex
	var wg sync.WaitGroup

	for _, movie := range movies {
		if movie.ID != targetMovie.ID {
			wg.Add(1)
			go func(movie types.Movie) {
				defer wg.Done()
				features := r.GetFeatureVector(movie.Keywords + ", " + movie.Characters + ", " + movie.Actors + ", " + movie.Director + ", " + movie.Crew + ", " + movie.Genres + ", " + movie.Overview)
				similarity := r.CosineSimilarity(targetFeatures, features)
				mu.Lock()
				similarities[movie.ID] = similarity
				mu.Unlock()
			}(movie)
		}
	}

	wg.Wait()

	type movieSimilarity struct {
		movieID    int
		similarity float64
	}

	var sortedMovies []movieSimilarity
	for movieID, similarity := range similarities {
		sortedMovies = append(sortedMovies, movieSimilarity{movieID, similarity})
	}

	sort.Slice(sortedMovies, func(i, j int) bool {
		return sortedMovies[i].similarity > sortedMovies[j].similarity
	})

	var result []types.SimilarMovie
	for _, movie := range sortedMovies {
		for _, m := range movies {
			if m.ID == movie.movieID {
				result = append(result, types.SimilarMovie{ID: m.ID, Similarity: movie.similarity})
				break
			}
		}
	}
	fmt.Printf("Processed in %s\n", time.Since(start))
	return result
}
