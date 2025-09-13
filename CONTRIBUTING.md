# Contributing to dotman

Thank you for your interest in contributing to dotman! This document provides guidelines and information for contributors.

## 🚀 Getting Started

### Prerequisites
- Go 1.23 or later
- Git
- Make (optional, for using Makefile commands)

### Development Setup

1. **Fork and clone the repository**
   ```bash
   git clone https://github.com/YOUR_USERNAME/dotman.git
   cd dotman
   ```

2. **Install dependencies**
   ```bash
   make deps
   # or manually:
   go mod tidy
   go mod download
   ```

3. **Run tests to ensure everything works**
   ```bash
   make test
   ```

4. **Build the project**
   ```bash
   make build
   ```

## 🧪 Testing

We have a comprehensive testing suite. Please ensure all tests pass before submitting a PR.

### Running Tests
```bash
# Run all tests
make test

# Run tests with coverage
make test-coverage

# Run tests with verbose output
make test-verbose

# Run only unit tests
make test-unit

# Run only integration tests
make test-integration
```

### Writing Tests
- Follow the existing test patterns in `*_test.go` files
- Use the `testutils` package for common testing utilities
- Test both success and failure cases
- Use descriptive test names

See [TESTING.md](TESTING.md) for detailed testing guidelines.

## 🔍 Code Quality

### Linting
We use `golangci-lint` for code quality checks:

```bash
make lint
```

### Code Style
- Follow Go standard formatting (`go fmt`)
- Use meaningful variable and function names
- Add comments for exported functions and types
- Keep functions focused and reasonably sized

### Pre-commit Checks
Before committing, ensure:
- [ ] All tests pass (`make test`)
- [ ] Code is formatted (`make fmt`)
- [ ] Linter passes (`make lint`)
- [ ] No new warnings or errors

## 📝 Submitting Changes

### Pull Request Process

1. **Create a feature branch**
   ```bash
   git checkout -b feature/your-feature-name
   ```

2. **Make your changes**
   - Write code following our style guidelines
   - Add tests for new functionality
   - Update documentation if needed

3. **Test your changes**
   ```bash
   make test
   make lint
   ```

4. **Commit your changes**
   ```bash
   git add .
   git commit -m "feat: add your feature description"
   ```

5. **Push to your fork**
   ```bash
   git push origin feature/your-feature-name
   ```

6. **Create a Pull Request**
   - Use the PR template
   - Provide a clear description
   - Link any related issues

### Commit Message Format
We follow conventional commits:
- `feat:` new features
- `fix:` bug fixes
- `docs:` documentation changes
- `style:` formatting changes
- `refactor:` code refactoring
- `test:` test additions/changes
- `chore:` maintenance tasks

Examples:
```
feat: add support for Windows symlinks
fix: resolve path resolution issue in remove command
docs: update installation instructions
```

## 🐛 Reporting Issues

### Bug Reports
Use the bug report template and include:
- Clear description of the issue
- Steps to reproduce
- Expected vs actual behavior
- Environment details (OS, Go version, etc.)
- Configuration file (if relevant)

### Feature Requests
Use the feature request template and include:
- Problem description
- Proposed solution
- Use cases
- Implementation ideas (if any)

## 🏗️ Project Structure

```
dotman/
├── cmd/                 # CLI commands
├── internal/            # Internal packages
│   ├── config/         # Configuration management
│   ├── files/          # File operations
│   ├── git/            # Git operations
│   ├── manager/        # Core business logic
│   └── testutils/      # Test utilities
├── .github/            # GitHub workflows and templates
├── Makefile           # Build and test commands
└── README.md          # Project documentation
```

## 🔄 Release Process

Releases are automated via GitHub Actions:
1. Create and push a tag: `git tag v1.0.0 && git push origin v1.0.0`
2. GitHub Actions will automatically build and create a release
3. Binaries for Linux, macOS, and Windows will be attached

## 📋 Development Guidelines

### Adding New Features
1. Start with tests (TDD approach)
2. Implement the feature
3. Update documentation
4. Ensure all tests pass
5. Submit PR with clear description

### Modifying Existing Features
1. Understand the current implementation
2. Add tests for the changes
3. Make minimal, focused changes
4. Update documentation if behavior changes
5. Ensure backward compatibility when possible

### Error Handling
- Use Go's standard error wrapping (`fmt.Errorf("...: %w", err)`)
- Provide meaningful error messages
- Log errors appropriately
- Handle edge cases gracefully

## 🤝 Community Guidelines

- Be respectful and inclusive
- Help others learn and grow
- Provide constructive feedback
- Follow the [Code of Conduct](CODE_OF_CONDUCT.md)

## 📞 Getting Help

- Check existing [Issues](https://github.com/ZonCen/dotman/issues)
- Join discussions in [Discussions](https://github.com/ZonCen/dotman/discussions)
- Create an issue for questions or problems

## 🎯 Areas for Contribution

- **Testing**: Improve test coverage
- **Documentation**: Improve docs and examples
- **Performance**: Optimize operations
- **Features**: Add new functionality
- **Bug fixes**: Fix reported issues
- **Cross-platform**: Improve Windows/macOS support

Thank you for contributing to dotman! 🎉
