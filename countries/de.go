package countries

import (
	"time"
)

// DEProvider implements holiday calculations for Germany
type DEProvider struct {
	*BaseProvider
}

// NewDEProvider creates a new German holiday provider
func NewDEProvider() *DEProvider {
	base := NewBaseProvider("DE")
	base.subdivisions = []string{
		"BW", "BY", "BE", "BB", "HB", "HH", "HE", "MV", "NI", "NW", "RP", "SL", "SN", "ST", "SH", "TH",
		// Baden-Württemberg, Bayern, Berlin, Brandenburg, Bremen, Hamburg, Hessen, 
		// Mecklenburg-Vorpommern, Niedersachsen, Nordrhein-Westfalen, Rheinland-Pfalz, 
		// Saarland, Sachsen, Sachsen-Anhalt, Schleswig-Holstein, Thüringen
	}
	base.categories = []string{"public", "religious", "regional"}
	
	return &DEProvider{BaseProvider: base}
}

// LoadHolidays loads all German holidays for a given year
func (de *DEProvider) LoadHolidays(year int) map[time.Time]*Holiday {
	holidays := make(map[time.Time]*Holiday)
	
	// Fixed date holidays
	
	// New Year's Day - January 1
	newYear := time.Date(year, 1, 1, 0, 0, 0, 0, time.UTC)
	holidays[newYear] = de.CreateHoliday(
		"Neujahr",
		newYear,
		"public",
		map[string]string{
			"de": "Neujahr",
			"en": "New Year's Day",
		},
	)
	
	// Epiphany - January 6 (BW, BY, ST)
	epiphany := time.Date(year, 1, 6, 0, 0, 0, 0, time.UTC)
	holidays[epiphany] = de.CreateHoliday(
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
	holidays[labourDay] = de.CreateHoliday(
		"Tag der Arbeit",
		labourDay,
		"public",
		map[string]string{
			"de": "Tag der Arbeit",
			"en": "Labour Day",
		},
	)
	
	// German Unity Day - October 3 (since 1990)
	if year >= 1990 {
		unityDay := time.Date(year, 10, 3, 0, 0, 0, 0, time.UTC)
		holidays[unityDay] = de.CreateHoliday(
			"Tag der Deutschen Einheit",
			unityDay,
			"public",
			map[string]string{
				"de": "Tag der Deutschen Einheit",
				"en": "German Unity Day",
			},
		)
	}
	
	// All Saints' Day - November 1 (BW, BY, NW, RP, SL)
	allSaints := time.Date(year, 11, 1, 0, 0, 0, 0, time.UTC)
	holidays[allSaints] = de.CreateHoliday(
		"Allerheiligen",
		allSaints,
		"religious",
		map[string]string{
			"de": "Allerheiligen",
			"en": "All Saints' Day",
		},
	)
	
	// Christmas Eve - December 24
	christmasEve := time.Date(year, 12, 24, 0, 0, 0, 0, time.UTC)
	holidays[christmasEve] = de.CreateHoliday(
		"Heiligabend",
		christmasEve,
		"religious",
		map[string]string{
			"de": "Heiligabend",
			"en": "Christmas Eve",
		},
	)
	
	// Christmas Day - December 25
	christmas := time.Date(year, 12, 25, 0, 0, 0, 0, time.UTC)
	holidays[christmas] = de.CreateHoliday(
		"1. Weihnachtsfeiertag",
		christmas,
		"public",
		map[string]string{
			"de": "1. Weihnachtsfeiertag",
			"en": "Christmas Day",
		},
	)
	
	// Boxing Day - December 26
	boxingDay := time.Date(year, 12, 26, 0, 0, 0, 0, time.UTC)
	holidays[boxingDay] = de.CreateHoliday(
		"2. Weihnachtsfeiertag",
		boxingDay,
		"public",
		map[string]string{
			"de": "2. Weihnachtsfeiertag",
			"en": "Boxing Day",
		},
	)
	
	// New Year's Eve - December 31
	newYearEve := time.Date(year, 12, 31, 0, 0, 0, 0, time.UTC)
	holidays[newYearEve] = de.CreateHoliday(
		"Silvester",
		newYearEve,
		"public",
		map[string]string{
			"de": "Silvester",
			"en": "New Year's Eve",
		},
	)
	
	// Easter-based holidays
	easter := EasterSunday(year)
	
	// Good Friday
	goodFriday := easter.AddDate(0, 0, -2)
	holidays[goodFriday] = de.CreateHoliday(
		"Karfreitag",
		goodFriday,
		"public",
		map[string]string{
			"de": "Karfreitag",
			"en": "Good Friday",
		},
	)
	
	// Easter Sunday
	holidays[easter] = de.CreateHoliday(
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
	holidays[easterMonday] = de.CreateHoliday(
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
	holidays[ascension] = de.CreateHoliday(
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
	holidays[whitSunday] = de.CreateHoliday(
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
	holidays[whitMonday] = de.CreateHoliday(
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
	holidays[corpusChristi] = de.CreateHoliday(
		"Fronleichnam",
		corpusChristi,
		"religious",
		map[string]string{
			"de": "Fronleichnam",
			"en": "Corpus Christi",
		},
	)
	
	return holidays
}

// GetRegionalHolidays returns state-specific holidays
func (de *DEProvider) GetRegionalHolidays(year int, states []string) map[time.Time]*Holiday {
	holidays := make(map[time.Time]*Holiday)
	
	for _, state := range states {
		switch state {
		case "BY": // Bavaria
			// Assumption of Mary - August 15
			assumption := time.Date(year, 8, 15, 0, 0, 0, 0, time.UTC)
			holidays[assumption] = de.CreateHoliday(
				"Mariä Himmelfahrt",
				assumption,
				"religious",
				map[string]string{
					"de": "Mariä Himmelfahrt",
					"en": "Assumption of Mary",
				},
			)
		}
		
		// Reformation Day for Protestant states
		if state == "SN" || state == "ST" || state == "TH" || state == "BB" || state == "MV" {
			reformation := time.Date(year, 10, 31, 0, 0, 0, 0, time.UTC)
			holidays[reformation] = de.CreateHoliday(
				"Reformationstag",
				reformation,
				"religious",
				map[string]string{
					"de": "Reformationstag",
					"en": "Reformation Day",
				},
			)
		}
		
		// Repentance and Prayer Day for Protestant states (only SN still observes it)
		if state == "SN" {
			repentance := de.getRepentanceDay(year)
			holidays[repentance] = de.CreateHoliday(
				"Buß- und Bettag",
				repentance,
				"religious",
				map[string]string{
					"de": "Buß- und Bettag",
					"en": "Repentance and Prayer Day",
				},
			)
		}
	}
	
	return holidays
}

// getRepentanceDay calculates Repentance and Prayer Day (Buß- und Bettag)
// It's the Wednesday before November 23
func (de *DEProvider) getRepentanceDay(year int) time.Time {
	nov23 := time.Date(year, 11, 23, 0, 0, 0, 0, time.UTC)
	
	// Find the Wednesday before November 23
	daysBack := int(nov23.Weekday()) - int(time.Wednesday)
	if daysBack <= 0 {
		daysBack += 7
	}
	
	return nov23.AddDate(0, 0, -daysBack)
}

// GetSpecialObservances returns non-public observances
func (de *DEProvider) GetSpecialObservances(year int) map[time.Time]*Holiday {
	holidays := make(map[time.Time]*Holiday)
	
	// Carnival Monday (48 days before Easter)
	easter := EasterSunday(year)
	carnivalMonday := easter.AddDate(0, 0, -48)
	holidays[carnivalMonday] = de.CreateHoliday(
		"Rosenmontag",
		carnivalMonday,
		"regional",
		map[string]string{
			"de": "Rosenmontag",
			"en": "Carnival Monday",
		},
	)
	
	// Shrove Tuesday (47 days before Easter)
	shroveTuesday := easter.AddDate(0, 0, -47)
	holidays[shroveTuesday] = de.CreateHoliday(
		"Fastnacht",
		shroveTuesday,
		"regional",
		map[string]string{
			"de": "Fastnacht",
			"en": "Shrove Tuesday",
		},
	)
	
	// Ash Wednesday (46 days before Easter)
	ashWednesday := easter.AddDate(0, 0, -46)
	holidays[ashWednesday] = de.CreateHoliday(
		"Aschermittwoch",
		ashWednesday,
		"religious",
		map[string]string{
			"de": "Aschermittwoch",
			"en": "Ash Wednesday",
		},
	)
	
	// World War II Remembrance - May 8
	if year >= 1945 {
		remembrance := time.Date(year, 5, 8, 0, 0, 0, 0, time.UTC)
		holidays[remembrance] = de.CreateHoliday(
			"Tag der Befreiung",
			remembrance,
			"regional",
			map[string]string{
				"de": "Tag der Befreiung",
				"en": "Liberation Day",
			},
		)
	}
	
	return holidays
}
