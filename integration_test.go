package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"ratemykb/classification"
	"ratemykb/config"
	"ratemykb/output"
	"ratemykb/scanner"
	"ratemykb/state"
)

// TestIntegration is an end-to-end test of the entire application
// It creates a temporary folder with sample Markdown files, runs the application,
// and checks that the output report is generated correctly
func TestIntegration(t *testing.T) {
	// Create a temporary directory
	tempDir, err := os.MkdirTemp("", "ratemykb-test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Create sample Markdown files
	createSampleFiles(t, tempDir)

	// Create a mock config file
	configPath := filepath.Join(tempDir, "config.yaml")
	createMockConfig(t, configPath)

	// Load configuration
	cfg, err := config.LoadConfig(configPath)
	if err != nil {
		t.Fatalf("Failed to load configuration: %v", err)
	}

	// Initialize state manager
	stateManager, err := state.New(tempDir)
	if err != nil {
		t.Fatalf("Failed to initialize state manager: %v", err)
	}

	// Initialize scanner
	fileScanner, err := scanner.New(cfg)
	if err != nil {
		t.Fatalf("Failed to initialize scanner: %v", err)
	}

	// Scan the target folder
	files, err := fileScanner.ScanDirectory(tempDir)
	if err != nil {
		t.Fatalf("Failed to scan directory: %v", err)
	}

	// Initialize classifier
	classifier, err := classification.New(cfg)
	if err != nil {
		t.Fatalf("Failed to initialize classifier: %v", err)
	}

	// Process each file
	for _, file := range files {
		// Create a result file with default classification
		result := output.ResultFile{
			Path:           file.Path,
			Status:         file.Status,
			Classification: classification.Classification("Unknown"),
		}

		// Classify files that need review
		if file.Status == scanner.StatusNeedsReview {
			// Read the content of the file
			content, err := scanner.ReadFileContent(file.Path)
			if err != nil {
				t.Fatalf("Failed to read file content: %v", err)
			}

			// Classify the content
			result.Classification, err = classifier.ClassifyContent(content)
			if err != nil {
				t.Fatalf("Failed to classify content: %v", err)
			}
		} else if file.Status == scanner.StatusEmpty {
			// Map scanner status to classification
			result.Classification = classification.Classification("Empty")
		} else if file.Status == scanner.StatusFrontmatterOnly {
			// Frontmatter-only files are considered low quality
			result.Classification = classification.Classification("Low quality")
		} else if file.Status == scanner.StatusExcluded {
			// Skip excluded files
			continue
		}

		// Add processed file to state and update report
		if err := stateManager.AddProcessedFile(result); err != nil {
			t.Fatalf("Failed to add processed file: %v", err)
		}
	}

	// Check that the report was generated
	reportPath := filepath.Join(tempDir, "vault-quality-report.md")
	if _, err := os.Stat(reportPath); os.IsNotExist(err) {
		t.Fatalf("Report file was not generated")
	}

	// Read the report content
	reportContent, err := os.ReadFile(reportPath)
	if err != nil {
		t.Fatalf("Failed to read report: %v", err)
	}

	// Check that the report contains the expected sections
	report := string(reportContent)
	baseSections := []string{
		"# Vault Quality Report",
		"## Statistics",
		"## Empty Files",
		"## Files with Frontmatter Only",
	}

	for _, section := range baseSections {
		if !strings.Contains(report, section) {
			t.Errorf("Report missing section: %s", section)
		}
	}

	// Check that each file type is correctly categorized
	checkReportContains(t, report, "Empty Files", "[[empty-file]]")
	checkReportContains(t, report, "Files with Frontmatter Only", "[[frontmatter-only]]")

	// Test resumable processing by adding a new file
	newFile := filepath.Join(tempDir, "new-file.md")
	newContent := `---
title: New File
---

This is a new file added after the initial run.
`
	if err := os.WriteFile(newFile, []byte(newContent), 0644); err != nil {
		t.Fatalf("Failed to create new file: %v", err)
	}

	// Scan the directory again
	files, err = fileScanner.ScanDirectory(tempDir)
	if err != nil {
		t.Fatalf("Failed to scan directory: %v", err)
	}

	// Process only the new file
	for _, file := range files {
		// Skip already processed files
		if stateManager.IsFileProcessed(file.Path) {
			continue
		}

		// Create a result file with default classification
		result := output.ResultFile{
			Path:           file.Path,
			Status:         file.Status,
			Classification: classification.Classification("Unknown"),
		}

		// Classify files that need review
		if file.Status == scanner.StatusNeedsReview {
			// Read the content of the file
			content, err := scanner.ReadFileContent(file.Path)
			if err != nil {
				t.Fatalf("Failed to read file content: %v", err)
			}

			// Classify the content
			result.Classification, err = classifier.ClassifyContent(content)
			if err != nil {
				t.Fatalf("Failed to classify content: %v", err)
			}
		} else if file.Status == scanner.StatusEmpty {
			// Map scanner status to classification
			result.Classification = classification.Classification("Empty")
		} else if file.Status == scanner.StatusFrontmatterOnly {
			// Frontmatter-only files are considered low quality
			result.Classification = classification.Classification("Low quality")
		} else if file.Status == scanner.StatusExcluded {
			// Skip excluded files
			continue
		}

		// Add processed file to state and update report
		if err := stateManager.AddProcessedFile(result); err != nil {
			t.Fatalf("Failed to add processed file: %v", err)
		}
	}

	// Read the updated report
	reportContent, err = os.ReadFile(reportPath)
	if err != nil {
		t.Fatalf("Failed to read updated report: %v", err)
	}

	// Check that the new file is included in the report
	report = string(reportContent)
	if !strings.Contains(report, "[[new-file]]") {
		t.Errorf("Updated report does not contain the new file")
	}

	// Check that the original files are still in the report
	checkReportContains(t, report, "Empty Files", "[[empty-file]]")
	checkReportContains(t, report, "Files with Frontmatter Only", "[[frontmatter-only]]")
}

// TestIncrementalProcessing tests that files are processed incrementally
// and that the report is updated after each file
func TestIncrementalProcessing(t *testing.T) {
	// Create a temporary directory
	tempDir, err := os.MkdirTemp("", "ratemykb-incremental-test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Create a mock config file
	configPath := filepath.Join(tempDir, "config.yaml")
	createMockConfig(t, configPath)

	// Load configuration
	cfg, err := config.LoadConfig(configPath)
	if err != nil {
		t.Fatalf("Failed to load configuration: %v", err)
	}

	// Initialize state manager
	stateManager, err := state.New(tempDir)
	if err != nil {
		t.Fatalf("Failed to initialize state manager: %v", err)
	}

	// Initialize scanner
	fileScanner, err := scanner.New(cfg)
	if err != nil {
		t.Fatalf("Failed to initialize scanner: %v", err)
	}

	// Create the first file
	emptyFile := filepath.Join(tempDir, "empty-file.md")
	if err := os.WriteFile(emptyFile, []byte(""), 0644); err != nil {
		t.Fatalf("Failed to create empty file: %v", err)
	}

	// Scan the directory
	files, err := fileScanner.ScanDirectory(tempDir)
	if err != nil {
		t.Fatalf("Failed to scan directory: %v", err)
	}

	// Process the first file
	for _, file := range files {
		// Create a result file
		result := output.ResultFile{
			Path:           file.Path,
			Status:         file.Status,
			Classification: classification.Classification("Empty"),
		}

		// Add processed file to state and update report
		if err := stateManager.AddProcessedFile(result); err != nil {
			t.Fatalf("Failed to add processed file: %v", err)
		}
	}

	// Check that the report was generated and contains the empty file
	reportPath := filepath.Join(tempDir, "vault-quality-report.md")
	reportContent, err := os.ReadFile(reportPath)
	if err != nil {
		t.Fatalf("Failed to read report: %v", err)
	}
	report := string(reportContent)
	checkReportContains(t, report, "Empty Files", "[[empty-file]]")

	// Create a second file
	frontmatterFile := filepath.Join(tempDir, "frontmatter-only.md")
	frontmatterContent := `---
title: Test
date: 2023-01-01
---
`
	if err := os.WriteFile(frontmatterFile, []byte(frontmatterContent), 0644); err != nil {
		t.Fatalf("Failed to create frontmatter-only file: %v", err)
	}

	// Scan the directory again
	files, err = fileScanner.ScanDirectory(tempDir)
	if err != nil {
		t.Fatalf("Failed to scan directory: %v", err)
	}

	// Process only the new file
	for _, file := range files {
		// Skip already processed files
		if stateManager.IsFileProcessed(file.Path) {
			continue
		}

		// Create a result file
		result := output.ResultFile{
			Path:           file.Path,
			Status:         file.Status,
			Classification: classification.Classification("Low quality"),
		}

		// Add processed file to state and update report
		if err := stateManager.AddProcessedFile(result); err != nil {
			t.Fatalf("Failed to add processed file: %v", err)
		}
	}

	// Check that the report was updated and contains both files
	reportContent, err = os.ReadFile(reportPath)
	if err != nil {
		t.Fatalf("Failed to read updated report: %v", err)
	}
	report = string(reportContent)
	checkReportContains(t, report, "Empty Files", "[[empty-file]]")
	checkReportContains(t, report, "Files with Frontmatter Only", "[[frontmatter-only]]")
}

// createSampleFiles creates test Markdown files in the temporary directory
func createSampleFiles(t *testing.T, dir string) {
	// Create a subdirectory
	subDir := filepath.Join(dir, "notes")
	if err := os.Mkdir(subDir, 0755); err != nil {
		t.Fatalf("Failed to create subdirectory: %v", err)
	}

	// 1. Empty file
	emptyFile := filepath.Join(dir, "empty-file.md")
	if err := os.WriteFile(emptyFile, []byte(""), 0644); err != nil {
		t.Fatalf("Failed to create empty file: %v", err)
	}

	// 2. File with frontmatter only
	frontmatterOnly := filepath.Join(dir, "frontmatter-only.md")
	frontmatterContent := `---
title: Test
date: 2023-01-01
---
`
	if err := os.WriteFile(frontmatterOnly, []byte(frontmatterContent), 0644); err != nil {
		t.Fatalf("Failed to create frontmatter-only file: %v", err)
	}

	// 3. Good enough file
	goodFile := filepath.Join(dir, "good-file.md")
	goodContent := `---
title: Good File
date: 2023-01-01
---

# Good File

This is a good file with plenty of content.

## Section 1

Lorem ipsum dolor sit amet, consectetur adipiscing elit.

## Section 2

More content here that makes this a decent file.
`
	if err := os.WriteFile(goodFile, []byte(goodContent), 0644); err != nil {
		t.Fatalf("Failed to create good file: %v", err)
	}

	// 4. Low quality file
	lowQualityFile := filepath.Join(subDir, "low-quality.md")
	lowQualityContent := `---
title: Low Quality
---

TODO: Add content later
`
	if err := os.WriteFile(lowQualityFile, []byte(lowQualityContent), 0644); err != nil {
		t.Fatalf("Failed to create low quality file: %v", err)
	}

	// 5. Exclusion list file
	exclusionFile := filepath.Join(dir, "quality_exclude_links.md")
	exclusionContent := `# Files to Exclude

The following files should be excluded from quality checks:

- [[good-file]]
`
	if err := os.WriteFile(exclusionFile, []byte(exclusionContent), 0644); err != nil {
		t.Fatalf("Failed to create exclusion file: %v", err)
	}
}

// createMockConfig creates a mock configuration file for testing
func createMockConfig(t *testing.T, path string) {
	config := `ai_engine:
  url: "http://localhost:11434/"
  model: "mock-model" # We'll use a mock classifier in tests

scan_settings:
  file_extension: ".md"
  exclude_directories: []

prompt_config:
  quality_classification_prompt: "Review the content and determine if it's: 'Empty', 'Low quality/low effort', or 'Good enough'."

exclusion_file:
  path: "quality_exclude_links.md"
`
	if err := os.WriteFile(path, []byte(config), 0644); err != nil {
		t.Fatalf("Failed to create mock config: %v", err)
	}
}

// checkReportContains checks if a specific section of the report contains the expected text
func checkReportContains(t *testing.T, report, section, expectedText string) {
	sectionStart := strings.Index(report, section)
	if sectionStart == -1 {
		t.Fatalf("Section not found: %s", section)
	}

	// Find the next section or end of file
	nextSectionStart := len(report)
	for _, s := range []string{"## Empty Files", "## Files with Frontmatter Only", "## "} {
		if s != section {
			idx := strings.Index(report[sectionStart:], s)
			if idx != -1 && sectionStart+idx > sectionStart {
				idx += sectionStart
				if idx < nextSectionStart {
					nextSectionStart = idx
				}
			}
		}
	}

	// Check if the section contains the expected text
	sectionContent := report[sectionStart:nextSectionStart]
	if !strings.Contains(sectionContent, expectedText) {
		t.Errorf("Section %s does not contain expected text: %s", section, expectedText)
	}
}
