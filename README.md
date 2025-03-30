# Rate My KB

A Go CLI tool for scanning and evaluating the quality of Markdown files in a knowledge base. It scans a directory of Markdown files, determines their quality, and generates a report categorizing them as empty, frontmatter-only, or low quality.

## Project Structure

- `cli`: Command-line interface handling using Cobra
- `config`: Configuration management with Viper
- `scanner`: File scanning and parsing with pre-checks
- `classification`: Quality classification logic with GenAI integration
- `output`: Report generation in Markdown format

## Features

- Identifies empty files and files containing only frontmatter
- Uses AI to classify the quality of files with content
- Supports exclusion of specific files and directories
- Generates a detailed Markdown report with categorized files
- Fully configurable via YAML configuration

## Building the Project

```bash
go build
```

## Usage

```bash
# Basic usage with default settings
./ratemykb /path/to/knowledge-base

# Specify a configuration file
./ratemykb --config /path/to/config.yaml /path/to/knowledge-base

# Alternative target folder specification
./ratemykb --target /path/to/knowledge-base
```

## Configuration

Create a `config.yaml` file to customize the behavior:

```yaml
ai_engine:
  url: "http://localhost:11434/"  # Ollama server URL
  model: "gemma3:1b"              # GenAI model to use

scan_settings:
  file_extension: ".md"           # File extension to scan
  exclude_directories:            # Directories to exclude
    - ".obsidian"
    - ".git"
    - "templates"

prompt_config:
  quality_classification_prompt: "Review the content and determine if it's: 'Empty', 'Low quality/low effort', or 'Good enough'."

exclusion_file:
  path: "quality_exclude_links.md"  # File containing links to exclude
```

## Exclusion File Format

The exclusion file should contain Obsidian-style links to files that should be excluded from quality checks:

```markdown
# Files to Exclude

- [[file-to-exclude]]
- [[another/file-to-exclude]]
```

## Generated Report

The report is generated in the target folder as `vault-quality-report.md` and contains:

1. **Statistics** - Overview of all scanned files
2. **Empty Files** - Files with no content
3. **Files with Frontmatter Only** - Files containing only YAML frontmatter
4. **Low Quality/Low Effort Files** - Files classified as low quality by the AI

## Running Tests

```bash
# Run all tests
go test ./...

# Run tests with verbose output
go test -v ./...

# Run integration test
go test -run TestIntegration
```

## Requirements

- Go 1.18 or higher
- Ollama server running locally (or accessible via network) for AI classification

## License

MIT 