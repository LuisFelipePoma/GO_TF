package model

import (
	"math"
	"sort"
	"sync"

	"github.com/LuisFelipePoma/Movies_Recomender_With_Golang/src/backend/types"
)

// Función para encontrar los usuarios más similares a un usuario dado
func mostSimilarUsersC(users map[int]types.User, userIndex int) []int {
	similarities := make(map[int]float64)
	var wg sync.WaitGroup
	mu := &sync.Mutex{}

	for i := 0; i < len(users); i++ {
		if i != userIndex {
			wg.Add(1)
			go func(i int) {
				defer wg.Done()
				similarity := pearsonCorrelation(users[userIndex].Ratings, users[i].Ratings)
				mu.Lock()
				similarities[i] = similarity
				mu.Unlock()
			}(i)
		}
	}

	wg.Wait()

	// Ordenar los usuarios por similitud
	type kv struct {
		Key   int
		Value float64
	}
	var sortedSimilarities []kv
	for k, v := range similarities {
		sortedSimilarities = append(sortedSimilarities, kv{k, v})
	}
	// Ordenar en orden descendente
	sort.Slice(sortedSimilarities, func(i, j int) bool {
		return sortedSimilarities[i].Value > sortedSimilarities[j].Value
	})

	// Devolver los índices de los usuarios más similares
	var mostSimilar []int
	for _, kv := range sortedSimilarities {
		mostSimilar = append(mostSimilar, kv.Key)
	}
	return mostSimilar
}

// Función para recomendar ítems a un usuario basado en usuarios similares
func RecommendItemsC(users map[int]types.User, userIndex int, numRecommendations int) []int {
	similarUsers := mostSimilarUsersC(users, userIndex)

	recommendations := make(map[int]float64)
	var wg sync.WaitGroup
	mu := &sync.Mutex{}

	// Calculate similarity scores
	similarities := make(map[int]float64)
	for _, similarUser := range similarUsers {
		similarity := calculateCosineSimilarity(users[userIndex].Ratings, users[similarUser].Ratings)
		similarities[similarUser] = similarity
	}

	for _, similarUser := range similarUsers {
		wg.Add(1)
		go func(similarUser int) {
			defer wg.Done()
			similarity := similarities[similarUser]
			for itemID, rating := range users[similarUser].Ratings {
				// If the user has not rated this item
				if _, exists := users[userIndex].Ratings[itemID]; !exists {
					mu.Lock()
					recommendations[itemID] += similarity * rating
					mu.Unlock()
				}
			}
		}(similarUser)
	}

	wg.Wait()

	// Sort the recommendations by the weighted scores
	type kv struct {
		Key   int
		Value float64
	}
	var sortedRecommendations []kv
	for k, v := range recommendations {
		sortedRecommendations = append(sortedRecommendations, kv{k, v})
	}
	sort.Slice(sortedRecommendations, func(i, j int) bool {
		return sortedRecommendations[i].Value > sortedRecommendations[j].Value
	})

	// Return the top N recommended item IDs
	var recommendedItems []int
	for i := 0; i < numRecommendations && i < len(sortedRecommendations); i++ {
		recommendedItems = append(recommendedItems, sortedRecommendations[i].Key)
	}

	return recommendedItems
}

// FUNCTIONS
func pearsonCorrelation(user1, user2 map[int]float64) float64 {
	var sum1, sum2, sum1Sq, sum2Sq, pSum float64
	var n int

	// Encontrar ítems comunes
	for itemID, rating1 := range user1 {
		if rating2, ok := user2[itemID]; ok {
			sum1 += rating1
			sum2 += rating2
			sum1Sq += rating1 * rating1
			sum2Sq += rating2 * rating2
			pSum += rating1 * rating2
			n++
		}
	}

	if n == 0 {
		return 0
	}

	numerator := pSum - (sum1 * sum2 / float64(n))
	denominator := math.Sqrt((sum1Sq - sum1*sum1/float64(n)) * (sum2Sq - sum2*sum2/float64(n)))
	if denominator == 0 {
		return 0
	}
	return numerator / denominator
}

// Calcula la similitud coseno entre dos conjuntos de valoraciones de usuario
func calculateCosineSimilarity(ratings1 map[int]float64, ratings2 map[int]float64) float64 {
	dotProduct := 0.0
	sumSquares1 := 0.0
	sumSquares2 := 0.0

	for itemID, rating1 := range ratings1 {
		if rating2, exists := ratings2[itemID]; exists {
			dotProduct += rating1 * rating2
			sumSquares1 += rating1 * rating1
			sumSquares2 += rating2 * rating2
		}

	}

	if sumSquares1 == 0 || sumSquares2 == 0 {
		return 0.0
	}
	return dotProduct / (math.Sqrt(sumSquares1) * math.Sqrt(sumSquares2))
}
