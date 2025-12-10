## Archived

This project is no longer in use has been personally replaced by Claude Code

```zsh
helpme() { claude -p "$*" }
```

# helpme

A smart command-line assistant that generates shell commands using AI. Simply describe what you want to do, and `helpme` will suggest the appropriate command.

## Features

- Supports multiple AI providers:
  - Anthropic Claude (default)
  - OpenAI ChatGPT
  - Ollama (local AI models)
- Stream-based output for real-time responses
- Customizable system prompts
- Cross-platform support

## Installation

```bash
go install dbut.dev/helpme@latest
```

## Configuration

The tool can be configured using environment variables:

### General Configuration

- `HELPME_TOOL`: Choose the AI provider ("CLAUDE", "CHATGPT", or "OLLAMA")
- `HELPME_SYSTEM_PROMPT`: Custom system prompt for the AI
- `HELPME_DEBUG`: Enable debug mode (default: false)

### Provider-specific Configuration

#### Claude (Default)

```bash
export ANTHROPIC_API_TOKEN="your-api-key"
```

#### OpenAI

```bash
export OPENAI_TOKEN="your-api-key"
export OPENAI_MODEL="gpt-4" # Optional
```

#### Ollama

```bash
export HELPME_MODEL="codegemma:instruct" # Optional, this is the default model
```

## Usage

```bash
helpme <your command description>
```

### Examples

```bash
# Find large files
helpme find files larger than 1GB in the current directory

# Process text
helpme count number of lines in all python files recursively

# System administration
helpme show me system memory usage in a human readable format

# Git operations
helpme undo my last git commit but keep the changes
```

## Building from Source

1. Clone the repository:

```bash
git clone https://github.com/dbut2/helpme.git
cd helpme
```

2. Build the project:

```bash
go build
```

3. (Optional) Install locally:

```bash
go install
```

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.
