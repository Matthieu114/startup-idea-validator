package main

import (
	"log"
	"net/http"

	handler "github.com/Matthieu114/startup-idea-validator/internal/handlers"
	"github.com/joho/godotenv"
)

func main() {

	router := http.NewServeMux()
	router.HandleFunc("POST /validate", handler.ValidateIdeaHandler)

	err := http.ListenAndServe(":8000", router)

	if err != nil {
		log.Fatal(err)
	}

}

func init() {
    // Try to load the .env file, but ignore errors if it doesn't exist
    err := godotenv.Load()
    if err != nil {
        log.Println("No .env file found (this is normal in production)")
    }
}

