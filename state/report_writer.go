package state

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"ratemykb/output"
	"ratemykb/scanner"
)

// updateReport regenerates the report with all processed files
func (ps *ProcessingState) updateReport() error {
	// Create a temporary file for writing
	tempFile := ps.ReportPath + ".tmp"
	file, err := os.Create(tempFile)
	if err != nil {
		return fmt.Errorf("failed to create temp report file: %w", err)
	}

	// Generate report content
	var content strings.Builder

	// Add header
	content.WriteString("# Vault Quality Report\n\n")
	content.WriteString(fmt.Sprintf("Generated on: %s\n\n", time.Now().Format("2006-01-02 15:04:05")))
	content.WriteString(fmt.Sprintf("Target folder: `%s`\n\n", ps.TargetFolder))

	// Categorize files
	var emptyFiles, frontmatterOnlyFiles []output.ResultFile
	classificationMap := make(map[string][]output.ResultFile)

	for _, file := range ps.ProcessedFiles {
		if file.Status == scanner.StatusEmpty {
			emptyFiles = append(emptyFiles, file)
		} else if file.Status == scanner.StatusFrontmatterOnly {
			frontmatterOnlyFiles = append(frontmatterOnlyFiles, file)
		} else if file.Classification != "" {
			classStr := string(file.Classification)
			classificationMap[classStr] = append(classificationMap[classStr], file)
		}
	}

	// Add statistics
	content.WriteString("## Statistics\n\n")
	content.WriteString(fmt.Sprintf("- Total files processed: %d\n", len(ps.ProcessedFiles)))
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
		// Sort for consistent output
		sort.Slice(emptyFiles, func(i, j int) bool {
			return emptyFiles[i].Path < emptyFiles[j].Path
		})

		for _, file := range emptyFiles {
			link := formatObsidianLink(ps.TargetFolder, file.Path)
			content.WriteString(fmt.Sprintf("- %s\n", link))
		}
		content.WriteString("\n")
	}

	// Add frontmatter-only files section
	content.WriteString("## Files with Frontmatter Only\n\n")
	if len(frontmatterOnlyFiles) == 0 {
		content.WriteString("No files with frontmatter only found.\n\n")
	} else {
		// Sort for consistent output
		sort.Slice(frontmatterOnlyFiles, func(i, j int) bool {
			return frontmatterOnlyFiles[i].Path < frontmatterOnlyFiles[j].Path
		})

		for _, file := range frontmatterOnlyFiles {
			link := formatObsidianLink(ps.TargetFolder, file.Path)
			content.WriteString(fmt.Sprintf("- %s\n", link))
		}
		content.WriteString("\n")
	}

	// Add sections for each classification type
	var classTypes []string
	for classType := range classificationMap {
		classTypes = append(classTypes, classType)
	}
	sort.Strings(classTypes)

	for _, classType := range classTypes {
		classFiles := classificationMap[classType]
		content.WriteString(fmt.Sprintf("## %s Files\n\n", classType))
		if len(classFiles) == 0 {
			content.WriteString(fmt.Sprintf("No %s files found.\n\n", strings.ToLower(classType)))
		} else {
			// Sort for consistent output
			sort.Slice(classFiles, func(i, j int) bool {
				return classFiles[i].Path < classFiles[j].Path
			})

			for _, file := range classFiles {
				link := formatObsidianLink(ps.TargetFolder, file.Path)
				content.WriteString(fmt.Sprintf("- %s\n", link))
			}
			content.WriteString("\n")
		}
	}

	// Write content to temporary file
	_, err = file.WriteString(content.String())
	if err != nil {
		file.Close()
		os.Remove(tempFile)
		return fmt.Errorf("failed to write to temp report: %w", err)
	}

	// Close the file
	if err := file.Close(); err != nil {
		os.Remove(tempFile)
		return fmt.Errorf("failed to close temp report file: %w", err)
	}

	// Atomically replace the existing report
	if err := os.Rename(tempFile, ps.ReportPath); err != nil {
		os.Remove(tempFile)
		return fmt.Errorf("failed to replace report: %w", err)
	}

	return nil
}

// formatObsidianLink converts a file path to an Obsidian link format [[link-to-page]]
func formatObsidianLink(targetFolder, filePath string) string {
	// Make path relative to target folder
	relPath, err := filepath.Rel(targetFolder, filePath)
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