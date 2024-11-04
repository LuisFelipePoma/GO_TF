package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net"
	"sort"
	"sync"
	"time"

	Error "github.com/LuisFelipePoma/Movies_Recomender_With_Golang/src/api/errors"
	"github.com/LuisFelipePoma/Movies_Recomender_With_Golang/src/api/services"
	"github.com/LuisFelipePoma/Movies_Recomender_With_Golang/src/api/types"
	"github.com/LuisFelipePoma/Movies_Recomender_With_Golang/src/api/utils"
)

var slaveNodes = []string{
	"localhost:8082",
	"localhost:8083",
	"localhost:8084",
}

var moviesService = services.NewMovies()

const TIMEOUT = 5 * time.Second

func main() {
	// Create a listen
	listener, err := net.Listen("tcp", "localhost:8081")
	if err != nil {
		log.Fatalf("Error al crear el servidor: %v", err)
	}
	defer listener.Close()

	// Load the movies
	err = moviesService.LoadMovies("../database/data/data_clean.json")
	if err != nil {
		log.Fatalf("Error al cargar las películas: %v", err)
	}

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

func handleRequest(conn net.Conn) {
	defer conn.Close()

	// Decodificar el JSON recibido en una estructura Task
	var task types.Request
	decoder := json.NewDecoder(conn)
	if err := decoder.Decode(&task); err != nil {
		fmt.Println("Error al decodificar JSON:", err)
		Error.ReturnError(conn, "Error al decodificar JSON")
		return
	}
	fmt.Println("Leyendo tarea....")
	fmt.Printf("%+v\n", task)

	// Procesar la solicitud
	switch task.Option {
	case 1:
		fmt.Println("Option 1")
		movie := moviesService.GetMovieByTitle(task.Data)
		if movie == nil {
			fmt.Println("Película no encontrada")
			Error.ReturnError(conn, "Película no encontrada")
			return
		}
		response := similarMoviesHandler(movie.ID)
		if err := Error.SendJSONResponse(conn, response); err != nil {
			Error.ReturnError(conn, err.Error())
			return
		}
		fmt.Println("Nodo Master envió resultado")
	case 2:
		fmt.Println("Option 2")
		recomendations := moviesService.Recommendations
		if recomendations == nil {
			fmt.Println("No hay recomendaciones disponibles")
			Error.ReturnError(conn, "No hay recomendaciones disponibles")
			return
		}
		response := types.Response{
			Error:         "",
			MovieResponse: recomendations,
			TargetMovie:   moviesService.LastRecomendation,
		}
		if err := Error.SendJSONResponse(conn, response); err != nil {
			Error.ReturnError(conn, err.Error())
			return
		}
		fmt.Println("Nodo Master envió resultado")
	case 3:
		fmt.Println("Option 3")
		genre := task.Data
		movies := moviesService.GetRecomendationsByGenre(genre)
		if len(movies) == 0 {
			fmt.Println("No se encontraron películas con el género especificado")
			Error.ReturnError(conn, "No se encontraron películas con el género especificado")
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
	case 4:
		fmt.Println("Option 4")
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

func similarMoviesHandler(movieID int) types.Response {
	start := time.Now()
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
	response := types.Response{
		Error:         "",
		MovieResponse: movieResponses,
		TargetMovie:   moviesService.GetMovieByID(movieID).Title,
	}
	moviesService.Recommendations = movieResponses
	fmt.Printf("Request processed in %s\n", time.Since(start))
	return response
}

func getSimilarMoviesFromNode(node string, startIdx, endIdx, movieID int) ([]types.SimilarMovie, error) {
	const maxRetries = 3
	const retryDelay = 2 * time.Second

	var conn net.Conn
	var err error

	for i := 0; i < maxRetries; i++ {
		conn, err = net.DialTimeout("tcp", node, TIMEOUT)
		if err == nil {
			break
		}
		log.Printf("Error al conectar con el nodo %s, reintentando (%d/%d)\n", node, i+1, maxRetries)
		time.Sleep(retryDelay)
	}

	if err != nil {
		log.Printf("Error al conectar con el nodo %s\n", node)
		return reassignTask(struct {
			Movies      []types.Movie `json:"movies"`
			TargetMovie types.Movie   `json:"target_movie"`
		}{
			Movies:      moviesService.Movies[startIdx:endIdx],
			TargetMovie: *moviesService.GetMovieByID(movieID),
		}, node)
	}
	defer conn.Close()

	fmt.Printf("-- Enviando tarea al nodo %s --\n", node)
	task := struct {
		Movies      []types.Movie `json:"movies"`
		TargetMovie types.Movie   `json:"target_movie"`
	}{
		Movies:      moviesService.Movies[startIdx:endIdx],
		TargetMovie: *moviesService.GetMovieByID(movieID),
	}

	data, err := json.Marshal(task)
	if err != nil {
		return nil, err
	}

	_, err = conn.Write(data)
	if err != nil {
		log.Printf("Error al enviar la tarea al nodo %s: %v\n", node, err)
		return reassignTask(task, node)
	}

	// Establecer un límite de tiempo para la lectura de la respuesta
	conn.SetReadDeadline(time.Now().Add(5 * time.Second))

	response, err := io.ReadAll(conn)
	if err != nil {
		log.Printf("Error al leer la respuesta del nodo %s: %v\n", node, err)
		return reassignTask(task, node)
	}

	var similarMovies []types.SimilarMovie
	err = json.Unmarshal(response, &similarMovies)
	if err != nil {
		return nil, err
	}
	moviesService.LastRecomendation = moviesService.GetMovieByID(movieID).Title
	return similarMovies, nil
}

func reassignTask(task interface{}, failedNode string) ([]types.SimilarMovie, error) {
	for _, node := range slaveNodes {
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

func sendTaskToNode(node string, task interface{}) ([]types.SimilarMovie, error) {
	conn, err := net.DialTimeout("tcp", node, 5*time.Second)
	if err != nil {
		return nil, fmt.Errorf("error al conectar con el nodo %s: %v", node, err)
	}
	defer conn.Close()

	data, err := json.Marshal(task)
	if err != nil {
		return nil, err
	}

	_, err = conn.Write(data)
	if err != nil {
		return nil, fmt.Errorf("error al enviar la tarea al nodo %s: %v", node, err)
	}

	response, err := io.ReadAll(conn)
	if err != nil {
		return nil, fmt.Errorf("error al leer la respuesta del nodo %s: %v", node, err)
	}

	var similarMovies []types.SimilarMovie
	err = json.Unmarshal(response, &similarMovies)
	if err != nil {
		return nil, err
	}

	return similarMovies, nil
}
