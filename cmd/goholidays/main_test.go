package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"io"
	"os"
	"strings"
	"testing"
	"time"

	goholidays "github.com/coredds/goholiday"
)

// Helper function to capture stdout
func captureOutput(fn func()) string {
	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	fn()

	w.Close()
	os.Stdout = old

	var buf bytes.Buffer
	_, _ = io.Copy(&buf, r)
	return buf.String()
}

func TestCheckSpecificDate(t *testing.T) {
	country := goholidays.NewCountry("US")

	// Test table output
	t.Run("Table Output", func(t *testing.T) {
		output := captureOutput(func() {
			checkSpecificDate(country, "2024-07-04", "table", false)
		})

		if !strings.Contains(output, "Independence Day") {
			t.Error("Output should contain holiday name")
		}
	})

	// Test JSON output
	t.Run("JSON Output", func(t *testing.T) {
		output := captureOutput(func() {
			checkSpecificDate(country, "2024-07-04", "json", false)
		})

		var result map[string]interface{}
		if err := json.Unmarshal([]byte(output), &result); err != nil {
			t.Fatalf("Failed to parse JSON output: %v", err)
		}

		if !result["is_holiday"].(bool) {
			t.Error("July 4th should be marked as holiday")
		}
	})

	// Test business day info
	t.Run("Business Day Info", func(t *testing.T) {
		output := captureOutput(func() {
			checkSpecificDate(country, "2024-07-04", "table", true)
		})

		if !strings.Contains(output, "business day") {
			t.Error("Output should contain business day information")
		}
	})

	// Test non-holiday date
	t.Run("Non-Holiday Date", func(t *testing.T) {
		output := captureOutput(func() {
			checkSpecificDate(country, "2024-03-15", "table", false)
		})

		if !strings.Contains(output, "not a holiday") {
			t.Error("Output should indicate non-holiday")
		}
	})

	// Test invalid date format
	t.Run("Invalid Date Format", func(t *testing.T) {
		exitCalled := false
		osExit = func(code int) {
			exitCalled = true
		}
		defer func() {
			osExit = os.Exit // Reset
		}()

		checkSpecificDate(country, "invalid-date", "table", false)

		if !exitCalled {
			t.Error("Invalid date should cause exit")
		}
	})
}

func TestListHolidaysForYear(t *testing.T) {
	country := goholidays.NewCountry("US")
	year := 2024

	// Test table output
	t.Run("Table Output", func(t *testing.T) {
		output := captureOutput(func() {
			listHolidaysForYear(country, year, "table")
		})

		expectedHeaders := []string{"Date", "Holiday", "Category", "Observed"}
		for _, header := range expectedHeaders {
			if !strings.Contains(output, header) {
				t.Errorf("Output should contain header '%s'", header)
			}
		}

		expectedHolidays := []string{"New Year's Day", "Independence Day", "Christmas Day"}
		for _, holiday := range expectedHolidays {
			if !strings.Contains(output, holiday) {
				t.Errorf("Output should contain holiday '%s'", holiday)
			}
		}
	})

	// Test JSON output
	t.Run("JSON Output", func(t *testing.T) {
		output := captureOutput(func() {
			listHolidaysForYear(country, year, "json")
		})

		var holidays map[string]interface{}
		if err := json.Unmarshal([]byte(output), &holidays); err != nil {
			t.Fatalf("Failed to parse JSON output: %v", err)
		}

		if len(holidays) == 0 {
			t.Error("JSON output should contain holidays")
		}
	})

	// Test CSV output
	t.Run("CSV Output", func(t *testing.T) {
		output := captureOutput(func() {
			listHolidaysForYear(country, year, "csv")
		})

		lines := strings.Split(output, "\n")
		if len(lines) < 2 {
			t.Error("CSV output should have header and data rows")
		}

		if !strings.HasPrefix(lines[0], "Date,Name,Category,Observed") {
			t.Error("CSV should have correct header")
		}
	})
}

func TestShowCalendar(t *testing.T) {
	country := goholidays.NewCountry("US")
	year := 2024
	month := time.July

	output := captureOutput(func() {
		showCalendar(country, year, month)
	})

	// Check calendar formatting
	if !strings.Contains(output, "July 2024") {
		t.Error("Calendar should show month and year")
	}

	if !strings.Contains(output, "Su Mo Tu We Th Fr Sa") {
		t.Error("Calendar should show weekday headers")
	}

	// Check holiday marking
	if !strings.Contains(output, "4") { // July 4th
		t.Error("Calendar should show Independence Day")
	}
}

func TestListSupportedCountries(t *testing.T) {
	output := captureOutput(func() {
		listSupportedCountries()
	})

	// Check header
	if !strings.Contains(output, "Supported Countries:") {
		t.Error("Output should have header")
	}

	// Check format
	if !strings.Contains(output, "Code  Name") {
		t.Error("Output should have column headers")
	}

	// Check some countries
	expectedCountries := []string{"US", "GB", "CA", "AU"}
	for _, country := range expectedCountries {
		if !strings.Contains(output, country) {
			t.Errorf("Output should list country '%s'", country)
		}
	}
}

func TestMainFunctionality(t *testing.T) {
	// Save original args
	oldArgs := os.Args
	defer func() { os.Args = oldArgs }()

	// Test version flag
	t.Run("Version Flag", func(t *testing.T) {
		// Reset flag package state
		flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError)

		os.Args = []string{"cmd", "-version"}
		output := captureOutput(func() {
			main()
		})

		if !strings.Contains(output, "goholidays CLI") {
			t.Error("Version output should contain CLI version")
		}
	})

	// Test list flag
	t.Run("List Flag", func(t *testing.T) {
		// Reset flag package state
		flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError)

		os.Args = []string{"cmd", "-list"}
		output := captureOutput(func() {
			main()
		})

		if !strings.Contains(output, "Supported Countries:") {
			t.Error("List output should show supported countries")
		}
	})

	// Test missing country
	t.Run("Missing Country", func(t *testing.T) {
		// Reset flag package state
		flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError)

		exitCalled := false
		osExit = func(code int) {
			exitCalled = true
		}
		defer func() {
			osExit = os.Exit // Reset
		}()

		os.Args = []string{"cmd"}
		output := captureOutput(func() {
			main()
		})

		if !exitCalled {
			t.Error("Missing country should cause exit")
		}

		if !strings.Contains(output, "Error: country code is required") {
			t.Error("Should show error for missing country")
		}
	})

	// Test valid country with options
	t.Run("Valid Country", func(t *testing.T) {
		// Reset flag package state
		flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError)

		os.Args = []string{
			"cmd",
			"-country", "US",
			"-year", "2024",
			"-subdivisions", "CA,NY",
			"-language", "es",
			"-format", "json",
		}
		output := captureOutput(func() {
			main()
		})

		var holidays map[string]interface{}
		if err := json.Unmarshal([]byte(output), &holidays); err != nil {
			t.Fatalf("Failed to parse JSON output: %v", err)
		}

		if len(holidays) == 0 {
			t.Error("Should return holidays for valid country")
		}
	})

	// Test calendar view
	t.Run("Calendar View", func(t *testing.T) {
		// Reset flag package state
		flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError)

		os.Args = []string{
			"cmd",
			"-country", "US",
			"-calendar",
			"-month", "7",
			"-year", "2024",
		}
		output := captureOutput(func() {
			main()
		})

		if !strings.Contains(output, "July 2024") {
			t.Error("Should show calendar for specified month")
		}
	})

	// Test specific date check
	t.Run("Date Check", func(t *testing.T) {
		// Reset flag package state
		flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError)

		os.Args = []string{
			"cmd",
			"-country", "US",
			"-date", "2024-07-04",
			"-format", "json",
		}
		output := captureOutput(func() {
			main()
		})

		var result map[string]interface{}
		if err := json.Unmarshal([]byte(output), &result); err != nil {
			t.Fatalf("Failed to parse JSON output: %v", err)
		}

		if !result["is_holiday"].(bool) {
			t.Error("Should identify holiday correctly")
		}
	})
}

func TestEdgeCases(t *testing.T) {
	// Test invalid date format
	t.Run("Invalid Date Format", func(t *testing.T) {
		exitCalled := false
		osExit = func(code int) {
			exitCalled = true
		}
		defer func() {
			osExit = os.Exit // Reset
		}()

		country := goholidays.NewCountry("US")
		checkSpecificDate(country, "invalid", "table", false)

		if !exitCalled {
			t.Error("Invalid date should cause exit")
		}
	})

	// Test invalid format
	t.Run("Invalid Format", func(t *testing.T) {
		country := goholidays.NewCountry("US")
		output := captureOutput(func() {
			listHolidaysForYear(country, 2024, "invalid")
		})

		// Should default to table format
		if !strings.Contains(output, "Date") || !strings.Contains(output, "Holiday") {
			t.Error("Should default to table format for invalid format")
		}
	})

	// Test invalid month
	t.Run("Invalid Month", func(t *testing.T) {
		country := goholidays.NewCountry("US")
		output := captureOutput(func() {
			showCalendar(country, 2024, 13)
		})

		if !strings.Contains(output, "Error: Invalid month") {
			t.Error("Invalid month should produce error message")
		}
	})

	// Test empty subdivisions
	t.Run("Empty Subdivisions", func(t *testing.T) {
		// Reset flag package state
		flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError)

		os.Args = []string{
			"cmd",
			"-country", "US",
			"-subdivisions", "",
		}
		output := captureOutput(func() {
			main()
		})

		if strings.Contains(output, "error") {
			t.Error("Empty subdivisions should not cause error")
		}
	})

	// Test whitespace in subdivisions
	t.Run("Whitespace Subdivisions", func(t *testing.T) {
		// Reset flag package state
		flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError)

		os.Args = []string{
			"cmd",
			"-country", "US",
			"-subdivisions", " CA , NY ",
		}
		output := captureOutput(func() {
			main()
		})

		if strings.Contains(output, "error") {
			t.Error("Whitespace in subdivisions should be handled")
		}
	})
}

func TestPerformance(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping performance test in short mode")
	}

	country := goholidays.NewCountry("US")

	// Test large year range
	t.Run("Large Year Range", func(t *testing.T) {
		start := time.Now()
		for year := 2000; year <= 2050; year++ {
			listHolidaysForYear(country, year, "json")
		}
		duration := time.Since(start)
		t.Logf("Processing 50 years took: %v", duration)
	})

	// Test concurrent calendar views
	t.Run("Concurrent Calendar Views", func(t *testing.T) {
		start := time.Now()
		for month := time.January; month <= time.December; month++ {
			showCalendar(country, 2024, month)
		}
		duration := time.Since(start)
		t.Logf("Processing 12 months took: %v", duration)
	})
}
