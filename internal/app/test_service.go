package app

import (
	"context"
	"fmt"
	"log"

	"github.com/google/generative-ai-go/genai"
	"google.golang.org/api/iterator"
)

type TestService struct {
	geminiClient *genai.Client
}

func NewTestService(client *genai.Client) *TestService {
	return &TestService{
		geminiClient: client,
	}
}

func (s *TestService) Prompt(body string) ([]genai.Part, *genai.GenerateContentResponse) {
	ctx := context.Background()

	// promptText := `
	// 	Saya memiliki data penjualan dummy untuk perusahaan logistik:
	// 	- Bulan Januari: 100 paket, Omset Rp 5.000.000
	// 	- Bulan Februari: 150 paket, Omset Rp 7.500.000
	// 	- Bulan Maret: 120 paket, Omset Rp 6.000.000
	// 	Tolong analisis singkat tren penjualan ini dan berikan 1 saran strategi pemasaran.
	// `

	// Generate model
	model := s.geminiClient.GenerativeModel("gemini-2.5-flash")

	// Send request to Gemini
	fmt.Println("Processing request to Gemini...")
	resp, err := model.GenerateContent(ctx, genai.Text(body))

	if err != nil {
		log.Fatalf("Failed to get respond from Gemini: %v", err)
	}

	// Get result
	result := summarizeRespond(resp)

	return result, resp
}

func summarizeRespond(resp *genai.GenerateContentResponse) []genai.Part {
	var result []genai.Part

	for _, cand := range resp.Candidates {
		if cand.Content != nil {
			for _, part := range cand.Content.Parts {
				// fmt.Printf("%v", part)
				result = append(result, part)
			}
		}
	}

	return result
}

func (s *TestService) AvailModel() {
	ctx := context.Background()

	fmt.Println("Available models:")

	// Use the client that have injected
	iter := s.geminiClient.ListModels(ctx)

	for {
		m, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			log.Fatalf("Failed to retrive model list: %v", err)
		}

		// Only show models that capable for GenerateContent function
		for _, method := range m.SupportedGenerationMethods {
			if method == "generateContent" {
				fmt.Println("-", m.Name)
			}
		}
	}

	// Available models:
	// 	- models/gemini-2.5-flash
	// 	- models/gemini-2.5-pro
	// 	- models/gemini-2.0-flash
	// 	- models/gemini-2.0-flash-001
	// 	- models/gemini-2.0-flash-lite-001
	// 	- models/gemini-2.0-flash-lite
	// 	- models/gemini-2.5-flash-preview-tts
	// 	- models/gemini-2.5-pro-preview-tts
	// 	- models/gemma-4-26b-a4b-it
	// 	- models/gemma-4-31b-it
	// 	- models/gemini-flash-latest
	// 	- models/gemini-flash-lite-latest
	// 	- models/gemini-pro-latest
	// 	- models/gemini-2.5-flash-lite
	// 	- models/gemini-2.5-flash-image
	// 	- models/gemini-3-pro-preview
	// 	- models/gemini-3-flash-preview
	// 	- models/gemini-3.1-pro-preview
	// 	- models/gemini-3.1-pro-preview-customtools
	// 	- models/gemini-3.1-flash-lite-preview
	// 	- models/gemini-3.1-flash-lite
	// 	- models/gemini-3-pro-image-preview
	// 	- models/gemini-3-pro-image
	// 	- models/nano-banana-pro-preview
	// 	- models/gemini-3.1-flash-image-preview
	// 	- models/gemini-3.1-flash-image
	// 	- models/gemini-3.1-flash-lite-image
	// 	- models/gemini-3.5-flash
	// 	- models/gemini-omni-flash-preview
	// 	- models/lyria-3-clip-preview
	// 	- models/lyria-3-pro-preview
	// 	- models/gemini-3.1-flash-tts-preview
	// 	- models/gemini-robotics-er-1.5-preview
	// 	- models/gemini-robotics-er-1.6-preview
	// 	- models/gemini-2.5-computer-use-preview-10-2025
	// 	- models/antigravity-preview-05-2026
	// 	- models/deep-research-max-preview-04-2026
	// 	- models/deep-research-preview-04-2026
	// 	- models/deep-research-pro-preview-12-2025
}
