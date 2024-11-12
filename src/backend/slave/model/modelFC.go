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
				similarity := cosineSimilarity(users[userIndex].Ratings, users[i].Ratings)
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

	for _, similarUser := range similarUsers {
		wg.Add(1)
		go func(similarUser int) {
			defer wg.Done()
			for itemID, rating := range users[similarUser].Ratings {
				// Si el usuario no ha calificado este ítem
				if _, exists := users[userIndex].Ratings[itemID]; !exists {
					mu.Lock()
					recommendations[itemID] += rating
					mu.Unlock()
				}
			}
		}(similarUser)
	}

	wg.Wait()

	// Ordenar las recomendaciones por las calificaciones acumuladas
	type kv struct {
		Key   int
		Value float64
	}
	var sortedRecommendations []kv
	for k, v := range recommendations {
		sortedRecommendations = append(sortedRecommendations, kv{k, v})
	}
	// Ordenar en orden descendente
	sort.Slice(sortedRecommendations, func(i, j int) bool {
		return sortedRecommendations[i].Value > sortedRecommendations[j].Value
	})

	// Devolver los índices de los ítems recomendados
	var recommendedItems []int
	for i := 0; i < numRecommendations && i < len(sortedRecommendations); i++ {
		recommendedItems = append(recommendedItems, sortedRecommendations[i].Key)
	}

	return recommendedItems
}

// FUNCTIONS

// Función para calcular la similitud coseno entre dos usuarios
func cosineSimilarity(user1, user2 map[int]float64) float64 {
	dotProduct := 0.0
	normA := 0.0
	normB := 0.0
	for itemID, rating1 := range user1 {
		if rating2, exists := user2[itemID]; exists {
			dotProduct += rating1 * rating2
			normA += rating1 * rating1
			normB += rating2 * rating2
		}
	}
	if normA == 0 || normB == 0 { // Evitar división por cero
		return 0
	}
	return dotProduct / (math.Sqrt(normA) * math.Sqrt(normB))
}

// Función para encontrar los usuarios más similares a un usuario dado
func mostSimilarUsers(users []types.User, userIndex int) []int {
	similarities := make(map[int]float64)
	for i := 0; i < len(users); i++ {
		if i != userIndex {
			similarity := cosineSimilarity(users[userIndex].Ratings, users[i].Ratings)
			similarities[i] = similarity
		}
	}

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
