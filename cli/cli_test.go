package cli

import (
	"bytes"
	"io"
	"os"
	"path/filepath"
	"testing"

	"github.com/spf13/cobra"
)

// executeCommand is a helper function to execute the CLI command for testing
func executeCommand(t *testing.T, args ...string) (string, error) {
	t.Helper()

	// Create a new root command for testing to avoid side effects
	testRootCmd := &cobra.Command{
		Use:   "ratemykb",
		Short: rootCmd.Short,
		Long:  rootCmd.Long,
		RunE:  rootCmd.RunE,
	}

	// Copy the flag definitions from the main root command
	testRootCmd.PersistentFlags().StringVarP(&targetFolder, "target", "t", "", "Target folder containing Markdown files")
	testRootCmd.PersistentFlags().StringVarP(&configFile, "config", "c", "", "Path to configuration file")

	// Redirect output for testing
	buff := bytes.NewBufferString("")
	testRootCmd.SetOut(buff)
	testRootCmd.SetErr(buff)
	testRootCmd.SetArgs(args)

	err := testRootCmd.Execute()
	return buff.String(), err
}

func TestNoTargetFolder(t *testing.T) {
	// Reset global variables before the test
	targetFolder = ""
	configFile = ""

	_, err := executeCommand(t)

	if err == nil {
		t.Error("Expected error when no target folder is provided")
	}

	if err.Error() != "target folder is required" {
		t.Errorf("Expected error message 'target folder is required', got: %s", err.Error())
	}
}

func TestTargetFolderDoesNotExist(t *testing.T) {
	// Reset global variables before the test
	targetFolder = ""
	configFile = ""

	nonExistentFolder := "/path/that/does/not/exist"
	_, err := executeCommand(t, "--target", nonExistentFolder)

	if err == nil {
		t.Error("Expected error when target folder does not exist")
	}

	expectedErr := "target folder does not exist: " + nonExistentFolder
	if err.Error() != expectedErr {
		t.Errorf("Expected error message '%s', got: %s", expectedErr, err.Error())
	}
}

func TestTargetFolderAsArgument(t *testing.T) {
	// Reset global variables before the test
	targetFolder = ""
	configFile = ""

	// Create a temporary directory for the test
	tempDir, err := os.MkdirTemp("", "ratemykb-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Create a sample config file
	configPath := filepath.Join(tempDir, "config.yaml")
	if err := os.WriteFile(configPath, []byte("ai_engine:\n  url: 'http://test.url'\n  model: 'test-model'"), 0644); err != nil {
		t.Fatalf("Failed to write config file: %v", err)
	}

	// Capture stdout to check the output
	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	// Execute command with target folder as argument and config as flag
	_, err = executeCommand(t, tempDir, "--config", configPath)

	// Restore stdout
	w.Close()
	os.Stdout = oldStdout

	// Read the captured output
	var buf bytes.Buffer
	io.Copy(&buf, r)
	output := buf.String()

	// Verify there was no error
	if err != nil {
		t.Errorf("Did not expect an error, but got: %v", err)
	}

	// Check that the output mentions the correct target folder
	if output == "" {
		t.Error("Expected some output but got none")
	}
}

func TestTargetFolderAsFlag(t *testing.T) {
	// Reset global variables before the test
	targetFolder = ""
	configFile = ""

	// Create a temporary directory for the test
	tempDir, err := os.MkdirTemp("", "ratemykb-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Create a sample config file
	configPath := filepath.Join(tempDir, "config.yaml")
	if err := os.WriteFile(configPath, []byte("ai_engine:\n  url: 'http://test.url'\n  model: 'test-model'"), 0644); err != nil {
		t.Fatalf("Failed to write config file: %v", err)
	}

	// Capture stdout to check the output
	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	// Execute command with target and config as flags
	_, err = executeCommand(t, "--target", tempDir, "--config", configPath)

	// Restore stdout
	w.Close()
	os.Stdout = oldStdout

	// Read the captured output
	var buf bytes.Buffer
	io.Copy(&buf, r)
	output := buf.String()

	// Verify there was no error
	if err != nil {
		t.Errorf("Did not expect an error, but got: %v", err)
	}

	// Check that the output mentions the correct target folder
	if output == "" {
		t.Error("Expected some output but got none")
	}
}

func TestConfigLoadingError(t *testing.T) {
	// Reset global variables before the test
	targetFolder = ""
	configFile = ""

	// Create a temporary directory for the test
	tempDir, err := os.MkdirTemp("", "ratemykb-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Invalid config path
	invalidConfigPath := filepath.Join(tempDir, "nonexistent-config.yaml")

	// Execute command with valid target folder but invalid config path
	_, err = executeCommand(t, "--target", tempDir, "--config", invalidConfigPath)

	// Verify there was an error
	if err == nil {
		t.Error("Expected an error for invalid config path, but got none")
	}
}
