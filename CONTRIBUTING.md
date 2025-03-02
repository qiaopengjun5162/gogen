# Contributing to gogen

Thank you for your interest in contributing to `gogen`! This document outlines the process for submitting bug reports,
feature requests, and code contributions.

## Reporting Issues

1. Check the [Issues](https://github.com/yourusername/gogen/issues) page to see if your issue has already been reported.
2. If not, create a new issue with:
    - A clear title and description.
    - Steps to reproduce the problem (if applicable).
    - Expected and actual behavior.
    - Any relevant logs or screenshots.

## Requesting Features

- Open an issue with the `[Feature Request]` prefix in the title.
- Describe the feature, its use case, and why it would benefit `gogen`.

## Submitting Pull Requests

1. **Fork the Repository**:
   ```bash
   git clone https://github.com/yourusername/gogen.git
   cd gogen
   ```
2. **Create a Branch**:
   ```bash
   git checkout -b feature/your-feature-name
   ```
3. **Make Changes**:
    - Ensure your code follows the existing style (e.g., consistent error handling, ANSI color usage).
    - Add or update tests in the `test/` directory if applicable.
    - Update `CHANGELOG.md` with your changes under an `[Unreleased]` section.
4. **Test Your Changes**:
   ```bash
   go build -v -ldflags "-X main.GitCommit=$(git rev-parse --short HEAD)" -o gogen ./main.go
   ./gogen --git=https://github.com/example/template
   ```
5. **Commit and Push**:
   ```bash
   git add .
   git commit -m "Add your feature or fix"
   git push origin feature/your-feature-name
   ```
6. **Open a Pull Request**:
    - Go to the repository on GitHub and create a pull request from your branch.
    - Reference any related issues (e.g., `Fixes #123`).

## Code Guidelines

- Use Go idiomatic conventions (e.g., error handling with `if err != nil`).
- Keep CLI output consistent with existing ANSI color codes (e.g., `ColorCyan` for info).
- Add comments for complex logic, especially in template processing.

## Questions?

Feel free to reach out by opening an issue or contacting the maintainers directly.

Happy contributing!
