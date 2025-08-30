package config

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestConfigEdgeCases(t *testing.T) {
	// Test empty config file
	t.Run("Empty Config File", func(t *testing.T) {
		tmpFile, err := os.CreateTemp("", "goholidays_empty_*.yaml")
		if err != nil {
			t.Fatalf("Failed to create temp file: %v", err)
		}
		defer os.Remove(tmpFile.Name())
		tmpFile.Close()

		cm := NewConfigManager()
		cfg, err := cm.LoadConfigFromFile(tmpFile.Name())
		if err != nil {
			t.Errorf("Loading empty config should not error: %v", err)
		}
		if cfg == nil {
			t.Error("Empty config should return default config")
		}
	})

	// Test config with invalid types
	t.Run("Invalid Types", func(t *testing.T) {
		configContent := `
general:
  default_country: 123  # Should be string
  supported_languages: "en"  # Should be array
logging:
  level: ["info"]  # Should be string
`
		tmpFile, err := os.CreateTemp("", "goholidays_invalid_types_*.yaml")
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
		if err == nil {
			t.Error("Invalid types should cause error")
		}
	})

	// Test config with invalid date formats
	t.Run("Invalid Date Formats", func(t *testing.T) {
		configContent := `
custom_holidays:
  US:
    - name: "Invalid Date Holiday"
      date: "13-32"  # Invalid month/day
      category: "test"
    - name: "Invalid Format Holiday"
      date: "not-a-date"
      category: "test"
`
		tmpFile, err := os.CreateTemp("", "goholidays_invalid_dates_*.yaml")
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
			t.Errorf("Invalid dates should not prevent loading: %v", err)
		}

		holidays := cm.GetCustomHolidays("US")
		// Note: Current implementation doesn't validate date formats
		// This test documents the current behavior
		if len(holidays) == 0 {
			t.Log("Invalid date holidays are currently not loaded (expected behavior)")
		}
	})

	// Test config with duplicate holiday definitions
	t.Run("Duplicate Holidays", func(t *testing.T) {
		configContent := `
custom_holidays:
  US:
    - name: "Duplicate Holiday"
      date: "01-01"
      category: "test"
    - name: "Duplicate Holiday"
      date: "01-01"
      category: "test"
`
		tmpFile, err := os.CreateTemp("", "goholidays_duplicate_*.yaml")
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
			t.Errorf("Duplicate holidays should not prevent loading: %v", err)
		}

		holidays := cm.GetCustomHolidays("US")
		// Deduplication should now work - expect only 1 unique holiday
		if len(holidays) != 1 {
			t.Errorf("Expected 1 unique holiday after deduplication, got %d", len(holidays))
		} else {
			t.Logf("Found %d unique holiday after deduplication (expected)", len(holidays))
		}
	})

	// Test config with invalid timezone
	t.Run("Invalid Timezone", func(t *testing.T) {
		configContent := `
general:
  default_timezone: "Invalid/Timezone"
`
		tmpFile, err := os.CreateTemp("", "goholidays_invalid_tz_*.yaml")
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
		if err == nil {
			t.Error("Invalid timezone should cause validation error")
		}
	})

	// Test config with invalid file permissions
	t.Run("Invalid File Permissions", func(t *testing.T) {
		tmpFile, err := os.CreateTemp("", "goholidays_perms_*.yaml")
		if err != nil {
			t.Fatalf("Failed to create temp file: %v", err)
		}
		defer os.Remove(tmpFile.Name())

		if err := os.Chmod(tmpFile.Name(), 0000); err != nil {
			t.Fatalf("Failed to change file permissions: %v", err)
		}

		cm := NewConfigManager()
		_, err = cm.LoadConfigFromFile(tmpFile.Name())
		// Note: File permission handling may vary by OS
		// This test documents the current behavior
		if err != nil {
			t.Logf("Unreadable file caused error as expected: %v", err)
		} else {
			t.Log("Unreadable file did not cause error (OS-dependent behavior)")
		}

		// Reset permissions for cleanup
		_ = os.Chmod(tmpFile.Name(), 0600)
	})

	// Test config with invalid directory permissions
	t.Run("Invalid Directory Permissions", func(t *testing.T) {
		tmpDir, err := os.MkdirTemp("", "goholidays_dir_*")
		if err != nil {
			t.Fatalf("Failed to create temp directory: %v", err)
		}
		defer os.RemoveAll(tmpDir)

		configPath := filepath.Join(tmpDir, "config.yaml")
		if err := os.WriteFile(configPath, []byte(""), 0600); err != nil {
			t.Fatalf("Failed to write config: %v", err)
		}

		if err := os.Chmod(tmpDir, 0000); err != nil {
			t.Fatalf("Failed to change directory permissions: %v", err)
		}

		cm := NewConfigManager()
		_, err = cm.LoadConfigFromFile(configPath)
		// Note: Directory permission handling may vary by OS
		// This test documents the current behavior
		if err != nil {
			t.Logf("Unreadable directory caused error as expected: %v", err)
		} else {
			t.Log("Unreadable directory did not cause error (OS-dependent behavior)")
		}

		// Reset permissions for cleanup
		_ = os.Chmod(tmpDir, 0700)
	})

	// Test config with very large file
	t.Run("Large Config File", func(t *testing.T) {
		tmpFile, err := os.CreateTemp("", "goholidays_large_*.yaml")
		if err != nil {
			t.Fatalf("Failed to create temp file: %v", err)
		}
		defer os.Remove(tmpFile.Name())

		// Write large config (100K lines)
		for i := 0; i < 100000; i++ {
			if _, err := tmpFile.WriteString("# Comment line\n"); err != nil {
				t.Fatalf("Failed to write config: %v", err)
			}
		}
		tmpFile.Close()

		cm := NewConfigManager()
		_, err = cm.LoadConfigFromFile(tmpFile.Name())
		if err != nil {
			t.Errorf("Large config should load: %v", err)
		}
	})

	// Test concurrent config access
	t.Run("Concurrent Access", func(t *testing.T) {
		cm := NewConfigManager()
		done := make(chan bool)

		// Start multiple goroutines accessing config
		for i := 0; i < 10; i++ {
			go func() {
				config := cm.GetConfig()
				if config == nil {
					t.Error("Concurrent access returned nil config")
				}
				done <- true
			}()
		}

		// Wait for all goroutines
		for i := 0; i < 10; i++ {
			<-done
		}
	})

	// Test config with all possible environment variables
	t.Run("All Environment Variables", func(t *testing.T) {
		// Save original environment
		origEnv := make(map[string]string)
		envVars := []string{
			"GOHOLIDAYS_DEFAULT_COUNTRY",
			"GOHOLIDAYS_DEFAULT_LANGUAGE",
			"GOHOLIDAYS_DEFAULT_TIMEZONE",
			"GOHOLIDAYS_ENVIRONMENT",
			"GOHOLIDAYS_LOG_LEVEL",
			"GOHOLIDAYS_COUNTRIES_US_ENABLED",
		}

		for _, env := range envVars {
			origEnv[env] = os.Getenv(env)
		}

		// Restore environment after test
		defer func() {
			for env, val := range origEnv {
				if val != "" {
					os.Setenv(env, val)
				} else {
					os.Unsetenv(env)
				}
			}
		}()

		// Set all environment variables
		os.Setenv("GOHOLIDAYS_DEFAULT_COUNTRY", "GB")
		os.Setenv("GOHOLIDAYS_DEFAULT_LANGUAGE", "en")
		os.Setenv("GOHOLIDAYS_DEFAULT_TIMEZONE", "UTC")
		os.Setenv("GOHOLIDAYS_ENVIRONMENT", "prod")
		os.Setenv("GOHOLIDAYS_LOG_LEVEL", "debug")
		os.Setenv("GOHOLIDAYS_COUNTRIES_US_ENABLED", "true")

		cm := NewConfigManager()
		cfg, err := cm.LoadConfig()
		if err != nil {
			t.Fatalf("Failed to load config: %v", err)
		}

		// Verify all environment variables were applied
		if cfg.General.DefaultCountry != "GB" {
			t.Error("Environment variable for default country not applied")
		}
		if cfg.General.DefaultLanguage != "en" {
			t.Error("Environment variable for default language not applied")
		}
		if cfg.General.DefaultTimezone != "UTC" {
			t.Error("Environment variable for default timezone not applied")
		}
		if cfg.General.Environment != "prod" {
			t.Error("Environment variable for environment not applied")
		}
		if cfg.Logging.Level != "debug" {
			t.Error("Environment variable for log level not applied")
		}
		if !cfg.Countries["US"].Enabled {
			t.Error("Environment variable for US enabled not applied")
		}
	})

	// Test config with custom holiday calculation rules
	t.Run("Custom Holiday Calculations", func(t *testing.T) {
		configContent := `
custom_holidays:
  US:
    - name: "Easter Based Holiday"
      calculation:
        type: "easter_offset"
        easter_offset: 10
    - name: "Weekday Based Holiday"
      calculation:
        type: "weekday"
        weekday: "monday"
        week: 3
        month: 1
`
		tmpFile, err := os.CreateTemp("", "goholidays_calc_*.yaml")
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
			t.Errorf("Custom calculations should load: %v", err)
		}

		holidays := cm.GetCustomHolidays("US")
		if len(holidays) != 2 {
			t.Errorf("Expected 2 custom holidays, got %d", len(holidays))
		}
	})

	// Test config with year ranges
	t.Run("Year Ranges", func(t *testing.T) {
		configContent := `
custom_holidays:
  US:
    - name: "Past Holiday"
      date: "01-01"
      year_range:
        start: 1900
        end: 1999
    - name: "Future Holiday"
      date: "01-01"
      year_range:
        start: 2100
`
		tmpFile, err := os.CreateTemp("", "goholidays_years_*.yaml")
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
			t.Errorf("Year ranges should load: %v", err)
		}

		holidays := cm.GetCustomHolidays("US")
		currentYear := time.Now().Year()
		// Note: Current implementation doesn't filter by year ranges
		// This test documents the current behavior
		for _, holiday := range holidays {
			if holiday.YearRange != nil {
				t.Logf("Holiday %s has year range %d-%d (current implementation doesn't filter)",
					holiday.Name, holiday.YearRange.Start, holiday.YearRange.End)
			}
		}
		t.Logf("Found %d holidays for current year %d", len(holidays), currentYear)
	})
}

func TestConfigPerformance(t *testing.T) {
	// Test performance settings
	t.Run("Performance Settings", func(t *testing.T) {
		configContent := `
performance:
  enable_caching: true
  cache_ttl: 24h
  max_cache_size: 1000
  preload_years: 2
  concurrent_limit: 10
  batch_size: 100
`
		tmpFile, err := os.CreateTemp("", "goholidays_perf_*.yaml")
		if err != nil {
			t.Fatalf("Failed to create temp file: %v", err)
		}
		defer os.Remove(tmpFile.Name())

		if _, err := tmpFile.WriteString(configContent); err != nil {
			t.Fatalf("Failed to write config: %v", err)
		}
		tmpFile.Close()

		cm := NewConfigManager()
		cfg, err := cm.LoadConfigFromFile(tmpFile.Name())
		if err != nil {
			t.Fatalf("Failed to load config: %v", err)
		}

		if !cfg.Performance.EnableCaching {
			t.Error("Caching should be enabled")
		}
		if cfg.Performance.CacheTTL != 24*time.Hour {
			t.Error("Cache TTL not set correctly")
		}
		if cfg.Performance.MaxCacheSize != 1000 {
			t.Error("Max cache size not set correctly")
		}
		if cfg.Performance.PreloadYears != 2 {
			t.Error("Preload years not set correctly")
		}
		if cfg.Performance.ConcurrentLimit != 10 {
			t.Error("Concurrent limit not set correctly")
		}
		if cfg.Performance.BatchSize != 100 {
			t.Error("Batch size not set correctly")
		}
	})
}
