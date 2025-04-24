package main

import (
	"log"
	"net/http"

	validate "github.com/Matthieu114/startup-idea-validator/internal/handlers"
)

func main() {
	router := http.NewServeMux()
	router.HandleFunc("POST /validate", validate.ReceiveAndLogJSON)

	err := http.ListenAndServe(":8000", router)

	if err != nil {
		log.Fatal(err)
	}

}
