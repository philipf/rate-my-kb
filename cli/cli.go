package cli

import (
	"fmt"
	"os"
	"ratemykb/classification"
	"ratemykb/config"
	"ratemykb/output"
	"ratemykb/scanner"
	"ratemykb/state"

	"github.com/spf13/cobra"
)

var (
	// Used for flags
	configFile   string
	targetFolder string
	rootCmd      = &cobra.Command{
		Use:   "ratemykb",
		Short: "Rate My Knowledge Base - Evaluate Markdown files quality",
		Long: `Rate My Knowledge Base is a CLI tool that evaluates the quality of Markdown files
in an Obsidian vault or any directory containing Markdown files.
It classifies files as Empty, Low quality/low effort, or Good enough,
and generates a report in Markdown format.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			// If target folder not provided as a flag, check if it's provided as an argument
			if targetFolder == "" && len(args) > 0 {
				targetFolder = args[0]
			}

			// Validate that target folder is provided
			if targetFolder == "" {
				return fmt.Errorf("target folder is required")
			}

			// Check if target folder exists
			if _, err := os.Stat(targetFolder); os.IsNotExist(err) {
				return fmt.Errorf("target folder does not exist: %s", targetFolder)
			}

			// Load configuration
			cfg, err := config.LoadConfig(configFile)
			if err != nil {
				return fmt.Errorf("failed to load configuration: %w", err)
			}

			// Print the LLM model and endpoint
			fmt.Printf("LLM model: %s\n", cfg.AIEngine.Model)
			fmt.Printf("LLM endpoint: %s\n", cfg.AIEngine.URL)

			// Initialize state manager
			stateManager, err := state.New(targetFolder)
			if err != nil {
				return fmt.Errorf("failed to initialize state manager: %w", err)
			}

			// Initialize scanner
			fileScanner, err := scanner.New(cfg)
			if err != nil {
				return fmt.Errorf("failed to initialize scanner: %w", err)
			}

			// Scan the target folder
			fmt.Printf("Scanning %s for Markdown files...\n", targetFolder)
			files, err := fileScanner.ScanDirectory(targetFolder)
			if err != nil {
				return fmt.Errorf("failed to scan directory: %w", err)
			}
			fmt.Printf("Found %d Markdown files\n", len(files))

			// Initialize classifier
			classifier, err := classification.New(cfg)
			if err != nil {
				return fmt.Errorf("failed to initialize classifier: %w", err)
			}

			// Get total number of files to process
			totalFiles := len(files)
			totalAlreadyProcessed := 0
			fmt.Printf("Processing %d files...\n", totalFiles)

			// Helper function to show progress
			showProgress := func(i int, action, details string) {
				filesProcessed := i + 1
				percentComplete := float64(filesProcessed) / float64(totalFiles) * 100
				fmt.Printf("[%d/%d - %.1f%%] %s %s\n", filesProcessed, totalFiles, percentComplete, action, details)
			}

			// Process each file
			for i, file := range files {
				// Check if file has already been processed
				if stateManager.IsFileProcessed(file.Path) {
					totalAlreadyProcessed++
					showProgress(i, "Skipping (already processed)", file.Path)
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
						fmt.Printf("Warning: Could not read file %s: %v\n", file.Path, err)
						continue
					}

					// Classify the content
					showProgress(i, "Classifying", file.Path)
					result.Classification, err = classifier.ClassifyContent(content)

					if err != nil {
						fmt.Printf("Warning: Could not classify file %s: %v\n", file.Path, err)
						continue
					}

					// Print the classification result
					fmt.Printf("Classification result: %s\n", result.Classification)

				} else if file.Status == scanner.StatusEmpty {
					// Map scanner status to classification
					result.Classification = classification.Classification("Empty")
					showProgress(i, "Skipping classification for", file.Path+" (Empty)")
				} else if file.Status == scanner.StatusFrontmatterOnly {
					// Frontmatter-only files are considered low quality
					result.Classification = classification.Classification("Low quality")
					showProgress(i, "Skipping classification for", file.Path+" (Frontmatter-only)")
				} else if file.Status == scanner.StatusExcluded {
					// Show progress for excluded files
					showProgress(i, "Skipping", file.Path+" (Excluded)")
					continue // Don't add excluded files to the report
				}

				// Add processed file to state and update report
				if err := stateManager.AddProcessedFile(result); err != nil {
					fmt.Printf("Warning: Could not update report for %s: %v\n", file.Path, err)
				}
			}

			fmt.Printf("Processing complete: %d new files processed, %d already processed, %d total\n",
				len(stateManager.GetProcessedFiles())-totalAlreadyProcessed,
				totalAlreadyProcessed,
				len(stateManager.GetProcessedFiles()))

			// No need to generate a final report as it's been updated incrementally
			fmt.Printf("Report available at %s/vault-quality-report.md\n", targetFolder)
			return nil
		},
	}
)

// Execute is the entry point for the CLI application
// It handles command-line arguments and initiates the scanning process
func Execute() {
	// Add flags
	rootCmd.PersistentFlags().StringVarP(&targetFolder, "target", "t", "", "Target folder containing Markdown files")
	rootCmd.PersistentFlags().StringVarP(&configFile, "config", "c", "", "Path to configuration file")

	// Execute the command
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

// ResetForTesting resets the CLI state for testing purposes
// This is needed to avoid flag redefinition errors when running multiple tests
func ResetForTesting() {
	// Create a new root command
	rootCmd = &cobra.Command{
		Use:   "ratemykb",
		Short: "Rate My Knowledge Base - Evaluate Markdown files quality",
		Long: `Rate My Knowledge Base is a CLI tool that evaluates the quality of Markdown files
in an Obsidian vault or any directory containing Markdown files.
It classifies files as Empty, Low quality/low effort, or Good enough,
and generates a report in Markdown format.`,
		RunE: rootCmd.RunE,
	}

	// Add flags
	rootCmd.PersistentFlags().StringVarP(&targetFolder, "target", "t", "", "Target folder containing Markdown files")
	rootCmd.PersistentFlags().StringVarP(&configFile, "config", "c", "", "Path to configuration file")
}
