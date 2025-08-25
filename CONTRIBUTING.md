# Contributing to GoHolidays

Thank you for your interest in contributing to GoHolidays! This document provides guidelines for contributing to the project.

## Development Setup

1. **Prerequisites**
   - Go 1.21 or later
   - Git

2. **Clone the repository**
   ```bash
   git clone https://github.com/your-username/goholidays.git
   cd goholidays
   ```

3. **Install dependencies**
   ```bash
   go mod download
   ```

4. **Run tests**
   ```bash
   go test ./...
   ```

## Project Structure

```
GoHolidays/
├── goholidays.go           # Main library code
├── goholidays_test.go      # Main tests
├── countries/              # Country-specific implementations
│   ├── base.go            # Base provider and utilities
│   ├── us.go              # United States holidays
│   └── ...                # Other countries
├── updater/               # Data synchronization
│   └── sync.go            # Python holidays sync
├── cmd/                   # Command-line tools
│   └── goholidays/        # CLI application
├── examples/              # Usage examples
└── docs/                  # Documentation
```

## Adding New Countries

To add support for a new country:

1. **Create a country file** in the `countries/` directory (e.g., `countries/ca.go` for Canada)

2. **Implement the HolidayProvider interface**:
   ```go
   type CAProvider struct {
       *BaseProvider
   }
   
   func NewCAProvider() *CAProvider {
       base := NewBaseProvider("CA")
       // Configure subdivisions, categories, etc.
       return &CAProvider{BaseProvider: base}
   }
   
   func (ca *CAProvider) LoadHolidays(year int) map[time.Time]*Holiday {
       // Implement holiday calculations
   }
   ```

3. **Add tests** in a corresponding test file

4. **Update the main library** to recognize the new country

5. **Update documentation** including supported country lists

## Holiday Calculation Guidelines

- **Fixed date holidays**: Use `time.Date(year, month, day, 0, 0, 0, 0, time.UTC)`
- **Variable holidays**: Use helper functions like `NthWeekdayOfMonth()`, `EasterSunday()`, etc.
- **Observed dates**: Use `BaseProvider.CalculateObservedDate()` for weekend shifts
- **Multi-language names**: Always provide at least English ("en") names
- **Categories**: Use standard categories (public, bank, school, etc.)

## Testing

- Write unit tests for all new functionality
- Include edge cases and boundary conditions
- Test multiple years to ensure calculations work correctly
- Add benchmark tests for performance-critical code

Example test structure:
```go
func TestCAHolidays(t *testing.T) {
    ca := NewCAProvider()
    holidays := ca.LoadHolidays(2024)
    
    // Test specific holiday
    canadaDay := time.Date(2024, 7, 1, 0, 0, 0, 0, time.UTC)
    if _, exists := holidays[canadaDay]; !exists {
        t.Error("Canada Day should be a holiday")
    }
}
```

## Code Style

- Follow standard Go conventions
- Use `gofmt` to format code
- Add meaningful comments for complex calculations
- Use descriptive variable and function names
- Keep functions focused and single-purpose

## Documentation

- Update README.md with new features
- Add inline documentation for public APIs
- Include usage examples for new functionality
- Update CLI help text when adding new options

## Submitting Changes

1. **Fork the repository**
2. **Create a feature branch**: `git checkout -b feature/new-country`
3. **Make your changes** following the guidelines above
4. **Add tests** for new functionality
5. **Run the test suite**: `go test ./...`
6. **Run linting**: `go vet ./...`
7. **Commit your changes** with descriptive messages
8. **Push to your fork**: `git push origin feature/new-country`
9. **Create a Pull Request** with a clear description of changes

## Pull Request Guidelines

- Provide a clear description of what the PR does
- Reference any related issues
- Include test results
- Update documentation as needed
- Keep PRs focused on a single feature or fix

## Data Sources

When adding new countries, use official government sources for holiday information:

- Government websites
- Official calendar publications
- Legal documents defining public holidays
- The Python holidays library as a reference (but verify independently)

## Performance Considerations

- Holiday lookups should be O(1) time complexity
- Memory usage should be minimal
- Consider lazy loading for countries with many subdivisions
- Benchmark performance-critical changes

## Questions or Issues?

- Check existing issues on GitHub
- Create a new issue for bugs or feature requests
- Use discussions for questions about implementation

Thank you for contributing to GoHolidays!
