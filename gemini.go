package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/google/generative-ai-go/genai"
	"google.golang.org/api/option"
)

func GeminiImage(imgData []byte) string {
	ctx := context.Background()
	client, err := genai.NewClient(ctx, option.WithAPIKey(geminiKey))
	if err != nil {
		log.Fatal(err)
	}
	defer client.Close()

	model := client.GenerativeModel("gemini-pro-vision")
	prompt := []genai.Part{
		genai.ImageData("jpg", imgData),
		genai.Text("Describe this image with scientific detail, reply in zh-TW:"),
	}
	resp, err := model.GenerateContent(ctx, prompt...)

	if err != nil {
		log.Fatal(err)
	}

	bs, _ := json.MarshalIndent(resp, "", "    ")
	fmt.Println(string(bs))
	return string(bs)
}

// Gemini Chat Complete: Iput a prompt and get the response string.
func GeminiChatComplete(req string) string {
	ctx := context.Background()
	client, err := genai.NewClient(ctx, option.WithAPIKey(geminiKey))
	if err != nil {
		log.Fatal(err)
	}
	defer client.Close()
	model := client.GenerativeModel("gemini-pro")
	cs := model.StartChat()

	send := func(msg string) *genai.GenerateContentResponse {
		fmt.Printf("== Me: %s\n== Model:\n", msg)
		res, err := cs.SendMessage(ctx, genai.Text(msg))
		if err != nil {
			log.Fatal(err)
		}
		return res
	}

	res := send(req)
	return printResponse(res)
}

func printResponse(resp *genai.GenerateContentResponse) string {
	var ret string
	for _, cand := range resp.Candidates {
		for _, part := range cand.Content.Parts {
			ret = ret + fmt.Sprintf("%v", part)
			fmt.Println(part)
		}
	}
	fmt.Println("---")
	return ret + "\n---"
}
