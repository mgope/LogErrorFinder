package api

import (
	"context"
	"encoding/json"
	"fmt"
	"logerrorfinder/model"
	"os"

	"github.com/joho/godotenv"
	"google.golang.org/genai"
)

type GeminiService struct {
	client *genai.Client
	model  string
}

func NewGeminiService(ctx context.Context) (*GeminiService, error) {

	_ = godotenv.Load()

	apiKey := os.Getenv("GEMINI_API_KEY")

	if apiKey == "" {
		return nil, fmt.Errorf("GEMINI_API_KEY is not set")
	}

	client, err := genai.NewClient(ctx, &genai.ClientConfig{
		APIKey:  apiKey,
		Backend: genai.BackendGeminiAPI,
	})

	if err != nil {
		return nil, err
	}

	return &GeminiService{
		client: client,
		model:  "gemini-3.5-flash-lite",
	}, nil
}

func (g *GeminiService) AnalyzeLogs(
	ctx context.Context,
	logData model.LogForLLM,
) (*model.LLMResponse, error) {

	logJSON, err := json.MarshalIndent(logData, "", "  ")
	if err != nil {
		return nil, err
	}

	prompt := fmt.Sprintf(`
You are a production log analysis assistant.

Analyze the following logs.

The selected log is the ERROR that needs investigation.

The previous logs provide context.

Your job is to determine:

1. Root cause
2. Recommended Actions
4. Evidence from the provided logs
3. Confidence in percentage (0-100) 

IMPORTANT:
- Do not invent evidence.
- Only use information present in the logs.
- If the evidence is insufficient, say so.
- Return ONLY valid JSON.
- Do not use markdown.
- Do not wrap the JSON in a code block.

Required JSON format:

{
	"root_cause": "string",
  	"recommended_actions": [
    	"string",
    	"string"
  	],
  	"evidence": [
    	"string",
    	"string"
	],
	"confidence": 0.0
}

Logs:

%s
`, string(logJSON))

	result, err := g.client.Models.GenerateContent(
		ctx,
		g.model,
		genai.Text(prompt),
		&genai.GenerateContentConfig{
			ResponseMIMEType: "application/json",
		},
	)

	if err != nil {
		return nil, err
	}

	responseText := result.Text()

	var response model.LLMResponse

	if err := json.Unmarshal(
		[]byte(responseText),
		&response,
	); err != nil {
		return nil, fmt.Errorf(
			"failed to parse Gemini JSON: %w\nresponse: %s",
			err,
			responseText,
		)
	}

	return &response, nil
}
