## goholiday v0.6.4 - Holiday Rule Updates, Accuracy Improvements & Go Naming Convention Compliance

### Module Naming Convention
- **Go Best Practices**: Module renamed from `github.com/coredds/GoHoliday` to `github.com/coredds/goholiday`
- **All Lowercase**: Follows idiomatic Go naming conventions for packages and modules
- **Full Compatibility**: All imports, documentation, and examples updated
- **Breaking Change**: Users must update import paths to `github.com/coredds/goholiday`

### Enhanced Holiday Accuracy
- **Rule Alignment**: Updated holiday rules to match the authoritative Python holidays reference library
- **Historical Transitions**: Proper handling of holiday name changes and date adjustments over time
- **Olympic Adjustments**: Accurate implementation of Tokyo Olympics 2020-2021 holiday date changes

### Updated Countries (4)
- **Canada (CA)**: Added National Day for Truth and Reconciliation (September 30, since 2021)
  - Bilingual English/French support
  - Historical accuracy with year-based implementation
- **Australia (AU)**: King's Birthday vs Queen's Birthday naming transition
  - Queen's Birthday for 2022 and earlier
  - King's Birthday for 2023 onwards
  - Applied to both national and Queensland variations
- **Japan (JP)**: Tokyo Olympics holiday date adjustments (2020-2021)
  - Marine Day, Mountain Day, and Sports Day moved for Olympic years
  - Normal scheduling maintained for all other years
- **Netherlands (NL)**: Historical Queen's Day to King's Day transition
  - Queen's Day (April 30) before 2014
  - King's Day (April 27) from 2014 onwards
  - Proper weekend shifting rules for both holidays

### Performance & Quality
- **Sub-microsecond Performance**: Maintained ~49ns holiday lookups
- **100% Test Coverage**: All updated countries fully tested
- **Backward Compatibility**: All existing functionality preserved
- **chronogo Integration**: Updated with all 34 supported countries

### Installation
```bash
go get github.com/coredds/goholiday@v0.6.4
```

### Quick Example - New Features
```go
// Canada's Truth and Reconciliation Day (since 2021)
canada := goholidays.NewCountry("CA")
if holiday, isHoliday := canada.IsHoliday(time.Date(2021, 9, 30, 0, 0, 0, 0, time.UTC)); isHoliday {
    fmt.Printf("Holiday: %s\n", holiday.Name) // "National Day for Truth and Reconciliation"
    fmt.Printf("French: %s\n", holiday.Languages["fr"]) // "Journée nationale de la vérité et de la réconciliation"
}

// Australia's King's Birthday (2023+)
australia := goholidays.NewCountry("AU")
if holiday, isHoliday := australia.IsHoliday(time.Date(2023, 6, 12, 0, 0, 0, 0, time.UTC)); isHoliday {
    fmt.Printf("Holiday: %s\n", holiday.Name) // "King's Birthday" (was "Queen's Birthday" before 2023)
}

// Japan's Olympic year adjustments
japan := goholidays.NewCountry("JP")
if holiday, isHoliday := japan.IsHoliday(time.Date(2020, 7, 23, 0, 0, 0, 0, time.UTC)); isHoliday {
    fmt.Printf("Holiday: %s\n", holiday.Name) // "Marine Day" (moved for Olympics)
}
```

---

## goholiday v0.6.3 - Major Country Expansion

### New Countries Added (3)
- **Chile (CL)**: 16 holidays with unique legal framework
  - Variable holidays that move to Monday per Chilean "Ley de Feriados"
  - Regional holidays: Battle of Arica (2020+), Chillan Foundation Day (2019+)
  - Spanish/English bilingual support
- **Ireland (IE)**: 14 holidays with Celtic and Catholic traditions
  - Celtic festivals: Saint Brigid's Day, May Day, Lughnasadh, Samhain
  - Bank holidays on first/last Mondays of specific months
  - Irish Gaelic/English bilingual support
- **Israel (IL)**: 11 holidays with Hebrew calendar system
  - Major Jewish holidays: Rosh Hashanah, Yom Kippur, Passover, Shavuot
  - Memorial sequence: Holocaust Remembrance Day, Memorial Day, Independence Day
  - Hebrew/English bilingual support

### Enhanced Features
- **Multi-language Support**: Added Spanish, Irish Gaelic (GA), and Hebrew (HE)
- **Variable Holiday Laws**: Support for Chilean holiday movement and Irish bank holiday patterns
- **Hebrew Calendar Integration**: Accurate calculations for Jewish holidays (2023-2026)
- **Regional Variations**: Enhanced support for region-specific holiday introductions
- **Comprehensive Error Handling**: Structured errors with context support
- **Test Coverage**: 100% test coverage for all new countries

### Library Statistics
- **Total Countries**: 37 (was 34, +3 new)
- **Regional Subdivisions**: 600+ supported
- **Languages**: 15+ languages supported
- **Holiday Coverage**: 1000+ holidays across all countries

### Performance
- Sub-microsecond holiday lookups with intelligent caching
- Thread-safe concurrent operations
- Memory efficient with lazy loading

### Compatibility
- Fully backward compatible with existing API
- Enhanced API with error handling and context support
- chronogo integration maintained

### Installation
```bash
go get github.com/coredds/goholiday@v0.6.3
```

### Quick Example
```go
// Chile Independence Day
chile := goholidays.NewCountry("CL")
if holiday, ok := chile.IsHoliday(time.Date(2024, 9, 18, 0, 0, 0, 0, time.UTC)); ok {
    fmt.Println(holiday.Languages["es"]) // "Dia de la Independencia"
}

// Ireland Saint Patrick's Day  
ireland := goholidays.NewCountry("IE")
if holiday, ok := ireland.IsHoliday(time.Date(2024, 3, 17, 0, 0, 0, 0, time.UTC)); ok {
    fmt.Println(holiday.Languages["ga"]) // "La Fheile Padraig"
}

// Israel Passover
israel := goholidays.NewCountry("IL")
if holiday, ok := israel.IsHoliday(time.Date(2024, 4, 23, 0, 0, 0, 0, time.UTC)); ok {
    fmt.Println(holiday.Languages["he"]) // Hebrew text
}
```

See the [CHANGELOG](https://github.com/coredds/goholiday/blob/main/CHANGELOG.md) for complete details.
