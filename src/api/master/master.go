package master

import (
	// "encoding/json"
	// "fmt"
	// "io"
	"log"
	// t "master/types"
	"net"
	// // "sort"
	// "sync"
	// "time"
)

// var slaveNodes = []string{
// 	"localhost:8082",
// 	"localhost:8083",
// 	"localhost:8084",
// }

type NodeMaster struct {
	Port string
}

func NewNodeMaster(port string) *NodeMaster {
	return &NodeMaster{
		Port: port,
	}
}

func (n *NodeMaster) Run() {
	ln, err := net.Listen("tcp", n.Port)
	if err != nil {
		log.Fatalf("Error al iniciar el servidor: %v\n", err)
	}
	defer ln.Close()

	log.Printf("Servidor iniciado en el puerto %s\n", n.Port)

	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Printf("Error al aceptar la conexión: %v\n", err)
			continue
		}
		conn.Close()
		// go handleConnection(conn)
	}
}

// func similarMoviesHandler(movieID int) []t.MovieResponse {
// 	start := time.Now()
// 	// Distribute the task to the slave nodes
// 	numSlaves := len(slaveNodes)
// 	ranges := splitRanges(len(moviesService.Movies), numSlaves)
// 	// Create a goroutine for each slave node
// 	var wg sync.WaitGroup
// 	// Channel to receive the results from the slaves
// 	results := make(chan []t.SimilarMovie, numSlaves)

// 	for i, node := range slaveNodes {
// 		wg.Add(1)
// 		go func(node string, startIdx, endIdx, movieID int) {
// 			defer wg.Done()
// 			result, err := getSimilarMoviesFromNode(node, startIdx, endIdx, movieID)
// 			if err == nil {
// 				results <- result
// 			}
// 		}(node, ranges[i][0], ranges[i][1], movieID)
// 	}
// 	// Wait for all goroutines to finish
// 	wg.Wait()
// 	close(results)

// 	// Combine the results from all the slaves
// 	var combinedResults []t.SimilarMovie
// 	for result := range results {
// 		combinedResults = append(combinedResults, result...)
// 	}

// 	// Sort the combined results by similarity
// 	sort.Slice(combinedResults, func(i, j int) bool {
// 		return combinedResults[i].Similarity > combinedResults[j].Similarity
// 	})

// 	// Limit the number of results to 10
// 	if len(combinedResults) > 10 {
// 		combinedResults = combinedResults[:10]
// 	}

// 	// Map similar movie IDs to movie details
// 	var movieResponses []t.MovieResponse
// 	for _, similarMovie := range combinedResults {
// 		for _, movie := range moviesService.Movies {
// 			if movie.ID == similarMovie.ID {
// 				movieResponses = append(movieResponses, t.MovieResponse{
// 					ID:          similarMovie.ID,
// 					Title:       movie.Title,
// 					Characters:  movie.Characters,
// 					Actors:      movie.Actors,
// 					Director:    movie.Director,
// 					Genres:      movie.Genres,
// 					ImdbId:      movie.ImdbId,
// 					VoteAverage: movie.VoteAverage,
// 				})
// 				break
// 			}
// 		}
// 	}
// 	fmt.Printf("Request processed in %s\n", time.Since(start))
// 	return movieResponses
// }

// func splitRanges(totalMovies, numParts int) [][2]int {
// 	chunkSize := (totalMovies + numParts - 1) / numParts
// 	var ranges [][2]int
// 	for i := 0; i < totalMovies; i += chunkSize {
// 		end := i + chunkSize
// 		if end > totalMovies {
// 			end = totalMovies
// 		}
// 		ranges = append(ranges, [2]int{i, end})
// 	}
// 	return ranges
// }

// func getSimilarMoviesFromNode(node string, startIdx, endIdx, movieID int) ([]t.SimilarMovie, error) {
// 	conn, err := net.DialTimeout("tcp", node, 15*time.Second)
// 	if err != nil {
// 		log.Printf("Error al conectar con el nodo %s: %v\n", node, err)
// 		return reassignTask(struct {
// 			Movies      []t.Movie `json:"movies"`
// 			TargetMovie t.Movie   `json:"target_movie"`
// 		}{
// 			Movies:      moviesService.Movies[startIdx:endIdx],
// 			TargetMovie: *moviesService.GetMovieByID(movieID),
// 		}, node)
// 	}
// 	defer conn.Close()

// 	task := struct {
// 		Movies      []t.Movie `json:"movies"`
// 		TargetMovie t.Movie   `json:"target_movie"`
// 	}{
// 		Movies:      moviesService.Movies[startIdx:endIdx],
// 		TargetMovie: *moviesService.GetMovieByID(movieID),
// 	}

// 	data, err := json.Marshal(task)
// 	if err != nil {
// 		return nil, err
// 	}

// 	_, err = conn.Write(data)
// 	if err != nil {
// 		log.Printf("Error al enviar la tarea al nodo %s: %v\n", node, err)
// 		return reassignTask(task, node)
// 	}

// 	response, err := io.ReadAll(conn)
// 	if err != nil {
// 		log.Printf("Error al leer la respuesta del nodo %s: %v\n", node, err)
// 		return reassignTask(task, node)
// 	}

// 	var similarMovies []t.SimilarMovie
// 	err = json.Unmarshal(response, &similarMovies)
// 	if err != nil {
// 		return nil, err
// 	}

// 	return similarMovies, nil
// }

// func reassignTask(task interface{}, failedNode string) ([]t.SimilarMovie, error) {
// 	for _, node := range slaveNodes {
// 		if node != failedNode {
// 			result, err := sendTaskToNode(node, task)
// 			if err == nil {
// 				return result, nil
// 			}
// 			log.Printf("Error al intentar reasignar al nodo %s: %v\n", node, err)
// 		}
// 	}
// 	return nil, fmt.Errorf("no hay nodos disponibles para reasignar la tarea")
// }

// func sendTaskToNode(node string, task interface{}) ([]t.SimilarMovie, error) {
// 	conn, err := net.DialTimeout("tcp", node, 3*time.Second)
// 	if err != nil {
// 		return nil, fmt.Errorf("error al conectar con el nodo %s: %v", node, err)
// 	}
// 	defer conn.Close()

// 	data, err := json.Marshal(task)
// 	if err != nil {
// 		return nil, err
// 	}

// 	_, err = conn.Write(data)
// 	if err != nil {
// 		return nil, fmt.Errorf("error al enviar la tarea al nodo %s: %v", node, err)
// 	}

// 	response, err := io.ReadAll(conn)
// 	if err != nil {
// 		return nil, fmt.Errorf("error al leer la respuesta del nodo %s: %v", node, err)
// 	}

// 	var similarMovies []t.SimilarMovie
// 	err = json.Unmarshal(response, &similarMovies)
// 	if err != nil {
// 		return nil, err
// 	}

// 	return similarMovies, nil
// }
