package cli

import (
	"fmt"
	"os"
	"ratemykb/classification"
	"ratemykb/config"
	"ratemykb/output"
	"ratemykb/scanner"

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

			// Print the configuration
			fmt.Printf("Configuration: %+v\n", cfg)

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

			// Prepare results for output
			var results []output.ResultFile

			// Process each file
			for _, file := range files {
				// Create a result file with default classification
				result := output.ResultFile{
					Path:           file.Path,
					Status:         file.Status,
					Classification: classification.ClassificationUnknown,
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
					fmt.Printf("Classifying %s...\n", file.Path)
					result.Classification, err = classifier.ClassifyContent(content)

					if err != nil {
						fmt.Printf("Warning: Could not classify file %s: %v\n", file.Path, err)
						continue
					}

					// Print the classification result
					fmt.Printf("Classification result: %s\n", result.Classification)

				} else if file.Status == scanner.StatusEmpty {
					// Map scanner status to classification
					result.Classification = classification.ClassificationEmpty
				} else if file.Status == scanner.StatusFrontmatterOnly {
					// Frontmatter-only files are considered low quality
					result.Classification = classification.ClassificationLowQuality
				}

				// Add to results
				results = append(results, result)
			}

			// Generate the report
			outputGenerator := output.New(targetFolder)
			if err := outputGenerator.CreateReport(results); err != nil {
				return fmt.Errorf("failed to generate report: %w", err)
			}

			fmt.Printf("Report generated successfully in %s/vault-quality-report.md\n", targetFolder)
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
