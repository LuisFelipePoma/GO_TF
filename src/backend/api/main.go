package main

import (
	"encoding/json"
	"fmt"
	"github.com/LuisFelipePoma/Movies_Recomender_With_Golang/src/backend/types"
	"io"
	"net"
	"net/http"
	"os"
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
	http.HandleFunc("/api/movies/similar", getMovies)
}

// HANDLERS
func getMovies(w http.ResponseWriter, r *http.Request) {
	// Get title from args
	fmt.Println("GET /api/movies/similar")
	title := r.URL.Query().Get("title")
	fmt.Println(title)
	response, errorMessage := handleOption(1, title)
	if errorMessage != "" {
		fmt.Println(errorMessage)
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// Conect to the master node
func handleOption(option int, data string) (types.Response, string) {
	fmt.Println("Connecting to master node : ", "master:"+nodeMasterPort)
	conn, err := net.Dial("tcp", "master:"+nodeMasterPort) // Connect to the master node
	if err != nil {
		return types.Response{}, "Error al conectar con el nodo maestro."
	}
	defer conn.Close()

	// Create a request
	request := types.Request{
		Option: option,
		Data:   data,
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
