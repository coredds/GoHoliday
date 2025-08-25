package config

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"gopkg.in/yaml.v3"
)

// Config represents the main configuration structure
type Config struct {
	// General settings
	General GeneralConfig `yaml:"general"`
	
	// Country-specific overrides
	Countries map[string]CountryConfig `yaml:"countries"`
	
	// Custom holidays by country
	CustomHolidays map[string][]CustomHoliday `yaml:"custom_holidays"`
	
	// Output formatting
	Output OutputConfig `yaml:"output"`
	
	// Performance settings
	Performance PerformanceConfig `yaml:"performance"`
	
	// Logging configuration
	Logging LoggingConfig `yaml:"logging"`
}

// GeneralConfig contains general application settings
type GeneralConfig struct {
	DefaultCountry   string   `yaml:"default_country"`
	DefaultLanguage  string   `yaml:"default_language"`
	DefaultTimezone  string   `yaml:"default_timezone"`
	SupportedLanguages []string `yaml:"supported_languages"`
	Environment      string   `yaml:"environment"` // dev, staging, prod
}

// CountryConfig allows overriding country-specific settings
type CountryConfig struct {
	Enabled      bool              `yaml:"enabled"`
	Subdivisions []string          `yaml:"subdivisions"`
	Categories   []string          `yaml:"categories"`
	Overrides    map[string]string `yaml:"overrides"` // Holiday name overrides
	ExcludedHolidays []string      `yaml:"excluded_holidays"`
	AdditionalHolidays []string    `yaml:"additional_holidays"`
}

// CustomHoliday allows users to define their own holidays
type CustomHoliday struct {
	Name         string            `yaml:"name"`
	Date         string            `yaml:"date"`         // YYYY-MM-DD or calculation rule
	Countries    []string          `yaml:"countries"`    // Which countries it applies to
	Subdivisions []string          `yaml:"subdivisions"` // Which subdivisions
	Category     string            `yaml:"category"`
	Languages    map[string]string `yaml:"languages"`
	YearRange    *YearRange        `yaml:"year_range,omitempty"`
	Calculation  *CalculationRule  `yaml:"calculation,omitempty"`
}

// YearRange defines when a holiday is valid
type YearRange struct {
	Start int `yaml:"start,omitempty"`
	End   int `yaml:"end,omitempty"`
}

// CalculationRule defines how to calculate variable holidays
type CalculationRule struct {
	Type        string `yaml:"type"`         // "easter_offset", "weekday", "fixed"
	EasterOffset int   `yaml:"easter_offset,omitempty"`
	Month       int    `yaml:"month,omitempty"`
	WeekdayRule *WeekdayRule `yaml:"weekday_rule,omitempty"`
}

// WeekdayRule defines weekday-based holiday calculations
type WeekdayRule struct {
	Weekday   string `yaml:"weekday"`    // "monday", "tuesday", etc.
	Week      int    `yaml:"week"`       // 1=first, -1=last, etc.
	Month     int    `yaml:"month"`
}

// OutputConfig controls how holidays are formatted and returned
type OutputConfig struct {
	DateFormat      string            `yaml:"date_format"`      // "2006-01-02", "January 2, 2006", etc.
	TimeFormat      string            `yaml:"time_format"`      // "15:04:05", "3:04 PM", etc.
	Timezone        string            `yaml:"timezone"`         // "UTC", "Local", "America/New_York"
	IncludeMetadata bool              `yaml:"include_metadata"` // Include category, subdivisions, etc.
	Languages       []string          `yaml:"languages"`        // Preferred language order
	Formats         map[string]string `yaml:"formats"`          // Custom format strings
}

// PerformanceConfig controls performance-related settings
type PerformanceConfig struct {
	EnableCaching    bool          `yaml:"enable_caching"`
	CacheTTL         time.Duration `yaml:"cache_ttl"`
	MaxCacheSize     int           `yaml:"max_cache_size"`
	PreloadYears     int           `yaml:"preload_years"`     // How many years to preload
	ConcurrentLimit  int           `yaml:"concurrent_limit"`  // Max concurrent operations
	BatchSize        int           `yaml:"batch_size"`        // Batch size for bulk operations
}

// LoggingConfig controls logging behavior
type LoggingConfig struct {
	Level      string `yaml:"level"`       // "debug", "info", "warn", "error"
	Format     string `yaml:"format"`      // "json", "text"
	Output     string `yaml:"output"`      // "stdout", "stderr", file path
	EnableFile bool   `yaml:"enable_file"`
	MaxSize    int    `yaml:"max_size"`    // Max log file size in MB
}

// ConfigManager handles configuration loading and management
type ConfigManager struct {
	config *Config
	paths  []string
}

// NewConfigManager creates a new configuration manager
func NewConfigManager() *ConfigManager {
	return &ConfigManager{
		paths: []string{
			"goholidays.yaml",
			"goholidays.yml",
			"config/goholidays.yaml",
			"config/goholidays.yml",
			"/etc/goholidays/config.yaml",
			filepath.Join(os.Getenv("HOME"), ".goholidays.yaml"),
		},
	}
}

// LoadConfig loads configuration from various sources
func (cm *ConfigManager) LoadConfig() (*Config, error) {
	// Start with default configuration
	config := cm.getDefaultConfig()
	
	// Try to load from files
	for _, path := range cm.paths {
		if err := cm.loadFromFile(path, config); err == nil {
			break // Successfully loaded from this file
		}
	}
	
	// Override with environment variables
	cm.loadFromEnvironment(config)
	
	// Validate configuration
	if err := cm.validateConfig(config); err != nil {
		return nil, fmt.Errorf("invalid configuration: %w", err)
	}
	
	cm.config = config
	return config, nil
}

// LoadConfigFromFile loads configuration from a specific file
func (cm *ConfigManager) LoadConfigFromFile(filePath string) (*Config, error) {
	// Start with default configuration
	config := cm.getDefaultConfig()
	
	// Load from the specified file
	if err := cm.loadFromFile(filePath, config); err != nil {
		return nil, fmt.Errorf("failed to load config from %s: %w", filePath, err)
	}
	
	// Override with environment variables
	cm.loadFromEnvironment(config)
	
	// Validate configuration
	if err := cm.validateConfig(config); err != nil {
		return nil, fmt.Errorf("invalid configuration: %w", err)
	}
	
	// Set as current config
	cm.config = config
	return config, nil
}

// GetConfig returns the current configuration
func (cm *ConfigManager) GetConfig() *Config {
	if cm.config == nil {
		config, _ := cm.LoadConfig() // Use default on error
		return config
	}
	return cm.config
}

// getDefaultConfig returns the default configuration
func (cm *ConfigManager) getDefaultConfig() *Config {
	return &Config{
		General: GeneralConfig{
			DefaultCountry:     "US",
			DefaultLanguage:    "en",
			DefaultTimezone:    "UTC",
			SupportedLanguages: []string{"en", "es", "fr", "de"},
			Environment:        "prod",
		},
		Countries: make(map[string]CountryConfig),
		CustomHolidays: make(map[string][]CustomHoliday),
		Output: OutputConfig{
			DateFormat:      "2006-01-02",
			TimeFormat:      "15:04:05",
			Timezone:        "UTC",
			IncludeMetadata: true,
			Languages:       []string{"en"},
			Formats:         make(map[string]string),
		},
		Performance: PerformanceConfig{
			EnableCaching:   true,
			CacheTTL:        24 * time.Hour,
			MaxCacheSize:    1000,
			PreloadYears:    2,
			ConcurrentLimit: 10,
			BatchSize:       100,
		},
		Logging: LoggingConfig{
			Level:      "info",
			Format:     "text",
			Output:     "stdout",
			EnableFile: false,
			MaxSize:    100,
		},
	}
}

// loadFromFile loads configuration from a YAML file
func (cm *ConfigManager) loadFromFile(path string, config *Config) error {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return err
	}
	
	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	
	return yaml.Unmarshal(data, config)
}

// loadFromEnvironment loads configuration from environment variables
func (cm *ConfigManager) loadFromEnvironment(config *Config) {
	// General settings
	if env := os.Getenv("GOHOLIDAYS_DEFAULT_COUNTRY"); env != "" {
		config.General.DefaultCountry = env
	}
	if env := os.Getenv("GOHOLIDAYS_DEFAULT_LANGUAGE"); env != "" {
		config.General.DefaultLanguage = env
	}
	if env := os.Getenv("GOHOLIDAYS_DEFAULT_TIMEZONE"); env != "" {
		config.General.DefaultTimezone = env
	}
	if env := os.Getenv("GOHOLIDAYS_ENVIRONMENT"); env != "" {
		config.General.Environment = env
	}
	
	// Output settings
	if env := os.Getenv("GOHOLIDAYS_DATE_FORMAT"); env != "" {
		config.Output.DateFormat = env
	}
	if env := os.Getenv("GOHOLIDAYS_TIMEZONE"); env != "" {
		config.Output.Timezone = env
	}
	
	// Performance settings
	if env := os.Getenv("GOHOLIDAYS_ENABLE_CACHING"); env != "" {
		config.Performance.EnableCaching = strings.ToLower(env) == "true"
	}
	
	// Logging settings
	if env := os.Getenv("GOHOLIDAYS_LOG_LEVEL"); env != "" {
		config.Logging.Level = env
	}
}

// validateConfig validates the configuration
func (cm *ConfigManager) validateConfig(config *Config) error {
	// Validate timezone
	if config.General.DefaultTimezone != "" {
		if _, err := time.LoadLocation(config.General.DefaultTimezone); err != nil {
			return fmt.Errorf("invalid default timezone: %w", err)
		}
	}
	
	// Validate output timezone
	if config.Output.Timezone != "" && config.Output.Timezone != "Local" {
		if _, err := time.LoadLocation(config.Output.Timezone); err != nil {
			return fmt.Errorf("invalid output timezone: %w", err)
		}
	}
	
	// Validate environment
	validEnvs := []string{"dev", "development", "staging", "prod", "production"}
	valid := false
	for _, env := range validEnvs {
		if config.General.Environment == env {
			valid = true
			break
		}
	}
	if !valid {
		return fmt.Errorf("invalid environment: %s (must be one of: %v)", 
			config.General.Environment, validEnvs)
	}
	
	// Validate logging level
	validLevels := []string{"debug", "info", "warn", "error"}
	valid = false
	for _, level := range validLevels {
		if config.Logging.Level == level {
			valid = true
			break
		}
	}
	if !valid {
		return fmt.Errorf("invalid log level: %s (must be one of: %v)", 
			config.Logging.Level, validLevels)
	}
	
	return nil
}

// SaveConfig saves the current configuration to a file
func (cm *ConfigManager) SaveConfig(path string) error {
	if cm.config == nil {
		return fmt.Errorf("no configuration loaded")
	}
	
	data, err := yaml.Marshal(cm.config)
	if err != nil {
		return err
	}
	
	// Create directory if it doesn't exist
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}
	
	return os.WriteFile(path, data, 0644)
}

// GetCountryConfig returns configuration for a specific country
func (cm *ConfigManager) GetCountryConfig(countryCode string) CountryConfig {
	if cm.config == nil {
		return CountryConfig{Enabled: true}
	}
	
	if config, exists := cm.config.Countries[countryCode]; exists {
		return config
	}
	
	return CountryConfig{Enabled: true}
}

// IsCountryEnabled checks if a country is enabled
func (cm *ConfigManager) IsCountryEnabled(countryCode string) bool {
	config := cm.GetCountryConfig(countryCode)
	return config.Enabled
}

// GetCustomHolidays returns custom holidays for a country
func (cm *ConfigManager) GetCustomHolidays(countryCode string) []CustomHoliday {
	if cm.config == nil {
		return []CustomHoliday{}
	}
	
	// Get holidays for the specific country
	if holidays, exists := cm.config.CustomHolidays[countryCode]; exists {
		return holidays
	}
	
	// Also check for global holidays (if any are marked with "*")
	var globalHolidays []CustomHoliday
	if holidays, exists := cm.config.CustomHolidays["*"]; exists {
		globalHolidays = append(globalHolidays, holidays...)
	}
	
	return globalHolidays
}

// Global configuration manager instance
var DefaultConfigManager = NewConfigManager()
