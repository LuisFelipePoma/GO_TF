package main

import (
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"os"
	"strconv"
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

var slaveNodes = []string{
	"localhost:8082",
	"localhost:8083",
	// Agrega más nodos esclavos según sea necesario
}

var movies []Movie

func main() {
	// Leer el archivo JSON con las películas una vez
	data, err := os.ReadFile("../../database/data/data_clean.json")
	if err != nil {
		fmt.Println("Error reading JSON file:", err)
		return
	}

	if err := json.Unmarshal(data, &movies); err != nil {
		fmt.Println("Error deserializing JSON:", err)
		return
	}

	http.HandleFunc("/similar_movies", similarMoviesHandler)
	http.ListenAndServe(":8080", nil)
}

func similarMoviesHandler(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	movieIDStr := r.URL.Query().Get("id")
	movieID, err := strconv.Atoi(movieIDStr)
	if err != nil {
		http.Error(w, "Invalid movie ID", http.StatusBadRequest)
		return
	}

	// Dividir películas entre los esclavos
	numSlaves := len(slaveNodes)
	movieChunks := splitMovies(movies, numSlaves)

	var wg sync.WaitGroup
	results := make(chan []SimilarMovie, numSlaves)
	//
	fmt.Println("Handling petition for:", movieID)
	for i, node := range slaveNodes {
		wg.Add(1)
		go func(node string, chunk []Movie) {
			// func(node string, chunk []Movie) {
			defer wg.Done()
			result, err := getSimilarMoviesFromNode(node, chunk, movieID)
			if err == nil {
				results <- result
			}
		}(node, movieChunks[i])
	}

	wg.Wait()
	close(results)

	var combinedResults []SimilarMovie
	for result := range results {
		combinedResults = append(combinedResults, result...)
	}

	// Enviar la respuesta al cliente
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(combinedResults)
	fmt.Printf("Request processed in %s\n", time.Since(start))
}

func splitMovies(movies []Movie, numParts int) [][]Movie {
	chunkSize := (len(movies) + numParts - 1) / numParts
	var chunks [][]Movie
	for i := 0; i < len(movies); i += chunkSize {
		end := i + chunkSize
		if end > len(movies) {
			end = len(movies)
		}
		chunks = append(chunks, movies[i:end])
	}
	return chunks
}

func getSimilarMoviesFromNode(node string, movies []Movie, movieID int) ([]SimilarMovie, error) {
	conn, err := net.Dial("tcp", node)
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	task := struct {
		Movies      []Movie `json:"movies"`
		TargetMovie int     `json:"target_movie"`
	}{
		Movies:      movies,
		TargetMovie: movieID,
	}
	// Codificar la tarea en JSON y enviarla al servidor
	encoder := json.NewEncoder(conn)
	if err := encoder.Encode(task); err != nil {
		fmt.Println("Error al codificar JSON:", err)
		return nil, err
	}
	fmt.Printf("Nodo maestro envió la pelicula: %+v\n", task)

	// Decodificar la respuesta del servidor
	var result []SimilarMovie
	decoder := json.NewDecoder(conn)
	if err := decoder.Decode(&result); err != nil {
		fmt.Println("Error al decodificar JSON:", err)
		return nil, err
	}
	fmt.Printf("Nodo maestro recibió resultados")

	return result, nil
}
