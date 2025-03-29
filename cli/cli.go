package cli

import (
	"fmt"
	"os"
	"ratemykb/config"

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

			// Here we would continue with the scanning and classification process
			// For now, just print a message with the loaded configuration
			fmt.Printf("Successfully loaded configuration for target folder: %s\n", targetFolder)
			fmt.Printf("AI Engine: %s, Model: %s\n", cfg.AIEngine.URL, cfg.AIEngine.Model)

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
