package updater

import (
	"context"
)

// Syncer defines the interface for syncing holiday data from external sources
type Syncer interface {
	// FetchCountryList retrieves the list of available country modules
	FetchCountryList(ctx context.Context) ([]string, error)

	// FetchCountryFile retrieves the source file for a specific country
	FetchCountryFile(ctx context.Context, countryCode string) (string, error)

	// ParseHolidayDefinitions parses source content into structured holiday data
	ParseHolidayDefinitions(source string) (*CountryData, error)

	// ValidatePythonContent validates Python source content
	ValidatePythonContent(content string) error
}

// Ensure GitHubSyncer implements the Syncer interface
var _ Syncer = (*GitHubSyncer)(nil)
