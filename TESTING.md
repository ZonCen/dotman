# Testing Guide for dotman

This document describes the testing strategy and how to run tests for the dotman project.

## Testing Philosophy

Our testing approach follows the **"Right-sized Testing"** principle - not too granular, not too broad, but focused on:

1. **Core Functionality** - Test the essential business logic
2. **Edge Cases** - Test error conditions and boundary cases  
3. **Integration Points** - Test how components work together
4. **User Workflows** - Test complete user scenarios

## Test Structure

```
internal/
├── testutils/           # Test utilities and helpers
│   └── testutils.go
├── helpers_test.go      # Unit tests for helper functions
├── config/
│   └── config_test.go   # Unit tests for config management
├── files/
│   └── json_files_test.go # Unit tests for file operations
└── manager/
    ├── add_test.go      # Integration tests for add functionality
    └── remove_test.go   # Integration tests for remove functionality

cmd/
└── integration_test.go  # CLI integration tests

.github/workflows/
└── test.yml            # CI/CD test pipeline
```

## Test Categories

### 1. Unit Tests (47.5% coverage)
- **Purpose**: Test individual functions in isolation
- **Location**: `internal/*_test.go`
- **Examples**: Path resolution, file existence checks, JSON operations

### 2. Integration Tests (27.5% coverage)
- **Purpose**: Test complete workflows with real file system
- **Location**: `internal/manager/*_test.go`
- **Examples**: Add/remove file workflows, symlink creation

### 3. CLI Tests (21.7% coverage)
- **Purpose**: Test command-line interface behavior
- **Location**: `cmd/integration_test.go`
- **Examples**: Command execution, flag handling

## Running Tests

### Basic Commands
```bash
# Run all tests
go test ./...

# Run tests with verbose output
go test -v ./...

# Run tests with coverage
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out -o coverage.html
```

### Using Makefile
```bash
# Run all tests
make test

# Run tests with verbose output
make test-verbose

# Run tests with coverage report
make test-coverage

# Run only unit tests
make test-unit

# Run only integration tests
make test-integration

# Format code
make fmt

# Run linter
make lint
```

## Test Utilities

The `internal/testutils` package provides:

- **TestDir()** - Creates temporary directories for testing
- **CreateTestFile()** - Creates test files with content
- **CreateTestSymlink()** - Creates test symlinks
- **AssertFileExists()** - Asserts file existence
- **AssertSymlink()** - Asserts symlink properties
- **SetupTestEnvironment()** - Creates complete test environment

## Coverage Goals

Current coverage:
- **Config**: 84.6% (Excellent)
- **Files**: 82.0% (Excellent)  
- **Internal**: 47.5% (Good)
- **Manager**: 27.5% (Needs improvement)
- **CMD**: 21.7% (Needs improvement)

### Target Coverage
- **Core packages** (config, files): 80%+
- **Manager package**: 60%+
- **CLI package**: 40%+

## Writing New Tests

### Unit Test Example
```go
func TestResolvePath(t *testing.T) {
    tests := []struct {
        name     string
        input    string
        expected string
    }{
        {
            name:     "absolute path",
            input:    "/absolute/path",
            expected: "/absolute/path",
        },
        // ... more test cases
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            result, err := ResolvePath(tt.input)
            if err != nil {
                t.Errorf("ResolvePath() error = %v", err)
                return
            }
            if result != tt.expected {
                t.Errorf("ResolvePath() = %v, want %v", result, tt.expected)
            }
        })
    }
}
```

### Integration Test Example
```go
func TestAddFile(t *testing.T) {
    testDir, repoDir, symlinkDir := testutils.SetupTestEnvironment(t)
    defer testutils.CleanupTestEnvironment(t, testDir)
    
    // Create test file
    testFile := filepath.Join(symlinkDir, ".zshrc")
    testutils.CreateTestFile(t, testFile, "zsh configuration")
    
    // Create empty info.json
    infoPath := filepath.Join(repoDir, "info.json")
    testutils.CreateTestFile(t, infoPath, "{}")
    
    // Test adding file
    err := AddFile(testFile, repoDir, false)
    if err != nil {
        t.Errorf("AddFile() error = %v", err)
    }
    
    // Verify results
    repoFile := filepath.Join(repoDir, ".zshrc")
    testutils.AssertFileExists(t, repoFile)
    testutils.AssertSymlink(t, testFile, repoFile)
}
```

## Test Best Practices

1. **Use Table-Driven Tests** for multiple test cases
2. **Test Both Success and Failure Cases**
3. **Clean Up Test Resources** with defer statements
4. **Use Descriptive Test Names** that explain what's being tested
5. **Test Edge Cases** like empty inputs, non-existent files
6. **Mock External Dependencies** when appropriate
7. **Keep Tests Independent** - no shared state between tests

## CI/CD Integration

Tests run automatically on:
- Push to main/develop branches
- Pull requests to main branch

The CI pipeline:
1. Sets up Go environment
2. Installs dependencies
3. Runs all tests with coverage
4. Uploads coverage to Codecov
5. Runs linter
6. Builds the application
7. Tests the binary

## Debugging Tests

### Verbose Output
```bash
go test -v ./internal/manager/...
```

### Run Specific Test
```bash
go test -run TestAddFile ./internal/manager/...
```

### Test with Race Detection
```bash
go test -race ./...
```

### Profile Test Performance
```bash
go test -bench=. -cpuprofile=cpu.prof ./...
go tool pprof cpu.prof
```

## Future Improvements

1. **Add More Manager Tests** - Test sync, status, init functions
2. **Add Git Package Tests** - Mock git operations
3. **Add CLI Tests** - Test actual command execution
4. **Add Benchmark Tests** - Performance testing
5. **Add Property-Based Tests** - Using quickcheck-style testing
6. **Improve Test Data** - More realistic test scenarios

## Contributing

When adding new features:
1. Write tests first (TDD approach)
2. Ensure tests pass locally
3. Update this documentation if needed
4. Ensure CI passes before merging

Remember: **Good tests are an investment in code quality and maintainability!**
