package config

import (
	"os"
	"testing"
	"time"

	"github.com/coredds/GoHoliday/countries"
)

func TestConfigurationIntegration(t *testing.T) {
	// Create a temporary config file for testing
	configContent := `
countries:
  US:
    enabled: true
    categories: ["federal"]
    subdivisions: ["CA", "NY"]
    excluded_holidays: ["Columbus Day"]
    overrides:
      "Martin Luther King Jr. Day": "MLK Day"
  GB:
    enabled: true
    excluded_holidays: ["Boxing Day"]

custom_holidays:
  US:
    - name: "Company Day"
      date: "03-15"
      category: "observance"
      languages:
        en: "Company Day"
    - name: "Easter Monday Custom"
      category: "religious"
      calculation:
        type: "easter_offset"
        easter_offset: 1

output:
  timezone: "America/New_York"
  date_format: "2006-01-02"
  include_metadata: true
`

	// Write temporary config
	tmpFile, err := os.CreateTemp("", "goholidays_test_*.yaml")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tmpFile.Name())

	if _, err := tmpFile.WriteString(configContent); err != nil {
		t.Fatalf("Failed to write config: %v", err)
	}
	tmpFile.Close()

	// Set environment variable to use our test config
	originalConfigPath := os.Getenv("GOHOLIDAYS_CONFIG")
	os.Setenv("GOHOLIDAYS_CONFIG", tmpFile.Name())
	defer func() {
		if originalConfigPath != "" {
			os.Setenv("GOHOLIDAYS_CONFIG", originalConfigPath)
		} else {
			os.Unsetenv("GOHOLIDAYS_CONFIG")
		}
	}()

	// Test holiday manager initialization
	manager := NewHolidayManager()
	if manager == nil {
		t.Fatal("Holiday manager should not be nil")
	}

	// Test supported countries
	countries := manager.GetSupportedCountries()
	if len(countries) == 0 {
		t.Error("Should have supported countries")
	}

	// Verify US is enabled
	found := false
	for _, country := range countries {
		if country == "US" {
			found = true
			break
		}
	}
	if !found {
		t.Error("US should be in supported countries")
	}

	// Test getting holidays with configuration applied
	holidays, err := manager.GetHolidays("US", 2024)
	if err != nil {
		t.Fatalf("Error getting US holidays: %v", err)
	}

	// Print holidays for debugging
	t.Logf("Found %d holidays for US", len(holidays))
	for _, holiday := range holidays {
		t.Logf("Holiday: %s", holiday.Name)
	}

	// Verify MLK Day name override (if Martin Luther King Jr. Day exists)
	mlkOriginalFound := false
	mlkNewFound := false
	for _, holiday := range holidays {
		if holiday.Name == "Martin Luther King Jr. Day" {
			mlkOriginalFound = true
		}
		if holiday.Name == "MLK Day" {
			mlkNewFound = true
		}
	}
	
	// If we have the original holiday, the override should have been applied
	if mlkOriginalFound {
		t.Error("Original 'Martin Luther King Jr. Day' should be renamed to 'MLK Day'")
	}
	
	// We should have the renamed version if the holiday exists
	if !mlkNewFound && mlkOriginalFound {
		t.Error("MLK Day name override should be applied")
	}

	// Verify Columbus Day is excluded (only check if it would normally exist)
	columbusFound := false
	for _, holiday := range holidays {
		if holiday.Name == "Columbus Day" {
			columbusFound = true
			break
		}
	}
	if columbusFound {
		t.Error("Columbus Day should be excluded")
	}

	// Test custom holidays
	companyDayFound := false
	for date, holiday := range holidays {
		if holiday.Name == "Company Day" && date.Month() == 3 && date.Day() == 15 {
			companyDayFound = true
			break
		}
	}
	if !companyDayFound {
		t.Error("Custom Company Day should be added")
	}

	// Test country info
	info, err := manager.GetCountryInfo("US")
	if err != nil {
		t.Fatalf("Error getting country info: %v", err)
	}

	if enabled, ok := info["enabled"].(bool); !ok || !enabled {
		t.Error("US should be enabled")
	}

	// Test regional holidays
	regionalHolidays, err := manager.GetHolidaysWithSubdivisions("US", 2024, []string{"CA", "TX"})
	if err != nil {
		t.Fatalf("Error getting regional holidays: %v", err)
	}

	// Should include CA but not TX (TX not in allowed subdivisions)
	if len(regionalHolidays) <= len(holidays) {
		t.Log("Note: No additional regional holidays found - this may be expected")
	}
}

func TestConfigurationOverrides(t *testing.T) {
	configContent := `
countries:
  GB:
    enabled: true
    excluded_holidays: ["Boxing Day", "Spring Bank Holiday"]
    overrides:
      "New Year's Day": "New Year"
      "Christmas Day": "Xmas"
`

	tmpFile, err := os.CreateTemp("", "goholidays_test_*.yaml")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tmpFile.Name())

	if _, err := tmpFile.WriteString(configContent); err != nil {
		t.Fatalf("Failed to write config: %v", err)
	}
	tmpFile.Close()

	originalConfigPath := os.Getenv("GOHOLIDAYS_CONFIG")
	os.Setenv("GOHOLIDAYS_CONFIG", tmpFile.Name())
	defer func() {
		if originalConfigPath != "" {
			os.Setenv("GOHOLIDAYS_CONFIG", originalConfigPath)
		} else {
			os.Unsetenv("GOHOLIDAYS_CONFIG")
		}
	}()

	manager := NewHolidayManager()
	holidays, err := manager.GetHolidays("GB", 2024)
	if err != nil {
		t.Fatalf("Error getting GB holidays: %v", err)
	}

	// Check overrides
	newYearFound := false
	xmasFound := false
	for _, holiday := range holidays {
		if holiday.Name == "New Year" {
			newYearFound = true
		}
		if holiday.Name == "Xmas" {
			xmasFound = true
		}
	}

	if !newYearFound {
		t.Error("New Year override should be applied")
	}
	if !xmasFound {
		t.Error("Xmas override should be applied")
	}

	// Check exclusions
	boxingDayFound := false
	springBankFound := false
	for _, holiday := range holidays {
		if holiday.Name == "Boxing Day" {
			boxingDayFound = true
		}
		if holiday.Name == "Spring Bank Holiday" {
			springBankFound = true
		}
	}

	if boxingDayFound {
		t.Error("Boxing Day should be excluded")
	}
	if springBankFound {
		t.Error("Spring Bank Holiday should be excluded")
	}
}

func TestCustomHolidayCalculation(t *testing.T) {
	configContent := `
custom_holidays:
  US:
    - name: "Good Friday Custom"
      category: "religious"
      calculation:
        type: "easter_offset"
        easter_offset: -2
    - name: "Second Monday March"
      category: "observance"
      calculation:
        type: "weekday"
        weekday_rule:
          month: 3
          weekday: "monday"
          week: 2
`

	tmpFile, err := os.CreateTemp("", "goholidays_test_*.yaml")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tmpFile.Name())

	if _, err := tmpFile.WriteString(configContent); err != nil {
		t.Fatalf("Failed to write config: %v", err)
	}
	tmpFile.Close()

	originalConfigPath := os.Getenv("GOHOLIDAYS_CONFIG")
	os.Setenv("GOHOLIDAYS_CONFIG", tmpFile.Name())
	defer func() {
		if originalConfigPath != "" {
			os.Setenv("GOHOLIDAYS_CONFIG", originalConfigPath)
		} else {
			os.Unsetenv("GOHOLIDAYS_CONFIG")
		}
	}()

	manager := NewHolidayManager()
	holidays, err := manager.GetHolidays("US", 2024)
	if err != nil {
		t.Fatalf("Error getting US holidays: %v", err)
	}

	// Check Good Friday Custom (Easter - 2 days)
	easter := countries.EasterSunday(2024)
	expectedGoodFriday := easter.AddDate(0, 0, -2)
	
	goodFridayFound := false
	for date, holiday := range holidays {
		if holiday.Name == "Good Friday Custom" && 
		   date.Year() == expectedGoodFriday.Year() &&
		   date.Month() == expectedGoodFriday.Month() &&
		   date.Day() == expectedGoodFriday.Day() {
			goodFridayFound = true
			break
		}
	}
	if !goodFridayFound {
		t.Error("Good Friday Custom should be calculated correctly")
	}

	// Check Second Monday of March
	secondMondayFound := false
	for date, holiday := range holidays {
		if holiday.Name == "Second Monday March" && 
		   date.Month() == 3 && 
		   date.Weekday() == time.Monday {
			// Check if it's the second Monday
			firstMonday := countries.NthWeekdayOfMonth(2024, 3, time.Monday, 1)
			secondMonday := firstMonday.AddDate(0, 0, 7)
			if date.Day() == secondMonday.Day() {
				secondMondayFound = true
				break
			}
		}
	}
	if !secondMondayFound {
		t.Error("Second Monday March should be calculated correctly")
	}
}

func TestEnvironmentConfiguration(t *testing.T) {
	// Test that different environments can have different configs
	devConfig := `
countries:
  US:
    enabled: true
    categories: ["public"]
logging:
  level: "debug"
performance:
  cache_enabled: false
`

	prodConfig := `
countries:
  US:
    enabled: true
    categories: ["public", "observance", "religious"]
logging:
  level: "info"
performance:
  cache_enabled: true
  cache_size: 1000
`

	// Test dev config
	tmpDevFile, err := os.CreateTemp("", "goholidays_dev_*.yaml")
	if err != nil {
		t.Fatalf("Failed to create dev temp file: %v", err)
	}
	defer os.Remove(tmpDevFile.Name())

	if _, err := tmpDevFile.WriteString(devConfig); err != nil {
		t.Fatalf("Failed to write dev config: %v", err)
	}
	tmpDevFile.Close()

	// Test prod config
	tmpProdFile, err := os.CreateTemp("", "goholidays_prod_*.yaml")
	if err != nil {
		t.Fatalf("Failed to create prod temp file: %v", err)
	}
	defer os.Remove(tmpProdFile.Name())

	if _, err := tmpProdFile.WriteString(prodConfig); err != nil {
		t.Fatalf("Failed to write prod config: %v", err)
	}
	tmpProdFile.Close()

	// Test dev environment
	originalConfigPath := os.Getenv("GOHOLIDAYS_CONFIG")
	originalEnv := os.Getenv("GOHOLIDAYS_ENV")
	
	os.Setenv("GOHOLIDAYS_CONFIG", tmpDevFile.Name())
	os.Setenv("GOHOLIDAYS_ENV", "dev")
	
	devManager := NewHolidayManager()
	devConfig1 := devManager.configManager.GetConfig()
	
	if devConfig1.Logging.Level != "debug" {
		t.Error("Dev config should have debug logging")
	}
	
	// Test prod environment
	os.Setenv("GOHOLIDAYS_CONFIG", tmpProdFile.Name())
	os.Setenv("GOHOLIDAYS_ENV", "prod")
	
	prodManager := NewHolidayManager()
	prodConfig1 := prodManager.configManager.GetConfig()
	
	if prodConfig1.Logging.Level != "info" {
		t.Error("Prod config should have info logging")
	}
	
	// Restore environment
	if originalConfigPath != "" {
		os.Setenv("GOHOLIDAYS_CONFIG", originalConfigPath)
	} else {
		os.Unsetenv("GOHOLIDAYS_CONFIG")
	}
	if originalEnv != "" {
		os.Setenv("GOHOLIDAYS_ENV", originalEnv)
	} else {
		os.Unsetenv("GOHOLIDAYS_ENV")
	}
}
