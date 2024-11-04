package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"github.com/LuisFelipePoma/Movies_Recomender_With_Golang/src/api/slave/model"
	"github.com/LuisFelipePoma/Movies_Recomender_With_Golang/src/api/types"
	"net"
	"os"
	"strconv"
	"strings"
)

func main() {
	// Leer el puerto del usuario
	port := ""
	if len(os.Args) > 1 {
		port = os.Args[1]
	} else {
		reader := bufio.NewReader(os.Stdin)
		fmt.Print("Enter port: ")
		port, _ = reader.ReadString('\n')
		port = strings.TrimSpace(port)

		// Validar que el puerto sea un número válido
		if _, err := strconv.Atoi(port); err != nil {
			fmt.Println("Invalid port number")
			return
		}
	}

	// Iniciar el servidor TCP
	ln, err := net.Listen("tcp", ":"+port)
	if err != nil {
		fmt.Println("Error starting TCP server:", err)
		return
	}
	defer ln.Close()

	fmt.Println("Slave node listening on port", port)

	for {
		// Aceptar conexiones entrantes
		conn, err := ln.Accept()
		if err != nil {
			fmt.Println("Error accepting connection:", err)
			continue
		}
		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {
	defer conn.Close()

	fmt.Println("Leyendo los datos entrantes....")
	// Decodificar el JSON recibido en una estructura Task
	var task struct {
		Movies      []types.Movie `json:"movies"`
		TargetMovie types.Movie   `json:"target_movie"`
	}
	decoder := json.NewDecoder(conn)
	if err := decoder.Decode(&task); err != nil {
		fmt.Println("Error al decodificar JSON:", err)
		return
	}
	fmt.Println("Nodo Esclavo recibió tarea...")

	fmt.Println("Calculando las peliculas similares....")

	// Crear una instancia de Recommender
	recommender := model.NewRecommender()

	// Obtener películas similares
	similarMovies := recommender.GetSimilarMovies(task.Movies, task.TargetMovie)

	// Imprimir la cantidad de películas similares encontradas
	fmt.Println("Se encontro", len(similarMovies), "similar movies")

	// Codificar el resultado en JSON y enviarlo de vuelta al cliente
	encoder := json.NewEncoder(conn)
	if err := encoder.Encode(similarMovies); err != nil {
		fmt.Println("Error al codificar JSON:", err)
		return
	}
	fmt.Println("Nodo Esclavo envió resultado")
}
