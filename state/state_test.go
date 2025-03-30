package state

import (
	"os"
	"path/filepath"
	"testing"

	"ratemykb/classification"
	"ratemykb/output"
	"ratemykb/scanner"
)

func TestNew(t *testing.T) {
	// Create a temporary directory
	tempDir, err := os.MkdirTemp("", "state-test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Create a new state
	state, err := New(tempDir)
	if err != nil {
		t.Fatalf("Failed to create state: %v", err)
	}

	// Check that the state was created correctly
	if state.TargetFolder != tempDir {
		t.Errorf("Expected target folder %s, got %s", tempDir, state.TargetFolder)
	}

	if state.ReportPath != filepath.Join(tempDir, "vault-quality-report.md") {
		t.Errorf("Expected report path %s, got %s", filepath.Join(tempDir, "vault-quality-report.md"), state.ReportPath)
	}

	if len(state.ProcessedFiles) != 0 {
		t.Errorf("Expected 0 processed files, got %d", len(state.ProcessedFiles))
	}
}

func TestIsFileProcessed(t *testing.T) {
	// Create a temporary directory
	tempDir, err := os.MkdirTemp("", "state-test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Create a new state
	state, err := New(tempDir)
	if err != nil {
		t.Fatalf("Failed to create state: %v", err)
	}

	// Check that a file is not processed
	filePath := filepath.Join(tempDir, "test.md")
	if state.IsFileProcessed(filePath) {
		t.Errorf("Expected file %s to not be processed", filePath)
	}

	// Add a processed file
	state.ProcessedFiles[filePath] = output.ResultFile{
		Path:           filePath,
		Status:         scanner.StatusNeedsReview,
		Classification: classification.Classification("Good enough"),
	}

	// Check that the file is now processed
	if !state.IsFileProcessed(filePath) {
		t.Errorf("Expected file %s to be processed", filePath)
	}
}

func TestAddProcessedFile(t *testing.T) {
	// Create a temporary directory
	tempDir, err := os.MkdirTemp("", "state-test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Create a new state
	state, err := New(tempDir)
	if err != nil {
		t.Fatalf("Failed to create state: %v", err)
	}

	// Add a processed file
	filePath := filepath.Join(tempDir, "test.md")
	result := output.ResultFile{
		Path:           filePath,
		Status:         scanner.StatusNeedsReview,
		Classification: classification.Classification("Good enough"),
	}

	err = state.AddProcessedFile(result)
	if err != nil {
		t.Fatalf("Failed to add processed file: %v", err)
	}

	// Check that the file was added
	if !state.IsFileProcessed(filePath) {
		t.Errorf("Expected file %s to be processed", filePath)
	}

	// Check that the report was created
	if _, err := os.Stat(state.ReportPath); os.IsNotExist(err) {
		t.Errorf("Expected report file %s to exist", state.ReportPath)
	}
}

func TestLoadExistingReport(t *testing.T) {
	// Create a temporary directory
	tempDir, err := os.MkdirTemp("", "state-test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Create a test report file
	reportPath := filepath.Join(tempDir, "vault-quality-report.md")
	reportContent := `# Vault Quality Report

Generated on: 2023-01-01 12:00:00

Target folder: ` + "`" + tempDir + "`" + `

## Statistics

- Total files processed: 3
- Empty files: 1
- Files with frontmatter only: 1
- Good enough files: 1

## Empty Files

- [[empty-file]]

## Files with Frontmatter Only

- [[frontmatter-only]]

## Good enough Files

- [[good-file]]
`

	err = os.WriteFile(reportPath, []byte(reportContent), 0644)
	if err != nil {
		t.Fatalf("Failed to create test report: %v", err)
	}

	// Create a new state, which should load the existing report
	state, err := New(tempDir)
	if err != nil {
		t.Fatalf("Failed to create state: %v", err)
	}

	// Check that the processed files were loaded
	if len(state.ProcessedFiles) != 3 {
		t.Errorf("Expected 3 processed files, got %d", len(state.ProcessedFiles))
	}

	// Check specific files
	emptyFilePath := filepath.Join(tempDir, "empty-file.md")
	if !state.IsFileProcessed(emptyFilePath) {
		t.Errorf("Expected file %s to be processed", emptyFilePath)
	}

	frontmatterFilePath := filepath.Join(tempDir, "frontmatter-only.md")
	if !state.IsFileProcessed(frontmatterFilePath) {
		t.Errorf("Expected file %s to be processed", frontmatterFilePath)
	}

	goodFilePath := filepath.Join(tempDir, "good-file.md")
	if !state.IsFileProcessed(goodFilePath) {
		t.Errorf("Expected file %s to be processed", goodFilePath)
	}

	// Check classifications
	if state.ProcessedFiles[emptyFilePath].Classification != classification.Classification("Empty") {
		t.Errorf("Expected classification Empty, got %s", state.ProcessedFiles[emptyFilePath].Classification)
	}

	if state.ProcessedFiles[frontmatterFilePath].Classification != classification.Classification("Low quality") {
		t.Errorf("Expected classification Low quality, got %s", state.ProcessedFiles[frontmatterFilePath].Classification)
	}

	if state.ProcessedFiles[goodFilePath].Classification != classification.Classification("Good enough") {
		t.Errorf("Expected classification Good enough, got %s", state.ProcessedFiles[goodFilePath].Classification)
	}
}