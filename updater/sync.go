// Package updater provides functionality to sync holiday data from the Python holidays repository
package updater

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// PythonHolidaysSync handles synchronization with the Python holidays repository
type PythonHolidaysSync struct {
	repoURL    string
	dataDir    string
	httpClient *http.Client
}

// NewPythonHolidaysSync creates a new sync instance
func NewPythonHolidaysSync(dataDir string) *PythonHolidaysSync {
	return &PythonHolidaysSync{
		repoURL:    "https://api.github.com/repos/vacanza/holidays",
		dataDir:    dataDir,
		httpClient: &http.Client{Timeout: 30 * time.Second},
	}
}

// CountryData represents holiday data for a country
type CountryData struct {
	CountryCode  string                       `json:"country_code"`
	Name         string                       `json:"name"`
	Subdivisions map[string]string            `json:"subdivisions,omitempty"`
	Categories   []string                     `json:"categories"`
	Languages    []string                     `json:"languages"`
	Holidays     map[string]HolidayDefinition `json:"holidays"`
	UpdatedAt    time.Time                    `json:"updated_at"`
}

// HolidayDefinition represents a holiday definition from Python source
type HolidayDefinition struct {
	Name         string            `json:"name"`
	Category     string            `json:"category"`
	Languages    map[string]string `json:"languages"`
	Calculation  string            `json:"calculation"` // "fixed", "easter_based", "weekday_based"
	Month        int               `json:"month,omitempty"`
	Day          int               `json:"day,omitempty"`
	EasterOffset int               `json:"easter_offset,omitempty"`
	WeekdayRule  *WeekdayRule      `json:"weekday_rule,omitempty"`
	YearRange    *YearRange        `json:"year_range,omitempty"`
	Subdivisions []string          `json:"subdivisions,omitempty"`
}

// WeekdayRule defines rules for weekday-based holidays
type WeekdayRule struct {
	Month      int          `json:"month"`
	Weekday    time.Weekday `json:"weekday"`
	Occurrence int          `json:"occurrence"` // 1=first, 2=second, -1=last
}

// YearRange defines the valid year range for a holiday
type YearRange struct {
	Start int `json:"start,omitempty"`
	End   int `json:"end,omitempty"`
}

// SyncData synchronizes holiday data from the Python holidays repository
func (phs *PythonHolidaysSync) SyncData() error {
	fmt.Println("Starting sync with Python holidays repository...")

	// Create data directory if it doesn't exist
	if err := os.MkdirAll(phs.dataDir, 0755); err != nil {
		return fmt.Errorf("failed to create data directory: %w", err)
	}

	// Get list of supported countries from the repository
	countries, err := phs.getSupportedCountries()
	if err != nil {
		return fmt.Errorf("failed to get supported countries: %w", err)
	}

	fmt.Printf("Found %d countries to sync\n", len(countries))

	// Sync each country
	for _, country := range countries {
		fmt.Printf("Syncing %s (%s)...\n", country.Name, country.CountryCode)

		countryData, err := phs.fetchCountryData(country.CountryCode)
		if err != nil {
			fmt.Printf("Warning: failed to sync %s: %v\n", country.CountryCode, err)
			continue
		}

		if err := phs.saveCountryData(countryData); err != nil {
			fmt.Printf("Warning: failed to save %s: %v\n", country.CountryCode, err)
			continue
		}
	}

	fmt.Println("Sync completed!")
	return nil
}

// getSupportedCountries retrieves the list of supported countries
func (phs *PythonHolidaysSync) getSupportedCountries() ([]CountryInfo, error) {
	// This is a placeholder implementation
	// In a real implementation, this would parse the Python source code or
	// use a metadata file to get the list of supported countries

	countries := []CountryInfo{
		{CountryCode: "US", Name: "United States"},
		{CountryCode: "GB", Name: "United Kingdom"},
		{CountryCode: "CA", Name: "Canada"},
		{CountryCode: "AU", Name: "Australia"},
		{CountryCode: "DE", Name: "Germany"},
		{CountryCode: "FR", Name: "France"},
	}

	return countries, nil
}

// CountryInfo holds basic country information
type CountryInfo struct {
	CountryCode string `json:"country_code"`
	Name        string `json:"name"`
}

// fetchCountryData fetches holiday data for a specific country
func (phs *PythonHolidaysSync) fetchCountryData(countryCode string) (*CountryData, error) {
	// This is a placeholder implementation
	// In a real implementation, this would:
	// 1. Fetch the Python source file for the country
	// 2. Parse the Python AST to extract holiday definitions
	// 3. Convert the definitions to our JSON format

	countryData := &CountryData{
		CountryCode: countryCode,
		Name:        phs.getCountryName(countryCode),
		Categories:  []string{"public", "bank", "government"},
		Languages:   []string{"en"},
		Holidays:    make(map[string]HolidayDefinition),
		UpdatedAt:   time.Now(),
	}

	// Add some sample holidays based on country
	switch countryCode {
	case "US":
		countryData.Holidays["new_years_day"] = HolidayDefinition{
			Name:        "New Year's Day",
			Category:    "public",
			Languages:   map[string]string{"en": "New Year's Day", "es": "Año Nuevo"},
			Calculation: "fixed",
			Month:       1,
			Day:         1,
		}
		countryData.Holidays["independence_day"] = HolidayDefinition{
			Name:        "Independence Day",
			Category:    "public",
			Languages:   map[string]string{"en": "Independence Day", "es": "Día de la Independencia"},
			Calculation: "fixed",
			Month:       7,
			Day:         4,
		}
		countryData.Holidays["thanksgiving"] = HolidayDefinition{
			Name:        "Thanksgiving Day",
			Category:    "public",
			Languages:   map[string]string{"en": "Thanksgiving Day", "es": "Día de Acción de Gracias"},
			Calculation: "weekday_based",
			WeekdayRule: &WeekdayRule{
				Month:      11,
				Weekday:    time.Thursday,
				Occurrence: 4,
			},
		}
	}

	return countryData, nil
}

// saveCountryData saves country data to a JSON file
func (phs *PythonHolidaysSync) saveCountryData(data *CountryData) error {
	filename := filepath.Join(phs.dataDir, strings.ToLower(data.CountryCode)+".json")

	file, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("failed to create file %s: %w", filename, err)
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")

	if err := encoder.Encode(data); err != nil {
		return fmt.Errorf("failed to encode data: %w", err)
	}

	return nil
}

// LoadCountryData loads country data from a JSON file
func (phs *PythonHolidaysSync) LoadCountryData(countryCode string) (*CountryData, error) {
	filename := filepath.Join(phs.dataDir, strings.ToLower(countryCode)+".json")

	file, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to open file %s: %w", filename, err)
	}
	defer file.Close()

	var data CountryData
	decoder := json.NewDecoder(file)

	if err := decoder.Decode(&data); err != nil {
		return nil, fmt.Errorf("failed to decode data: %w", err)
	}

	return &data, nil
}

// getCountryName returns the full name for a country code
func (phs *PythonHolidaysSync) getCountryName(countryCode string) string {
	names := map[string]string{
		"US": "United States",
		"GB": "United Kingdom",
		"CA": "Canada",
		"AU": "Australia",
		"DE": "Germany",
		"FR": "France",
	}

	if name, exists := names[countryCode]; exists {
		return name
	}

	return countryCode
}

// CheckForUpdates checks if there are updates available
func (phs *PythonHolidaysSync) CheckForUpdates() (bool, error) {
	// This would check the last commit date of the Python holidays repository
	// and compare it with our last sync time

	req, err := http.NewRequest("GET", phs.repoURL+"/commits", nil)
	if err != nil {
		return false, err
	}

	resp, err := phs.httpClient.Do(req)
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return false, fmt.Errorf("HTTP %d", resp.StatusCode)
	}

	// Parse response and check dates
	// This is a simplified implementation
	return true, nil
}
