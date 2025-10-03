package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/coredds/goholiday/config"
	"gopkg.in/yaml.v3"
)

func main() {
	fmt.Println("goholiday Configuration Management")
	fmt.Println("=================================")

	// 1. Basic Configuration Loading
	fmt.Println("\n1. Loading Configuration")
	cfg := &config.Config{
		General: config.GeneralConfig{
			DefaultCountry:     "US",
			DefaultLanguage:    "en",
			DefaultTimezone:    "UTC",
			SupportedLanguages: []string{"en", "es", "fr"},
			Environment:        "dev",
		},
		Performance: config.PerformanceConfig{
			EnableCaching: true,
			CacheTTL:      24 * time.Hour,
			MaxCacheSize:  1000,
			PreloadYears:  2,
		},
		Logging: config.LoggingConfig{
			Level:  "info",
			Format: "text",
			Output: "stdout",
		},
	}
	fmt.Printf("Default configuration created: %+v\n", cfg)

	// 2. Custom Configuration File
	fmt.Println("\n2. Creating and Loading Custom Configuration")
	customConfig := `
general:
  default_country: US
  default_language: en
  default_timezone: UTC
  supported_languages: [en, es, fr]
  environment: prod

performance:
  enable_caching: true
  cache_ttl: 24h
  max_cache_size: 1000
  preload_years: 2

logging:
  level: info
  format: json
  output: ./logs/holiday.log
  enable_file: true
  max_size: 100
`
	// Write custom config
	customPath := filepath.Join(".", "custom_config.yaml")
	err := os.WriteFile(customPath, []byte(customConfig), 0644)
	if err != nil {
		log.Fatal("Failed to write custom config:", err)
	}
	defer os.Remove(customPath) // Clean up after demo

	// Load custom config
	customCfg := &config.Config{}
	err = yaml.Unmarshal([]byte(customConfig), customCfg)
	if err != nil {
		log.Fatal("Failed to parse custom config:", err)
	}
	fmt.Printf("Custom configuration loaded: %+v\n", customCfg)

	// 3. Environment-Based Configuration
	fmt.Println("\n3. Environment-Specific Configuration")

	// Create development config
	devConfig := `
general:
  environment: development
  default_country: US
  default_language: en

performance:
  enable_caching: true
  cache_ttl: 1h
  preload_years: 1

logging:
  level: debug
  format: text
  output: stdout
  enable_file: false
`
	devPath := filepath.Join(".", "dev_config.yaml")
	err = os.WriteFile(devPath, []byte(devConfig), 0644)
	if err != nil {
		log.Fatal("Failed to write dev config:", err)
	}
	defer os.Remove(devPath)

	// Create production config
	prodConfig := `
general:
  environment: production
  default_country: US
  default_language: en

performance:
  enable_caching: true
  cache_ttl: 168h
  max_cache_size: 10000
  preload_years: 5

logging:
  level: warn
  format: json
  output: /var/log/holidays.log
  enable_file: true
  max_size: 500
`
	prodPath := filepath.Join(".", "prod_config.yaml")
	err = os.WriteFile(prodPath, []byte(prodConfig), 0644)
	if err != nil {
		log.Fatal("Failed to write prod config:", err)
	}
	defer os.Remove(prodPath)

	// Load based on environment
	envConfigs := map[string]string{
		"development": devPath,
		"production":  prodPath,
	}

	for env, path := range envConfigs {
		data, err := os.ReadFile(path)
		if err != nil {
			log.Printf("Failed to read %s config: %v\n", env, err)
			continue
		}

		cfg := &config.Config{}
		err = yaml.Unmarshal(data, cfg)
		if err != nil {
			log.Printf("Failed to parse %s config: %v\n", env, err)
			continue
		}

		fmt.Printf("\n%s configuration:\n", env)
		fmt.Printf("- Environment: %s\n", cfg.General.Environment)
		fmt.Printf("- Default Language: %s\n", cfg.General.DefaultLanguage)
		fmt.Printf("- Cache TTL: %v\n", cfg.Performance.CacheTTL)
		fmt.Printf("- Log Level: %s\n", cfg.Logging.Level)
		fmt.Printf("- Log Output: %s\n", cfg.Logging.Output)
	}

	// 4. Configuration Validation
	fmt.Println("\n4. Configuration Validation")
	invalidConfig := `
general:
  environment: invalid_env
  default_timezone: invalid_tz

logging:
  level: unknown_level
  format: invalid_format
`
	invalidPath := filepath.Join(".", "invalid_config.yaml")
	err = os.WriteFile(invalidPath, []byte(invalidConfig), 0644)
	if err != nil {
		log.Fatal("Failed to write invalid config:", err)
	}
	defer os.Remove(invalidPath)

	invalidCfg := &config.Config{}
	err = yaml.Unmarshal([]byte(invalidConfig), invalidCfg)
	fmt.Printf("Loading invalid config result: %v\n", err)

	// 5. Configuration Merging
	fmt.Println("\n5. Configuration Merging")
	baseConfig := `
general:
  default_country: US
  default_language: en

performance:
  enable_caching: true
  cache_ttl: 24h

logging:
  level: info
  format: text
`
	overrideConfig := `
general:
  default_language: es

performance:
  cache_ttl: 48h
  max_cache_size: 2000

logging:
  format: json
  output: ./custom.log
`
	basePath := filepath.Join(".", "base_config.yaml")
	overridePath := filepath.Join(".", "override_config.yaml")

	err = os.WriteFile(basePath, []byte(baseConfig), 0644)
	if err != nil {
		log.Fatal("Failed to write base config:", err)
	}
	defer os.Remove(basePath)

	err = os.WriteFile(overridePath, []byte(overrideConfig), 0644)
	if err != nil {
		log.Fatal("Failed to write override config:", err)
	}
	defer os.Remove(overridePath)

	baseCfg := &config.Config{}
	overrideCfg := &config.Config{}
	_ = yaml.Unmarshal([]byte(baseConfig), baseCfg)
	_ = yaml.Unmarshal([]byte(overrideConfig), overrideCfg)

	// Merge manually since we don't have a Merge function
	merged := &config.Config{
		General: config.GeneralConfig{
			DefaultCountry:  baseCfg.General.DefaultCountry,
			DefaultLanguage: baseCfg.General.DefaultLanguage,
		},
		Performance: config.PerformanceConfig{
			EnableCaching: baseCfg.Performance.EnableCaching,
			CacheTTL:      baseCfg.Performance.CacheTTL,
		},
		Logging: config.LoggingConfig{
			Level:  baseCfg.Logging.Level,
			Format: baseCfg.Logging.Format,
		},
	}

	// Override with values from override config
	if overrideCfg.General.DefaultLanguage != "" {
		merged.General.DefaultLanguage = overrideCfg.General.DefaultLanguage
	}
	if overrideCfg.Performance.CacheTTL != 0 {
		merged.Performance.CacheTTL = overrideCfg.Performance.CacheTTL
	}
	if overrideCfg.Performance.MaxCacheSize != 0 {
		merged.Performance.MaxCacheSize = overrideCfg.Performance.MaxCacheSize
	}
	if overrideCfg.Logging.Format != "" {
		merged.Logging.Format = overrideCfg.Logging.Format
	}
	if overrideCfg.Logging.Output != "" {
		merged.Logging.Output = overrideCfg.Logging.Output
	}

	fmt.Println("Merged configuration:")
	fmt.Printf("- Default Language: %s\n", merged.General.DefaultLanguage)
	fmt.Printf("- Cache TTL: %v\n", merged.Performance.CacheTTL)
	fmt.Printf("- Max Cache Size: %d\n", merged.Performance.MaxCacheSize)
	fmt.Printf("- Log Format: %s\n", merged.Logging.Format)
	fmt.Printf("- Log Output: %s\n", merged.Logging.Output)

	fmt.Println("\nConfiguration management example completed!")
}
