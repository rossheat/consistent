# consistent

Want to visualise the consistency of an LLM's answers? 

consistent is a CLI tool that allows users to ask an LLM the same yes/no question multiple times and visualise its consistency as a bar chart. 

This tool currently supports models available via the Anthropic API.

![Demo](demo.gif)

## Table of Contents

- [Installation](#installation)
- [Usage](#usage)
- [CLI Arguments](#cli-arguments)
- [Contributing](#contributing)
- [License](#license)

## Installation

To install consistent, follow these steps:

1. Ensure you have Go installed on your system (version 1.16 or later).
2. Clone this repository:
   ```
   git clone https://github.com/rossheat/consistent.git
   ```
3. Navigate to the project directory:
   ```
   cd consistent
   ```
4. Build the binary:
   ```
   go build -o consistent
   ```
5. (Optional) Move the binary to a directory in your PATH for easy access:
   ```
   sudo mv consistent /usr/local/bin/
   ```

## Usage

To use consistent, you'll need an Anthropic API key. Run the tool with the following command:

```
consistent -key YOUR_ANTHROPIC_API_KEY
```

After running the command, you'll be prompted to enter a yes/no question. The tool will then send this question to multiple instances of the specified Anthropic model and visualise the consistency of the responses.

## CLI Arguments

consistent supports the following command-line arguments:

- `-key` (required): Your Anthropic API key.
- `-debug` (optional): Start the program in debug mode. Default: false.
- `-model` (optional): The name of the Anthropic model you'd like to question. Default: "claude-3-5-sonnet-20240620".
- `-instances` (optional): The number of times your question is sent to the model API. Default: 50.
- `-delay` (optional): Milliseconds of delay between calling the Anthropic API. Default: 500.

Example with all arguments:

```
consistent -key YOUR_API_KEY -debug -model claude-3-5-sonnet-20240620 -instances 25 -delay 1000
```

## Contributing

Contributions to consistent are welcome! Please feel free to submit a Pull Request.

## License

This project is made available under the terms of the [MIT License](LICENSE.md)