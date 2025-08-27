package goholidays

import (
	"fmt"
	"log"
	"time"
)

// APIVersion represents the API version
type APIVersion string

const (
	// APIVersionV1 represents version 1.0 of the API
	APIVersionV1 APIVersion = "v1.0"
	
	// CurrentAPIVersion is the current stable API version
	CurrentAPIVersion = APIVersionV1
)

// DeprecationLevel indicates the level of deprecation
type DeprecationLevel int

const (
	// DeprecationNone indicates no deprecation
	DeprecationNone DeprecationLevel = iota
	
	// DeprecationWarning indicates the feature is deprecated but still works
	DeprecationWarning
	
	// DeprecationError indicates the feature will be removed soon
	DeprecationError
	
	// DeprecationRemoved indicates the feature has been removed
	DeprecationRemoved
)

// DeprecationInfo holds information about deprecated features
type DeprecationInfo struct {
	Level       DeprecationLevel
	Message     string
	Replacement string
	RemovalDate *time.Time
}

// deprecationRegistry tracks deprecated features
var deprecationRegistry = make(map[string]DeprecationInfo)

// RegisterDeprecation registers a deprecated feature
func RegisterDeprecation(feature string, info DeprecationInfo) {
	deprecationRegistry[feature] = info
}

// CheckDeprecation checks if a feature is deprecated and logs warnings
func CheckDeprecation(feature string) {
	if info, exists := deprecationRegistry[feature]; exists {
		switch info.Level {
		case DeprecationWarning:
			log.Printf("DEPRECATION WARNING: %s is deprecated. %s", feature, info.Message)
			if info.Replacement != "" {
				log.Printf("Use %s instead.", info.Replacement)
			}
		case DeprecationError:
			log.Printf("DEPRECATION ERROR: %s will be removed soon. %s", feature, info.Message)
			if info.RemovalDate != nil {
				log.Printf("Scheduled for removal on %s", info.RemovalDate.Format("2006-01-02"))
			}
		case DeprecationRemoved:
			panic(fmt.Sprintf("DEPRECATED FEATURE REMOVED: %s has been removed. %s", feature, info.Message))
		}
	}
}

// APIStabilityLevel indicates the stability level of an API feature
type APIStabilityLevel int

const (
	// StabilityExperimental indicates experimental features that may change
	StabilityExperimental APIStabilityLevel = iota
	
	// StabilityBeta indicates beta features that are mostly stable
	StabilityBeta
	
	// StabilityStable indicates stable features with backward compatibility
	StabilityStable
	
	// StabilityFrozen indicates frozen features that will never change
	StabilityFrozen
)

// APIFeature represents an API feature with stability information
type APIFeature struct {
	Name        string
	Stability   APIStabilityLevel
	Since       APIVersion
	Description string
}

// featureRegistry tracks API features and their stability
var featureRegistry = make(map[string]APIFeature)

// RegisterAPIFeature registers an API feature
func RegisterAPIFeature(feature APIFeature) {
	featureRegistry[feature.Name] = feature
}

// GetAPIFeature retrieves information about an API feature
func GetAPIFeature(name string) (APIFeature, bool) {
	feature, exists := featureRegistry[name]
	return feature, exists
}

// ValidateAPIUsage validates that API features are used correctly
func ValidateAPIUsage(featureName string, requiredStability APIStabilityLevel) error {
	feature, exists := featureRegistry[featureName]
	if !exists {
		return fmt.Errorf("unknown API feature: %s", featureName)
	}
	
	if feature.Stability < requiredStability {
		stabilityNames := []string{"experimental", "beta", "stable", "frozen"}
		return fmt.Errorf("feature %s has %s stability, but %s required",
			featureName, stabilityNames[feature.Stability], stabilityNames[requiredStability])
	}
	
	return nil
}

// BackwardCompatibility provides backward compatibility helpers
type BackwardCompatibility struct {
	enabled bool
}

// NewBackwardCompatibility creates a new backward compatibility manager
func NewBackwardCompatibility(enabled bool) *BackwardCompatibility {
	return &BackwardCompatibility{enabled: enabled}
}

// IsEnabled returns whether backward compatibility is enabled
func (bc *BackwardCompatibility) IsEnabled() bool {
	return bc.enabled
}

// SetEnabled enables or disables backward compatibility
func (bc *BackwardCompatibility) SetEnabled(enabled bool) {
	bc.enabled = enabled
}

// GlobalBackwardCompatibility is the global backward compatibility manager
var GlobalBackwardCompatibility = NewBackwardCompatibility(true)

// init registers core API features
func init() {
	// Register core stable API features
	RegisterAPIFeature(APIFeature{
		Name:        "NewCountry",
		Stability:   StabilityStable,
		Since:       APIVersionV1,
		Description: "Creates a new country holiday provider",
	})
	
	RegisterAPIFeature(APIFeature{
		Name:        "IsHoliday",
		Stability:   StabilityStable,
		Since:       APIVersionV1,
		Description: "Checks if a given date is a holiday",
	})
	
	RegisterAPIFeature(APIFeature{
		Name:        "HolidaysForYear",
		Stability:   StabilityStable,
		Since:       APIVersionV1,
		Description: "Returns all holidays for a specific year",
	})
	
	RegisterAPIFeature(APIFeature{
		Name:        "HolidaysForDateRange",
		Stability:   StabilityStable,
		Since:       APIVersionV1,
		Description: "Returns holidays within a date range",
	})
	
	// Register beta features
	RegisterAPIFeature(APIFeature{
		Name:        "OptimizedHoliday",
		Stability:   StabilityBeta,
		Since:       APIVersionV1,
		Description: "Memory-optimized holiday creation",
	})
	
	RegisterAPIFeature(APIFeature{
		Name:        "HolidayCache",
		Stability:   StabilityBeta,
		Since:       APIVersionV1,
		Description: "LRU cache for holiday data",
	})
}

// VersionInfo provides version information
type VersionInfo struct {
	LibraryVersion string     `json:"library_version"`
	APIVersion     APIVersion `json:"api_version"`
	GoVersion      string     `json:"go_version"`
	BuildTime      string     `json:"build_time,omitempty"`
	GitCommit      string     `json:"git_commit,omitempty"`
}

// GetVersionInfo returns version information
func GetVersionInfo() VersionInfo {
	return VersionInfo{
		LibraryVersion: Version,
		APIVersion:     CurrentAPIVersion,
		GoVersion:      "1.23+",
	}
}

// APIContract defines the contract for country providers
type APIContract interface {
	// Core methods that must be implemented
	LoadHolidays(year int) map[time.Time]*Holiday
	GetCountryCode() string
	GetSupportedSubdivisions() []string
	GetSupportedCategories() []string
	
	// Optional methods for enhanced functionality
	GetName() string
	GetLanguages() []string
	IsSubdivisionSupported(subdivision string) bool
	IsCategorySupported(category string) bool
}

// CompatibilityValidator validates API compatibility
type CompatibilityValidator struct {
	version APIVersion
}

// NewCompatibilityValidator creates a new compatibility validator
func NewCompatibilityValidator(version APIVersion) *CompatibilityValidator {
	return &CompatibilityValidator{version: version}
}

// ValidateProvider validates that a provider meets the API contract
func (cv *CompatibilityValidator) ValidateProvider(provider APIContract) error {
	// Validate required methods
	if provider.GetCountryCode() == "" {
		return fmt.Errorf("provider must return a valid country code")
	}
	
	if provider.GetSupportedSubdivisions() == nil {
		return fmt.Errorf("provider must return supported subdivisions list")
	}
	
	if provider.GetSupportedCategories() == nil {
		return fmt.Errorf("provider must return supported categories list")
	}
	
	// Test holiday loading
	testYear := time.Now().Year()
	holidays := provider.LoadHolidays(testYear)
	if holidays == nil {
		return fmt.Errorf("provider must return a valid holidays map")
	}
	
	return nil
}

// MigrationGuide provides guidance for API migrations
type MigrationGuide struct {
	FromVersion APIVersion
	ToVersion   APIVersion
	Steps       []MigrationStep
}

// MigrationStep represents a single migration step
type MigrationStep struct {
	Description string
	OldCode     string
	NewCode     string
	Required    bool
}

// GetMigrationGuide returns migration guidance between API versions
func GetMigrationGuide(from, to APIVersion) *MigrationGuide {
	// For now, return empty guide as we only have v1.0
	return &MigrationGuide{
		FromVersion: from,
		ToVersion:   to,
		Steps:       []MigrationStep{},
	}
}
