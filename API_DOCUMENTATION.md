# GoHoliday API Documentation
*Comprehensive Guide to the GoHoliday Library v1.0*

## Table of Contents
1. [Quick Start](#quick-start)
2. [Core API Reference](#core-api-reference)
3. [Country Coverage](#country-coverage)
4. [Advanced Features](#advanced-features)
5. [Performance Guide](#performance-guide)
6. [Migration Guide](#migration-guide)
7. [Examples](#examples)

---

## Quick Start

### Installation
```bash
go get github.com/coredds/GoHoliday
```

### Basic Usage
```go
package main

import (
    "fmt"
    "time"
    "github.com/coredds/GoHoliday"
)

func main() {
    // Create a country holiday provider
    us := goholidays.NewCountry("US")
    
    // Check if a date is a holiday
    date := time.Date(2025, 7, 4, 0, 0, 0, 0, time.UTC)
    if holiday, isHoliday := us.IsHoliday(date); isHoliday {
        fmt.Printf("%s is %s\n", date.Format("January 2"), holiday.Name)
        // Output: July 4 is Independence Day
    }
    
    // Get all holidays for a year
    holidays := us.HolidaysForYear(2025)
    fmt.Printf("US has %d holidays in 2025\n", len(holidays))
}
```

---

## Core API Reference

### Country Creation

#### `NewCountry(countryCode string, options ...CountryOptions) *Country`
Creates a new country holiday provider.

**Parameters:**
- `countryCode`: ISO 3166-1 alpha-2 country code (e.g., "US", "FR", "IN")
- `options`: Optional configuration

**Supported Countries:**
- `US` - United States
- `CA` - Canada  
- `GB` - United Kingdom
- `AU` - Australia
- `NZ` - New Zealand
- `JP` - Japan
- `IN` - India
- `FR` - France
- `DE` - Germany

**Example:**
```go
// Basic usage
us := goholidays.NewCountry("US")

// With options
options := goholidays.CountryOptions{
    Subdivisions: []string{"CA", "NY"},
    Categories:   []goholidays.HolidayCategory{goholidays.CategoryPublic},
    Language:     "en",
}
us := goholidays.NewCountry("US", options)
```

### Holiday Lookup

#### `IsHoliday(date time.Time) (*Holiday, bool)`
Checks if a given date is a holiday. **Thread-safe**.

**Returns:**
- `*Holiday`: Holiday information if found
- `bool`: Whether the date is a holiday

**Example:**
```go
date := time.Date(2025, 12, 25, 0, 0, 0, 0, time.UTC)
if holiday, isHoliday := country.IsHoliday(date); isHoliday {
    fmt.Printf("Holiday: %s (%s)\n", holiday.Name, holiday.Category)
}
```

#### `HolidaysForYear(year int) map[time.Time]*Holiday`
Returns all holidays for a specific year. **Thread-safe**.

**Example:**
```go
holidays := country.HolidaysForYear(2025)
for date, holiday := range holidays {
    fmt.Printf("%s: %s\n", date.Format("Jan 2"), holiday.Name)
}
```

#### `HolidaysForDateRange(start, end time.Time) map[time.Time]*Holiday`
Returns holidays within a date range.

**Example:**
```go
start := time.Date(2025, 6, 1, 0, 0, 0, 0, time.UTC)
end := time.Date(2025, 8, 31, 0, 0, 0, 0, time.UTC)
summerHolidays := country.HolidaysForDateRange(start, end)
```

### Holiday Structure

```go
type Holiday struct {
    Name       string                `json:"name"`
    Date       time.Time             `json:"date"`
    Category   HolidayCategory       `json:"category"`
    Observed   *time.Time            `json:"observed,omitempty"`
    Languages  map[string]string     `json:"languages,omitempty"`
    IsObserved bool                  `json:"is_observed"`
}
```

**Fields:**
- `Name`: Primary name of the holiday
- `Date`: Date of the holiday
- `Category`: Type of holiday (public, religious, etc.)
- `Observed`: Alternative observed date if different
- `Languages`: Translations in different languages
- `IsObserved`: Whether this is an observed date

### Holiday Categories

```go
type HolidayCategory string

const (
    CategoryPublic      HolidayCategory = "public"
    CategoryBank        HolidayCategory = "bank"
    CategorySchool      HolidayCategory = "school"
    CategoryGovernment  HolidayCategory = "government"
    CategoryReligious   HolidayCategory = "religious"
    CategoryOptional    HolidayCategory = "optional"
    CategoryHalfDay     HolidayCategory = "half_day"
    CategoryArmedForces HolidayCategory = "armed_forces"
    CategoryWorkday     HolidayCategory = "workday"
)
```

---

## Country Coverage

### United States (US)
**Federal Holidays:** New Year's Day, MLK Day, Presidents' Day, Memorial Day, Juneteenth, Independence Day, Labor Day, Veterans Day, Thanksgiving, Christmas

**State Support:** All 50 states + DC
**Languages:** English, Spanish

### Canada (CA)  
**National Holidays:** New Year's Day, Good Friday, Victoria Day, Canada Day, Labour Day, Thanksgiving, Christmas

**Province Support:** All provinces and territories
**Languages:** English, French

### United Kingdom (GB)
**National Holidays:** New Year's Day, Good Friday, Easter Monday, Christmas Day, Boxing Day

**Region Support:** England, Scotland, Wales, Northern Ireland
**Languages:** English

### Australia (AU)
**National Holidays:** New Year's Day, Australia Day, Good Friday, Easter Monday, ANZAC Day, Christmas Day, Boxing Day

**State Support:** All states and territories
**Languages:** English

### New Zealand (NZ)
**National Holidays:** New Year's Day, Waitangi Day, Good Friday, Easter Monday, ANZAC Day, Queen's Birthday, Labour Day, Christmas Day, Boxing Day

**Regional Support:** All regions
**Languages:** English, Māori

### Japan (JP)
**National Holidays:** New Year's Day, Coming of Age Day, National Foundation Day, Emperor's Birthday, Showa Day, Constitution Memorial Day, Greenery Day, Children's Day, Marine Day, Mountain Day, Respect for the Aged Day, Autumnal Equinox Day, Health and Sports Day, Culture Day, Labour Thanksgiving Day

**Languages:** English, Japanese

### India (IN) ✨ *New*
**National Holidays:** Republic Day, Independence Day, Gandhi Jayanti

**Religious Festivals:** Diwali, Holi, Eid al-Fitr, Christmas, Good Friday

**State Support:** All 36 states and union territories
**Languages:** English, Hindi, Arabic

### France (FR) ✨ *New*
**National Holidays:** New Year's Day, Labour Day, Victory in Europe Day, Bastille Day, Assumption of Mary, All Saints' Day, Armistice Day, Christmas Day

**Religious Holidays:** Easter Monday, Ascension Day, Whit Monday

**Regional Support:** All regions including overseas territories
**Languages:** French, English

### Germany (DE)
**National Holidays:** New Year's Day, Good Friday, Easter Monday, Labour Day, Ascension Day, Whit Monday, German Unity Day, Christmas Day, Boxing Day

**State Support:** All 16 federal states  
**Languages:** German, English

---

## Advanced Features

### Multi-Language Support

All holidays include translations where culturally appropriate:

```go
holiday, _ := country.IsHoliday(date)
if holiday.Languages != nil {
    fmt.Printf("English: %s\n", holiday.Languages["en"])
    fmt.Printf("Local: %s\n", holiday.Languages["fr"]) // or "hi", "ja", etc.
}
```

### Category Filtering

Filter holidays by category:

```go
options := goholidays.CountryOptions{
    Categories: []goholidays.HolidayCategory{
        goholidays.CategoryPublic,
        goholidays.CategoryReligious,
    },
}
country := goholidays.NewCountry("FR", options)
```

### Regional/State Holidays

Access regional holidays:

```go
options := goholidays.CountryOptions{
    Subdivisions: []string{"CA", "NY"}, // California and New York
}
us := goholidays.NewCountry("US", options)
```

### Performance Optimization

#### Object Pooling
```go
// Create memory-optimized holidays
holiday := goholidays.OptimizedHoliday(
    "Holiday Name",
    date,
    goholidays.CategoryPublic,
    map[string]string{"en": "English Name"},
)

// Return to pool when done
goholidays.GlobalHolidayPool.Put(holiday)
```

#### Caching
```go
// Create LRU cache for computed holidays
cache := goholidays.NewHolidayCache(100) // Max 100 entries
cache.Set("key", holidays)
if cached, exists := cache.Get("key"); exists {
    // Use cached holidays
}
```

### API Stability

#### Version Management
```go
// Check API version
info := goholidays.GetVersionInfo()
fmt.Printf("Library: %s, API: %s\n", info.LibraryVersion, info.APIVersion)

// Validate feature stability
err := goholidays.ValidateAPIUsage("NewCountry", goholidays.StabilityStable)
if err != nil {
    log.Fatal("API compatibility issue:", err)
}
```

#### Deprecation Handling
```go
// Check for deprecated features
goholidays.CheckDeprecation("feature_name") // Logs warnings if deprecated
```

---

## Performance Guide

### Benchmarks
Current performance metrics on AMD Ryzen 7 5700G:

| Operation | Performance | Memory |
|-----------|-------------|---------|
| Holiday Lookup | ~104 ns/op | 0 B/op |
| Concurrent Access | 40.93M ops/sec | Thread-safe |
| Holiday Loading | ~1,442 ns/op | 3,920 B/op |
| Memory per Country | 9.74 KB | Optimized |

### Best Practices

1. **Reuse Country Instances**: Create once, use many times
```go
// Good
us := goholidays.NewCountry("US")
for _, date := range dates {
    us.IsHoliday(date)
}

// Avoid
for _, date := range dates {
    us := goholidays.NewCountry("US") // Inefficient
    us.IsHoliday(date)
}
```

2. **Pre-load Years**: Load commonly used years upfront
```go
options := goholidays.CountryOptions{
    Years: []int{2024, 2025, 2026}, // Pre-load these years
}
us := goholidays.NewCountry("US", options)
```

3. **Use Concurrent Access**: Library is thread-safe
```go
var wg sync.WaitGroup
for i := 0; i < 100; i++ {
    wg.Add(1)
    go func() {
        defer wg.Done()
        us.IsHoliday(date) // Safe concurrent access
    }()
}
wg.Wait()
```

---

## Migration Guide

### From v0.3.0 to v1.0.0

#### Breaking Changes
- None! v1.0.0 is fully backward compatible with v0.3.0

#### New Features in v1.0.0
1. **Enhanced Global Coverage**: Additional country implementations beyond 15 countries
2. **Advanced Performance**: Further optimization beyond 275K+ ops/sec
3. **API Stability**: Mature version management and deprecation system
4. **Enterprise Features**: Advanced configuration and monitoring

### From v0.2.2 to v0.3.0

#### Breaking Changes
- None! v0.3.0 is fully backward compatible with v0.2.2

#### New Features in v0.3.0
1. **Four New Countries**: Italy (IT), Spain (ES), Netherlands (NL), South Korea (KR)
2. **Enhanced Coverage**: 15 countries with 200+ regional subdivisions
3. **Cultural Accuracy**: Native language support for all new countries
4. **Maintained Performance**: 275K+ operations/second with expanded coverage
5. **Comprehensive Testing**: 82.7% test coverage in countries package

### From v0.2.2 to v1.0.0

#### Breaking Changes
- None! v1.0.0 is fully backward compatible with v0.2.2

#### New Features in v1.0.0
1. **Enhanced Global Coverage**: Additional country implementations
2. **Advanced Performance**: Further optimization beyond 40M+ ops/sec
3. **API Stability**: Mature version management and deprecation system
4. **Enterprise Features**: Advanced configuration and monitoring

### From v0.1.2 to v0.2.2

#### Breaking Changes
- None! v0.2.2 is fully backward compatible with v0.1.2

#### New Features in v0.2.2
1. **Brazil Support**: Complete Brazilian holiday implementation with Carnival
2. **Mexico Support**: Mexican holidays with constitutional reforms
3. **Latin American Coverage**: Comprehensive South/North American representation
4. **Cultural Accuracy**: Real-world holiday calculations and traditions

#### Migration Steps
1. Update import if needed:
```go
// Old (still works)
import "github.com/coredds/GoHoliday"

// New (recommended)
import "github.com/coredds/GoHoliday"
```

2. Add France support:
```go
// New capability
france := goholidays.NewCountry("FR")
```

3. Optional: Use new performance features:
```go
// Optional optimization
holiday := goholidays.OptimizedHoliday(name, date, category, languages)
```

---

## Examples

### Multi-Country Application
```go
package main

import (
    "fmt"
    "time"
    "github.com/coredds/GoHoliday"
)

func main() {
    countries := map[string]*goholidays.Country{
        "🇺🇸 US": goholidays.NewCountry("US"),
        "🇫🇷 FR": goholidays.NewCountry("FR"),
        "🇮🇳 IN": goholidays.NewCountry("IN"),
        "🇯🇵 JP": goholidays.NewCountry("JP"),
    }
    
    date := time.Date(2025, 12, 25, 0, 0, 0, 0, time.UTC)
    
    for name, country := range countries {
        if holiday, isHoliday := country.IsHoliday(date); isHoliday {
            fmt.Printf("%s: %s\n", name, holiday.Name)
        }
    }
}
```

### Business Day Calculator
```go
package main

import (
    "time"
    "github.com/coredds/GoHoliday"
)

func isBusinessDay(date time.Time, country *goholidays.Country) bool {
    // Weekend check
    if date.Weekday() == time.Saturday || date.Weekday() == time.Sunday {
        return false
    }
    
    // Holiday check
    _, isHoliday := country.IsHoliday(date)
    return !isHoliday
}

func main() {
    us := goholidays.NewCountry("US")
    date := time.Date(2025, 7, 4, 0, 0, 0, 0, time.UTC)
    
    if isBusinessDay(date, us) {
        fmt.Println("Business day")
    } else {
        fmt.Println("Non-business day")
    }
}
```

### Holiday Calendar
```go
package main

import (
    "fmt"
    "time"
    "github.com/coredds/GoHoliday"
)

func printHolidayCalendar(country *goholidays.Country, year int) {
    holidays := country.HolidaysForYear(year)
    
    for date, holiday := range holidays {
        fmt.Printf("%-15s | %-25s | %s\n", 
            date.Format("Jan 2, 2006"),
            holiday.Name,
            holiday.Category)
    }
}

func main() {
    france := goholidays.NewCountry("FR")
    fmt.Println("🇫🇷 France Holidays 2025")
    fmt.Println("=" * 50)
    printHolidayCalendar(france, 2025)
}
```

---

## Support & Contributing

### Documentation
- **API Reference**: [GoDoc](https://pkg.go.dev/github.com/coredds/GoHoliday)
- **Examples**: `/examples` directory
- **Performance**: `/examples/performance-analysis`

### Community
- **Issues**: [GitHub Issues](https://github.com/coredds/GoHoliday/issues)
- **Discussions**: [GitHub Discussions](https://github.com/coredds/GoHoliday/discussions)
- **Contributing**: See `CONTRIBUTING.md`

### Performance Analysis
Run the built-in performance analysis:
```bash
go run examples/performance-analysis/main.go
```

---

*GoHoliday v1.0.0 - Production-ready holiday data for Go applications*
*Performance: 40M+ ops/sec | Coverage: 9 countries | Thread-safe: ✅*
