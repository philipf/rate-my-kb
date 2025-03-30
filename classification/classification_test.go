package classification

import (
	"context"
	"ratemykb/config"
	"strings"
	"testing"

	"github.com/tmc/langchaingo/llms"
)

func TestClassifyContent_WithMockClassifier(t *testing.T) {
	tests := []struct {
		name           string
		content        string
		mockClassifier Classification
		want           Classification
		wantErr        bool
	}{
		{
			name:           "Empty content",
			content:        "",
			mockClassifier: Classification("Empty"),
			want:           Classification("Empty"),
			wantErr:        false,
		},
		{
			name:           "Empty content with whitespace",
			content:        "   \n   \t   ",
			mockClassifier: Classification("Empty"),
			want:           Classification("Empty"),
			wantErr:        false,
		},
		{
			name:           "Low quality content",
			content:        "This is some content",
			mockClassifier: Classification("Low quality"),
			want:           Classification("Low quality"),
			wantErr:        false,
		},
		{
			name:           "Good enough content",
			content:        "This is some good enough content",
			mockClassifier: Classification("Good enough"),
			want:           Classification("Good enough"),
			wantErr:        false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			classifier := NewMockClassifier(tt.mockClassifier)

			got, err := classifier.ClassifyContent(tt.content)
			if (err != nil) != tt.wantErr {
				t.Errorf("ClassifyContent() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			// For empty content, the function should return early without calling the LLM
			if tt.content == "" || len(strings.TrimSpace(tt.content)) == 0 {
				if got != Classification("Empty") {
					t.Errorf("ClassifyContent() with empty content = %v, want %v", got, Classification("Empty"))
				}
				return
			}

			if got != tt.want {
				t.Errorf("ClassifyContent() = %v, want %v", got, tt.want)
			}
		})
	}
}

// mixedResponseLLM is a mock LLM that returns responses with text surrounding JSON content
type mixedResponseLLM struct {
	classification string
	responseType   string
}

// Call implements the llms.Model interface
func (m *mixedResponseLLM) Call(ctx context.Context, prompt string, options ...llms.CallOption) (string, error) {
	return "", nil // Not used in this test
}

// GenerateContent implements the llms.Model interface
func (m *mixedResponseLLM) GenerateContent(ctx context.Context, messages []llms.MessageContent, options ...llms.CallOption) (*llms.ContentResponse, error) {
	var content string
	
	switch m.responseType {
	case "text_before_json":
		content = "The content provides specific information about a Machine Learning Guru and suggests watching certain videos, indicating substance without excessive detail. It's clear and informative.\n\n```json\n{\n  \"classification\": \"" + m.classification + "\"\n}\n```"
	case "text_after_json":
		content = "```json\n{\n  \"classification\": \"" + m.classification + "\"\n}\n```\n\nThis classification was determined based on the content's structure and information density."
	case "text_surrounding_json":
		content = "Analysis: The content is well-structured.\n\n{\n  \"classification\": \"" + m.classification + "\"\n}\n\nAdditional notes: The formatting could be improved."
	default:
		content = "{\n  \"classification\": \"" + m.classification + "\"\n}"
	}
	
	return &llms.ContentResponse{
		Choices: []*llms.ContentChoice{
			{
				Content: content,
				// No function call in this case
			},
		},
	}, nil
}

// TestJSONExtractionFromMixedContent tests the ability to extract JSON from responses with additional text
func TestJSONExtractionFromMixedContent(t *testing.T) {
	tests := []struct {
		name         string
		responseType string
		expected     Classification
	}{
		{
			name:         "Text before JSON",
			responseType: "text_before_json",
			expected:     Classification("Good enough"),
		},
		{
			name:         "Text after JSON",
			responseType: "text_after_json",
			expected:     Classification("Good enough"),
		},
		{
			name:         "Text surrounding JSON",
			responseType: "text_surrounding_json",
			expected:     Classification("Good enough"),
		},
		{
			name:         "Clean JSON only",
			responseType: "clean_json",
			expected:     Classification("Good enough"),
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a minimal config for testing
			cfg := &config.Config{
				PromptConfig: config.PromptConfig{
					QualityClassificationPrompt: "Here is the content to review: {{ content }}",
				},
			}
			
			// Create a classifier with our custom mock LLM
			classifier := &Classifier{
				config: cfg,
				llm:    &mixedResponseLLM{classification: "Good enough", responseType: tt.responseType},
			}
			
			// Test with some non-empty content
			got, err := classifier.ClassifyContent("Some test content")
			
			if err != nil {
				t.Errorf("ClassifyContent() error = %v, expected no error", err)
				return
			}
			
			if got != tt.expected {
				t.Errorf("ClassifyContent() = %v, want %v", got, tt.expected)
			}
		})
	}
}
