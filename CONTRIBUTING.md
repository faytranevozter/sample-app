# Contributing to Go Template

Thank you for considering contributing to this project! We welcome contributions of all kinds — bug fixes, new features, documentation improvements, and more.

## Table of Contents

- [Code of Conduct](#code-of-conduct)
- [Getting Started](#getting-started)
- [How to Contribute](#how-to-contribute)
- [Development Workflow](#development-workflow)
- [Code Style](#code-style)
- [Testing](#testing)
- [Submitting Changes](#submitting-changes)
- [Reporting Bugs](#reporting-bugs)
- [Requesting Features](#requesting-features)

## Code of Conduct

Please be respectful and constructive in all interactions. We are committed to providing a welcoming and inclusive environment for everyone.

## Getting Started

1. **Fork** the repository on GitHub
2. **Clone** your fork locally:
   ```bash
   git clone https://github.com/<your-username>/go-template.git
   cd go-template
   ```
3. **Set up** your environment:
   ```bash
   cp .env.example .env
   # Edit .env with your local configuration
   ```
4. **Install** dependencies:
   ```bash
   go mod download
   ```
5. **Run** the application:
   ```bash
   make run
   ```

## How to Contribute

### Bug Fixes

- Check existing [issues](../../issues) to see if the bug has already been reported.
- If not, open a new issue using the **Bug Report** template.
- Reference the issue number in your pull request.

### New Features

- Open a **Feature Request** issue first to discuss the proposal.
- Wait for feedback before starting implementation.
- Keep changes focused and avoid unrelated modifications.

### Documentation

- Fix typos, improve clarity, or add missing documentation.
- Documentation PRs are always welcome.

## Development Workflow

1. **Create a branch** from `main`:
   ```bash
   git checkout -b feature/your-feature-name
   ```

2. **Make your changes** following the project's code style and architecture.

3. **Write tests** for new functionality.

4. **Run tests** to ensure nothing is broken:
   ```bash
   make test
   ```

5. **Commit** your changes with clear, descriptive messages:
   ```bash
   git commit -m "feat: add user profile endpoint"
   ```

6. **Push** to your fork and open a pull request.

### Commit Message Convention

We follow a simple convention for commit messages:

- `feat:` — New feature
- `fix:` — Bug fix
- `docs:` — Documentation changes
- `test:` — Adding or updating tests
- `refactor:` — Code refactoring without behavior changes
- `chore:` — Maintenance tasks (dependencies, CI, etc.)

## Code Style

- Follow standard Go conventions and idioms.
- Run `go vet ./...` to check for common issues.
- Use `gofmt` or `goimports` to format your code.
- Keep functions focused and concise.
- Use meaningful variable and function names.
- Add comments for exported functions and complex logic.

### Architecture Guidelines

This project follows **Clean Architecture**. When adding new modules:

1. Define domain entities and interfaces in `domain/`
2. Implement business logic in a dedicated module directory (e.g., `user/`)
3. Implement data access in `internal/repository/`
4. Create HTTP handlers in `internal/rest/`
5. Wire everything together in `app/main.go`

## Testing

- Write tests using Go's standard `testing` package.
- Use manual mock implementations (no external mocking frameworks).
- Place test files alongside the source code (e.g., `service.go` → `service_test.go`).
- Ensure all tests pass before submitting a PR:
  ```bash
  make test
  ```
- Check coverage with:
  ```bash
  make testcoverage
  ```

## Submitting Changes

1. Ensure your code follows the project's style guidelines.
2. Ensure all tests pass.
3. Open a pull request against the `main` branch.
4. Provide a clear description of your changes in the PR.
5. Reference any related issues.
6. Be responsive to review feedback.

### Pull Request Checklist

- [ ] Code follows the project's style guidelines
- [ ] Tests are added or updated for new functionality
- [ ] All existing tests pass (`make test`)
- [ ] Documentation is updated if needed
- [ ] Commit messages follow the convention

## Reporting Bugs

Use the [Bug Report](../../issues/new?template=bug_report.md) template when filing issues. Please include:

- A clear description of the bug
- Steps to reproduce
- Expected vs. actual behavior
- Your environment details (Go version, OS, etc.)

## Requesting Features

Use the [Feature Request](../../issues/new?template=feature_request.md) template. Please include:

- A clear description of the feature
- The problem it solves
- Any alternative solutions you've considered

---

Thank you for helping make this project better! 🎉
