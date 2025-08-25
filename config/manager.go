package config

import (
	"fmt"
	"strings"
	"time"

	"github.com/coredds/GoHoliday/countries"
)

// HolidayManager integrates configuration with holiday providers
type HolidayManager struct {
	configManager *ConfigManager
	providers     map[string]countries.HolidayProvider
}

// NewHolidayManager creates a new holiday manager with configuration
func NewHolidayManager() *HolidayManager {
	cm := NewConfigManager()
	config, _ := cm.LoadConfig()

	hm := &HolidayManager{
		configManager: cm,
		providers:     make(map[string]countries.HolidayProvider),
	}

	// Initialize providers based on configuration
	hm.initializeProviders(config)

	return hm
}

// initializeProviders creates holiday providers for enabled countries
func (hm *HolidayManager) initializeProviders(config *Config) {
	// Always create providers, but respect enabled status when using them
	hm.providers["AU"] = countries.NewAUProvider()
	hm.providers["CA"] = countries.NewCAProvider()
	hm.providers["NZ"] = countries.NewNZProvider()
	hm.providers["GB"] = countries.NewGBProvider()
	hm.providers["US"] = countries.NewUSProvider()
	hm.providers["DE"] = countries.NewDEProvider()
	hm.providers["FR"] = countries.NewFRProvider()
}

// GetHolidays returns holidays for a country with configuration applied
func (hm *HolidayManager) GetHolidays(countryCode string, year int) (map[time.Time]*countries.Holiday, error) {
	config := hm.configManager.GetConfig()

	// Check if country is enabled
	if !hm.configManager.IsCountryEnabled(countryCode) {
		return nil, fmt.Errorf("country %s is not enabled", countryCode)
	}

	// Get the provider
	provider, exists := hm.providers[countryCode]
	if !exists {
		return nil, fmt.Errorf("no provider available for country %s", countryCode)
	}

	// Load base holidays
	holidays := provider.LoadHolidays(year)

	// Apply configuration overrides
	holidays = hm.applyCountryConfig(holidays, countryCode, config)

	// Add custom holidays
	customHolidays := hm.getCustomHolidays(countryCode, year, config)
	for date, holiday := range customHolidays {
		holidays[date] = holiday
	}

	// Apply output formatting
	holidays = hm.applyOutputFormatting(holidays, config)

	return holidays, nil
}

// GetHolidaysWithSubdivisions returns holidays including regional ones
func (hm *HolidayManager) GetHolidaysWithSubdivisions(countryCode string, year int, subdivisions []string) (map[time.Time]*countries.Holiday, error) {
	// Get base holidays
	holidays, err := hm.GetHolidays(countryCode, year)
	if err != nil {
		return nil, err
	}

	// Get country configuration
	countryConfig := hm.configManager.GetCountryConfig(countryCode)

	// Filter subdivisions based on configuration
	allowedSubdivisions := hm.filterSubdivisions(subdivisions, countryConfig.Subdivisions)

	if len(allowedSubdivisions) == 0 {
		return holidays, nil
	}

	// Get regional holidays based on provider type
	var regionalHolidays map[time.Time]*countries.Holiday

	switch provider := hm.providers[countryCode].(type) {
	case *countries.USProvider:
		regionalHolidays = provider.GetStateHolidays(year, allowedSubdivisions)
	case *countries.GBProvider:
		regionalHolidays = provider.GetRegionalHolidays(year, allowedSubdivisions)
	case *countries.DEProvider:
		regionalHolidays = provider.GetRegionalHolidays(year, allowedSubdivisions)
	case *countries.FRProvider:
		regionalHolidays = provider.GetRegionalHolidays(year, allowedSubdivisions)
		// Add other providers as needed
	}

	// Merge regional holidays
	for date, holiday := range regionalHolidays {
		holidays[date] = holiday
	}

	return holidays, nil
}

// applyCountryConfig applies country-specific configuration
func (hm *HolidayManager) applyCountryConfig(holidays map[time.Time]*countries.Holiday, countryCode string, config *Config) map[time.Time]*countries.Holiday {
	countryConfig := hm.configManager.GetCountryConfig(countryCode)

	// Apply holiday name overrides
	for _, holiday := range holidays {
		if newName, exists := countryConfig.Overrides[holiday.Name]; exists {
			holiday.Name = newName
		}
	}

	// Exclude holidays
	for _, excludedName := range countryConfig.ExcludedHolidays {
		for date, holiday := range holidays {
			if holiday.Name == excludedName {
				delete(holidays, date)
			}
		}
	}

	// Filter by categories
	if len(countryConfig.Categories) > 0 {
		filteredHolidays := make(map[time.Time]*countries.Holiday)
		for date, holiday := range holidays {
			for _, allowedCategory := range countryConfig.Categories {
				if holiday.Category == allowedCategory {
					filteredHolidays[date] = holiday
					break
				}
			}
		}
		holidays = filteredHolidays
	}

	return holidays
}

// getCustomHolidays processes custom holidays from configuration
func (hm *HolidayManager) getCustomHolidays(countryCode string, year int, config *Config) map[time.Time]*countries.Holiday {
	holidays := make(map[time.Time]*countries.Holiday)
	customHolidays := hm.configManager.GetCustomHolidays(countryCode)

	for _, custom := range customHolidays {
		// Check year range
		if custom.YearRange != nil {
			if custom.YearRange.Start > 0 && year < custom.YearRange.Start {
				continue
			}
			if custom.YearRange.End > 0 && year > custom.YearRange.End {
				continue
			}
		}

		// Calculate the date
		date, err := hm.calculateCustomHolidayDate(custom, year)
		if err != nil {
			continue // Skip invalid dates
		}

		// Create the holiday
		holiday := &countries.Holiday{
			Name:      custom.Name,
			Date:      date,
			Category:  custom.Category,
			Languages: custom.Languages,
		}

		holidays[date] = holiday
	}

	return holidays
}

// calculateCustomHolidayDate calculates the date for a custom holiday
func (hm *HolidayManager) calculateCustomHolidayDate(custom CustomHoliday, year int) (time.Time, error) {
	if custom.Date != "" {
		// Fixed date - parse YYYY-MM-DD or MM-DD
		if strings.Contains(custom.Date, fmt.Sprintf("%d-", year)) {
			return time.Parse("2006-01-02", custom.Date)
		} else if len(custom.Date) == 5 { // MM-DD format
			return time.Parse("2006-01-02", fmt.Sprintf("%d-%s", year, custom.Date))
		}
	}

	if custom.Calculation != nil {
		switch custom.Calculation.Type {
		case "easter_offset":
			easter := countries.EasterSunday(year)
			return easter.AddDate(0, 0, custom.Calculation.EasterOffset), nil

		case "weekday":
			if custom.Calculation.WeekdayRule != nil {
				weekday := parseWeekday(custom.Calculation.WeekdayRule.Weekday)
				return countries.NthWeekdayOfMonth(year, time.Month(custom.Calculation.WeekdayRule.Month),
					weekday, custom.Calculation.WeekdayRule.Week), nil
			}

		case "fixed":
			return time.Date(year, time.Month(custom.Calculation.Month), 1, 0, 0, 0, 0, time.UTC), nil
		}
	}

	return time.Time{}, fmt.Errorf("unable to calculate date for custom holiday %s", custom.Name)
}

// parseWeekday converts string weekday to time.Weekday
func parseWeekday(weekdayStr string) time.Weekday {
	switch strings.ToLower(weekdayStr) {
	case "sunday":
		return time.Sunday
	case "monday":
		return time.Monday
	case "tuesday":
		return time.Tuesday
	case "wednesday":
		return time.Wednesday
	case "thursday":
		return time.Thursday
	case "friday":
		return time.Friday
	case "saturday":
		return time.Saturday
	default:
		return time.Monday // Default fallback
	}
}

// filterSubdivisions filters subdivisions based on configuration
func (hm *HolidayManager) filterSubdivisions(requested, allowed []string) []string {
	if len(allowed) == 0 {
		return requested // No restrictions
	}

	var filtered []string
	for _, req := range requested {
		for _, allow := range allowed {
			if req == allow {
				filtered = append(filtered, req)
				break
			}
		}
	}

	return filtered
}

// applyOutputFormatting applies output formatting configuration
func (hm *HolidayManager) applyOutputFormatting(holidays map[time.Time]*countries.Holiday, config *Config) map[time.Time]*countries.Holiday {
	// Note: We don't apply timezone conversion to holiday dates as they should remain
	// as calendar dates (midnight UTC) to avoid changing the actual holiday date.
	// Timezone conversion should only be applied when displaying times, not dates.
	
	return holidays
}

// GetSupportedCountries returns list of supported countries
func (hm *HolidayManager) GetSupportedCountries() []string {
	var countries []string

	for countryCode := range hm.providers {
		if hm.configManager.IsCountryEnabled(countryCode) {
			countries = append(countries, countryCode)
		}
	}

	return countries
}

// GetCountryInfo returns information about a country's configuration
func (hm *HolidayManager) GetCountryInfo(countryCode string) (map[string]interface{}, error) {
	provider, exists := hm.providers[countryCode]
	if !exists {
		return nil, fmt.Errorf("country %s not supported", countryCode)
	}

	countryConfig := hm.configManager.GetCountryConfig(countryCode)

	info := map[string]interface{}{
		"country_code":            countryCode,
		"enabled":                 countryConfig.Enabled,
		"supported_subdivisions":  provider.GetSupportedSubdivisions(),
		"supported_categories":    provider.GetSupportedCategories(),
		"configured_subdivisions": countryConfig.Subdivisions,
		"configured_categories":   countryConfig.Categories,
		"excluded_holidays":       countryConfig.ExcludedHolidays,
		"holiday_overrides":       countryConfig.Overrides,
	}

	return info, nil
}
