package countries

import (
	"time"
)

// SEProvider implements holiday calculations for Sweden
type SEProvider struct {
	*BaseProvider
}

// NewSEProvider creates a new Swedish holiday provider
func NewSEProvider() *SEProvider {
	base := NewBaseProvider("SE")
	base.subdivisions = []string{
		// 21 counties (län)
		"AB", "AC", "BD", "C", "D", "E", "F", "G", "H", "I", "K",
		"M", "N", "O", "S", "T", "U", "W", "X", "Y", "Z",
	}
	base.categories = []string{"public", "religious", "cultural", "traditional"}

	return &SEProvider{BaseProvider: base}
}

// LoadHolidays loads all Swedish holidays for a given year
func (se *SEProvider) LoadHolidays(year int) map[time.Time]*Holiday {
	holidays := make(map[time.Time]*Holiday)

	// Fixed national holidays

	// New Year's Day - January 1
	newYear := time.Date(year, 1, 1, 0, 0, 0, 0, time.UTC)
	holidays[newYear] = se.CreateHoliday(
		"Nyårsdagen",
		newYear,
		"public",
		map[string]string{
			"sv": "Nyårsdagen",
			"en": "New Year's Day",
		},
	)

	// Epiphany - January 6
	epiphany := time.Date(year, 1, 6, 0, 0, 0, 0, time.UTC)
	holidays[epiphany] = se.CreateHoliday(
		"Trettondedag jul",
		epiphany,
		"religious",
		map[string]string{
			"sv": "Trettondedag jul",
			"en": "Epiphany",
		},
	)

	// Labour Day - May 1
	labourDay := time.Date(year, 5, 1, 0, 0, 0, 0, time.UTC)
	holidays[labourDay] = se.CreateHoliday(
		"Första maj",
		labourDay,
		"public",
		map[string]string{
			"sv": "Första maj",
			"en": "Labour Day",
		},
	)

	// National Day of Sweden - June 6
	nationalDay := time.Date(year, 6, 6, 0, 0, 0, 0, time.UTC)
	holidays[nationalDay] = se.CreateHoliday(
		"Sveriges nationaldag",
		nationalDay,
		"public",
		map[string]string{
			"sv": "Sveriges nationaldag",
			"en": "National Day of Sweden",
		},
	)

	// Midsummer Eve - Friday between June 19-25
	midsummerEve := se.calculateMidsummerEve(year)
	holidays[midsummerEve] = se.CreateHoliday(
		"Midsommarafton",
		midsummerEve,
		"traditional",
		map[string]string{
			"sv": "Midsommarafton",
			"en": "Midsummer Eve",
		},
	)

	// Midsummer Day - Saturday between June 20-26
	midsummerDay := midsummerEve.AddDate(0, 0, 1)
	holidays[midsummerDay] = se.CreateHoliday(
		"Midsommardagen",
		midsummerDay,
		"traditional",
		map[string]string{
			"sv": "Midsommardagen",
			"en": "Midsummer Day",
		},
	)

	// All Saints' Day - Saturday between October 31 - November 6
	allSaints := se.calculateAllSaintsDay(year)
	holidays[allSaints] = se.CreateHoliday(
		"Alla helgons dag",
		allSaints,
		"religious",
		map[string]string{
			"sv": "Alla helgons dag",
			"en": "All Saints' Day",
		},
	)

	// Christmas Eve - December 24
	christmasEve := time.Date(year, 12, 24, 0, 0, 0, 0, time.UTC)
	holidays[christmasEve] = se.CreateHoliday(
		"Julafton",
		christmasEve,
		"traditional",
		map[string]string{
			"sv": "Julafton",
			"en": "Christmas Eve",
		},
	)

	// Christmas Day - December 25
	christmas := time.Date(year, 12, 25, 0, 0, 0, 0, time.UTC)
	holidays[christmas] = se.CreateHoliday(
		"Juldagen",
		christmas,
		"religious",
		map[string]string{
			"sv": "Juldagen",
			"en": "Christmas Day",
		},
	)

	// Boxing Day - December 26
	boxingDay := time.Date(year, 12, 26, 0, 0, 0, 0, time.UTC)
	holidays[boxingDay] = se.CreateHoliday(
		"Annandag jul",
		boxingDay,
		"religious",
		map[string]string{
			"sv": "Annandag jul",
			"en": "Boxing Day",
		},
	)

	// New Year's Eve - December 31
	newYearEve := time.Date(year, 12, 31, 0, 0, 0, 0, time.UTC)
	holidays[newYearEve] = se.CreateHoliday(
		"Nyårsafton",
		newYearEve,
		"traditional",
		map[string]string{
			"sv": "Nyårsafton",
			"en": "New Year's Eve",
		},
	)

	// Easter-based holidays
	easterDate := se.CalculateEaster(year)

	// Good Friday
	goodFriday := easterDate.AddDate(0, 0, -2)
	holidays[goodFriday] = se.CreateHoliday(
		"Långfredagen",
		goodFriday,
		"religious",
		map[string]string{
			"sv": "Långfredagen",
			"en": "Good Friday",
		},
	)

	// Easter Sunday
	holidays[easterDate] = se.CreateHoliday(
		"Påskdagen",
		easterDate,
		"religious",
		map[string]string{
			"sv": "Påskdagen",
			"en": "Easter Sunday",
		},
	)

	// Easter Monday
	easterMonday := easterDate.AddDate(0, 0, 1)
	holidays[easterMonday] = se.CreateHoliday(
		"Annandag påsk",
		easterMonday,
		"religious",
		map[string]string{
			"sv": "Annandag påsk",
			"en": "Easter Monday",
		},
	)

	// Ascension Day (39 days after Easter)
	ascension := easterDate.AddDate(0, 0, 39)
	holidays[ascension] = se.CreateHoliday(
		"Kristi himmelsfärdsdag",
		ascension,
		"religious",
		map[string]string{
			"sv": "Kristi himmelsfärdsdag",
			"en": "Ascension Day",
		},
	)

	// Whit Sunday (49 days after Easter)
	whitSunday := easterDate.AddDate(0, 0, 49)
	holidays[whitSunday] = se.CreateHoliday(
		"Pingstdagen",
		whitSunday,
		"religious",
		map[string]string{
			"sv": "Pingstdagen",
			"en": "Whit Sunday",
		},
	)

	return holidays
}

// calculateMidsummerEve calculates Midsummer Eve (Friday between June 19-25)
func (se *SEProvider) calculateMidsummerEve(year int) time.Time {
	// Start from June 19
	date := time.Date(year, 6, 19, 0, 0, 0, 0, time.UTC)

	// Find the first Friday on or after June 19
	for date.Weekday() != time.Friday {
		date = date.AddDate(0, 0, 1)
	}

	return date
}

// calculateAllSaintsDay calculates All Saints' Day (Saturday between October 31 - November 6)
func (se *SEProvider) calculateAllSaintsDay(year int) time.Time {
	// Start from October 31
	date := time.Date(year, 10, 31, 0, 0, 0, 0, time.UTC)

	// Find the first Saturday on or after October 31
	for date.Weekday() != time.Saturday {
		date = date.AddDate(0, 0, 1)
	}

	return date
}

// CreateHoliday creates a new holiday with Swedish localization
func (se *SEProvider) CreateHoliday(name string, date time.Time, category string, languages map[string]string) *Holiday {
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
func (se *SEProvider) CalculateEaster(year int) time.Time {
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
