package countries

import (
	"time"
)

// NLProvider implements holiday calculations for Netherlands
type NLProvider struct {
	*BaseProvider
}

// NewNLProvider creates a new Dutch holiday provider
func NewNLProvider() *NLProvider {
	base := NewBaseProvider("NL")
	base.subdivisions = []string{
		// 12 provinces
		"DR", "FL", "FR", "GE", "GR", "LI", "NB", "NH", "OV", "UT", "ZE", "ZH",
	}
	base.categories = []string{"public", "national", "religious", "royal"}

	return &NLProvider{BaseProvider: base}
}

// LoadHolidays loads all Dutch holidays for a given year
func (nl *NLProvider) LoadHolidays(year int) map[time.Time]*Holiday {
	holidays := make(map[time.Time]*Holiday)

	// Fixed national holidays

	// New Year's Day - January 1
	newYear := time.Date(year, 1, 1, 0, 0, 0, 0, time.UTC)
	holidays[newYear] = nl.CreateHoliday(
		"Nieuwjaarsdag",
		newYear,
		"public",
		map[string]string{
			"nl": "Nieuwjaarsdag",
			"en": "New Year's Day",
		},
	)

	// King's Day - April 27 (or April 26 if April 27 is Sunday)
	kingsDay := time.Date(year, 4, 27, 0, 0, 0, 0, time.UTC)
	if kingsDay.Weekday() == time.Sunday {
		kingsDay = time.Date(year, 4, 26, 0, 0, 0, 0, time.UTC)
	}
	holidays[kingsDay] = nl.CreateHoliday(
		"Koningsdag",
		kingsDay,
		"royal",
		map[string]string{
			"nl": "Koningsdag",
			"en": "King's Day",
		},
	)

	// Liberation Day - May 5
	liberation := time.Date(year, 5, 5, 0, 0, 0, 0, time.UTC)
	holidays[liberation] = nl.CreateHoliday(
		"Bevrijdingsdag",
		liberation,
		"public",
		map[string]string{
			"nl": "Bevrijdingsdag",
			"en": "Liberation Day",
		},
	)

	// Christmas Day - December 25
	christmas := time.Date(year, 12, 25, 0, 0, 0, 0, time.UTC)
	holidays[christmas] = nl.CreateHoliday(
		"Eerste Kerstdag",
		christmas,
		"religious",
		map[string]string{
			"nl": "Eerste Kerstdag",
			"en": "Christmas Day",
		},
	)

	// Boxing Day - December 26
	boxingDay := time.Date(year, 12, 26, 0, 0, 0, 0, time.UTC)
	holidays[boxingDay] = nl.CreateHoliday(
		"Tweede Kerstdag",
		boxingDay,
		"religious",
		map[string]string{
			"nl": "Tweede Kerstdag",
			"en": "Boxing Day",
		},
	)

	// Easter-based holidays
	easterDate := nl.CalculateEaster(year)

	// Easter Sunday (Pasen)
	holidays[easterDate] = nl.CreateHoliday(
		"Eerste Paasdag",
		easterDate,
		"religious",
		map[string]string{
			"nl": "Eerste Paasdag",
			"en": "Easter Sunday",
		},
	)

	// Easter Monday (Tweede Paasdag)
	easterMonday := easterDate.AddDate(0, 0, 1)
	holidays[easterMonday] = nl.CreateHoliday(
		"Tweede Paasdag",
		easterMonday,
		"religious",
		map[string]string{
			"nl": "Tweede Paasdag",
			"en": "Easter Monday",
		},
	)

	// Good Friday (Goede Vrijdag)
	goodFriday := easterDate.AddDate(0, 0, -2)
	holidays[goodFriday] = nl.CreateHoliday(
		"Goede Vrijdag",
		goodFriday,
		"religious",
		map[string]string{
			"nl": "Goede Vrijdag",
			"en": "Good Friday",
		},
	)

	// Ascension Day (Hemelvaartsdag) - 39 days after Easter
	ascension := easterDate.AddDate(0, 0, 39)
	holidays[ascension] = nl.CreateHoliday(
		"Hemelvaartsdag",
		ascension,
		"religious",
		map[string]string{
			"nl": "Hemelvaartsdag",
			"en": "Ascension Day",
		},
	)

	// Whit Sunday (Pinksteren) - 49 days after Easter
	whitSunday := easterDate.AddDate(0, 0, 49)
	holidays[whitSunday] = nl.CreateHoliday(
		"Eerste Pinksterdag",
		whitSunday,
		"religious",
		map[string]string{
			"nl": "Eerste Pinksterdag",
			"en": "Whit Sunday",
		},
	)

	// Whit Monday (Tweede Pinksterdag) - 50 days after Easter
	whitMonday := easterDate.AddDate(0, 0, 50)
	holidays[whitMonday] = nl.CreateHoliday(
		"Tweede Pinksterdag",
		whitMonday,
		"religious",
		map[string]string{
			"nl": "Tweede Pinksterdag",
			"en": "Whit Monday",
		},
	)

	return holidays
}

// CreateHoliday creates a new holiday with Dutch localization
func (nl *NLProvider) CreateHoliday(name string, date time.Time, category string, languages map[string]string) *Holiday {
	return &Holiday{
		Name:         name,
		Date:         date,
		Category:     category,
		Languages:    languages,
		IsObserved:   true,
		Subdivisions: []string{},
	}
}

// CalculateEaster calculates Easter date for a given year using the Western (Gregorian) calculation
func (nl *NLProvider) CalculateEaster(year int) time.Time {
	// Using the anonymous Gregorian algorithm
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
