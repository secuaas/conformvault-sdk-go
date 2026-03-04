# Contributing to ConformVault Go SDK

Thank you for your interest in contributing!

## Getting Started

1. Fork the repository
2. Clone your fork: `git clone https://github.com/YOUR_USERNAME/conformvault-sdk-go.git`
3. Create a feature branch: `git checkout -b feature/your-feature`
4. Make your changes
5. Run tests: `go test ./...`
6. Run linter: `golangci-lint run`
7. Commit and push your changes
8. Open a Pull Request

## Development

```bash
# Build
go build ./...

# Test
go test ./... -race

# Lint
golangci-lint run
```

## Code Style

- Follow standard Go conventions (`gofmt`, `go vet`)
- Add godoc comments on all exported types and functions
- Write tests for new functionality

## Reporting Issues

Please open an issue on GitHub with:
- A clear description of the problem
- Steps to reproduce
- Expected vs actual behavior
- Go version and OS

## License

By contributing, you agree that your contributions will be licensed under the MIT License.
