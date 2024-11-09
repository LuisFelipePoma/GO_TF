package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
	"os"
	"github.com/LuisFelipePoma/Movies_Recomender_With_Golang/src/backend/slave/model"
	"github.com/LuisFelipePoma/Movies_Recomender_With_Golang/src/backend/types"
)

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

// handleConnection handles incoming connections
func handleConnection(conn net.Conn) {
	defer conn.Close()

	fmt.Println("Leyendo los datos entrantes....")
	// Decodificate the JSON data
	var task struct {
		Movies      []types.Movie `json:"movies"`
		TargetMovie types.Movie   `json:"target_movie"`
	}
	decoder := json.NewDecoder(conn) // Create a JSON decoder that reads from
	// Parse the JSON data
	if err := decoder.Decode(&task); err != nil {
		fmt.Println("Error al decodificar JSON:", err)
		return
	}
	fmt.Println("Nodo Esclavo recibió tarea...")
	fmt.Println("Calculando las peliculas similares....")

	// Create a new recommender instance
	recommender := model.NewRecommender()

	// Get similar movies
	similarMovies := recommender.GetSimilarMovies(task.Movies, task.TargetMovie)
	fmt.Println("Se encontro", len(similarMovies), "similar movies")

	// Send the result back to the master node
	encoder := json.NewEncoder(conn)
	if err := encoder.Encode(similarMovies); err != nil {
		fmt.Println("Error al codificar JSON:", err)
		return
	}
	fmt.Println("Nodo Esclavo envió resultado")
}
