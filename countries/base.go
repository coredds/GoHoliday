// Package countries provides country-specific holiday implementations
package countries

import (
	"time"
)

// HolidayProvider defines the interface for country-specific holiday providers
type HolidayProvider interface {
	LoadHolidays(year int) map[time.Time]*Holiday
	GetCountryCode() string
	GetSupportedSubdivisions() []string
	GetSupportedCategories() []string
}

// Holiday represents a holiday with all its properties
type Holiday struct {
	Name        string            `json:"name"`
	Date        time.Time         `json:"date"`
	Category    string            `json:"category"`
	Observed    *time.Time        `json:"observed,omitempty"`
	Languages   map[string]string `json:"languages,omitempty"`
	IsObserved  bool              `json:"is_observed"`
	Subdivisions []string         `json:"subdivisions,omitempty"`
}

// BaseProvider provides common functionality for holiday providers
type BaseProvider struct {
	countryCode   string
	subdivisions  []string
	categories    []string
	observedShift bool
}

// NewBaseProvider creates a new base provider
func NewBaseProvider(countryCode string) *BaseProvider {
	return &BaseProvider{
		countryCode:   countryCode,
		subdivisions:  []string{},
		categories:    []string{"public"},
		observedShift: true,
	}
}

// GetCountryCode returns the country code
func (bp *BaseProvider) GetCountryCode() string {
	return bp.countryCode
}

// GetSupportedSubdivisions returns supported subdivisions
func (bp *BaseProvider) GetSupportedSubdivisions() []string {
	return bp.subdivisions
}

// GetSupportedCategories returns supported categories
func (bp *BaseProvider) GetSupportedCategories() []string {
	return bp.categories
}

// CalculateObservedDate calculates the observed date for a holiday
func (bp *BaseProvider) CalculateObservedDate(date time.Time) *time.Time {
	if !bp.observedShift {
		return nil
	}
	
	weekday := date.Weekday()
	var observed time.Time
	
	switch weekday {
	case time.Saturday:
		observed = date.AddDate(0, 0, -1) // Friday
	case time.Sunday:
		observed = date.AddDate(0, 0, 1) // Monday
	default:
		return nil // No shift needed
	}
	
	return &observed
}

// CreateHoliday creates a new holiday with standard properties
func (bp *BaseProvider) CreateHoliday(name string, date time.Time, category string, languages map[string]string) *Holiday {
	holiday := &Holiday{
		Name:      name,
		Date:      date,
		Category:  category,
		Languages: languages,
	}
	
	if observed := bp.CalculateObservedDate(date); observed != nil {
		holiday.Observed = observed
		holiday.IsObserved = true
	}
	
	return holiday
}

// EasterSunday calculates Easter Sunday for a given year using the Western calendar
func EasterSunday(year int) time.Time {
	// Anonymous Gregorian algorithm
	a := year % 19
	b := year / 100
	c := year % 100
	d := b / 4
	e := b % 4
	f := (b + 8) / 25
	g := (b - f + 1) / 3
	h := (19*a + b - d - g + 15) % 30
	i := c / 4
	k := c % 4
	l := (32 + 2*e + 2*i - h - k) % 7
	m := (a + 11*h + 22*l) / 451
	month := (h + l - 7*m + 114) / 31
	day := ((h + l - 7*m + 114) % 31) + 1
	
	return time.Date(year, time.Month(month), day, 0, 0, 0, 0, time.UTC)
}

// GoodFriday calculates Good Friday (2 days before Easter)
func GoodFriday(year int) time.Time {
	easter := EasterSunday(year)
	return easter.AddDate(0, 0, -2)
}

// EasterMonday calculates Easter Monday (1 day after Easter)
func EasterMonday(year int) time.Time {
	easter := EasterSunday(year)
	return easter.AddDate(0, 0, 1)
}

// NthWeekdayOfMonth calculates the nth occurrence of a weekday in a month
// n=1 for first, n=2 for second, etc. Use n=-1 for last occurrence
func NthWeekdayOfMonth(year int, month time.Month, weekday time.Weekday, n int) time.Time {
	if n > 0 {
		// Find the first occurrence of the weekday in the month
		firstDay := time.Date(year, month, 1, 0, 0, 0, 0, time.UTC)
		daysToWeekday := (int(weekday) - int(firstDay.Weekday()) + 7) % 7
		firstOccurrence := firstDay.AddDate(0, 0, daysToWeekday)
		
		// Add weeks to get the nth occurrence
		return firstOccurrence.AddDate(0, 0, (n-1)*7)
	} else if n == -1 {
		// Find the last occurrence of the weekday in the month
		lastDay := time.Date(year, month+1, 0, 0, 0, 0, 0, time.UTC) // Last day of current month
		daysBack := (int(lastDay.Weekday()) - int(weekday) + 7) % 7
		return lastDay.AddDate(0, 0, -daysBack)
	}
	
	panic("Invalid n value for NthWeekdayOfMonth")
}
