package countries

import (
	"time"
)

// PLProvider implements holiday calculations for Poland
type PLProvider struct {
	*BaseProvider
}

// NewPLProvider creates a new Polish holiday provider
func NewPLProvider() *PLProvider {
	base := NewBaseProvider("PL")
	base.subdivisions = []string{
		"DS", "KP", "LB", "LD", "LU", "MA", "MZ", "OP", "PK", "PD", "PM", "SL", "SK", "WN", "WP", "ZP",
		// Lower Silesian, Kuyavian-Pomeranian, Lublin, Łódź, Lubusz, Lesser Poland, Masovian,
		// Opole, Subcarpathian, Podlaskie, Pomeranian, Silesian, Świętokrzyskie, Warmian-Masurian,
		// Greater Poland, West Pomeranian
	}
	base.categories = []string{"public", "religious", "national"}

	return &PLProvider{BaseProvider: base}
}

// LoadHolidays loads all Polish holidays for a given year
func (pl *PLProvider) LoadHolidays(year int) map[time.Time]*Holiday {
	holidays := make(map[time.Time]*Holiday)

	// Fixed date holidays

	// New Year's Day - January 1
	newYear := time.Date(year, 1, 1, 0, 0, 0, 0, time.UTC)
	holidays[newYear] = pl.CreateHoliday(
		"Nowy Rok",
		newYear,
		"public",
		map[string]string{
			"pl": "Nowy Rok",
			"en": "New Year's Day",
		},
	)

	// Epiphany - January 6 (public holiday since 2011)
	if year >= 2011 {
		epiphany := time.Date(year, 1, 6, 0, 0, 0, 0, time.UTC)
		holidays[epiphany] = pl.CreateHoliday(
			"Święto Trzech Króli",
			epiphany,
			"religious",
			map[string]string{
				"pl": "Święto Trzech Króli",
				"en": "Epiphany",
			},
		)
	}

	// Labour Day - May 1
	labourDay := time.Date(year, 5, 1, 0, 0, 0, 0, time.UTC)
	holidays[labourDay] = pl.CreateHoliday(
		"Święto Pracy",
		labourDay,
		"public",
		map[string]string{
			"pl": "Święto Pracy",
			"en": "Labour Day",
		},
	)

	// Constitution Day - May 3
	constitutionDay := time.Date(year, 5, 3, 0, 0, 0, 0, time.UTC)
	holidays[constitutionDay] = pl.CreateHoliday(
		"Święto Konstytucji 3 Maja",
		constitutionDay,
		"national",
		map[string]string{
			"pl": "Święto Konstytucji 3 Maja",
			"en": "Constitution Day",
		},
	)

	// Assumption of Mary - August 15
	assumption := time.Date(year, 8, 15, 0, 0, 0, 0, time.UTC)
	holidays[assumption] = pl.CreateHoliday(
		"Wniebowzięcie Najświętszej Maryi Panny",
		assumption,
		"religious",
		map[string]string{
			"pl": "Wniebowzięcie Najświętszej Maryi Panny",
			"en": "Assumption of Mary",
		},
	)

	// All Saints' Day - November 1
	allSaints := time.Date(year, 11, 1, 0, 0, 0, 0, time.UTC)
	holidays[allSaints] = pl.CreateHoliday(
		"Wszystkich Świętych",
		allSaints,
		"religious",
		map[string]string{
			"pl": "Wszystkich Świętych",
			"en": "All Saints' Day",
		},
	)

	// Independence Day - November 11
	independenceDay := time.Date(year, 11, 11, 0, 0, 0, 0, time.UTC)
	holidays[independenceDay] = pl.CreateHoliday(
		"Narodowe Święto Niepodległości",
		independenceDay,
		"national",
		map[string]string{
			"pl": "Narodowe Święto Niepodległości",
			"en": "Independence Day",
		},
	)

	// Christmas Eve - December 24 (not a public holiday but widely observed)
	christmasEve := time.Date(year, 12, 24, 0, 0, 0, 0, time.UTC)
	holidays[christmasEve] = pl.CreateHoliday(
		"Wigilia",
		christmasEve,
		"religious",
		map[string]string{
			"pl": "Wigilia",
			"en": "Christmas Eve",
		},
	)

	// Christmas Day - December 25
	christmas := time.Date(year, 12, 25, 0, 0, 0, 0, time.UTC)
	holidays[christmas] = pl.CreateHoliday(
		"Boże Narodzenie",
		christmas,
		"public",
		map[string]string{
			"pl": "Boże Narodzenie",
			"en": "Christmas Day",
		},
	)

	// Boxing Day - December 26
	boxingDay := time.Date(year, 12, 26, 0, 0, 0, 0, time.UTC)
	holidays[boxingDay] = pl.CreateHoliday(
		"Drugi dzień Bożego Narodzenia",
		boxingDay,
		"public",
		map[string]string{
			"pl": "Drugi dzień Bożego Narodzenia",
			"en": "Boxing Day",
		},
	)

	// Easter-based holidays
	easter := EasterSunday(year)

	// Easter Sunday
	holidays[easter] = pl.CreateHoliday(
		"Niedziela Wielkanocna",
		easter,
		"religious",
		map[string]string{
			"pl": "Niedziela Wielkanocna",
			"en": "Easter Sunday",
		},
	)

	// Easter Monday
	easterMonday := easter.AddDate(0, 0, 1)
	holidays[easterMonday] = pl.CreateHoliday(
		"Poniedziałek Wielkanocny",
		easterMonday,
		"public",
		map[string]string{
			"pl": "Poniedziałek Wielkanocny",
			"en": "Easter Monday",
		},
	)

	// Whit Sunday (49 days after Easter)
	whitSunday := easter.AddDate(0, 0, 49)
	holidays[whitSunday] = pl.CreateHoliday(
		"Zielone Świątki",
		whitSunday,
		"religious",
		map[string]string{
			"pl": "Zielone Świątki",
			"en": "Whit Sunday",
		},
	)

	// Corpus Christi (60 days after Easter) - public holiday
	corpusChristi := easter.AddDate(0, 0, 60)
	holidays[corpusChristi] = pl.CreateHoliday(
		"Boże Ciało",
		corpusChristi,
		"public",
		map[string]string{
			"pl": "Boże Ciało",
			"en": "Corpus Christi",
		},
	)

	return holidays
}

// GetRegionalHolidays returns voivodeship-specific holidays
func (pl *PLProvider) GetRegionalHolidays(year int, voivodeships []string) map[time.Time]*Holiday {
	holidays := make(map[time.Time]*Holiday)

	for _, voivodeship := range voivodeships {
		switch voivodeship {
		case "SL": // Silesian Voivodeship
			// Silesian Uprising Day - May 3 (regional observance)
			uprising := time.Date(year, 5, 3, 0, 0, 0, 0, time.UTC)
			holidays[uprising] = pl.CreateHoliday(
				"Dzień Powstań Śląskich",
				uprising,
				"regional",
				map[string]string{
					"pl": "Dzień Powstań Śląskich",
					"en": "Silesian Uprisings Day",
				},
			)

		case "MA": // Lesser Poland (Małopolska)
			// Saint Stanislaus Day - May 8 (regional patron saint)
			stanislaus := time.Date(year, 5, 8, 0, 0, 0, 0, time.UTC)
			holidays[stanislaus] = pl.CreateHoliday(
				"Święty Stanisław",
				stanislaus,
				"regional",
				map[string]string{
					"pl": "Święty Stanisław",
					"en": "Saint Stanislaus Day",
				},
			)

		case "PD": // Podlaskie Voivodeship
			// Blessed Virgin Mary Day - September 8
			maryDay := time.Date(year, 9, 8, 0, 0, 0, 0, time.UTC)
			holidays[maryDay] = pl.CreateHoliday(
				"Narodzenie Najświętszej Maryi Panny",
				maryDay,
				"regional",
				map[string]string{
					"pl": "Narodzenie Najświętszej Maryi Panny",
					"en": "Nativity of Mary",
				},
			)
		}
	}

	return holidays
}

// GetSpecialObservances returns non-public observances
func (pl *PLProvider) GetSpecialObservances(year int) map[time.Time]*Holiday {
	holidays := make(map[time.Time]*Holiday)

	// Fat Thursday (52 days before Easter)
	easter := EasterSunday(year)
	fatThursday := easter.AddDate(0, 0, -52)
	holidays[fatThursday] = pl.CreateHoliday(
		"Tłusty Czwartek",
		fatThursday,
		"cultural",
		map[string]string{
			"pl": "Tłusty Czwartek",
			"en": "Fat Thursday",
		},
	)

	// Ash Wednesday (46 days before Easter)
	ashWednesday := easter.AddDate(0, 0, -46)
	holidays[ashWednesday] = pl.CreateHoliday(
		"Środa Popielcowa",
		ashWednesday,
		"religious",
		map[string]string{
			"pl": "Środa Popielcowa",
			"en": "Ash Wednesday",
		},
	)

	// Palm Sunday (7 days before Easter)
	palmSunday := easter.AddDate(0, 0, -7)
	holidays[palmSunday] = pl.CreateHoliday(
		"Niedziela Palmowa",
		palmSunday,
		"religious",
		map[string]string{
			"pl": "Niedziela Palmowa",
			"en": "Palm Sunday",
		},
	)

	// Good Friday (2 days before Easter)
	goodFriday := easter.AddDate(0, 0, -2)
	holidays[goodFriday] = pl.CreateHoliday(
		"Wielki Piątek",
		goodFriday,
		"religious",
		map[string]string{
			"pl": "Wielki Piątek",
			"en": "Good Friday",
		},
	)

	// Saint Nicholas Day - December 6
	stNicholas := time.Date(year, 12, 6, 0, 0, 0, 0, time.UTC)
	holidays[stNicholas] = pl.CreateHoliday(
		"Mikołajki",
		stNicholas,
		"cultural",
		map[string]string{
			"pl": "Mikołajki",
			"en": "Saint Nicholas Day",
		},
	)

	// New Year's Eve - December 31
	newYearEve := time.Date(year, 12, 31, 0, 0, 0, 0, time.UTC)
	holidays[newYearEve] = pl.CreateHoliday(
		"Sylwester",
		newYearEve,
		"cultural",
		map[string]string{
			"pl": "Sylwester",
			"en": "New Year's Eve",
		},
	)

	// Warsaw Uprising Day - August 1 (commemorative)
	if year >= 1944 {
		warsawUprising := time.Date(year, 8, 1, 0, 0, 0, 0, time.UTC)
		holidays[warsawUprising] = pl.CreateHoliday(
			"Dzień Powstania Warszawskiego",
			warsawUprising,
			"commemorative",
			map[string]string{
				"pl": "Dzień Powstania Warszawskiego",
				"en": "Warsaw Uprising Day",
			},
		)
	}

	// Holocaust Remembrance Day - April 19 (commemorative)
	holocaustDay := time.Date(year, 4, 19, 0, 0, 0, 0, time.UTC)
	holidays[holocaustDay] = pl.CreateHoliday(
		"Dzień Pamięci o Holokauście",
		holocaustDay,
		"commemorative",
		map[string]string{
			"pl": "Dzień Pamięci o Holokauście",
			"en": "Holocaust Remembrance Day",
		},
	)

	return holidays
}

