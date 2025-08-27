# GoHoliday

[![CI](https://github.com/coredds/GoHoliday/workflows/CI/badge.svg)](https://github.com/coredds/GoHoliday/actions/workflows/ci.yml)
[![codecov](https://codecov.io/gh/coredds/GoHoliday/branch/main/graph/badge.svg)](https://codecov.io/gh/coredds/GoHoliday)
[![Go Version](https://img.shields.io/github/go-mod/go-version/coredds/GoHoliday?v=1.23)](https://golang.org/)
[![License](https://img.shields.io/github/license/coredds/GoHoliday)](LICENSE)

A comprehensive Go library for holiday data and business day calculations, designed as a high-performance backend for date/time applications including [ChronoGo](https://github.com/davidhintelmann/ChronoGo).

## Features

- **Multi-Country Support**: 15 countries with 200+ regional subdivisions
- **High Performance**: Sub-microsecond holiday lookups with intelligent caching  
- **ChronoGo Integration**: Native HolidayChecker interface implementation
- **Enterprise Configuration**: YAML-based configuration system with environment support
- **Cultural Accuracy**: Multi-language holiday names and historical transitions
- **Thread-Safe**: Concurrent operations with built-in safety
- **Business Intelligence**: Holiday categories, regional variations, and custom overrides

Current version: **0.3.0** with Italy, Spain, Netherlands, and South Korea support

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
| India | IN | National & state holidays | Native |
| Brazil | BR | 27 (states & federal district) | Native |
| Mexico | MX | 32 (states & federal district) | Native |
| Italy | IT | 20 (regions) | Native |
| Spain | ES | 19 (autonomous communities) | Native |
| Netherlands | NL | 12 (provinces) | Native |
| South Korea | KR | 17 (provinces & cities) | Native |

**Total**: 15 countries with 200+ regional subdivisions

The library architecture supports rapid expansion to additional countries including Singapore and additional Latin American countries.

## Features

- **Multi-Country Support**: 15 countries with 200+ regional subdivisions
- **High Performance**: Sub-microsecond holiday lookups with intelligent caching  
- **ChronoGo Integration**: Native HolidayChecker interface implementation
- **Enterprise Configuration**: YAML-based configuration system with environment support
- **Cultural Accuracy**: Multi-language holiday names and historical transitions
- **Thread-Safe**: Concurrent operations with built-in safety
- **Business Intelligence**: Holiday categories, regional variations, and custom overrides

Current version: **0.3.0** with Italy, Spain, Netherlands, and South Korea support

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

### Multi-Country and Multilingual Support

```go
package main

import (
    "fmt"
    "time"
    
    "github.com/coredds/GoHoliday"
)

func main() {
    // Brazil with Portuguese support
    br := goholidays.NewCountry("BR")
    
    // Check Carnival (complex calculation based on Easter)
    carnival2024 := time.Date(2024, 2, 13, 0, 0, 0, 0, time.UTC)
    if holiday, ok := br.IsHoliday(carnival2024); ok {
        fmt.Printf("Brazil: %s\n", holiday.Name)
        if ptName, exists := holiday.Languages["pt"]; exists {
            fmt.Printf("Em português: %s\n", ptName)
        }
    }
    
    // Mexico with Spanish support  
    mx := goholidays.NewCountry("MX")
    
    // Check Constitution Day (variable Monday holiday)
    constDay2024 := time.Date(2024, 2, 5, 0, 0, 0, 0, time.UTC)
    if holiday, ok := mx.IsHoliday(constDay2024); ok {
        fmt.Printf("Mexico: %s\n", holiday.Name)
        if esName, exists := holiday.Languages["es"]; exists {
            fmt.Printf("En español: %s\n", esName)
        }
    }
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

### Performance Characteristics
- **First lookup**: Loads and caches year data (~100μs)
- **Subsequent lookups**: O(1) cache hits (<50ns)
- **Thread-safe**: Concurrent operations with automatic memory management

## Architecture

```
GoHoliday/
├── chronogo/           # ChronoGo integration layer
├── config/             # Configuration system
├── countries/          # Country implementations
├── updater/            # Python holidays syncer
└── examples/           # Demo applications
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

# Country-specific demos
cd examples/japan-demo && go run main.go
cd examples/brazil-demo && go run main.go
cd examples/mexico-demo && go run main.go

# Multi-country comparisons
cd examples/br-mx-demo && go run main.go
cd examples/latam-comparison && go run main.go

# Performance analysis
cd examples/performance-analysis && go run main.go
```

## Attribution

This project builds upon the [Python holidays library](https://github.com/vacanza/holidays) and the work of the [Vacanza organization](https://github.com/vacanza), providing a Go-native implementation optimized for performance and ChronoGo integration.

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Contributing

We welcome contributions! Please see our [Contributing Guidelines](CONTRIBUTING.md) for details.
