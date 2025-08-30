package main

import (
	"fmt"
	"runtime"
	"sync"
	"time"

	goholidays "github.com/coredds/GoHoliday"
)

func main() {
	fmt.Println("GoHoliday Performance Analysis")
	fmt.Println("=============================")

	// 1. Basic Performance Benchmarks
	fmt.Println("\n1. Basic Operations Performance")

	us := goholidays.NewCountry("US")
	benchmarkBasicOperations(us)

	// 2. Memory Usage Analysis
	fmt.Println("\n2. Memory Usage Analysis")
	analyzeMemoryUsage(us)

	// 3. Concurrent Access Performance
	fmt.Println("\n3. Concurrent Access Performance")
	benchmarkConcurrentAccess(us)

	// 4. Multi-Country Performance
	fmt.Println("\n4. Multi-Country Performance")
	benchmarkMultiCountry()

	// 5. Date Range Performance
	fmt.Println("\n5. Date Range Performance")
	benchmarkDateRanges(us)

	// 6. Cache Performance
	fmt.Println("\n6. Cache Performance")
	benchmarkCachePerformance(us)

	fmt.Println("\nPerformance analysis completed!")
}

func benchmarkBasicOperations(country *goholidays.Country) {
	operations := []struct {
		name string
		fn   func()
	}{
		{
			name: "Holiday Check",
			fn: func() {
				country.IsHoliday(time.Now())
			},
		},
		{
			name: "Year Holidays",
			fn: func() {
				country.HolidaysForYear(2024)
			},
		},
		{
			name: "Date Range",
			fn: func() {
				start := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
				end := time.Date(2024, 12, 31, 0, 0, 0, 0, time.UTC)
				country.HolidaysForDateRange(start, end)
			},
		},
	}

	iterations := 1000
	for _, op := range operations {
		start := time.Now()
		for i := 0; i < iterations; i++ {
			op.fn()
		}
		duration := time.Since(start)
		fmt.Printf("%s: %v total, %v per operation\n",
			op.name, duration, duration/time.Duration(iterations))
	}
}

func analyzeMemoryUsage(country *goholidays.Country) {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	fmt.Printf("Initial memory usage: %v MB\n", m.Alloc/1024/1024)

	// Load multiple years
	years := []int{2020, 2021, 2022, 2023, 2024, 2025}
	for _, year := range years {
		country.HolidaysForYear(year)
	}

	runtime.ReadMemStats(&m)
	fmt.Printf("After loading %d years: %v MB\n", len(years), m.Alloc/1024/1024)
}

func benchmarkConcurrentAccess(country *goholidays.Country) {
	workers := 10
	iterations := 1000
	var wg sync.WaitGroup

	start := time.Now()
	for w := 0; w < workers; w++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for i := 0; i < iterations; i++ {
				country.IsHoliday(time.Now())
			}
		}()
	}
	wg.Wait()
	duration := time.Since(start)

	totalOps := workers * iterations
	fmt.Printf("Concurrent operations (%d workers, %d iterations each):\n", workers, iterations)
	fmt.Printf("Total time: %v\n", duration)
	fmt.Printf("Average per operation: %v\n", duration/time.Duration(totalOps))
}

func benchmarkMultiCountry() {
	countries := []string{"US", "GB", "JP", "AU", "DE", "FR", "ES", "IT"}
	providers := make([]*goholidays.Country, len(countries))

	// Initialize providers
	start := time.Now()
	for i, code := range countries {
		providers[i] = goholidays.NewCountry(code)
	}
	initDuration := time.Since(start)
	fmt.Printf("Provider initialization (%d countries): %v\n", len(countries), initDuration)

	// Test holiday checks
	date := time.Date(2024, 12, 25, 0, 0, 0, 0, time.UTC)
	start = time.Now()
	for _, provider := range providers {
		provider.IsHoliday(date)
	}
	checkDuration := time.Since(start)
	fmt.Printf("Holiday check across all countries: %v\n", checkDuration)
}

func benchmarkDateRanges(country *goholidays.Country) {
	ranges := []struct {
		name  string
		start time.Time
		end   time.Time
	}{
		{
			name:  "Month",
			start: time.Date(2024, 12, 1, 0, 0, 0, 0, time.UTC),
			end:   time.Date(2024, 12, 31, 0, 0, 0, 0, time.UTC),
		},
		{
			name:  "Quarter",
			start: time.Date(2024, 10, 1, 0, 0, 0, 0, time.UTC),
			end:   time.Date(2024, 12, 31, 0, 0, 0, 0, time.UTC),
		},
		{
			name:  "Year",
			start: time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
			end:   time.Date(2024, 12, 31, 0, 0, 0, 0, time.UTC),
		},
	}

	iterations := 100
	for _, r := range ranges {
		start := time.Now()
		for i := 0; i < iterations; i++ {
			country.HolidaysForDateRange(r.start, r.end)
		}
		duration := time.Since(start)
		fmt.Printf("%s range: %v total, %v per operation\n",
			r.name, duration, duration/time.Duration(iterations))
	}
}

func benchmarkCachePerformance(country *goholidays.Country) {
	// First access (cold cache)
	start := time.Now()
	country.HolidaysForYear(2024)
	coldDuration := time.Since(start)
	fmt.Printf("Cold cache access: %v\n", coldDuration)

	// Second access (warm cache)
	start = time.Now()
	country.HolidaysForYear(2024)
	warmDuration := time.Since(start)
	fmt.Printf("Warm cache access: %v\n", warmDuration)
	fmt.Printf("Cache performance improvement: %.2fx\n", float64(coldDuration)/float64(warmDuration))
}
