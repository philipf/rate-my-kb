package state

import (
	"fmt"
	"os"
	"path/filepath"

	"ratemykb/output"
)

// ProcessingState manages the state of file processing
type ProcessingState struct {
	TargetFolder   string
	ReportPath     string
	ProcessedFiles map[string]output.ResultFile
}

// New creates a new ProcessingState and loads existing state if a report exists
func New(targetFolder string) (*ProcessingState, error) {
	ps := &ProcessingState{
		TargetFolder:   targetFolder,
		ReportPath:     filepath.Join(targetFolder, "vault-quality-report.md"),
		ProcessedFiles: make(map[string]output.ResultFile),
	}

	// Load existing state from report if it exists
	if _, err := os.Stat(ps.ReportPath); err == nil {
		if err := ps.loadExistingReport(); err != nil {
			return nil, fmt.Errorf("failed to load existing report: %w", err)
		}
		fmt.Printf("Found existing report with %d processed files\n", len(ps.ProcessedFiles))
	}

	return ps, nil
}

// IsFileProcessed checks if a file has already been processed
func (ps *ProcessingState) IsFileProcessed(filePath string) bool {
	_, exists := ps.ProcessedFiles[filePath]
	return exists
}

// AddProcessedFile adds a processed file to the state and updates the report
func (ps *ProcessingState) AddProcessedFile(file output.ResultFile) error {
	// Add to processed files map
	ps.ProcessedFiles[file.Path] = file

	// Update the report
	return ps.updateReport()
}

// GetProcessedFiles returns the map of processed files
func (ps *ProcessingState) GetProcessedFiles() map[string]output.ResultFile {
	return ps.ProcessedFiles
}