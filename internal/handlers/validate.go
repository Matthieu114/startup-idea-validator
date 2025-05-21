package handler

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/Matthieu114/startup-idea-validator/internal/model"
	"github.com/Matthieu114/startup-idea-validator/internal/openai"
	"github.com/Matthieu114/startup-idea-validator/internal/scraper"
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

	// Get OpenAI analysis
	openAiResponse := openai.GetOpenAiApiResponse(ideaValue)

	// Scrape Product Hunt for similar products
	products, err := scraper.ScrapeProductHunt()
	if err != nil {
		log.Printf("Error scraping Product Hunt: %v", err)
		// Continue with OpenAI response even if scraping fails
	}

	// Compare idea with scraped products
	similarProducts := scraper.CompareIdeaWithProducts(ideaValue, products)

	// Combine OpenAI response with Product Hunt results
	response := struct {
		OpenAIResponse  interface{}                  `json:"openai_analysis"`
		SimilarProducts []scraper.ProductHuntProduct `json:"similar_products"`
	}{
		OpenAIResponse:  openAiResponse,
		SimilarProducts: similarProducts,
	}

	jsonContent, err := json.Marshal(response)
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
