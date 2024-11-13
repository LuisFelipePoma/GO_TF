package main

import (
	"encoding/json"
	"fmt"
	Error "github.com/LuisFelipePoma/Movies_Recomender_With_Golang/src/backend/errors"
	"github.com/LuisFelipePoma/Movies_Recomender_With_Golang/src/backend/services"
	"github.com/LuisFelipePoma/Movies_Recomender_With_Golang/src/backend/types"
	"github.com/LuisFelipePoma/Movies_Recomender_With_Golang/src/backend/utils"
	"io"
	"log"
	"net"
	"os"
	"sort"
	"sync"
	"time"
)

var slaveNodes = []string{
	os.Getenv("SLAVE1"),
	os.Getenv("SLAVE2"),
	os.Getenv("SLAVE3"),
}

// Task represents a task to be processed by a slave node
type Task struct {
	Slave    string
	TaskData types.TaskDistributed
}

// Result represents the result from a slave node
type Result struct {
	Response types.Response
	Error    error
}

var moviesService = services.NewMovies()
var movies []types.Movie

// Distribute the task to the slave nodes
var numSlaves = len(slaveNodes)
var rangesMovies [][2]int

const TIMEOUT = 10 * time.Second
const MAX_RETRIES = 2

// 500ms
const RETRY_DELAY = 150 * time.Millisecond

// ENTRYPOINT
func main() {
	// Leer el puerto desde la variable de entorno
	port := os.Getenv("PORT")
	name := os.Getenv("NODE_NAME")
	if port == "" {
		log.Fatal("El puerto no está configurado en la variable de entorno PORT")
	}

	// Cargar Peliculas
	err := moviesService.LoadMovies("movies.json")
	if err != nil {
		fmt.Println("No se pudo cargar la informacion de las peliculas.", err)
		return
	}

	// Cargar Ratings
	err = moviesService.LoadRatings("ratings.csv")
	if err != nil {
		fmt.Println("No se pudo cargar los ratings.", err)
		return
	}

	// Split the movies into ranges
	movies = moviesService.Movies
	rangesMovies = utils.SplitRanges(len(movies), numSlaves)

	// Create a listener
	listener, err := net.Listen("tcp", ":"+port)
	if err != nil {
		log.Fatalf("Error al crear el servidor: %v", err)
	}
	defer listener.Close()

	// Listening
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

var dictFunction = map[types.TaskType]func(types.TaskDistributed) types.Response{
	types.TaskRecomend:     similarMoviesHandler,
	types.TaskSearch:       searchMoviesHandler,
	types.TaskGet:          getNMoviesHandler,
	types.TaskUserRecomend: sendUserRecomendations,
}

// Handle the incoming requests
func handleRequest(conn net.Conn) {
	defer conn.Close()
	// Decode the request
	var task types.TaskDistributed
	decoder := json.NewDecoder(conn)
	if err := decoder.Decode(&task); err != nil {
		fmt.Println("Error al decodificar JSON:", err)
		Error.ReturnError(conn, "Error al decodificar JSON")
		return
	}
	fmt.Printf("Se realizara la tarea: %v\n", task.Type)
	// Start the timer
	start := time.Now()
	response := dictFunction[task.Type](task)
	fmt.Printf("Tarea procesada en %s\n", time.Since(start))
	// Send the response
	if err := Error.SendJSONResponse(conn, response); err != nil {
		Error.ReturnError(conn, err.Error())
		return
	}
	fmt.Println("Nodo Master envió resultado")
}

// <----------- HANDLERS TASK FOR THE NODES
// SimilarMoviesHandler returns a list of similar movies based
func similarMoviesHandler(task types.TaskDistributed) types.Response {
	// get data from task
	data := task.Data.TaskRecomendations
	fmt.Printf("Recomendacion para %+v\n", data.MovieId)

	// Get the data from the task
	n_recomendations := task.Data.Quantity
	movieTarget, err := moviesService.GetMovieById(data.MovieId)
	if err != nil {
		return types.Response{Error: fmt.Sprintf("Error getting movie by ID: %v", err)}
	}

	// update the task to the new ranges for each node
	var tasks []types.TaskDistributed
	for _, r := range rangesMovies {
		newTask := types.TaskDistributed{
			Type: types.TaskRecomend,
			Data: types.TaskData{
				Movies: movies[r[0]:r[1]],
				TaskRecomendations: &types.TaskRecomendations{
					MovieId:     movieTarget.ID,
					TargetMovie: *movieTarget,
				},
			},
		}
		tasks = append(tasks, newTask)
	}

	// SEND THE TASK TO THE NODES
	results := sendTasksToNodes(tasks)

	// Sort the combined results by similarity
	sort.Slice(results, func(i, j int) bool {
		return results[i].Similarity > results[j].Similarity
	})

	// Limit the number of results to n_recomendations
	if len(results) > n_recomendations {
		results = results[:n_recomendations]
	}

	// Create the response
	response := types.Response{
		Error:         "",
		MovieResponse: results,
		TargetMovie:   movieTarget.Title,
	}

	return response
}

// SearchMoviesHandler returns a list of movies based on a search query
func searchMoviesHandler(task types.TaskDistributed) types.Response {
	data := task.Data.TaskSearch
	query := data.Query

	fmt.Println("Buscando peliculas...")
	// update the task to the new ranges for each node
	var tasks []types.TaskDistributed
	for _, r := range rangesMovies {
		newTask := types.TaskDistributed{
			Type: types.TaskSearch,
			Data: types.TaskData{
				Movies: movies[r[0]:r[1]],
				TaskSearch: &types.TaskSearchQuery{
					Query: query,
				},
			},
		}
		tasks = append(tasks, newTask)
	}

	// Send the task to the nodes
	results := sendTasksToNodes(tasks) // get the combined results

	// Sort based on the similarity (importance) and voteAverage
	sort.Slice(results, func(i, j int) bool {
		if results[i].Similarity == results[j].Similarity {
			return results[i].VoteAverage > results[j].VoteAverage
		}
		return results[i].Similarity > results[j].Similarity
	})

	// Create the response
	response := types.Response{
		Error:         "",
		MovieResponse: results,
		TargetMovie:   query,
	}

	return response
}

// GetNMoviesHandler returns a list of n movies
func getNMoviesHandler(task types.TaskDistributed) types.Response {
	// Get the data from the task
	n := task.Data.Quantity

	// Update the task to the new ranges for each node
	var tasks []types.TaskDistributed
	for _, r := range rangesMovies {
		newTask := types.TaskDistributed{
			Type: types.TaskGet,
			Data: types.TaskData{
				Movies:   movies[r[0]:r[1]],
				Quantity: n,
			},
		}
		tasks = append(tasks, newTask)
	}

	// Send the task to the nodes
	results := sendTasksToNodes(tasks)

	// Combine the results from all the slaves
	var combinedResults []types.MovieResponse
	combinedResults = append(combinedResults, results...)

	// Limit the number of results to n
	if len(combinedResults) > n {
		combinedResults = combinedResults[:n]
	}

	// Create the response
	response := types.Response{
		Error:         "",
		MovieResponse: combinedResults,
	}

	return response
}

// DivideUsers splits the users into specified number of groups, excluding the given userId
func SplitUsers(usersRatings map[int]types.User, excludeUserId int, numGroups int) []map[int]types.User {
	groups := make([]map[int]types.User, numGroups)
	for i := range groups {
		groups[i] = make(map[int]types.User)
	}
	currentGroup := 0
	for id, user := range usersRatings {
		if user.ID != excludeUserId {
			groups[currentGroup][id] = user
			currentGroup = (currentGroup + 1) % numGroups
		}
	}
	return groups
}

// sendUserRecomendations returns a list of movies recommended for a user
func sendUserRecomendations(task types.TaskDistributed) types.Response {
	// Get the data from the task
	userID := moviesService.GetRandomUserID()
	fmt.Printf("Recomendaciones para el usuario %d\n", userID)

	// Divide users into groups excluding the current user
	userGroups := SplitUsers(moviesService.UserRatings, userID, numSlaves)

	// Update the task
	var tasks []types.TaskDistributed
	for _, group := range userGroups {
		newTask := types.TaskDistributed{
			Type: types.TaskUserRecomend,
			Data: types.TaskData{
				Quantity: 5,
				TaskUserRecomendations: &types.TaskUserRecomendations{
					UserID:      userID,
					User:        moviesService.UserRatings[userID].Ratings,
					UserRatings: group, // Ensure this method processes the group
				},
			},
		}
		tasks = append(tasks, newTask)
	}

	// Send the task to the nodes
	results := sendTasksToNodes(tasks)

	// Combine the results from all the slaves
	var combinedResults []types.MovieResponse
	combinedResults = append(combinedResults, results...)

	// Sort the combined results by similarity
	sort.Slice(combinedResults, func(i, j int) bool {
		return combinedResults[i].Similarity > combinedResults[j].Similarity
	})

	// Create the response
	response := types.Response{
		Error:         "",
		MovieResponse: combinedResults,
		TargetMovie:   fmt.Sprintf("%d", userID),
	}
	return response
}

// <------------ Function to handle the connection with the nodes
func sendTasksToNodes(tasks []types.TaskDistributed) []types.MovieResponse {
	// Create a goroutine for each slave node
	var wg sync.WaitGroup
	// Channel to receive the results from the slaves
	results := make(chan []types.MovieResponse, numSlaves)

	// Create a goroutine for each slave node
	for i, node := range slaveNodes {
		wg.Add(1)
		go func(node string, tasks []types.TaskDistributed) {
			defer wg.Done()
			result, err := senTaskToNode(node, tasks[i])
			if err == nil {
				results <- result
			} else {
				fmt.Println(err)
			}
		}(node, tasks)
	}
	// Wait for all goroutines to finish
	wg.Wait()
	close(results)

	// Combine the results from all the slaves
	var combinedResults []types.MovieResponse
	for result := range results {
		combinedResults = append(combinedResults, result...)
	}
	return combinedResults
}

// Funtion to Get movies from the nodes
func senTaskToNode(node string, task types.TaskDistributed) ([]types.MovieResponse, error) {
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
		return reassignTask(task, node)
	}
	defer conn.Close()

	// If the connection was successful, send the task to the node
	fmt.Printf("-- Enviando tarea al nodo %s --\n", node)

	// Send the task to the node
	data, err := json.Marshal(task)
	if err != nil {
		return nil, err
	}

	// Set time out for the connection of 35s
	conn.SetDeadline(time.Now().Add(20 * time.Second))

	// Send the data to the node
	_, err = conn.Write(data)
	if err != nil {
		log.Printf("Error al enviar la tarea al nodo %s: %v\n", node, err)
		return reassignTask(task, node)
	}

	// Read the response from the node
	response, err := io.ReadAll(conn)
	if err != nil {
		log.Printf("Error al leer la respuesta del nodo %s: %v\n", node, err)
		return reassignTask(task, node)
	}

	// Parse the response
	var similarMovies []types.MovieResponse
	err = json.Unmarshal(response, &similarMovies)
	if err != nil {
		return nil, err
	}
	return similarMovies, nil
}

// Reassign the task to another node
func reassignTask(task interface{}, failedNode string) ([]types.MovieResponse, error) {
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
func sendTaskToNode(node string, task interface{}) ([]types.MovieResponse, error) {
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
	var similarMovies []types.MovieResponse
	err = json.Unmarshal(response, &similarMovies)
	if err != nil {
		return nil, err
	}
	return similarMovies, nil
}
