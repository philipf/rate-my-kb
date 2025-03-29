package scanner

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"ratemykb/config"
)

// FileStatus represents the pre-check status of a markdown file
type FileStatus string

const (
	// StatusEmpty indicates the file is empty (no content after trimming whitespace)
	StatusEmpty FileStatus = "Empty"

	// StatusFrontmatterOnly indicates the file contains only frontmatter
	StatusFrontmatterOnly FileStatus = "Frontmatter-only"

	// StatusNeedsReview indicates the file has content and should be checked by the AI
	StatusNeedsReview FileStatus = "Needs-review"

	// StatusExcluded indicates the file is in the exclusion list
	StatusExcluded FileStatus = "Excluded"
)

// File represents a markdown file with its path and status
type File struct {
	Path   string     // Path to the file
	Status FileStatus // Status of the file based on pre-checks
}

// Scanner handles the scanning of markdown files in a directory
type Scanner struct {
	config      *config.Config
	excludeList map[string]bool // Map of files to exclude
}

// New creates a new Scanner with the provided configuration
func New(cfg *config.Config) (*Scanner, error) {
	scanner := &Scanner{
		config:      cfg,
		excludeList: make(map[string]bool),
	}

	// Parse exclusion file if it exists
	if cfg.ExclusionFile.Path != "" {
		if err := scanner.parseExclusionFile(cfg.ExclusionFile.Path); err != nil {
			return nil, fmt.Errorf("failed to parse exclusion file: %w", err)
		}
	}

	return scanner, nil
}

// ScanDirectory recursively scans the target directory for markdown files
// and returns a list of files with their pre-check status
func (s *Scanner) ScanDirectory(targetDir string) ([]File, error) {
	var files []File

	// Walk through the directory tree
	err := filepath.Walk(targetDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Skip directories
		if info.IsDir() {
			// Check if this directory should be excluded
			for _, excludeDir := range s.config.ScanSettings.ExcludeDirectories {
				if info.Name() == excludeDir || (strings.HasPrefix(excludeDir, "/") &&
					strings.HasPrefix(filepath.ToSlash(path), filepath.ToSlash(filepath.Join(targetDir, strings.TrimPrefix(excludeDir, "/"))))) {
					return filepath.SkipDir
				}
			}
			return nil
		}

		// Process only files with the configured extension
		if filepath.Ext(path) == s.config.ScanSettings.FileExtension {
			// Normalize path for exclusion check
			normalizedPath := s.normalizePathForExclusionCheck(path)

			// Skip if file is in exclusion list
			if s.excludeList[normalizedPath] {
				files = append(files, File{
					Path:   path,
					Status: StatusExcluded,
				})
				return nil
			}

			// Perform pre-checks on the file
			status, err := s.checkFileStatus(path)
			if err != nil {
				// Log error but continue processing other files
				fmt.Printf("Warning: Error checking file %s: %v\n", path, err)
				return nil
			}

			// Add file with its status to the result
			files = append(files, File{
				Path:   path,
				Status: status,
			})
		}

		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("error scanning directory: %w", err)
	}

	return files, nil
}

// checkFileStatus performs pre-checks on a file and returns its status
func (s *Scanner) checkFileStatus(filePath string) (FileStatus, error) {
	content, err := os.ReadFile(filePath)
	if err != nil {
		return "", fmt.Errorf("failed to read file: %w", err)
	}

	// Check if file is empty
	trimmedContent := strings.TrimSpace(string(content))
	if trimmedContent == "" {
		return StatusEmpty, nil
	}

	// Check if file contains only frontmatter
	if s.isFrontmatterOnly(trimmedContent) {
		return StatusFrontmatterOnly, nil
	}

	return StatusNeedsReview, nil
}

// isFrontmatterOnly checks if the content contains only YAML frontmatter
func (s *Scanner) isFrontmatterOnly(content string) bool {
	lines := strings.Split(content, "\n")

	// Check for YAML frontmatter
	if len(lines) < 2 || lines[0] != "---" {
		return false
	}

	// Find the end of frontmatter
	endIndex := -1
	for i := 1; i < len(lines); i++ {
		if lines[i] == "---" {
			endIndex = i
			break
		}
	}

	// No end marker found
	if endIndex == -1 {
		return false
	}

	// Check if there's any content after the frontmatter
	for i := endIndex + 1; i < len(lines); i++ {
		if strings.TrimSpace(lines[i]) != "" {
			return false
		}
	}

	return true
}

// parseExclusionFile reads the exclusion file and extracts Obsidian links
func (s *Scanner) parseExclusionFile(filePath string) error {
	file, err := os.Open(filePath)
	if err != nil {
		// If file doesn't exist, just return without error
		if os.IsNotExist(err) {
			return nil
		}
		return fmt.Errorf("failed to open exclusion file: %w", err)
	}
	defer file.Close()

	// Regular expression to match Obsidian links [[link-to-page]]
	linkPattern := regexp.MustCompile(`\[\[([^\]]+)\]\]`)

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		matches := linkPattern.FindAllStringSubmatch(line, -1)

		for _, match := range matches {
			if len(match) >= 2 {
				// Add the link to the exclusion list
				linkText := match[1]
				s.excludeList[linkText] = true

				// Also add with .md extension if it doesn't have one
				if !strings.HasSuffix(linkText, ".md") {
					s.excludeList[linkText+".md"] = true
				}
			}
		}
	}

	if err := scanner.Err(); err != nil {
		return fmt.Errorf("error reading exclusion file: %w", err)
	}

	return nil
}

// normalizePathForExclusionCheck converts a file path to the format used in Obsidian links
func (s *Scanner) normalizePathForExclusionCheck(path string) string {
	// Extract just the filename without extension
	filename := filepath.Base(path)
	fileExt := filepath.Ext(filename)
	filenameWithoutExt := strings.TrimSuffix(filename, fileExt)

	return filenameWithoutExt
}

// ReadFileContent reads and returns the content of a file
func ReadFileContent(filePath string) (string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return "", fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	content, err := io.ReadAll(file)
	if err != nil {
		return "", fmt.Errorf("failed to read file content: %w", err)
	}

	return string(content), nil
}
