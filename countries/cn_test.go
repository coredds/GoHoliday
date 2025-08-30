package countries

import (
	"fmt"
	"testing"
	"time"
)

func TestCNProvider_NewCNProvider(t *testing.T) {
	provider := NewCNProvider()

	if provider == nil {
		t.Fatal("Expected provider to be created")
	}

	if provider.GetCountryCode() != "CN" {
		t.Errorf("Expected country code CN, got %s", provider.GetCountryCode())
	}

	subdivisions := provider.GetSupportedSubdivisions()
	if len(subdivisions) != 34 {
		t.Errorf("Expected 34 subdivisions, got %d", len(subdivisions))
	}

	categories := provider.GetSupportedCategories()
	expectedCategories := []string{"public", "traditional", "lunar"}
	if len(categories) != len(expectedCategories) {
		t.Errorf("Expected %d categories, got %d", len(expectedCategories), len(categories))
	}
}

func TestCNProvider_LoadHolidays(t *testing.T) {
	provider := NewCNProvider()
	year := 2024

	holidays := provider.LoadHolidays(year)

	// Test some key Chinese holidays
	testCases := []struct {
		date     time.Time
		name     string
		category string
	}{
		{
			date:     time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
			name:     "元旦",
			category: "public",
		},
		{
			date:     time.Date(2024, 5, 1, 0, 0, 0, 0, time.UTC),
			name:     "劳动节",
			category: "public",
		},
		{
			date:     time.Date(2024, 10, 1, 0, 0, 0, 0, time.UTC),
			name:     "国庆节",
			category: "public",
		},
		{
			date:     time.Date(2024, 10, 2, 0, 0, 0, 0, time.UTC),
			name:     "国庆节第二天",
			category: "public",
		},
		{
			date:     time.Date(2024, 10, 3, 0, 0, 0, 0, time.UTC),
			name:     "国庆节第三天",
			category: "public",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			holiday, exists := holidays[tc.date]
			if !exists {
				t.Errorf("Expected holiday %s on %s to exist", tc.name, tc.date.Format("2006-01-02"))
				return
			}

			if holiday.Name != tc.name {
				t.Errorf("Expected holiday name %s, got %s", tc.name, holiday.Name)
			}

			if holiday.Category != tc.category {
				t.Errorf("Expected category %s, got %s", tc.category, holiday.Category)
			}
		})
	}
}

func TestCNProvider_LunarHolidays2024(t *testing.T) {
	provider := NewCNProvider()
	year := 2024

	holidays := provider.LoadHolidays(year)

	// Test lunar calendar-based holidays for 2024
	testCases := []struct {
		date time.Time
		name string
	}{
		{
			// Spring Festival Eve (Chinese New Year's Eve)
			date: time.Date(2024, 2, 9, 0, 0, 0, 0, time.UTC),
			name: "除夕",
		},
		{
			// Spring Festival Day 1 (Chinese New Year)
			date: time.Date(2024, 2, 10, 0, 0, 0, 0, time.UTC),
			name: "春节第一天",
		},
		{
			// Qingming Festival
			date: time.Date(2024, 4, 4, 0, 0, 0, 0, time.UTC),
			name: "清明节",
		},
		{
			// Dragon Boat Festival
			date: time.Date(2024, 6, 10, 0, 0, 0, 0, time.UTC),
			name: "端午节",
		},
		{
			// Mid-Autumn Festival
			date: time.Date(2024, 9, 17, 0, 0, 0, 0, time.UTC),
			name: "中秋节",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			holiday, exists := holidays[tc.date]
			if !exists {
				t.Errorf("Expected lunar holiday %s on %s to exist", tc.name, tc.date.Format("2006-01-02"))
				return
			}

			if holiday.Name != tc.name {
				t.Errorf("Expected holiday name %s, got %s", tc.name, holiday.Name)
			}
		})
	}
}

func TestCNProvider_SpringFestivalDates(t *testing.T) {
	provider := NewCNProvider()

	testCases := []struct {
		year          int
		expectedDates []time.Time
	}{
		{
			year: 2024,
			expectedDates: []time.Time{
				time.Date(2024, 2, 9, 0, 0, 0, 0, time.UTC),  // Eve
				time.Date(2024, 2, 10, 0, 0, 0, 0, time.UTC), // Day 1
				time.Date(2024, 2, 11, 0, 0, 0, 0, time.UTC), // Day 2
				time.Date(2024, 2, 12, 0, 0, 0, 0, time.UTC), // Day 3
				time.Date(2024, 2, 13, 0, 0, 0, 0, time.UTC), // Day 4
				time.Date(2024, 2, 14, 0, 0, 0, 0, time.UTC), // Day 5
				time.Date(2024, 2, 15, 0, 0, 0, 0, time.UTC), // Day 6
			},
		},
		{
			year: 2025,
			expectedDates: []time.Time{
				time.Date(2025, 1, 28, 0, 0, 0, 0, time.UTC), // Eve
				time.Date(2025, 1, 29, 0, 0, 0, 0, time.UTC), // Day 1
				time.Date(2025, 1, 30, 0, 0, 0, 0, time.UTC), // Day 2
				time.Date(2025, 1, 31, 0, 0, 0, 0, time.UTC), // Day 3
				time.Date(2025, 2, 1, 0, 0, 0, 0, time.UTC),  // Day 4
				time.Date(2025, 2, 2, 0, 0, 0, 0, time.UTC),  // Day 5
				time.Date(2025, 2, 3, 0, 0, 0, 0, time.UTC),  // Day 6
			},
		},
	}

	for _, tc := range testCases {
		t.Run(fmt.Sprintf("year_%d", tc.year), func(t *testing.T) {
			dates := provider.getSpringFestivalDates(tc.year)

			if len(dates) != len(tc.expectedDates) {
				t.Errorf("Expected %d Spring Festival dates, got %d", len(tc.expectedDates), len(dates))
				return
			}

			for i, expectedDate := range tc.expectedDates {
				if !dates[i].Equal(expectedDate) {
					t.Errorf("Expected Spring Festival date %s, got %s", expectedDate.Format("2006-01-02"), dates[i].Format("2006-01-02"))
				}
			}
		})
	}
}

func TestCNProvider_MultilingualSupport(t *testing.T) {
	provider := NewCNProvider()
	year := 2024

	holidays := provider.LoadHolidays(year)

	// Test National Day multilingual support
	nationalDay := time.Date(2024, 10, 1, 0, 0, 0, 0, time.UTC)
	holiday, exists := holidays[nationalDay]
	if !exists {
		t.Fatal("Expected National Day to exist")
	}

	expectedLanguages := map[string]string{
		"zh": "国庆节",
		"en": "National Day",
	}

	for lang, expectedName := range expectedLanguages {
		if name, exists := holiday.Languages[lang]; !exists {
			t.Errorf("Expected language %s to exist for National Day", lang)
		} else if name != expectedName {
			t.Errorf("Expected %s name %s, got %s", lang, expectedName, name)
		}
	}
}

func TestCNProvider_GetRegionalHolidays(t *testing.T) {
	provider := NewCNProvider()
	year := 2024

	// Test Hong Kong SAR Establishment Day
	regionalHolidays := provider.GetRegionalHolidays(year, []string{"91"})

	hkSARDay := time.Date(2024, 7, 1, 0, 0, 0, 0, time.UTC)
	holiday, exists := regionalHolidays[hkSARDay]
	if !exists {
		t.Errorf("Expected Hong Kong SAR Establishment Day on %s to exist", hkSARDay.Format("2006-01-02"))
		return
	}

	if holiday.Name != "香港特别行政区成立纪念日" {
		t.Errorf("Expected holiday name 香港特别行政区成立纪念日, got %s", holiday.Name)
	}

	if holiday.Category != "regional" {
		t.Errorf("Expected category regional, got %s", holiday.Category)
	}

	// Test Macau SAR Establishment Day
	regionalHolidays = provider.GetRegionalHolidays(year, []string{"92"})

	macauSARDay := time.Date(2024, 12, 20, 0, 0, 0, 0, time.UTC)
	holiday, exists = regionalHolidays[macauSARDay]
	if !exists {
		t.Errorf("Expected Macau SAR Establishment Day on %s to exist", macauSARDay.Format("2006-01-02"))
		return
	}

	if holiday.Name != "澳门特别行政区成立纪念日" {
		t.Errorf("Expected holiday name 澳门特别行政区成立纪念日, got %s", holiday.Name)
	}
}

func TestCNProvider_GetSpecialObservances(t *testing.T) {
	provider := NewCNProvider()
	year := 2024

	observances := provider.GetSpecialObservances(year)

	// Test Teachers' Day
	teachersDay := time.Date(2024, 9, 10, 0, 0, 0, 0, time.UTC)
	holiday, exists := observances[teachersDay]
	if !exists {
		t.Errorf("Expected Teachers' Day on %s to exist", teachersDay.Format("2006-01-02"))
		return
	}

	if holiday.Name != "教师节" {
		t.Errorf("Expected holiday name 教师节, got %s", holiday.Name)
	}

	// Test Constitution Day
	constitutionDay := time.Date(2024, 12, 4, 0, 0, 0, 0, time.UTC)
	holiday, exists = observances[constitutionDay]
	if !exists {
		t.Errorf("Expected Constitution Day on %s to exist", constitutionDay.Format("2006-01-02"))
		return
	}

	if holiday.Name != "宪法日" {
		t.Errorf("Expected holiday name 宪法日, got %s", holiday.Name)
	}
}

func TestCNProvider_HolidayCount(t *testing.T) {
	provider := NewCNProvider()
	year := 2024

	holidays := provider.LoadHolidays(year)

	// China should have 19 main holidays:
	// - 6 fixed public holidays
	// - 7 Spring Festival days (Eve + 6 days)
	// - 3 National Day holidays
	// - 3 other traditional lunar holidays (Qingming, Dragon Boat, Mid-Autumn)
	expectedCount := 19
	if len(holidays) != expectedCount {
		t.Errorf("Expected %d holidays for China in %d, got %d", expectedCount, year, len(holidays))
	}
}

// Benchmark tests
func BenchmarkCNProvider_LoadHolidays(b *testing.B) {
	provider := NewCNProvider()
	year := 2024

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		provider.LoadHolidays(year)
	}
}


