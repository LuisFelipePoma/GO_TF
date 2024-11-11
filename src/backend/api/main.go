package main

import (
	"encoding/json"
	"fmt"
	"github.com/LuisFelipePoma/Movies_Recomender_With_Golang/src/backend/types"
	"io"
	"net"
	"net/http"
	"os"
	"strconv"
)

var nodeMasterPort = os.Getenv("MASTER_NODE") // Port of the master node

func main() {
	// Read
	port := os.Getenv("PORT")

	// Configurar los manejadores HTTP
	setupRoutes()

	fmt.Println("Server is running on port ", port)
	http.ListenAndServe(":"+port, nil)

}

// Configurar los manejadores HTTP
func setupRoutes() {
	http.HandleFunc("/api/movies/similar", corsMiddleware(getSimilarMovies)) // GET
	http.HandleFunc("/api/movies/id", corsMiddleware(getById))               // GET
	http.HandleFunc("/api/movies/search", corsMiddleware(getMoviesBySearch)) // GET
	http.HandleFunc("/api/movies", corsMiddleware(getAllMovies))             // GET
}

// Middleware para manejar CORS
func corsMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}
		next(w, r)
	}
}

// HANDLERS
func getSimilarMovies(w http.ResponseWriter, r *http.Request) {
	fmt.Println("GET /api/movies/similar")
	// Get title from args
	title := r.URL.Query().Get("title")
	quantity := r.URL.Query().Get("n")

	// handle errors
	if title == "" {
		http.Error(w, "Title is required", http.StatusBadRequest)
		return
	}

	if quantity == "" {
		http.Error(w, "Quantity is required", http.StatusBadRequest)
		return
	}
	// Convert quantity to int
	quantityInt, err := strconv.Atoi(quantity)
	if err != nil {
		http.Error(w, "Invalid quantity format", http.StatusBadRequest)
		return
	}

	// Create a request
	request := types.TaskDistributed{
		Type: types.TaskRecomend,
		Data: types.TaskData{
			Quantity: quantityInt,
			TaskRecomendations: &types.TaskRecomendations{
				Title: title,
			},
		},
	}
	// Handle the connection to the master node
	response, errorMessage := handleMasterConection(request)
	if errorMessage != "" {
		fmt.Println(errorMessage)
	}
	// Send the response
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func getAllMovies(w http.ResponseWriter, r *http.Request) {
	// fmt.Println("GET /api/movies")
	// nStr := r.URL.Query().Get("n")
	// if nStr == "" {
	// 	nStr = "10"
	// }
	// // Convert n to int
	// n, err := strconv.Atoi(nStr)

	// if err != nil {
	// 	http.Error(w, "Invalid number format", http.StatusBadRequest)
	// 	return
	// }
	// response := err
	// w.Header().Set("Content-Type", "application/json")
	// json.NewEncoder(w).Encode(response)
}

func getById(w http.ResponseWriter, r *http.Request) {
	// fmt.Println("GET /api/movies")
	// idStr := r.URL.Query().Get("id")
	// if idStr == "" {
	// 	idStr = "1"
	// }
	// // Convert n to int
	// id, err := strconv.Atoi(idStr)

	// if err != nil {
	// 	http.Error(w, "Invalid number format", http.StatusBadRequest)
	// 	return
	// }
	// response := err
	// w.Header().Set("Content-Type", "application/json")
	// json.NewEncoder(w).Encode(response)
}

func getMoviesBySearch(w http.ResponseWriter, r *http.Request) {
	fmt.Println("GET /api/movies/search")
	query := r.URL.Query().Get("query")
	// Create a request
	request := types.TaskDistributed{
		Type: types.TaskSearch,
		Data: types.TaskData{
			TaskSearch: &types.TaskMasterSearch{
				Query: query,
			},
		},
	}

	// Handle the connection to the master node
	response, errorMessage := handleMasterConection(request)
	if errorMessage != "" {
		fmt.Println(errorMessage)
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// Conect to the master node
func handleMasterConection(taskMaster types.TaskDistributed) (types.Response, string) {
	conn, err := net.Dial("tcp", nodeMasterPort) // Connect to the master node
	if err != nil {
		return types.Response{}, "Error al conectar con el nodo maestro."
	}
	defer conn.Close()

	requestBytes, err := json.Marshal(taskMaster)
	if err != nil {
		return types.Response{}, "Error al serializar la petición."
	}

	// Send the request
	_, err = conn.Write(requestBytes)
	if err != nil {
		return types.Response{}, "Error al enviar la petición."
	}

	// Receive and deserialize the response
	response, errMsg := receiveAndDeserializeResponse(conn)
	if errMsg != "" {
		return types.Response{}, errMsg
	}

	// Check if there are no recomendations
	if len(response.MovieResponse) == 0 {
		return types.Response{}, "No se encontraron recomendaciones."
	}

	return response, ""
}

func receiveAndDeserializeResponse(conn net.Conn) (types.Response, string) {
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

	return response, ""
}
