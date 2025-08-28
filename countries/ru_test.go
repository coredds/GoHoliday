package countries

import (
	"testing"
	"time"
)

func TestRUProvider(t *testing.T) {
	provider := NewRUProvider()

	// Test basic provider properties
	if provider.GetCountryCode() != "RU" {
		t.Errorf("Expected country code 'RU', got '%s'", provider.GetCountryCode())
	}

	// Test subdivisions (85 federal subjects)
	subdivisions := provider.GetSupportedSubdivisions()
	if len(subdivisions) != 85 {
		t.Errorf("Expected 85 subdivisions for Russia, got %d", len(subdivisions))
	}

	// Test categories
	categories := provider.GetSupportedCategories()
	expectedCategories := []string{"national", "religious", "commemorative", "orthodox"}
	if len(categories) != len(expectedCategories) {
		t.Errorf("Expected %d categories, got %d", len(expectedCategories), len(categories))
	}
}

func TestRUHolidays2024(t *testing.T) {
	provider := NewRUProvider()
	holidays := provider.LoadHolidays(2024)

	// Test some key Russian holidays for 2024
	testCases := []struct {
		date     time.Time
		name     string
		category string
	}{
		{time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC), "Новый год", "national"},
		{time.Date(2024, 1, 7, 0, 0, 0, 0, time.UTC), "Рождество Христово", "orthodox"},
		{time.Date(2024, 2, 23, 0, 0, 0, 0, time.UTC), "День защитника Отечества", "national"},
		{time.Date(2024, 3, 8, 0, 0, 0, 0, time.UTC), "Международный женский день", "national"},
		{time.Date(2024, 5, 1, 0, 0, 0, 0, time.UTC), "Праздник Весны и Труда", "national"},
		{time.Date(2024, 5, 5, 0, 0, 0, 0, time.UTC), "Пасха", "orthodox"}, // Orthodox Easter 2024
		{time.Date(2024, 5, 9, 0, 0, 0, 0, time.UTC), "День Победы", "commemorative"},
		{time.Date(2024, 6, 12, 0, 0, 0, 0, time.UTC), "День России", "national"},
		{time.Date(2024, 11, 4, 0, 0, 0, 0, time.UTC), "День народного единства", "national"},
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
	if len(holidays) < 15 {
		t.Errorf("Expected at least 15 holidays for Russia in 2024, got %d", len(holidays))
	}
}

func TestRULanguageSupport(t *testing.T) {
	provider := NewRUProvider()
	holidays := provider.LoadHolidays(2024)

	// Test Victory Day in Russian and English
	victoryDay := time.Date(2024, 5, 9, 0, 0, 0, 0, time.UTC)
	holiday, exists := holidays[victoryDay]
	if !exists {
		t.Fatal("Victory Day not found")
	}

	// Check Russian translation
	if russian, ok := holiday.Languages["ru"]; !ok || russian != "День Победы" {
		t.Errorf("Expected Russian translation 'День Победы', got '%s'", russian)
	}

	// Check English translation
	if english, ok := holiday.Languages["en"]; !ok || english != "Victory Day" {
		t.Errorf("Expected English translation 'Victory Day', got '%s'", english)
	}
}

func TestRUNewYearHolidays(t *testing.T) {
	provider := NewRUProvider()
	holidays := provider.LoadHolidays(2024)

	// Test New Year holidays (January 1-8)
	// Note: January 7 is also Orthodox Christmas, so we count both as separate holidays
	newYearHolidaysCount := 0
	for i := 1; i <= 8; i++ {
		date := time.Date(2024, 1, i, 0, 0, 0, 0, time.UTC)
		if holiday, exists := holidays[date]; exists {
			if holiday.Category == "national" || (i == 7 && holiday.Category == "orthodox") {
				newYearHolidaysCount++
			}
		}
	}

	if newYearHolidaysCount != 8 {
		t.Errorf("Expected 8 holidays in New Year period (including Orthodox Christmas), got %d", newYearHolidaysCount)
	}

	// Check New Year's Day specifically
	newYearDay := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	holiday, exists := holidays[newYearDay]
	if !exists {
		t.Fatal("New Year's Day not found")
	}

	if holiday.Name != "Новый год" {
		t.Errorf("Expected New Year's Day name 'Новый год', got '%s'", holiday.Name)
	}
}

func TestRUOrthodoxChristmas(t *testing.T) {
	provider := NewRUProvider()
	holidays := provider.LoadHolidays(2024)

	// Test Orthodox Christmas (January 7)
	christmas := time.Date(2024, 1, 7, 0, 0, 0, 0, time.UTC)
	holiday, exists := holidays[christmas]
	if !exists {
		t.Fatal("Orthodox Christmas not found")
	}

	if holiday.Name != "Рождество Христово" {
		t.Errorf("Expected Orthodox Christmas name 'Рождество Христово', got '%s'", holiday.Name)
	}

	if holiday.Category != "orthodox" {
		t.Errorf("Expected Orthodox Christmas to be orthodox category, got '%s'", holiday.Category)
	}
}

func TestRUOrthodoxEaster(t *testing.T) {
	provider := NewRUProvider()
	holidays := provider.LoadHolidays(2024)

	// Test Orthodox Easter (May 5, 2024)
	easter := time.Date(2024, 5, 5, 0, 0, 0, 0, time.UTC)
	holiday, exists := holidays[easter]
	if !exists {
		t.Fatal("Orthodox Easter not found")
	}

	if holiday.Name != "Пасха" {
		t.Errorf("Expected Orthodox Easter name 'Пасха', got '%s'", holiday.Name)
	}

	if holiday.Category != "orthodox" {
		t.Errorf("Expected Orthodox Easter to be orthodox category, got '%s'", holiday.Category)
	}

	// Test related Orthodox holidays
	palmSunday := time.Date(2024, 4, 28, 0, 0, 0, 0, time.UTC) // 1 week before Easter
	if _, exists := holidays[palmSunday]; !exists {
		t.Error("Orthodox Palm Sunday should exist")
	}

	easterMonday := time.Date(2024, 5, 6, 0, 0, 0, 0, time.UTC) // Day after Easter
	if _, exists := holidays[easterMonday]; !exists {
		t.Error("Orthodox Easter Monday should exist")
	}
}

func TestRUVictoryDay(t *testing.T) {
	provider := NewRUProvider()
	holidays := provider.LoadHolidays(2024)

	// Test Victory Day (May 9) - most important commemorative holiday
	victoryDay := time.Date(2024, 5, 9, 0, 0, 0, 0, time.UTC)
	holiday, exists := holidays[victoryDay]
	if !exists {
		t.Fatal("Victory Day not found")
	}

	if holiday.Name != "День Победы" {
		t.Errorf("Expected Victory Day name 'День Победы', got '%s'", holiday.Name)
	}

	if holiday.Category != "commemorative" {
		t.Errorf("Expected Victory Day to be commemorative category, got '%s'", holiday.Category)
	}
}

func TestRUDefenderDay(t *testing.T) {
	provider := NewRUProvider()
	holidays := provider.LoadHolidays(2024)

	// Test Defender of the Fatherland Day (February 23)
	defenderDay := time.Date(2024, 2, 23, 0, 0, 0, 0, time.UTC)
	holiday, exists := holidays[defenderDay]
	if !exists {
		t.Fatal("Defender of the Fatherland Day not found")
	}

	if holiday.Name != "День защитника Отечества" {
		t.Errorf("Expected Defender Day name 'День защитника Отечества', got '%s'", holiday.Name)
	}

	if holiday.Category != "national" {
		t.Errorf("Expected Defender Day to be national category, got '%s'", holiday.Category)
	}
}

func TestRUWomensDay(t *testing.T) {
	provider := NewRUProvider()
	holidays := provider.LoadHolidays(2024)

	// Test International Women's Day (March 8)
	womensDay := time.Date(2024, 3, 8, 0, 0, 0, 0, time.UTC)
	holiday, exists := holidays[womensDay]
	if !exists {
		t.Fatal("International Women's Day not found")
	}

	if holiday.Name != "Международный женский день" {
		t.Errorf("Expected Women's Day name 'Международный женский день', got '%s'", holiday.Name)
	}

	if holiday.Category != "national" {
		t.Errorf("Expected Women's Day to be national category, got '%s'", holiday.Category)
	}
}

func TestRURussiaDay(t *testing.T) {
	provider := NewRUProvider()
	holidays := provider.LoadHolidays(2024)

	// Test Russia Day (June 12)
	russiaDay := time.Date(2024, 6, 12, 0, 0, 0, 0, time.UTC)
	holiday, exists := holidays[russiaDay]
	if !exists {
		t.Fatal("Russia Day not found")
	}

	if holiday.Name != "День России" {
		t.Errorf("Expected Russia Day name 'День России', got '%s'", holiday.Name)
	}

	if holiday.Category != "national" {
		t.Errorf("Expected Russia Day to be national category, got '%s'", holiday.Category)
	}
}

func TestRUUnityDay(t *testing.T) {
	provider := NewRUProvider()
	holidays := provider.LoadHolidays(2024)

	// Test Unity Day (November 4)
	unityDay := time.Date(2024, 11, 4, 0, 0, 0, 0, time.UTC)
	holiday, exists := holidays[unityDay]
	if !exists {
		t.Fatal("Unity Day not found")
	}

	if holiday.Name != "День народного единства" {
		t.Errorf("Expected Unity Day name 'День народного единства', got '%s'", holiday.Name)
	}

	if holiday.Category != "national" {
		t.Errorf("Expected Unity Day to be national category, got '%s'", holiday.Category)
	}
}

func TestRUConstitutionDay(t *testing.T) {
	provider := NewRUProvider()

	// Test that Constitution Day exists for years <= 2004
	holidays2004 := provider.LoadHolidays(2004)
	constitutionDay2004 := time.Date(2004, 12, 12, 0, 0, 0, 0, time.UTC)
	if _, exists := holidays2004[constitutionDay2004]; !exists {
		t.Error("Constitution Day should exist in 2004")
	}

	// Test that Constitution Day doesn't exist for years > 2004
	holidays2005 := provider.LoadHolidays(2005)
	constitutionDay2005 := time.Date(2005, 12, 12, 0, 0, 0, 0, time.UTC)
	if _, exists := holidays2005[constitutionDay2005]; exists {
		t.Error("Constitution Day should not exist in 2005")
	}
}

func TestRUOrthodoxEasterCalculation(t *testing.T) {
	provider := NewRUProvider()

	// Test Orthodox Easter dates for known years
	testCases := []struct {
		year  int
		month time.Month
		day   int
	}{
		{2024, time.May, 5},    // Orthodox Easter 2024
		{2025, time.April, 20}, // Orthodox Easter 2025
		{2026, time.April, 12}, // Orthodox Easter 2026
		{2027, time.May, 2},    // Orthodox Easter 2027
	}

	for _, tc := range testCases {
		holidays := provider.LoadHolidays(tc.year)
		expected := time.Date(tc.year, tc.month, tc.day, 0, 0, 0, 0, time.UTC)

		// Find Orthodox Easter in holidays
		var easterFound bool
		for date, holiday := range holidays {
			if holiday.Name == "Пасха" && date.Equal(expected) {
				easterFound = true
				break
			}
		}

		if !easterFound {
			t.Errorf("Expected Orthodox Easter %d to be %s, but not found in holidays", tc.year, expected.Format("2006-01-02"))
		}
	}
}

func BenchmarkRUProvider(b *testing.B) {
	provider := NewRUProvider()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = provider.LoadHolidays(2024)
	}
}

func BenchmarkRUOrthodoxEasterCalculation(b *testing.B) {
	provider := NewRUProvider()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = provider.calculateOrthodoxEaster(2024)
	}
}
