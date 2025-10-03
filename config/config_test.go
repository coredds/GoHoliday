package config

import (
	"os"
	"testing"
	"time"

	"github.com/coredds/goholiday/countries"
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

// TestConfigValidation tests configuration validation
func TestConfigValidation(t *testing.T) {
	cm := NewConfigManager()

	// Test valid config
	validConfig := &Config{
		General: GeneralConfig{
			DefaultCountry:  "US",
			DefaultLanguage: "en",
			Environment:     "dev",
		},
		Logging: LoggingConfig{
			Level: "info",
		},
	}

	err := cm.validateConfig(validConfig)
	if err != nil {
		t.Errorf("Valid config should not produce validation error: %v", err)
	}

	// Test invalid environment
	invalidConfig := &Config{
		General: GeneralConfig{
			DefaultCountry:  "US",
			DefaultLanguage: "en",
			Environment:     "invalid",
		},
	}

	err = cm.validateConfig(invalidConfig)
	if err == nil {
		t.Error("Invalid environment should produce validation error")
	}
}

// TestConfigSaving tests saving configuration to file
func TestConfigSaving(t *testing.T) {
	cm := NewConfigManager()

	// Load a config first
	_, err := cm.LoadConfig()
	if err != nil {
		t.Logf("Could not load config: %v, using default", err)
	}

	// Create temporary file
	tmpFile, err := os.CreateTemp("", "goholidays_save_test_*.yaml")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	tmpFile.Close()
	defer os.Remove(tmpFile.Name())

	// Save config
	err = cm.SaveConfig(tmpFile.Name())
	if err != nil {
		t.Errorf("Failed to save config: %v", err)
	}

	// Verify file exists and has content
	if _, err := os.Stat(tmpFile.Name()); os.IsNotExist(err) {
		t.Error("Config file should exist after saving")
	}
}

// TestCountryConfigMethods tests country-specific config methods
func TestCountryConfigMethods(t *testing.T) {
	configContent := `
countries:
  US:
    enabled: true
    subdivisions: ["CA", "NY"]
  GB:
    enabled: false
    
custom_holidays:
  US:
    - name: "Test Holiday"
      date: "06-15"
      category: "test"
`

	tmpFile, err := os.CreateTemp("", "goholidays_country_test_*.yaml")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tmpFile.Name())

	if _, err := tmpFile.WriteString(configContent); err != nil {
		t.Fatalf("Failed to write config: %v", err)
	}
	tmpFile.Close()

	cm := NewConfigManager()
	_, err = cm.LoadConfigFromFile(tmpFile.Name())
	if err != nil {
		t.Fatalf("Failed to load config: %v", err)
	}

	// Test IsCountryEnabled
	if !cm.IsCountryEnabled("US") {
		t.Error("US should be enabled")
	}

	if cm.IsCountryEnabled("GB") {
		t.Error("GB should be disabled")
	}

	// Test undefined country (should be enabled by default according to GetCountryConfig)
	if !cm.IsCountryEnabled("FR") {
		t.Error("FR should be enabled (default behavior)")
	}

	// Test GetCountryConfig
	usConfig := cm.GetCountryConfig("US")
	if !usConfig.Enabled {
		t.Error("US config should be enabled")
	}

	if len(usConfig.Subdivisions) != 2 {
		t.Errorf("Expected 2 subdivisions, got %d", len(usConfig.Subdivisions))
	}

	// Test GetCustomHolidays
	customHolidays := cm.GetCustomHolidays("US")
	if len(customHolidays) != 1 {
		t.Errorf("Expected 1 custom holiday, got %d", len(customHolidays))
	}

	if customHolidays[0].Name != "Test Holiday" {
		t.Errorf("Expected 'Test Holiday', got '%s'", customHolidays[0].Name)
	}
}

// TestErrorHandling tests various error conditions
func TestErrorHandling(t *testing.T) {
	cm := NewConfigManager()

	// Test loading non-existent file
	_, err := cm.LoadConfigFromFile("/non/existent/file.yaml")
	if err == nil {
		t.Error("Loading non-existent file should return error")
	}

	// Test loading invalid YAML
	tmpFile, err := os.CreateTemp("", "goholidays_invalid_test_*.yaml")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tmpFile.Name())

	invalidYAML := "invalid: yaml: content: [unclosed bracket"
	if _, err := tmpFile.WriteString(invalidYAML); err != nil {
		t.Fatalf("Failed to write invalid YAML: %v", err)
	}
	tmpFile.Close()

	_, err = cm.LoadConfigFromFile(tmpFile.Name())
	if err == nil {
		t.Error("Loading invalid YAML should return error")
	}
}

// TestDefaultConfig tests default configuration generation
func TestDefaultConfig(t *testing.T) {
	cm := NewConfigManager()
	config := cm.getDefaultConfig()

	if config == nil {
		t.Fatal("Default config should not be nil")
	}

	if config.General.DefaultLanguage == "" {
		t.Error("Default config should have default language set")
	}

	if config.General.Environment == "" {
		t.Error("Default config should have environment set")
	}
}

// BenchmarkConfigLoading benchmarks configuration loading
func BenchmarkConfigLoading(b *testing.B) {
	configContent := `
countries:
  US:
    enabled: true
  GB:
    enabled: true
`

	tmpFile, err := os.CreateTemp("", "goholidays_bench_test_*.yaml")
	if err != nil {
		b.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tmpFile.Name())

	if _, err := tmpFile.WriteString(configContent); err != nil {
		b.Fatalf("Failed to write config: %v", err)
	}
	tmpFile.Close()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		cm := NewConfigManager()
		_, err := cm.LoadConfigFromFile(tmpFile.Name())
		if err != nil {
			b.Fatalf("Failed to load config: %v", err)
		}
	}
}
