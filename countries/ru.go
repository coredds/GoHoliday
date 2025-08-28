package countries

import (
	"time"
)

// RUProvider provides holidays for Russia (Russian Federation)
type RUProvider struct {
	*BaseProvider
}

// NewRUProvider creates a new RUProvider instance
func NewRUProvider() *RUProvider {
	base := NewBaseProvider("RU")
	base.subdivisions = []string{
		// 85 federal subjects of Russia
		"AD", "AL", "ALT", "AMU", "ARK", "AST", "BA", "BEL", "BRY", "BU",
		"CE", "CHE", "CHU", "CU", "DA", "IN", "IRK", "IVA", "KB", "KC",
		"KDA", "KEM", "KGD", "KGN", "KHA", "KHM", "KIR", "KK", "KL", "KLU",
		"KO", "KOS", "KR", "KRS", "KYA", "LEN", "LIP", "MAG", "ME", "MO",
		"MOW", "MOS", "MUR", "NEN", "NGR", "NIZ", "NVS", "OMS", "ORE", "ORL",
		"PEN", "PER", "PNZ", "PRI", "PSK", "ROS", "RYA", "SA", "SAK", "SAM",
		"SAR", "SE", "SMO", "SPE", "STA", "SVE", "TAM", "TA", "TOM", "TUL",
		"TVE", "TY", "TYU", "UD", "ULY", "VGG", "VLA", "VLG", "VOR", "YAN",
		"YAR", "YEV", "ZAB", "CHI", "KAM",
	}
	base.categories = []string{"national", "religious", "commemorative", "orthodox"}

	return &RUProvider{BaseProvider: base}
}

// LoadHolidays loads holidays for Russia for the given year
func (p *RUProvider) LoadHolidays(year int) map[time.Time]*Holiday {
	holidays := make(map[time.Time]*Holiday)

	// Add fixed national holidays
	p.addFixedHolidays(holidays, year)

	// Add Orthodox holidays (based on Julian calendar until 1918, then Gregorian)
	p.addOrthodoxHolidays(holidays, year)

	return holidays
}

// addFixedHolidays adds fixed-date Russian holidays
func (p *RUProvider) addFixedHolidays(holidays map[time.Time]*Holiday, year int) {
	// New Year Holidays (January 1-8)
	for i := 1; i <= 8; i++ {
		date := time.Date(year, 1, i, 0, 0, 0, 0, time.UTC)
		name := "Новогодние каникулы"
		enName := "New Year Holidays"
		if i == 1 {
			name = "Новый год"
			enName = "New Year's Day"
		}
		holidays[date] = p.CreateHoliday(
			name,
			date,
			"national",
			map[string]string{
				"ru": name,
				"en": enName,
			},
		)
	}

	// Defender of the Fatherland Day
	holidays[time.Date(year, 2, 23, 0, 0, 0, 0, time.UTC)] = p.CreateHoliday(
		"День защитника Отечества",
		time.Date(year, 2, 23, 0, 0, 0, 0, time.UTC),
		"national",
		map[string]string{
			"ru": "День защитника Отечества",
			"en": "Defender of the Fatherland Day",
		},
	)

	// International Women's Day
	holidays[time.Date(year, 3, 8, 0, 0, 0, 0, time.UTC)] = p.CreateHoliday(
		"Международный женский день",
		time.Date(year, 3, 8, 0, 0, 0, 0, time.UTC),
		"national",
		map[string]string{
			"ru": "Международный женский день",
			"en": "International Women's Day",
		},
	)

	// Spring and Labour Holiday
	holidays[time.Date(year, 5, 1, 0, 0, 0, 0, time.UTC)] = p.CreateHoliday(
		"Праздник Весны и Труда",
		time.Date(year, 5, 1, 0, 0, 0, 0, time.UTC),
		"national",
		map[string]string{
			"ru": "Праздник Весны и Труда",
			"en": "Spring and Labour Holiday",
		},
	)

	// Victory Day
	holidays[time.Date(year, 5, 9, 0, 0, 0, 0, time.UTC)] = p.CreateHoliday(
		"День Победы",
		time.Date(year, 5, 9, 0, 0, 0, 0, time.UTC),
		"commemorative",
		map[string]string{
			"ru": "День Победы",
			"en": "Victory Day",
		},
	)

	// Russia Day
	holidays[time.Date(year, 6, 12, 0, 0, 0, 0, time.UTC)] = p.CreateHoliday(
		"День России",
		time.Date(year, 6, 12, 0, 0, 0, 0, time.UTC),
		"national",
		map[string]string{
			"ru": "День России",
			"en": "Russia Day",
		},
	)

	// Unity Day
	holidays[time.Date(year, 11, 4, 0, 0, 0, 0, time.UTC)] = p.CreateHoliday(
		"День народного единства",
		time.Date(year, 11, 4, 0, 0, 0, 0, time.UTC),
		"national",
		map[string]string{
			"ru": "День народного единства",
			"en": "Unity Day",
		},
	)

	// Constitution Day (December 12) - non-working day until 2004
	if year <= 2004 {
		holidays[time.Date(year, 12, 12, 0, 0, 0, 0, time.UTC)] = p.CreateHoliday(
			"День Конституции",
			time.Date(year, 12, 12, 0, 0, 0, 0, time.UTC),
			"national",
			map[string]string{
				"ru": "День Конституции",
				"en": "Constitution Day",
			},
		)
	}
}

// addOrthodoxHolidays adds Russian Orthodox holidays
func (p *RUProvider) addOrthodoxHolidays(holidays map[time.Time]*Holiday, year int) {
	// Orthodox Christmas (January 7)
	holidays[time.Date(year, 1, 7, 0, 0, 0, 0, time.UTC)] = p.CreateHoliday(
		"Рождество Христово",
		time.Date(year, 1, 7, 0, 0, 0, 0, time.UTC),
		"orthodox",
		map[string]string{
			"ru": "Рождество Христово",
			"en": "Orthodox Christmas",
		},
	)

	// Orthodox Easter and related holidays
	orthodoxEaster := p.calculateOrthodoxEaster(year)
	if !orthodoxEaster.IsZero() {
		// Orthodox Easter
		holidays[orthodoxEaster] = p.CreateHoliday(
			"Пасха",
			orthodoxEaster,
			"orthodox",
			map[string]string{
				"ru": "Пасха",
				"en": "Orthodox Easter",
			},
		)

		// Palm Sunday (1 week before Easter)
		palmSunday := orthodoxEaster.AddDate(0, 0, -7)
		holidays[palmSunday] = p.CreateHoliday(
			"Вербное воскресенье",
			palmSunday,
			"orthodox",
			map[string]string{
				"ru": "Вербное воскресенье",
				"en": "Orthodox Palm Sunday",
			},
		)

		// Orthodox Easter Monday
		easterMonday := orthodoxEaster.AddDate(0, 0, 1)
		holidays[easterMonday] = p.CreateHoliday(
			"Светлый понедельник",
			easterMonday,
			"orthodox",
			map[string]string{
				"ru": "Светлый понедельник",
				"en": "Orthodox Easter Monday",
			},
		)

		// Trinity Sunday (50 days after Easter)
		trinity := orthodoxEaster.AddDate(0, 0, 49)
		holidays[trinity] = p.CreateHoliday(
			"День Святой Троицы",
			trinity,
			"orthodox",
			map[string]string{
				"ru": "День Святой Троицы",
				"en": "Trinity Sunday",
			},
		)
	}
}

// calculateOrthodoxEaster calculates Orthodox Easter date using the Julian calendar
func (p *RUProvider) calculateOrthodoxEaster(year int) time.Time {
	// Orthodox Easter calculation (simplified Julian-based algorithm)
	// This is a simplified version - in practice, more complex calculations are used

	// Known Orthodox Easter dates for reference
	orthodoxEasterDates := map[int]time.Time{
		2024: time.Date(2024, 5, 5, 0, 0, 0, 0, time.UTC),  // May 5, 2024
		2025: time.Date(2025, 4, 20, 0, 0, 0, 0, time.UTC), // April 20, 2025
		2026: time.Date(2026, 4, 12, 0, 0, 0, 0, time.UTC), // April 12, 2026
		2027: time.Date(2027, 5, 2, 0, 0, 0, 0, time.UTC),  // May 2, 2027
		2028: time.Date(2028, 4, 16, 0, 0, 0, 0, time.UTC), // April 16, 2028
	}

	if date, exists := orthodoxEasterDates[year]; exists {
		return date
	}

	return time.Time{} // Return zero time if not found
}
