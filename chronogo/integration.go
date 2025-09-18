package chronogo

import (
	"sync"
	"time"

	"github.com/coredds/GoHoliday/countries"
)

// HolidayChecker provides the minimal interface ChronoGo needs
type HolidayChecker interface {
	IsHoliday(date time.Time) bool
	GetHolidayName(date time.Time) string
	AreHolidays(dates []time.Time) []bool
}

// FastCountryChecker - optimized for ChronoGo's business day loops
// Provides O(1) holiday lookups with intelligent caching
type FastCountryChecker struct {
	countryCode string
	provider    countries.HolidayProvider
	yearCache   map[int]map[time.Time]*countries.Holiday
	mutex       sync.RWMutex
}

// Checker creates an optimized holiday checker for a specific country
func Checker(countryCode string) *FastCountryChecker {
	var provider countries.HolidayProvider

	switch countryCode {
	case "US":
		provider = countries.NewUSProvider()
	case "CA":
		provider = countries.NewCAProvider()
	case "GB":
		provider = countries.NewGBProvider()
	case "AU":
		provider = countries.NewAUProvider()
	case "NZ":
		provider = countries.NewNZProvider()
	case "DE":
		provider = countries.NewDEProvider()
	case "FR":
		provider = countries.NewFRProvider()
	case "JP":
		provider = countries.NewJPProvider()
	case "UA":
		provider = countries.NewUAProvider()
	case "CL":
		provider = countries.NewCLProvider()
	case "IE":
		provider = countries.NewIEProvider()
	case "IL":
		provider = countries.NewILProvider()
	default:
		// Fallback to US if country not supported
		provider = countries.NewUSProvider()
	}

	return &FastCountryChecker{
		countryCode: countryCode,
		provider:    provider,
		yearCache:   make(map[int]map[time.Time]*countries.Holiday),
	}
}

// IsHoliday performs fast O(1) holiday lookup optimized for business day calculations
func (f *FastCountryChecker) IsHoliday(date time.Time) bool {
	// Normalize date to UTC midnight for consistent lookups
	normalizedDate := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, time.UTC)
	year := normalizedDate.Year()

	// Check cache first (read lock)
	f.mutex.RLock()
	if yearHolidays, exists := f.yearCache[year]; exists {
		_, isHoliday := yearHolidays[normalizedDate]
		f.mutex.RUnlock()
		return isHoliday
	}
	f.mutex.RUnlock()

	// Cache miss - load year data (write lock)
	f.loadYearData(year)

	// Check again after loading
	f.mutex.RLock()
	yearHolidays := f.yearCache[year]
	_, isHoliday := yearHolidays[normalizedDate]
	f.mutex.RUnlock()

	return isHoliday
}

// GetHolidayName returns the name of the holiday on the given date, empty string if not a holiday
func (f *FastCountryChecker) GetHolidayName(date time.Time) string {
	// Normalize date to UTC midnight
	normalizedDate := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, time.UTC)
	year := normalizedDate.Year()

	// Ensure year data is loaded
	f.mutex.RLock()
	yearHolidays, exists := f.yearCache[year]
	f.mutex.RUnlock()

	if !exists {
		f.loadYearData(year)
		f.mutex.RLock()
		yearHolidays = f.yearCache[year]
		f.mutex.RUnlock()
	}

	if holiday, isHoliday := yearHolidays[normalizedDate]; isHoliday {
		return holiday.Name
	}

	return ""
}

// AreHolidays performs batch holiday checking for efficient range operations
// Optimized for ChronoGo's bulk date processing
func (f *FastCountryChecker) AreHolidays(dates []time.Time) []bool {
	results := make([]bool, len(dates))

	// Group dates by year for efficient batch loading
	yearGroups := make(map[int][]int) // year -> slice of indices
	for i, date := range dates {
		year := date.Year()
		yearGroups[year] = append(yearGroups[year], i)
	}

	// Pre-load all required years
	for year := range yearGroups {
		f.ensureYearLoaded(year)
	}

	// Process all dates with cached data
	f.mutex.RLock()
	for i, date := range dates {
		normalizedDate := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, time.UTC)
		yearHolidays := f.yearCache[normalizedDate.Year()]
		_, isHoliday := yearHolidays[normalizedDate]
		results[i] = isHoliday
	}
	f.mutex.RUnlock()

	return results
}

// GetHolidaysInRange returns all holidays in the specified date range
// Useful for ChronoGo's calendar operations
func (f *FastCountryChecker) GetHolidaysInRange(start, end time.Time) map[time.Time]string {
	holidays := make(map[time.Time]string)

	// Determine year range
	startYear := start.Year()
	endYear := end.Year()

	// Pre-load all years in range
	for year := startYear; year <= endYear; year++ {
		f.ensureYearLoaded(year)
	}

	// Collect holidays in range
	f.mutex.RLock()
	for year := startYear; year <= endYear; year++ {
		yearHolidays := f.yearCache[year]
		for date, holiday := range yearHolidays {
			if (date.Equal(start) || date.After(start)) && (date.Equal(end) || date.Before(end)) {
				holidays[date] = holiday.Name
			}
		}
	}
	f.mutex.RUnlock()

	return holidays
}

// CountHolidaysInRange counts holidays in a date range without allocating holiday data
// Optimized for ChronoGo's business day counting
func (f *FastCountryChecker) CountHolidaysInRange(start, end time.Time) int {
	count := 0

	// Iterate through range checking each date
	for d := start; !d.After(end); d = d.AddDate(0, 0, 1) {
		if f.IsHoliday(d) {
			count++
		}
	}

	return count
}

// GetCountryCode returns the country code this checker is configured for
func (f *FastCountryChecker) GetCountryCode() string {
	return f.countryCode
}

// ClearCache clears the holiday cache to free memory
// Useful for long-running applications
func (f *FastCountryChecker) ClearCache() {
	f.mutex.Lock()
	f.yearCache = make(map[int]map[time.Time]*countries.Holiday)
	f.mutex.Unlock()
}

// loadYearData loads holiday data for a specific year (internal method)
func (f *FastCountryChecker) loadYearData(year int) {
	f.mutex.Lock()
	defer f.mutex.Unlock()

	// Double-check after acquiring write lock
	if _, exists := f.yearCache[year]; exists {
		return
	}

	// Load holidays for the year
	holidays := f.provider.LoadHolidays(year)
	f.yearCache[year] = holidays
}

// ensureYearLoaded ensures year data is loaded (internal method)
func (f *FastCountryChecker) ensureYearLoaded(year int) {
	f.mutex.RLock()
	_, exists := f.yearCache[year]
	f.mutex.RUnlock()

	if !exists {
		f.loadYearData(year)
	}
}
