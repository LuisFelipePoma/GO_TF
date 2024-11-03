package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net"
	"os"
	"sort"
	"strconv"
	"sync"
	"time"
)

// Movie representa la estructura de una película.
type Movie struct {
	ID          int     `json:"id"`
	Keywords    string  `json:"keywords"`
	Characters  string  `json:"characters"`
	Actors      string  `json:"actors"`
	Director    string  `json:"director"`
	Crew        string  `json:"crew"`
	Genres      string  `json:"genres"`
	Overview    string  `json:"overview"`
	Title       string  `json:"title"`
	ImdbId      string  `json:"imdb_id"`
	VoteAverage float64 `json:"vote_average"`
}

// SimilarMovie representa una película similar con su ID y similaridad.
type SimilarMovie struct {
	ID         int     `json:"id"`
	Similarity float64 `json:"similarity"`
}

type MovieResponse struct {
	ID          int     `json:"id"`
	Title       string  `json:"title"`
	Characters  string  `json:"characters"`
	Actors      string  `json:"actors"`
	Director    string  `json:"director"`
	Genres      string  `json:"genres"`
	ImdbId      string  `json:"imdb_id"`
	VoteAverage float64 `json:"vote_average"`
}

var slaveNodes = []string{
	"localhost:8082",
	"localhost:8083",
	// Agrega más nodos esclavos según sea necesario
}

var movies []Movie

func main() {
	// Leer el archivo JSON con las películas una vez
	if err := loadMovies("../../database/data/data_clean.json"); err != nil {
		fmt.Println("Error loading movies:", err)
		return
	}

	// Iniciar la interfaz de usuario en consola
	for {
		fmt.Println("Ingrese el ID de la película para encontrar similares (o 'exit' para salir):")
		var input string
		fmt.Scanln(&input)
		if input == "exit" {
			break
		}

		movieID, err := strconv.Atoi(input)
		if err != nil {
			fmt.Println("ID de película inválido")
			continue
		}

		similarMoviesHandler(movieID)
	}
}

func loadMovies(filePath string) error {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("error reading JSON file: %w", err)
	}

	if err := json.Unmarshal(data, &movies); err != nil {
		return fmt.Errorf("error deserializing JSON: %w", err)
	}

	return nil
}

func similarMoviesHandler(movieID int) {
	start := time.Now()

	// Dividir índices de películas entre los esclavos
	numSlaves := len(slaveNodes)
	ranges := splitRanges(len(movies), numSlaves)

	var wg sync.WaitGroup
	results := make(chan []SimilarMovie, numSlaves)

	for i, node := range slaveNodes {
		wg.Add(1)
		go func(node string, startIdx, endIdx, movieID int) {
			defer wg.Done()
			result, err := getSimilarMoviesFromNode(node, startIdx, endIdx, movieID)
			if err == nil {
				results <- result
			}
		}(node, ranges[i][0], ranges[i][1], movieID)
	}

	wg.Wait()
	close(results)

	var combinedResults []SimilarMovie
	for result := range results {
		combinedResults = append(combinedResults, result...)
	}
	// Sort by similarity and return just the 10 first
	sort.Slice(combinedResults, func(i, j int) bool {
		return combinedResults[i].Similarity > combinedResults[j].Similarity
	})
	if len(combinedResults) > 10 {
		combinedResults = combinedResults[:10]
	}
	// Parse from SimilarMovie to MovieResponse
	var movieResponses []MovieResponse
	for _, similarMovie := range combinedResults {
		for _, movie := range movies {
			if movie.ID == similarMovie.ID {
				movieResponses = append(movieResponses, MovieResponse{
					ID:          similarMovie.ID,
					Title:       movie.Title,
					Characters:  movie.Characters,
					Actors:      movie.Actors,
					Director:    movie.Director,
					Genres:      movie.Genres,
					ImdbId:      movie.ImdbId,
					VoteAverage: movie.VoteAverage,
				})
				break
			}
		}
	}

	// Mostrar la respuesta en la consola
	fmt.Println("Películas similares encontradas:")
	for _, movieResponse := range movieResponses {
		fmt.Printf("Título: %s\n Descripción: %s\n Voto promedio: %.2f\n\n",
			movieResponse.Title, movieResponse.Genres, movieResponse.VoteAverage)
	}
	fmt.Printf("Request processed in %s\n", time.Since(start))
}

func splitRanges(totalMovies, numParts int) [][2]int {
	chunkSize := (totalMovies + numParts - 1) / numParts
	var ranges [][2]int
	for i := 0; i < totalMovies; i += chunkSize {
		end := i + chunkSize
		if end > totalMovies {
			end = totalMovies
		}
		ranges = append(ranges, [2]int{i, end})
	}
	return ranges
}

func getSimilarMoviesFromNode(node string, startIdx, endIdx, movieID int) ([]SimilarMovie, error) {
	conn, err := net.Dial("tcp", node)
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	task := struct {
		Movies      []Movie `json:"movies"`
		TargetMovie Movie   `json:"target_movie"`
	}{
		Movies:      movies[startIdx:endIdx],
		TargetMovie: movies[movieID],
	}

	data, err := json.Marshal(task)
	if err != nil {
		return nil, err
	}

	_, err = conn.Write(data)
	if err != nil {
		return nil, err
	}

	response, err := io.ReadAll(conn)
	if err != nil {
		return nil, err
	}

	var similarMovies []SimilarMovie
	err = json.Unmarshal(response, &similarMovies)
	if err != nil {
		return nil, err
	}

	return similarMovies, nil
}
