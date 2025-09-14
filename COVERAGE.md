# Code Coverage Guide

This project includes comprehensive code coverage functionality to help you review total coverage and identify uncovered code paths.

## Quick Start

### Using Make Commands

The project includes several make targets for code coverage:

```bash
# Generate full coverage report with HTML
make coverage

# Show coverage summary only
make coverage-summary

# Generate HTML report only
make coverage-html

# Check if coverage meets threshold (70%)
make coverage-check
```

### Using Coverage Script

For more options, use the dedicated coverage script:

```bash
# Full coverage analysis (recommended)
./coverage.sh

# Show coverage summary only
./coverage.sh summary

# Generate HTML report only
./coverage.sh html

# Check coverage threshold
./coverage.sh check
```

## Understanding Coverage Output

### Text Summary
The text summary shows coverage per function:
```
github.com/sean9999/polity/network/connection.go:15:	Connect		75.0%
github.com/sean9999/polity/network/connection.go:25:	Disconnect	50.0%
total:							(statements)	85.2%
```

### HTML Report
The HTML report (`coverage.html`) provides:
- Visual highlighting of covered/uncovered code
- Interactive navigation through files
- Line-by-line coverage details
- Overall statistics

Open the HTML report with: `open coverage.html` (macOS) or your browser.

## Coverage Threshold

The project is configured with a **70% coverage threshold**.

- ✅ Coverage above 70% passes
- ❌ Coverage below 70% fails the check

To adjust the threshold, edit the `THRESHOLD` variable in:
- `coverage.sh` (line 8)
- Makefile coverage-check targets

## Identifying Uncovered Code

### Using the Coverage Script
```bash
./coverage.sh
```
Shows files with uncovered code paths at the bottom of the output.

### Using Go Tools Directly
```bash
# Generate coverage data
go test -coverprofile=coverage.out ./...

# Show uncovered code
go tool cover -func=coverage.out | grep -v "100.0%"

# Detailed view of a specific file
go tool cover -func=coverage.out | grep "filename.go"
```

### Using HTML Report
1. Generate HTML: `make coverage-html`
2. Open `coverage.html` in your browser
3. Red highlighting shows uncovered lines
4. Green highlighting shows covered lines

## Integration with Development Workflow

### Pre-commit Checks
Add to your development routine:
```bash
# Run tests and check coverage before committing
make test && make coverage-check
```

### IDE Integration
Most Go IDEs support coverage visualization:
- **VS Code**: Use Go extension's coverage features
- **GoLand**: Built-in coverage runner
- **Vim/Neovim**: Use vim-go with coverage support

## Advanced Usage

### Coverage by Package
```bash
go test -coverprofile=coverage.out ./...
go tool cover -func=coverage.out | sort -k3 -nr
```

### Coverage with Specific Packages Only
```bash
go test -coverprofile=coverage.out ./network/...
go tool cover -html=coverage.out -o network-coverage.html
```

### Detailed Coverage Modes
```bash
# Count mode (default)
go test -covermode=count -coverprofile=coverage.out ./...

# Atomic mode (for concurrent tests)
go test -covermode=atomic -coverprofile=coverage.out ./...
```

## Generated Files

The following files are created and should not be committed:
- `coverage.out` - Raw coverage data
- `coverage.html` - HTML coverage report

These are automatically ignored via `.gitignore`.

## Troubleshooting

### No Test Files Found
If you see "no test files found":
1. Ensure test files end with `_test.go`
2. Check test files are in the same package as code being tested

### Low Coverage Warnings
If coverage is unexpectedly low:
1. Check if all test files are being discovered
2. Verify test functions start with `Test`
3. Use the HTML report to identify specific uncovered lines

### Build Failures
If tests fail to build:
1. Run `go mod tidy` to ensure dependencies are correct
2. Check for import issues in test files
3. Verify Go version compatibility

## Best Practices

1. **Aim for meaningful coverage** - Focus on critical code paths
2. **Use HTML reports** - Visual feedback helps identify gaps
3. **Regular monitoring** - Check coverage with each significant change
4. **Don't chase 100%** - Focus on testing business logic and error paths
5. **Document exceptions** - Some code (like simple getters) may not need tests

## CI/CD Integration

To integrate coverage into your CI pipeline, add these commands:

```yaml
# Example for GitHub Actions
- name: Run tests with coverage
  run: make coverage-check

# Example for GitLab CI
test:
  script:
    - make coverage-check
```

The `coverage-check` target will fail the build if coverage is below the threshold.