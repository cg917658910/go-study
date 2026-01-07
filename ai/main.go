package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"google.golang.org/genai"
)

var model = "gemini-3-flash-preview"

func main() {
	ctx := context.Background()

	clientConfig := &genai.ClientConfig{
		APIKey:  "AIzaSyAWljYwSfLje5NidjtIkHLDvcStYW8h4BA",
		Backend: genai.BackendGeminiAPI,
	}
	_ = clientConfig
	client, err := genai.NewClient(ctx, clientConfig)
	if err != nil {
		log.Fatal(err)
	}

	result, err := client.Models.GenerateContent(
		ctx,
		"gemini-2.5-flash-image",
		genai.Text("Create a picture of a nano banana dish in a "+
			" fancy restaurant with a Gemini theme"), nil,
	)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("result:", result)
	for _, part := range result.Candidates[0].Content.Parts {
		if part.Text != "" {
			fmt.Println(part.Text)
		} else if part.InlineData != nil {
			imageBytes := part.InlineData.Data
			outputFilename := "gemini_generated_image.png"
			_ = os.WriteFile(outputFilename, imageBytes, 0644)
		}
	}
}
