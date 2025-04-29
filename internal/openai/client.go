package openai

import (
	"context"
	"fmt"
	"os"

	"github.com/openai/openai-go"
	"github.com/openai/openai-go/option"
)

func GetOpenAiApiResponse(message string) string {

	client := openai.NewClient(
		option.WithAPIKey(os.Getenv("OPEN_AI_API_KEY")),
	)

	prompt := fmt.Sprintf(`
You are a startup idea evaluator.

Given a business idea, evaluate it based on the following criteria for both Europe and the USA:

1. Size of the Market (small, medium, large, very large) + a short explanation. Also, provide a confidence score (0-100%%) on the size of the market estimate for each region.
2. Amount of Competition (low, medium, high) + a short explanation. Also, provide a confidence score (0-100%%) on the competition estimate for each region.
3. Growth of the Market (declining, stagnant, growing, rapidly growing) + a short explanation. Also, provide a confidence score (0-100%%) on the growth estimate for each region.
4. Purchasing Power of the Target Customers (low, medium, high) + a short explanation. Also, provide a confidence score (0-100%%) on the purchasing power estimate for each region.
5. Painfulness of the Problem (minor inconvenience, moderate pain, very painful, urgent critical need) + a short explanation. Also, provide a confidence score (0-100%%) on the pain estimate for each region.

After that:
- Suggest 3 improvements or new directions based on the idea.
- List 3 competitors or similar products if they exist (name + short description).

Return the result as a JSON object in the following structure:

{
  "market_size": {
    "rating": "large",
    "explanation": "Farmers markets are increasingly popular in urban areas, especially in the U.S. and Europe.",
    "confidence_score": 90
  },
  "competition": {
    "rating": "medium",
    "explanation": "Several apps and websites already list farmers markets, but few offer real-time or localized recommendations.",
    "confidence_score": 80
  },
  "market_growth": {
    "rating": "growing",
    "explanation": "Interest in local food sourcing and sustainability is growing steadily.",
    "confidence_score": 85
  },
  "purchasing_power": {
    "rating": "medium",
    "explanation": "Customers at farmers markets tend to have moderate to high disposable income depending on the region.",
    "confidence_score": 75
  },
  "problem_painfulness": {
    "rating": "moderate pain",
    "explanation": "People want to find farmers markets easily but existing solutions are fragmented.",
    "confidence_score": 90
  },
  "suggestions": [
    "Allow farmers to update real-time product availability.",
    "Include recipes based on available products at markets.",
    "Offer loyalty rewards for visiting multiple markets."
  ],
  "competitors": [
    {
      "name": "LocalHarvest",
      "description": "Directory of farmers markets and local food sources across the U.S."
    },
    {
      "name": "Farmstand App",
      "description": "Mobile app showing nearby farmers markets worldwide."
    },
    {
      "name": "USDA Farmers Market Directory",
      "description": "Official USDA listing of markets in the United States."
    }
  ]
}`)

	chatCompletion, err := client.Chat.Completions.New(context.Background(), openai.ChatCompletionNewParams{
		Messages: []openai.ChatCompletionMessageParamUnion{
			openai.SystemMessage(prompt),
			openai.UserMessage(message),
		},
		Model: openai.ChatModelGPT4o,
	})

	if err != nil {
		panic(err.Error())
	}

	return chatCompletion.Choices[0].Message.Content
}
