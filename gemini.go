package main

import (
	"context"
	"fmt"
	"log"

	"github.com/google/generative-ai-go/genai"
	"google.golang.org/api/option"
)

func GeminiImage(imgData []byte) (string, error) {
	ctx := context.Background()
	client, err := genai.NewClient(ctx, option.WithAPIKey(geminiKey))
	if err != nil {
		log.Fatal(err)
	}
	defer client.Close()

	model := client.GenerativeModel("gemini-pro-vision")
	model.Temperature = 0.8
	prompt := []genai.Part{
		genai.ImageData("png", imgData),
		genai.Text("Describe this image with scientific detail, reply in zh-TW:"),
	}
	log.Println("Begin processing image...")
	resp, err := model.GenerateContent(ctx, prompt...)
	log.Println("Finished processing image...", resp)
	if err != nil {
		log.Fatal(err)
		return "", err
	}

	return printResponse(resp), nil
}

func startNewChatSession() *genai.ChatSession {
    ctx := context.Background()
    client, err := genai.NewClient(ctx, option.WithAPIKey(geminiKey))
    if err != nil {
        log.Fatal(err)
    }
    model := client.GenerativeModel("gemini-pro")
    model.Temperature = 0.3
	cs := model.StartChat()
	return cs
}

func send(cs *genai.ChatSession, msg string) *genai.GenerateContentResponse {
    if cs == nil {
        cs = startNewChatSession()
    }

    ctx := context.Background()
    fmt.Printf("== Me: %s\n== Model:\n", msg)
    res, err := cs.SendMessage(ctx, genai.Text(msg))
    if err != nil {
        log.Fatal(err)
    }
    return res
}


func printResponse(resp *genai.GenerateContentResponse) string {
	var ret string
	for _, cand := range resp.Candidates {
		for _, part := range cand.Content.Parts {
			ret = ret + fmt.Sprintf("%v", part)
			fmt.Println(part)
		}
	}
	return ret
}

