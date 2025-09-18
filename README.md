# GoHoliday

[![CI](https://github.com/coredds/GoHoliday/workflows/CI/badge.svg)](https://github.com/coredds/GoHoliday/actions/workflows/ci.yml)
[![codecov](https://codecov.io/gh/coredds/GoHoliday/branch/main/graph/badge.svg)](https://codecov.io/gh/coredds/GoHoliday)
[![Go Version](https://img.shields.io/github/go-mod/go-version/coredds/GoHoliday?v=1.23)](https://golang.org/)
[![License](https://img.shields.io/github/license/coredds/GoHoliday)](LICENSE)

A comprehensive Go library for holiday data and business day calculations. Provides high-performance holiday checking with multi-country support, designed for integration with date/time applications including [ChronoGo](https://github.com/davidhintelmann/ChronoGo).

**Current Version**: 0.5.3

## Features

- **33 Countries**: Complete coverage with 600+ regional subdivisions
- **High Performance**: Sub-microsecond holiday lookups with intelligent caching
- **Multi-Language**: Native language support for holiday names
- **Thread-Safe**: Concurrent operations with built-in safety
- **ChronoGo Integration**: Native HolidayChecker interface implementation
- **Enterprise Configuration**: YAML-based configuration with environment support
- **Regional Variations**: State, province, and regional holiday support
- **Historical Accuracy**: Proper handling of holiday transitions and changes
- **Robust Error Handling**: Structured errors with context support and validation
- **Backward Compatible**: V2 API alongside original API

## Supported Countries

| Country | Code | Subdivisions | Languages | Features |
|---------|------|-------------|-----------|----------|
| United States | US | 56 (states, territories) | EN, ES | Federal and state holidays |
| United Kingdom | GB | 4 (England, Scotland, Wales, NI) | EN | Bank holidays, regional variations |
| Canada | CA | 13 (provinces, territories) | EN, FR | Federal and provincial holidays |
| Australia | AU | 8 (states, territories) | EN | National and state holidays |
| New Zealand | NZ | 17 (regions) | EN, MI | National and regional holidays |
| Germany | DE | 16 (states) | DE, EN | Federal and state holidays |
| France | FR | Regions & territories | FR, EN | National and regional holidays |
| Japan | JP | National | JA, EN | Public holidays |
| India | IN | 36 (states, territories) | HI, EN | Multi-religious festivals |
| Brazil | BR | 27 (states, federal district) | PT, EN | National and state holidays |
| Mexico | MX | 32 (states, federal district) | ES, EN | National and state holidays |
| Italy | IT | 20 (regions) | IT, EN | National and patron saint holidays |
| Spain | ES | 19 (autonomous communities) | ES, EN | National and regional holidays |
| Netherlands | NL | 12 (provinces) | NL, EN | National holidays |
| South Korea | KR | 17 (provinces, cities) | KO, EN | National holidays |
| Singapore | SG | 5 (regions) | EN, ZH, MS, TA | Multi-cultural holidays |
| Switzerland | CH | 26 (cantons) | DE, FR, IT, RM | Federal and cantonal holidays |
| Sweden | SE | 21 (counties) | SV, EN | National holidays |
| Argentina | AR | 24 (provinces) | ES, EN | National and provincial holidays |
| Thailand | TH | 77 (provinces) | TH, EN | Buddhist and royal holidays |
| Norway | NO | 11 (counties) | NO, EN | National holidays |
| Turkey | TR | 81 (provinces) | TR, EN | Secular and religious holidays |
| Russia | RU | 85 (federal subjects) | RU, EN | Orthodox calendar holidays |
| Indonesia | ID | 38 (provinces) | ID, EN | Multi-religious holidays |
| Belgium | BE | 3 (regions) | NL, FR, DE, EN | National and regional holidays |
| Poland | PL | 16 (voivodeships) | PL, EN | National and regional holidays |
| Austria | AT | 9 (states) | DE, EN | Federal and state holidays |
| Finland | FI | 19 (regions) | FI, SV, EN | National holidays |
| China | CN | 34 (provinces, regions) | ZH, EN | Lunar calendar holidays |
| Portugal | PT | 20 (districts, regions) | PT, EN | National and regional holidays |

## Installation

```bash
go get github.com/coredds/GoHoliday
```

For ChronoGo integration:
```bash
go get github.com/coredds/GoHoliday/chronogo
```

## Quick Start

### Basic Usage

```go
package main

import (
    "fmt"
    "time"
    
    "github.com/coredds/GoHoliday"
)

func main() {
    // Create country instance
    us := goholidays.NewCountry("US")
    
    // Check if date is a holiday
    date := time.Date(2024, 7, 4, 0, 0, 0, 0, time.UTC)
    if holiday, isHoliday := us.IsHoliday(date); isHoliday {
        fmt.Printf("Holiday: %s\n", holiday.Name)
        fmt.Printf("Category: %s\n", holiday.Category)
    }
    
    // Get all holidays for a year
    holidays := us.GetHolidays(2024)
    fmt.Printf("Found %d holidays in 2024\n", len(holidays))
}
```

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
    // Create holiday checker for ChronoGo
    holidayChecker := chronogo.Checker("US")
    
    // Use with ChronoGo for business day calculations
    dt := chronogo_lib.Now()
    
    // Check if today is a holiday
    if holidayChecker.IsHoliday(dt.Time) {
        name := holidayChecker.GetHolidayName(dt.Time)
        fmt.Printf("Today is a holiday: %s\n", name)
    }
    
    // Calculate next business day
    nextBusiness := dt.NextBusinessDay(holidayChecker)
    fmt.Printf("Next business day: %s\n", nextBusiness.Format("2006-01-02"))
}
```

### Multi-Language Support

```go
package main

import (
    "fmt"
    "time"
    
    "github.com/coredds/GoHoliday"
)

func main() {
    // Singapore with 4 official languages
    sg := goholidays.NewCountry("SG")
    
    date := time.Date(2024, 2, 10, 0, 0, 0, 0, time.UTC)
    if holiday, ok := sg.IsHoliday(date); ok {
        fmt.Printf("English: %s\n", holiday.Name)
        
        // Display in all available languages
        for lang, name := range holiday.Languages {
            fmt.Printf("%s: %s\n", lang, name)
        }
    }
}
```

### Error Handling

```go
package main

import (
    "context"
    "errors"
    "fmt"
    "time"
    
    "github.com/coredds/GoHoliday"
)

func main() {
    // Create country with validation
    country, err := goholidays.NewCountryWithError("US")
    if err != nil {
        var holidayErr *goholidays.HolidayError
        if errors.As(err, &holidayErr) {
            switch holidayErr.Code {
            case goholidays.ErrInvalidCountry:
                fmt.Printf("Unsupported country: %s\n", holidayErr.Country)
            default:
                fmt.Printf("Error: %v\n", err)
            }
        }
        return
    }
    
    // Check holiday with error handling
    date := time.Date(2024, 7, 4, 0, 0, 0, 0, time.UTC)
    holiday, isHoliday, err := country.IsHolidayWithError(date)
    if err != nil {
        fmt.Printf("Error: %v\n", err)
        return
    }
    
    if isHoliday {
        fmt.Printf("Holiday: %s\n", holiday.Name)
    }
    
    // Use context for cancellation/timeout
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()
    
    holidays, err := country.HolidaysForYearWithContext(ctx, 2024)
    if err != nil {
        if goholidays.IsContextCancelled(err) {
            fmt.Println("Operation timed out")
        } else {
            fmt.Printf("Error: %v\n", err)
        }
        return
    }
    
    fmt.Printf("Found %d holidays\n", len(holidays))
}
```

## Performance

| Operation | Duration | Throughput |
|-----------|----------|------------|
| Holiday Check | ~50ns | 24M ops/sec |
| Batch Operations | ~60ns/date | 16M ops/sec |
| Range Counting | ~26μs | 40K ranges/sec |

### Characteristics

- **First lookup**: Loads and caches year data (~100μs)
- **Subsequent lookups**: O(1) cache hits (<50ns)
- **Thread-safe**: Concurrent operations with automatic memory management
- **Memory efficient**: Lazy loading and intelligent caching

## Configuration

GoHoliday supports YAML-based configuration for enterprise deployments:

```yaml
# config/goholidays.yaml
countries:
  US:
    subdivisions: ["CA", "NY", "TX"]
    categories: ["federal", "state"]
    custom_holidays:
      - name: "Company Day"
        date: "2024-03-15"
        category: "company"

performance:
  cache_size: 1000
  preload_years: [2024, 2025]
```

## Architecture

```
GoHoliday/
├── chronogo/           # ChronoGo integration layer
├── config/             # Configuration system
├── countries/          # Country implementations
├── updater/            # Python holidays syncer
└── examples/           # Demo applications
```

## Error Handling

### API Methods

**Original API (backward compatible):**
- `NewCountry(countryCode string, options ...CountryOptions) *Country`
- `IsHoliday(date time.Time) (*Holiday, bool)`
- `HolidaysForYear(year int) map[time.Time]*Holiday`
- `HolidaysForDateRange(start, end time.Time) map[time.Time]*Holiday`

**Enhanced API (with error handling):**
- `NewCountryWithError(countryCode string, options ...CountryOptions) (*Country, error)`
- `IsHolidayWithError(date time.Time) (*Holiday, bool, error)`
- `HolidaysForYearWithError(year int) (map[time.Time]*Holiday, error)`
- `HolidaysForDateRangeWithError(start, end time.Time) (map[time.Time]*Holiday, error)`

**Context API (with cancellation/timeout support):**
- `IsHolidayWithContext(ctx context.Context, date time.Time) (*Holiday, bool, error)`
- `HolidaysForYearWithContext(ctx context.Context, year int) (map[time.Time]*Holiday, error)`
- `HolidaysForDateRangeWithContext(ctx context.Context, start, end time.Time) (map[time.Time]*Holiday, error)`

GoHoliday provides comprehensive error handling with structured errors, context support, and input validation while maintaining full backward compatibility.

Key features:
- **Structured Errors**: Typed errors with specific error codes
- **Context Support**: Cancellation and timeout handling  
- **Input Validation**: Country codes, years, and date ranges
- **Backward Compatibility**: Original API remains unchanged

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

### Adding New Countries

1. Implement the `HolidayProvider` interface in `countries/`
2. Add comprehensive tests
3. Update documentation
4. Submit pull request

See existing implementations for reference patterns.

## Recent Changes

### Version 0.5.3 (2025-08-30)
- Added Portugal, Italy, and India support
- Enhanced GitHub API integration with token authentication
- Improved sync system with better error handling
- Comprehensive test coverage for all new countries

### Version 0.5.0 (2025-08-28)
- Added Norway, Turkey, Russia, and Indonesia
- Orthodox calendar support
- Multi-religious holiday systems
- Enhanced regional subdivision support

## Attribution

This project builds upon the [Python holidays library](https://github.com/vacanza/holidays) and the work of the [Vacanza organization](https://github.com/vacanza), providing a Go-native implementation optimized for performance and ChronoGo integration.

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Contributing

We welcome contributions! Please see our [Contributing Guidelines](CONTRIBUTING.md) for details.