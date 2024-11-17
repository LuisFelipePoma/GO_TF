package main

import (
	"bufio"
	"fmt"
	"log"
	"math/rand"
	"net"
	"os"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/LuisFelipePoma/Movies_Recomender_With_Golang/src/backend/slave/model"
	"github.com/LuisFelipePoma/Movies_Recomender_With_Golang/src/backend/types"
	jsoniter "github.com/json-iterator/go"
)

var json = jsoniter.ConfigFastest        // Use jsoniter for faster performance
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
	fmt.Printf("Slave %s listening on port %s\n", name, port) // Show the port
	fmt.Println("Local address:", ln.Addr())                  // Show local addres

	// Accept incoming connections
	for {
		conn, err := ln.Accept()
		if err != nil {
			fmt.Println("Error accepting connection:", err)
			continue
		}
		go handleConnection(conn) // Handle the connection
	}
}

// Map of functions to execute based on the task type
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

	// Read the data from the connection
	reader := bufio.NewReader(conn)
	decoder := json.NewDecoder(reader)

	// Decode the JSON data
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
	query := strings.ToLower(data.TaskSearch.Query) // Convert query to lowercase
	movies := taskMaster.Data.Movies                // Get movies from task

	fmt.Println("Buscando peliculas...")
	// Inicializar la semilla de aleatoriedad
	rand.Seed(time.Now().UnixNano())

	// Define the result struct
	type result struct {
		movie      types.MovieResponse
		similarity float64
	}

	numWorkers := 5                             // Number of workers
	jobs := make(chan types.Movie, len(movies)) // Jobs channel
	results := make(chan result, len(movies))   // Results channel

	var wg sync.WaitGroup

	// Define the field weights
	fieldWeights := []struct {
		field  func(types.Movie) string
		weight float64
	}{
		{func(m types.Movie) string { return m.Title }, 4},
		{func(m types.Movie) string { return m.Genres }, 3},
		{func(m types.Movie) string { return m.Actors }, 2},
		{func(m types.Movie) string { return m.Keywords }, 1},
	}

	// Start workers to process
	for w := 0; w < numWorkers; w++ {
		wg.Add(1) // Increment the wait group
		go func() {
			defer wg.Done()
			for movie := range jobs { // Loop over the jobs
				totalSimilarity := 0.0            // Initialize the total similarity
				for _, fw := range fieldWeights { // Loop over the field weights
					if strings.Contains(strings.ToLower(fw.field(movie)), query) { // Check if the field contains the query
						totalSimilarity += fw.weight // Increment the total similarity
					}
				}
				// If the total similarity is greater than 0, send the result
				if totalSimilarity > 0 {
					results <- result{ // Send the result
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

	// Send jobs to the workers
	for _, movie := range movies {
		jobs <- movie
	}
	close(jobs)    //  Close the jobs channel
	wg.Wait()      // Wait for all workers to finish
	close(results) // Close the results channel

	// Collect the results from the channel
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

// getNMovies returns the random N movies with concurrency using Reservoir Sampling
func getNMovies(taskMaster types.TaskDistributed) []types.MovieResponse {
	data := taskMaster.Data         // Get data from task
	n := data.Quantity              // Get quantity
	movies := data.Movies           // Get movies
	genres := data.TaskSearch.Query // Get genres

	fmt.Println("Obteniendo", n, "películas aleatorias")

	totalVote := 0.0 // Initialize the total vote
	for _, movie := range movies {
		totalVote += movie.VoteAverage // Sum the vote average
	}
	// Calculate the average vote of the movies
	averageVote := totalVote / float64(len(movies))

	var (
		reservoir     []types.MovieResponse
		reservoirLock sync.Mutex
		count         int
	)

	numWorkers := 5                             // Number of workers
	jobs := make(chan types.Movie, len(movies)) // Jobs channel
	var wg sync.WaitGroup

	// Start workers
	for w := 0; w < numWorkers; w++ {
		wg.Add(1)
		go func(workerID int) {
			defer wg.Done()
			for movie := range jobs { // Loop over the jobs
				// Filtrar por VoteAverage
				if movie.VoteAverage <= averageVote {
					continue
				}
				// Filtrar por géneros si no contiene "All"
				if !strings.Contains(genres, "All") {
					genreList := strings.Split(genres, ",") // Split the genres
					match := true                           // Initialize the match
					for _, genre := range genreList {
						if !strings.Contains(movie.Genres, strings.TrimSpace(genre)) {
							match = false
							break
						}
					}
					if !match {
						continue
					}
				}

				// Reservoir Sampling
				reservoirLock.Lock()
				if count < n {
					// Reservorio no lleno, agregar directamente
					reservoir = append(reservoir, types.MovieResponse{
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
					count++
				} else {
					// Reservorio lleno, reemplazar con probabilidad
					j := rand.Intn(count + 1)
					if j < n {
						reservoir[j] = types.MovieResponse{
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
						}
					}
					count++
				}
				reservoirLock.Unlock()
			}
		}(w)
	}

	// Enviar trabajos
	go func() {
		for _, movie := range movies {
			jobs <- movie
		}
		close(jobs)
	}()

	// Esperar a que todos los workers terminen
	wg.Wait()

	return reservoir
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
