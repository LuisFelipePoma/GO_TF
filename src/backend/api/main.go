package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"
	"os"
	"strconv"

	"github.com/LuisFelipePoma/Movies_Recomender_With_Golang/src/backend/services"
	"github.com/LuisFelipePoma/Movies_Recomender_With_Golang/src/backend/types"
)

var nodeMasterPort = os.Getenv("MASTER_NODE") // Port of the master node
var moviesService = services.NewMovies()

func main() {
	// Read
	port := os.Getenv("PORT")
	// Cargar Peliculas
	err := moviesService.LoadMovies("movies.json")
	if err != nil {
		fmt.Println(err)
		return
	}

	// Configurar los manejadores HTTP
	setupRoutes()

	fmt.Println("Server is running on port ", port)
	http.ListenAndServe(":"+port, nil)

}

// Configurar los manejadores HTTP
func setupRoutes() {
	http.HandleFunc("/api/movies/similar", getSimilarMovies) // GET
	http.HandleFunc("/api/movies/id", getById)               // GET
	http.HandleFunc("/api/movies/search", getMoviesBySearch) // GET
	http.HandleFunc("/api/movies", getAllMovies)             // GET
}

// HANDLERS
func getSimilarMovies(w http.ResponseWriter, r *http.Request) {
	fmt.Println("GET /api/movies/similar")
	// Get title from args

	title := r.URL.Query().Get("title")
	response, errorMessage := handleMasterConection(title)
	if errorMessage != "" {
		fmt.Println(errorMessage)
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func getAllMovies(w http.ResponseWriter, r *http.Request) {
	fmt.Println("GET /api/movies")
	nStr := r.URL.Query().Get("n")
	if nStr == "" {
		nStr = "10"
	}
	// Convert n to int
	n, err := strconv.Atoi(nStr)

	if err != nil {
		http.Error(w, "Invalid number format", http.StatusBadRequest)
		return
	}
	response := moviesService.GetAllMovies(n)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func getById(w http.ResponseWriter, r *http.Request) {
	fmt.Println("GET /api/movies")
	idStr := r.URL.Query().Get("id")
	if idStr == "" {
		idStr = "1"
	}
	// Convert n to int
	id, err := strconv.Atoi(idStr)

	if err != nil {
		http.Error(w, "Invalid number format", http.StatusBadRequest)
		return
	}
	response := moviesService.GetMovieByID(id)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func getMoviesBySearch(w http.ResponseWriter, r *http.Request) {
	fmt.Println("GET /api/movies/search")
	query := r.URL.Query().Get("query")
	response := moviesService.GetMoviesBySearch(query)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// Conect to the master node
func handleMasterConection(data string) (types.Response, string) {
	conn, err := net.Dial("tcp", "master:"+nodeMasterPort) // Connect to the master node
	if err != nil {
		return types.Response{}, "Error al conectar con el nodo maestro."
	}
	defer conn.Close()

	movie := moviesService.GetMovieByTitle(data)
	// Create a request
	request := types.Request{
		TargetMovie: *movie,
		Movies:      moviesService.Movies,
	}

	// Serialize the request
	requestBytes, err := json.Marshal(request)
	if err != nil {
		return types.Response{}, "Error al serializar la petición."
	}

	// Send the request
	_, err = conn.Write(requestBytes)
	if err != nil {
		return types.Response{}, "Error al enviar la petición."
	}

	// Receive the response
	responseBytes, err := io.ReadAll(conn)
	if err != nil {
		return types.Response{}, "Error al recibir la respuesta."
	}

	// Deserialize the response
	var response types.Response // Create a response
	err = json.Unmarshal(responseBytes, &response)
	if err != nil {
		return types.Response{}, "Error al deserializar la respuesta."
	}

	// Check if there was an error
	if response.Error != "" {
		return types.Response{}, response.Error
	}

	// Check if there are no recomendations
	if len(response.MovieResponse) == 0 {
		return types.Response{}, "No se encontraron recomendaciones."
	}

	return response, ""
}
