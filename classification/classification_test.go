package classification

import (
	"strings"
	"testing"
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
