# goholiday

[![CI](https://github.com/coredds/goholiday/workflows/CI/badge.svg)](https://github.com/coredds/goholiday/actions/workflows/ci.yml)
[![codecov](https://codecov.io/gh/coredds/goholiday/branch/main/graph/badge.svg)](https://codecov.io/gh/coredds/goholiday)
[![Go Version](https://img.shields.io/github/go-mod/go-version/coredds/goholiday?v=1.23)](https://golang.org/)
[![License](https://img.shields.io/github/license/coredds/goholiday)](LICENSE)

A comprehensive Go library for holiday data and business day calculations. Provides high-performance holiday checking with multi-country support, designed for integration with date/time applications including [chronogo](https://github.com/davidhintelmann/chronogo).

**Current Version**: 0.6.4

## Features

- **37 Countries**: Complete coverage with 600+ regional subdivisions
- **High Performance**: Sub-microsecond holiday lookups with intelligent caching
- **Multi-Language**: Native language support for holiday names
- **Thread-Safe**: Concurrent operations with built-in safety
- **chronogo Integration**: Native HolidayChecker interface implementation
- **Enterprise Configuration**: YAML-based configuration with environment support
- **Regional Variations**: State, province, and regional holiday support
- **Historical Accuracy**: Proper handling of holiday transitions and changes
- **Robust Error Handling**: Structured errors with context support and validation
- **Backward Compatible**: Enhanced API alongside original API

## Supported Countries

**37 countries** with comprehensive holiday coverage:

| Country | Code | Subdivisions | Languages | Key Features |
|---------|------|-------------|-----------|--------------|
| Argentina | AR | 24 provinces | ES, EN | National and provincial holidays |
| Australia | AU | 8 states/territories | EN | National and state holidays |
| Austria | AT | 9 states | DE, EN | Federal and state holidays |
| Belgium | BE | 3 regions | NL, FR, DE, EN | National and regional holidays |
| Brazil | BR | 27 states | PT, EN | National and state holidays |
| Canada | CA | 13 provinces/territories | EN, FR | Federal and provincial holidays |
| Chile | CL | 16 regions | ES, EN | Variable holidays, regional laws |
| China | CN | 34 provinces/regions | ZH, EN | Lunar calendar holidays |
| Finland | FI | 19 regions | FI, SV, EN | National holidays |
| France | FR | Regions & territories | FR, EN | National and regional holidays |
| Germany | DE | 16 states | DE, EN | Federal and state holidays |
| India | IN | 36 states/territories | HI, EN | Multi-religious festivals |
| Indonesia | ID | 38 provinces | ID, EN | Multi-religious holidays |
| Ireland | IE | 30 counties/provinces | EN, GA | Celtic festivals, bank holidays |
| Israel | IL | 6 districts | EN, HE | Hebrew calendar, memorial days |
| Italy | IT | 20 regions | IT, EN | National and patron saint holidays |
| Japan | JP | National | JA, EN | Public holidays |
| Mexico | MX | 32 states | ES, EN | National and state holidays |
| Netherlands | NL | 12 provinces | NL, EN | National holidays |
| New Zealand | NZ | 17 regions | EN, MI | National and regional holidays |
| Norway | NO | 11 counties | NO, EN | National holidays |
| Poland | PL | 16 voivodeships | PL, EN | National and regional holidays |
| Portugal | PT | 20 districts/regions | PT, EN | National and regional holidays |
| Russia | RU | 85 federal subjects | RU, EN | Orthodox calendar holidays |
| Singapore | SG | 5 regions | EN, ZH, MS, TA | Multi-cultural holidays |
| South Korea | KR | 17 provinces/cities | KO, EN | National holidays |
| Spain | ES | 19 autonomous communities | ES, EN | National and regional holidays |
| Sweden | SE | 21 counties | SV, EN | National holidays |
| Switzerland | CH | 26 cantons | DE, FR, IT, RM | Federal and cantonal holidays |
| Thailand | TH | 77 provinces | TH, EN | Buddhist and royal holidays |
| Turkey | TR | 81 provinces | TR, EN | Secular and religious holidays |
| Ukraine | UA | 27 regions | UK, EN, RU | Orthodox calendar, historical holidays |
| United Kingdom | GB | 4 countries | EN | Bank holidays, regional variations |
| United States | US | 56 states/territories | EN, ES | Federal and state holidays |

## Installation

```bash
go get github.com/coredds/goholiday
```

For chronogo integration:
```bash
go get github.com/coredds/goholiday/chronogo
```

## Quick Start

### Basic Usage

```go
package main

import (
    "fmt"
    "time"
    
    "github.com/coredds/goholiday"
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
    holidays := us.HolidaysForYear(2024)
    fmt.Printf("Found %d holidays in 2024\n", len(holidays))
}
```

### chronogo Integration

```go
package main

import (
    "fmt"
    "time"
    
    "github.com/coredds/goholiday/chronogo"
    chronogo_lib "github.com/davidhintelmann/chronogo"
)

func main() {
    // Create holiday checker for chronogo
    holidayChecker := chronogo.Checker("US")
    
    // Use with chronogo for business day calculations
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
// Chile with Spanish support
chile := goholidays.NewCountry("CL")
date := time.Date(2024, 9, 18, 0, 0, 0, 0, time.UTC)
if holiday, ok := chile.IsHoliday(date); ok {
    fmt.Printf("English: %s\n", holiday.Name)           // Independence Day
    fmt.Printf("Spanish: %s\n", holiday.Languages["es"]) // Día de la Independencia
}

// Ireland with Irish Gaelic support
ireland := goholidays.NewCountry("IE")
date = time.Date(2024, 3, 17, 0, 0, 0, 0, time.UTC)
if holiday, ok := ireland.IsHoliday(date); ok {
    fmt.Printf("English: %s\n", holiday.Name)           // Saint Patrick's Day
    fmt.Printf("Irish: %s\n", holiday.Languages["ga"])  // Lá Fhéile Pádraig
}

// Israel with Hebrew support
israel := goholidays.NewCountry("IL")
date = time.Date(2024, 4, 23, 0, 0, 0, 0, time.UTC)
if holiday, ok := israel.IsHoliday(date); ok {
    fmt.Printf("English: %s\n", holiday.Name)           // Passover
    fmt.Printf("Hebrew: %s\n", holiday.Languages["he"]) // פסח
}
```

### Error Handling

```go
// Create country with validation
country, err := goholidays.NewCountryWithError("US")
if err != nil {
    var holidayErr *goholidays.HolidayError
    if errors.As(err, &holidayErr) {
        fmt.Printf("Error: %s (code: %d)\n", holidayErr.Message, holidayErr.Code)
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

// Use context for cancellation/timeout
ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
defer cancel()

holidays, err := country.HolidaysForYearWithContext(ctx, 2024)
if err != nil {
    if goholidays.IsContextCancelled(err) {
        fmt.Println("Operation timed out")
    }
    return
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

YAML-based configuration for enterprise deployments:

```yaml
countries:
  US:
    subdivisions: ["CA", "NY", "TX"]
    categories: ["federal", "state"]
    custom_holidays:
      - name: "Company Day"
        date: "2024-03-15"
        category: "company"
```

## Architecture

```
goholiday/
├── chronogo/           # chronogo integration layer
├── config/             # Configuration system
├── countries/          # Country implementations
├── updater/            # Python holidays syncer
└── examples/           # Demo applications
```

## API Reference

### Core Methods

**Original API (backward compatible):**
- `NewCountry(countryCode)` - Create country instance
- `IsHoliday(date)` - Check if date is holiday
- `HolidaysForYear(year)` - Get all holidays for year
- `HolidaysForDateRange(start, end)` - Get holidays in range

**Enhanced API (with error handling):**
- `NewCountryWithError(countryCode)` - Create with validation
- `IsHolidayWithError(date)` - Check with error handling
- `HolidaysForYearWithError(year)` - Get holidays with validation
- `HolidaysForDateRangeWithError(start, end)` - Get range with validation

**Context API (with cancellation/timeout):**
- `IsHolidayWithContext(ctx, date)` - Check with context
- `HolidaysForYearWithContext(ctx, year)` - Get holidays with context
- `HolidaysForDateRangeWithContext(ctx, start, end)` - Get range with context

Features: Structured errors, context support, input validation, full backward compatibility.

## Testing

```bash
go test ./...                  # Run all tests
go test ./countries -v         # Country provider tests
go test ./chronogo -bench=.    # Performance benchmarks
```

## Development

```bash
git clone https://github.com/coredds/goholiday
cd goholiday
go mod download
go test ./...
```

### Adding New Countries

1. Implement `HolidayProvider` interface in `countries/`
2. Add comprehensive tests
3. Update documentation and integration points
4. Submit pull request

## Recent Changes

### Version 0.6.3 (2025-09-18)
- Added Chile, Ireland, and Israel support
- Enhanced multi-language support (Spanish, Irish Gaelic, Hebrew)
- Variable holiday laws and regional variations
- Hebrew calendar and Celtic festival support
- Comprehensive error handling and context support
- Improved codebase organization and test coverage

### Version 0.5.3 (2025-08-30)
- Added Portugal, Italy, India, and Ukraine support
- Enhanced GitHub API integration with token authentication
- Improved sync system with better error handling
- Orthodox calendar support and multi-religious holiday systems

## Attribution

This project builds upon the [Python holidays library](https://github.com/vacanza/holidays) and the work of the [Vacanza organization](https://github.com/vacanza), providing a Go-native implementation optimized for performance and chronogo integration.

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Contributing

We welcome contributions! Please see our [Contributing Guidelines](CONTRIBUTING.md) for details.