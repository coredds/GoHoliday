package countries

import (
	"time"
)

// ATProvider implements holiday calculations for Austria
type ATProvider struct {
	*BaseProvider
}

// NewATProvider creates a new Austrian holiday provider
func NewATProvider() *ATProvider {
	base := NewBaseProvider("AT")
	base.subdivisions = []string{
		"1", "2", "3", "4", "5", "6", "7", "8", "9",
		// Burgenland, Kärnten, Niederösterreich, Oberösterreich, Salzburg,
		// Steiermark, Tirol, Vorarlberg, Wien
	}
	base.categories = []string{"public", "religious", "regional"}

	return &ATProvider{BaseProvider: base}
}

// LoadHolidays loads all Austrian holidays for a given year
func (at *ATProvider) LoadHolidays(year int) map[time.Time]*Holiday {
	holidays := make(map[time.Time]*Holiday)

	// Fixed date holidays

	// New Year's Day - January 1
	newYear := time.Date(year, 1, 1, 0, 0, 0, 0, time.UTC)
	holidays[newYear] = at.CreateHoliday(
		"Neujahr",
		newYear,
		"public",
		map[string]string{
			"de": "Neujahr",
			"en": "New Year's Day",
		},
	)

	// Epiphany - January 6
	epiphany := time.Date(year, 1, 6, 0, 0, 0, 0, time.UTC)
	holidays[epiphany] = at.CreateHoliday(
		"Heilige Drei Könige",
		epiphany,
		"religious",
		map[string]string{
			"de": "Heilige Drei Könige",
			"en": "Epiphany",
		},
	)

	// Labour Day - May 1
	labourDay := time.Date(year, 5, 1, 0, 0, 0, 0, time.UTC)
	holidays[labourDay] = at.CreateHoliday(
		"Staatsfeiertag",
		labourDay,
		"public",
		map[string]string{
			"de": "Staatsfeiertag",
			"en": "Labour Day",
		},
	)

	// Assumption of Mary - August 15
	assumption := time.Date(year, 8, 15, 0, 0, 0, 0, time.UTC)
	holidays[assumption] = at.CreateHoliday(
		"Mariä Himmelfahrt",
		assumption,
		"religious",
		map[string]string{
			"de": "Mariä Himmelfahrt",
			"en": "Assumption of Mary",
		},
	)

	// Austrian National Day - October 26 (since 1965)
	if year >= 1965 {
		nationalDay := time.Date(year, 10, 26, 0, 0, 0, 0, time.UTC)
		holidays[nationalDay] = at.CreateHoliday(
			"Nationalfeiertag",
			nationalDay,
			"public",
			map[string]string{
				"de": "Nationalfeiertag",
				"en": "Austrian National Day",
			},
		)
	}

	// All Saints' Day - November 1
	allSaints := time.Date(year, 11, 1, 0, 0, 0, 0, time.UTC)
	holidays[allSaints] = at.CreateHoliday(
		"Allerheiligen",
		allSaints,
		"religious",
		map[string]string{
			"de": "Allerheiligen",
			"en": "All Saints' Day",
		},
	)

	// Immaculate Conception - December 8
	immaculateConception := time.Date(year, 12, 8, 0, 0, 0, 0, time.UTC)
	holidays[immaculateConception] = at.CreateHoliday(
		"Mariä Empfängnis",
		immaculateConception,
		"religious",
		map[string]string{
			"de": "Mariä Empfängnis",
			"en": "Immaculate Conception",
		},
	)

	// Christmas Day - December 25
	christmas := time.Date(year, 12, 25, 0, 0, 0, 0, time.UTC)
	holidays[christmas] = at.CreateHoliday(
		"Christtag",
		christmas,
		"public",
		map[string]string{
			"de": "Christtag",
			"en": "Christmas Day",
		},
	)

	// Boxing Day - December 26
	boxingDay := time.Date(year, 12, 26, 0, 0, 0, 0, time.UTC)
	holidays[boxingDay] = at.CreateHoliday(
		"Stefanitag",
		boxingDay,
		"public",
		map[string]string{
			"de": "Stefanitag",
			"en": "Boxing Day",
		},
	)

	// Easter-based holidays
	easter := EasterSunday(year)

	// Easter Sunday
	holidays[easter] = at.CreateHoliday(
		"Ostersonntag",
		easter,
		"religious",
		map[string]string{
			"de": "Ostersonntag",
			"en": "Easter Sunday",
		},
	)

	// Easter Monday
	easterMonday := easter.AddDate(0, 0, 1)
	holidays[easterMonday] = at.CreateHoliday(
		"Ostermontag",
		easterMonday,
		"public",
		map[string]string{
			"de": "Ostermontag",
			"en": "Easter Monday",
		},
	)

	// Ascension Day (39 days after Easter)
	ascension := easter.AddDate(0, 0, 39)
	holidays[ascension] = at.CreateHoliday(
		"Christi Himmelfahrt",
		ascension,
		"public",
		map[string]string{
			"de": "Christi Himmelfahrt",
			"en": "Ascension Day",
		},
	)

	// Whit Sunday (49 days after Easter)
	whitSunday := easter.AddDate(0, 0, 49)
	holidays[whitSunday] = at.CreateHoliday(
		"Pfingstsonntag",
		whitSunday,
		"religious",
		map[string]string{
			"de": "Pfingstsonntag",
			"en": "Whit Sunday",
		},
	)

	// Whit Monday (50 days after Easter)
	whitMonday := easter.AddDate(0, 0, 50)
	holidays[whitMonday] = at.CreateHoliday(
		"Pfingstmontag",
		whitMonday,
		"public",
		map[string]string{
			"de": "Pfingstmontag",
			"en": "Whit Monday",
		},
	)

	// Corpus Christi (60 days after Easter)
	corpusChristi := easter.AddDate(0, 0, 60)
	holidays[corpusChristi] = at.CreateHoliday(
		"Fronleichnam",
		corpusChristi,
		"public",
		map[string]string{
			"de": "Fronleichnam",
			"en": "Corpus Christi",
		},
	)

	return holidays
}

// GetRegionalHolidays returns state-specific holidays
func (at *ATProvider) GetRegionalHolidays(year int, states []string) map[time.Time]*Holiday {
	holidays := make(map[time.Time]*Holiday)

	for _, state := range states {
		switch state {
		case "1": // Burgenland
			// Saint Martin's Day - November 11
			martinDay := time.Date(year, 11, 11, 0, 0, 0, 0, time.UTC)
			holidays[martinDay] = at.CreateHoliday(
				"Martinstag",
				martinDay,
				"regional",
				map[string]string{
					"de": "Martinstag",
					"en": "Saint Martin's Day",
				},
			)

		case "2": // Carinthia (Kärnten)
			// Carinthian Plebiscite Day - October 10
			plebiscite := time.Date(year, 10, 10, 0, 0, 0, 0, time.UTC)
			holidays[plebiscite] = at.CreateHoliday(
				"Tag der Volksabstimmung",
				plebiscite,
				"regional",
				map[string]string{
					"de": "Tag der Volksabstimmung",
					"en": "Carinthian Plebiscite Day",
				},
			)

		case "7": // Tyrol (Tirol)
			// Sacred Heart of Jesus - June 19 (Friday after Corpus Christi)
			easter := EasterSunday(year)
			corpusChristi := easter.AddDate(0, 0, 60)
			// Find the Friday after Corpus Christi
			daysToFriday := (5 - int(corpusChristi.Weekday()) + 7) % 7
			if daysToFriday == 0 {
				daysToFriday = 7
			}
			sacredHeart := corpusChristi.AddDate(0, 0, daysToFriday)
			holidays[sacredHeart] = at.CreateHoliday(
				"Herz-Jesu-Fest",
				sacredHeart,
				"regional",
				map[string]string{
					"de": "Herz-Jesu-Fest",
					"en": "Sacred Heart of Jesus",
				},
			)

		case "8": // Vorarlberg
			// Saint Joseph's Day - March 19
			josephDay := time.Date(year, 3, 19, 0, 0, 0, 0, time.UTC)
			holidays[josephDay] = at.CreateHoliday(
				"Josefitag",
				josephDay,
				"regional",
				map[string]string{
					"de": "Josefitag",
					"en": "Saint Joseph's Day",
				},
			)
		}
	}

	return holidays
}

// GetSpecialObservances returns non-public observances
func (at *ATProvider) GetSpecialObservances(year int) map[time.Time]*Holiday {
	holidays := make(map[time.Time]*Holiday)

	// Carnival Monday (48 days before Easter)
	easter := EasterSunday(year)
	carnivalMonday := easter.AddDate(0, 0, -48)
	holidays[carnivalMonday] = at.CreateHoliday(
		"Rosenmontag",
		carnivalMonday,
		"cultural",
		map[string]string{
			"de": "Rosenmontag",
			"en": "Carnival Monday",
		},
	)

	// Shrove Tuesday (47 days before Easter)
	shroveTuesday := easter.AddDate(0, 0, -47)
	holidays[shroveTuesday] = at.CreateHoliday(
		"Fasching",
		shroveTuesday,
		"cultural",
		map[string]string{
			"de": "Fasching",
			"en": "Shrove Tuesday",
		},
	)

	// Ash Wednesday (46 days before Easter)
	ashWednesday := easter.AddDate(0, 0, -46)
	holidays[ashWednesday] = at.CreateHoliday(
		"Aschermittwoch",
		ashWednesday,
		"religious",
		map[string]string{
			"de": "Aschermittwoch",
			"en": "Ash Wednesday",
		},
	)

	// Palm Sunday (7 days before Easter)
	palmSunday := easter.AddDate(0, 0, -7)
	holidays[palmSunday] = at.CreateHoliday(
		"Palmsonntag",
		palmSunday,
		"religious",
		map[string]string{
			"de": "Palmsonntag",
			"en": "Palm Sunday",
		},
	)

	// Good Friday (2 days before Easter)
	goodFriday := easter.AddDate(0, 0, -2)
	holidays[goodFriday] = at.CreateHoliday(
		"Karfreitag",
		goodFriday,
		"religious",
		map[string]string{
			"de": "Karfreitag",
			"en": "Good Friday",
		},
	)

	// Saint Nicholas Day - December 6
	stNicholas := time.Date(year, 12, 6, 0, 0, 0, 0, time.UTC)
	holidays[stNicholas] = at.CreateHoliday(
		"Nikolaus",
		stNicholas,
		"cultural",
		map[string]string{
			"de": "Nikolaus",
			"en": "Saint Nicholas Day",
		},
	)

	// Christmas Eve - December 24
	christmasEve := time.Date(year, 12, 24, 0, 0, 0, 0, time.UTC)
	holidays[christmasEve] = at.CreateHoliday(
		"Heiligabend",
		christmasEve,
		"cultural",
		map[string]string{
			"de": "Heiligabend",
			"en": "Christmas Eve",
		},
	)

	// New Year's Eve - December 31
	newYearEve := time.Date(year, 12, 31, 0, 0, 0, 0, time.UTC)
	holidays[newYearEve] = at.CreateHoliday(
		"Silvester",
		newYearEve,
		"cultural",
		map[string]string{
			"de": "Silvester",
			"en": "New Year's Eve",
		},
	)

	// Austrian State Treaty Day - May 15 (commemorative, since 1955)
	if year >= 1955 {
		stateTreaty := time.Date(year, 5, 15, 0, 0, 0, 0, time.UTC)
		holidays[stateTreaty] = at.CreateHoliday(
			"Staatsvertragsunterzeichnung",
			stateTreaty,
			"commemorative",
			map[string]string{
				"de": "Staatsvertragsunterzeichnung",
				"en": "Austrian State Treaty Day",
			},
		)
	}

	// Anschluss Remembrance Day - March 12 (commemorative)
	if year >= 1938 {
		anschluss := time.Date(year, 3, 12, 0, 0, 0, 0, time.UTC)
		holidays[anschluss] = at.CreateHoliday(
			"Gedenktag gegen Gewalt und Rassismus",
			anschluss,
			"commemorative",
			map[string]string{
				"de": "Gedenktag gegen Gewalt und Rassismus",
				"en": "Day against Violence and Racism",
			},
		)
	}

	return holidays
}
