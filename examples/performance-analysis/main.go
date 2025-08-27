package main

import (
	"fmt"
	"runtime"
	"time"

	"github.com/coredds/GoHoliday"
)

func main() {
	fmt.Println("🚀 GoHoliday Performance Analysis")
	fmt.Println("=================================")

	// Test countries
	countries := []string{"US", "CA", "GB", "AU", "NZ", "JP", "IN"}
	
	// Test years
	years := []int{2020, 2021, 2022, 2023, 2024, 2025, 2026}
	
	// Benchmark country creation
	fmt.Println("\n📊 Country Creation Performance:")
	fmt.Println("--------------------------------")
	
	start := time.Now()
	countryInstances := make(map[string]*goholidays.Country)
	for _, code := range countries {
		countryInstances[code] = goholidays.NewCountry(code)
	}
	duration := time.Since(start)
	fmt.Printf("Created %d countries in %v (%.2f µs per country)\n", 
		len(countries), duration, float64(duration.Nanoseconds())/float64(len(countries))/1000)

	// Benchmark holiday loading
	fmt.Println("\n📅 Holiday Loading Performance:")
	fmt.Println("-------------------------------")
	
	for _, code := range countries {
		country := countryInstances[code]
		start := time.Now()
		totalHolidays := 0
		
		for _, year := range years {
			holidays := country.HolidaysForYear(year)
			totalHolidays += len(holidays)
		}
		
		duration := time.Since(start)
		fmt.Printf("%-2s: %d holidays across %d years in %v (%.2f µs per holiday)\n",
			code, totalHolidays, len(years), duration,
			float64(duration.Nanoseconds())/float64(totalHolidays)/1000)
	}

	// Benchmark date lookups
	fmt.Println("\n🔍 Date Lookup Performance:")
	fmt.Println("----------------------------")
	
	testDates := []time.Time{
		time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),   // New Year
		time.Date(2024, 7, 4, 0, 0, 0, 0, time.UTC),   // US Independence Day
		time.Date(2024, 12, 25, 0, 0, 0, 0, time.UTC), // Christmas
		time.Date(2024, 6, 15, 0, 0, 0, 0, time.UTC),  // Random date
		time.Date(2024, 10, 31, 0, 0, 0, 0, time.UTC), // Halloween
	}
	
	for _, code := range countries {
		country := countryInstances[code]
		start := time.Now()
		holidayCount := 0
		
		for i := 0; i < 1000; i++ {
			for _, date := range testDates {
				if _, isHoliday := country.IsHoliday(date); isHoliday {
					holidayCount++
				}
			}
		}
		
		duration := time.Since(start)
		totalLookups := 1000 * len(testDates)
		fmt.Printf("%-2s: %d/%d holiday hits in %v (%.2f ns per lookup)\n",
			code, holidayCount, totalLookups, duration,
			float64(duration.Nanoseconds())/float64(totalLookups))
	}

	// Memory usage analysis
	fmt.Println("\n💾 Memory Usage Analysis:")
	fmt.Println("-------------------------")
	
	var m1, m2 runtime.MemStats
	runtime.GC()
	runtime.ReadMemStats(&m1)
	
	// Create additional countries and load holidays
	additionalCountries := make([]*goholidays.Country, 100)
	for i := 0; i < 100; i++ {
		additionalCountries[i] = goholidays.NewCountry("US")
		for _, year := range []int{2020, 2021, 2022, 2023, 2024} {
			additionalCountries[i].HolidaysForYear(year)
		}
	}
	
	runtime.GC()
	runtime.ReadMemStats(&m2)
	
	allocatedMB := float64(m2.TotalAlloc-m1.TotalAlloc) / 1024 / 1024
	fmt.Printf("Memory allocated for 100 countries: %.2f MB\n", allocatedMB)
	fmt.Printf("Average memory per country: %.2f KB\n", allocatedMB*1024/100)
	
	// Concurrent access test
	fmt.Println("\n🔄 Concurrent Access Performance:")
	fmt.Println("----------------------------------")
	
	country := goholidays.NewCountry("US")
	country.HolidaysForYear(2024) // Pre-load
	
	start = time.Now()
	numGoroutines := 100
	numOperations := 1000
	done := make(chan bool, numGoroutines)
	
	testDate := time.Date(2024, 7, 4, 0, 0, 0, 0, time.UTC)
	
	for i := 0; i < numGoroutines; i++ {
		go func() {
			for j := 0; j < numOperations; j++ {
				country.IsHoliday(testDate)
			}
			done <- true
		}()
	}
	
	for i := 0; i < numGoroutines; i++ {
		<-done
	}
	
	duration = time.Since(start)
	totalOps := numGoroutines * numOperations
	fmt.Printf("Concurrent access: %d goroutines × %d operations = %d total ops in %v\n",
		numGoroutines, numOperations, totalOps, duration)
	fmt.Printf("Rate: %.2f million ops/second\n", float64(totalOps)/duration.Seconds()/1_000_000)

	// Feature comparison
	fmt.Println("\n🌍 Country Feature Comparison:")
	fmt.Println("-------------------------------")
	
	for _, code := range countries {
		country := countryInstances[code]
		holidays2024 := country.HolidaysForYear(2024)
		
		categories := make(map[string]int)
		multiLanguage := 0
		
		for _, holiday := range holidays2024 {
			categories[string(holiday.Category)]++
			if len(holiday.Languages) > 1 {
				multiLanguage++
			}
		}
		
		fmt.Printf("%-2s: %2d holidays", code, len(holidays2024))
		if multiLanguage > 0 {
			fmt.Printf(", %d multilingual", multiLanguage)
		}
		
		var catList []string
		for cat := range categories {
			catList = append(catList, cat)
		}
		if len(catList) > 0 {
			fmt.Printf(", categories: %v", catList)
		}
		fmt.Println()
	}

	fmt.Println("\n✅ Performance analysis completed!")
	fmt.Printf("📈 System: %s/%s, Go %s\n", runtime.GOOS, runtime.GOARCH, runtime.Version())
}
