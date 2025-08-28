package countries

import (
	"testing"
	"time"
)

func TestIDProvider(t *testing.T) {
	provider := NewIDProvider()

	// Test basic provider properties
	if provider.GetCountryCode() != "ID" {
		t.Errorf("Expected country code 'ID', got '%s'", provider.GetCountryCode())
	}

	// Test subdivisions (38 provinces)
	subdivisions := provider.GetSupportedSubdivisions()
	if len(subdivisions) != 38 {
		t.Errorf("Expected 38 subdivisions for Indonesia, got %d", len(subdivisions))
	}

	// Test categories
	categories := provider.GetSupportedCategories()
	expectedCategories := []string{"national", "religious", "islamic", "christian", "buddhist", "hindu", "chinese"}
	if len(categories) != len(expectedCategories) {
		t.Errorf("Expected %d categories, got %d", len(expectedCategories), len(categories))
	}
}

func TestIDHolidays2024(t *testing.T) {
	provider := NewIDProvider()
	holidays := provider.LoadHolidays(2024)

	// Test some key Indonesian holidays for 2024
	testCases := []struct {
		date     time.Time
		name     string
		category string
	}{
		{time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC), "Tahun Baru Masehi", "national"},
		{time.Date(2024, 2, 10, 0, 0, 0, 0, time.UTC), "Tahun Baru Imlek", "chinese"},
		{time.Date(2024, 3, 11, 0, 0, 0, 0, time.UTC), "Hari Raya Nyepi", "hindu"},
		{time.Date(2024, 3, 29, 0, 0, 0, 0, time.UTC), "Wafat Isa Al Masih", "christian"}, // Good Friday 2024
		{time.Date(2024, 4, 10, 0, 0, 0, 0, time.UTC), "Hari Raya Idul Fitri", "islamic"},
		{time.Date(2024, 5, 1, 0, 0, 0, 0, time.UTC), "Hari Buruh Internasional", "national"},
		{time.Date(2024, 5, 9, 0, 0, 0, 0, time.UTC), "Kenaikan Isa Al Masih", "christian"}, // Ascension Day 2024
		{time.Date(2024, 5, 23, 0, 0, 0, 0, time.UTC), "Hari Raya Waisak", "buddhist"},
		{time.Date(2024, 6, 1, 0, 0, 0, 0, time.UTC), "Hari Lahir Pancasila", "national"},
		{time.Date(2024, 6, 16, 0, 0, 0, 0, time.UTC), "Hari Raya Idul Adha", "islamic"},
		{time.Date(2024, 8, 17, 0, 0, 0, 0, time.UTC), "Hari Kemerdekaan Republik Indonesia", "national"},
		{time.Date(2024, 11, 10, 0, 0, 0, 0, time.UTC), "Hari Pahlawan", "national"},
		{time.Date(2024, 12, 25, 0, 0, 0, 0, time.UTC), "Hari Raya Natal", "christian"},
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

	// Check that we have a reasonable number of holidays (should be diverse)
	if len(holidays) < 15 {
		t.Errorf("Expected at least 15 holidays for Indonesia in 2024, got %d", len(holidays))
	}
}

func TestIDLanguageSupport(t *testing.T) {
	provider := NewIDProvider()
	holidays := provider.LoadHolidays(2024)

	// Test Independence Day in Indonesian and English
	independenceDay := time.Date(2024, 8, 17, 0, 0, 0, 0, time.UTC)
	holiday, exists := holidays[independenceDay]
	if !exists {
		t.Fatal("Independence Day not found")
	}

	// Check Indonesian translation
	if indonesian, ok := holiday.Languages["id"]; !ok || indonesian != "Hari Kemerdekaan Republik Indonesia" {
		t.Errorf("Expected Indonesian translation 'Hari Kemerdekaan Republik Indonesia', got '%s'", indonesian)
	}

	// Check English translation
	if english, ok := holiday.Languages["en"]; !ok || english != "Independence Day" {
		t.Errorf("Expected English translation 'Independence Day', got '%s'", english)
	}
}

func TestIDIndependenceDay(t *testing.T) {
	provider := NewIDProvider()
	holidays := provider.LoadHolidays(2024)

	// Test Independence Day (August 17) - Indonesia's National Day
	independenceDay := time.Date(2024, 8, 17, 0, 0, 0, 0, time.UTC)
	holiday, exists := holidays[independenceDay]
	if !exists {
		t.Fatal("Independence Day not found")
	}

	if holiday.Name != "Hari Kemerdekaan Republik Indonesia" {
		t.Errorf("Expected Independence Day name 'Hari Kemerdekaan Republik Indonesia', got '%s'", holiday.Name)
	}

	if holiday.Category != "national" {
		t.Errorf("Expected Independence Day to be national category, got '%s'", holiday.Category)
	}
}

func TestIDPancasilaDay(t *testing.T) {
	provider := NewIDProvider()
	holidays := provider.LoadHolidays(2024)

	// Test Pancasila Day (June 1)
	pancasilaDay := time.Date(2024, 6, 1, 0, 0, 0, 0, time.UTC)
	holiday, exists := holidays[pancasilaDay]
	if !exists {
		t.Fatal("Pancasila Day not found")
	}

	if holiday.Name != "Hari Lahir Pancasila" {
		t.Errorf("Expected Pancasila Day name 'Hari Lahir Pancasila', got '%s'", holiday.Name)
	}

	if holiday.Category != "national" {
		t.Errorf("Expected Pancasila Day to be national category, got '%s'", holiday.Category)
	}
}

func TestIDIslamicHolidays(t *testing.T) {
	provider := NewIDProvider()
	holidays := provider.LoadHolidays(2024)

	// Count Islamic holidays
	islamicCount := 0
	islamicNames := []string{}
	for _, holiday := range holidays {
		if holiday.Category == "islamic" {
			islamicCount++
			islamicNames = append(islamicNames, holiday.Name)
		}
	}

	// Should have multiple Islamic holidays (Idul Fitri 2 days + others)
	if islamicCount < 4 {
		t.Errorf("Expected at least 4 Islamic holidays for Indonesia in 2024, got %d: %v", islamicCount, islamicNames)
	}

	// Test Idul Fitri (should have 2 days)
	iduFitriFirst := time.Date(2024, 4, 10, 0, 0, 0, 0, time.UTC)
	if _, exists := holidays[iduFitriFirst]; !exists {
		t.Error("Idul Fitri first day should exist")
	}

	iduFitriSecond := time.Date(2024, 4, 11, 0, 0, 0, 0, time.UTC)
	if _, exists := holidays[iduFitriSecond]; !exists {
		t.Error("Idul Fitri second day should exist")
	}
}

func TestIDChristianHolidays(t *testing.T) {
	provider := NewIDProvider()
	holidays := provider.LoadHolidays(2024)

	// Test Christian holidays
	christianHolidays := []struct {
		date time.Time
		name string
	}{
		{time.Date(2024, 3, 29, 0, 0, 0, 0, time.UTC), "Wafat Isa Al Masih"},   // Good Friday
		{time.Date(2024, 5, 9, 0, 0, 0, 0, time.UTC), "Kenaikan Isa Al Masih"}, // Ascension Day
		{time.Date(2024, 12, 25, 0, 0, 0, 0, time.UTC), "Hari Raya Natal"},     // Christmas
	}

	for _, tc := range christianHolidays {
		holiday, exists := holidays[tc.date]
		if !exists {
			t.Errorf("Christian holiday %s not found on %s", tc.name, tc.date.Format("2006-01-02"))
			continue
		}

		if holiday.Name != tc.name {
			t.Errorf("Expected Christian holiday name '%s', got '%s'", tc.name, holiday.Name)
		}

		if holiday.Category != "christian" {
			t.Errorf("Expected Christian holiday category, got '%s'", holiday.Category)
		}
	}
}

func TestIDMultiReligiousHolidays(t *testing.T) {
	provider := NewIDProvider()
	holidays := provider.LoadHolidays(2024)

	// Test that Indonesia has holidays from multiple religions
	religiousCategories := map[string]int{
		"islamic":   0,
		"christian": 0,
		"buddhist":  0,
		"hindu":     0,
		"chinese":   0,
	}

	for _, holiday := range holidays {
		if count, exists := religiousCategories[holiday.Category]; exists {
			religiousCategories[holiday.Category] = count + 1
		}
	}

	// Should have holidays from each religion
	for religion, count := range religiousCategories {
		if count == 0 {
			t.Errorf("Expected at least one %s holiday in Indonesia", religion)
		}
	}

	// Islamic should have the most holidays (largest population)
	if religiousCategories["islamic"] < religiousCategories["christian"] {
		t.Error("Expected more Islamic holidays than Christian holidays in Indonesia")
	}
}

func TestIDChineseNewYear(t *testing.T) {
	provider := NewIDProvider()
	holidays := provider.LoadHolidays(2024)

	// Test Chinese New Year (February 10, 2024)
	chineseNewYear := time.Date(2024, 2, 10, 0, 0, 0, 0, time.UTC)
	holiday, exists := holidays[chineseNewYear]
	if !exists {
		t.Fatal("Chinese New Year not found")
	}

	if holiday.Name != "Tahun Baru Imlek" {
		t.Errorf("Expected Chinese New Year name 'Tahun Baru Imlek', got '%s'", holiday.Name)
	}

	if holiday.Category != "chinese" {
		t.Errorf("Expected Chinese New Year to be chinese category, got '%s'", holiday.Category)
	}
}

func TestIDVesakDay(t *testing.T) {
	provider := NewIDProvider()
	holidays := provider.LoadHolidays(2024)

	// Test Vesak Day (May 23, 2024)
	vesakDay := time.Date(2024, 5, 23, 0, 0, 0, 0, time.UTC)
	holiday, exists := holidays[vesakDay]
	if !exists {
		t.Fatal("Vesak Day not found")
	}

	if holiday.Name != "Hari Raya Waisak" {
		t.Errorf("Expected Vesak Day name 'Hari Raya Waisak', got '%s'", holiday.Name)
	}

	if holiday.Category != "buddhist" {
		t.Errorf("Expected Vesak Day to be buddhist category, got '%s'", holiday.Category)
	}
}

func TestIDNyepi(t *testing.T) {
	provider := NewIDProvider()
	holidays := provider.LoadHolidays(2024)

	// Test Nyepi (March 11, 2024)
	nyepi := time.Date(2024, 3, 11, 0, 0, 0, 0, time.UTC)
	holiday, exists := holidays[nyepi]
	if !exists {
		t.Fatal("Nyepi not found")
	}

	if holiday.Name != "Hari Raya Nyepi" {
		t.Errorf("Expected Nyepi name 'Hari Raya Nyepi', got '%s'", holiday.Name)
	}

	if holiday.Category != "hindu" {
		t.Errorf("Expected Nyepi to be hindu category, got '%s'", holiday.Category)
	}
}

func TestIDHeroesDay(t *testing.T) {
	provider := NewIDProvider()
	holidays := provider.LoadHolidays(2024)

	// Test Heroes' Day (November 10)
	heroesDay := time.Date(2024, 11, 10, 0, 0, 0, 0, time.UTC)
	holiday, exists := holidays[heroesDay]
	if !exists {
		t.Fatal("Heroes' Day not found")
	}

	if holiday.Name != "Hari Pahlawan" {
		t.Errorf("Expected Heroes' Day name 'Hari Pahlawan', got '%s'", holiday.Name)
	}

	if holiday.Category != "national" {
		t.Errorf("Expected Heroes' Day to be national category, got '%s'", holiday.Category)
	}
}

func BenchmarkIDProvider(b *testing.B) {
	provider := NewIDProvider()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = provider.LoadHolidays(2024)
	}
}

func BenchmarkIDMultiReligiousCalculation(b *testing.B) {
	provider := NewIDProvider()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		holidays := make(map[time.Time]*Holiday)
		provider.addIslamicHolidays(holidays, 2024)
		provider.addChristianHolidays(holidays, 2024)
		provider.addOtherReligiousHolidays(holidays, 2024)
	}
}
