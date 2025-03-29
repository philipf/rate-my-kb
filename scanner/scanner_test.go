package scanner

import (
	"os"
	"path/filepath"
	"testing"

	"ratemykb/config"
)

func TestScannerNew(t *testing.T) {
	cfg := config.GetDefaultConfig()

	scanner, err := New(cfg)
	if err != nil {
		t.Fatalf("Failed to create scanner: %v", err)
	}

	if scanner == nil {
		t.Fatal("Expected scanner to be non-nil")
	}

	if scanner.config != cfg {
		t.Errorf("Expected scanner config to be the same as provided config")
	}
}

func TestEmptyFileCheck(t *testing.T) {
	// Create a temporary directory
	tempDir, err := os.MkdirTemp("", "scanner-test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Create an empty file
	emptyFilePath := filepath.Join(tempDir, "empty.md")
	if err := os.WriteFile(emptyFilePath, []byte(""), 0644); err != nil {
		t.Fatalf("Failed to create empty file: %v", err)
	}

	// Create a file with only whitespace
	whitespaceFilePath := filepath.Join(tempDir, "whitespace.md")
	if err := os.WriteFile(whitespaceFilePath, []byte("  \n  \t  \n"), 0644); err != nil {
		t.Fatalf("Failed to create whitespace file: %v", err)
	}

	// Create a file with content
	contentFilePath := filepath.Join(tempDir, "content.md")
	if err := os.WriteFile(contentFilePath, []byte("# This is a heading\n\nThis is content."), 0644); err != nil {
		t.Fatalf("Failed to create content file: %v", err)
	}

	// Create scanner with default config
	cfg := config.GetDefaultConfig()
	scanner, err := New(cfg)
	if err != nil {
		t.Fatalf("Failed to create scanner: %v", err)
	}

	// Test empty file
	status, err := scanner.checkFileStatus(emptyFilePath)
	if err != nil {
		t.Errorf("Failed to check empty file status: %v", err)
	}
	if status != StatusEmpty {
		t.Errorf("Expected empty file status to be %s, got %s", StatusEmpty, status)
	}

	// Test whitespace file
	status, err = scanner.checkFileStatus(whitespaceFilePath)
	if err != nil {
		t.Errorf("Failed to check whitespace file status: %v", err)
	}
	if status != StatusEmpty {
		t.Errorf("Expected whitespace file status to be %s, got %s", StatusEmpty, status)
	}

	// Test content file
	status, err = scanner.checkFileStatus(contentFilePath)
	if err != nil {
		t.Errorf("Failed to check content file status: %v", err)
	}
	if status != StatusNeedsReview {
		t.Errorf("Expected content file status to be %s, got %s", StatusNeedsReview, status)
	}
}

func TestFrontmatterOnlyCheck(t *testing.T) {
	// Create a temporary directory
	tempDir, err := os.MkdirTemp("", "scanner-test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Create a file with only frontmatter
	frontmatterOnlyPath := filepath.Join(tempDir, "frontmatter_only.md")
	frontmatterContent := "---\ntitle: Test\ndate: 2023-01-01\n---\n"
	if err := os.WriteFile(frontmatterOnlyPath, []byte(frontmatterContent), 0644); err != nil {
		t.Fatalf("Failed to create frontmatter-only file: %v", err)
	}

	// Create a file with frontmatter and content
	frontmatterAndContentPath := filepath.Join(tempDir, "frontmatter_and_content.md")
	mixedContent := "---\ntitle: Test\ndate: 2023-01-01\n---\n\n# Heading\n\nContent here."
	if err := os.WriteFile(frontmatterAndContentPath, []byte(mixedContent), 0644); err != nil {
		t.Fatalf("Failed to create frontmatter-and-content file: %v", err)
	}

	// Create a file with invalid frontmatter (missing end marker)
	invalidFrontmatterPath := filepath.Join(tempDir, "invalid_frontmatter.md")
	invalidContent := "---\ntitle: Test\ndate: 2023-01-01\n"
	if err := os.WriteFile(invalidFrontmatterPath, []byte(invalidContent), 0644); err != nil {
		t.Fatalf("Failed to create invalid-frontmatter file: %v", err)
	}

	// Create scanner with default config
	cfg := config.GetDefaultConfig()
	scanner, err := New(cfg)
	if err != nil {
		t.Fatalf("Failed to create scanner: %v", err)
	}

	// Test frontmatter-only file
	status, err := scanner.checkFileStatus(frontmatterOnlyPath)
	if err != nil {
		t.Errorf("Failed to check frontmatter-only file status: %v", err)
	}
	if status != StatusFrontmatterOnly {
		t.Errorf("Expected frontmatter-only file status to be %s, got %s", StatusFrontmatterOnly, status)
	}

	// Test frontmatter-and-content file
	status, err = scanner.checkFileStatus(frontmatterAndContentPath)
	if err != nil {
		t.Errorf("Failed to check frontmatter-and-content file status: %v", err)
	}
	if status != StatusNeedsReview {
		t.Errorf("Expected frontmatter-and-content file status to be %s, got %s", StatusNeedsReview, status)
	}

	// Test invalid frontmatter file
	status, err = scanner.checkFileStatus(invalidFrontmatterPath)
	if err != nil {
		t.Errorf("Failed to check invalid-frontmatter file status: %v", err)
	}
	if status != StatusNeedsReview {
		t.Errorf("Expected invalid-frontmatter file status to be %s, got %s", StatusNeedsReview, status)
	}
}

func TestExclusionFileHandling(t *testing.T) {
	// Create a temporary directory
	tempDir, err := os.MkdirTemp("", "scanner-test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Create an exclusion file with Obsidian links
	exclusionPath := filepath.Join(tempDir, "quality_exclude_links.md")
	exclusionContent := "# Excluded Files\n\n- [[excluded-file]]\n- [[another-excluded]]\n"
	if err := os.WriteFile(exclusionPath, []byte(exclusionContent), 0644); err != nil {
		t.Fatalf("Failed to create exclusion file: %v", err)
	}

	// Create test files
	excludedPath := filepath.Join(tempDir, "excluded-file.md")
	if err := os.WriteFile(excludedPath, []byte("# This should be excluded"), 0644); err != nil {
		t.Fatalf("Failed to create excluded file: %v", err)
	}

	anotherExcludedPath := filepath.Join(tempDir, "another-excluded.md")
	if err := os.WriteFile(anotherExcludedPath, []byte("# This should also be excluded"), 0644); err != nil {
		t.Fatalf("Failed to create another excluded file: %v", err)
	}

	includedPath := filepath.Join(tempDir, "included-file.md")
	if err := os.WriteFile(includedPath, []byte("# This should be included"), 0644); err != nil {
		t.Fatalf("Failed to create included file: %v", err)
	}

	// Create scanner with custom config pointing to our exclusion file
	cfg := config.GetDefaultConfig()
	cfg.ExclusionFile.Path = exclusionPath

	scanner, err := New(cfg)
	if err != nil {
		t.Fatalf("Failed to create scanner: %v", err)
	}

	// Check that the exclusion list contains the expected entries
	if !scanner.excludeList["excluded-file"] {
		t.Errorf("Expected 'excluded-file' to be in the exclusion list")
	}

	if !scanner.excludeList["another-excluded"] {
		t.Errorf("Expected 'another-excluded' to be in the exclusion list")
	}

	if scanner.excludeList["included-file"] {
		t.Errorf("Did not expect 'included-file' to be in the exclusion list")
	}

	// Test scanning the directory
	files, err := scanner.ScanDirectory(tempDir)
	if err != nil {
		t.Fatalf("Failed to scan directory: %v", err)
	}

	// Check that files are classified correctly
	fileStatuses := make(map[string]FileStatus)
	for _, file := range files {
		baseFile := filepath.Base(file.Path)
		fileStatuses[baseFile] = file.Status
	}

	// Verify exclusions
	if status, ok := fileStatuses["excluded-file.md"]; !ok || status != StatusExcluded {
		t.Errorf("Expected 'excluded-file.md' to have status %s, got %s", StatusExcluded, status)
	}

	if status, ok := fileStatuses["another-excluded.md"]; !ok || status != StatusExcluded {
		t.Errorf("Expected 'another-excluded.md' to have status %s, got %s", StatusExcluded, status)
	}

	// Verify inclusion
	if status, ok := fileStatuses["included-file.md"]; !ok || status != StatusNeedsReview {
		t.Errorf("Expected 'included-file.md' to have status %s, got %s", StatusNeedsReview, status)
	}

	// The exclusion file itself should be processed normally
	if status, ok := fileStatuses["quality_exclude_links.md"]; !ok || status != StatusNeedsReview {
		t.Errorf("Expected 'quality_exclude_links.md' to have status %s, got %s", StatusNeedsReview, status)
	}
}

func TestDirectoryExclusion(t *testing.T) {
	// Create a temporary directory structure
	tempDir, err := os.MkdirTemp("", "scanner-test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Create a directory to exclude
	excludeDir := filepath.Join(tempDir, ".obsidian")
	if err := os.Mkdir(excludeDir, 0755); err != nil {
		t.Fatalf("Failed to create exclude dir: %v", err)
	}

	// Create a file in the excluded directory
	excludedFilePath := filepath.Join(excludeDir, "excluded.md")
	if err := os.WriteFile(excludedFilePath, []byte("# This should be excluded"), 0644); err != nil {
		t.Fatalf("Failed to create file in excluded dir: %v", err)
	}

	// Create a file in the main directory
	includedFilePath := filepath.Join(tempDir, "included.md")
	if err := os.WriteFile(includedFilePath, []byte("# This should be included"), 0644); err != nil {
		t.Fatalf("Failed to create included file: %v", err)
	}

	// Create scanner with config that excludes .obsidian directory
	cfg := config.GetDefaultConfig()
	cfg.ScanSettings.ExcludeDirectories = []string{".obsidian"}

	scanner, err := New(cfg)
	if err != nil {
		t.Fatalf("Failed to create scanner: %v", err)
	}

	// Scan the directory
	files, err := scanner.ScanDirectory(tempDir)
	if err != nil {
		t.Fatalf("Failed to scan directory: %v", err)
	}

	// Check that only the included file is in the result
	if len(files) != 1 {
		t.Errorf("Expected 1 file, got %d", len(files))
	}

	if len(files) > 0 && filepath.Base(files[0].Path) != "included.md" {
		t.Errorf("Expected to find 'included.md', got '%s'", filepath.Base(files[0].Path))
	}
}

func TestReadFileContent(t *testing.T) {
	// Create a temporary directory
	tempDir, err := os.MkdirTemp("", "scanner-test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Create a test file
	testFilePath := filepath.Join(tempDir, "test.md")
	expectedContent := "# Test Content\n\nThis is a test file."
	if err := os.WriteFile(testFilePath, []byte(expectedContent), 0644); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	// Read the file content
	content, err := ReadFileContent(testFilePath)
	if err != nil {
		t.Errorf("Failed to read file content: %v", err)
	}

	// Verify the content
	if content != expectedContent {
		t.Errorf("Expected content '%s', got '%s'", expectedContent, content)
	}

	// Test reading a non-existent file
	_, err = ReadFileContent(filepath.Join(tempDir, "nonexistent.md"))
	if err == nil {
		t.Errorf("Expected error when reading non-existent file, got nil")
	}
}
