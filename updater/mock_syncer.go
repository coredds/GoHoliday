package updater

import (
	"context"
	"fmt"
	"strings"
	"time"
)

// MockSyncer is a mock implementation of the Syncer interface for testing
type MockSyncer struct {
	countries    []string
	countryFiles map[string]string
	shouldError  bool
	errorMessage string
}

// NewMockSyncer creates a new mock syncer with default test data
func NewMockSyncer() *MockSyncer {
	return &MockSyncer{
		countries: []string{"US", "GB", "CA", "AU", "DE", "FR"},
		countryFiles: map[string]string{
			"US": mockUSPythonSource,
			"GB": mockGBPythonSource,
			"CA": mockCAPythonSource,
		},
		shouldError: false,
	}
}

// SetError configures the mock to return an error
func (m *MockSyncer) SetError(shouldError bool, message string) {
	m.shouldError = shouldError
	m.errorMessage = message
}

// AddCountry adds a country to the mock data
func (m *MockSyncer) AddCountry(code, source string) {
	// Add to countries list if not already present
	found := false
	for _, c := range m.countries {
		if c == code {
			found = true
			break
		}
	}
	if !found {
		m.countries = append(m.countries, code)
	}

	// Set the source
	m.countryFiles[code] = source
}

// FetchCountryList returns the mock list of countries
func (m *MockSyncer) FetchCountryList(ctx context.Context) ([]string, error) {
	if m.shouldError {
		return nil, fmt.Errorf("mock error: %s", m.errorMessage)
	}

	// Simulate some processing time
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	case <-time.After(10 * time.Millisecond):
		// Continue
	}

	return append([]string{}, m.countries...), nil
}

// FetchCountryFile returns mock Python source for a country
func (m *MockSyncer) FetchCountryFile(ctx context.Context, countryCode string) (string, error) {
	if m.shouldError {
		return "", fmt.Errorf("mock error: %s", m.errorMessage)
	}

	// Simulate some processing time
	select {
	case <-ctx.Done():
		return "", ctx.Err()
	case <-time.After(10 * time.Millisecond):
		// Continue
	}

	source, exists := m.countryFiles[strings.ToUpper(countryCode)]
	if !exists {
		return "", fmt.Errorf("country %s not found", countryCode)
	}

	return source, nil
}

// ParseHolidayDefinitions parses mock Python source into holiday data
func (m *MockSyncer) ParseHolidayDefinitions(source string) (*CountryData, error) {
	if m.shouldError {
		return nil, fmt.Errorf("mock error: %s", m.errorMessage)
	}

	// Simple mock parsing - extract country code from source
	countryCode := "US" // default
	if strings.Contains(source, "class GB") {
		countryCode = "GB"
	} else if strings.Contains(source, "class CA") {
		countryCode = "CA"
	}

	return &CountryData{
		CountryCode:  countryCode,
		Name:         getCountryName(countryCode),
		Subdivisions: map[string]string{},
		Categories:   []string{"public"},
		Languages:    []string{"en"},
		Holidays: map[string]HolidayDefinition{
			"new_years_day": {
				Name:        "New Year's Day",
				Category:    "public",
				Languages:   map[string]string{"en": "New Year's Day"},
				Calculation: "fixed",
				Month:       1,
				Day:         1,
			},
		},
		UpdatedAt: time.Now(),
	}, nil
}

// ValidatePythonContent validates Python source content (mock implementation)
func (m *MockSyncer) ValidatePythonContent(content string) error {
	if m.shouldError {
		return fmt.Errorf("mock validation error: %s", m.errorMessage)
	}

	// Simple validation - check for basic Python class structure
	if !strings.Contains(content, "class") {
		return fmt.Errorf("invalid Python content: no class definition found")
	}

	return nil
}

// getCountryName returns a human-readable country name
func getCountryName(code string) string {
	names := map[string]string{
		"US": "United States",
		"GB": "United Kingdom",
		"CA": "Canada",
		"AU": "Australia",
		"DE": "Germany",
		"FR": "France",
	}
	if name, exists := names[code]; exists {
		return name
	}
	return code
}

// Mock Python source files for testing
const mockUSPythonSource = `
from datetime import date
from dateutil.easter import easter

from holidays.calendars import _islamic_to_gre
from holidays.constants import JAN, FEB, MAR, APR, MAY, JUN, JUL, AUG, SEP, OCT, NOV, DEC
from holidays.holiday_base import HolidayBase


class US(HolidayBase):
    country = "US"
    subdivisions = ["AL", "AK", "AS", "AZ", "AR", "CA", "CO", "CT", "DE", "DC", "FL", "GA"]

    def __init__(self, **kwargs):
        super().__init__(**kwargs)

    def _populate(self, year):
        # New Year's Day
        self._add_holiday("New Year's Day", date(year, JAN, 1))
        
        # Independence Day
        self._add_holiday("Independence Day", date(year, JUL, 4))
        
        # Christmas Day
        self._add_holiday("Christmas Day", date(year, DEC, 25))
`

const mockGBPythonSource = `
from datetime import date
from dateutil.easter import easter

from holidays.constants import JAN, FEB, MAR, APR, MAY, JUN, JUL, AUG, SEP, OCT, NOV, DEC
from holidays.holiday_base import HolidayBase


class GB(HolidayBase):
    country = "GB"
    subdivisions = ["ENG", "NIR", "SCT", "WLS"]

    def __init__(self, **kwargs):
        super().__init__(**kwargs)

    def _populate(self, year):
        # New Year's Day
        self._add_holiday("New Year's Day", date(year, JAN, 1))
        
        # Christmas Day
        self._add_holiday("Christmas Day", date(year, DEC, 25))
`

const mockCAPythonSource = `
from datetime import date
from dateutil.easter import easter

from holidays.constants import JAN, FEB, MAR, APR, MAY, JUN, JUL, AUG, SEP, OCT, NOV, DEC
from holidays.holiday_base import HolidayBase


class CA(HolidayBase):
    country = "CA"
    subdivisions = ["AB", "BC", "MB", "NB", "NL", "NS", "NT", "NU", "ON", "PE", "QC", "SK", "YT"]

    def __init__(self, **kwargs):
        super().__init__(**kwargs)

    def _populate(self, year):
        # New Year's Day
        self._add_holiday("New Year's Day", date(year, JAN, 1))
        
        # Canada Day
        self._add_holiday("Canada Day", date(year, JUL, 1))
        
        # Christmas Day
        self._add_holiday("Christmas Day", date(year, DEC, 25))
`
