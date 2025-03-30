package classification

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"ratemykb/config"
	"strings"

	"github.com/tmc/langchaingo/jsonschema"
	"github.com/tmc/langchaingo/llms"
	"github.com/tmc/langchaingo/llms/ollama"
)

// Package classification will handle the quality classification of scanned files

// Classification represents the quality classification of a file
type Classification string

// Classifier handles the quality classification of files using a GenAI engine
type Classifier struct {
	config *config.Config
	llm    llms.Model
}

// New creates a new Classifier with the provided configuration
func New(cfg *config.Config) (*Classifier, error) {
	// Special case for tests: if the model name is "mock-model", use a test classifier
	if cfg.AIEngine.Model == "mock-model" {
		// Create a test LLM that uses simple heuristics
		return &Classifier{
			config: cfg,
			llm:    &testLLM{},
		}, nil
	}

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
// It returns the classification as provided by the LLM
func (c *Classifier) ClassifyContent(content string) (Classification, error) {
	// Early checks for empty content
	if strings.TrimSpace(content) == "" {
		return Classification("Empty"), nil
	}

	// If this is a mock classifier (used in tests), return the mock classification directly
	if mockLLM, ok := c.llm.(*mockLLM); ok {
		return mockLLM.classification, nil
	}

	ctx := context.Background()

	// Create the prompt by replacing the template variable in the configuration prompt
	prompt := strings.Replace(c.config.PromptConfig.QualityClassificationPrompt, "{{ content }}", content, 1)

	// Call the LLM with function calling
	resp, err := c.llm.GenerateContent(ctx,
		[]llms.MessageContent{
			llms.TextParts(llms.ChatMessageTypeHuman, prompt),
		},
		llms.WithFunctions(classificationFunctions),
	)
	if err != nil {
		return Classification("Unknown"), fmt.Errorf("error calling GenAI engine: %w", err)
	}

	// Check if we have a function call response
	if len(resp.Choices) > 0 && resp.Choices[0].FuncCall != nil {
		// print the function call response
		// fmt.Println("Function call response:", resp.Choices[0].FuncCall.Arguments)

		var classificationResponse struct {
			Classification string `json:"classification"`
		}

		err = json.Unmarshal([]byte(resp.Choices[0].FuncCall.Arguments), &classificationResponse)
		if err != nil {
			return Classification("Unknown"), fmt.Errorf("error parsing function call response: %w", err)
		}

		// Use the classification directly from the LLM
		if classificationResponse.Classification != "" {
			return Classification(classificationResponse.Classification), nil
		}
	}

	// If no function call, try to parse from the content directly
	if len(resp.Choices) > 0 && resp.Choices[0].Content != "" {
		// Try to parse the content as JSON
		var classificationResponse struct {
			Classification string `json:"classification"`
		}

		content := resp.Choices[0].Content

		// Clean up the content if it contains markdown code blocks
		content = strings.TrimSpace(content)

		// Remove <think> XML tags section if present (for deepseek model)
		if thinkStart := strings.Index(content, "<think>"); thinkStart != -1 {
			if thinkEnd := strings.Index(content, "</think>"); thinkEnd != -1 {
				beforeThink := content[:thinkStart]
				afterThink := content[thinkEnd+8:] // 8 is the length of "</think>"
				content = beforeThink + afterThink
				content = strings.TrimSpace(content)
			}
		}

		if strings.HasPrefix(content, "```") {
			// Remove markdown code block formatting
			content = strings.TrimPrefix(content, "```json")
			content = strings.TrimPrefix(content, "```")
			content = strings.TrimSuffix(content, "```")
			content = strings.TrimSpace(content)
		}

		err := json.Unmarshal([]byte(content), &classificationResponse)
		if err == nil && classificationResponse.Classification != "" {
			// Successfully parsed JSON, use the classification
			return Classification(classificationResponse.Classification), nil
		} else {
			// print the error
			fmt.Println("Error parsing JSON:", err)
		}

		// If not valid JSON or missing classification, use the raw content
		return Classification(strings.TrimSpace(content)), nil
	}

	return Classification("Unknown"), errors.New("no valid response from GenAI engine")
}

// Define the classification function for the LLM
var classificationFunctions = []llms.FunctionDefinition{
	{
		Name:        "classifyContent",
		Description: "Classify the quality of content",
		Parameters: jsonschema.Definition{
			Type: jsonschema.Object,
			Properties: map[string]jsonschema.Definition{
				"classification": {
					Type:        jsonschema.String,
					Description: "The classification of the content describing its quality",
				},
			},
			Required: []string{"classification"},
		},
	},
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
	// Create a JSON string with the classification
	args := fmt.Sprintf(`{"classification": "%s"}`, m.classification)

	return &llms.ContentResponse{
		Choices: []*llms.ContentChoice{
			{
				Content: string(m.classification),
				FuncCall: &llms.FunctionCall{
					Name:      "classifyContent",
					Arguments: args,
				},
			},
		},
	}, nil
}

// testLLM is a test implementation of the llms.Model interface for simple classification
type testLLM struct{}

// Call implements the llms.Model interface for testing
func (m *testLLM) Call(ctx context.Context, prompt string, options ...llms.CallOption) (string, error) {
	// Extract content from the prompt
	contentIndex := strings.Index(prompt, "Here is the content to review:")
	if contentIndex == -1 {
		return string(Classification("Unknown")), nil
	}

	content := prompt[contentIndex+len("Here is the content to review:"):]

	// Simple classification logic for tests
	content = strings.TrimSpace(content)
	if content == "" {
		return string(Classification("Empty")), nil
	}

	if len(content) < 100 || strings.Contains(content, "TODO") {
		return string(Classification("Low quality")), nil
	}

	return string(Classification("Good enough")), nil
}

// GenerateContent implements the llms.Model interface for testing
func (m *testLLM) GenerateContent(ctx context.Context, messages []llms.MessageContent, options ...llms.CallOption) (*llms.ContentResponse, error) {
	// Extract the prompt from the messages
	var prompt string
	if len(messages) > 0 {
		var parts []string
		for _, part := range messages[0].Parts {
			if textPart, ok := part.(llms.TextContent); ok {
				parts = append(parts, textPart.Text)
			}
		}
		prompt = strings.Join(parts, "")
	}

	// Extract content from the prompt
	contentIndex := strings.Index(prompt, "Here is the content to review:")
	if contentIndex == -1 {
		return simpleResponse(Classification("Unknown")), nil
	}

	content := prompt[contentIndex+len("Here is the content to review:"):]

	// Simple classification logic for tests
	content = strings.TrimSpace(content)
	if content == "" {
		return simpleResponse(Classification("Empty")), nil
	}

	if len(content) < 100 || strings.Contains(content, "TODO") {
		return simpleResponse(Classification("Low quality")), nil
	}

	return simpleResponse(Classification("Good enough")), nil
}

// simpleResponse creates a ContentResponse with both regular content and function call
func simpleResponse(classification Classification) *llms.ContentResponse {
	args := fmt.Sprintf(`{"classification": "%s"}`, classification)

	return &llms.ContentResponse{
		Choices: []*llms.ContentChoice{
			{
				Content: string(classification),
				FuncCall: &llms.FunctionCall{
					Name:      "classifyContent",
					Arguments: args,
				},
			},
		},
	}
}
