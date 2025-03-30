package config

import (
	"os"
	"path/filepath"
	"reflect"
	"testing"
)

func TestLoadConfig(t *testing.T) {
	// Test loading default configuration
	t.Run("Default Configuration", func(t *testing.T) {
		config, err := LoadConfig("")
		if err != nil {
			t.Fatalf("LoadConfig() error = %v", err)
		}

		// Verify default values
		if config.AIEngine.URL != "http://localhost:11434/" {
			t.Errorf("Expected default AIEngine.URL to be 'http://localhost:11434/', got %s", config.AIEngine.URL)
		}

		if config.AIEngine.Model != "gemma3:1b" {
			t.Errorf("Expected default AIEngine.Model to be 'gemma3:1b', got %s", config.AIEngine.Model)
		}

		if config.ScanSettings.FileExtension != ".md" {
			t.Errorf("Expected default ScanSettings.FileExtension to be '.md', got %s", config.ScanSettings.FileExtension)
		}

		if !reflect.DeepEqual(config.ScanSettings.ExcludeDirectories, []string{}) {
			t.Errorf("Expected default ScanSettings.ExcludeDirectories to be empty, got %v", config.ScanSettings.ExcludeDirectories)
		}

		if config.ExclusionFile.Path != "quality_exclude_links.md" {
			t.Errorf("Expected default ExclusionFile.Path to be 'quality_exclude_links.md', got %s", config.ExclusionFile.Path)
		}
	})

	// Test loading custom configuration
	t.Run("Custom Configuration", func(t *testing.T) {
		// Create a temporary directory
		tempDir, err := os.MkdirTemp("", "config_test")
		if err != nil {
			t.Fatalf("Failed to create temp directory: %v", err)
		}
		defer os.RemoveAll(tempDir)

		// Create a test configuration file
		configPath := filepath.Join(tempDir, "config.yaml")
		configContent := `
ai_engine:
  url: "https://api.openai.com/v1/"
  model: "gpt-4"
scan_settings:
  file_extension: ".md"
  exclude_directories:
    - ".obsidian"
    - ".git"
prompt_config:
  quality_classification_prompt: "Custom prompt for classification"
exclusion_file:
  path: "custom_exclusion_file.md"
`
		err = os.WriteFile(configPath, []byte(configContent), 0644)
		if err != nil {
			t.Fatalf("Failed to write test config file: %v", err)
		}

		// Load the custom configuration
		config, err := LoadConfig(configPath)
		if err != nil {
			t.Fatalf("LoadConfig() error = %v", err)
		}

		// Verify custom values
		if config.AIEngine.URL != "https://api.openai.com/v1/" {
			t.Errorf("Expected AIEngine.URL to be 'https://api.openai.com/v1/', got %s", config.AIEngine.URL)
		}

		if config.AIEngine.Model != "gpt-4" {
			t.Errorf("Expected AIEngine.Model to be 'gpt-4', got %s", config.AIEngine.Model)
		}

		expectedExcludeDirs := []string{".obsidian", ".git"}
		if !reflect.DeepEqual(config.ScanSettings.ExcludeDirectories, expectedExcludeDirs) {
			t.Errorf("Expected ScanSettings.ExcludeDirectories to be %v, got %v", expectedExcludeDirs, config.ScanSettings.ExcludeDirectories)
		}

		if config.PromptConfig.QualityClassificationPrompt != "Custom prompt for classification" {
			t.Errorf("Expected PromptConfig.QualityClassificationPrompt to be 'Custom prompt for classification', got %s", config.PromptConfig.QualityClassificationPrompt)
		}

		if config.ExclusionFile.Path != "custom_exclusion_file.md" {
			t.Errorf("Expected ExclusionFile.Path to be 'custom_exclusion_file.md', got %s", config.ExclusionFile.Path)
		}
	})

	// Test loading from a directory path
	t.Run("Directory Path", func(t *testing.T) {
		// Create a temporary directory
		tempDir, err := os.MkdirTemp("", "config_test_dir")
		if err != nil {
			t.Fatalf("Failed to create temp directory: %v", err)
		}
		defer os.RemoveAll(tempDir)

		// Create a test configuration file in that directory
		configPath := filepath.Join(tempDir, "config.yaml")
		configContent := `
ai_engine:
  url: "https://api.custom.com/v1/"
  model: "custom-model"
`
		err = os.WriteFile(configPath, []byte(configContent), 0644)
		if err != nil {
			t.Fatalf("Failed to write test config file: %v", err)
		}

		// Load the configuration by providing the directory path
		config, err := LoadConfig(tempDir)
		if err != nil {
			t.Fatalf("LoadConfig() error = %v", err)
		}

		// Verify custom values
		if config.AIEngine.URL != "https://api.custom.com/v1/" {
			t.Errorf("Expected AIEngine.URL to be 'https://api.custom.com/v1/', got %s", config.AIEngine.URL)
		}

		if config.AIEngine.Model != "custom-model" {
			t.Errorf("Expected AIEngine.Model to be 'custom-model', got %s", config.AIEngine.Model)
		}
	})

	// Test error handling for non-existent configuration file
	t.Run("Non-existent File", func(t *testing.T) {
		_, err := LoadConfig("/non/existent/path.yaml")
		if err == nil {
			t.Errorf("Expected an error when loading non-existent config file, got nil")
		}
	})

	// Test error handling for invalid YAML
	t.Run("Invalid YAML", func(t *testing.T) {
		// Create a temporary file with invalid YAML
		tempFile, err := os.CreateTemp("", "invalid_config*.yaml")
		if err != nil {
			t.Fatalf("Failed to create temp file: %v", err)
		}
		defer os.Remove(tempFile.Name())

		// Write invalid YAML content
		invalidContent := `
ai_engine:
  url: "http://localhost:11434/
  model: gemma3:1b
this is not valid yaml
`
		err = os.WriteFile(tempFile.Name(), []byte(invalidContent), 0644)
		if err != nil {
			t.Fatalf("Failed to write invalid config file: %v", err)
		}

		// Try to load the invalid configuration
		_, err = LoadConfig(tempFile.Name())
		if err == nil {
			t.Errorf("Expected an error when loading invalid YAML, got nil")
		}
	})
}

func TestGetDefaultConfig(t *testing.T) {
	config := GetDefaultConfig()

	if config.AIEngine.URL != "http://localhost:11434/" {
		t.Errorf("Expected default AIEngine.URL to be 'http://localhost:11434/', got %s", config.AIEngine.URL)
	}

	if config.AIEngine.Model != "gemma3:1b" {
		t.Errorf("Expected default AIEngine.Model to be 'gemma3:1b', got %s", config.AIEngine.Model)
	}

	if config.ScanSettings.FileExtension != ".md" {
		t.Errorf("Expected default ScanSettings.FileExtension to be '.md', got %s", config.ScanSettings.FileExtension)
	}
}
