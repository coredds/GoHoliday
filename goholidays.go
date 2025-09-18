// Package goholidays provides comprehensive holiday data for countries and their subdivisions.
// It offers a Go-native implementation inspired by the Python holidays library, providing
// accurate and up-to-date holiday information with high performance and low memory footprint.
package goholidays

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/coredds/GoHoliday/countries"
)

// Version represents the current version of the GoHoliday library
const Version = "0.5.3"

// ErrorCode represents different types of errors that can occur
type ErrorCode int

const (
	// ErrInvalidCountry indicates an unsupported or invalid country code
	ErrInvalidCountry ErrorCode = iota

	// ErrInvalidYear indicates an invalid year value
	ErrInvalidYear

	// ErrDataLoadFailed indicates failure to load holiday data
	ErrDataLoadFailed

	// ErrCancelled indicates the operation was cancelled via context
	ErrCancelled

	// ErrInvalidDate indicates an invalid date parameter
	ErrInvalidDate

	// ErrProviderNotFound indicates no provider exists for the country
	ErrProviderNotFound
)

// HolidayError represents a structured error with context about what went wrong
type HolidayError struct {
	Code    ErrorCode
	Country string
	Year    int
	Date    string
	Message string
	Cause   error
}

// Error implements the error interface
func (e *HolidayError) Error() string {
	if e.Cause != nil {
		return fmt.Sprintf("%s: %v", e.Message, e.Cause)
	}
	return e.Message
}

// Unwrap returns the underlying cause error
func (e *HolidayError) Unwrap() error {
	return e.Cause
}

// Is checks if the error matches a specific error code
func (e *HolidayError) Is(target error) bool {
	if he, ok := target.(*HolidayError); ok {
		return e.Code == he.Code
	}
	return false
}

// NewHolidayError creates a new HolidayError
func NewHolidayError(code ErrorCode, message string) *HolidayError {
	return &HolidayError{
		Code:    code,
		Message: message,
	}
}

// NewHolidayErrorWithCause creates a new HolidayError with an underlying cause
func NewHolidayErrorWithCause(code ErrorCode, message string, cause error) *HolidayError {
	return &HolidayError{
		Code:    code,
		Message: message,
		Cause:   cause,
	}
}

// NewCountryError creates a country-specific error
func NewCountryError(code ErrorCode, country, message string) *HolidayError {
	return &HolidayError{
		Code:    code,
		Country: country,
		Message: message,
	}
}

// NewYearError creates a year-specific error
func NewYearError(code ErrorCode, country string, year int, message string) *HolidayError {
	return &HolidayError{
		Code:    code,
		Country: country,
		Year:    year,
		Message: message,
	}
}

// SupportedCountries contains all countries that have holiday providers
var SupportedCountries = map[string]bool{
	"AR": true, "AT": true, "AU": true, "BE": true, "BR": true, "CA": true,
	"CH": true, "CN": true, "DE": true, "ES": true, "FI": true, "FR": true,
	"GB": true, "ID": true, "IN": true, "IT": true, "JP": true, "KR": true,
	"MX": true, "NL": true, "NO": true, "NZ": true, "PL": true, "PT": true,
	"RU": true, "SE": true, "SG": true, "TH": true, "TR": true, "UA": true,
	"US": true,
}

// ValidateCountryCode checks if a country code is supported
func ValidateCountryCode(code string) error {
	if code == "" {
		return NewCountryError(ErrInvalidCountry, code, "country code cannot be empty")
	}

	if !SupportedCountries[code] {
		return NewCountryError(ErrInvalidCountry, code,
			fmt.Sprintf("country code '%s' is not supported", code))
	}

	return nil
}

// ValidateYear checks if a year is valid for holiday calculations
func ValidateYear(year int) error {
	// Reasonable bounds for holiday calculations
	if year < 1900 || year > 2200 {
		return NewHolidayError(ErrInvalidYear,
			fmt.Sprintf("year %d is outside valid range (1900-2200)", year))
	}
	return nil
}

// IsContextCancelled checks if an error is due to context cancellation
func IsContextCancelled(err error) bool {
	if err == context.Canceled || err == context.DeadlineExceeded {
		return true
	}

	if he, ok := err.(*HolidayError); ok {
		return he.Code == ErrCancelled
	}

	return false
}

// WrapContextError wraps a context error as a HolidayError
func WrapContextError(err error) error {
	if err == nil {
		return nil
	}

	if err == context.Canceled {
		return NewHolidayErrorWithCause(ErrCancelled, "operation was cancelled", err)
	}

	if err == context.DeadlineExceeded {
		return NewHolidayErrorWithCause(ErrCancelled, "operation timed out", err)
	}

	return err
}

// IsValidCountry checks if a country code is supported
func IsValidCountry(countryCode string) bool {
	return SupportedCountries[countryCode]
}

// GetSupportedCountries returns a list of all supported country codes
func GetSupportedCountries() []string {
	countries := make([]string, 0, len(SupportedCountries))
	for code := range SupportedCountries {
		countries = append(countries, code)
	}
	return countries
}

// HolidayCategory represents different types of holidays
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

// Holiday represents a single holiday with its properties
type Holiday struct {
	Name       string            `json:"name"`
	Date       time.Time         `json:"date"`
	Category   HolidayCategory   `json:"category"`
	Observed   *time.Time        `json:"observed,omitempty"`
	Languages  map[string]string `json:"languages,omitempty"`
	IsObserved bool              `json:"is_observed"`
}

// Country represents a country's holiday provider with thread-safe caching
type Country struct {
	code         string
	subdivisions []string
	years        map[int]map[time.Time]*Holiday
	categories   []HolidayCategory
	language     string
	mu           sync.RWMutex // Protects concurrent access to years map
}

// CountryOptions provides configuration options for creating a Country
type CountryOptions struct {
	Subdivisions []string
	Categories   []HolidayCategory
	Language     string
	Years        []int
}

// NewCountry creates a new Country holiday provider
// Note: For error handling, use NewCountryWithError instead
func NewCountry(countryCode string, options ...CountryOptions) *Country {
	c := &Country{
		code:       countryCode,
		years:      make(map[int]map[time.Time]*Holiday),
		categories: []HolidayCategory{CategoryPublic},
		language:   "en",
	}

	if len(options) > 0 {
		opt := options[0]
		if opt.Subdivisions != nil {
			c.subdivisions = opt.Subdivisions
		}
		if opt.Categories != nil {
			c.categories = opt.Categories
		}
		if opt.Language != "" {
			c.language = opt.Language
		}
		if opt.Years != nil {
			c.loadYears(opt.Years)
		}
	}

	return c
}

// IsHoliday checks if the given date is a holiday (thread-safe)
func (c *Country) IsHoliday(date time.Time) (*Holiday, bool) {
	year := date.Year()

	// First, try to read with read lock
	c.mu.RLock()
	holidays, exists := c.years[year]
	c.mu.RUnlock()

	if exists {
		// Normalize date to compare only year, month, day
		dateKey := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, time.UTC)
		if holiday, found := holidays[dateKey]; found {
			return holiday, true
		}
	} else {
		// Load holidays for this year if not already loaded
		c.loadYear(year)
		return c.IsHoliday(date)
	}
	return nil, false
}

// HolidaysForYear returns all holidays for a specific year (thread-safe)
func (c *Country) HolidaysForYear(year int) map[time.Time]*Holiday {
	c.mu.RLock()
	holidays, exists := c.years[year]
	c.mu.RUnlock()

	if exists {
		// Return a copy to prevent external modification
		result := make(map[time.Time]*Holiday, len(holidays))
		for k, v := range holidays {
			result[k] = v
		}
		return result
	}

	c.loadYear(year)

	c.mu.RLock()
	defer c.mu.RUnlock()

	// Return a copy to prevent external modification
	result := make(map[time.Time]*Holiday, len(c.years[year]))
	for k, v := range c.years[year] {
		result[k] = v
	}
	return result
}

// HolidaysForDateRange returns all holidays within a date range
func (c *Country) HolidaysForDateRange(start, end time.Time) map[time.Time]*Holiday {
	result := make(map[time.Time]*Holiday)

	startYear := start.Year()
	endYear := end.Year()

	for year := startYear; year <= endYear; year++ {
		yearHolidays := c.HolidaysForYear(year)
		for date, holiday := range yearHolidays {
			if (date.After(start) || date.Equal(start)) && (date.Before(end) || date.Equal(end)) {
				result[date] = holiday
			}
		}
	}

	return result
}

// GetCountryCode returns the country code
func (c *Country) GetCountryCode() string {
	return c.code
}

// GetSubdivisions returns the subdivisions
func (c *Country) GetSubdivisions() []string {
	return c.subdivisions
}

// GetCategories returns the holiday categories
func (c *Country) GetCategories() []HolidayCategory {
	return c.categories
}

// GetLanguage returns the current language
func (c *Country) GetLanguage() string {
	return c.language
}

// loadYear loads holidays for a specific year (thread-safe)
func (c *Country) loadYear(year int) {
	// Double-checked locking pattern for performance
	c.mu.RLock()
	_, exists := c.years[year]
	c.mu.RUnlock()

	if exists {
		return // Already loaded
	}

	c.mu.Lock()
	defer c.mu.Unlock()

	// Check again after acquiring write lock
	if c.years[year] == nil {
		c.years[year] = make(map[time.Time]*Holiday)
		c.loadCountryHolidays(year)
	}
}

// loadYears loads holidays for multiple years
func (c *Country) loadYears(years []int) {
	for _, year := range years {
		c.loadYear(year)
	}
}

// loadCountryHolidays loads country-specific holidays using the countries package
func (c *Country) loadCountryHolidays(year int) {
	// This will be integrated with the countries package in the full implementation
	// For now, we use the placeholder implementation
	switch c.code {
	case "US":
		c.loadUSHolidays(year)
	case "GB":
		c.loadGBHolidays(year)
	case "CA":
		c.loadCAHolidays(year)
	case "AU":
		c.loadAUHolidays(year)
	case "NZ":
		c.loadNZHolidays(year)
	case "JP":
		c.loadJPHolidays(year)
	case "IN":
		c.loadINHolidays(year)
	case "FR":
		c.loadFRHolidays(year)
	case "DE":
		c.loadDEHolidays(year)
	case "BR":
		c.loadBRHolidays(year)
	case "MX":
		c.loadMXHolidays(year)
	case "IT":
		c.loadITHolidays(year)
	case "ES":
		c.loadESHolidays(year)
	case "NL":
		c.loadNLHolidays(year)
	case "KR":
		c.loadKRHolidays(year)
	case "UA":
		c.loadUAHolidays(year)
	// Add more countries as needed
	default:
		// Load from generic holiday data or return empty
	}
}

// loadUSHolidays loads US holidays using the US provider
func (c *Country) loadUSHolidays(year int) {
	provider := countries.NewUSProvider()
	holidayMap := provider.LoadHolidays(year)

	for date, holiday := range holidayMap {
		c.years[year][date] = &Holiday{
			Name:       holiday.Name,
			Date:       holiday.Date,
			Category:   HolidayCategory(holiday.Category),
			Languages:  holiday.Languages,
			Observed:   holiday.Observed,
			IsObserved: holiday.IsObserved,
		}
	}
}

func (c *Country) loadGBHolidays(year int) {
	holidays := c.years[year]

	// New Year's Day
	holidays[time.Date(year, 1, 1, 0, 0, 0, 0, time.UTC)] = &Holiday{
		Name:     "New Year's Day",
		Date:     time.Date(year, 1, 1, 0, 0, 0, 0, time.UTC),
		Category: CategoryPublic,
		Languages: map[string]string{
			"en": "New Year's Day",
		},
	}

	// Christmas Day
	holidays[time.Date(year, 12, 25, 0, 0, 0, 0, time.UTC)] = &Holiday{
		Name:     "Christmas Day",
		Date:     time.Date(year, 12, 25, 0, 0, 0, 0, time.UTC),
		Category: CategoryPublic,
		Languages: map[string]string{
			"en": "Christmas Day",
		},
	}
}

func (c *Country) loadCAHolidays(year int) {
	holidays := c.years[year]

	// New Year's Day
	holidays[time.Date(year, 1, 1, 0, 0, 0, 0, time.UTC)] = &Holiday{
		Name:     "New Year's Day",
		Date:     time.Date(year, 1, 1, 0, 0, 0, 0, time.UTC),
		Category: CategoryPublic,
		Languages: map[string]string{
			"en": "New Year's Day",
			"fr": "Jour de l'An",
		},
	}

	// Canada Day
	holidays[time.Date(year, 7, 1, 0, 0, 0, 0, time.UTC)] = &Holiday{
		Name:     "Canada Day",
		Date:     time.Date(year, 7, 1, 0, 0, 0, 0, time.UTC),
		Category: CategoryPublic,
		Languages: map[string]string{
			"en": "Canada Day",
			"fr": "Fête du Canada",
		},
	}

	// Christmas Day
	holidays[time.Date(year, 12, 25, 0, 0, 0, 0, time.UTC)] = &Holiday{
		Name:     "Christmas Day",
		Date:     time.Date(year, 12, 25, 0, 0, 0, 0, time.UTC),
		Category: CategoryPublic,
		Languages: map[string]string{
			"en": "Christmas Day",
			"fr": "Noël",
		},
	}

	// Thanksgiving Day - 2nd Monday in October
	thanksgiving := c.getNthWeekdayOfMonth(year, 10, time.Monday, 2)
	holidays[thanksgiving] = &Holiday{
		Name:     "Thanksgiving Day",
		Date:     thanksgiving,
		Category: CategoryPublic,
		Languages: map[string]string{
			"en": "Thanksgiving Day",
			"fr": "Action de grâce",
		},
	}
}

func (c *Country) loadAUHolidays(year int) {
	holidays := c.years[year]

	// New Year's Day
	holidays[time.Date(year, 1, 1, 0, 0, 0, 0, time.UTC)] = &Holiday{
		Name:     "New Year's Day",
		Date:     time.Date(year, 1, 1, 0, 0, 0, 0, time.UTC),
		Category: CategoryPublic,
		Languages: map[string]string{
			"en": "New Year's Day",
		},
	}

	// Australia Day
	holidays[time.Date(year, 1, 26, 0, 0, 0, 0, time.UTC)] = &Holiday{
		Name:     "Australia Day",
		Date:     time.Date(year, 1, 26, 0, 0, 0, 0, time.UTC),
		Category: CategoryPublic,
		Languages: map[string]string{
			"en": "Australia Day",
		},
	}

	// ANZAC Day
	holidays[time.Date(year, 4, 25, 0, 0, 0, 0, time.UTC)] = &Holiday{
		Name:     "ANZAC Day",
		Date:     time.Date(year, 4, 25, 0, 0, 0, 0, time.UTC),
		Category: CategoryPublic,
		Languages: map[string]string{
			"en": "ANZAC Day",
		},
	}

	// Queen's Birthday - 2nd Monday in June (most states)
	queensBirthday := c.getNthWeekdayOfMonth(year, 6, time.Monday, 2)
	holidays[queensBirthday] = &Holiday{
		Name:     "Queen's Birthday",
		Date:     queensBirthday,
		Category: CategoryPublic,
		Languages: map[string]string{
			"en": "Queen's Birthday",
		},
	}

	// Labour Day - 1st Monday in October (most states)
	labourDay := c.getNthWeekdayOfMonth(year, 10, time.Monday, 1)
	holidays[labourDay] = &Holiday{
		Name:     "Labour Day",
		Date:     labourDay,
		Category: CategoryPublic,
		Languages: map[string]string{
			"en": "Labour Day",
		},
	}

	// Christmas Day
	holidays[time.Date(year, 12, 25, 0, 0, 0, 0, time.UTC)] = &Holiday{
		Name:     "Christmas Day",
		Date:     time.Date(year, 12, 25, 0, 0, 0, 0, time.UTC),
		Category: CategoryPublic,
		Languages: map[string]string{
			"en": "Christmas Day",
		},
	}

	// Boxing Day
	holidays[time.Date(year, 12, 26, 0, 0, 0, 0, time.UTC)] = &Holiday{
		Name:     "Boxing Day",
		Date:     time.Date(year, 12, 26, 0, 0, 0, 0, time.UTC),
		Category: CategoryPublic,
		Languages: map[string]string{
			"en": "Boxing Day",
		},
	}
}

func (c *Country) loadNZHolidays(year int) {
	holidays := c.years[year]

	// New Year's Day
	holidays[time.Date(year, 1, 1, 0, 0, 0, 0, time.UTC)] = &Holiday{
		Name:     "New Year's Day",
		Date:     time.Date(year, 1, 1, 0, 0, 0, 0, time.UTC),
		Category: CategoryPublic,
		Languages: map[string]string{
			"en": "New Year's Day",
			"mi": "Te Rā Tau Hou",
		},
	}

	// Day after New Year's Day
	holidays[time.Date(year, 1, 2, 0, 0, 0, 0, time.UTC)] = &Holiday{
		Name:     "Day after New Year's Day",
		Date:     time.Date(year, 1, 2, 0, 0, 0, 0, time.UTC),
		Category: CategoryPublic,
		Languages: map[string]string{
			"en": "Day after New Year's Day",
			"mi": "Te Rā i muri i te Tau Hou",
		},
	}

	// Waitangi Day
	holidays[time.Date(year, 2, 6, 0, 0, 0, 0, time.UTC)] = &Holiday{
		Name:     "Waitangi Day",
		Date:     time.Date(year, 2, 6, 0, 0, 0, 0, time.UTC),
		Category: CategoryPublic,
		Languages: map[string]string{
			"en": "Waitangi Day",
			"mi": "Te Rā o Waitangi",
		},
	}

	// Good Friday (Easter-based)
	easter := c.easterSunday(year)
	goodFriday := easter.AddDate(0, 0, -2)
	holidays[goodFriday] = &Holiday{
		Name:     "Good Friday",
		Date:     goodFriday,
		Category: CategoryPublic,
		Languages: map[string]string{
			"en": "Good Friday",
			"mi": "Rārā Pai",
		},
	}

	// Easter Monday
	easterMonday := easter.AddDate(0, 0, 1)
	holidays[easterMonday] = &Holiday{
		Name:     "Easter Monday",
		Date:     easterMonday,
		Category: CategoryPublic,
		Languages: map[string]string{
			"en": "Easter Monday",
			"mi": "Rā Aranga Rērā",
		},
	}

	// ANZAC Day
	holidays[time.Date(year, 4, 25, 0, 0, 0, 0, time.UTC)] = &Holiday{
		Name:     "ANZAC Day",
		Date:     time.Date(year, 4, 25, 0, 0, 0, 0, time.UTC),
		Category: CategoryPublic,
		Languages: map[string]string{
			"en": "ANZAC Day",
			"mi": "Te Rā o nga Hoia",
		},
	}

	// Queen's Birthday - First Monday in June
	queensBirthday := c.getNthWeekdayOfMonth(year, 6, time.Monday, 1)
	holidays[queensBirthday] = &Holiday{
		Name:     "Queen's Birthday",
		Date:     queensBirthday,
		Category: CategoryPublic,
		Languages: map[string]string{
			"en": "Queen's Birthday",
			"mi": "Te Rā Whānau o te Kuini",
		},
	}

	// Matariki - Known astronomical dates for certain years
	matarikiDates := map[int]time.Time{
		2022: time.Date(2022, 6, 24, 0, 0, 0, 0, time.UTC),
		2023: time.Date(2023, 7, 14, 0, 0, 0, 0, time.UTC),
		2024: time.Date(2024, 6, 28, 0, 0, 0, 0, time.UTC),
		2025: time.Date(2025, 6, 20, 0, 0, 0, 0, time.UTC),
		2026: time.Date(2026, 7, 10, 0, 0, 0, 0, time.UTC),
		2027: time.Date(2027, 6, 25, 0, 0, 0, 0, time.UTC),
		2028: time.Date(2028, 7, 14, 0, 0, 0, 0, time.UTC),
		2029: time.Date(2029, 7, 6, 0, 0, 0, 0, time.UTC),
		2030: time.Date(2030, 6, 21, 0, 0, 0, 0, time.UTC),
	}

	if matarikiDate, exists := matarikiDates[year]; exists {
		holidays[matarikiDate] = &Holiday{
			Name:     "Matariki",
			Date:     matarikiDate,
			Category: CategoryPublic,
			Languages: map[string]string{
				"en": "Matariki",
				"mi": "Matariki",
			},
		}
	}

	// Labour Day - Fourth Monday in October
	labourDay := c.getNthWeekdayOfMonth(year, 10, time.Monday, 4)
	holidays[labourDay] = &Holiday{
		Name:     "Labour Day",
		Date:     labourDay,
		Category: CategoryPublic,
		Languages: map[string]string{
			"en": "Labour Day",
			"mi": "Te Rā Whakatōhea",
		},
	}

	// Christmas Day
	holidays[time.Date(year, 12, 25, 0, 0, 0, 0, time.UTC)] = &Holiday{
		Name:     "Christmas Day",
		Date:     time.Date(year, 12, 25, 0, 0, 0, 0, time.UTC),
		Category: CategoryPublic,
		Languages: map[string]string{
			"en": "Christmas Day",
			"mi": "Te Rā Kirihimete",
		},
	}

	// Boxing Day
	holidays[time.Date(year, 12, 26, 0, 0, 0, 0, time.UTC)] = &Holiday{
		Name:     "Boxing Day",
		Date:     time.Date(year, 12, 26, 0, 0, 0, 0, time.UTC),
		Category: CategoryPublic,
		Languages: map[string]string{
			"en": "Boxing Day",
			"mi": "Te Rā Pouaka",
		},
	}
}

func (c *Country) loadJPHolidays(year int) {
	holidays := c.years[year]

	// New Year's Day (元日, Ganjitsu)
	holidays[time.Date(year, 1, 1, 0, 0, 0, 0, time.UTC)] = &Holiday{
		Name:     "New Year's Day",
		Date:     time.Date(year, 1, 1, 0, 0, 0, 0, time.UTC),
		Category: CategoryPublic,
		Languages: map[string]string{
			"en": "New Year's Day",
			"ja": "元日",
		},
	}

	// Coming of Age Day (成人の日, Seijin no Hi) - Second Monday of January
	comingOfAge := c.getNthWeekdayOfMonth(year, 1, time.Monday, 2)
	holidays[comingOfAge] = &Holiday{
		Name:     "Coming of Age Day",
		Date:     comingOfAge,
		Category: CategoryPublic,
		Languages: map[string]string{
			"en": "Coming of Age Day",
			"ja": "成人の日",
		},
	}

	// National Foundation Day (建国記念の日, Kenkoku Kinen no Hi)
	holidays[time.Date(year, 2, 11, 0, 0, 0, 0, time.UTC)] = &Holiday{
		Name:     "National Foundation Day",
		Date:     time.Date(year, 2, 11, 0, 0, 0, 0, time.UTC),
		Category: CategoryPublic,
		Languages: map[string]string{
			"en": "National Foundation Day",
			"ja": "建国記念の日",
		},
	}

	// Emperor's Birthday (天皇誕生日, Tennō Tanjōbi)
	var emperorBirthday time.Time
	if year >= 2020 {
		emperorBirthday = time.Date(year, 2, 23, 0, 0, 0, 0, time.UTC) // Emperor Naruhito
	} else {
		emperorBirthday = time.Date(year, 12, 23, 0, 0, 0, 0, time.UTC) // Emperor Akihito
	}
	holidays[emperorBirthday] = &Holiday{
		Name:     "Emperor's Birthday",
		Date:     emperorBirthday,
		Category: CategoryPublic,
		Languages: map[string]string{
			"en": "Emperor's Birthday",
			"ja": "天皇誕生日",
		},
	}

	// Constitution Memorial Day (憲法記念日, Kenpō Kinenbi)
	holidays[time.Date(year, 5, 3, 0, 0, 0, 0, time.UTC)] = &Holiday{
		Name:     "Constitution Memorial Day",
		Date:     time.Date(year, 5, 3, 0, 0, 0, 0, time.UTC),
		Category: CategoryPublic,
		Languages: map[string]string{
			"en": "Constitution Memorial Day",
			"ja": "憲法記念日",
		},
	}

	// Greenery Day (みどりの日, Midori no Hi)
	holidays[time.Date(year, 5, 4, 0, 0, 0, 0, time.UTC)] = &Holiday{
		Name:     "Greenery Day",
		Date:     time.Date(year, 5, 4, 0, 0, 0, 0, time.UTC),
		Category: CategoryPublic,
		Languages: map[string]string{
			"en": "Greenery Day",
			"ja": "みどりの日",
		},
	}

	// Children's Day (こどもの日, Kodomo no Hi)
	holidays[time.Date(year, 5, 5, 0, 0, 0, 0, time.UTC)] = &Holiday{
		Name:     "Children's Day",
		Date:     time.Date(year, 5, 5, 0, 0, 0, 0, time.UTC),
		Category: CategoryPublic,
		Languages: map[string]string{
			"en": "Children's Day",
			"ja": "こどもの日",
		},
	}

	// Marine Day (海の日, Umi no Hi) - Third Monday of July
	marine := c.getNthWeekdayOfMonth(year, 7, time.Monday, 3)
	holidays[marine] = &Holiday{
		Name:     "Marine Day",
		Date:     marine,
		Category: CategoryPublic,
		Languages: map[string]string{
			"en": "Marine Day",
			"ja": "海の日",
		},
	}

	// Sports Day (スポーツの日, Supōtsu no Hi) - Second Monday of October
	sports := c.getNthWeekdayOfMonth(year, 10, time.Monday, 2)
	sportsName := "Sports Day"
	if year < 2020 {
		sportsName = "Health and Sports Day"
	}
	holidays[sports] = &Holiday{
		Name:     sportsName,
		Date:     sports,
		Category: CategoryPublic,
		Languages: map[string]string{
			"en": sportsName,
			"ja": "スポーツの日",
		},
	}

	// Culture Day (文化の日, Bunka no Hi)
	holidays[time.Date(year, 11, 3, 0, 0, 0, 0, time.UTC)] = &Holiday{
		Name:     "Culture Day",
		Date:     time.Date(year, 11, 3, 0, 0, 0, 0, time.UTC),
		Category: CategoryPublic,
		Languages: map[string]string{
			"en": "Culture Day",
			"ja": "文化の日",
		},
	}
}

// getNthWeekdayOfMonth is a helper method for calculating variable holidays
func (c *Country) getNthWeekdayOfMonth(year int, month time.Month, weekday time.Weekday, n int) time.Time {
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

	panic("Invalid n value for getNthWeekdayOfMonth")
}

// easterSunday calculates Easter Sunday for a given year using the Western (Gregorian) algorithm
func (c *Country) easterSunday(year int) time.Time {
	// Anonymous Gregorian algorithm (Western Easter)
	a := year % 19
	b := year / 100
	c2 := year % 100
	d := b / 4
	e := b % 4
	f := (b + 8) / 25
	g := (b - f + 1) / 3
	h := (19*a + b - d - g + 15) % 30
	i := c2 / 4
	k := c2 % 4
	l := (32 + 2*e + 2*i - h - k) % 7
	m := (a + 11*h + 22*l) / 451
	n := (h + l - 7*m + 114) / 31
	p := (h + l - 7*m + 114) % 31

	return time.Date(year, time.Month(n), p+1, 0, 0, 0, 0, time.UTC)
}

// loadINHolidays loads holidays specific to India
func (c *Country) loadINHolidays(year int) {
	holidays := c.years[year]

	// New Year's Day
	holidays[time.Date(year, 1, 1, 0, 0, 0, 0, time.UTC)] = &Holiday{
		Name:     "New Year's Day",
		Date:     time.Date(year, 1, 1, 0, 0, 0, 0, time.UTC),
		Category: CategoryPublic,
		Languages: map[string]string{
			"en": "New Year's Day",
			"hi": "नव वर्ष दिवस",
		},
	}

	// Republic Day
	holidays[time.Date(year, 1, 26, 0, 0, 0, 0, time.UTC)] = &Holiday{
		Name:     "Republic Day",
		Date:     time.Date(year, 1, 26, 0, 0, 0, 0, time.UTC),
		Category: CategoryPublic,
		Languages: map[string]string{
			"en": "Republic Day",
			"hi": "गणतंत्र दिवस",
		},
	}

	// Independence Day
	holidays[time.Date(year, 8, 15, 0, 0, 0, 0, time.UTC)] = &Holiday{
		Name:     "Independence Day",
		Date:     time.Date(year, 8, 15, 0, 0, 0, 0, time.UTC),
		Category: CategoryPublic,
		Languages: map[string]string{
			"en": "Independence Day",
			"hi": "स्वतंत्रता दिवस",
		},
	}

	// Gandhi Jayanti
	holidays[time.Date(year, 10, 2, 0, 0, 0, 0, time.UTC)] = &Holiday{
		Name:     "Gandhi Jayanti",
		Date:     time.Date(year, 10, 2, 0, 0, 0, 0, time.UTC),
		Category: CategoryPublic,
		Languages: map[string]string{
			"en": "Gandhi Jayanti",
			"hi": "गांधी जयंती",
		},
	}

	// Christmas Day
	holidays[time.Date(year, 12, 25, 0, 0, 0, 0, time.UTC)] = &Holiday{
		Name:     "Christmas Day",
		Date:     time.Date(year, 12, 25, 0, 0, 0, 0, time.UTC),
		Category: CategoryReligious,
		Languages: map[string]string{
			"en": "Christmas Day",
			"hi": "क्रिसमस",
		},
	}

	// Good Friday (varies each year based on Easter calculation)
	easter := c.easterSunday(year)
	goodFriday := easter.AddDate(0, 0, -2)
	holidays[goodFriday] = &Holiday{
		Name:     "Good Friday",
		Date:     goodFriday,
		Category: CategoryReligious,
		Languages: map[string]string{
			"en": "Good Friday",
			"hi": "गुड फ्राइडे",
		},
	}

	// Note: Religious festivals like Diwali, Holi, Eid are approximated
	// In a full implementation, these would use proper lunar calendar calculations

	// Diwali (approximate - typically October/November)
	// This is a simplified calculation; actual dates vary based on lunar calendar
	diwaliDate := c.approximateDiwali(year)
	holidays[diwaliDate] = &Holiday{
		Name:     "Diwali",
		Date:     diwaliDate,
		Category: CategoryReligious,
		Languages: map[string]string{
			"en": "Diwali",
			"hi": "दीवाली",
		},
	}

	// Holi (approximate - typically March)
	// This is a simplified calculation; actual dates vary based on lunar calendar
	holiDate := c.approximateHoli(year)
	holidays[holiDate] = &Holiday{
		Name:     "Holi",
		Date:     holiDate,
		Category: CategoryReligious,
		Languages: map[string]string{
			"en": "Holi",
			"hi": "होली",
		},
	}
}

// approximateDiwali provides an approximate date for Diwali
// Note: In a production system, this should use proper lunar calendar calculations
func (c *Country) approximateDiwali(year int) time.Time {
	// Diwali typically falls in October/November
	// This is a very rough approximation for demonstration
	switch year {
	case 2024:
		return time.Date(year, 11, 1, 0, 0, 0, 0, time.UTC)
	case 2025:
		return time.Date(year, 10, 20, 0, 0, 0, 0, time.UTC)
	case 2026:
		return time.Date(year, 11, 8, 0, 0, 0, 0, time.UTC)
	default:
		// Default approximation: third week of October
		return time.Date(year, 10, 21, 0, 0, 0, 0, time.UTC)
	}
}

// approximateHoli provides an approximate date for Holi
// Note: In a production system, this should use proper lunar calendar calculations
func (c *Country) approximateHoli(year int) time.Time {
	// Holi typically falls in March
	// This is a very rough approximation for demonstration
	switch year {
	case 2024:
		return time.Date(year, 3, 25, 0, 0, 0, 0, time.UTC)
	case 2025:
		return time.Date(year, 3, 14, 0, 0, 0, 0, time.UTC)
	case 2026:
		return time.Date(year, 3, 3, 0, 0, 0, 0, time.UTC)
	default:
		// Default approximation: second week of March
		return time.Date(year, 3, 14, 0, 0, 0, 0, time.UTC)
	}
}

// loadFRHolidays loads holidays specific to France
func (c *Country) loadFRHolidays(year int) {
	holidays := c.years[year]

	// New Year's Day
	holidays[time.Date(year, 1, 1, 0, 0, 0, 0, time.UTC)] = &Holiday{
		Name:     "Jour de l'An",
		Date:     time.Date(year, 1, 1, 0, 0, 0, 0, time.UTC),
		Category: CategoryPublic,
		Languages: map[string]string{
			"en": "New Year's Day",
			"fr": "Jour de l'An",
		},
	}

	// Labour Day
	holidays[time.Date(year, 5, 1, 0, 0, 0, 0, time.UTC)] = &Holiday{
		Name:     "Fête du Travail",
		Date:     time.Date(year, 5, 1, 0, 0, 0, 0, time.UTC),
		Category: CategoryPublic,
		Languages: map[string]string{
			"en": "Labour Day",
			"fr": "Fête du Travail",
		},
	}

	// Victory in Europe Day
	holidays[time.Date(year, 5, 8, 0, 0, 0, 0, time.UTC)] = &Holiday{
		Name:     "Fête de la Victoire",
		Date:     time.Date(year, 5, 8, 0, 0, 0, 0, time.UTC),
		Category: CategoryPublic,
		Languages: map[string]string{
			"en": "Victory in Europe Day",
			"fr": "Fête de la Victoire",
		},
	}

	// Bastille Day
	holidays[time.Date(year, 7, 14, 0, 0, 0, 0, time.UTC)] = &Holiday{
		Name:     "Fête nationale",
		Date:     time.Date(year, 7, 14, 0, 0, 0, 0, time.UTC),
		Category: CategoryPublic,
		Languages: map[string]string{
			"en": "Bastille Day",
			"fr": "Fête nationale",
		},
	}

	// Assumption of Mary
	holidays[time.Date(year, 8, 15, 0, 0, 0, 0, time.UTC)] = &Holiday{
		Name:     "Assomption",
		Date:     time.Date(year, 8, 15, 0, 0, 0, 0, time.UTC),
		Category: CategoryReligious,
		Languages: map[string]string{
			"en": "Assumption of Mary",
			"fr": "Assomption",
		},
	}

	// All Saints' Day
	holidays[time.Date(year, 11, 1, 0, 0, 0, 0, time.UTC)] = &Holiday{
		Name:     "Toussaint",
		Date:     time.Date(year, 11, 1, 0, 0, 0, 0, time.UTC),
		Category: CategoryReligious,
		Languages: map[string]string{
			"en": "All Saints' Day",
			"fr": "Toussaint",
		},
	}

	// Armistice Day
	holidays[time.Date(year, 11, 11, 0, 0, 0, 0, time.UTC)] = &Holiday{
		Name:     "Armistice",
		Date:     time.Date(year, 11, 11, 0, 0, 0, 0, time.UTC),
		Category: CategoryPublic,
		Languages: map[string]string{
			"en": "Armistice Day",
			"fr": "Armistice",
		},
	}

	// Christmas Day
	holidays[time.Date(year, 12, 25, 0, 0, 0, 0, time.UTC)] = &Holiday{
		Name:     "Noël",
		Date:     time.Date(year, 12, 25, 0, 0, 0, 0, time.UTC),
		Category: CategoryReligious,
		Languages: map[string]string{
			"en": "Christmas Day",
			"fr": "Noël",
		},
	}

	// Easter-based holidays
	easter := c.easterSunday(year)

	// Easter Monday
	easterMonday := easter.AddDate(0, 0, 1)
	holidays[easterMonday] = &Holiday{
		Name:     "Lundi de Pâques",
		Date:     easterMonday,
		Category: CategoryReligious,
		Languages: map[string]string{
			"en": "Easter Monday",
			"fr": "Lundi de Pâques",
		},
	}

	// Ascension Day (39 days after Easter)
	ascension := easter.AddDate(0, 0, 39)
	holidays[ascension] = &Holiday{
		Name:     "Ascension",
		Date:     ascension,
		Category: CategoryReligious,
		Languages: map[string]string{
			"en": "Ascension Day",
			"fr": "Ascension",
		},
	}

	// Whit Monday (50 days after Easter)
	whitMonday := easter.AddDate(0, 0, 50)
	holidays[whitMonday] = &Holiday{
		Name:     "Lundi de Pentecôte",
		Date:     whitMonday,
		Category: CategoryReligious,
		Languages: map[string]string{
			"en": "Whit Monday",
			"fr": "Lundi de Pentecôte",
		},
	}
}

// loadDEHolidays loads Germany holidays using the DE provider
func (c *Country) loadDEHolidays(year int) {
	provider := countries.NewDEProvider()
	holidayMap := provider.LoadHolidays(year)

	for date, holiday := range holidayMap {
		c.years[year][date] = &Holiday{
			Name:      holiday.Name,
			Date:      holiday.Date,
			Category:  HolidayCategory(holiday.Category),
			Languages: holiday.Languages,
		}
	}
}

// loadBRHolidays loads Brazil holidays using the BR provider
func (c *Country) loadBRHolidays(year int) {
	provider := countries.NewBRProvider()
	holidayMap := provider.LoadHolidays(year)

	for date, holiday := range holidayMap {
		c.years[year][date] = &Holiday{
			Name:      holiday.Name,
			Date:      holiday.Date,
			Category:  HolidayCategory(holiday.Category),
			Languages: holiday.Languages,
		}
	}
}

// loadMXHolidays loads Mexico holidays using the MX provider
func (c *Country) loadMXHolidays(year int) {
	provider := countries.NewMXProvider()
	holidayMap := provider.LoadHolidays(year)

	for date, holiday := range holidayMap {
		c.years[year][date] = &Holiday{
			Name:      holiday.Name,
			Date:      holiday.Date,
			Category:  HolidayCategory(holiday.Category),
			Languages: holiday.Languages,
		}
	}
}

// loadITHolidays loads Italy holidays using the IT provider
func (c *Country) loadITHolidays(year int) {
	provider := countries.NewITProvider()
	holidayMap := provider.LoadHolidays(year)

	for date, holiday := range holidayMap {
		c.years[year][date] = &Holiday{
			Name:      holiday.Name,
			Date:      holiday.Date,
			Category:  HolidayCategory(holiday.Category),
			Languages: holiday.Languages,
		}
	}
}

// loadESHolidays loads Spain holidays using the ES provider
func (c *Country) loadESHolidays(year int) {
	provider := countries.NewESProvider()
	holidayMap := provider.LoadHolidays(year)

	for date, holiday := range holidayMap {
		c.years[year][date] = &Holiday{
			Name:      holiday.Name,
			Date:      holiday.Date,
			Category:  HolidayCategory(holiday.Category),
			Languages: holiday.Languages,
		}
	}
}

// loadNLHolidays loads Netherlands holidays using the NL provider
func (c *Country) loadNLHolidays(year int) {
	provider := countries.NewNLProvider()
	holidayMap := provider.LoadHolidays(year)

	for date, holiday := range holidayMap {
		c.years[year][date] = &Holiday{
			Name:      holiday.Name,
			Date:      holiday.Date,
			Category:  HolidayCategory(holiday.Category),
			Languages: holiday.Languages,
		}
	}
}

// loadKRHolidays loads South Korea holidays using the KR provider
func (c *Country) loadKRHolidays(year int) {
	provider := countries.NewKRProvider()
	holidayMap := provider.LoadHolidays(year)

	for date, holiday := range holidayMap {
		c.years[year][date] = &Holiday{
			Name:      holiday.Name,
			Date:      holiday.Date,
			Category:  HolidayCategory(holiday.Category),
			Languages: holiday.Languages,
		}
	}
}

// loadUAHolidays loads Ukraine holidays using the UA provider
func (c *Country) loadUAHolidays(year int) {
	provider := countries.NewUAProvider()
	holidayMap := provider.LoadHolidays(year)

	for date, holiday := range holidayMap {
		c.years[year][date] = &Holiday{
			Name:      holiday.Name,
			Date:      holiday.Date,
			Category:  HolidayCategory(holiday.Category),
			Languages: holiday.Languages,
		}
	}
}

// ============================================================================
// Enhanced API Methods with Error Handling and Context Support
// ============================================================================

// NewCountryWithError creates a new Country with validation
// This is the recommended way to create countries with proper error handling
func NewCountryWithError(countryCode string, options ...CountryOptions) (*Country, error) {
	// Validate country code
	if err := ValidateCountryCode(countryCode); err != nil {
		return nil, err
	}

	// Use existing NewCountry function
	country := NewCountry(countryCode, options...)
	return country, nil
}

// IsHolidayWithError checks if the given date is a holiday with error handling
func (c *Country) IsHolidayWithError(date time.Time) (*Holiday, bool, error) {
	year := date.Year()

	// Validate year
	if err := ValidateYear(year); err != nil {
		return nil, false, err
	}

	// Use existing IsHoliday method
	holiday, isHoliday := c.IsHoliday(date)
	return holiday, isHoliday, nil
}

// IsHolidayWithContext checks if the given date is a holiday with context support
func (c *Country) IsHolidayWithContext(ctx context.Context, date time.Time) (*Holiday, bool, error) {
	// Check for context cancellation
	select {
	case <-ctx.Done():
		return nil, false, WrapContextError(ctx.Err())
	default:
	}

	year := date.Year()

	// Validate year
	if err := ValidateYear(year); err != nil {
		return nil, false, err
	}

	// Load year with context if needed
	if err := c.loadYearWithContext(ctx, year); err != nil {
		return nil, false, err
	}

	// Use existing IsHoliday method for the actual lookup
	holiday, isHoliday := c.IsHoliday(date)
	return holiday, isHoliday, nil
}

// HolidaysForYearWithError returns all holidays for a specific year with error handling
func (c *Country) HolidaysForYearWithError(year int) (map[time.Time]*Holiday, error) {
	// Validate year
	if err := ValidateYear(year); err != nil {
		return nil, err
	}

	// Use existing HolidaysForYear method
	holidays := c.HolidaysForYear(year)
	return holidays, nil
}

// HolidaysForYearWithContext returns all holidays for a specific year with context support
func (c *Country) HolidaysForYearWithContext(ctx context.Context, year int) (map[time.Time]*Holiday, error) {
	// Check for context cancellation
	select {
	case <-ctx.Done():
		return nil, WrapContextError(ctx.Err())
	default:
	}

	// Validate year
	if err := ValidateYear(year); err != nil {
		return nil, err
	}

	// Load year with context
	if err := c.loadYearWithContext(ctx, year); err != nil {
		return nil, err
	}

	// Use existing HolidaysForYear method
	holidays := c.HolidaysForYear(year)
	return holidays, nil
}

// HolidaysForDateRangeWithError returns all holidays within a date range with error handling
func (c *Country) HolidaysForDateRangeWithError(start, end time.Time) (map[time.Time]*Holiday, error) {
	// Validate date range
	if start.After(end) {
		return nil, NewHolidayError(ErrInvalidDate, "start date cannot be after end date")
	}

	startYear := start.Year()
	endYear := end.Year()

	// Validate years
	if err := ValidateYear(startYear); err != nil {
		return nil, err
	}
	if err := ValidateYear(endYear); err != nil {
		return nil, err
	}

	// Use existing HolidaysForDateRange method
	holidays := c.HolidaysForDateRange(start, end)
	return holidays, nil
}

// HolidaysForDateRangeWithContext returns all holidays within a date range with context support
func (c *Country) HolidaysForDateRangeWithContext(ctx context.Context, start, end time.Time) (map[time.Time]*Holiday, error) {
	// Check for context cancellation
	select {
	case <-ctx.Done():
		return nil, WrapContextError(ctx.Err())
	default:
	}

	// Validate date range
	if start.After(end) {
		return nil, NewHolidayError(ErrInvalidDate, "start date cannot be after end date")
	}

	startYear := start.Year()
	endYear := end.Year()

	// Validate and load all required years with context
	for year := startYear; year <= endYear; year++ {
		select {
		case <-ctx.Done():
			return nil, WrapContextError(ctx.Err())
		default:
		}

		if err := ValidateYear(year); err != nil {
			return nil, err
		}

		if err := c.loadYearWithContext(ctx, year); err != nil {
			return nil, err
		}
	}

	// Use existing HolidaysForDateRange method
	holidays := c.HolidaysForDateRange(start, end)
	return holidays, nil
}

// GetHolidayCount returns the number of holidays for a given year
func (c *Country) GetHolidayCount(year int) (int, error) {
	holidays, err := c.HolidaysForYearWithError(year)
	if err != nil {
		return 0, err
	}
	return len(holidays), nil
}

// GetHolidayCountWithContext returns the number of holidays for a given year with context
func (c *Country) GetHolidayCountWithContext(ctx context.Context, year int) (int, error) {
	holidays, err := c.HolidaysForYearWithContext(ctx, year)
	if err != nil {
		return 0, err
	}
	return len(holidays), nil
}

// loadYearWithContext loads holidays for a specific year with context support
func (c *Country) loadYearWithContext(ctx context.Context, year int) error {
	// Check if already loaded
	c.mu.RLock()
	_, exists := c.years[year]
	c.mu.RUnlock()

	if exists {
		return nil // Already loaded
	}

	// Check for context cancellation before acquiring write lock
	select {
	case <-ctx.Done():
		return WrapContextError(ctx.Err())
	default:
	}

	c.mu.Lock()
	defer c.mu.Unlock()

	// Double-check after acquiring write lock
	if c.years[year] != nil {
		return nil // Already loaded by another goroutine
	}

	// Check for context cancellation again
	select {
	case <-ctx.Done():
		return WrapContextError(ctx.Err())
	default:
	}

	// Use existing loadYear method but with error handling
	c.years[year] = make(map[time.Time]*Holiday)

	// Validate country code and load holidays
	if err := ValidateCountryCode(c.code); err != nil {
		delete(c.years, year)
		return err
	}

	// Use existing loadCountryHolidays method
	c.loadCountryHolidays(year)

	return nil
}
