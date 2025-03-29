Below is a comprehensive blueprint that starts from a high-level plan and iteratively breaks the project into small, manageable, testable pieces. Each section includes context, a detailed description, and ends by wiring together previous parts. At the end, you’ll find a series of prompts (each wrapped in markdown code blocks tagged as text) that you can feed to a code-generation LLM in sequence.

---

## 1. High-Level Blueprint

### 1.1. Project Setup and Structure
- **Initialize the Project:**  
  - Create a Go module (using `go mod init`) and set up the repository structure.
- **Directory Structure:**  
  - **Main:** `main.go`
  - **Modules/Packages:**  
    - `cli` – for Cobra command definitions and command-line argument parsing.
    - `config` – for YAML configuration loading via Viper.
    - `scanner` – for file scanning, pre-checks (empty, frontmatter-only), and exclusion handling.
    - `classification` – for integrating with the GenAI engine (via LangChainGo) and handling classification responses.
    - `output` – for generating the markdown report.

### 1.2. Core Functionality Overview
- **Input & Scanning:**  
  - Accept a target folder and scan for Markdown (`.md`) files.
  - Skip excluded directories and files listed in `quality_exclude_links.md`.
  - Check files for being empty or containing only frontmatter.
- **Classification Process:**  
  - For files with content, call the GenAI engine using a configurable prompt.
  - Classify files as “Empty”, “Low quality/low effort”, or “Good enough.”
- **Output Generation:**  
  - Produce a markdown report in the target folder with Obsidian link formatting.
- **Error Handling:**  
  - Log configuration, I/O, and API errors.
- **Testing:**  
  - Write unit tests for each module and integration tests for end-to-end functionality.

---

## 2. Iterative Breakdown Into Chunks

### Step 1: Initial Project Setup
- **Task 1.1:** Initialize a new Go module.
- **Task 1.2:** Create the repository structure with directories for `cli`, `config`, `scanner`, `classification`, and `output`.
- **Task 1.3:** Write a minimal `main.go` that calls a dummy initialization function.

### Step 2: Configuration Module
- **Task 2.1:** Define a YAML configuration structure that includes:
  - AI engine settings (`url`, `model`)
  - Scan settings (`file_extension`, `exclude_directories`)
  - Prompt configuration (`quality_classification_prompt`)
  - Exclusion file path.
- **Task 2.2:** Use Viper to load and validate the configuration from a file.
- **Task 2.3:** Write unit tests to verify configuration loading and default values.

### Step 3: CLI Module Using Cobra
- **Task 3.1:** Set up a basic Cobra command that accepts a target folder as a flag or argument.
- **Task 3.2:** Integrate the configuration module to load the YAML file.
- **Task 3.3:** Write tests to verify CLI flag parsing and configuration loading.

### Step 4: File Scanner Module
- **Task 4.1:** Implement recursive file scanning in the target folder for `.md` files.
- **Task 4.2:** Exclude directories based on the configuration.
- **Task 4.3:** Implement pre-checks:
  - **Empty File Check:** Verify files after trimming whitespace.
  - **Frontmatter-only Check:** Detect if a file contains only a YAML block.
- **Task 4.4:** Parse `quality_exclude_links.md` to extract Obsidian links and exclude these files.
- **Task 4.5:** Write unit tests for file scanning, pre-checks, and exclusion handling.

### Step 5: Classification Module
- **Task 5.1:** Create a function to send file content along with the prompt from the config to the GenAI engine (simulate using LangChainGo).
- **Task 5.2:** Implement parsing of the GenAI response (expected: Empty, Low quality, Good enough).
- **Task 5.3:** Write unit tests simulating different responses (using mocks).

### Step 6: Output Module
- **Task 6.1:** Develop a function to generate a markdown report with three sections:
  - Empty files
  - Files with frontmatter only
  - Low quality/low effort files
- **Task 6.2:** Format each entry using the Obsidian link syntax (`[[link-to-page]]`).
- **Task 6.3:** Write tests to validate the report content and format.

### Step 7: Integration and Wiring Everything Together
- **Task 7.1:** In `main.go`, wire the CLI module, configuration loading, file scanning, classification, and output generation.
- **Task 7.2:** Ensure proper error logging and handling through the full workflow.
- **Task 7.3:** Write integration tests simulating a full run on a temporary folder with test Markdown files.

---

## 3. Final Series of Prompts for a Code-Generation LLM

Each prompt is designed to build on previous work, include testing, and ensure no orphaned code is left unintegrated. You can feed these one-by-one to generate the full project.

---

### **Prompt 1: Project Setup & Minimal Main**

```text
Create a minimal Go module for the CLI tool. Initialize a new Go module (e.g., with "go mod init ratemykb") and set up the repository structure with the following directories:
- cli
- config
- scanner
- classification
- output

Then, write a minimal `main.go` file that imports a dummy initialization function from the CLI module (for example, a function called `Execute()`) and calls it. Ensure the project builds successfully. Include comments to explain each section.

Refer to todo.md and make sure thiss step has been validated.
```

---

### **Prompt 2: Configuration Module**

```text
Develop the configuration module in the `config` package. Define a configuration struct that matches the following YAML structure:
- ai_engine:
    - url (default: "http://localhost:11434/")
    - model (default: "gemma:12b")
- scan_settings:
    - file_extension (set to ".md")
    - exclude_directories (a list of directories to ignore)
- prompt_config:
    - quality_classification_prompt (the GenAI prompt)
- exclusion_file:
    - path (default: "quality_exclude_links.md" in the target folder)

Use Viper to load this YAML configuration from a file. Include error handling for missing or invalid configuration files. Also, add unit tests that verify configuration loading, defaults being applied, and proper error reporting. End the prompt by wiring this configuration module so that it can be imported by other modules.

Refer to todo.md and make sure thiss step has been validated.
```

---

### **Prompt 3: CLI Module with Cobra**

```text
Create the CLI module in the `cli` package using Cobra. Implement a root command that accepts a target folder as an argument or flag. In the command's run function, load the configuration using the configuration module from the previous step. Ensure proper error handling if the configuration fails to load. Write tests for the CLI module to verify that command-line arguments are parsed correctly and that the configuration is loaded. End the prompt with a call to the configuration module to ensure integration.
```

---

### **Prompt 4: File Scanner Module**

```text
Develop the file scanner module in the `scanner` package. This module should:
1. Recursively scan the specified target folder for files with the `.md` extension.
2. Exclude directories specified in the configuration.
3. For each file, perform two pre-checks:
   - **Empty File Check:** Read file content, trim whitespace, and classify as "empty" if no content remains.
   - **Frontmatter-only Check:** Detect if a file contains only a YAML block delimited by "---" markers, with no content outside.
4. Parse an exclusion file (`quality_exclude_links.md`) to extract Obsidian links (formatted as `[[link-to-page]]`) and skip files that are already flagged.

Include unit tests for each of these functionalities. End the prompt by ensuring that the scanner module outputs a list of file objects along with their pre-check statuses, ready for further processing.
```

---

### **Prompt 5: Classification Module**

```text
Implement the classification module in the `classification` package. This module should:
1. Define a function that accepts file content and the GenAI prompt (from the configuration).
2. Simulate sending this data to a GenAI engine (using LangChainGo) to classify the file as one of: "Empty", "Low quality/low effort", or "Good enough".
3. Parse and return the classification result.

Write unit tests that simulate different responses from the GenAI engine using mocking. End the prompt with the module exposing a function that can be called from the main workflow.
```

---

### **Prompt 6: Output Module**

```text
Create the output module in the `output` package. This module should:
1. Accept categorized results from the scanner and classification modules.
2. Generate a markdown report that includes three sections:
   - **Empty Files**
   - **Files with Frontmatter Only**
   - **Low Quality/Low Effort Files**
3. Format each file entry using the Obsidian link format (e.g., `[[link-to-page]]`).
4. Write the report to a markdown file in the root of the target folder.

Include unit tests to verify the formatting and content of the report. End the prompt by ensuring that the output module integrates with the rest of the project.
```

---

### **Prompt 7: Integration & End-to-End Wiring**

```text
Wire together all the modules in `main.go` to create a fully integrated CLI application. The integration should follow these steps:
1. Use the CLI module to parse the command-line arguments and load the target folder.
2. Load the configuration using the configuration module.
3. Scan the target folder using the file scanner module, applying pre-checks and exclusion logic.
4. For files that require further classification, call the classification module to determine their quality.
5. Generate a final markdown report using the output module, categorizing files as:
   - Empty Files
   - Files with Frontmatter Only
   - Low Quality/Low Effort Files
6. Ensure that all errors are properly logged and that the process is testable end-to-end.

Write integration tests that simulate a run on a temporary folder structure with sample Markdown files covering all edge cases. End the prompt with a summary of the integration and ensure that there is no orphaned code.
```

---

These prompts are designed to be incremental, test-driven, and integrated—ensuring that each step builds upon the previous ones while adhering to best practices and a solid development/testing strategy.