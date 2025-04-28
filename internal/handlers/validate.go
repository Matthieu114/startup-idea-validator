package handler

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/Matthieu114/startup-idea-validator/internal/model"
	"github.com/Matthieu114/startup-idea-validator/internal/openai"
)

func ValidateIdeaHandler(w http.ResponseWriter, r *http.Request) {
	var i model.IdeaRequest

	err := json.NewDecoder(r.Body).Decode(&i)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	ideaValue := i.Value

	if len(ideaValue) == 0 {
		http.Error(w, "Invalid input", http.StatusConflict)
		return
	}

	openAiResponse := openai.GetOpenAiApiResponse(ideaValue)

	jsonContent, err := json.Marshal(openAiResponse)

	if err != nil {
		log.Fatal(err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonContent)

}
/*
	receive an idea for a company 

	"app that regroups farmers markets around you"

	get JSON result 

	1. returns --> score of how good the app idea is (this will be based on a few criteria)

		1. Size of Market (niche or not)
		2. Amount of competition ? not sure
		3. Growth of market
		4. Purchasing Power of customers
		5. How painful is the problem you are trying to resolve for the customers

	2. new Ideas + suggestions --> given by open AI
	
	3. Competition / products that already exist (this will need to be scraped / found off the web idk how yet)

*/
