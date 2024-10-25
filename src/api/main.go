package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/joho/godotenv"
)

// Define a struct for the JSON response
type Response struct {
	Message string `json:"message"`
}

func main() {
	// Load env variable for the port
	godotenv.Load()
	PORT := os.Getenv("PORT")
	if PORT == "" {
		PORT = "8080"
	}

	// Set router
	router := http.NewServeMux()

	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// Create a response object
		response := Response{Message: "Hello World!"}
		fmt.Println("Request received")
		// sleep for two seconds
		time.Sleep(10 * time.Second)
		fmt.Println("Request processed")

		// Convert the response object to JSON
		jsonResponse, err := json.Marshal(response)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Set the Content-Type header to application/json
		w.Header().Set("Content-Type", "application/json")
		w.Write(jsonResponse)
	})

	// Start server
	fmt.Printf("Server is listening on port %s\n", PORT)
	err := http.ListenAndServe(":"+PORT, router)
	if err != nil {
		fmt.Printf("Error starting server: %s\n", err)
	}
}
