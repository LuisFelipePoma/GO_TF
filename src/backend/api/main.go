// go
package main

import (
	"encoding/json"
	"fmt"
	"github.com/LuisFelipePoma/Movies_Recomender_With_Golang/src/backend/types"
	"github.com/gorilla/websocket"
	"io"
	"net"
	"net/http"
	"os"
	"strconv"
	"sync"
	"time"
)

var nodeMasterPort = os.Getenv("MASTER_NODE") // Puerto del nodo maestro

var upgrader = websocket.Upgrader{} // Utilizado para actualizar la conexión HTTP a WebSocket
var clients = make(map[*websocket.Conn]int)
var clientsMutex = sync.Mutex{}

// Message represents the recommendation message structure
type Message struct {
	User            string   `json:"user"`
	Recommendations []string `json:"recommendations"`
}

func main() {
	// Leer el puerto desde la variable de entorno
	port := os.Getenv("PORT")
	if port == "" {
		port = "8000" // Puerto por defecto si no está configurado
	}
	// Configurar los manejadores HTTP
	setupRoutes()

	// Iniciar una rutina para enviar recomendaciones periódicamente
	go sendPeriodicRecommendations()

	fmt.Println("Server is running on port", port)
	http.ListenAndServe(":"+port, nil)
}

// Configurar los manejadores HTTP
func setupRoutes() {
	http.HandleFunc("/api/movies/similar", corsMiddleware(getSimilarMovies)) // GET
	http.HandleFunc("/api/movies/search", corsMiddleware(getMoviesBySearch)) // GET
	http.HandleFunc("/api/movies", corsMiddleware(getAllMovies))             // GET
	http.HandleFunc("/ws", handleWebSocketConnections)                       // WebSocket endpoint
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
// Handler to send recommendations
func handleWebSocketConnections(w http.ResponseWriter, r *http.Request) {
	// Update the upgrader to allow all connections
	upgrader.CheckOrigin = func(r *http.Request) bool { return true } // Permitir todas las conexiones
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println("Error al actualizar a WebSocket:", err)
		return
	}
	defer ws.Close()

	// Registrar el nuevo cliente con un valor por defecto de 'n'
	clientsMutex.Lock()
	clients[ws] = 5 // Valor por defecto inicial
	clientsMutex.Unlock()
	fmt.Println("Nuevo cliente conectado.")

	// Listen for new messages
	for {
		_, message, err := ws.ReadMessage()
		if err != nil {
			fmt.Println("Cliente desconectado:", err)
			delete(clients, ws)
			break
		}
		// Parsear el mensaje para obtener 'n'
		var msg struct {
			N int `json:"n"`
		}
		err = json.Unmarshal(message, &msg)
		if err != nil {
			fmt.Println("Error al parsear el mensaje:", err)
			continue
		}

		// Actualizar 'n' para este cliente
		clientsMutex.Lock()
		clients[ws] = msg.N
		clientsMutex.Unlock()
		fmt.Printf("El cliente %s solicitó %d recomendaciones.\n", ws.RemoteAddr(), msg.N)
	}
}

// Function for sending periodic recommendations
func sendPeriodicRecommendations() {
	for {
		time.Sleep(25 * time.Second) // Wait 10 seconds

		clientsMutex.Lock()
		if len(clients) == 0 {
			clientsMutex.Unlock()
			continue
		}

		// Make a copy of the current clients to iterate
		currentClients := make(map[*websocket.Conn]int)
		for client, n := range clients {
			currentClients[client] = n
		}
		clientsMutex.Unlock()

		fmt.Println("Enviando recomendaciones a los clientes.")
		// Send the recommendations
		for client, n := range clients {
			// Create a request
			request := types.TaskDistributed{
				Type: types.TaskUserRecomend,
				Data: types.TaskData{
					Quantity: n,
				},
			}
			// Get response from the master node
			response, errorMessage := handleMasterConection(request)
			if errorMessage != "" {
				fmt.Println(errorMessage)
				continue
			}
			fmt.Println("Enviando recomendaciones al cliente: ", client.RemoteAddr())
			// Send the response to the client
			err := client.WriteJSON(response)
			if err != nil {
				fmt.Println("Error al enviar recomendaciones:", err)
				client.Close()
				// remove the client from the list
				clientsMutex.Lock()
				delete(clients, client)
				clientsMutex.Unlock()
			}
		}
	}
}

// Handle the similar movies request
func getSimilarMovies(w http.ResponseWriter, r *http.Request) {
	fmt.Println("GET /api/movies/similar")
	// Leer el cuerpo de la petición
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Unable to read request body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	var data map[string]interface{}
	if err := json.Unmarshal(body, &data); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	// Obtener y validar 'id'
	id, ok := data["id"].(float64)
	if !ok {
		http.Error(w, "ID is required and must be a number", http.StatusBadRequest)
		return
	}

	// Obtener y validar 'n'
	nFloat, ok := data["n"].(float64)
	if !ok {
		http.Error(w, "Quantity 'n' is required and must be a number", http.StatusBadRequest)
		return
	}
	quantity := int(nFloat)
	movieId := int(id)
	// Create a request
	request := types.TaskDistributed{
		Type: types.TaskRecomend,
		Data: types.TaskData{
			Quantity: quantity,
			TaskRecomendations: &types.TaskRecomendations{
				MovieId: movieId,
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
	fmt.Println("GET /api/movies")
	genre := r.URL.Query().Get("genre")
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

	// Create Request
	request := types.TaskDistributed{
		Type: types.TaskGet,
		Data: types.TaskData{
			TaskSearch: &types.TaskSearchQuery{
				Query: genre,
			},
			Quantity: n,
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

func getMoviesBySearch(w http.ResponseWriter, r *http.Request) {
	fmt.Println("GET /api/movies/search")
	query := r.URL.Query().Get("query")
	n := r.URL.Query().Get("n")

	// Handle errors
	if query == "" {
		http.Error(w, "Query is required", http.StatusBadRequest)
		return
	}
	if n == "" {
		n = "10"
	}

	// Convert n to int
	nInt, err := strconv.Atoi(n)
	if err != nil {
		http.Error(w, "Invalid number format", http.StatusBadRequest)
		return
	}
	// Create a request
	request := types.TaskDistributed{
		Type: types.TaskSearch,
		Data: types.TaskData{
			Quantity: nInt,
			TaskSearch: &types.TaskSearchQuery{
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
