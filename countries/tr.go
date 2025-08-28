package countries

import (
	"time"
)

// TRProvider provides holidays for Turkey (Turkish Republic)
type TRProvider struct {
	*BaseProvider
}

// NewTRProvider creates a new TRProvider instance
func NewTRProvider() *TRProvider {
	return &TRProvider{
		BaseProvider: &BaseProvider{
			countryCode: "TR",
			subdivisions: []string{
				// 81 provinces of Turkey
				"01", "02", "03", "04", "05", "06", "07", "08", "09", "10",
				"11", "12", "13", "14", "15", "16", "17", "18", "19", "20",
				"21", "22", "23", "24", "25", "26", "27", "28", "29", "30",
				"31", "32", "33", "34", "35", "36", "37", "38", "39", "40",
				"41", "42", "43", "44", "45", "46", "47", "48", "49", "50",
				"51", "52", "53", "54", "55", "56", "57", "58", "59", "60",
				"61", "62", "63", "64", "65", "66", "67", "68", "69", "70",
				"71", "72", "73", "74", "75", "76", "77", "78", "79", "80", "81",
			},
			categories: []string{"national", "religious", "commemorative", "seasonal"},
		},
	}
}

// LoadHolidays loads holidays for Turkey for the given year
func (p *TRProvider) LoadHolidays(year int) map[time.Time]*Holiday {
	holidays := make(map[time.Time]*Holiday)

	// Add fixed national holidays
	p.addFixedHolidays(holidays, year)

	// Add Islamic holidays (based on lunar calendar)
	p.addIslamicHolidays(holidays, year)

	return holidays
}

// addFixedHolidays adds fixed-date Turkish holidays
func (p *TRProvider) addFixedHolidays(holidays map[time.Time]*Holiday, year int) {
	// New Year's Day
	holidays[time.Date(year, 1, 1, 0, 0, 0, 0, time.UTC)] = p.CreateHoliday(
		"Yılbaşı",
		time.Date(year, 1, 1, 0, 0, 0, 0, time.UTC),
		"national",
		map[string]string{
			"tr": "Yılbaşı",
			"en": "New Year's Day",
		},
	)

	// National Sovereignty and Children's Day
	holidays[time.Date(year, 4, 23, 0, 0, 0, 0, time.UTC)] = p.CreateHoliday(
		"Ulusal Egemenlik ve Çocuk Bayramı",
		time.Date(year, 4, 23, 0, 0, 0, 0, time.UTC),
		"national",
		map[string]string{
			"tr": "Ulusal Egemenlik ve Çocuk Bayramı",
			"en": "National Sovereignty and Children's Day",
		},
	)

	// Labour and Solidarity Day
	holidays[time.Date(year, 5, 1, 0, 0, 0, 0, time.UTC)] = p.CreateHoliday(
		"Emek ve Dayanışma Günü",
		time.Date(year, 5, 1, 0, 0, 0, 0, time.UTC),
		"national",
		map[string]string{
			"tr": "Emek ve Dayanışma Günü",
			"en": "Labour and Solidarity Day",
		},
	)

	// Commemoration of Atatürk, Youth and Sports Day
	holidays[time.Date(year, 5, 19, 0, 0, 0, 0, time.UTC)] = p.CreateHoliday(
		"Atatürk'ü Anma, Gençlik ve Spor Bayramı",
		time.Date(year, 5, 19, 0, 0, 0, 0, time.UTC),
		"commemorative",
		map[string]string{
			"tr": "Atatürk'ü Anma, Gençlik ve Spor Bayramı",
			"en": "Commemoration of Atatürk, Youth and Sports Day",
		},
	)

	// Democracy and National Unity Day (since 2017)
	if year >= 2017 {
		holidays[time.Date(year, 7, 15, 0, 0, 0, 0, time.UTC)] = p.CreateHoliday(
			"Demokrasi ve Milli Birlik Günü",
			time.Date(year, 7, 15, 0, 0, 0, 0, time.UTC),
			"commemorative",
			map[string]string{
				"tr": "Demokrasi ve Milli Birlik Günü",
				"en": "Democracy and National Unity Day",
			},
		)
	}

	// Victory Day
	holidays[time.Date(year, 8, 30, 0, 0, 0, 0, time.UTC)] = p.CreateHoliday(
		"Zafer Bayramı",
		time.Date(year, 8, 30, 0, 0, 0, 0, time.UTC),
		"national",
		map[string]string{
			"tr": "Zafer Bayramı",
			"en": "Victory Day",
		},
	)

	// Republic Day
	holidays[time.Date(year, 10, 29, 0, 0, 0, 0, time.UTC)] = p.CreateHoliday(
		"Cumhuriyet Bayramı",
		time.Date(year, 10, 29, 0, 0, 0, 0, time.UTC),
		"national",
		map[string]string{
			"tr": "Cumhuriyet Bayramı",
			"en": "Republic Day",
		},
	)
}

// addIslamicHolidays adds Islamic holidays based on lunar calendar
func (p *TRProvider) addIslamicHolidays(holidays map[time.Time]*Holiday, year int) {
	// Note: Islamic holidays are based on lunar calendar and shift each year
	// These are approximations for common years. In practice, these would be
	// calculated using precise astronomical calculations or official announcements.

	// Ramadan Festival (Eid al-Fitr) - 3 days
	ramadanFestival := p.getIslamicHolidayDate(year, "ramadan")
	if !ramadanFestival.IsZero() {
		for i := 0; i < 3; i++ {
			date := ramadanFestival.AddDate(0, 0, i)
			dayName := "Ramazan Bayramı"
			enName := "Ramadan Festival"
			if i == 0 {
				dayName += " 1. Gün"
				enName += " Day 1"
			} else if i == 1 {
				dayName += " 2. Gün"
				enName += " Day 2"
			} else {
				dayName += " 3. Gün"
				enName += " Day 3"
			}

			holidays[date] = p.CreateHoliday(
				dayName,
				date,
				"religious",
				map[string]string{
					"tr": dayName,
					"en": enName,
				},
			)
		}
	}

	// Sacrifice Festival (Eid al-Adha) - 4 days
	sacrificeFestival := p.getIslamicHolidayDate(year, "sacrifice")
	if !sacrificeFestival.IsZero() {
		for i := 0; i < 4; i++ {
			date := sacrificeFestival.AddDate(0, 0, i)
			dayName := "Kurban Bayramı"
			enName := "Sacrifice Festival"
			if i == 0 {
				dayName += " 1. Gün"
				enName += " Day 1"
			} else if i == 1 {
				dayName += " 2. Gün"
				enName += " Day 2"
			} else if i == 2 {
				dayName += " 3. Gün"
				enName += " Day 3"
			} else {
				dayName += " 4. Gün"
				enName += " Day 4"
			}

			holidays[date] = p.CreateHoliday(
				dayName,
				date,
				"religious",
				map[string]string{
					"tr": dayName,
					"en": enName,
				},
			)
		}
	}
}

// getIslamicHolidayDate returns approximate dates for Islamic holidays
// In practice, these would be calculated using precise astronomical calculations
func (p *TRProvider) getIslamicHolidayDate(year int, holiday string) time.Time {
	// Approximate dates for major Islamic holidays
	// These shift about 11 days earlier each year due to lunar calendar

	islamicHolidays := map[int]map[string]time.Time{
		2024: {
			"ramadan":   time.Date(2024, 4, 10, 0, 0, 0, 0, time.UTC), // Eid al-Fitr
			"sacrifice": time.Date(2024, 6, 16, 0, 0, 0, 0, time.UTC), // Eid al-Adha
		},
		2025: {
			"ramadan":   time.Date(2025, 3, 30, 0, 0, 0, 0, time.UTC), // Eid al-Fitr
			"sacrifice": time.Date(2025, 6, 6, 0, 0, 0, 0, time.UTC),  // Eid al-Adha
		},
		2026: {
			"ramadan":   time.Date(2026, 3, 20, 0, 0, 0, 0, time.UTC), // Eid al-Fitr
			"sacrifice": time.Date(2026, 5, 26, 0, 0, 0, 0, time.UTC), // Eid al-Adha
		},
		2027: {
			"ramadan":   time.Date(2027, 3, 9, 0, 0, 0, 0, time.UTC),  // Eid al-Fitr
			"sacrifice": time.Date(2027, 5, 16, 0, 0, 0, 0, time.UTC), // Eid al-Adha
		},
	}

	if yearHolidays, exists := islamicHolidays[year]; exists {
		if date, exists := yearHolidays[holiday]; exists {
			return date
		}
	}

	return time.Time{} // Return zero time if not found
}
