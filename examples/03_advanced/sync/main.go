package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/coredds/GoHoliday/updater"
)

func main() {
	fmt.Println("GoHoliday Sync System Example")
	fmt.Println("============================")

	// Create temporary directory for demo
	tempDir, err := os.MkdirTemp("", "holiday-sync-demo")
	if err != nil {
		log.Fatal("Failed to create temp directory:", err)
	}
	defer os.RemoveAll(tempDir)

	// 1. Basic Sync Setup
	fmt.Println("\n1. Setting up Sync System")
	sync := updater.NewPythonHolidaysSync(tempDir)
	fmt.Printf("Sync system initialized with data directory: %s\n", tempDir)

	// 2. Check for Updates
	fmt.Println("\n2. Checking for Updates")
	hasUpdates, err := sync.CheckForUpdates()
	if err != nil {
		fmt.Printf("Error checking for updates: %v\n", err)
	} else {
		fmt.Printf("Updates available: %v\n", hasUpdates)
	}

	// 3. Manual Data Management
	fmt.Println("\n3. Manual Data Management")

	// Create sample holiday data
	sampleData := `{
		"US": {
			"2024-01-01": "New Year's Day",
			"2024-07-04": "Independence Day",
			"2024-12-25": "Christmas Day"
		}
	}`

	dataPath := filepath.Join(tempDir, "holiday_data.json")
	err = os.WriteFile(dataPath, []byte(sampleData), 0644)
	if err != nil {
		log.Fatal("Failed to write sample data:", err)
	}
	fmt.Printf("Sample holiday data written to: %s\n", dataPath)

	// 4. Sync Configuration
	fmt.Println("\n4. Sync Configuration")
	fmt.Printf("Sync configuration:\n")
	fmt.Printf("- Data directory: %s\n", tempDir)
	fmt.Printf("- Default update interval: 24h\n")
	fmt.Printf("- Default retry attempts: 3\n")

	// 5. Data Verification
	fmt.Println("\n5. Data Verification")
	verifyData(dataPath)

	// 6. Error Handling Demo
	fmt.Println("\n6. Error Handling")
	demonstrateErrorHandling(sync)

	// 7. Sync Status
	fmt.Println("\n7. Sync Status")
	printSyncStatus(sync)

	fmt.Println("\nSync system example completed!")
}

func verifyData(path string) {
	// Read and verify the data file
	data, err := os.ReadFile(path)
	if err != nil {
		fmt.Printf("Error reading data file: %v\n", err)
		return
	}

	fmt.Printf("Data file size: %d bytes\n", len(data))
	fmt.Println("Data verification completed")
}

func demonstrateErrorHandling(sync *updater.PythonHolidaysSync) {
	// Demonstrate various error scenarios
	scenarios := []struct {
		name string
		fn   func() error
	}{
		{
			name: "Invalid data directory",
			fn: func() error {
				invalidSync := updater.NewPythonHolidaysSync("/nonexistent")
				_, err := invalidSync.CheckForUpdates()
				return err
			},
		},
		{
			name: "Network timeout simulation",
			fn: func() error {
				// Simulate a network timeout
				time.Sleep(100 * time.Millisecond)
				return fmt.Errorf("network timeout")
			},
		},
	}

	for _, scenario := range scenarios {
		fmt.Printf("\nTesting scenario: %s\n", scenario.name)
		err := scenario.fn()
		if err != nil {
			fmt.Printf("Expected error occurred: %v\n", err)
		} else {
			fmt.Println("No error occurred")
		}
	}
}

func printSyncStatus(sync *updater.PythonHolidaysSync) {
	// Print various status information
	fmt.Println("Current sync status:")
	fmt.Printf("- Status: Active\n")
	fmt.Printf("- Next sync: Scheduled within 24h\n")
}
