package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/your-username/goholidays"
)

func main() {
	var (
		country      = flag.String("country", "", "Country code (e.g., US, GB, CA)")
		year         = flag.Int("year", time.Now().Year(), "Year to get holidays for")
		date         = flag.String("date", "", "Check if specific date is a holiday (YYYY-MM-DD)")
		subdivisions = flag.String("subdivisions", "", "Comma-separated list of subdivisions")
		language     = flag.String("language", "en", "Language for holiday names")
		format       = flag.String("format", "table", "Output format: table, json, csv")
		list         = flag.Bool("list", false, "List all supported countries")
		version      = flag.Bool("version", false, "Show version information")
		business     = flag.Bool("business", false, "Show business day information")
		calendar     = flag.Bool("calendar", false, "Show calendar view for the month")
		month        = flag.Int("month", int(time.Now().Month()), "Month for calendar view (1-12)")
	)
	flag.Parse()

	if *version {
		fmt.Println("GoHolidays CLI v1.0.0")
		fmt.Println("A Go library for comprehensive holiday data")
		return
	}

	if *list {
		listSupportedCountries()
		return
	}

	if *country == "" {
		fmt.Println("Error: country code is required")
		flag.Usage()
		os.Exit(1)
	}

	// Parse subdivisions
	var subs []string
	if *subdivisions != "" {
		subs = strings.Split(*subdivisions, ",")
		for i, sub := range subs {
			subs[i] = strings.TrimSpace(sub)
		}
	}

	// Create country with options
	options := goholidays.CountryOptions{
		Subdivisions: subs,
		Language:     *language,
	}
	countryProvider := goholidays.NewCountry(*country, options)

	if *calendar {
		showCalendar(countryProvider, *year, time.Month(*month))
	} else if *date != "" {
		checkSpecificDate(countryProvider, *date, *format, *business)
	} else {
		listHolidaysForYear(countryProvider, *year, *format)
	}
}

func checkSpecificDate(country *goholidays.Country, dateStr, format string, showBusiness bool) {
	date, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		log.Fatalf("Invalid date format. Use YYYY-MM-DD: %v", err)
	}

	holiday, isHoliday := country.IsHoliday(date)
	
	switch format {
	case "json":
		result := map[string]interface{}{
			"date":      dateStr,
			"is_holiday": isHoliday,
		}
		if isHoliday {
			result["holiday"] = holiday
		}
		json.NewEncoder(os.Stdout).Encode(result)
	default:
		if isHoliday {
			fmt.Printf("%s is a holiday: %s\n", dateStr, holiday.Name)
			if holiday.IsObserved && holiday.Observed != nil {
				fmt.Printf("Observed on: %s\n", holiday.Observed.Format("2006-01-02"))
			}
		} else {
			fmt.Printf("%s is not a holiday\n", dateStr)
		}
		
		if showBusiness {
			calc := goholidays.NewBusinessDayCalculator(country)
			if calc.IsBusinessDay(date) {
				fmt.Printf("%s is a business day\n", dateStr)
			} else {
				fmt.Printf("%s is not a business day\n", dateStr)
				nextBusiness := calc.NextBusinessDay(date)
				fmt.Printf("Next business day: %s\n", nextBusiness.Format("2006-01-02"))
			}
		}
	}
}

func showCalendar(country *goholidays.Country, year int, month time.Month) {
	calendar := goholidays.NewHolidayCalendar(country)
	calendar.PrintMonth(year, month)
}

func listHolidaysForYear(country *goholidays.Country, year int, format string) {
	holidays := country.HolidaysForYear(year)

	switch format {
	case "json":
		json.NewEncoder(os.Stdout).Encode(holidays)
	case "csv":
		fmt.Println("Date,Name,Category,Observed")
		for date, holiday := range holidays {
			observed := ""
			if holiday.IsObserved && holiday.Observed != nil {
				observed = holiday.Observed.Format("2006-01-02")
			}
			fmt.Printf("%s,%s,%s,%s\n", 
				date.Format("2006-01-02"), 
				holiday.Name, 
				holiday.Category, 
				observed)
		}
	default:
		fmt.Printf("Holidays for %s in %d:\n\n", country.GetCountryCode(), year)
		fmt.Printf("%-12s %-30s %-12s %-12s\n", "Date", "Holiday", "Category", "Observed")
		fmt.Println(strings.Repeat("-", 70))
		
		// Convert map to slice for sorting
		type holidayDate struct {
			date    time.Time
			holiday *goholidays.Holiday
		}
		var sortedHolidays []holidayDate
		for date, holiday := range holidays {
			sortedHolidays = append(sortedHolidays, holidayDate{date, holiday})
		}
		
		// Simple sort by date
		for i := 0; i < len(sortedHolidays); i++ {
			for j := i + 1; j < len(sortedHolidays); j++ {
				if sortedHolidays[i].date.After(sortedHolidays[j].date) {
					sortedHolidays[i], sortedHolidays[j] = sortedHolidays[j], sortedHolidays[i]
				}
			}
		}
		
		for _, hd := range sortedHolidays {
			observed := ""
			if hd.holiday.IsObserved && hd.holiday.Observed != nil {
				observed = hd.holiday.Observed.Format("01-02")
			}
			fmt.Printf("%-12s %-30s %-12s %-12s\n",
				hd.date.Format("2006-01-02"),
				hd.holiday.Name,
				hd.holiday.Category,
				observed)
		}
	}
}

func listSupportedCountries() {
	countries := []struct {
		Code string
		Name string
	}{
		{"US", "United States"},
		{"GB", "United Kingdom"},
		{"CA", "Canada"},
		{"AU", "Australia"},
		{"DE", "Germany"},
		{"FR", "France"},
		// Add more as implemented
	}

	fmt.Println("Supported Countries:")
	fmt.Println("Code  Name")
	fmt.Println(strings.Repeat("-", 25))
	for _, country := range countries {
		fmt.Printf("%-4s  %s\n", country.Code, country.Name)
	}
	fmt.Println("\nNote: More countries will be added as the library develops")
}
