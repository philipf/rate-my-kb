# TODO Checklist: Golang CLI Obsidian Quality Scanner

This checklist breaks down the project into thorough, incremental tasks. Check off items as you complete each step.

---

## 1. Project Setup & Repository Structure
- [x] **Initialize Project**
  - [x] Run `go mod init ratemykb` to initialize the Go module.
  - [x] Set up the basic repository structure.
- [x] **Create Directories**
  - [x] Create a `cli` directory.
  - [x] Create a `config` directory.
  - [x] Create a `scanner` directory.
  - [x] Create a `classification` directory.
  - [x] Create an `output` directory.
- [x] **Minimal Main File**
  - [x] Write `main.go` that calls a dummy `Execute()` function from the `cli` package.
  - [x] Ensure the project compiles without errors.

---

## 2. Configuration Module
- [x] **Define Configuration Struct**
  - [x] Create a struct to hold:
    - `ai_engine` (with `url` defaulting to `"http://localhost:11434/"` and `model` defaulting to `"gemma3:1b"`)
    - `scan_settings` (with `file_extension` set to `.md` and a list of `exclude_directories`)
    - `prompt_config` (containing `quality_classification_prompt`)
    - `exclusion_file` (with `path` for `quality_exclude_links.md`)
- [x] **Implement Viper Loading**
  - [x] Use Viper to read the YAML configuration file.
  - [x] Apply default values where necessary.
  - [x] Add error handling for missing/invalid configuration.
- [x] **Unit Tests for Configuration**
  - [x] Write tests to verify proper loading of configuration.
  - [x] Test that defaults are correctly applied.
  - [x] Verify error messages when configuration is missing or invalid.

---

## 3. CLI Module with Cobra
- [x] **Set Up Cobra Command**
  - [x] Create a root command that accepts a target folder as an argument or flag.
  - [x] Integrate configuration loading within the command execution.
- [x] **Error Handling**
  - [x] Ensure that errors in configuration loading are caught and logged.
- [x] **Unit Tests for CLI**
  - [x] Write tests to verify command-line argument parsing.
  - [x] Test that the configuration is loaded when the CLI command is executed.

---

## 4. File Scanner Module
- [x] **Recursive File Scanning**
  - [x] Implement functionality to scan the target folder recursively.
  - [x] Filter to only process files with the `.md` extension.
- [x] **Directory Exclusion**
  - [x] Exclude directories as specified in the configuration.
- [x] **Pre-checks Implementation**
  - [x] **Empty File Check:** 
    - [x] Read file content, trim whitespace, and flag as empty if no content remains.
  - [x] **Frontmatter-only Check:**
    - [x] Detect if a file starts with `---` and ends with `---` with no additional content.
- [x] **Exclusion File Parsing**
  - [x] Parse `quality_exclude_links.md` to extract Obsidian links (e.g., `[[link-to-page]]`).
  - [x] Skip files that appear in the exclusion list.
- [x] **Unit Tests for Scanner**
  - [x] Write tests for file scanning.
  - [x] Test pre-checks for empty and frontmatter-only files.
  - [x] Validate exclusion file parsing and matching.

---

## 5. Classification Module
- [x] **Integrate with GenAI Engine**
  - [x] Create a function to send file content and the classification prompt to the GenAI engine.
  - [x] Integrate with LangChainGo to connect to Ollama.
- [x] **Response Handling**
  - [x] Parse responses to classify files as "Empty", "Low quality/low effort", or "Good enough".
- [x] **Unit Tests for Classification**
  - [x] Use mocks to simulate various GenAI engine responses.
  - [x] Verify that the function correctly interprets and returns classifications.

---

## 6. Output Module
- [x] **Report Generation**
  - [x] Develop a function to generate a markdown report.
  - [x] Create three sections in the report:
    - **Empty Files**
    - **Files with Frontmatter Only**
    - **Low Quality/Low Effort Files**
- [x] **Formatting**
  - [x] Format each entry using Obsidian-style links (e.g., `[[link-to-page]]`).
  - [x] Write the report to a markdown file in the target folder.
- [x] **Unit Tests for Output**
  - [x] Test that the report is generated with correct sections and formatting.
  - [x] Validate that file entries match expected Obsidian link format.

---

## 7. Integration & End-to-End Wiring
- [x] **Wire Modules in `main.go`**
  - [x] Use the CLI module to parse command-line arguments.
  - [x] Load the configuration using the configuration module.
  - [x] Scan the target folder using the file scanner module.
  - [x] For files requiring classification, call the classification module.
  - [x] Generate the final report using the output module.
- [x] **Error Logging and Handling**
  - [x] Ensure all modules log errors appropriately.
  - [x] Verify the workflow handles errors gracefully.
- [x] **Integration Testing**
  - [x] Write integration tests that simulate a full run on a temporary folder.
  - [x] Include tests for all edge cases (empty files, frontmatter-only files, files to be classified).
  - [x] Confirm that the final markdown report correctly categorizes and formats all entries.

---

## 8. Documentation & Code Quality
- [x] **In-Code Documentation**
  - [x] Add comments to explain the purpose of functions and modules.
- [x] **External Documentation**
  - [x] Create a README that explains:
    - Configuration setup
    - Command-line usage
    - Running tests
    - Interpreting the generated report
- [x] **Code Review**
  - [x] Review all code for adherence to Go best practices.
  - [x] Ensure there are no orphaned or unused code sections.

---

## 9. Final Steps
- [x] **Final Code Review and Testing**
  - [x] Run all unit and integration tests to ensure everything works.
  - [x] Check that the application builds and runs without errors.
- [x] **Prepare for Deployment**
  - [x] Update the README with build and deployment instructions.
  - [x] Ensure logging and error messages are clear for users.

---

Happy coding! Use this checklist to track progress and ensure every step is thoroughly tested and integrated.
