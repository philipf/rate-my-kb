package main

import (
	"ratemykb/cli"
)

// main is the entry point for the application
// It calls the cli.Execute function which handles the command-line interface
// and orchestrates the entire workflow:
// 1. Parse command-line arguments and validate the target folder
// 2. Load the configuration
// 3. Scan the target folder for Markdown files
// 4. Classify files requiring quality assessment
// 5. Generate the final markdown report
func main() {
	// Call the Execute function from the cli package
	cli.Execute()
}
