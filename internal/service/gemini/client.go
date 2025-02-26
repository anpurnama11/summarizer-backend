package gemini

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"anpurnama/summarizer-backend/internal/repository"

	"github.com/google/generative-ai-go/genai"
	"github.com/joho/godotenv"
	"google.golang.org/api/option"
)

type Client struct {
	model           *genai.Client
	modelPool       chan *genai.GenerativeModel
	styleRepository repository.StyleRepository
}

func NewClient(styleRepo repository.StyleRepository) (*Client, error) {
	if err := godotenv.Load(); err != nil {
		return nil, fmt.Errorf("error loading .env file: %w", err)
	}

	ctx := context.Background()
	client, err := genai.NewClient(ctx, option.WithAPIKey(os.Getenv("GEMINI_API_KEY")))
	if err != nil {
		return nil, fmt.Errorf("failed to create Gemini client: %w", err)
	}

	return &Client{
		model:           client,
		modelPool:       make(chan *genai.GenerativeModel, 10), // buffer size of 10
		styleRepository: styleRepo,
	}, nil
}

func (c *Client) Summarize(ctx context.Context, content string, styleName string) (string, error) {
	start := time.Now()
	model := c.getModelFromPool()

	// Get style from repository
	style, err := c.styleRepository.GetByName(ctx, styleName)
	if err != nil {
		return "", fmt.Errorf("failed to get style: %w", err)
	}
	if style == nil {
		return "", fmt.Errorf("style '%s' not found", styleName)
	}

	prompt := fmt.Sprintf("%s\n\n%s", style.PromptTemplate, content)

	resp, err := model.GenerateContent(ctx, genai.Text(prompt))
	if err != nil {
		return "", fmt.Errorf("failed to generate summary: %w", err)
	}

	if len(resp.Candidates) == 0 {
		return "", fmt.Errorf("no response received from Gemini")
	}

	candidate := resp.Candidates[0]
	if candidate.Content == nil {
		return "", fmt.Errorf("empty content in response")
	}

	var summary string
	for _, part := range candidate.Content.Parts {
		summary += fmt.Sprint(part)
	}

	c.returnModelToPool(model)

	log.Printf("Process completed in %s", time.Since(start))
	return summary, nil
}

func (c *Client) getModelFromPool() *genai.GenerativeModel {
	select {
	case model := <-c.modelPool:
		return model
	default:
		return c.model.GenerativeModel("gemini-2.0-flash-001")
	}
}

func (c *Client) returnModelToPool(model *genai.GenerativeModel) {
	select {
	case c.modelPool <- model:
	default:
		// do nothing if the pool is full
	}
}
