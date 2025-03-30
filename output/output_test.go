package output

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"ratemykb/classification"
	"ratemykb/scanner"
)

func TestFormatObsidianLink(t *testing.T) {
	tests := []struct {
		name         string
		targetFolder string
		filePath     string
		expected     string
	}{
		{
			name:         "basic path",
			targetFolder: "/root",
			filePath:     "/root/file.md",
			expected:     "[[file]]",
		},
		{
			name:         "nested path",
			targetFolder: "/root",
			filePath:     "/root/folder/subfolder/file.md",
			expected:     "[[folder/subfolder/file]]",
		},
		{
			name:         "with spaces",
			targetFolder: "/root",
			filePath:     "/root/folder with spaces/file with spaces.md",
			expected:     "[[folder with spaces/file with spaces]]",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			generator := New(tc.targetFolder)
			result := generator.formatObsidianLink(tc.filePath)
			if result != tc.expected {
				t.Errorf("expected %s, got %s", tc.expected, result)
			}
		})
	}
}

func TestCreateReport(t *testing.T) {
	// Create a temporary directory for the test
	tempDir, err := os.MkdirTemp("", "output-test-")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Create test data
	files := []ResultFile{
		{
			Path:   filepath.Join(tempDir, "empty.md"),
			Status: scanner.StatusEmpty,
		},
		{
			Path:   filepath.Join(tempDir, "frontmatter-only.md"),
			Status: scanner.StatusFrontmatterOnly,
		},
		{
			Path:           filepath.Join(tempDir, "low-quality.md"),
			Status:         scanner.StatusNeedsReview,
			Classification: classification.Classification("Low quality"),
		},
		{
			Path:           filepath.Join(tempDir, "good-enough.md"),
			Status:         scanner.StatusNeedsReview,
			Classification: classification.Classification("Good enough"),
		},
	}

	// Create the report
	generator := New(tempDir)
	err = generator.CreateReport(files)
	if err != nil {
		t.Fatalf("CreateReport returned error: %v", err)
	}

	// Check that the report file exists
	reportPath := filepath.Join(tempDir, "vault-quality-report.md")
	if _, err := os.Stat(reportPath); os.IsNotExist(err) {
		t.Fatalf("report file was not created")
	}

	// Read the report content
	content, err := os.ReadFile(reportPath)
	if err != nil {
		t.Fatalf("failed to read report: %v", err)
	}

	// Verify content
	contentStr := string(content)

	// Check section headers
	if !strings.Contains(contentStr, "# Vault Quality Report") {
		t.Error("report missing main header")
	}
	if !strings.Contains(contentStr, "## Empty Files") {
		t.Error("report missing empty files section")
	}
	if !strings.Contains(contentStr, "## Files with Frontmatter Only") {
		t.Error("report missing frontmatter-only files section")
	}
	if !strings.Contains(contentStr, "## Low quality Files") {
		t.Error("report missing low quality files section")
	}
	if !strings.Contains(contentStr, "## Good enough Files") {
		t.Error("report missing good enough files section")
	}

	// Check statistics
	if !strings.Contains(contentStr, "Total files scanned: 4") {
		t.Error("report missing or incorrect total files count")
	}
	if !strings.Contains(contentStr, "Empty files: 1") {
		t.Error("report missing or incorrect empty files count")
	}
	if !strings.Contains(contentStr, "Files with frontmatter only: 1") {
		t.Error("report missing or incorrect frontmatter-only files count")
	}
	if !strings.Contains(contentStr, "Low quality files: 1") {
		t.Error("report missing or incorrect low quality files count")
	}
	if !strings.Contains(contentStr, "Good enough files: 1") {
		t.Error("report missing or incorrect good enough files count")
	}

	// Check file links
	if !strings.Contains(contentStr, "[[empty]]") {
		t.Error("report missing empty file link")
	}
	if !strings.Contains(contentStr, "[[frontmatter-only]]") {
		t.Error("report missing frontmatter-only file link")
	}
	if !strings.Contains(contentStr, "[[low-quality]]") {
		t.Error("report missing low quality file link")
	}
	if !strings.Contains(contentStr, "[[good-enough]]") {
		t.Error("report missing good enough file link")
	}
}

func TestEmptySections(t *testing.T) {
	// Create a temporary directory for the test
	tempDir, err := os.MkdirTemp("", "output-test-empty-")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Create a generator with empty files
	generator := New(tempDir)
	err = generator.CreateReport([]ResultFile{})
	if err != nil {
		t.Fatalf("CreateReport returned error: %v", err)
	}

	// Read the report content
	reportPath := filepath.Join(tempDir, "vault-quality-report.md")
	content, err := os.ReadFile(reportPath)
	if err != nil {
		t.Fatalf("failed to read report: %v", err)
	}

	// Verify empty sections
	contentStr := string(content)
	if !strings.Contains(contentStr, "No empty files found.") {
		t.Error("report missing empty files message")
	}
	if !strings.Contains(contentStr, "No files with frontmatter only found.") {
		t.Error("report missing empty frontmatter-only files message")
	}
}
