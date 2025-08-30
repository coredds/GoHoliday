package countries

import (
	"time"
)

// PTProvider implements holiday calculations for Portugal
type PTProvider struct {
	*BaseProvider
}

// NewPTProvider creates a new Portuguese holiday provider
func NewPTProvider() *PTProvider {
	base := NewBaseProvider("PT")
	base.subdivisions = []string{
		// Districts (18 mainland districts)
		"01", "02", "03", "04", "05", "06", "07", "08", "09", "10",
		"11", "12", "13", "14", "15", "16", "17", "18",
		// Autonomous regions
		"20", "30", // Azores, Madeira
	}
	base.categories = []string{"public", "religious", "regional", "municipal"}

	return &PTProvider{BaseProvider: base}
}

// LoadHolidays loads all Portuguese holidays for a given year
func (pt *PTProvider) LoadHolidays(year int) map[time.Time]*Holiday {
	holidays := make(map[time.Time]*Holiday)

	// Fixed date holidays
	pt.addFixedHolidays(holidays, year)

	// Easter-based holidays
	pt.addEasterHolidays(holidays, year)

	return holidays
}

// addFixedHolidays adds fixed-date Portuguese holidays
func (pt *PTProvider) addFixedHolidays(holidays map[time.Time]*Holiday, year int) {
	fixedHolidays := []struct {
		month    int
		day      int
		name     string
		nameEn   string
		category string
	}{
		{1, 1, "Ano Novo", "New Year's Day", "public"},
		{4, 25, "Dia da Liberdade", "Freedom Day", "public"},
		{5, 1, "Dia do Trabalhador", "Labour Day", "public"},
		{6, 10, "Dia de Portugal", "Portugal Day", "public"},
		{8, 15, "Assunção de Nossa Senhora", "Assumption of Mary", "religious"},
		{10, 5, "Implantação da República", "Republic Day", "public"},
		{11, 1, "Dia de Todos os Santos", "All Saints' Day", "religious"},
		{12, 1, "Restauração da Independência", "Restoration of Independence", "public"},
		{12, 8, "Imaculada Conceição", "Immaculate Conception", "religious"},
		{12, 25, "Natal", "Christmas Day", "religious"},
	}

	for _, h := range fixedHolidays {
		date := time.Date(year, time.Month(h.month), h.day, 0, 0, 0, 0, time.UTC)
		holidays[date] = &Holiday{
			Name:     h.nameEn,
			Date:     date,
			Category: h.category,
			Languages: map[string]string{
				"en": h.nameEn,
				"pt": h.name,
			},
			IsObserved: true,
		}
	}
}

// addEasterHolidays adds Easter-based Portuguese holidays
func (pt *PTProvider) addEasterHolidays(holidays map[time.Time]*Holiday, year int) {
	easter := pt.calculateEaster(year)

	easterHolidays := []struct {
		offset   int
		name     string
		nameEn   string
		category string
	}{
		{-47, "Carnaval", "Carnival Tuesday", "public"}, // Shrove Tuesday
		{-2, "Sexta-feira Santa", "Good Friday", "religious"},
		{0, "Páscoa", "Easter Sunday", "religious"},
		{60, "Corpo de Deus", "Corpus Christi", "religious"}, // 60 days after Easter
	}

	for _, h := range easterHolidays {
		date := easter.AddDate(0, 0, h.offset)
		holidays[date] = &Holiday{
			Name:     h.nameEn,
			Date:     date,
			Category: h.category,
			Languages: map[string]string{
				"en": h.nameEn,
				"pt": h.name,
			},
			IsObserved: true,
		}
	}
}

// calculateEaster calculates Easter Sunday for a given year using the Western (Gregorian) calendar
func (pt *PTProvider) calculateEaster(year int) time.Time {
	// Using the algorithm for Western Easter (Gregorian calendar)
	// This is the standard algorithm used in Western Christianity

	a := year % 19
	b := year / 100
	c := year % 100
	d := b / 4
	e := b % 4
	f := (b + 8) / 25
	g := (b - f + 1) / 3
	h := (19*a + b - d - g + 15) % 30
	i := c / 4
	k := c % 4
	l := (32 + 2*e + 2*i - h - k) % 7
	m := (a + 11*h + 22*l) / 451
	month := (h + l - 7*m + 114) / 31
	day := ((h + l - 7*m + 114) % 31) + 1

	return time.Date(year, time.Month(month), day, 0, 0, 0, 0, time.UTC)
}
