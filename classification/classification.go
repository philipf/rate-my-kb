package classification

import (
	"context"
	"errors"
	"fmt"
	"ratemykb/config"
	"strings"

	"github.com/tmc/langchaingo/llms"
	"github.com/tmc/langchaingo/llms/ollama"
)

// Package classification will handle the quality classification of scanned files

// Classification represents the quality classification of a file
type Classification string

const (
	// ClassificationEmpty indicates the file is effectively empty
	ClassificationEmpty Classification = "Empty"

	// ClassificationLowQuality indicates the file is of low quality or low effort
	ClassificationLowQuality Classification = "Low quality/low effort"

	// ClassificationGoodEnough indicates the file has good enough content
	ClassificationGoodEnough Classification = "Good enough"

	// ClassificationUnknown indicates the classification could not be determined
	ClassificationUnknown Classification = "Unknown"
)

// Classifier handles the quality classification of files using a GenAI engine
type Classifier struct {
	config *config.Config
	llm    llms.Model
}

// New creates a new Classifier with the provided configuration
func New(cfg *config.Config) (*Classifier, error) {
	// Initialize Ollama client
	ollamaOpts := []ollama.Option{
		ollama.WithServerURL(cfg.AIEngine.URL),
		ollama.WithModel(cfg.AIEngine.Model),
	}

	llm, err := ollama.New(ollamaOpts...)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize Ollama client: %w", err)
	}

	return &Classifier{
		config: cfg,
		llm:    llm,
	}, nil
}

// ClassifyContent classifies the content of a file using the GenAI engine
// It returns the classification as one of: Empty, Low quality/low effort, or Good enough
func (c *Classifier) ClassifyContent(content string) (Classification, error) {
	// Early checks for empty content
	if strings.TrimSpace(content) == "" {
		return ClassificationEmpty, nil
	}

	// If this is a mock classifier (used in tests), return the mock classification directly
	if mockLLM, ok := c.llm.(*mockLLM); ok {
		return mockLLM.classification, nil
	}

	ctx := context.Background()

	// Create the prompt by combining the configuration prompt and the file content
	prompt := fmt.Sprintf("%s\n\nHere is the content to review:\n%s",
		c.config.PromptConfig.QualityClassificationPrompt, content)

	// Call the LLM to get a classification
	resp, err := c.llm.Call(ctx, prompt, llms.WithMaxTokens(100))
	if err != nil {
		return ClassificationUnknown, fmt.Errorf("error calling GenAI engine: %w", err)
	}

	// Process the response to extract the classification
	return parseClassification(resp)
}

// parseClassification extracts a classification from the GenAI response
func parseClassification(response string) (Classification, error) {
	// Convert to lowercase and trim spaces for consistent matching
	normalized := strings.ToLower(strings.TrimSpace(response))

	if strings.Contains(normalized, "empty") {
		return ClassificationEmpty, nil
	}

	if strings.Contains(normalized, "low quality") || strings.Contains(normalized, "low effort") {
		return ClassificationLowQuality, nil
	}

	if strings.Contains(normalized, "good enough") {
		return ClassificationGoodEnough, nil
	}

	// If we can't determine a classification, return an error
	return ClassificationUnknown, errors.New("could not determine classification from response")
}

// NewMockClassifier creates a classifier that always returns a predefined classification
// This is useful for testing purposes
func NewMockClassifier(fixedClassification Classification) *Classifier {
	return &Classifier{
		config: nil,
		llm:    &mockLLM{classification: fixedClassification},
	}
}

// mockLLM is a mock implementation of the llms.Model interface for testing
type mockLLM struct {
	classification Classification
}

// Call implements the llms.Model interface for testing
func (m *mockLLM) Call(ctx context.Context, prompt string, options ...llms.CallOption) (string, error) {
	return string(m.classification), nil
}

// GenerateContent implements the llms.Model interface for testing
func (m *mockLLM) GenerateContent(ctx context.Context, messages []llms.MessageContent, options ...llms.CallOption) (*llms.ContentResponse, error) {
	return &llms.ContentResponse{
		Choices: []*llms.ContentChoice{
			{
				Content: string(m.classification),
			},
		},
	}, nil
}
