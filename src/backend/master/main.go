package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"sort"
	"sync"
	"time"

	Error "github.com/LuisFelipePoma/Movies_Recomender_With_Golang/src/backend/errors"
	"github.com/LuisFelipePoma/Movies_Recomender_With_Golang/src/backend/services"
	"github.com/LuisFelipePoma/Movies_Recomender_With_Golang/src/backend/types"
	"github.com/LuisFelipePoma/Movies_Recomender_With_Golang/src/backend/utils"
)

var slaveNodes = []string{
	"slave1:8082",
	"slave2:8083",
	"slave3:8084",
}

var moviesService = services.NewMovies()

const TIMEOUT = 5 * time.Second
const MAX_RETRIES = 3

// 500ms
const RETRY_DELAY = 500 * time.Millisecond

// ENTRYPOINT
func main() {
	// Leer el puerto desde la variable de entorno
	port := os.Getenv("PORT")
	name := os.Getenv("NODE_NAME")
	if port == "" {
		log.Fatal("El puerto no está configurado en la variable de entorno PORT")
	}
	// Create a listener
	listener, err := net.Listen("tcp", ":"+port)
	if err != nil {
		log.Fatalf("Error al crear el servidor: %v", err)
	}
	defer listener.Close()

	// Load the movies from JSON
	err = moviesService.LoadMovies("./movies.json")
	if err != nil {
		log.Fatalf("Error al cargar las películas: %v", err)
	}
	// LIstening
	fmt.Printf("Nodo %s escuchando en el puerto %s\n", name, port)
	// Start the server
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Printf("Error al aceptar la conexión: %v\n", err)
			continue
		}

		go handleRequest(conn)
	}
}

// Handle the incoming requests
func handleRequest(conn net.Conn) {
	defer conn.Close()
	// Decode the request
	var task types.Request
	decoder := json.NewDecoder(conn)
	if err := decoder.Decode(&task); err != nil {
		fmt.Println("Error al decodificar JSON:", err)
		Error.ReturnError(conn, "Error al decodificar JSON")
		return
	}
	fmt.Println("Leyendo tarea....")
	fmt.Printf("%+v\n", task)

	// Process based ont the task OPTION
	switch task.Option {
	// CASE 1: Get similar movies based on one movie
	case 1:
		fmt.Println("Opcion 1")
		movie := moviesService.GetMovieByTitle(task.Data) // Get the movie
		if movie == nil {
			fmt.Println("Película no encontrada")
			Error.ReturnError(conn, "Película no encontrada")
			return
		}
		// Process the task
		response := similarMoviesHandler(movie.ID)
		// Send the response
		if err := Error.SendJSONResponse(conn, response); err != nil {
			Error.ReturnError(conn, err.Error())
			return
		}
		fmt.Println("Nodo Master envió resultado")
	// CASE 2: Get the last recommendations made based on the last movie
	case 2:
		fmt.Println("Opcion 2")
		recomendations := moviesService.Recommendations // Get the last recommendations
		if recomendations == nil {
			fmt.Println("No hay recomendaciones disponibles")
			Error.ReturnError(conn, "No hay recomendaciones disponibles")
			return
		}
		// Create the response
		response := types.Response{
			Error:         "",
			MovieResponse: recomendations,
			TargetMovie:   moviesService.LastRecomendation,
		}
		// Send the response
		if err := Error.SendJSONResponse(conn, response); err != nil {
			Error.ReturnError(conn, err.Error())
			return
		}
		fmt.Println("Nodo Master envió resultado")

	// CASE 3: Filter the recomendations by genre
	case 3:
		fmt.Println("Opcion 3")
		genre := task.Data                                      // Get the genre from the task data
		movies := moviesService.GetRecomendationsByGenre(genre) // Get the movies
		if len(movies) == 0 {
			fmt.Println("No se encontraron películas con el género especificado")
			Error.ReturnError(conn, "No se encontraron películas con el género especificado")
			return
		}
		// Create the response
		response := types.Response{
			Error:         "",
			MovieResponse: movies,
			TargetMovie:   moviesService.LastRecomendation,
		}
		// Send the response
		if err := Error.SendJSONResponse(conn, response); err != nil {
			Error.ReturnError(conn, err.Error())
			return
		}
		fmt.Println("Nodo Master envió resultado")
	// CASE 4: Filter the recomendations by vote average
	case 4:
		fmt.Println("Opcion 4")
		voteAverage := task.Data
		movies := moviesService.GetMoviesByVoteAverage(voteAverage)
		if len(movies) == 0 {
			fmt.Println("No se encontraron películas con el voteAverage especificado")
			Error.ReturnError(conn, "No se encontraron películas con el voteAverage especificado")
			return
		}
		response := types.Response{
			Error:         "",
			MovieResponse: movies,
			TargetMovie:   moviesService.LastRecomendation,
		}
		if err := Error.SendJSONResponse(conn, response); err != nil {
			Error.ReturnError(conn, err.Error())
			return
		}
		fmt.Println("Nodo Master envió resultado")

	}
}

// SimilarMoviesHandler returns a list of similar movies based
func similarMoviesHandler(movieID int) types.Response {
	start := time.Now() // Start the timer
	// Distribute the task to the slave nodes
	numSlaves := len(slaveNodes)
	ranges := utils.SplitRanges(len(moviesService.Movies), numSlaves)
	// Create a goroutine for each slave node
	var wg sync.WaitGroup
	// Channel to receive the results from the slaves
	results := make(chan []types.SimilarMovie, numSlaves)

	for i, node := range slaveNodes {
		wg.Add(1)
		go func(node string, startIdx, endIdx, movieID int) {
			defer wg.Done()
			result, err := getSimilarMoviesFromNode(node, startIdx, endIdx, movieID)
			if err == nil {
				results <- result
			} else {
				fmt.Println(err)
			}
		}(node, ranges[i][0], ranges[i][1], movieID)
	}
	// Wait for all goroutines to finish
	wg.Wait()
	close(results)

	// Combine the results from all the slaves
	var combinedResults []types.SimilarMovie
	for result := range results {
		combinedResults = append(combinedResults, result...)
	}

	// Sort the combined results by similarity
	sort.Slice(combinedResults, func(i, j int) bool {
		return combinedResults[i].Similarity > combinedResults[j].Similarity
	})

	// Limit the number of results to 10
	if len(combinedResults) > 10 {
		combinedResults = combinedResults[:10]
	}

	// Map similar movie IDs to movie details
	var movieResponses []types.MovieResponse
	for _, similarMovie := range combinedResults {
		for _, movie := range moviesService.Movies {
			if movie.ID == similarMovie.ID {
				movieResponses = append(movieResponses, types.MovieResponse{
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
	// Create the response
	response := types.Response{
		Error:         "",
		MovieResponse: movieResponses,
		TargetMovie:   moviesService.GetMovieByID(movieID).Title,
	}
	moviesService.Recommendations = movieResponses // Save the recommendations
	fmt.Printf("Request processed in %s\n", time.Since(start))
	return response
}

// Get the similar movies from a slave node
func getSimilarMoviesFromNode(node string, startIdx, endIdx, movieID int) ([]types.SimilarMovie, error) {
	// Connect to the slave node
	var conn net.Conn
	var err error

	// Try to connect to the node with retries
	for i := 0; i < MAX_RETRIES; i++ {
		conn, err = net.DialTimeout("tcp", node, TIMEOUT)
		if err == nil {
			break
		}
		log.Printf("Error al conectar con el nodo %s, reintentando (%d/%d)\n", node, i+1, MAX_RETRIES)
		time.Sleep(RETRY_DELAY)
	}

	// Check if there was an error connecting to the node
	if err != nil {
		log.Printf("Error al conectar con el nodo %s\n", node)
		// Reassign the task to another node
		return reassignTask(struct {
			Movies      []types.Movie `json:"movies"`
			TargetMovie types.Movie   `json:"target_movie"`
		}{
			Movies:      moviesService.Movies[startIdx:endIdx],
			TargetMovie: *moviesService.GetMovieByID(movieID),
		}, node)
	}
	defer conn.Close()

	// If the connection was successful, send the task to the node
	fmt.Printf("-- Enviando tarea al nodo %s --\n", node)

	// Create the task
	task := struct {
		Movies      []types.Movie `json:"movies"`
		TargetMovie types.Movie   `json:"target_movie"`
	}{
		Movies:      moviesService.Movies[startIdx:endIdx],
		TargetMovie: *moviesService.GetMovieByID(movieID),
	}

	// Send the task to the node
	data, err := json.Marshal(task)
	if err != nil {
		return nil, err
	}

	// Send the data to the node
	_, err = conn.Write(data)
	if err != nil {
		log.Printf("Error al enviar la tarea al nodo %s: %v\n", node, err)
		return reassignTask(task, node)
	}

	// Stablish a read deadline
	conn.SetReadDeadline(time.Now().Add(5 * time.Second))

	// Read the response from the node
	response, err := io.ReadAll(conn)
	if err != nil {
		log.Printf("Error al leer la respuesta del nodo %s: %v\n", node, err)
		return reassignTask(task, node)
	}

	// Parse the response
	var similarMovies []types.SimilarMovie
	err = json.Unmarshal(response, &similarMovies)
	if err != nil {
		return nil, err
	}

	// Save the last recommendation
	moviesService.LastRecomendation = moviesService.GetMovieByID(movieID).Title
	return similarMovies, nil
}

// Reassign the task to another node
func reassignTask(task interface{}, failedNode string) ([]types.SimilarMovie, error) {
	// Try to reassign the task to another node
	for _, node := range slaveNodes {
		// Check if the node is not the one that failed
		if node != failedNode {
			result, err := sendTaskToNode(node, task)
			if err == nil {
				fmt.Printf("-- Reasignada la tarea del nodo %s al nodo %s --\n", failedNode, node)
				return result, nil
			}
			log.Printf("Error al intentar la tarea del nodo %s reasignar al nodo %s\n", failedNode, node)
		}
	}
	return nil, fmt.Errorf("<-- No hay nodos disponibles para reasignar la tarea del nodo %s -->", failedNode)
}

// Send the task to a node
func sendTaskToNode(node string, task interface{}) ([]types.SimilarMovie, error) {
	conn, err := net.DialTimeout("tcp", node, TIMEOUT)
	if err != nil {
		return nil, fmt.Errorf("error al conectar con el nodo %s: %v", node, err)
	}
	defer conn.Close()

	// Create the task
	data, err := json.Marshal(task)
	if err != nil {
		return nil, err
	}
	// Send the data to the node
	_, err = conn.Write(data)
	if err != nil {
		return nil, fmt.Errorf("error al enviar la tarea al nodo %s: %v", node, err)
	}
	// Read the response
	response, err := io.ReadAll(conn)
	if err != nil {
		return nil, fmt.Errorf("error al leer la respuesta del nodo %s: %v", node, err)
	}

	// Parse the response
	var similarMovies []types.SimilarMovie
	err = json.Unmarshal(response, &similarMovies)
	if err != nil {
		return nil, err
	}
	return similarMovies, nil
}
