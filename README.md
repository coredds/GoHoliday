# GoHoliday

[![CI](https://github.com/davidhintelmann/GoHoliday/workflows/CI/badge.svg)](https://github.com/davidhintelmann/GoHoliday/actions/workflows/ci.yml)
[![codecov](https://codecov.io/gh/davidhintelmann/GoHoliday/branch/main/graph/badge.svg)](https://codecov.io/gh/davidhintelmann/GoHoliday)
[![Go Version](https://img.shields.io/github/go-mod/go-version/davidhintelmann/GoHoliday)](https://golang.org/)
[![License](https://img.shields.io/github/license/davidhintelmann/GoHoliday)](LICENSE)
[![Go Report Card](https://goreportcard.com/badge/github.com/davidhintelmann/GoHoliday)](https://goreportcard.com/report/github.com/davidhintelmann/GoHoliday)

A comprehensive Go library for holiday data and business day calculations, designed as a high-performance backend for date/time applications including [ChronoGo](https://github.com/davidhintelmann/ChronoGo).

## Project Status

**Current Version**: 0.1.2  
**Supported Countries**: 8 countries with comprehensive regional subdivision support  
**Performance**: Sub-microsecond holiday lookups with O(1) caching  
**Integration**: Native ChronoGo HolidayChecker interface implementation  

### Recent Additions
- **Japan (JP)** support with 16 public holidays including cultural accuracy and Emperor transition handling
- Enhanced ChronoGo integration with FastCountryChecker performance layer
- Comprehensive test coverage with benchmark validation

## Objectives

1. **Primary**: Provide enterprise-grade holiday data for ChronoGo business day calculations
2. **Secondary**: Support international business operations with accurate regional holiday data  
3. **Tertiary**: Enable rapid expansion to additional countries via Python holidays library integration

## Supported Countries

| Country | Code | Subdivisions | Implementation |
|---------|------|-------------|----------------|
| United States | US | 56 (states, territories) | Native |
| United Kingdom | GB | 4 (England, Scotland, Wales, NI) | Native |
| Canada | CA | 13 (provinces, territories) | Native |
| Australia | AU | 8 (states, territories) | Native |
| New Zealand | NZ | 17 (regions) | Native |
| Germany | DE | 16 (states) | Native |
| France | FR | Regions & territories | Native |
| Japan | JP | National public holidays | Native |

**Total**: 8 countries with 100+ regional subdivisions

### Future Expansion
The library architecture supports rapid country expansion using data from the Python holidays ecosystem. Target countries include India, Brazil, Mexico, Italy, Spain, Singapore, and South Korea.

## Features

- **Enterprise Configuration**: YAML-based configuration system with environment support
- **Multi-Country Support**: 8 countries with regional subdivisions
- **High Performance**: Sub-microsecond holiday lookups with intelligent caching
- **ChronoGo Integration**: Native HolidayChecker interface implementation
- **Business Intelligence**: Holiday categories, regional variations, and custom overrides
- **Thread-Safe**: Concurrent operations with built-in safety
- **Cultural Accuracy**: Multi-language holiday names and historical transitions

## Quick Start

### ChronoGo Integration

```go
package main

import (
    "fmt"
    "time"
    
    "github.com/coredds/GoHoliday/chronogo"
    chronogo_lib "github.com/davidhintelmann/chronogo"
)

func main() {
    // Create optimized holiday checker for ChronoGo
    holidayChecker := chronogo.Checker("US")
    
    // Use with ChronoGo for business day calculations
    dt := chronogo_lib.Now()
    
    // Check if today is a holiday
    if holidayChecker.IsHoliday(dt.Time) {
        name := holidayChecker.GetHolidayName(dt.Time)
        fmt.Printf("Today is a holiday: %s\n", name)
    }
    
    // Calculate next business day (ChronoGo handles this with the checker)
    nextBusiness := dt.NextBusinessDay(holidayChecker)
    fmt.Printf("Next business day: %s\n", nextBusiness.Format("2006-01-02"))
}
```

### Direct Library Usage

```go
package main

import (
    "fmt"
    "time"
    
    "github.com/coredds/GoHoliday"
)

func main() {
    // Create a holiday provider for Japan
    jp := goholidays.NewCountry("JP")
    
    // Check if a specific date is a holiday
    newYear := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
    if holiday, ok := jp.IsHoliday(newYear); ok {
        fmt.Printf("%s is %s\n", newYear.Format("2006-01-02"), holiday.Name)
        // Output: 2024-01-01 is New Year's Day
        
        // Check for Japanese translation
        if japaneseName, exists := holiday.Languages["ja"]; exists {
            fmt.Printf("In Japanese: %s\n", japaneseName)
            // Output: In Japanese: 元日
        }
    }
    
    // Get all holidays for a year
    holidays := jp.HolidaysForYear(2024)
    fmt.Printf("Japan has %d public holidays in 2024\n", len(holidays))
}
```

## Installation

```bash
go get github.com/coredds/GoHoliday
```

For ChronoGo integration:
```bash
go get github.com/coredds/GoHoliday/chronogo
```

## Performance

| Operation | Duration | Throughput |
|-----------|----------|------------|
| Holiday Check | ~50ns | 24M ops/sec |
| Batch Operations | ~60ns/date | 16M ops/sec |
| Range Counting | ~26μs | 40K ranges/sec |

### Caching Benefits
- First check: loads and caches year data (~100μs)
- Subsequent checks: O(1) cache hits (<50ns)
- Thread-safe concurrent operations
- Automatic memory management

## Architecture

```
GoHoliday/
├── chronogo/           # ChronoGo integration layer
│   ├── integration.go  # FastCountryChecker implementation
│   └── integration_test.go
├── config/             # Configuration system
│   ├── config.go       # YAML configuration
│   ├── manager.go      # Holiday manager
│   └── *.yaml          # Environment configs
├── countries/          # Country implementations
│   ├── base.go         # Base provider interface
│   ├── us.go           # United States
│   ├── jp.go           # Japan
│   └── ...             # Other countries
└── updater/            # Python holidays syncer
    └── python_ast_parser.go
```

## Testing

```bash
# Run all tests
go test ./...

# Test specific components
go test ./chronogo -v          # ChronoGo integration tests
go test ./config -v            # Configuration system tests
go test ./countries -v         # Country provider tests

# Run benchmarks
go test ./chronogo -bench=.    # Performance benchmarks
```

## Development

### Setup
```bash
git clone https://github.com/coredds/GoHoliday
cd GoHoliday
go mod download
go test ./...
```

### Examples
```bash
# ChronoGo integration demo
cd examples/chronogo && go run main.go

# Japan holidays demo  
cd examples/chronogo-japan && go run main.go
```

## Attribution

This project builds upon and extends the excellent work of the Python holidays ecosystem:

- **Python holidays library**: [github.com/python-holidays/python-holidays](https://github.com/python-holidays/python-holidays) - Comprehensive holiday data for 200+ countries
- **Vacanza organization**: [github.com/vacanza](https://github.com/vacanza) - Current maintainers of the Python holidays ecosystem including the vacanza/python-holidays fork
- **Original python-holidays contributors**: The community that built the foundational holiday calculation algorithms and country-specific implementations

GoHoliday provides a Go-native implementation optimized for performance and ChronoGo integration while maintaining compatibility with the data structures and calculations established by the Python holidays community.

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Contributing

We welcome contributions! Please see our [Contributing Guidelines](CONTRIBUTING.md) for details.

---

**Ready to enhance your Go applications with comprehensive holiday data?**

```bash
go get github.com/coredds/GoHoliday
```
