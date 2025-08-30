package countries

import (
	"time"
)

// BEProvider implements holiday calculations for Belgium
type BEProvider struct {
	*BaseProvider
}

// NewBEProvider creates a new Belgian holiday provider
func NewBEProvider() *BEProvider {
	base := NewBaseProvider("BE")
	base.subdivisions = []string{
		"BRU", "VLG", "WAL",
		// Brussels-Capital Region, Flemish Region, Walloon Region
	}
	base.categories = []string{"public", "religious", "regional"}

	return &BEProvider{BaseProvider: base}
}

// LoadHolidays loads all Belgian holidays for a given year
func (be *BEProvider) LoadHolidays(year int) map[time.Time]*Holiday {
	holidays := make(map[time.Time]*Holiday)

	// Fixed date holidays

	// New Year's Day - January 1
	newYear := time.Date(year, 1, 1, 0, 0, 0, 0, time.UTC)
	holidays[newYear] = be.CreateHoliday(
		"Nieuwjaar",
		newYear,
		"public",
		map[string]string{
			"nl": "Nieuwjaar",
			"fr": "Nouvel An",
			"de": "Neujahr",
			"en": "New Year's Day",
		},
	)

	// Labour Day - May 1
	labourDay := time.Date(year, 5, 1, 0, 0, 0, 0, time.UTC)
	holidays[labourDay] = be.CreateHoliday(
		"Dag van de Arbeid",
		labourDay,
		"public",
		map[string]string{
			"nl": "Dag van de Arbeid",
			"fr": "Fête du Travail",
			"de": "Tag der Arbeit",
			"en": "Labour Day",
		},
	)

	// Belgian National Day - July 21
	nationalDay := time.Date(year, 7, 21, 0, 0, 0, 0, time.UTC)
	holidays[nationalDay] = be.CreateHoliday(
		"Nationale Feestdag",
		nationalDay,
		"public",
		map[string]string{
			"nl": "Nationale Feestdag",
			"fr": "Fête nationale",
			"de": "Nationalfeiertag",
			"en": "Belgian National Day",
		},
	)

	// Assumption of Mary - August 15
	assumption := time.Date(year, 8, 15, 0, 0, 0, 0, time.UTC)
	holidays[assumption] = be.CreateHoliday(
		"Onze-Lieve-Vrouw-Hemelvaart",
		assumption,
		"religious",
		map[string]string{
			"nl": "Onze-Lieve-Vrouw-Hemelvaart",
			"fr": "Assomption",
			"de": "Mariä Himmelfahrt",
			"en": "Assumption of Mary",
		},
	)

	// All Saints' Day - November 1
	allSaints := time.Date(year, 11, 1, 0, 0, 0, 0, time.UTC)
	holidays[allSaints] = be.CreateHoliday(
		"Allerheiligen",
		allSaints,
		"religious",
		map[string]string{
			"nl": "Allerheiligen",
			"fr": "Toussaint",
			"de": "Allerheiligen",
			"en": "All Saints' Day",
		},
	)

	// Armistice Day - November 11
	armistice := time.Date(year, 11, 11, 0, 0, 0, 0, time.UTC)
	holidays[armistice] = be.CreateHoliday(
		"Wapenstilstand",
		armistice,
		"public",
		map[string]string{
			"nl": "Wapenstilstand",
			"fr": "Armistice",
			"de": "Waffenstillstand",
			"en": "Armistice Day",
		},
	)

	// Christmas Day - December 25
	christmas := time.Date(year, 12, 25, 0, 0, 0, 0, time.UTC)
	holidays[christmas] = be.CreateHoliday(
		"Kerstmis",
		christmas,
		"public",
		map[string]string{
			"nl": "Kerstmis",
			"fr": "Noël",
			"de": "Weihnachten",
			"en": "Christmas Day",
		},
	)

	// Easter-based holidays
	easter := EasterSunday(year)

	// Easter Sunday
	holidays[easter] = be.CreateHoliday(
		"Pasen",
		easter,
		"religious",
		map[string]string{
			"nl": "Pasen",
			"fr": "Pâques",
			"de": "Ostern",
			"en": "Easter Sunday",
		},
	)

	// Easter Monday
	easterMonday := easter.AddDate(0, 0, 1)
	holidays[easterMonday] = be.CreateHoliday(
		"Paasmaandag",
		easterMonday,
		"public",
		map[string]string{
			"nl": "Paasmaandag",
			"fr": "Lundi de Pâques",
			"de": "Ostermontag",
			"en": "Easter Monday",
		},
	)

	// Ascension Day (39 days after Easter)
	ascension := easter.AddDate(0, 0, 39)
	holidays[ascension] = be.CreateHoliday(
		"Onze-Lieve-Heer-Hemelvaart",
		ascension,
		"public",
		map[string]string{
			"nl": "Onze-Lieve-Heer-Hemelvaart",
			"fr": "Ascension",
			"de": "Christi Himmelfahrt",
			"en": "Ascension Day",
		},
	)

	// Whit Sunday (49 days after Easter)
	whitSunday := easter.AddDate(0, 0, 49)
	holidays[whitSunday] = be.CreateHoliday(
		"Pinksteren",
		whitSunday,
		"religious",
		map[string]string{
			"nl": "Pinksteren",
			"fr": "Pentecôte",
			"de": "Pfingsten",
			"en": "Whit Sunday",
		},
	)

	// Whit Monday (50 days after Easter)
	whitMonday := easter.AddDate(0, 0, 50)
	holidays[whitMonday] = be.CreateHoliday(
		"Pinkstermaandag",
		whitMonday,
		"public",
		map[string]string{
			"nl": "Pinkstermaandag",
			"fr": "Lundi de Pentecôte",
			"de": "Pfingstmontag",
			"en": "Whit Monday",
		},
	)

	return holidays
}

// GetRegionalHolidays returns region-specific holidays
func (be *BEProvider) GetRegionalHolidays(year int, regions []string) map[time.Time]*Holiday {
	holidays := make(map[time.Time]*Holiday)

	for _, region := range regions {
		switch region {
		case "VLG": // Flemish Region
			// Day of the Flemish Community - July 11
			flemishDay := time.Date(year, 7, 11, 0, 0, 0, 0, time.UTC)
			holidays[flemishDay] = be.CreateHoliday(
				"Feest van de Vlaamse Gemeenschap",
				flemishDay,
				"regional",
				map[string]string{
					"nl": "Feest van de Vlaamse Gemeenschap",
					"en": "Day of the Flemish Community",
				},
			)

		case "WAL": // Walloon Region
			// Day of the French Community - September 27
			frenchDay := time.Date(year, 9, 27, 0, 0, 0, 0, time.UTC)
			holidays[frenchDay] = be.CreateHoliday(
				"Fête de la Communauté française",
				frenchDay,
				"regional",
				map[string]string{
					"fr": "Fête de la Communauté française",
					"en": "Day of the French Community",
				},
			)

		case "BRU": // Brussels-Capital Region
			// Iris Festival - May 8
			iris := time.Date(year, 5, 8, 0, 0, 0, 0, time.UTC)
			holidays[iris] = be.CreateHoliday(
				"Feest van het Brussels Hoofdstedelijk Gewest",
				iris,
				"regional",
				map[string]string{
					"nl": "Feest van het Brussels Hoofdstedelijk Gewest",
					"fr": "Fête de la Région de Bruxelles-Capitale",
					"en": "Brussels-Capital Region Day",
				},
			)
		}
	}

	return holidays
}

// GetSpecialObservances returns non-public observances
func (be *BEProvider) GetSpecialObservances(year int) map[time.Time]*Holiday {
	holidays := make(map[time.Time]*Holiday)

	// Carnival (Mardi Gras) - 47 days before Easter
	easter := EasterSunday(year)
	carnival := easter.AddDate(0, 0, -47)
	holidays[carnival] = be.CreateHoliday(
		"Carnaval",
		carnival,
		"regional",
		map[string]string{
			"nl": "Carnaval",
			"fr": "Carnaval",
			"en": "Carnival",
		},
	)

	// Saint Nicholas Day - December 6
	stNicholas := time.Date(year, 12, 6, 0, 0, 0, 0, time.UTC)
	holidays[stNicholas] = be.CreateHoliday(
		"Sinterklaas",
		stNicholas,
		"religious",
		map[string]string{
			"nl": "Sinterklaas",
			"fr": "Saint-Nicolas",
			"de": "Nikolaus",
			"en": "Saint Nicholas Day",
		},
	)

	// King's Day - November 15 (birthday of current King Philippe)
	if year >= 2013 { // Philippe became king in 2013
		kingsDay := time.Date(year, 11, 15, 0, 0, 0, 0, time.UTC)
		holidays[kingsDay] = be.CreateHoliday(
			"Koningsdag",
			kingsDay,
			"regional",
			map[string]string{
				"nl": "Koningsdag",
				"fr": "Fête du Roi",
				"en": "King's Day",
			},
		)
	}

	return holidays
}

