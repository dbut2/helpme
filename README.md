# helpme

helpme is a command-line interface (CLI) tool that leverages AI to generate commands based on user input. It uses the GPT-3.5 Turbo model from OpenAI to understand the user's needs and provides the corresponding command(s) to execute.

## Installation

To install helpme, use the following command:

```bash
go install github.com/dbut2/helpme@latest
```

## Usage

Before using helpme, ensure you have set the `OPENAI_TOKEN` environment variable with your OpenAI API key. Then, simply run the `helpme` command followed by your query.

```bash
helpme <your_query>
```

## Example

```bash
helpme create a new directory named mydir
```

This will generate the appropriate command, such as:

```bash
mkdir mydir
```
