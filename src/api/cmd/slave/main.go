package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"math"
	"net"
	"os"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"
)

// Movie representa la estructura de una película.
type Movie struct {
	ID         int    `json:"id"`
	Keywords   string `json:"keywords"`
	Characters string `json:"characters"`
	Actors     string `json:"actors"`
	Director   string `json:"director"`
	Crew       string `json:"crew"`
	Genres     string `json:"genres"`
	Overview   string `json:"overview"`
	Title      string `json:"title"`
}

// SimilarMovie representa una película similar con su ID y título.
type SimilarMovie struct {
	ID    int    `json:"id"`
	Title string `json:"title"`
}

// CosineSimilarity calcula la similitud del coseno entre dos vectores.
func CosineSimilarity(vec1, vec2 map[string]float64) float64 {
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
func GetFeatureVector(features string) map[string]float64 {
	featureVector := make(map[string]float64)
	featureList := strings.Split(features, ", ")
	for _, feature := range featureList {
		featureVector[feature]++
	}
	return featureVector
}

// GetSimilarMovies retorna una lista de películas similares a la película dada.
func GetSimilarMovies(movies []Movie, movieID int) []SimilarMovie {
	start := time.Now()
	var targetMovie *Movie
	for _, movie := range movies {
		if movie.ID == movieID {
			targetMovie = &movie
			break
		}
	}
	if targetMovie == nil {
		return nil
	}

	targetFeatures := GetFeatureVector(targetMovie.Keywords + ", " + targetMovie.Characters + ", " + targetMovie.Actors + ", " + targetMovie.Director + ", " + targetMovie.Crew + ", " + targetMovie.Genres + ", " + targetMovie.Overview)

	similarities := make(map[int]float64)
	var mu sync.Mutex
	var wg sync.WaitGroup

	for _, movie := range movies {
		if movie.ID != movieID {
			wg.Add(1)
			go func(movie Movie) {
				defer wg.Done()
				features := GetFeatureVector(movie.Keywords + ", " + movie.Characters + ", " + movie.Actors + ", " + movie.Director + ", " + movie.Crew + ", " + movie.Genres + ", " + movie.Overview)
				similarity := CosineSimilarity(targetFeatures, features)
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

	var result []SimilarMovie
	for _, movie := range sortedMovies {
		for _, m := range movies {
			if m.ID == movie.movieID {
				result = append(result, SimilarMovie{ID: m.ID, Title: m.Title})
				break
			}
		}
	}
	fmt.Printf("Processed in %s\n", time.Since(start))
	return result
}

func handleConnection(conn net.Conn) {
	defer conn.Close()

	fmt.Println("Reading data....")

	// Decodificar el JSON recibido en una estructura Task
	var task struct {
		Movies      []Movie `json:"movies"`
		TargetMovie int     `json:"target_movie"`
	}
	decoder := json.NewDecoder(conn)
	if err := decoder.Decode(&task); err != nil {
		fmt.Println("Error al decodificar JSON:", err)
		return
	}
	fmt.Printf("Nodo Esclavo recibió tarea...")

	fmt.Println("Getting similar movies....")
	// Obtener películas similares
	similarMovies := GetSimilarMovies(task.Movies, task.TargetMovie)

	// Codificar el resultado en JSON y enviarlo de vuelta al cliente
	encoder := json.NewEncoder(conn)
	if err := encoder.Encode(similarMovies); err != nil {
		fmt.Println("Error al codificar JSON:", err)
		return
	}
	fmt.Printf("Nodo Esclavo envió resultado")
}

func main() {
	// Leer el puerto del usuario
	port := ""
	if len(os.Args) > 1 {
		port = os.Args[1]
	} else {
		reader := bufio.NewReader(os.Stdin)
		fmt.Print("Enter port: ")
		port, _ = reader.ReadString('\n')
		port = strings.TrimSpace(port)

		// Validar que el puerto sea un número válido
		if _, err := strconv.Atoi(port); err != nil {
			fmt.Println("Invalid port number")
			return
		}
	}

	// Iniciar el servidor TCP
	ln, err := net.Listen("tcp", ":"+port)
	if err != nil {
		fmt.Println("Error starting TCP server:", err)
		return
	}
	defer ln.Close()

	fmt.Println("Slave node listening on port", port)

	for {
		// Aceptar conexiones entrantes
		conn, err := ln.Accept()
		if err != nil {
			fmt.Println("Error accepting connection:", err)
			continue
		}
		go handleConnection(conn)
	}
}
