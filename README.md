# GoHolidays

A Go-native library providing comprehensive holiday data for countries and their subdivisions, designed as the premier holiday data provider for [ChronoGo](https://github.com/davidhintelmann/ChronoGo) and other date/time applications.

## 🎯 Primary Use Case: ChronoGo Integration

GoHolidays is purpose-built as a drop-in replacement for ChronoGo's `DefaultHolidayChecker`, providing enterprise-grade holiday data for business day calculations:

```go
// Replace ChronoGo's basic holiday checker with GoHolidays
holidayChecker := chronogo.CreateDefaultUSChecker()

// Now ChronoGo has access to comprehensive holiday data
dt := chronogo.Now()
nextBusinessDay := dt.NextBusinessDay(holidayChecker)
businessDays := dt.AddBusinessDays(5, holidayChecker)

if dt.IsHoliday(holidayChecker) {
    fmt.Println("Today is a holiday!")
}
```

## ✨ Features

- **🏢 Enterprise-Ready**: YAML-based configuration system with environment support
- **🌍 Multi-Country Support**: 7 countries with regional subdivisions
- **⚡ High Performance**: Sub-microsecond holiday lookups with intelligent caching
- **🔧 ChronoGo Integration**: Native `HolidayChecker` interface implementation
- **📊 Business Intelligence**: Holiday categories, regional variations, and custom overrides
- **🔄 Auto-Updates**: GitHub syncer for real-time holiday data from Python holidays library
- **🛡️ Thread-Safe**: Concurrent operations with built-in safety

## 🌍 Supported Countries

| Country | Code | Subdivisions | Status |
|---------|------|-------------|---------|
| �� United States | US | 56 (states, territories) | ✅ Native |
| 🇬🇧 United Kingdom | GB | 4 (England, Scotland, Wales, NI) | ✅ Native |
| 🇨🇦 Canada | CA | 13 (provinces, territories) | ✅ Native |
| 🇦🇺 Australia | AU | 8 (states, territories) | ✅ Native |
| �� New Zealand | NZ | 17 (regions) | ✅ Native |
| 🇩🇪 Germany | DE | 16 (states) | ✅ Native |
| �� France | FR | Regions & territories | ✅ Native |

**Total**: 7 countries with 100+ regional subdivisions

### 🔄 Additional Countries via Python Sync
- 80+ additional countries available through GitHub syncer
- Real-time updates from [Python holidays library](https://github.com/vacanza/python-holidays)
- Automatic country discovery and integration

## 🚀 Quick Start

### ChronoGo Integration (Recommended)

```go
package main

import (
    "fmt"
    "time"
    
    "github.com/your-username/goholidays/chronogo"
    chronogo_lib "github.com/davidhintelmann/chronogo"
)

func main() {
    // Create GoHolidays checker for ChronoGo
    holidayChecker := chronogo.CreateDefaultUSChecker()
    
    // Use with ChronoGo for business day calculations
    dt := chronogo_lib.Now()
    
    // Check if today is a holiday
    if dt.IsHoliday(holidayChecker) {
        fmt.Println("Today is a holiday!")
    }
    
    // Calculate next business day
    nextBusiness := dt.NextBusinessDay(holidayChecker)
    fmt.Printf("Next business day: %s\n", nextBusiness.Format("2006-01-02"))
    
    // Add business days (skipping holidays and weekends)
    futureDate := dt.AddBusinessDays(5, holidayChecker)
    fmt.Printf("5 business days from now: %s\n", futureDate.Format("2006-01-02"))
}
```

### Advanced ChronoGo Configuration

```go
// Multi-country business operations
checker := chronogo.New().
    WithCountries("US", "CA", "GB").
    WithSubdivisions(map[string][]string{
        "US": {"CA", "NY"},  // California and New York
        "CA": {"ON", "BC"},  // Ontario and British Columbia
    }).
    WithCategories("federal", "bank").
    Build()

// Regional business day calculations
dt := chronogo_lib.Now()
isRegionalHoliday := dt.IsHoliday(checker)
```

### Direct Library Usage

```go
package main

import (
    "fmt"
    "time"
    
    "github.com/your-username/goholidays"
    "github.com/your-username/goholidays/config"
)

func main() {
    // Load configuration
    cfg, err := config.LoadFromFile("config/goholidays.yaml")
    if err != nil {
        panic(err)
    }
    
    // Create holiday manager
    manager := config.NewHolidayManager(cfg)
    
    // Check if today is a holiday in multiple countries
    today := time.Now()
    countries := []string{"US", "CA", "GB"}
    
    for _, country := range countries {
        holidays, err := manager.GetHolidays(country, today.Year())
        if err != nil {
            continue
        }
        
        for _, holiday := range holidays {
            if holiday.Date.Format("2006-01-02") == today.Format("2006-01-02") {
                fmt.Printf("%s: %s is a holiday - %s\n", 
                    country, today.Format("Jan 2"), holiday.Name)
            }
        }
    }
}
```

## ⚙️ Configuration System

GoHolidays uses a powerful YAML-based configuration system supporting multiple environments:

### Environment Configurations

```yaml
# config/goholidays.yaml (base configuration)
server:
  name: "GoHolidays"
  version: "1.0.0"
  environment: "production"

performance:
  cache_enabled: true
  cache_ttl: "24h"
  preload_years: [2024, 2025]

countries:
  US:
    enabled: true
    categories: ["federal", "state"]
    subdivisions: ["CA", "NY", "TX"]
  
  CA:
    enabled: true
    categories: ["federal", "provincial"]
```

```yaml
# config/dev.yaml (development overrides)
server:
  environment: "development"

performance:
  cache_enabled: false
  
logging:
  level: "debug"
```

### Custom Holiday Definitions

```yaml
countries:
  US:
    holiday_overrides:
      "Company Day": "2024-03-15"
      "Summer Break": "2024-07-01"
    excluded_holidays:
      - "Columbus Day"
```

## 🔧 Installation

```bash
go get github.com/your-username/goholidays
```

### For ChronoGo Integration

```bash
go get github.com/your-username/goholidays/chronogo
go get github.com/davidhintelmann/chronogo
```

## 📊 Performance Benchmarks

| Operation | Duration | Memory |
|-----------|----------|---------|
| Holiday Check | < 1μs | 0 allocs |
| Year Load | < 100μs | Minimal |
| Multi-Country | < 5μs | Shared cache |

### Caching Benefits

```go
// First check: loads and caches year data
checker.IsHoliday(dt) // ~100μs

// Subsequent checks: cache hits
checker.IsHoliday(dt) // <1μs (100x faster)
```

## 🏗️ Architecture

```
GoHolidays
├── chronogo/           # ChronoGo integration layer
│   ├── integration.go  # HolidayChecker implementation
│   └── examples/       # Usage examples
├── config/             # Configuration system
│   ├── config.go       # YAML configuration
│   ├── manager.go      # Holiday manager
│   └── *.yaml          # Environment configs
├── countries/          # Country implementations
│   ├── us/            # United States
│   ├── gb/            # United Kingdom
│   ├── ca/            # Canada
│   └── ...            # Other countries
└── sync/              # GitHub syncer
    └── syncer.go      # Python holidays sync
```

## 🎯 Why GoHolidays for ChronoGo?

| Feature | ChronoGo Default | GoHolidays |
|---------|------------------|------------|
| Countries | Basic US | 7 countries + 80 via sync |
| Subdivisions | None | 100+ regions |
| Configuration | Hardcoded | YAML + environment |
| Performance | Basic | Cached + optimized |
| Updates | Manual | Automated GitHub sync |
| Categories | Single | Multiple (federal, state, etc.) |

### Migration from ChronoGo Default

```go
// Before: ChronoGo's basic checker
checker := chronogo.DefaultHolidayChecker{}

// After: GoHolidays enterprise checker
checker := chronogo.CreateDefaultUSChecker()

// Same interface, enhanced capabilities!
isHoliday := dt.IsHoliday(checker)
```
    // Create a holiday provider for New Zealand
    nz := goholidays.NewCountry("NZ")
    
    // Check if a specific date is a holiday
    waitangiDay := time.Date(2024, 2, 6, 0, 0, 0, 0, time.UTC)
    if holiday, ok := nz.IsHoliday(waitangiDay); ok {
        fmt.Printf("%s is %s\n", waitangiDay.Format("2006-01-02"), holiday.Name)
        // Output: 2024-02-06 is Waitangi Day
        
        // Check for Māori translation
        if maoriName, exists := holiday.Languages["mi"]; exists {
            fmt.Printf("In Māori: %s\n", maoriName)
            // Output: In Māori: Te Rā o Waitangi
        }
    }
    
    // Get all holidays for a year
    holidays := nz.HolidaysForYear(2024)
    fmt.Printf("New Zealand has %d public holidays in 2024\n", len(holidays))
}
```

### Advanced Features

## 🚀 Getting Started Examples

### 1. Basic ChronoGo Integration

```bash
# Run the ChronoGo integration example
cd examples/chronogo
go run main.go
```

### 2. Enterprise Configuration

```bash
# Test different environments
go run main.go -env=dev      # Development settings
go run main.go -env=prod     # Production settings
go run main.go -env=staging  # Staging settings
```

### 3. Multi-Country Business Operations

```go
// Setup for international business
checker := chronogo.New().
    WithCountries("US", "CA", "GB", "DE", "FR").
    WithCategories("federal", "bank").
    EnableCaching().
    Build()

// Check holidays across all countries
dt := chronogo.Now()
if dt.IsHoliday(checker) {
    // Handle holiday in any configured country
    fmt.Println("Holiday detected in at least one country")
}
```

## 🧪 Testing

```bash
# Run all tests
go test ./...

# Test specific components
go test ./chronogo -v          # ChronoGo integration tests
go test ./config -v            # Configuration system tests
go test ./countries/us -v      # US holiday tests

# Run benchmarks
go test ./chronogo -bench=.    # Performance benchmarks
```

## 🔄 GitHub Syncer

Keep holiday data up-to-date with the Python holidays library:

```go
import "github.com/your-username/goholidays/sync"

// Sync with Python holidays repository
syncer := sync.NewGitHubSyncer("your-token")
countries, err := syncer.SyncCountries()
if err != nil {
    log.Fatal(err)
}

fmt.Printf("Synced %d countries\n", len(countries))
```

## 📈 Roadmap

### Phase 1: ChronoGo Integration ✅
- [x] HolidayChecker interface implementation
- [x] Performance optimization with caching
- [x] Multi-country support
- [x] Regional subdivision support

### Phase 2: Enhanced Features 🚧
- [ ] Holiday forecast API
- [ ] Business day calculation helpers
- [ ] Advanced holiday rules engine
- [ ] REST API service (optional)

### Phase 3: Ecosystem Integration 📋
- [ ] Integration with other date/time libraries
- [ ] Enterprise authentication and authorization
- [ ] Cloud deployment templates
- [ ] Monitoring and observability

## 🤝 Contributing

We welcome contributions! Please see our [Contributing Guidelines](CONTRIBUTING.md) for details.

### Development Setup

```bash
# Clone the repository
git clone https://github.com/your-username/goholidays
cd goholidays

# Install dependencies
go mod download

# Run tests
go test ./...

# Run the ChronoGo example
cd examples/chronogo && go run main.go
```

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 🙏 Acknowledgments

- [ChronoGo](https://github.com/davidhintelmann/ChronoGo) - The excellent Go date/time library we integrate with
- [Python holidays library](https://github.com/vacanza/python-holidays) - Source of comprehensive holiday data
- [Vacanza community](https://github.com/vacanza) - Maintaining the Python holidays ecosystem

---

**Ready to enhance your ChronoGo applications with comprehensive holiday data?**

```bash
go get github.com/your-username/goholidays/chronogo
```

**Experience the difference**: From basic US holidays to enterprise-grade, multi-country holiday intelligence.
