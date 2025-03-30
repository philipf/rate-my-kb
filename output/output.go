package output

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"ratemykb/classification"
	"ratemykb/scanner"
)

// ResultFile represents a file entry for the final report
type ResultFile struct {
	Path           string                        // Full path to the file
	Status         scanner.FileStatus            // Status from scanner pre-checks
	Classification classification.Classification // Classification from the AI
}

// Generator handles the generation of the final report
type Generator struct {
	targetFolder string // The root folder being scanned
}

// New creates a new Generator instance
func New(targetFolder string) *Generator {
	return &Generator{
		targetFolder: targetFolder,
	}
}

// CreateReport generates a markdown report from the scan results
// and writes it to a file in the target folder
func (g *Generator) CreateReport(files []ResultFile) error {
	// Categorize files
	var emptyFiles, frontmatterOnlyFiles []ResultFile

	// Map to store files by classification
	classificationMap := make(map[string][]ResultFile)

	for _, file := range files {
		if file.Status == scanner.StatusEmpty {
			emptyFiles = append(emptyFiles, file)
		} else if file.Status == scanner.StatusFrontmatterOnly {
			frontmatterOnlyFiles = append(frontmatterOnlyFiles, file)
		} else if file.Classification != "" {
			// Group files by their classification
			classStr := string(file.Classification)
			classificationMap[classStr] = append(classificationMap[classStr], file)
		}
	}

	// Generate report content
	var content strings.Builder

	// Add header
	content.WriteString("# Vault Quality Report\n\n")
	content.WriteString(fmt.Sprintf("Generated on: %s\n\n", time.Now().Format("2006-01-02 15:04:05")))
	content.WriteString(fmt.Sprintf("Target folder: `%s`\n\n", g.targetFolder))

	// Add statistics
	content.WriteString("## Statistics\n\n")
	content.WriteString(fmt.Sprintf("- Total files scanned: %d\n", len(files)))
	content.WriteString(fmt.Sprintf("- Empty files: %d\n", len(emptyFiles)))
	content.WriteString(fmt.Sprintf("- Files with frontmatter only: %d\n", len(frontmatterOnlyFiles)))

	// Add statistics for each classification type
	for classType, classFiles := range classificationMap {
		content.WriteString(fmt.Sprintf("- %s files: %d\n", classType, len(classFiles)))
	}
	content.WriteString("\n")

	// Add empty files section
	content.WriteString("## Empty Files\n\n")
	if len(emptyFiles) == 0 {
		content.WriteString("No empty files found.\n\n")
	} else {
		for _, file := range emptyFiles {
			link := g.formatObsidianLink(file.Path)
			content.WriteString(fmt.Sprintf("- %s\n", link))
		}
		content.WriteString("\n")
	}

	// Add frontmatter-only files section
	content.WriteString("## Files with Frontmatter Only\n\n")
	if len(frontmatterOnlyFiles) == 0 {
		content.WriteString("No files with frontmatter only found.\n\n")
	} else {
		for _, file := range frontmatterOnlyFiles {
			link := g.formatObsidianLink(file.Path)
			content.WriteString(fmt.Sprintf("- %s\n", link))
		}
		content.WriteString("\n")
	}

	// Add sections for each classification type
	for classType, classFiles := range classificationMap {
		content.WriteString(fmt.Sprintf("## %s Files\n\n", classType))
		if len(classFiles) == 0 {
			content.WriteString(fmt.Sprintf("No %s files found.\n\n", strings.ToLower(classType)))
		} else {
			for _, file := range classFiles {
				link := g.formatObsidianLink(file.Path)
				content.WriteString(fmt.Sprintf("- %s\n", link))
			}
			content.WriteString("\n")
		}
	}

	// Write report to file
	reportPath := filepath.Join(g.targetFolder, "vault-quality-report.md")
	err := os.WriteFile(reportPath, []byte(content.String()), 0644)
	if err != nil {
		return fmt.Errorf("failed to write report: %w", err)
	}

	return nil
}

// formatObsidianLink converts a file path to an Obsidian link format [[link-to-page]]
func (g *Generator) formatObsidianLink(filePath string) string {
	// Make path relative to target folder
	relPath, err := filepath.Rel(g.targetFolder, filePath)
	if err != nil {
		// Fallback to base name if relative path fails
		relPath = filepath.Base(filePath)
	}

	// Remove file extension
	baseName := strings.TrimSuffix(relPath, filepath.Ext(relPath))

	// Convert path separators to forward slashes for Obsidian format
	baseName = strings.ReplaceAll(baseName, string(filepath.Separator), "/")

	// Format as Obsidian link [[link-to-page]]
	return fmt.Sprintf("[[%s]]", baseName)
}
