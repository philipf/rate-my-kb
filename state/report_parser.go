package state

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"ratemykb/classification"
	"ratemykb/output"
	"ratemykb/scanner"
)

// loadExistingReport reads the existing report and populates the processed files map
func (ps *ProcessingState) loadExistingReport() error {
	file, err := os.Open(ps.ReportPath)
	if err != nil {
		return fmt.Errorf("failed to open report: %w", err)
	}
	defer file.Close()

	// Parse the report to extract processed files
	fileScanner := bufio.NewScanner(file)
	currentSection := ""
	obsidianLinkPattern := regexp.MustCompile(`\[\[([^\]]+)\]\]`)

	for fileScanner.Scan() {
		line := fileScanner.Text()

		// Identify sections
		if strings.HasPrefix(line, "## ") {
			currentSection = strings.TrimPrefix(line, "## ")
			continue
		}

		// Process file entries in each section
		if strings.HasPrefix(line, "- [[") && currentSection != "" {
			matches := obsidianLinkPattern.FindStringSubmatch(line)
			if len(matches) >= 2 {
				obsidianLink := matches[1]

				// Convert Obsidian link back to file path
				filePath := ps.convertObsidianLinkToPath(obsidianLink)

				// Determine classification based on section
				var classificationStr string
				var status scanner.FileStatus

				// Handle special known cases
				switch currentSection {
				case "Empty Files":
					classificationStr = "Empty"
					status = scanner.StatusEmpty
				case "Files with Frontmatter Only":
					classificationStr = "Low quality"
					status = scanner.StatusFrontmatterOnly
				default:
					// For all other sections, use the section name as the classification
					// This handles any LLM-generated classification dynamically
					if strings.HasSuffix(currentSection, " Files") {
						// Strip "Files" suffix if present
						classificationStr = strings.TrimSuffix(currentSection, " Files")
					} else {
						classificationStr = currentSection
					}
					status = scanner.StatusNeedsReview
				}

				// Add to processed files
				ps.ProcessedFiles[filePath] = output.ResultFile{
					Path:           filePath,
					Status:         status,
					Classification: classification.Classification(classificationStr),
				}
			}
		}
	}

	return fileScanner.Err()
}

// convertObsidianLinkToPath converts an Obsidian link back to a file path
func (ps *ProcessingState) convertObsidianLinkToPath(obsidianLink string) string {
	// Convert forward slashes to path separators
	pathWithoutExt := strings.ReplaceAll(obsidianLink, "/", string(filepath.Separator))

	// Add file extension and target folder path
	return filepath.Join(ps.TargetFolder, pathWithoutExt+".md")
}