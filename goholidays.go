// Package goholidays provides comprehensive holiday data for countries and their subdivisions.
// It offers a Go-native implementation inspired by the Python holidays library, providing
// accurate and up-to-date holiday information with high performance and low memory footprint.
package goholidays

import (
	"time"
)

// HolidayCategory represents different types of holidays
type HolidayCategory string

const (
	CategoryPublic     HolidayCategory = "public"
	CategoryBank       HolidayCategory = "bank"
	CategorySchool     HolidayCategory = "school"
	CategoryGovernment HolidayCategory = "government"
	CategoryReligious  HolidayCategory = "religious"
	CategoryOptional   HolidayCategory = "optional"
	CategoryHalfDay    HolidayCategory = "half_day"
	CategoryArmedForces HolidayCategory = "armed_forces"
	CategoryWorkday    HolidayCategory = "workday"
)

// Holiday represents a single holiday with its properties
type Holiday struct {
	Name        string            `json:"name"`
	Date        time.Time         `json:"date"`
	Category    HolidayCategory   `json:"category"`
	Observed    *time.Time        `json:"observed,omitempty"`
	Languages   map[string]string `json:"languages,omitempty"`
	IsObserved  bool              `json:"is_observed"`
}

// Country represents a country's holiday provider
type Country struct {
	code         string
	subdivisions []string
	years        map[int]map[time.Time]*Holiday
	categories   []HolidayCategory
	language     string
}

// CountryOptions provides configuration options for creating a Country
type CountryOptions struct {
	Subdivisions []string
	Categories   []HolidayCategory
	Language     string
	Years        []int
}

// NewCountry creates a new Country holiday provider
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

// IsHoliday checks if the given date is a holiday
func (c *Country) IsHoliday(date time.Time) (*Holiday, bool) {
	year := date.Year()
	if holidays, exists := c.years[year]; exists {
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

// HolidaysForYear returns all holidays for a specific year
func (c *Country) HolidaysForYear(year int) map[time.Time]*Holiday {
	if holidays, exists := c.years[year]; exists {
		return holidays
	}
	
	c.loadYear(year)
	return c.years[year]
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

// loadYear loads holidays for a specific year (placeholder implementation)
func (c *Country) loadYear(year int) {
	if c.years[year] == nil {
		c.years[year] = make(map[time.Time]*Holiday)
	}
	
	// TODO: Implement actual holiday loading logic
	// This will be replaced with country-specific holiday calculations
	c.loadCountryHolidays(year)
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
	// Add more countries as needed
	default:
		// Load from generic holiday data or return empty
	}
}

// Placeholder implementations for specific countries
func (c *Country) loadUSHolidays(year int) {
	holidays := c.years[year]
	
	// New Year's Day
	holidays[time.Date(year, 1, 1, 0, 0, 0, 0, time.UTC)] = &Holiday{
		Name:     "New Year's Day",
		Date:     time.Date(year, 1, 1, 0, 0, 0, 0, time.UTC),
		Category: CategoryPublic,
		Languages: map[string]string{
			"en": "New Year's Day",
			"es": "Año Nuevo",
		},
	}
	
	// Independence Day
	holidays[time.Date(year, 7, 4, 0, 0, 0, 0, time.UTC)] = &Holiday{
		Name:     "Independence Day",
		Date:     time.Date(year, 7, 4, 0, 0, 0, 0, time.UTC),
		Category: CategoryPublic,
		Languages: map[string]string{
			"en": "Independence Day",
			"es": "Día de la Independencia",
		},
	}
	
	// Christmas Day
	holidays[time.Date(year, 12, 25, 0, 0, 0, 0, time.UTC)] = &Holiday{
		Name:     "Christmas Day",
		Date:     time.Date(year, 12, 25, 0, 0, 0, 0, time.UTC),
		Category: CategoryPublic,
		Languages: map[string]string{
			"en": "Christmas Day",
			"es": "Navidad",
		},
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
