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
- [ ] **Define Configuration Struct**
  - [ ] Create a struct to hold:
    - `ai_engine` (with `url` defaulting to `"http://localhost:11434/"` and `model` defaulting to `"gemma:12b"`)
    - `scan_settings` (with `file_extension` set to `.md` and a list of `exclude_directories`)
    - `prompt_config` (containing `quality_classification_prompt`)
    - `exclusion_file` (with `path` for `quality_exclude_links.md`)
- [ ] **Implement Viper Loading**
  - [ ] Use Viper to read the YAML configuration file.
  - [ ] Apply default values where necessary.
  - [ ] Add error handling for missing/invalid configuration.
- [ ] **Unit Tests for Configuration**
  - [ ] Write tests to verify proper loading of configuration.
  - [ ] Test that defaults are correctly applied.
  - [ ] Verify error messages when configuration is missing or invalid.

---

## 3. CLI Module with Cobra
- [ ] **Set Up Cobra Command**
  - [ ] Create a root command that accepts a target folder as an argument or flag.
  - [ ] Integrate configuration loading within the command execution.
- [ ] **Error Handling**
  - [ ] Ensure that errors in configuration loading are caught and logged.
- [ ] **Unit Tests for CLI**
  - [ ] Write tests to verify command-line argument parsing.
  - [ ] Test that the configuration is loaded when the CLI command is executed.

---

## 4. File Scanner Module
- [ ] **Recursive File Scanning**
  - [ ] Implement functionality to scan the target folder recursively.
  - [ ] Filter to only process files with the `.md` extension.
- [ ] **Directory Exclusion**
  - [ ] Exclude directories as specified in the configuration.
- [ ] **Pre-checks Implementation**
  - [ ] **Empty File Check:** 
    - [ ] Read file content, trim whitespace, and flag as empty if no content remains.
  - [ ] **Frontmatter-only Check:**
    - [ ] Detect if a file starts with `---` and ends with `---` with no additional content.
- [ ] **Exclusion File Parsing**
  - [ ] Parse `quality_exclude_links.md` to extract Obsidian links (e.g., `[[link-to-page]]`).
  - [ ] Skip files that appear in the exclusion list.
- [ ] **Unit Tests for Scanner**
  - [ ] Write tests for file scanning.
  - [ ] Test pre-checks for empty and frontmatter-only files.
  - [ ] Validate exclusion file parsing and matching.

---

## 5. Classification Module
- [ ] **Integrate with GenAI Engine**
  - [ ] Create a function to send file content and the classification prompt to the GenAI engine.
  - [ ] Simulate interaction with LangChainGo for now.
- [ ] **Response Handling**
  - [ ] Parse responses to classify files as "Empty", "Low quality/low effort", or "Good enough".
- [ ] **Unit Tests for Classification**
  - [ ] Use mocks to simulate various GenAI engine responses.
  - [ ] Verify that the function correctly interprets and returns classifications.

---

## 6. Output Module
- [ ] **Report Generation**
  - [ ] Develop a function to generate a markdown report.
  - [ ] Create three sections in the report:
    - **Empty Files**
    - **Files with Frontmatter Only**
    - **Low Quality/Low Effort Files**
- [ ] **Formatting**
  - [ ] Format each entry using Obsidian-style links (e.g., `[[link-to-page]]`).
  - [ ] Write the report to a markdown file in the target folder.
- [ ] **Unit Tests for Output**
  - [ ] Test that the report is generated with correct sections and formatting.
  - [ ] Validate that file entries match expected Obsidian link format.

---

## 7. Integration & End-to-End Wiring
- [ ] **Wire Modules in `main.go`**
  - [ ] Use the CLI module to parse command-line arguments.
  - [ ] Load the configuration using the configuration module.
  - [ ] Scan the target folder using the file scanner module.
  - [ ] For files requiring classification, call the classification module.
  - [ ] Generate the final report using the output module.
- [ ] **Error Logging and Handling**
  - [ ] Ensure all modules log errors appropriately.
  - [ ] Verify the workflow handles errors gracefully.
- [ ] **Integration Testing**
  - [ ] Write integration tests that simulate a full run on a temporary folder.
  - [ ] Include tests for all edge cases (empty files, frontmatter-only files, files to be classified).
  - [ ] Confirm that the final markdown report correctly categorizes and formats all entries.

---

## 8. Documentation & Code Quality
- [ ] **In-Code Documentation**
  - [ ] Add comments to explain the purpose of functions and modules.
- [ ] **External Documentation**
  - [ ] Create a README that explains:
    - Configuration setup
    - Command-line usage
    - Running tests
    - Interpreting the generated report
- [ ] **Code Review**
  - [ ] Review all code for adherence to Go best practices.
  - [ ] Ensure there are no orphaned or unused code sections.

---

## 9. Final Steps
- [ ] **Final Code Review and Testing**
  - [ ] Run all unit and integration tests to ensure everything works.
  - [ ] Check that the application builds and runs without errors.
- [ ] **Prepare for Deployment**
  - [ ] Update the README with build and deployment instructions.
  - [ ] Ensure logging and error messages are clear for users.

---

Happy coding! Use this checklist to track progress and ensure every step is thoroughly tested and integrated.
