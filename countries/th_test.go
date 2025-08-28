package countries

import (
	"testing"
	"time"
)

func TestTHProvider(t *testing.T) {
	provider := NewTHProvider()

	// Test basic provider properties
	if provider.GetCountryCode() != "TH" {
		t.Errorf("Expected country code 'TH', got '%s'", provider.GetCountryCode())
	}

	// Test subdivisions (77 provinces + 1 special administrative area)
	subdivisions := provider.GetSupportedSubdivisions()
	if len(subdivisions) != 77 {
		t.Errorf("Expected 77 subdivisions for Thailand, got %d", len(subdivisions))
	}

	// Test categories
	categories := provider.GetSupportedCategories()
	expectedCategories := []string{"national", "religious", "royal", "buddhist", "cultural"}
	if len(categories) != len(expectedCategories) {
		t.Errorf("Expected %d categories, got %d", len(expectedCategories), len(categories))
	}
}

func TestTHHolidays2024(t *testing.T) {
	provider := NewTHProvider()
	holidays := provider.LoadHolidays(2024)

	// Test some key Thai holidays for 2024
	testCases := []struct {
		date     time.Time
		name     string
		category string
	}{
		{time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC), "วันขึ้นปีใหม่", "national"},
		{time.Date(2024, 2, 24, 0, 0, 0, 0, time.UTC), "วันมาฆบูชา", "buddhist"},                            // Magha Puja Day 2024
		{time.Date(2024, 4, 6, 0, 0, 0, 0, time.UTC), "วันจักรี", "royal"},                                  // Chakri Day
		{time.Date(2024, 4, 13, 0, 0, 0, 0, time.UTC), "วันสงกรานต์", "cultural"},                           // Songkran Festival
		{time.Date(2024, 4, 14, 0, 0, 0, 0, time.UTC), "วันสงกรานต์ (วันที่ 2)", "cultural"},                // Songkran Day 2
		{time.Date(2024, 4, 15, 0, 0, 0, 0, time.UTC), "วันสงกรานต์ (วันที่ 3)", "cultural"},                // Songkran Day 3
		{time.Date(2024, 5, 1, 0, 0, 0, 0, time.UTC), "วันแรงงานแห่งชาติ", "national"},                      // Labour Day
		{time.Date(2024, 5, 4, 0, 0, 0, 0, time.UTC), "วันฉัตรมงคล", "royal"},                               // Coronation Day
		{time.Date(2024, 5, 9, 0, 0, 0, 0, time.UTC), "วันพืชมงคล", "royal"},                                // Royal Ploughing Ceremony 2024
		{time.Date(2024, 5, 22, 0, 0, 0, 0, time.UTC), "วันวิสาขบูชา", "buddhist"},                          // Visakha Puja Day 2024
		{time.Date(2024, 6, 3, 0, 0, 0, 0, time.UTC), "วันเฉลิมพระชนมพรรษาสมเด็จพระนางเจ้าสุทิดา", "royal"}, // Queen Suthida's Birthday
		{time.Date(2024, 7, 21, 0, 0, 0, 0, time.UTC), "วันอาสาฬหบูชา", "buddhist"},                         // Asalha Puja Day 2024
		{time.Date(2024, 7, 22, 0, 0, 0, 0, time.UTC), "วันเข้าพรรษา", "buddhist"},                          // Khao Phansa 2024
		{time.Date(2024, 7, 28, 0, 0, 0, 0, time.UTC), "วันเฉลิมพระชนมพรรษาพระบาทสมเด็จพระปรเมนทรรามาธิบดีศรีสินทรมหาวชิราลงกรณ พระวชิรเกล้าเจ้าอยู่หัว", "royal"}, // King's Birthday
		{time.Date(2024, 8, 12, 0, 0, 0, 0, time.UTC), "วันเฉลิมพระชนมพรรษาสมเด็จพระนางเจ้าสิริกิติ์ พระบรมราชินีนาถ พระบรมราชชนนีพันปีหลวง", "royal"},             // Queen Mother's Birthday
		{time.Date(2024, 10, 13, 0, 0, 0, 0, time.UTC), "วันคล้ายวันสวรรคตพระบาทสมเด็จพระปรมินทรมหาภูมิพลอดุลยเดช บรมนาถบพิตร", "royal"},                           // King Bhumibol Memorial Day
		{time.Date(2024, 10, 23, 0, 0, 0, 0, time.UTC), "วันปิยมหาราช", "royal"}, // Chulalongkorn Day
		{time.Date(2024, 12, 5, 0, 0, 0, 0, time.UTC), "วันเฉลิมพระชนมพรรษาพระบาทสมเด็จพระปรมินทรมหาภูมิพลอดุลยเดช บรมนาถบพิตร", "royal"}, // King Bhumibol's Birthday
		{time.Date(2024, 12, 10, 0, 0, 0, 0, time.UTC), "วันรัฐธรรมนูญ", "national"}, // Constitution Day
		{time.Date(2024, 12, 31, 0, 0, 0, 0, time.UTC), "วันสิ้นปี", "national"},     // New Year's Eve
	}

	for _, tc := range testCases {
		holiday, exists := holidays[tc.date]
		if !exists {
			t.Errorf("Expected holiday on %s, but none found", tc.date.Format("2006-01-02"))
			continue
		}

		if holiday.Name != tc.name {
			t.Errorf("Expected holiday name '%s' on %s, got '%s'", tc.name, tc.date.Format("2006-01-02"), holiday.Name)
		}

		if holiday.Category != tc.category {
			t.Errorf("Expected category '%s' for %s, got '%s'", tc.category, tc.name, holiday.Category)
		}
	}

	// Check that we have a reasonable number of holidays
	if len(holidays) < 18 {
		t.Errorf("Expected at least 18 holidays for Thailand in 2024, got %d", len(holidays))
	}
}

func TestTHLanguageSupport(t *testing.T) {
	provider := NewTHProvider()
	holidays := provider.LoadHolidays(2024)

	// Test New Year's Day in Thai and English
	newYear := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	holiday, exists := holidays[newYear]
	if !exists {
		t.Fatal("New Year's Day not found")
	}

	// Check Thai translation
	if thai, ok := holiday.Languages["th"]; !ok || thai != "วันขึ้นปีใหม่" {
		t.Errorf("Expected Thai translation 'วันขึ้นปีใหม่', got '%s'", thai)
	}

	// Check English translation
	if english, ok := holiday.Languages["en"]; !ok || english != "New Year's Day" {
		t.Errorf("Expected English translation 'New Year's Day', got '%s'", english)
	}
}

func TestTHBuddhistHolidays(t *testing.T) {
	provider := NewTHProvider()
	holidays := provider.LoadHolidays(2024)

	// Count Buddhist holidays
	buddhistCount := 0
	buddhistHolidays := []string{}
	for _, holiday := range holidays {
		if holiday.Category == "buddhist" {
			buddhistCount++
			buddhistHolidays = append(buddhistHolidays, holiday.Name)
		}
	}

	// Should have Magha Puja, Visakha Puja, Asalha Puja, and Khao Phansa = 4 Buddhist holidays
	if buddhistCount != 4 {
		t.Errorf("Expected 4 Buddhist holidays for Thailand in 2024, got %d: %v", buddhistCount, buddhistHolidays)
	}
}

func TestTHRoyalHolidays(t *testing.T) {
	provider := NewTHProvider()
	holidays := provider.LoadHolidays(2024)

	// Count royal holidays
	royalCount := 0
	royalHolidays := []string{}
	for _, holiday := range holidays {
		if holiday.Category == "royal" {
			royalCount++
			royalHolidays = append(royalHolidays, holiday.Name)
		}
	}

	// Should have multiple royal holidays including Chakri Day, Coronation Day, Royal Ploughing, King/Queen birthdays, etc.
	if royalCount < 7 {
		t.Errorf("Expected at least 7 royal holidays for Thailand in 2024, got %d: %v", royalCount, royalHolidays)
	}
}

func TestTHSongkranFestival(t *testing.T) {
	provider := NewTHProvider()
	holidays := provider.LoadHolidays(2024)

	// Test Songkran Festival (3 days: April 13-15)
	songkranDates := []time.Time{
		time.Date(2024, 4, 13, 0, 0, 0, 0, time.UTC),
		time.Date(2024, 4, 14, 0, 0, 0, 0, time.UTC),
		time.Date(2024, 4, 15, 0, 0, 0, 0, time.UTC),
	}

	for i, date := range songkranDates {
		holiday, exists := holidays[date]
		if !exists {
			t.Errorf("Songkran day %d not found on %s", i+1, date.Format("2006-01-02"))
			continue
		}

		if holiday.Category != "cultural" {
			t.Errorf("Expected Songkran to be cultural category, got '%s'", holiday.Category)
		}

		// Check that it contains "สงกรานต์" in the name
		if len(holiday.Name) == 0 || holiday.Name[:len("วันสงกรานต์")] != "วันสงกรานต์" {
			t.Errorf("Expected Songkran holiday name to start with 'วันสงกรานต์', got '%s'", holiday.Name)
		}
	}
}

func TestTHMaghaPujaCalculation(t *testing.T) {
	provider := NewTHProvider()

	// Test Magha Puja dates for known years
	testCases := []struct {
		year  int
		month time.Month
		day   int
	}{
		{2024, time.February, 24}, // Magha Puja 2024
		{2025, time.February, 12}, // Magha Puja 2025
		{2026, time.March, 4},     // Magha Puja 2026
		{2027, time.February, 21}, // Magha Puja 2027
	}

	for _, tc := range testCases {
		maghaPuja := provider.calculateMaghaPuja(tc.year)
		expected := time.Date(tc.year, tc.month, tc.day, 0, 0, 0, 0, time.UTC)

		if !maghaPuja.Equal(expected) {
			t.Errorf("Expected Magha Puja %d to be %s, got %s", tc.year, expected.Format("2006-01-02"), maghaPuja.Format("2006-01-02"))
		}
	}
}

func TestTHVisakhaPujaCalculation(t *testing.T) {
	provider := NewTHProvider()

	// Test Visakha Puja dates for known years
	testCases := []struct {
		year  int
		month time.Month
		day   int
	}{
		{2024, time.May, 22}, // Visakha Puja 2024
		{2025, time.May, 12}, // Visakha Puja 2025
		{2026, time.May, 31}, // Visakha Puja 2026
		{2027, time.May, 21}, // Visakha Puja 2027
	}

	for _, tc := range testCases {
		visakhaPuja := provider.calculateVisakhaPuja(tc.year)
		expected := time.Date(tc.year, tc.month, tc.day, 0, 0, 0, 0, time.UTC)

		if !visakhaPuja.Equal(expected) {
			t.Errorf("Expected Visakha Puja %d to be %s, got %s", tc.year, expected.Format("2006-01-02"), visakhaPuja.Format("2006-01-02"))
		}
	}
}

func TestTHAsalhaPujaCalculation(t *testing.T) {
	provider := NewTHProvider()

	// Test Asalha Puja dates for known years
	testCases := []struct {
		year  int
		month time.Month
		day   int
	}{
		{2024, time.July, 21}, // Asalha Puja 2024
		{2025, time.July, 11}, // Asalha Puja 2025
		{2026, time.July, 30}, // Asalha Puja 2026
		{2027, time.July, 19}, // Asalha Puja 2027
	}

	for _, tc := range testCases {
		asalhaPuja := provider.calculateAsalhaPuja(tc.year)
		expected := time.Date(tc.year, tc.month, tc.day, 0, 0, 0, 0, time.UTC)

		if !asalhaPuja.Equal(expected) {
			t.Errorf("Expected Asalha Puja %d to be %s, got %s", tc.year, expected.Format("2006-01-02"), asalhaPuja.Format("2006-01-02"))
		}
	}
}

func TestTHKhaoPhansamFollowsAsalhaPuja(t *testing.T) {
	provider := NewTHProvider()
	holidays := provider.LoadHolidays(2024)

	// Find Asalha Puja and Khao Phansa
	asalhaPujaDate := time.Date(2024, 7, 21, 0, 0, 0, 0, time.UTC)
	khaoPhansamDate := time.Date(2024, 7, 22, 0, 0, 0, 0, time.UTC)

	asalhaPuja, asalhaExists := holidays[asalhaPujaDate]
	khaoPhansal, khaoExists := holidays[khaoPhansamDate]

	if !asalhaExists {
		t.Fatal("Asalha Puja not found")
	}
	if !khaoExists {
		t.Fatal("Khao Phansa not found")
	}

	// Verify Khao Phansa is the day after Asalha Puja
	expectedKhaoDate := asalhaPujaDate.AddDate(0, 0, 1)
	if !khaoPhansamDate.Equal(expectedKhaoDate) {
		t.Errorf("Expected Khao Phansa to be the day after Asalha Puja (%s), got %s", expectedKhaoDate.Format("2006-01-02"), khaoPhansamDate.Format("2006-01-02"))
	}

	// Verify both are Buddhist holidays
	if asalhaPuja.Category != "buddhist" {
		t.Errorf("Expected Asalha Puja to be Buddhist category, got '%s'", asalhaPuja.Category)
	}
	if khaoPhansal.Category != "buddhist" {
		t.Errorf("Expected Khao Phansa to be Buddhist category, got '%s'", khaoPhansal.Category)
	}
}

func BenchmarkTHProvider(b *testing.B) {
	provider := NewTHProvider()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = provider.LoadHolidays(2024)
	}
}

func BenchmarkTHBuddhistHolidayCalculations(b *testing.B) {
	provider := NewTHProvider()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = provider.calculateMaghaPuja(2024)
		_ = provider.calculateVisakhaPuja(2024)
		_ = provider.calculateAsalhaPuja(2024)
	}
}
