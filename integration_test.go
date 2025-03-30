package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"ratemykb/cli"
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

	// Save the original args and restore them later
	oldArgs := os.Args
	defer func() { os.Args = oldArgs }()

	// Set up command-line arguments for the test
	os.Args = []string{
		"ratemykb",
		"--target", tempDir,
		"--config", configPath,
	}

	// Run the CLI command (we'll use a separate function instead of main() directly)
	err = runCLI()
	if err != nil {
		t.Fatalf("Failed to run CLI: %v", err)
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

// runCLI executes the CLI command
func runCLI() error {
	// Call the Execute function from the cli package
	cli.Execute()
	return nil
}
