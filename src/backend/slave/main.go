package main

import (
	"bufio"
	"fmt"
	"github.com/LuisFelipePoma/Movies_Recomender_With_Golang/src/backend/slave/model"
	"github.com/LuisFelipePoma/Movies_Recomender_With_Golang/src/backend/types"
	jsoniter "github.com/json-iterator/go"
	"log"
	"math/rand"
	"net"
	"os"
	"sort"
	"strings"
	"sync"
	"time"
)

var json = jsoniter.ConfigFastest
var recommender = model.NewRecommender() // Create a new recommender instance

// Entry point of the program
func main() {
	// Read port from command line arguments or stdin
	port := os.Getenv("PORT")
	name := os.Getenv("NODE_NAME")
	if port == "" {
		log.Fatal("El puerto no está configurado en la variable de entorno PORT")
	}

	// Initialize TCP server
	ln, err := net.Listen("tcp", ":"+port)
	if err != nil {
		fmt.Println("Error starting TCP server:", err)
		return
	}
	defer ln.Close()
	fmt.Printf("Slave %s listening on port %s\n", name, port)
	// Show local addres
	fmt.Println("Local address:", ln.Addr())

	// Accept incoming connections
	for {
		conn, err := ln.Accept()
		if err != nil {
			fmt.Println("Error accepting connection:", err)
			continue
		}
		go handleConnection(conn)
	}
}

var dictFunction = map[types.TaskType]func(types.TaskDistributed) []types.MovieResponse{
	types.TaskRecomend:     getSimilarMovies,
	types.TaskSearch:       getMoviesSearch,
	types.TaskGet:          getNMovies,
	types.TaskUserRecomend: getUserRecommendations,
}

// handleConnection handles incoming connections
func handleConnection(conn net.Conn) {
	defer conn.Close()

	fmt.Println("Leyendo datos entrantes...")
	// Decode the JSON data using buffered reader
	var taskMaster types.TaskDistributed

	reader := bufio.NewReader(conn)
	decoder := json.NewDecoder(reader)

	if err := decoder.Decode(&taskMaster); err != nil {
		fmt.Println("Error en codificar JSON:", err)
		return
	}
	fmt.Printf("Esclavo recibio la tarea de: %s\n", taskMaster.Type)

	// Execute the task
	response := dictFunction[taskMaster.Type](taskMaster)

	// Send the result back using buffered writer
	writer := bufio.NewWriter(conn)
	encoder := json.NewEncoder(writer)
	if err := encoder.Encode(response); err != nil {
		fmt.Println("Error en parsear JSON:", err)
		return
	}
	writer.Flush()
}

// Services functions
func getSimilarMovies(taskMaster types.TaskDistributed) []types.MovieResponse {

	// get data from task
	data := taskMaster.Data

	fmt.Printf("Nodo Esclavo recibió %d\n peliculas", len(data.Movies))
	fmt.Println("Calculando las peliculas similares....")

	// Get similar movies
	similarMovies := recommender.GetSimilarMovies(data.Movies, data.TaskRecomendations.TargetMovie)
	fmt.Println("Se encontro", len(similarMovies), "similar movies")
	return similarMovies
}

// GetMoviesSearch search movies by query
func getMoviesSearch(taskMaster types.TaskDistributed) []types.MovieResponse {
	data := taskMaster.Data
	query := strings.ToLower(data.TaskSearch.Query)
	movies := taskMaster.Data.Movies

	fmt.Println("Buscando peliculas...")

	type result struct {
		movie      types.MovieResponse
		similarity float64
	}

	numWorkers := 5
	jobs := make(chan types.Movie, len(movies))
	results := make(chan result, len(movies))

	var wg sync.WaitGroup

	fieldWeights := []struct {
		field  func(types.Movie) string
		weight float64
	}{
		{func(m types.Movie) string { return m.Title }, 4},
		{func(m types.Movie) string { return m.Genres }, 3},
		{func(m types.Movie) string { return m.Actors }, 2},
		{func(m types.Movie) string { return m.Keywords }, 1},
	}

	for w := 0; w < numWorkers; w++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for movie := range jobs {
				totalSimilarity := 0.0
				for _, fw := range fieldWeights {
					if strings.Contains(strings.ToLower(fw.field(movie)), query) {
						totalSimilarity += fw.weight
					}
				}
				if totalSimilarity > 0 {
					results <- result{
						movie: types.MovieResponse{
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
							Similarity:  totalSimilarity,
						},
						similarity: totalSimilarity,
					}
				}
			}
		}()
	}

	for _, movie := range movies {
		jobs <- movie
	}
	close(jobs)
	wg.Wait()
	close(results)

	var resultsMovies []types.MovieResponse
	for res := range results {
		resultsMovies = append(resultsMovies, res.movie)
	}

	// Sort results by similarity descending
	sort.Slice(resultsMovies, func(i, j int) bool {
		return resultsMovies[i].Similarity > resultsMovies[j].Similarity
	})

	fmt.Println("Se encontraron", len(resultsMovies), "peliculas")
	return resultsMovies
}

// getNMovies returns the random N movies
func getNMovies(taskMaster types.TaskDistributed) []types.MovieResponse {
	data := taskMaster.Data
	n := data.Quantity
	movies := data.Movies
	genres := data.TaskSearch.Query

	fmt.Println("Obteniendo", n, "películas aleatorias")

	// Calcular el voto promedio
	totalVote := 0.0
	for _, movie := range movies {
		totalVote += movie.VoteAverage
	}
	averageVote := totalVote / float64(len(movies))

	var filteredMovies []types.MovieResponse

	for _, movie := range movies {
		// Verificar si la película tiene un voto mayor al promedio
		if movie.VoteAverage <= averageVote {
			continue
		}
		// example of genres => 'All, Comedy, Action'
		// example of genres => 'All, Comedy, Action'

		// If there is All, skip genres filter
		if !strings.Contains(genres, "All") {
			genres := strings.Split(genres, ",")
			// Verificar si la película tiene ambos de los géneros seleccionados
			genresMatch := true
			for _, genre := range genres {
				if !strings.Contains(movie.Genres, genre) {
					genresMatch = false
					break
				}
			}
			if !genresMatch {
				continue
			}
		}

		// Agregar la película a la lista filtrada
		filteredMovies = append(filteredMovies, types.MovieResponse{
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
		})
	}

	// Ajustar 'n' si es mayor que el número de películas filtradas
	if n > len(filteredMovies) {
		n = len(filteredMovies)
	}

	// Seleccionar aleatoriamente 'n' películas de las filtradas
	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(filteredMovies), func(i, j int) {
		filteredMovies[i], filteredMovies[j] = filteredMovies[j], filteredMovies[i]
	})

	selected := filteredMovies[:n]

	return selected
}

// getUserRecommendations returns the recommendations for a user
func getUserRecommendations(taskMaster types.TaskDistributed) []types.MovieResponse {
	data := taskMaster.Data
	// Get recommendations
	fmt.Println("Obteniendo recomendaciones para el usuario", data.TaskUserRecomendations.UserID)
	recommendations := model.RecommendItemsC(data.TaskUserRecomendations.UserRatings, data.TaskUserRecomendations.UserID, data.Quantity)
	// map ids of the movies
	movies := []types.MovieResponse{}
	for _, id := range recommendations {
		movieResponse := types.MovieResponse{
			ID: id,
		}
		movies = append(movies, movieResponse)
	}

	fmt.Println("Se encontraron", len(recommendations), "recomendaciones")
	return movies
}
