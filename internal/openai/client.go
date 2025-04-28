package openai

import (
	"context"
	"os"

	"github.com/openai/openai-go"
	"github.com/openai/openai-go/option"
)


func GetOpenAiApiResponse(message string) string {


	
	client := openai.NewClient(
		option.WithAPIKey(os.Getenv("OPEN_AI_API_KEY")),
	)

	chatCompletion, err := client.Chat.Completions.New(context.TODO(), openai.ChatCompletionNewParams{
		Messages: []openai.ChatCompletionMessageParamUnion{
			openai.UserMessage(message),
		},
		Model: openai.ChatModelGPT4o,
	})

	if err != nil {
		panic(err.Error())
	}

	return chatCompletion.Choices[0].Message.Content
}


