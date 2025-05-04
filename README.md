# What-To-Do

<div align="center">

[![License: MIT](https://img.shields.io/badge/License-MIT-blue.svg)](https://opensource.org/licenses/MIT)

</div>

A sleek, terminal-based bash command alias manager with a modern TUI. Simplify your workflow by organizing, managing, and executing your most used terminal commands.

## Features

- **Command Management**: Create, update, and delete bash command aliases with ease
- **Spotlight Feature**: Highlight your most important commands for quick access
- **Intuitive TUI**: Beautiful terminal interface built with Bubble Tea
- **Fast Execution**: Run your saved commands directly from the interface

## Installation

```bash
# Clone the repository
git clone https://github.com/yourusername/what-to-do.git
cd what-to-do

# Build the application
make build

# Install the binary (optional)
make install
```

## Usage

Run the application:

```bash
what-to-do
```

### Keyboard Shortcuts

- `a` or `n`: Add a new command
- Arrow keys: Navigate through the interface
- `Enter`: Execute selected command
- `d`: Delete selected command
- `u`: Update selected command
- `Esc`: Go back from form to main view
- `q` or `Ctrl+C`: Quit the application

## Development

```bash
# Run tests
make test

# Build for development
make build-dev
```

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.