package classification

import (
	"strings"
	"testing"
)

func TestParseClassification(t *testing.T) {
	tests := []struct {
		name     string
		response string
		want     Classification
		wantErr  bool
	}{
		{
			name:     "Empty classification",
			response: "This content is empty",
			want:     ClassificationEmpty,
			wantErr:  false,
		},
		{
			name:     "Low quality classification",
			response: "This content is of low quality",
			want:     ClassificationLowQuality,
			wantErr:  false,
		},
		{
			name:     "Low effort classification",
			response: "This is a low effort document",
			want:     ClassificationLowQuality,
			wantErr:  false,
		},
		{
			name:     "Good enough classification",
			response: "This content is good enough",
			want:     ClassificationGoodEnough,
			wantErr:  false,
		},
		{
			name:     "Unknown classification",
			response: "This is some other response that doesn't match any classification",
			want:     ClassificationUnknown,
			wantErr:  true,
		},
		{
			name:     "Case insensitive",
			response: "EMPTY",
			want:     ClassificationEmpty,
			wantErr:  false,
		},
		{
			name:     "With extra text",
			response: "I think this document is Good Enough based on its structure and content.",
			want:     ClassificationGoodEnough,
			wantErr:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parseClassification(tt.response)
			if (err != nil) != tt.wantErr {
				t.Errorf("parseClassification() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("parseClassification() = %v, want %v", got, tt.want)
			}
		})
	}
}

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
			mockClassifier: ClassificationEmpty,
			want:           ClassificationEmpty,
			wantErr:        false,
		},
		{
			name:           "Empty content with whitespace",
			content:        "   \n   \t   ",
			mockClassifier: ClassificationEmpty,
			want:           ClassificationEmpty,
			wantErr:        false,
		},
		{
			name:           "Low quality content",
			content:        "This is some content",
			mockClassifier: ClassificationLowQuality,
			want:           ClassificationLowQuality,
			wantErr:        false,
		},
		{
			name:           "Good enough content",
			content:        "This is some good enough content",
			mockClassifier: ClassificationGoodEnough,
			want:           ClassificationGoodEnough,
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
				if got != ClassificationEmpty {
					t.Errorf("ClassifyContent() with empty content = %v, want %v", got, ClassificationEmpty)
				}
				return
			}

			if got != tt.want {
				t.Errorf("ClassifyContent() = %v, want %v", got, tt.want)
			}
		})
	}
}
