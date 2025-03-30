# Rate My KB

A Go CLI tool for scanning and evaluating the quality of Markdown files in your knowledge base. It categorizes files as empty, frontmatter-only, or low quality—and leverages GenAI (via Ollama) to help classify content quality.

## Table of Contents

- [Overview](#overview)
- [Features](#features)
- [Installation](#installation)
  - [Installation from GitHub Releases](#installation-from-github-releases)
  - [Building from Source](#building-from-source)
- [Usage](#usage)
- [Configuration](#configuration)
- [Exclusion File Format](#exclusion-file-format)
- [Generated Report](#generated-report)
- [Running Tests](#running-tests)
- [Dependencies](#dependencies)
  - [Installing and Setting Up Ollama](#installing-and-setting-up-ollama)
- [Contributing](#contributing)
- [License](#license)

## Overview

Rate My KB scans a directory of Markdown files to evaluate their quality. It helps you quickly identify:
- Files that are empty or contain only YAML frontmatter
- Files that may need content improvement based on AI-driven classification

This tool is fully configurable via a YAML file and integrates with an Ollama server for GenAI-powered quality classification.

## Features

- **File Scanning:** Identifies empty files and files with only frontmatter.
- **AI Classification:** Uses GenAI (via Ollama) to classify file quality.
- **Exclusions:** Supports excluding specific files or directories.
- **Reporting:** Generates a detailed Markdown report with categorized files.
- **Configurability:** Easily customizable with a YAML configuration file.

## Installation

You can get Rate My KB in one of two ways: by installing a prebuilt binary from the GitHub releases or by building it from source.

### Installation from GitHub Releases

If you’d like to install without building from source, download the binary from the [GitHub Releases page](https://github.com/philipf/rate-my-kb/releases/tag/v0.1.0).

For example:

#### On Linux/macOS:
```bash
# Download the tarball (adjust the file name as needed for your OS/architecture)
wget https://github.com/philipf/rate-my-kb/releases/download/v0.1.0/ratemykb-linux-amd64.tar.gz
# Extract the archive
tar -xzvf ratemykb-linux-amd64.tar.gz
# Move the binary to a directory in your PATH
sudo mv ratemykb /usr/local/bin/
```

#### On Windows:
- Download the appropriate `.exe` file from the release page.
- Add the folder containing `ratemykb.exe` to your system PATH or run it directly.

### Building from Source

To build Rate My KB from source, ensure you have [Go 1.18 or higher](https://golang.org/dl/) installed. Then clone the repository and build:

```bash
git clone https://github.com/philipf/rate-my-kb.git
cd rate-my-kb
go build
```

## Usage

Rate My KB accepts flags for configuration and targeting the directory with Markdown files.

```
Usage:
  ratemykb [flags]

Flags:
  -c, --config string   Path to configuration file
  -h, --help            help for ratemykb
  -t, --target string   Target folder containing Markdown files
```

**Examples:**

- **Default Usage:**
  ```bash
  ./ratemykb -t /path/to/knowledge-base
  ```
- **Specifying a Custom Configuration File:**
  ```bash
  ./ratemykb -c /path/to/config.yaml -t /path/to/knowledge-base
  ```

## Configuration

Create a `config.yaml` file to customize the behavior:

```yaml
ai_engine:
  url: "http://localhost:11434/"  # Ollama server URL
  model: "deepseek-r1:8b"          # GenAI model to use
scan_settings:
  file_extension: ".md"            # File extension to scan
  exclude_directories:
    - ".obsidian"
    - ".git"
    - "templates"
prompt_config:
  quality_classification_prompt: "Review the content and determine if it's: 'Empty', 'Low quality/low effort', or 'Good enough'."
exclusion_file:
  path: "quality_exclude_links.md"  # File containing links to exclude
```

## Exclusion File Format

The exclusion file should contain Obsidian-style links to files that should be skipped during quality checks:

```markdown
# Files to Exclude
- [[file-to-exclude]]
- [[another/file-to-exclude]]
```

## Generated Report

After running Rate My KB, a report named `vault-quality-report.md` is generated in the target folder. The report includes:

1. **Statistics** – An overview of scanned files.
2. **Empty Files** – Files with no content.
3. **Files with Frontmatter Only** – Files containing only YAML frontmatter.
4. **Low Quality/Low Effort Files** – Files flagged by the AI as low quality.

## Running Tests

To run tests for the project:

```bash
# Run all tests
go test ./...

# Run tests with verbose output
go test -v ./...

# Run integration tests
go test -run TestIntegration
```

## Dependencies

### Installing and Setting Up Ollama

Rate My KB relies on the Ollama server for AI classification. Follow these steps to install and set up Ollama:

1. **Install Ollama:**
   - Visit the [Ollama website](https://ollama.ai/download) to download and install Ollama for your platform.

   **Linux:**
   ```bash
   curl -fsSL https://ollama.com/install.sh | sh
   ```

   **macOS:**
   - Download the `.dmg` file from the website and follow the installation instructions.

   **Windows:**
   - Download the installer from the website and follow the installation steps.

2. **Start the Ollama Server:**
   - **Linux & macOS:**
     ```bash
     ollama serve
     ```
   - **Windows:**
     - The server should start automatically after installation. If not, launch it from the Start menu.

3. **Pull the DeepSeek Model:**
   Open a new terminal and run:
   ```bash
   ollama pull deepseek-r1:8b
   ```
   (Note: The model is approximately 4-5GB, so the download may take some time.)

4. **Verify the Installation:**
   Test that the model is working:
   ```bash
   ollama run deepseek-r1:8b "Hello, how are you today?"
   ```

Ensure your `config.yaml` points to your Ollama server. For example:
```yaml
ai_engine:
  url: "http://localhost:11434/"
  model: "deepseek-r1:8b"
```

If your Ollama server is on a different machine or port, adjust the URL accordingly.

## Contributing

We welcome contributions! If you'd like to help improve Rate My KB, please:

1. Fork the repository.
2. Create a feature branch (`git checkout -b feature/my-new-feature`).
3. Commit your changes (`git commit -am 'Add some feature'`).
4. Push to the branch (`git push origin feature/my-new-feature`).
5. Open a Pull Request.

For more detailed contribution guidelines, refer to the `CONTRIBUTING.md` file (coming soon).

## License

This project is licensed under the [MIT License](LICENSE).
