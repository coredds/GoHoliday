# GoHoliday Examples

This directory contains examples demonstrating the features and usage of the GoHoliday library.

## Directory Structure

### 1. Quick Start
- `01_quickstart/`: A simple 5-minute introduction to GoHoliday's basic features
  - Basic holiday lookup
  - Year calendar
  - Multi-language support

### 2. Core Features
- `02_features/basic/`: Comprehensive examples of basic functionality
  - Holiday checks
  - Date ranges
  - Categories
- `02_features/multi_country/`: Cross-country holiday comparisons
  - Holiday differences
  - Regional variations
  - Calendar differences
- `02_features/business_days/`: Business day calculations
  - Working day counting
  - Business day adjustments
  - Holiday impact on schedules
- `02_features/localization/`: Language and regional features
  - Multi-language support
  - Regional holiday names
  - Country-specific formatting

### 3. Advanced Features
- `03_advanced/configuration/`: Configuration management
  - Custom settings
  - Environment-specific configs
- `03_advanced/sync/`: Update system usage
  - Holiday data updates
  - Synchronization features
- `03_advanced/performance/`: Performance testing
  - Benchmarks
  - Optimization examples

### Development Tools
- `dev_tools/ast_parser/`: Tools for library maintainers
  - AST parsing utilities
  - Development helpers

## Running the Examples

Each example can be run independently using:

```bash
go run examples/[directory]/main.go
```

For instance, to run the quickstart example:

```bash
go run examples/01_quickstart/main.go
```

## Learning Path

1. Start with `01_quickstart` to get a basic understanding
2. Move through the features in `02_features` to learn core functionality
3. Explore advanced features in `03_advanced` as needed
4. Check `dev_tools` if you're contributing to the library

## Notes

- Examples are self-contained and include necessary imports
- Each example demonstrates specific features and use cases
- Comments explain the code and concepts being demonstrated
- Advanced examples build on concepts from basic examples

