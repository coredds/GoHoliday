package countries

import (
	"time"
)

// CAProvider implements holiday calculations for Canada
type CAProvider struct {
	*BaseProvider
}

// NewCAProvider creates a new Canada holiday provider
func NewCAProvider() *CAProvider {
	base := NewBaseProvider("CA")
	base.subdivisions = []string{
		"AB", // Alberta
		"BC", // British Columbia
		"MB", // Manitoba
		"NB", // New Brunswick
		"NL", // Newfoundland and Labrador
		"NS", // Nova Scotia
		"NT", // Northwest Territories
		"NU", // Nunavut
		"ON", // Ontario
		"PE", // Prince Edward Island
		"QC", // Quebec
		"SK", // Saskatchewan
		"YT", // Yukon
	}
	base.categories = []string{"public", "bank", "government"}

	return &CAProvider{BaseProvider: base}
}

// LoadHolidays loads all Canadian holidays for a given year
func (ca *CAProvider) LoadHolidays(year int) map[time.Time]*Holiday {
	holidays := make(map[time.Time]*Holiday)

	// Fixed date holidays
	holidays[time.Date(year, 1, 1, 0, 0, 0, 0, time.UTC)] = ca.CreateHoliday(
		"New Year's Day",
		time.Date(year, 1, 1, 0, 0, 0, 0, time.UTC),
		"public",
		map[string]string{
			"en": "New Year's Day",
			"fr": "Jour de l'An",
		},
	)

	holidays[time.Date(year, 7, 1, 0, 0, 0, 0, time.UTC)] = ca.CreateHoliday(
		"Canada Day",
		time.Date(year, 7, 1, 0, 0, 0, 0, time.UTC),
		"public",
		map[string]string{
			"en": "Canada Day",
			"fr": "Fête du Canada",
		},
	)

	holidays[time.Date(year, 11, 11, 0, 0, 0, 0, time.UTC)] = ca.CreateHoliday(
		"Remembrance Day",
		time.Date(year, 11, 11, 0, 0, 0, 0, time.UTC),
		"public",
		map[string]string{
			"en": "Remembrance Day",
			"fr": "Jour du Souvenir",
		},
	)

	holidays[time.Date(year, 12, 25, 0, 0, 0, 0, time.UTC)] = ca.CreateHoliday(
		"Christmas Day",
		time.Date(year, 12, 25, 0, 0, 0, 0, time.UTC),
		"public",
		map[string]string{
			"en": "Christmas Day",
			"fr": "Noël",
		},
	)

	holidays[time.Date(year, 12, 26, 0, 0, 0, 0, time.UTC)] = ca.CreateHoliday(
		"Boxing Day",
		time.Date(year, 12, 26, 0, 0, 0, 0, time.UTC),
		"public",
		map[string]string{
			"en": "Boxing Day",
			"fr": "Lendemain de Noël",
		},
	)

	// Easter-based holidays
	easter := EasterSunday(year)

	// Good Friday
	goodFriday := easter.AddDate(0, 0, -2)
	holidays[goodFriday] = ca.CreateHoliday(
		"Good Friday",
		goodFriday,
		"public",
		map[string]string{
			"en": "Good Friday",
			"fr": "Vendredi saint",
		},
	)

	// Easter Monday (not all provinces)
	easterMonday := easter.AddDate(0, 0, 1)
	holidays[easterMonday] = ca.CreateHoliday(
		"Easter Monday",
		easterMonday,
		"public",
		map[string]string{
			"en": "Easter Monday",
			"fr": "Lundi de Pâques",
		},
	)

	// Variable date holidays

	// Family Day - 3rd Monday in February (varies by province)
	familyDay := NthWeekdayOfMonth(year, 2, time.Monday, 3)
	holidays[familyDay] = ca.CreateHoliday(
		"Family Day",
		familyDay,
		"public",
		map[string]string{
			"en": "Family Day",
			"fr": "Fête de la famille",
		},
	)

	// Victoria Day - Monday before May 25
	victoriaDay := ca.getVictoriaDay(year)
	holidays[victoriaDay] = ca.CreateHoliday(
		"Victoria Day",
		victoriaDay,
		"public",
		map[string]string{
			"en": "Victoria Day",
			"fr": "Fête de la Reine",
		},
	)

	// Labour Day - 1st Monday in September
	labourDay := NthWeekdayOfMonth(year, 9, time.Monday, 1)
	holidays[labourDay] = ca.CreateHoliday(
		"Labour Day",
		labourDay,
		"public",
		map[string]string{
			"en": "Labour Day",
			"fr": "Fête du Travail",
		},
	)

	// Thanksgiving Day - 2nd Monday in October
	thanksgiving := NthWeekdayOfMonth(year, 10, time.Monday, 2)
	holidays[thanksgiving] = ca.CreateHoliday(
		"Thanksgiving Day",
		thanksgiving,
		"public",
		map[string]string{
			"en": "Thanksgiving Day",
			"fr": "Action de grâce",
		},
	)

	// National Day for Truth and Reconciliation - September 30 (since 2021)
	if year >= 2021 {
		truthReconciliation := time.Date(year, 9, 30, 0, 0, 0, 0, time.UTC)
		holidays[truthReconciliation] = ca.CreateHoliday(
			"National Day for Truth and Reconciliation",
			truthReconciliation,
			"public",
			map[string]string{
				"en": "National Day for Truth and Reconciliation",
				"fr": "Journée nationale de la vérité et de la réconciliation",
			},
		)
	}

	return holidays
}

// getVictoriaDay calculates Victoria Day (Monday before May 25)
func (ca *CAProvider) getVictoriaDay(year int) time.Time {
	may25 := time.Date(year, 5, 25, 0, 0, 0, 0, time.UTC)

	// Find the Monday before May 25
	daysToSubtract := int(may25.Weekday()) - int(time.Monday)
	if daysToSubtract <= 0 {
		daysToSubtract += 7
	}

	return may25.AddDate(0, 0, -daysToSubtract)
}

// GetProvincialHolidays returns province-specific holidays
func (ca *CAProvider) GetProvincialHolidays(year int, provinces []string) map[time.Time]*Holiday {
	holidays := make(map[time.Time]*Holiday)

	for _, province := range provinces {
		switch province {
		case "AB": // Alberta
			// Alberta Family Day - 3rd Monday in February
			familyDay := NthWeekdayOfMonth(year, 2, time.Monday, 3)
			holidays[familyDay] = ca.CreateHoliday(
				"Family Day",
				familyDay,
				"public",
				map[string]string{
					"en": "Family Day",
				},
			)

		case "BC": // British Columbia
			// Family Day - 2nd Monday in February
			familyDay := NthWeekdayOfMonth(year, 2, time.Monday, 2)
			holidays[familyDay] = ca.CreateHoliday(
				"Family Day",
				familyDay,
				"public",
				map[string]string{
					"en": "Family Day",
				},
			)

		case "ON": // Ontario
			// Family Day - 3rd Monday in February
			familyDay := NthWeekdayOfMonth(year, 2, time.Monday, 3)
			holidays[familyDay] = ca.CreateHoliday(
				"Family Day",
				familyDay,
				"public",
				map[string]string{
					"en": "Family Day",
				},
			)

		case "QC": // Quebec
			// St. Jean Baptiste Day - June 24
			stJeanBaptiste := time.Date(year, 6, 24, 0, 0, 0, 0, time.UTC)
			holidays[stJeanBaptiste] = ca.CreateHoliday(
				"St. Jean Baptiste Day",
				stJeanBaptiste,
				"public",
				map[string]string{
					"en": "St. Jean Baptiste Day",
					"fr": "Fête nationale du Québec",
				},
			)

		case "NL": // Newfoundland and Labrador
			// St. Patrick's Day - March 17 (or nearest Monday)
			stPatricks := time.Date(year, 3, 17, 0, 0, 0, 0, time.UTC)
			observed := ca.getObservedDate(stPatricks)
			holidays[observed] = ca.CreateHoliday(
				"St. Patrick's Day",
				observed,
				"public",
				map[string]string{
					"en": "St. Patrick's Day",
				},
			)
		}
	}

	return holidays
}

// getObservedDate calculates the observed date for a holiday
func (ca *CAProvider) getObservedDate(date time.Time) time.Time {
	switch date.Weekday() {
	case time.Saturday:
		return date.AddDate(0, 0, 2) // Move to Monday
	case time.Sunday:
		return date.AddDate(0, 0, 1) // Move to Monday
	default:
		return date
	}
}
