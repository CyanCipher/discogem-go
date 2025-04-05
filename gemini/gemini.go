package gemini

import (
	"context"
	"errors"
	"log"
	"os"
	"path/filepath"

	"github.com/google/generative-ai-go/genai"
	"github.com/joho/godotenv"
	"google.golang.org/api/option"
)

func AskGemini(prompt string) (string, error) {

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file!")
	}

	// Setting up the client, and model
	ctx := context.Background()
	client, err := genai.NewClient(ctx, option.WithAPIKey(os.Getenv("GEMINI_API_KEY")))
	if err != nil {
		log.Fatal(err)
	}
	defer client.Close()
	model := client.GenerativeModel("gemini-2.0-flash")

	resp, err := model.GenerateContent(ctx, genai.Text(prompt))
	if err != nil {
		log.Fatal(err)
	}

	return GetResponse(resp)
}

func ImgQuery(prompt string, file_name string) (string, error) {

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file!")
	}

	//setting up the client and model
	ctx := context.Background()
	client, err := genai.NewClient(ctx, option.WithAPIKey(os.Getenv("GEMINI_API_KEY")))
	if err != nil {
		log.Fatal(err)
	}
	defer client.Close()
	model := client.GenerativeModel("gemini-2.0-flash")

	cwd, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	imgData, err := os.ReadFile(filepath.Join(cwd, "Media", file_name))
	if err != nil {
		log.Fatal(err)
	}

	resp, err := model.GenerateContent(ctx,
		genai.Text("Tell me something about this image"),
		genai.ImageData("png", imgData))

	if err != nil {
		log.Fatal(err)
	}

	return GetResponse(resp)
}

func GetResponse(resp *genai.GenerateContentResponse) (string, error) {
	var response string

	for _, cand := range resp.Candidates {
		if cand.Content != nil {
			for _, part := range cand.Content.Parts {
				if textpart, ok := part.(genai.Text); ok {
					response += string(textpart)
				}
			}
		} else {
			return "", errors.New("No text was generated as response!")
		}
	}

	return response, nil
}
