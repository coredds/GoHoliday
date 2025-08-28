package countries

import (
	"time"
)

// IDProvider provides holidays for Indonesia (Republic of Indonesia)
type IDProvider struct {
	*BaseProvider
}

// NewIDProvider creates a new IDProvider instance
func NewIDProvider() *IDProvider {
	base := NewBaseProvider("ID")
	base.subdivisions = []string{
		// 38 provinces of Indonesia
		"AC", "BA", "BB", "BT", "BE", "GO", "JA", "JB", "JT", "JI", "YO", "JK",
		"KS", "KB", "KT", "KI", "LA", "MA", "NB", "NT", "PA", "PB", "RI", "SR",
		"SN", "SS", "SB", "SG", "ST", "SU", "SL", "YO", "1024", "KU", "KR", "SG", "PE", "PP",
	}
	base.categories = []string{"national", "religious", "islamic", "christian", "buddhist", "hindu", "chinese"}

	return &IDProvider{BaseProvider: base}
}

// LoadHolidays loads holidays for Indonesia for the given year
func (p *IDProvider) LoadHolidays(year int) map[time.Time]*Holiday {
	holidays := make(map[time.Time]*Holiday)

	// Add fixed national holidays
	p.addFixedHolidays(holidays, year)

	// Add Islamic holidays (based on lunar calendar)
	p.addIslamicHolidays(holidays, year)

	// Add Christian holidays
	p.addChristianHolidays(holidays, year)

	// Add other religious holidays
	p.addOtherReligiousHolidays(holidays, year)

	return holidays
}

// addFixedHolidays adds fixed-date Indonesian holidays
func (p *IDProvider) addFixedHolidays(holidays map[time.Time]*Holiday, year int) {
	// New Year's Day
	holidays[time.Date(year, 1, 1, 0, 0, 0, 0, time.UTC)] = p.CreateHoliday(
		"Tahun Baru Masehi",
		time.Date(year, 1, 1, 0, 0, 0, 0, time.UTC),
		"national",
		map[string]string{
			"id": "Tahun Baru Masehi",
			"en": "New Year's Day",
		},
	)

	// Labor Day
	holidays[time.Date(year, 5, 1, 0, 0, 0, 0, time.UTC)] = p.CreateHoliday(
		"Hari Buruh Internasional",
		time.Date(year, 5, 1, 0, 0, 0, 0, time.UTC),
		"national",
		map[string]string{
			"id": "Hari Buruh Internasional",
			"en": "International Labor Day",
		},
	)

	// Pancasila Day
	holidays[time.Date(year, 6, 1, 0, 0, 0, 0, time.UTC)] = p.CreateHoliday(
		"Hari Lahir Pancasila",
		time.Date(year, 6, 1, 0, 0, 0, 0, time.UTC),
		"national",
		map[string]string{
			"id": "Hari Lahir Pancasila",
			"en": "Pancasila Day",
		},
	)

	// Independence Day
	holidays[time.Date(year, 8, 17, 0, 0, 0, 0, time.UTC)] = p.CreateHoliday(
		"Hari Kemerdekaan Republik Indonesia",
		time.Date(year, 8, 17, 0, 0, 0, 0, time.UTC),
		"national",
		map[string]string{
			"id": "Hari Kemerdekaan Republik Indonesia",
			"en": "Independence Day",
		},
	)

	// Heroes' Day
	holidays[time.Date(year, 11, 10, 0, 0, 0, 0, time.UTC)] = p.CreateHoliday(
		"Hari Pahlawan",
		time.Date(year, 11, 10, 0, 0, 0, 0, time.UTC),
		"national",
		map[string]string{
			"id": "Hari Pahlawan",
			"en": "Heroes' Day",
		},
	)
}

// addIslamicHolidays adds Islamic holidays (Indonesia has the largest Muslim population)
func (p *IDProvider) addIslamicHolidays(holidays map[time.Time]*Holiday, year int) {
	// Islamic New Year
	islamicNewYear := p.getIslamicHolidayDate(year, "islamic_new_year")
	if !islamicNewYear.IsZero() {
		holidays[islamicNewYear] = p.CreateHoliday(
			"Tahun Baru Islam",
			islamicNewYear,
			"islamic",
			map[string]string{
				"id": "Tahun Baru Islam",
				"en": "Islamic New Year",
			},
		)
	}

	// Maulid (Prophet Muhammad's Birthday)
	maulid := p.getIslamicHolidayDate(year, "maulid")
	if !maulid.IsZero() {
		holidays[maulid] = p.CreateHoliday(
			"Maulid Nabi Muhammad SAW",
			maulid,
			"islamic",
			map[string]string{
				"id": "Maulid Nabi Muhammad SAW",
				"en": "Prophet Muhammad's Birthday",
			},
		)
	}

	// Isra Mi'raj
	israMiraj := p.getIslamicHolidayDate(year, "isra_miraj")
	if !israMiraj.IsZero() {
		holidays[israMiraj] = p.CreateHoliday(
			"Isra Mi'raj",
			israMiraj,
			"islamic",
			map[string]string{
				"id": "Isra Mi'raj",
				"en": "Isra and Mi'raj",
			},
		)
	}

	// Eid al-Fitr (Idul Fitri) - 2 days
	iduFitri := p.getIslamicHolidayDate(year, "idul_fitri")
	if !iduFitri.IsZero() {
		for i := 0; i < 2; i++ {
			date := iduFitri.AddDate(0, 0, i)
			dayName := "Hari Raya Idul Fitri"
			if i == 1 {
				dayName += " Kedua"
			}
			holidays[date] = p.CreateHoliday(
				dayName,
				date,
				"islamic",
				map[string]string{
					"id": dayName,
					"en": "Eid al-Fitr Day " + string(rune('1'+i)),
				},
			)
		}
	}

	// Eid al-Adha (Idul Adha)
	iduAdha := p.getIslamicHolidayDate(year, "idul_adha")
	if !iduAdha.IsZero() {
		holidays[iduAdha] = p.CreateHoliday(
			"Hari Raya Idul Adha",
			iduAdha,
			"islamic",
			map[string]string{
				"id": "Hari Raya Idul Adha",
				"en": "Eid al-Adha",
			},
		)
	}
}

// addChristianHolidays adds Christian holidays
func (p *IDProvider) addChristianHolidays(holidays map[time.Time]*Holiday, year int) {
	// Good Friday
	easter := EasterSunday(year)
	goodFriday := easter.AddDate(0, 0, -2)
	holidays[goodFriday] = p.CreateHoliday(
		"Wafat Isa Al Masih",
		goodFriday,
		"christian",
		map[string]string{
			"id": "Wafat Isa Al Masih",
			"en": "Good Friday",
		},
	)

	// Christmas Day
	holidays[time.Date(year, 12, 25, 0, 0, 0, 0, time.UTC)] = p.CreateHoliday(
		"Hari Raya Natal",
		time.Date(year, 12, 25, 0, 0, 0, 0, time.UTC),
		"christian",
		map[string]string{
			"id": "Hari Raya Natal",
			"en": "Christmas Day",
		},
	)

	// Ascension Day
	ascension := easter.AddDate(0, 0, 39)
	holidays[ascension] = p.CreateHoliday(
		"Kenaikan Isa Al Masih",
		ascension,
		"christian",
		map[string]string{
			"id": "Kenaikan Isa Al Masih",
			"en": "Ascension Day",
		},
	)
}

// addOtherReligiousHolidays adds Buddhist, Hindu, and Chinese holidays
func (p *IDProvider) addOtherReligiousHolidays(holidays map[time.Time]*Holiday, year int) {
	// Chinese New Year (varies each year)
	chineseNewYear := p.getChineseNewYearDate(year)
	if !chineseNewYear.IsZero() {
		holidays[chineseNewYear] = p.CreateHoliday(
			"Tahun Baru Imlek",
			chineseNewYear,
			"chinese",
			map[string]string{
				"id": "Tahun Baru Imlek",
				"en": "Chinese New Year",
			},
		)
	}

	// Vesak Day (Buddhist holiday - varies each year)
	vesak := p.getVesakDate(year)
	if !vesak.IsZero() {
		holidays[vesak] = p.CreateHoliday(
			"Hari Raya Waisak",
			vesak,
			"buddhist",
			map[string]string{
				"id": "Hari Raya Waisak",
				"en": "Vesak Day",
			},
		)
	}

	// Hindu New Year (Nyepi) - varies each year
	nyepi := p.getNyepiDate(year)
	if !nyepi.IsZero() {
		holidays[nyepi] = p.CreateHoliday(
			"Hari Raya Nyepi",
			nyepi,
			"hindu",
			map[string]string{
				"id": "Hari Raya Nyepi",
				"en": "Day of Silence (Nyepi)",
			},
		)
	}
}

// getIslamicHolidayDate returns approximate dates for Islamic holidays
func (p *IDProvider) getIslamicHolidayDate(year int, holiday string) time.Time {
	islamicHolidays := map[int]map[string]time.Time{
		2024: {
			"islamic_new_year": time.Date(2024, 7, 7, 0, 0, 0, 0, time.UTC),
			"maulid":           time.Date(2024, 9, 15, 0, 0, 0, 0, time.UTC),
			"isra_miraj":       time.Date(2024, 2, 8, 0, 0, 0, 0, time.UTC),
			"idul_fitri":       time.Date(2024, 4, 10, 0, 0, 0, 0, time.UTC),
			"idul_adha":        time.Date(2024, 6, 16, 0, 0, 0, 0, time.UTC),
		},
		2025: {
			"islamic_new_year": time.Date(2025, 6, 26, 0, 0, 0, 0, time.UTC),
			"maulid":           time.Date(2025, 9, 5, 0, 0, 0, 0, time.UTC),
			"isra_miraj":       time.Date(2025, 1, 28, 0, 0, 0, 0, time.UTC),
			"idul_fitri":       time.Date(2025, 3, 30, 0, 0, 0, 0, time.UTC),
			"idul_adha":        time.Date(2025, 6, 6, 0, 0, 0, 0, time.UTC),
		},
	}

	if yearHolidays, exists := islamicHolidays[year]; exists {
		if date, exists := yearHolidays[holiday]; exists {
			return date
		}
	}
	return time.Time{}
}

// getChineseNewYearDate returns Chinese New Year date for the year
func (p *IDProvider) getChineseNewYearDate(year int) time.Time {
	chineseNewYearDates := map[int]time.Time{
		2024: time.Date(2024, 2, 10, 0, 0, 0, 0, time.UTC),
		2025: time.Date(2025, 1, 29, 0, 0, 0, 0, time.UTC),
		2026: time.Date(2026, 2, 17, 0, 0, 0, 0, time.UTC),
		2027: time.Date(2027, 2, 6, 0, 0, 0, 0, time.UTC),
	}

	if date, exists := chineseNewYearDates[year]; exists {
		return date
	}
	return time.Time{}
}

// getVesakDate returns Vesak Day date for the year
func (p *IDProvider) getVesakDate(year int) time.Time {
	vesakDates := map[int]time.Time{
		2024: time.Date(2024, 5, 23, 0, 0, 0, 0, time.UTC),
		2025: time.Date(2025, 5, 12, 0, 0, 0, 0, time.UTC),
		2026: time.Date(2026, 5, 31, 0, 0, 0, 0, time.UTC),
		2027: time.Date(2027, 5, 20, 0, 0, 0, 0, time.UTC),
	}

	if date, exists := vesakDates[year]; exists {
		return date
	}
	return time.Time{}
}

// getNyepiDate returns Nyepi (Hindu New Year) date for the year
func (p *IDProvider) getNyepiDate(year int) time.Time {
	nyepiDates := map[int]time.Time{
		2024: time.Date(2024, 3, 11, 0, 0, 0, 0, time.UTC),
		2025: time.Date(2025, 3, 29, 0, 0, 0, 0, time.UTC),
		2026: time.Date(2026, 3, 19, 0, 0, 0, 0, time.UTC),
		2027: time.Date(2027, 3, 9, 0, 0, 0, 0, time.UTC),
	}

	if date, exists := nyepiDates[year]; exists {
		return date
	}
	return time.Time{}
}
