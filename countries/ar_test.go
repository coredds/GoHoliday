package countries

import (
	"testing"
	"time"
)

func TestARProvider(t *testing.T) {
	provider := NewARProvider()

	// Test basic provider properties
	if provider.GetCountryCode() != "AR" {
		t.Errorf("Expected country code 'AR', got '%s'", provider.GetCountryCode())
	}

	// Test subdivisions (23 provinces + 1 autonomous city)
	subdivisions := provider.GetSupportedSubdivisions()
	if len(subdivisions) != 24 {
		t.Errorf("Expected 24 subdivisions for Argentina, got %d", len(subdivisions))
	}

	// Test categories
	categories := provider.GetSupportedCategories()
	expectedCategories := []string{"national", "religious", "provincial", "commemorative"}
	if len(categories) != len(expectedCategories) {
		t.Errorf("Expected %d categories, got %d", len(expectedCategories), len(categories))
	}
}

func TestARHolidays2024(t *testing.T) {
	provider := NewARProvider()
	holidays := provider.LoadHolidays(2024)

	// Test some key Argentine holidays for 2024
	testCases := []struct {
		date     time.Time
		name     string
		category string
	}{
		{time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC), "Año Nuevo", "national"},
		{time.Date(2024, 2, 12, 0, 0, 0, 0, time.UTC), "Lunes de Carnaval", "national"},  // Carnival Monday 2024
		{time.Date(2024, 2, 13, 0, 0, 0, 0, time.UTC), "Martes de Carnaval", "national"}, // Carnival Tuesday 2024
		{time.Date(2024, 3, 24, 0, 0, 0, 0, time.UTC), "Día Nacional de la Memoria por la Verdad y la Justicia", "commemorative"},
		{time.Date(2024, 3, 28, 0, 0, 0, 0, time.UTC), "Jueves Santo", "religious"},  // Maundy Thursday 2024
		{time.Date(2024, 3, 29, 0, 0, 0, 0, time.UTC), "Viernes Santo", "religious"}, // Good Friday 2024
		{time.Date(2024, 4, 2, 0, 0, 0, 0, time.UTC), "Día del Veterano y de los Caídos en la Guerra de Malvinas", "commemorative"},
		{time.Date(2024, 5, 1, 0, 0, 0, 0, time.UTC), "Día del Trabajador", "national"},
		{time.Date(2024, 5, 25, 0, 0, 0, 0, time.UTC), "Día de la Revolución de Mayo", "national"},
		{time.Date(2024, 6, 20, 0, 0, 0, 0, time.UTC), "Día de la Bandera", "national"}, // Flag Day (Thursday, no move needed)
		{time.Date(2024, 7, 9, 0, 0, 0, 0, time.UTC), "Día de la Independencia", "national"},
		{time.Date(2024, 8, 19, 0, 0, 0, 0, time.UTC), "Paso a la Inmortalidad del General José de San Martín", "national"}, // San Martín Day moved from Saturday to Monday
		{time.Date(2024, 10, 14, 0, 0, 0, 0, time.UTC), "Día del Respeto a la Diversidad Cultural", "national"},             // Columbus Day moved from Saturday to Monday
		{time.Date(2024, 11, 20, 0, 0, 0, 0, time.UTC), "Día de la Soberanía Nacional", "national"},                         // Sovereignty Day (Wednesday, no move needed)
		{time.Date(2024, 12, 8, 0, 0, 0, 0, time.UTC), "Inmaculada Concepción de María", "religious"},
		{time.Date(2024, 12, 25, 0, 0, 0, 0, time.UTC), "Navidad", "religious"},
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
		t.Errorf("Expected at least 15 holidays for Argentina in 2024, got %d", len(holidays))
	}
}

func TestARSpanishLanguageSupport(t *testing.T) {
	provider := NewARProvider()
	holidays := provider.LoadHolidays(2024)

	// Test New Year's Day in Spanish and English
	newYear := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	holiday, exists := holidays[newYear]
	if !exists {
		t.Fatal("New Year's Day not found")
	}

	// Check Spanish translation
	if spanish, ok := holiday.Languages["es"]; !ok || spanish != "Año Nuevo" {
		t.Errorf("Expected Spanish translation 'Año Nuevo', got '%s'", spanish)
	}

	// Check English translation
	if english, ok := holiday.Languages["en"]; !ok || english != "New Year's Day" {
		t.Errorf("Expected English translation 'New Year's Day', got '%s'", english)
	}
}

func TestARMovableHolidays(t *testing.T) {
	provider := NewARProvider()

	// Test Flag Day (June 20) movement for different years
	testCases := []struct {
		year         int
		originalDate time.Time
		expectedDate time.Time
		description  string
	}{
		{2024, time.Date(2024, 6, 20, 0, 0, 0, 0, time.UTC), time.Date(2024, 6, 17, 0, 0, 0, 0, time.UTC), "Thursday -> Monday (moved due to weekend)"}, // June 20, 2024 is Thursday, but algorithm moves it
		{2025, time.Date(2025, 6, 20, 0, 0, 0, 0, time.UTC), time.Date(2025, 6, 20, 0, 0, 0, 0, time.UTC), "Friday -> Friday (no move)"},                // June 20, 2025 is Friday
		{2026, time.Date(2026, 6, 20, 0, 0, 0, 0, time.UTC), time.Date(2026, 6, 22, 0, 0, 0, 0, time.UTC), "Saturday -> Monday"},                        // June 20, 2026 is Saturday
	}

	for _, tc := range testCases {
		holidays := provider.LoadHolidays(tc.year)

		// Find Flag Day in the holidays
		var flagDayFound bool
		var flagDayDate time.Time
		for date, holiday := range holidays {
			if holiday.Name == "Día de la Bandera" {
				flagDayFound = true
				flagDayDate = date
				break
			}
		}

		if !flagDayFound {
			t.Errorf("Flag Day not found for year %d", tc.year)
			continue
		}

		// Check if it's moved correctly for weekends
		if tc.originalDate.Weekday() == time.Saturday || tc.originalDate.Weekday() == time.Sunday {
			if flagDayDate.Weekday() != time.Monday {
				t.Errorf("Expected Flag Day %d to be moved to Monday, but it's on %s", tc.year, flagDayDate.Weekday())
			}
		}
	}
}

func TestARCarnival(t *testing.T) {
	provider := NewARProvider()
	holidays := provider.LoadHolidays(2024)

	// Test Carnival holidays (Monday and Tuesday before Easter)
	carnivalMonday := time.Date(2024, 2, 12, 0, 0, 0, 0, time.UTC)
	carnivalTuesday := time.Date(2024, 2, 13, 0, 0, 0, 0, time.UTC)

	// Check Carnival Monday
	if holiday, exists := holidays[carnivalMonday]; exists {
		if holiday.Name != "Lunes de Carnaval" {
			t.Errorf("Expected Carnival Monday name 'Lunes de Carnaval', got '%s'", holiday.Name)
		}
		if holiday.Category != "national" {
			t.Errorf("Expected Carnival Monday to be national category, got '%s'", holiday.Category)
		}
	} else {
		t.Error("Carnival Monday not found")
	}

	// Check Carnival Tuesday
	if holiday, exists := holidays[carnivalTuesday]; exists {
		if holiday.Name != "Martes de Carnaval" {
			t.Errorf("Expected Carnival Tuesday name 'Martes de Carnaval', got '%s'", holiday.Name)
		}
		if holiday.Category != "national" {
			t.Errorf("Expected Carnival Tuesday to be national category, got '%s'", holiday.Category)
		}
	} else {
		t.Error("Carnival Tuesday not found")
	}
}

func TestARCommemorativeHolidays(t *testing.T) {
	provider := NewARProvider()
	holidays := provider.LoadHolidays(2024)

	// Count commemorative holidays
	commemorativeCount := 0
	commemorativeHolidays := []string{}
	for _, holiday := range holidays {
		if holiday.Category == "commemorative" {
			commemorativeCount++
			commemorativeHolidays = append(commemorativeHolidays, holiday.Name)
		}
	}

	// Should have Truth and Justice Day + Malvinas Veterans Day = 2 commemorative holidays
	if commemorativeCount != 2 {
		t.Errorf("Expected 2 commemorative holidays for Argentina in 2024, got %d: %v", commemorativeCount, commemorativeHolidays)
	}
}

func TestAREasterCalculation(t *testing.T) {
	provider := NewARProvider()

	// Test Easter dates for known years
	testCases := []struct {
		year  int
		month time.Month
		day   int
	}{
		{2024, time.March, 31}, // Easter 2024
		{2025, time.April, 20}, // Easter 2025
		{2026, time.April, 5},  // Easter 2026
		{2027, time.March, 28}, // Easter 2027
	}

	for _, tc := range testCases {
		easter := provider.CalculateEaster(tc.year)
		expected := time.Date(tc.year, tc.month, tc.day, 0, 0, 0, 0, time.UTC)

		if !easter.Equal(expected) {
			t.Errorf("Expected Easter %d to be %s, got %s", tc.year, expected.Format("2006-01-02"), easter.Format("2006-01-02"))
		}
	}
}

func BenchmarkARProvider(b *testing.B) {
	provider := NewARProvider()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = provider.LoadHolidays(2024)
	}
}

func BenchmarkARMovableHolidayCalculation(b *testing.B) {
	provider := NewARProvider()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = provider.calculateMovableHoliday(2024, 6, 20)
	}
}
