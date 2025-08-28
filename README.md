# GoHoliday

[![CI](https://github.com/coredds/GoHoliday/workflows/CI/badge.svg)](https://github.com/coredds/GoHoliday/actions/workflows/ci.yml)
[![codecov](https://codecov.io/gh/coredds/GoHoliday/branch/main/graph/badge.svg)](https://codecov.io/gh/coredds/GoHoliday)
[![Security](https://github.com/coredds/GoHoliday/workflows/CI/badge.svg?job=security)](https://github.com/coredds/GoHoliday/actions/workflows/ci.yml)
[![Go Version](https://img.shields.io/github/go-mod/go-version/coredds/GoHoliday?v=1.23)](https://golang.org/)
[![License](https://img.shields.io/github/license/coredds/GoHoliday)](LICENSE)

A comprehensive Go library for holiday data and business day calculations, designed as a high-performance backend for date/time applications including [ChronoGo](https://github.com/davidhintelmann/ChronoGo).

## Features

- **Multi-Country Support**: 25 countries with 500+ regional subdivisions
- **High Performance**: Sub-microsecond holiday lookups with intelligent caching  
- **ChronoGo Integration**: Native HolidayChecker interface implementation
- **Enterprise Configuration**: YAML-based configuration system with environment support
- **Cultural Accuracy**: Multi-language holiday names and historical transitions
- **Thread-Safe**: Concurrent operations with built-in safety
- **Business Intelligence**: Holiday categories, regional variations, and custom overrides

Current version: **0.5.0** with Norway, Turkey, Russia, and Indonesia support

## Supported Countries

| Country | Code | Subdivisions | Languages | Implementation |
|---------|------|-------------|-----------|----------------|
| United States | US | 56 (states, territories) | EN, ES | Native |
| United Kingdom | GB | 4 (England, Scotland, Wales, NI) | EN | Native |
| Canada | CA | 13 (provinces, territories) | EN, FR | Native |
| Australia | AU | 8 (states, territories) | EN | Native |
| New Zealand | NZ | 17 (regions) | EN, MI | Native |
| Germany | DE | 16 (states) | DE, EN | Native |
| France | FR | Regions & territories | FR, EN | Native |
| Japan | JP | National public holidays | JA, EN | Native |
| India | IN | National & state holidays | HI, EN | Native |
| Brazil | BR | 27 (states & federal district) | PT, EN | Native |
| Mexico | MX | 32 (states & federal district) | ES, EN | Native |
| Italy | IT | 20 (regions) | IT, EN | Native |
| Spain | ES | 19 (autonomous communities) | ES, EN | Native |
| Netherlands | NL | 12 (provinces) | NL, EN | Native |
| South Korea | KR | 17 (provinces & cities) | KO, EN | Native |
| Singapore | SG | 5 (regions) | EN, ZH, MS, TA | Native |
| Switzerland | CH | 26 (cantons) | DE, FR, IT, RM | Native |
| Sweden | SE | 21 (counties) | SV, EN | Native |
| Argentina | AR | 24 (provinces) | ES, EN | Native |
| Thailand | TH | 77 (provinces) | TH, EN | Native |
| **Norway** | **NO** | **11 (counties)** | **NO, EN** | **Native** |
| **Turkey** | **TR** | **81 (provinces)** | **TR, EN** | **Native** |
| **Russia** | **RU** | **85 (federal subjects)** | **RU, EN** | **Native** |
| **Indonesia** | **ID** | **38 (provinces)** | **ID, EN** | **Native** |

**Total**: 25 countries with 500+ regional subdivisions

### Recent Additions (v0.5.0)
- **Norway**: Nordic country with Constitution Day (Grunnlovsdag) and traditional Norse holidays
- **Turkey**: Secular and Islamic holidays including Democracy Day and national observances  
- **Russia**: Orthodox calendar with extensive New Year celebration period and federal holidays
- **Indonesia**: Multi-religious society with Islamic, Christian, Buddhist, Hindu, and Chinese holidays

### Previous Additions (v0.4.0)
- **Singapore**: Multicultural holidays with 4 official languages (English, Chinese, Malay, Tamil)
- **Switzerland**: Federal and cantonal holidays with 4 national languages (German, French, Italian, Romansh)
- **Sweden**: Traditional Nordic holidays including Midsummer calculations
- **Argentina**: Latin American representation with Spanish localization and movable holidays
- **Thailand**: Southeast Asian Buddhist calendar with royal observances and lunar calculations

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

### Multi-Country and Multilingual Support

```go
package main

import (
    "fmt"
    "time"
    
    "github.com/coredds/GoHoliday"
)

func main() {
    // Thailand with Buddhist holidays and Thai language support
    th := goholidays.NewCountry("TH")
    
    // Check Songkran Festival (Thai New Year)
    songkran2024 := time.Date(2024, 4, 13, 0, 0, 0, 0, time.UTC)
    if holiday, ok := th.IsHoliday(songkran2024); ok {
        fmt.Printf("Thailand: %s\n", holiday.Name)
        if thName, exists := holiday.Languages["th"]; exists {
            fmt.Printf("ไทย: %s\n", thName)
        }
    }
    
    // Argentina with Spanish support and movable holidays
    ar := goholidays.NewCountry("AR")
    
    // Check Flag Day (movable to Monday if weekend)
    flagDay2024 := time.Date(2024, 6, 20, 0, 0, 0, 0, time.UTC)
    if holiday, ok := ar.IsHoliday(flagDay2024); ok {
        fmt.Printf("Argentina: %s\n", holiday.Name)
        if esName, exists := holiday.Languages["es"]; exists {
            fmt.Printf("Español: %s\n", esName)
        }
    }
    
    // Singapore with 4 official languages
    sg := goholidays.NewCountry("SG")
    
    // Check Chinese New Year
    cny2024 := time.Date(2024, 2, 10, 0, 0, 0, 0, time.UTC)
    if holiday, ok := sg.IsHoliday(cny2024); ok {
        fmt.Printf("Singapore: %s\n", holiday.Name)
        
        // Display in all 4 official languages
        languages := []string{"en", "zh", "ms", "ta"}
        for _, lang := range languages {
            if name, exists := holiday.Languages[lang]; exists {
                fmt.Printf("  %s: %s\n", lang, name)
            }
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

### New Country Performance (v0.4.0)
| Provider | Load Time | Memory | Features |
|----------|-----------|--------|----------|
| Singapore | ~2.9μs | Low | 4-language multicultural |
| Switzerland | ~3.2μs | Low | Federal/cantonal system |
| Sweden | ~4.3μs | Low | Traditional calculations |
| Argentina | ~3.6μs | Low | Movable holidays |
| Thailand | ~6.8μs | Medium | Buddhist lunar calendar |

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

# New country demos (v0.4.0)
cd examples/singapore-demo && go run main.go
cd examples/thailand-demo && go run main.go
cd examples/argentina-demo && go run main.go

# Multi-country comparisons
cd examples/multi-country && go run main.go
cd examples/new-countries-demo && go run main.go

# Performance analysis
cd examples/performance-analysis && go run main.go
```

## Attribution

This project builds upon the [Python holidays library](https://github.com/vacanza/holidays) and the work of the [Vacanza organization](https://github.com/vacanza), providing a Go-native implementation optimized for performance and ChronoGo integration.

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Contributing

We welcome contributions! Please see our [Contributing Guidelines](CONTRIBUTING.md) for details.
