package validate

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/Matthieu114/startup-idea-validator/internal/model"
	"github.com/Matthieu114/startup-idea-validator/internal/openai"
)

func ReceiveAndLogJSON(w http.ResponseWriter, r *http.Request) {
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
