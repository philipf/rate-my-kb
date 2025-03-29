package config

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/viper"
)

// Config represents the application configuration structure
type Config struct {
	AIEngine      AIEngineConfig      `mapstructure:"ai_engine"`
	ScanSettings  ScanSettingsConfig  `mapstructure:"scan_settings"`
	PromptConfig  PromptConfig        `mapstructure:"prompt_config"`
	ExclusionFile ExclusionFileConfig `mapstructure:"exclusion_file"`
}

// AIEngineConfig represents the AI engine configuration
type AIEngineConfig struct {
	URL   string `mapstructure:"url"`
	Model string `mapstructure:"model"`
}

// ScanSettingsConfig represents the scanning settings
type ScanSettingsConfig struct {
	FileExtension      string   `mapstructure:"file_extension"`
	ExcludeDirectories []string `mapstructure:"exclude_directories"`
}

// PromptConfig represents the configuration for the GenAI prompt
type PromptConfig struct {
	QualityClassificationPrompt string `mapstructure:"quality_classification_prompt"`
}

// ExclusionFileConfig represents the configuration for the exclusion file
type ExclusionFileConfig struct {
	Path string `mapstructure:"path"`
}

// LoadConfig loads the configuration from the specified path or uses default values
func LoadConfig(configPath string) (*Config, error) {
	v := viper.New()

	// Set default values
	setDefaults(v)

	// If configuration path was provided, use it
	if configPath != "" {
		// If the path is a directory, append the default config filename
		fileInfo, err := os.Stat(configPath)
		if err == nil && fileInfo.IsDir() {
			configPath = filepath.Join(configPath, "config.yaml")
		}

		// Set the path to the configuration file
		v.SetConfigFile(configPath)

		// Try to read the configuration file
		if err := v.ReadInConfig(); err != nil {
			if _, ok := err.(viper.ConfigFileNotFoundError); ok {
				return nil, fmt.Errorf("config file not found at %s: %w", configPath, err)
			}
			return nil, fmt.Errorf("error reading config file: %w", err)
		}
	}

	// Unmarshal the configuration into a Config struct
	var config Config
	if err := v.Unmarshal(&config); err != nil {
		return nil, fmt.Errorf("unable to decode config into struct: %w", err)
	}

	return &config, nil
}

// setDefaults sets the default values for the configuration
func setDefaults(v *viper.Viper) {
	// AI Engine defaults
	v.SetDefault("ai_engine.url", "http://localhost:11434/")
	v.SetDefault("ai_engine.model", "gemma:12b")

	// Scan Settings defaults
	v.SetDefault("scan_settings.file_extension", ".md")
	v.SetDefault("scan_settings.exclude_directories", []string{})

	// Prompt Config defaults
	v.SetDefault("prompt_config.quality_classification_prompt",
		"Review the content and determine if it's: 'Empty', 'Low quality/low effort', or 'Good enough'.")

	// Exclusion File defaults
	v.SetDefault("exclusion_file.path", "quality_exclude_links.md")
}

// GetDefaultConfig returns a config object with default values
func GetDefaultConfig() *Config {
	v := viper.New()
	setDefaults(v)

	var config Config
	_ = v.Unmarshal(&config)

	return &config
}
