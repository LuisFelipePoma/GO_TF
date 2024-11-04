package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"os"
	"os/exec"
	"sort"
	"strings"
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
	"localhost:8084",
	// Agrega más nodos esclavos según sea necesario
}

var movies []Movie
var recommendations []MovieResponse
var lastRecommendedMovieTitle string

func main() {
	// Leer el archivo JSON con las películas una vez
	if err := loadMovies("../../database/data/data_clean.json"); err != nil {
		fmt.Println("Error loading movies:", err)
		return
	}

	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Println("Seleccione una opción:")
		fmt.Println("1. Recomendar en base a una película")
		fmt.Println("2. Mostrar recientes recomendaciones")
		fmt.Println("3. Filtrar por género las últimas recomendaciones")
		fmt.Println("4. Filtrar por voteAverage las últimas recomendaciones")
		fmt.Println("5. Salir")
		fmt.Print("Opción: ")
		optionStr, _ := reader.ReadString('\n')
		optionStr = strings.TrimSpace(optionStr)
		option := 0
		fmt.Sscanf(optionStr, "%d", &option)

		switch option {
		case 1:
			fmt.Println("Ingrese el título de la película:")
			title, _ := reader.ReadString('\n')
			title = strings.TrimSpace(title)
			recommendations = recomendMoviesByTitle(title)
			fmt.Println("Recomendaciones:")
			printMovieDetails(recommendations)

		case 2:
			clearConsole()
			if len(recommendations) == 0 {
				fmt.Print("No hay recomendaciones recientes.\n\n")
				continue
			}
			fmt.Printf("Últimas recomendaciones basadas en la película: %s\n", lastRecommendedMovieTitle)
			printMovieDetails(recommendations)

		case 3:
			clearConsole()
			if len(recommendations) == 0 {
				fmt.Print("No hay recomendaciones para filtrar.\n\n")
				continue
			}
			fmt.Println("Ingrese el género:")
			genre, _ := reader.ReadString('\n')
			genre = strings.TrimSpace(genre)
			response := filterMoviesByGenre(recommendations, genre)
			fmt.Println("Recomendaciones filtradas por género:")
			printMovieDetails(response)

		case 4:
			clearConsole()
			if len(recommendations) == 0 {
				fmt.Print("No hay recomendaciones para filtrar.\n\n")
				continue
			}
			fmt.Println("Ingrese el voteAverage mínimo:")
			voteAverageStr, _ := reader.ReadString('\n')
			voteAverageStr = strings.TrimSpace(voteAverageStr)
			minVoteAverage := 0.0
			fmt.Sscanf(voteAverageStr, "%f", &minVoteAverage)
			response := filterMoviesByVoteAverage(recommendations, minVoteAverage)
			fmt.Println("Recomendaciones filtradas por voteAverage:")
			printMovieDetails(response)

		case 5:
			fmt.Println("Saliendo...")
			return
		default:
			fmt.Println("Opción no válida, intente de nuevo.")
		}
	}
}

// OPTIONS
func filterMoviesByVoteAverage(movies []MovieResponse, minVoteAverage float64) []MovieResponse {
	var filteredMovies []MovieResponse
	for _, movie := range movies {
		if movie.VoteAverage >= minVoteAverage {
			filteredMovies = append(filteredMovies, movie)
		}
	}
	return filteredMovies
}

func recomendMoviesByTitle(title string) []MovieResponse {
	var response []MovieResponse
	for _, movie := range movies {
		if strings.EqualFold(strings.ToLower(movie.Title), strings.ToLower(title)) {
			response = similarMoviesHandler(movie.ID)
			lastRecommendedMovieTitle = movie.Title
			break
		}
	}
	return response
}

func filterMoviesByGenre(movies []MovieResponse, genre string) []MovieResponse {
	var filteredMovies []MovieResponse
	for _, movie := range movies {
		if strings.Contains(strings.ToLower(movie.Genres), strings.ToLower(genre)) {
			filteredMovies = append(filteredMovies, movie)
		}
	}
	return filteredMovies
}

// OTHER

func clearConsole() {
	cmd := exec.Command("cmd", "/c", "cls")
	cmd.Stdout = os.Stdout
	cmd.Run()
}
func printMovieDetails(movies []MovieResponse) {
	for _, movie := range movies {
		fmt.Printf("Title: %s\n", movie.Title)
		fmt.Printf("Vote Average: %.2f\n", movie.VoteAverage)
		fmt.Printf("Genres: %s\n", movie.Genres)
		fmt.Println("-----------------------------")
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

func similarMoviesHandler(movieID int) []MovieResponse {
	start := time.Now()

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

	sort.Slice(combinedResults, func(i, j int) bool {
		return combinedResults[i].Similarity > combinedResults[j].Similarity
	})
	if len(combinedResults) > 10 {
		combinedResults = combinedResults[:10]
	}

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
	fmt.Printf("Request processed in %s\n", time.Since(start))
	return movieResponses
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

	// get movie by id
	var target Movie
	for _, movie := range movies {
		if movie.ID == movieID {
			target = movie
			break
		}
	}

	task := struct {
		Movies      []Movie `json:"movies"`
		TargetMovie Movie   `json:"target_movie"`
	}{
		Movies:      movies[startIdx:endIdx],
		TargetMovie: target}

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
