package main

import (
	"context"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/coredds/GoHoliday/updater"
)

func TestSyncFunctionality(t *testing.T) {
	// Create temporary directory
	tempDir, err := os.MkdirTemp("", "goholidays-sync-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Test sync with no existing data
	t.Run("Initial Sync", func(t *testing.T) {
		syncer := updater.NewMockSyncer()
		err := syncSingleCountry(context.Background(), syncer, "US", tempDir, false, false)
		if err != nil {
			t.Errorf("Initial sync failed: %v", err)
		}

		// Check that files were created
		files, err := os.ReadDir(tempDir)
		if err != nil {
			t.Fatalf("Failed to read directory: %v", err)
		}

		if len(files) == 0 {
			t.Error("Sync should create holiday data files")
		}
	})

	// Test sync with existing data
	t.Run("Update Sync", func(t *testing.T) {
		syncer := updater.NewMockSyncer()
		err := syncSingleCountry(context.Background(), syncer, "US", tempDir, false, false)
		if err != nil {
			t.Errorf("Update sync failed: %v", err)
		}

		// Check that files were updated
		files, err := os.ReadDir(tempDir)
		if err != nil {
			t.Fatalf("Failed to read directory: %v", err)
		}

		if len(files) == 0 {
			t.Error("Sync should update holiday data files")
		}
	})

	// Test sync specific country
	t.Run("Single Country Sync", func(t *testing.T) {
		syncer := updater.NewMockSyncer()
		err := syncSingleCountry(context.Background(), syncer, "US", tempDir, false, false)
		if err != nil {
			t.Errorf("Single country sync failed: %v", err)
		}

		// Check that country file exists
		countryFile := filepath.Join(tempDir, "US.json")
		if _, err := os.Stat(countryFile); os.IsNotExist(err) {
			t.Error("Country file should exist after sync")
		}
	})
}

func TestErrorHandling(t *testing.T) {
	// Test invalid directory
	t.Run("Invalid Directory", func(t *testing.T) {
		syncer := updater.NewMockSyncer()
		// Use a path with invalid characters that can't be created
		invalidPath := "invalid\x00path\x00with\x00nulls"
		err := syncSingleCountry(context.Background(), syncer, "US", invalidPath, false, false)
		if err == nil {
			t.Error("Should show error for invalid directory")
		}
	})

	// Test invalid country code
	t.Run("Invalid Country", func(t *testing.T) {
		tempDir, err := os.MkdirTemp("", "goholidays-sync-test-*")
		if err != nil {
			t.Fatalf("Failed to create temp directory: %v", err)
		}
		defer os.RemoveAll(tempDir)

		syncer := updater.NewMockSyncer()
		err = syncSingleCountry(context.Background(), syncer, "XX", tempDir, false, false)
		if err == nil {
			t.Error("Should show error for invalid country code")
		}
	})

	// Test unwritable directory
	t.Run("Unwritable Directory", func(t *testing.T) {
		tempDir, err := os.MkdirTemp("", "goholidays-sync-test-*")
		if err != nil {
			t.Fatalf("Failed to create temp directory: %v", err)
		}
		defer func() {
			// Reset permissions before cleanup
			_ = os.Chmod(tempDir, 0700)
			_ = os.RemoveAll(tempDir)
		}()

		// Create a file with the same name as the expected output file to block writing
		blockingFile := filepath.Join(tempDir, "US.json")
		if err := os.WriteFile(blockingFile, []byte("blocking"), 0400); err != nil {
			t.Fatalf("Failed to create blocking file: %v", err)
		}

		// Make the blocking file read-only
		if err := os.Chmod(blockingFile, 0400); err != nil {
			t.Fatalf("Failed to change file permissions: %v", err)
		}

		syncer := updater.NewMockSyncer()
		err = syncSingleCountry(context.Background(), syncer, "US", tempDir, false, false)
		if err == nil {
			t.Error("Should show error for unwritable directory")
		}
	})
}

func TestConcurrentSync(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping concurrent sync test in short mode")
	}

	tempDir, err := os.MkdirTemp("", "goholidays-sync-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Test concurrent syncs
	t.Run("Concurrent Syncs", func(t *testing.T) {
		done := make(chan bool)
		for i := 0; i < 3; i++ {
			go func() {
				syncer := updater.NewMockSyncer()
				err := syncSingleCountry(context.Background(), syncer, "US", tempDir, false, false)
				if err != nil {
					t.Errorf("Concurrent sync failed: %v", err)
				}
				done <- true
			}()
		}

		// Wait for all syncs
		for i := 0; i < 3; i++ {
			<-done
		}

		// Check that files are consistent
		files, err := os.ReadDir(tempDir)
		if err != nil {
			t.Fatalf("Failed to read directory: %v", err)
		}

		if len(files) == 0 {
			t.Error("Sync should create holiday data files")
		}
	})
}

func TestPerformance(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping performance test in short mode")
	}

	tempDir, err := os.MkdirTemp("", "goholidays-sync-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Test sync performance
	t.Run("Sync Performance", func(t *testing.T) {
		start := time.Now()
		syncer := updater.NewMockSyncer()
		err := syncSingleCountry(context.Background(), syncer, "US", tempDir, false, false)
		if err != nil {
			t.Errorf("Sync failed: %v", err)
		}
		duration := time.Since(start)

		t.Logf("Full sync took: %v", duration)
		if duration > 30*time.Second {
			t.Error("Sync took too long")
		}
	})

	// Test incremental sync performance
	t.Run("Incremental Sync Performance", func(t *testing.T) {
		start := time.Now()
		syncer := updater.NewMockSyncer()
		err := syncSingleCountry(context.Background(), syncer, "US", tempDir, false, false)
		if err != nil {
			t.Errorf("Incremental sync failed: %v", err)
		}
		duration := time.Since(start)

		t.Logf("Incremental sync took: %v", duration)
		if duration > 5*time.Second {
			t.Error("Incremental sync took too long")
		}
	})
}
