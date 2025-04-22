package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
)

type IdeaRequest struct {
	Value string `json:"idea"`
}

type IdeaReponse struct {
	Summary     string   `json:"summary"`
	Score       float32  `json:"score"`
	Suggestions []string `json:"suggestions"`
}

func main() {
	router := http.NewServeMux()
	router.HandleFunc("POST /validate", receiveAndLogJSON)

	err := http.ListenAndServe(":8000", router)

	if err != nil {
		log.Fatal(err)
	}
}

func receiveAndLogJSON(w http.ResponseWriter, r *http.Request) {

	var i IdeaRequest

	err := json.NewDecoder(r.Body).Decode(&i)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	fmt.Println("Received idea:", i.Value)

	file, err := os.Open("test.json")

	if err != nil {
		log.Fatal(err)
		return
	}

	defer file.Close()

	decoder := json.NewDecoder(file)
	var ideaResponse IdeaReponse
	err = decoder.Decode(&ideaResponse)

	if err != nil {
		log.Fatal(err)
		return
	}

	jsonContent, err := json.Marshal(ideaResponse)

	if err != nil {
		log.Fatal(err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonContent)

	// read i value and return new
	// fmt.Printf("%+v\n", ideaResponse)
	// fmt.Println(i.Value)
}
