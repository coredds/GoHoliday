package updater

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

// SyncAll synchronizes all countries
func (phs *PythonHolidaysSync) SyncAll(ctx context.Context) error {
	// Create data directory if it doesn't exist
	if err := os.MkdirAll(phs.dataDir, 0755); err != nil {
		return fmt.Errorf("failed to create data directory: %w", err)
	}

	// Get list of supported countries
	countries, err := phs.getSupportedCountries()
	if err != nil {
		return fmt.Errorf("failed to get supported countries: %w", err)
	}

	// Sync each country
	for _, country := range countries {
		if ctx != nil {
			select {
			case <-ctx.Done():
				return ctx.Err()
			default:
			}
		}

		if err := phs.SyncCountry(ctx, country.CountryCode); err != nil {
			return fmt.Errorf("failed to sync %s: %w", country.CountryCode, err)
		}
	}

	// Update last sync time
	if err := phs.SaveLastSyncTime(time.Now()); err != nil {
		return fmt.Errorf("failed to save last sync time: %w", err)
	}

	return nil
}

// SyncCountry synchronizes a specific country
func (phs *PythonHolidaysSync) SyncCountry(ctx context.Context, countryCode string) error {
	countryData, err := phs.fetchCountryData(countryCode)
	if err != nil {
		return fmt.Errorf("failed to fetch country data: %w", err)
	}

	if err := phs.SaveCountryData(countryCode, countryData); err != nil {
		return fmt.Errorf("failed to save country data: %w", err)
	}

	return nil
}

// SaveCountryData saves country data to disk
func (phs *PythonHolidaysSync) SaveCountryData(countryCode string, data *CountryData) error {
	if countryCode == "" {
		return fmt.Errorf("country code is required")
	}
	if data == nil {
		return fmt.Errorf("country data is required")
	}
	return phs.saveCountryData(data)
}

// GetLastSyncTime returns the time of the last successful sync
func (phs *PythonHolidaysSync) GetLastSyncTime() (time.Time, error) {
	return phs.getLastSyncTime()
}

// SaveLastSyncTime saves the time of the last successful sync
func (phs *PythonHolidaysSync) SaveLastSyncTime(t time.Time) error {
	return phs.saveLastSyncTime(t)
}

// getLastSyncTime reads the last sync time from disk
func (phs *PythonHolidaysSync) getLastSyncTime() (time.Time, error) {
	path := filepath.Join(phs.dataDir, "last_sync")
	data, err := os.ReadFile(path)
	if err != nil {
		return time.Time{}, err
	}
	return time.Parse(time.RFC3339, string(data))
}

// saveLastSyncTime writes the last sync time to disk
func (phs *PythonHolidaysSync) saveLastSyncTime(t time.Time) error {
	path := filepath.Join(phs.dataDir, "last_sync")
	return os.WriteFile(path, []byte(t.Format(time.RFC3339)), 0644)
}

