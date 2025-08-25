package config

import (
	"os"
	"testing"
	"time"

	"github.com/coredds/GoHoliday/countries"
)

func TestBasicConfigurationLoading(t *testing.T) {
	// Create a minimal test config
	configContent := `
countries:
  US:
    enabled: true
    overrides:
      "Martin Luther King Jr. Day": "MLK Day"
      
custom_holidays:
  US:
    - name: "Test Holiday"
      date: "06-15"
      category: "test"
`

	// Write temporary config
	tmpFile, err := os.CreateTemp("", "goholidays_basic_test_*.yaml")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tmpFile.Name())

	if _, err := tmpFile.WriteString(configContent); err != nil {
		t.Fatalf("Failed to write config: %v", err)
	}
	tmpFile.Close()

	// Create a new config manager that loads our test file
	cm := NewConfigManager()
	config, err := cm.LoadConfigFromFile(tmpFile.Name())
	if err != nil {
		t.Fatalf("Failed to load config: %v", err)
	}

	// Check that config was loaded
	if config == nil {
		t.Fatal("Config should not be nil")
	}

	// Check US country config
	usConfig := config.Countries["US"]
	if !usConfig.Enabled {
		t.Error("US should be enabled in test config")
	}

	// Check overrides
	if override, exists := usConfig.Overrides["Martin Luther King Jr. Day"]; !exists || override != "MLK Day" {
		t.Error("MLK Day override should be configured")
	}

	// Check custom holidays
	customHolidays, exists := config.CustomHolidays["US"]
	if !exists || len(customHolidays) == 0 {
		t.Error("US should have custom holidays")
	} else {
		found := false
		for _, holiday := range customHolidays {
			if holiday.Name == "Test Holiday" {
				found = true
				break
			}
		}
		if !found {
			t.Error("Test Holiday should be in custom holidays")
		}
	}
}

func TestHolidayManagerWithMockProvider(t *testing.T) {
	// Create a test config
	configContent := `
countries:
  TEST:
    enabled: true
    overrides:
      "Old Holiday": "New Holiday"
    excluded_holidays: ["Excluded Holiday"]
      
custom_holidays:
  TEST:
    - name: "Custom Holiday"
      date: "07-04"
      category: "custom"
`

	tmpFile, err := os.CreateTemp("", "goholidays_mock_test_*.yaml")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tmpFile.Name())

	if _, err := tmpFile.WriteString(configContent); err != nil {
		t.Fatalf("Failed to write config: %v", err)
	}
	tmpFile.Close()

	// Create config manager and load test config
	cm := NewConfigManager()
	config, err := cm.LoadConfigFromFile(tmpFile.Name())
	if err != nil {
		t.Fatalf("Failed to load config: %v", err)
	}

	// Set the loaded config
	cm.config = config

	// Create holiday manager with our test config
	hm := &HolidayManager{
		configManager: cm,
		providers:     make(map[string]countries.HolidayProvider),
	}

	// Create mock holidays to test configuration application
	testHolidays := make(map[time.Time]*countries.Holiday)

	date1 := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	date2 := time.Date(2024, 2, 1, 0, 0, 0, 0, time.UTC)
	date3 := time.Date(2024, 3, 1, 0, 0, 0, 0, time.UTC)

	testHolidays[date1] = &countries.Holiday{
		Name:     "Old Holiday",
		Date:     date1,
		Category: "test",
	}
	testHolidays[date2] = &countries.Holiday{
		Name:     "Keep Holiday",
		Date:     date2,
		Category: "test",
	}
	testHolidays[date3] = &countries.Holiday{
		Name:     "Excluded Holiday",
		Date:     date3,
		Category: "test",
	}

	// Apply configuration
	result := hm.applyCountryConfig(testHolidays, "TEST", config)

	// Check override was applied
	foundNewName := false
	foundOldName := false
	for _, holiday := range result {
		if holiday.Name == "New Holiday" {
			foundNewName = true
		}
		if holiday.Name == "Old Holiday" {
			foundOldName = true
		}
	}

	if !foundNewName {
		t.Error("Holiday name override should be applied")
	}
	if foundOldName {
		t.Error("Old holiday name should be replaced")
	}

	// Check exclusion was applied
	foundExcluded := false
	for _, holiday := range result {
		if holiday.Name == "Excluded Holiday" {
			foundExcluded = true
		}
	}
	if foundExcluded {
		t.Error("Excluded holiday should be removed")
	}

	// Check custom holidays
	customHolidays := hm.getCustomHolidays("TEST", 2024, config)
	if len(customHolidays) == 0 {
		t.Error("Should have custom holidays")
	}

	foundCustom := false
	for _, holiday := range customHolidays {
		if holiday.Name == "Custom Holiday" {
			foundCustom = true
		}
	}
	if !foundCustom {
		t.Error("Custom holiday should be added")
	}
}

func TestConfigurationPrecedence(t *testing.T) {
	// Test that environment variables override file settings
	configContent := `
countries:
  US:
    enabled: false
    
logging:
  level: "info"
`

	tmpFile, err := os.CreateTemp("", "goholidays_precedence_test_*.yaml")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tmpFile.Name())

	if _, err := tmpFile.WriteString(configContent); err != nil {
		t.Fatalf("Failed to write config: %v", err)
	}
	tmpFile.Close()

	// Set environment variables to override file settings
	originalUSEnabled := os.Getenv("GOHOLIDAYS_COUNTRIES_US_ENABLED")
	originalLogLevel := os.Getenv("GOHOLIDAYS_LOGGING_LEVEL")

	os.Setenv("GOHOLIDAYS_COUNTRIES_US_ENABLED", "true")
	os.Setenv("GOHOLIDAYS_LOGGING_LEVEL", "debug")

	defer func() {
		if originalUSEnabled != "" {
			os.Setenv("GOHOLIDAYS_COUNTRIES_US_ENABLED", originalUSEnabled)
		} else {
			os.Unsetenv("GOHOLIDAYS_COUNTRIES_US_ENABLED")
		}
		if originalLogLevel != "" {
			os.Setenv("GOHOLIDAYS_LOGGING_LEVEL", originalLogLevel)
		} else {
			os.Unsetenv("GOHOLIDAYS_LOGGING_LEVEL")
		}
	}()

	// Load config
	cm := NewConfigManager()
	config, err := cm.LoadConfigFromFile(tmpFile.Name())
	if err != nil {
		t.Fatalf("Failed to load config: %v", err)
	}

	// Check that environment variables took precedence
	if !config.Countries["US"].Enabled {
		t.Error("Environment variable should override file setting for US enabled")
	}

	if config.Logging.Level != "debug" {
		t.Error("Environment variable should override file setting for log level")
	}
}
