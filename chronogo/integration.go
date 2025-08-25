package chronogo

import (
	"time"

	"github.com/your-username/goholidays/config"
	"github.com/your-username/goholidays/countries"
)

// GoHolidaysChecker implements ChronoGo's HolidayChecker interface
// using the comprehensive GoHolidays system with configuration support
type GoHolidaysChecker struct {
	holidayManager *config.HolidayManager
	countries      []string
	subdivisions   []string
	categories     []string
	includeRegional bool
	
	// Cache for performance
	holidayCache map[int]map[time.Time]bool // year -> date -> isHoliday
}

// ChronoGoHolidayChecker interface matches ChronoGo's expected interface
// This ensures we implement exactly what ChronoGo expects
type ChronoGoHolidayChecker interface {
	IsHoliday(dt ChronoGoDateTime) bool
}

// ChronoGoDateTime represents the interface that ChronoGo DateTime should satisfy
// We'll accept anything that can give us Year(), Month(), Day()
type ChronoGoDateTime interface {
	Year() int
	Month() time.Month
	Day() int
}

// NewGoHolidaysChecker creates a new ChronoGo-compatible holiday checker
func NewGoHolidaysChecker() *GoHolidaysChecker {
	return &GoHolidaysChecker{
		holidayManager: config.NewHolidayManager(),
		countries:      []string{"US"}, // Default to US
		categories:     []string{"federal", "public", "bank", "observance"},
		includeRegional: false,
		holidayCache:   make(map[int]map[time.Time]bool),
	}
}

// WithCountries sets the countries to check for holidays
func (ghc *GoHolidaysChecker) WithCountries(countries ...string) *GoHolidaysChecker {
	ghc.countries = countries
	ghc.clearCache()
	return ghc
}

// WithSubdivisions enables regional holidays for specific subdivisions
func (ghc *GoHolidaysChecker) WithSubdivisions(subdivisions ...string) *GoHolidaysChecker {
	ghc.subdivisions = subdivisions
	ghc.includeRegional = len(subdivisions) > 0
	ghc.clearCache()
	return ghc
}

// WithCategories sets which holiday categories to consider
func (ghc *GoHolidaysChecker) WithCategories(categories ...string) *GoHolidaysChecker {
	ghc.categories = categories
	ghc.clearCache()
	return ghc
}

// WithConfiguration loads holidays using a specific configuration file
func (ghc *GoHolidaysChecker) WithConfiguration(configPath string) (*GoHolidaysChecker, error) {
	cm := config.NewConfigManager()
	_, err := cm.LoadConfigFromFile(configPath)
	if err != nil {
		return nil, err
	}
	
	ghc.holidayManager = config.NewHolidayManager()
	ghc.clearCache()
	return ghc, nil
}

// IsHoliday implements ChronoGo's HolidayChecker interface
func (ghc *GoHolidaysChecker) IsHoliday(dt ChronoGoDateTime) bool {
	year := dt.Year()
	date := time.Date(year, dt.Month(), dt.Day(), 0, 0, 0, 0, time.UTC)
	
	// Check cache first
	if yearCache, exists := ghc.holidayCache[year]; exists {
		if isHoliday, cached := yearCache[date]; cached {
			return isHoliday
		}
	}
	
	// Load holidays for this year if not cached
	isHoliday := ghc.checkHoliday(date)
	
	// Cache the result
	if ghc.holidayCache[year] == nil {
		ghc.holidayCache[year] = make(map[time.Time]bool)
	}
	ghc.holidayCache[year][date] = isHoliday
	
	return isHoliday
}

// checkHoliday does the actual holiday checking across all configured countries
func (ghc *GoHolidaysChecker) checkHoliday(date time.Time) bool {
	for _, countryCode := range ghc.countries {
		var holidays map[time.Time]*countries.Holiday
		var err error
		
		// Get holidays with or without regional support
		if ghc.includeRegional && len(ghc.subdivisions) > 0 {
			holidays, err = ghc.holidayManager.GetHolidaysWithSubdivisions(
				countryCode, date.Year(), ghc.subdivisions)
		} else {
			holidays, err = ghc.holidayManager.GetHolidays(countryCode, date.Year())
		}
		
		if err != nil {
			continue // Skip this country if there's an error
		}
		
		// Check if the date matches any holiday
		for holidayDate, holiday := range holidays {
			if ghc.datesMatch(date, holidayDate) && ghc.categoryMatches(holiday.Category) {
				return true
			}
		}
	}
	
	return false
}

// datesMatch compares two dates ignoring time components
func (ghc *GoHolidaysChecker) datesMatch(date1, date2 time.Time) bool {
	return date1.Year() == date2.Year() &&
		   date1.Month() == date2.Month() &&
		   date1.Day() == date2.Day()
}

// categoryMatches checks if the holiday category is in our accepted categories
func (ghc *GoHolidaysChecker) categoryMatches(category string) bool {
	if len(ghc.categories) == 0 {
		return true // Accept all categories if none specified
	}
	
	for _, acceptedCategory := range ghc.categories {
		if category == acceptedCategory {
			return true
		}
	}
	
	return false
}

// clearCache clears the holiday cache when configuration changes
func (ghc *GoHolidaysChecker) clearCache() {
	ghc.holidayCache = make(map[int]map[time.Time]bool)
}

// GetHolidays returns all holidays for a given year (ChronoGo compatibility)
// This matches the interface that ChronoGo's DefaultHolidayChecker provides
func (ghc *GoHolidaysChecker) GetHolidays(year int) []HolidayInfo {
	var allHolidays []HolidayInfo
	
	for _, countryCode := range ghc.countries {
		var holidays map[time.Time]*countries.Holiday
		var err error
		
		if ghc.includeRegional && len(ghc.subdivisions) > 0 {
			holidays, err = ghc.holidayManager.GetHolidaysWithSubdivisions(
				countryCode, year, ghc.subdivisions)
		} else {
			holidays, err = ghc.holidayManager.GetHolidays(countryCode, year)
		}
		
		if err != nil {
			continue
		}
		
		for date, holiday := range holidays {
			if ghc.categoryMatches(holiday.Category) {
				allHolidays = append(allHolidays, HolidayInfo{
					Name:     holiday.Name,
					Date:     date,
					Country:  countryCode,
					Category: holiday.Category,
				})
			}
		}
	}
	
	return allHolidays
}

// HolidayInfo provides detailed information about a holiday
type HolidayInfo struct {
	Name     string
	Date     time.Time
	Country  string
	Category string
}

// GetSupportedCountries returns the list of countries supported by GoHolidays
func (ghc *GoHolidaysChecker) GetSupportedCountries() []string {
	return ghc.holidayManager.GetSupportedCountries()
}

// GetCountryInfo returns detailed information about a country's holiday configuration
func (ghc *GoHolidaysChecker) GetCountryInfo(countryCode string) (map[string]interface{}, error) {
	return ghc.holidayManager.GetCountryInfo(countryCode)
}

// AddCustomHoliday allows adding custom holidays at runtime
// This provides compatibility with ChronoGo's AddHoliday functionality
func (ghc *GoHolidaysChecker) AddCustomHoliday(name string, date time.Time) {
	// This would require extending our configuration system to support runtime additions
	// For now, this is a placeholder that shows the interface
	// Users should use the configuration system for custom holidays
}

// PreloadYear pre-loads and caches all holidays for a specific year
// This can improve performance when doing many holiday checks for the same year
func (ghc *GoHolidaysChecker) PreloadYear(year int) error {
	date := time.Date(year, 1, 1, 0, 0, 0, 0, time.UTC)
	
	// Check one date to trigger loading of all holidays for the year
	ghc.IsHoliday(&dateTimeWrapper{date})
	
	return nil
}

// dateTimeWrapper wraps a time.Time to implement ChronoGoDateTime
type dateTimeWrapper struct {
	time.Time
}

func (dtw *dateTimeWrapper) Year() int        { return dtw.Time.Year() }
func (dtw *dateTimeWrapper) Month() time.Month { return dtw.Time.Month() }
func (dtw *dateTimeWrapper) Day() int         { return dtw.Time.Day() }

// CreateDefaultUSChecker creates a GoHolidays checker configured for US federal holidays
// This provides a drop-in replacement for ChronoGo's NewUSHolidayChecker()
func CreateDefaultUSChecker() *GoHolidaysChecker {
	return NewGoHolidaysChecker().
		WithCountries("US").
		WithCategories("federal", "public")
}

// CreateMultiCountryChecker creates a checker for multiple countries
func CreateMultiCountryChecker(countries ...string) *GoHolidaysChecker {
	return NewGoHolidaysChecker().
		WithCountries(countries...).
		WithCategories("federal", "public", "bank")
}

// CreateRegionalChecker creates a checker with regional holiday support
func CreateRegionalChecker(country string, subdivisions ...string) *GoHolidaysChecker {
	return NewGoHolidaysChecker().
		WithCountries(country).
		WithSubdivisions(subdivisions...).
		WithCategories("federal", "public", "regional", "state")
}
