package countries

import (
	"time"
)

// FRProvider implements holiday calculations for France
type FRProvider struct {
	*BaseProvider
}

// NewFRProvider creates a new French holiday provider
func NewFRProvider() *FRProvider {
	base := NewBaseProvider("FR")
	base.subdivisions = []string{
		"ARA", "BFC", "BRE", "CVL", "COR", "GES", "HDF", "IDF", "NOR", "NAQ", "OCC", "PDL", "PAC",
		"GP", "GF", "MQ", "RE", "YT", "NC", "PF", "BL", "MF", "PM", "WF",
		// Metropolitan regions + Overseas territories
		"67", "68", // Bas-Rhin, Haut-Rhin (Alsace-Moselle)
		"57", // Moselle
	}
	base.categories = []string{"public", "religious", "regional", "secular"}

	return &FRProvider{BaseProvider: base}
}

// LoadHolidays loads all French holidays for a given year
func (fr *FRProvider) LoadHolidays(year int) map[time.Time]*Holiday {
	holidays := make(map[time.Time]*Holiday)

	// Fixed date holidays

	// New Year's Day - January 1
	newYear := time.Date(year, 1, 1, 0, 0, 0, 0, time.UTC)
	holidays[newYear] = fr.CreateHoliday(
		"Jour de l'An",
		newYear,
		"public",
		map[string]string{
			"fr": "Jour de l'An",
			"en": "New Year's Day",
		},
	)

	// Labour Day - May 1
	labourDay := time.Date(year, 5, 1, 0, 0, 0, 0, time.UTC)
	holidays[labourDay] = fr.CreateHoliday(
		"Fête du Travail",
		labourDay,
		"public",
		map[string]string{
			"fr": "Fête du Travail",
			"en": "Labour Day",
		},
	)

	// Victory in Europe Day - May 8
	victoryDay := time.Date(year, 5, 8, 0, 0, 0, 0, time.UTC)
	holidays[victoryDay] = fr.CreateHoliday(
		"Fête de la Victoire",
		victoryDay,
		"public",
		map[string]string{
			"fr": "Fête de la Victoire",
			"en": "Victory in Europe Day",
		},
	)

	// Bastille Day - July 14
	bastilleDay := time.Date(year, 7, 14, 0, 0, 0, 0, time.UTC)
	holidays[bastilleDay] = fr.CreateHoliday(
		"Fête nationale",
		bastilleDay,
		"public",
		map[string]string{
			"fr": "Fête nationale",
			"en": "Bastille Day",
		},
	)

	// Assumption of Mary - August 15
	assumption := time.Date(year, 8, 15, 0, 0, 0, 0, time.UTC)
	holidays[assumption] = fr.CreateHoliday(
		"Assomption",
		assumption,
		"religious",
		map[string]string{
			"fr": "Assomption",
			"en": "Assumption of Mary",
		},
	)

	// All Saints' Day - November 1
	allSaints := time.Date(year, 11, 1, 0, 0, 0, 0, time.UTC)
	holidays[allSaints] = fr.CreateHoliday(
		"Toussaint",
		allSaints,
		"religious",
		map[string]string{
			"fr": "Toussaint",
			"en": "All Saints' Day",
		},
	)

	// Armistice Day - November 11
	armistice := time.Date(year, 11, 11, 0, 0, 0, 0, time.UTC)
	holidays[armistice] = fr.CreateHoliday(
		"Armistice",
		armistice,
		"public",
		map[string]string{
			"fr": "Armistice",
			"en": "Armistice Day",
		},
	)

	// Christmas Day - December 25
	christmas := time.Date(year, 12, 25, 0, 0, 0, 0, time.UTC)
	holidays[christmas] = fr.CreateHoliday(
		"Noël",
		christmas,
		"religious",
		map[string]string{
			"fr": "Noël",
			"en": "Christmas Day",
		},
	)

	// Easter-based holidays
	easter := EasterSunday(year)

	// Easter Sunday
	holidays[easter] = fr.CreateHoliday(
		"Pâques",
		easter,
		"religious",
		map[string]string{
			"fr": "Pâques",
			"en": "Easter Sunday",
		},
	)

	// Easter Monday
	easterMonday := easter.AddDate(0, 0, 1)
	holidays[easterMonday] = fr.CreateHoliday(
		"Lundi de Pâques",
		easterMonday,
		"religious",
		map[string]string{
			"fr": "Lundi de Pâques",
			"en": "Easter Monday",
		},
	)

	// Ascension Day (39 days after Easter)
	ascension := easter.AddDate(0, 0, 39)
	holidays[ascension] = fr.CreateHoliday(
		"Ascension",
		ascension,
		"religious",
		map[string]string{
			"fr": "Ascension",
			"en": "Ascension Day",
		},
	)

	// Whit Sunday (49 days after Easter)
	whitSunday := easter.AddDate(0, 0, 49)
	holidays[whitSunday] = fr.CreateHoliday(
		"Pentecôte",
		whitSunday,
		"religious",
		map[string]string{
			"fr": "Pentecôte",
			"en": "Whit Sunday",
		},
	)

	// Whit Monday (50 days after Easter)
	whitMonday := easter.AddDate(0, 0, 50)
	holidays[whitMonday] = fr.CreateHoliday(
		"Lundi de Pentecôte",
		whitMonday,
		"religious",
		map[string]string{
			"fr": "Lundi de Pentecôte",
			"en": "Whit Monday",
		},
	)

	return holidays
}

// GetRegionalHolidays returns region-specific holidays
func (fr *FRProvider) GetRegionalHolidays(year int, regions []string) map[time.Time]*Holiday {
	holidays := make(map[time.Time]*Holiday)

	for _, region := range regions {
		// Alsace-Moselle specific holidays (departments 57, 67, 68)
		if region == "57" || region == "67" || region == "68" {
			// Good Friday (Alsace-Moselle only)
			easter := EasterSunday(year)
			goodFriday := easter.AddDate(0, 0, -2)
			holidays[goodFriday] = fr.CreateHoliday(
				"Vendredi saint",
				goodFriday,
				"regional",
				map[string]string{
					"fr": "Vendredi saint",
					"en": "Good Friday",
				},
			)

			// St. Stephen's Day - December 26 (Alsace-Moselle only)
			stStephen := time.Date(year, 12, 26, 0, 0, 0, 0, time.UTC)
			holidays[stStephen] = fr.CreateHoliday(
				"Saint-Étienne",
				stStephen,
				"regional",
				map[string]string{
					"fr": "Saint-Étienne",
					"en": "St. Stephen's Day",
				},
			)
		}

		// Overseas territories specific holidays
		switch region {
		case "GP", "MQ": // Guadeloupe, Martinique
			// Slavery Abolition Day - May 27 (Guadeloupe) / May 22 (Martinique)
			var abolitionDay time.Time
			var name string
			if region == "GP" {
				abolitionDay = time.Date(year, 5, 27, 0, 0, 0, 0, time.UTC)
				name = "Abolition de l'esclavage (Guadeloupe)"
			} else {
				abolitionDay = time.Date(year, 5, 22, 0, 0, 0, 0, time.UTC)
				name = "Abolition de l'esclavage (Martinique)"
			}
			holidays[abolitionDay] = fr.CreateHoliday(
				name,
				abolitionDay,
				"regional",
				map[string]string{
					"fr": name,
					"en": "Slavery Abolition Day",
				},
			)

		case "GF": // French Guiana
			// Slavery Abolition Day - June 10
			abolition := time.Date(year, 6, 10, 0, 0, 0, 0, time.UTC)
			holidays[abolition] = fr.CreateHoliday(
				"Abolition de l'esclavage (Guyane)",
				abolition,
				"regional",
				map[string]string{
					"fr": "Abolition de l'esclavage (Guyane)",
					"en": "Slavery Abolition Day",
				},
			)

		case "RE": // Réunion
			// Slavery Abolition Day - December 20
			abolition := time.Date(year, 12, 20, 0, 0, 0, 0, time.UTC)
			holidays[abolition] = fr.CreateHoliday(
				"Abolition de l'esclavage (Réunion)",
				abolition,
				"regional",
				map[string]string{
					"fr": "Abolition de l'esclavage (Réunion)",
					"en": "Slavery Abolition Day",
				},
			)

		case "YT": // Mayotte
			// Mayotte Day - March 31
			mayotteDay := time.Date(year, 3, 31, 0, 0, 0, 0, time.UTC)
			holidays[mayotteDay] = fr.CreateHoliday(
				"Journée de Mayotte",
				mayotteDay,
				"regional",
				map[string]string{
					"fr": "Journée de Mayotte",
					"en": "Mayotte Day",
				},
			)
		}
	}

	return holidays
}

// GetSecularObservances returns non-religious observances
func (fr *FRProvider) GetSecularObservances(year int) map[time.Time]*Holiday {
	holidays := make(map[time.Time]*Holiday)

	// Epiphany - January 6
	epiphany := time.Date(year, 1, 6, 0, 0, 0, 0, time.UTC)
	holidays[epiphany] = fr.CreateHoliday(
		"Épiphanie",
		epiphany,
		"secular",
		map[string]string{
			"fr": "Épiphanie",
			"en": "Epiphany",
		},
	)

	// Candlemas - February 2
	candlemas := time.Date(year, 2, 2, 0, 0, 0, 0, time.UTC)
	holidays[candlemas] = fr.CreateHoliday(
		"Chandeleur",
		candlemas,
		"secular",
		map[string]string{
			"fr": "Chandeleur",
			"en": "Candlemas",
		},
	)

	// Valentine's Day - February 14
	valentine := time.Date(year, 2, 14, 0, 0, 0, 0, time.UTC)
	holidays[valentine] = fr.CreateHoliday(
		"Saint-Valentin",
		valentine,
		"secular",
		map[string]string{
			"fr": "Saint-Valentin",
			"en": "Valentine's Day",
		},
	)

	// Mother's Day - Last Sunday of May (or first Sunday of June if coincides with Whit Sunday)
	mothersDay := fr.getMothersDay(year)
	holidays[mothersDay] = fr.CreateHoliday(
		"Fête des Mères",
		mothersDay,
		"secular",
		map[string]string{
			"fr": "Fête des Mères",
			"en": "Mother's Day",
		},
	)

	// Father's Day - Third Sunday of June
	fathersDay := NthWeekdayOfMonth(year, 6, time.Sunday, 3)
	holidays[fathersDay] = fr.CreateHoliday(
		"Fête des Pères",
		fathersDay,
		"secular",
		map[string]string{
			"fr": "Fête des Pères",
			"en": "Father's Day",
		},
	)

	// Music Day - June 21
	musicDay := time.Date(year, 6, 21, 0, 0, 0, 0, time.UTC)
	holidays[musicDay] = fr.CreateHoliday(
		"Fête de la Musique",
		musicDay,
		"secular",
		map[string]string{
			"fr": "Fête de la Musique",
			"en": "Music Day",
		},
	)

	// All Souls' Day - November 2
	allSouls := time.Date(year, 11, 2, 0, 0, 0, 0, time.UTC)
	holidays[allSouls] = fr.CreateHoliday(
		"Jour des Morts",
		allSouls,
		"secular",
		map[string]string{
			"fr": "Jour des Morts",
			"en": "All Souls' Day",
		},
	)

	return holidays
}

// getMothersDay calculates French Mother's Day
// Last Sunday of May, or first Sunday of June if it coincides with Whit Sunday
func (fr *FRProvider) getMothersDay(year int) time.Time {
	// Last Sunday of May
	lastSundayMay := NthWeekdayOfMonth(year, 5, time.Sunday, -1)

	// Check if it coincides with Whit Sunday (49 days after Easter)
	easter := EasterSunday(year)
	whitSunday := easter.AddDate(0, 0, 49)

	if lastSundayMay.Equal(whitSunday) {
		// Move to first Sunday of June
		return NthWeekdayOfMonth(year, 6, time.Sunday, 1)
	}

	return lastSundayMay
}
