package main

import (
	"bufio"
	"fmt"
	m "github.com/LuisFelipePoma/Movies_Recomender_With_Golang/src/api/services"
	"os"
	"os/exec"
	"strings"
)

var moviesService = m.NewMovies()

func main() {
	// Leer el archivo JSON con las películas una vez
	if err := moviesService.LoadMovies("../../database/data/data_clean.json"); err != nil {
		fmt.Println("Error al leer el archivo JSON:", err)
		return
	}

	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Println("Seleccione una opción:")
		fmt.Println("1. Recomendar en base a una película")
		fmt.Println("2. Mostrar recientes recomendaciones")
		fmt.Println("3. Filtrar por género las últimas recomendaciones")
		fmt.Println("4. Filtrar por voteAverage las últimas recomendaciones")
		fmt.Println("5. Salir")
		fmt.Print("Opción: ")
		optionStr, _ := reader.ReadString('\n')
		optionStr = strings.TrimSpace(optionStr)
		option := 0
		fmt.Sscanf(optionStr, "%d", &option)

		switch option {
		case 1:
			fmt.Println("Ingrese el título de la película:")
			title, _ := reader.ReadString('\n')
			title = strings.TrimSpace(title)

			moviesService.Recommendations, moviesService.LastRecomendation = moviesService.RecomendMoviesByTitle(title)
			fmt.Println("Recomendaciones:")
			moviesService.PrintRecomendationsDetails()

		case 2:
			clearConsole()
			if moviesService.IsEmptyRecommendations() {
				continue
			}

			fmt.Printf("Últimas recomendaciones basadas en la película: %s\n", moviesService.LastRecomendation)
			moviesService.PrintRecomendationsDetails()

		case 3:
			clearConsole()
			if moviesService.IsEmptyRecommendations() {
				continue
			}

			fmt.Println("Ingrese el género:")
			genre, _ := reader.ReadString('\n')
			genre = strings.TrimSpace(genre)
			response := moviesService.GetRecomendationsByGenre(genre)
			fmt.Println("Recomendaciones filtradas por género:")
			moviesService.PrintMoviesDetails(response)

		case 4:
			clearConsole()
			if moviesService.IsEmptyRecommendations() {
				continue
			}

			fmt.Println("Ingrese el voteAverage mínimo:")
			voteAverageStr, _ := reader.ReadString('\n')
			voteAverageStr = strings.TrimSpace(voteAverageStr)
			minVoteAverage := 0.0
			fmt.Sscanf(voteAverageStr, "%f", &minVoteAverage)
			response := moviesService.GetMoviesByVoteAverage(minVoteAverage)
			fmt.Println("Recomendaciones filtradas por voteAverage:")
			moviesService.PrintMoviesDetails(response)

		case 5:
			fmt.Println("Saliendo...")
			return
		default:
			fmt.Println("Opción no válida, intente de nuevo.")
		}
	}
}

// OTHER
func clearConsole() {
	cmd := exec.Command("cmd", "/c", "cls")
	cmd.Stdout = os.Stdout
	cmd.Run()
}
