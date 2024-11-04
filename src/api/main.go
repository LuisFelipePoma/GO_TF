package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"github.com/LuisFelipePoma/Movies_Recomender_With_Golang/src/api/types"
	"io"
	"net"
	"os"
	"strings"
)

var nodeMasterPort = "localhost:8081"

func main() {
	reader := bufio.NewReader(os.Stdin) // Create a reader
	for {
		fmt.Println("Seleccione una opción:")
		fmt.Println("1. Recomendar en base a una película")
		fmt.Println("2. Mostrar recientes recomendaciones")
		fmt.Println("3. Filtrar por género las últimas recomendaciones")
		fmt.Println("4. Filtrar por voteAverage las últimas recomendaciones")
		fmt.Println("5. Salir")
		fmt.Print("Opción: ")
		// Read the option
		optionStr, _ := reader.ReadString('\n')  // Read until the newline character
		optionStr = strings.TrimSpace(optionStr) // Remove the newline character
		option := 0
		fmt.Sscanf(optionStr, "%d", &option) // Convert the string to an integer

		switch option {
		case 1:
			fmt.Println("Ingrese el título de la película:")
			title, _ := reader.ReadString('\n')
			title = strings.TrimSpace(title)
			response, errorMessage := handleOption(1, title)
			if errorMessage != "" {
				fmt.Println(errorMessage)
				continue
			}
			fmt.Println("Recomendaciones:")
			printRecomendationsDetails(response.MovieResponse)

		case 2:
			response, errorMessage := handleOption(2, "")
			if errorMessage != "" {
				fmt.Println(errorMessage)
				continue
			}
			fmt.Printf("\n\nÚltimas recomendaciones basadas en la ultima película ingresada: %s\n\n", response.TargetMovie)
			printRecomendationsDetails(response.MovieResponse)

		case 3:
			fmt.Println("Ingrese el género:")
			genre, _ := reader.ReadString('\n')
			genre = strings.TrimSpace(genre)
			response, errorMessage := handleOption(3, genre)
			if errorMessage != "" {
				fmt.Println(errorMessage)
				continue
			}
			fmt.Printf("Recomendaciones de %s filtradas por género:\n", response.TargetMovie)
			printRecomendationsDetails(response.MovieResponse)

		case 4:
			fmt.Println("Ingrese el voteAverage mínimo:")
			voteAverageStr, _ := reader.ReadString('\n')
			voteAverageStr = strings.TrimSpace(voteAverageStr)

			response, errorMessage := handleOption(4, voteAverageStr)
			if errorMessage != "" {
				fmt.Println(errorMessage)
				continue
			}
			fmt.Printf("Recomendaciones de %s filtradas por voteAverage:\n", response.TargetMovie)
			printRecomendationsDetails(response.MovieResponse)

		case 5:
			fmt.Println("Saliendo...")
			return
		default:
			fmt.Println("Opción no válida, intente de nuevo.")
		}
	}
}

// Print the details of the recomendations
func printRecomendationsDetails(recomendations []types.MovieResponse) {
	for _, movie := range recomendations {
		fmt.Printf("Title: %s\n", movie.Title)
		fmt.Printf("Vote Average: %.2f\n", movie.VoteAverage)
		fmt.Printf("Genres: %s\n", movie.Genres)
		fmt.Println("-----------------------------")
	}
}

// Conect to the master node
func handleOption(option int, data string) (types.Response, string) {
	conn, err := net.Dial("tcp", nodeMasterPort) // Connect to the master node
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
