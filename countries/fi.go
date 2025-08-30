package countries

import (
	"time"
)

// FIProvider implements holiday calculations for Finland
type FIProvider struct {
	*BaseProvider
}

// NewFIProvider creates a new Finnish holiday provider
func NewFIProvider() *FIProvider {
	base := NewBaseProvider("FI")
	base.categories = []string{"public", "religious"}

	return &FIProvider{BaseProvider: base}
}

// LoadHolidays loads all Finnish holidays for a given year
func (fi *FIProvider) LoadHolidays(year int) map[time.Time]*Holiday {
	holidays := make(map[time.Time]*Holiday)

	// Fixed date holidays

	// New Year's Day - January 1
	newYear := time.Date(year, 1, 1, 0, 0, 0, 0, time.UTC)
	holidays[newYear] = fi.CreateHoliday(
		"Uudenvuodenpäivä",
		newYear,
		"public",
		map[string]string{
			"fi": "Uudenvuodenpäivä",
			"sv": "Nyårsdagen",
			"en": "New Year's Day",
		},
	)

	// Epiphany - January 6
	epiphany := time.Date(year, 1, 6, 0, 0, 0, 0, time.UTC)
	holidays[epiphany] = fi.CreateHoliday(
		"Loppiainen",
		epiphany,
		"religious",
		map[string]string{
			"fi": "Loppiainen",
			"sv": "Trettondedag jul",
			"en": "Epiphany",
		},
	)

	// May Day - May 1
	mayDay := time.Date(year, 5, 1, 0, 0, 0, 0, time.UTC)
	holidays[mayDay] = fi.CreateHoliday(
		"Vappu",
		mayDay,
		"public",
		map[string]string{
			"fi": "Vappu",
			"sv": "Första maj",
			"en": "May Day",
		},
	)

	// Independence Day - December 6
	independenceDay := time.Date(year, 12, 6, 0, 0, 0, 0, time.UTC)
	holidays[independenceDay] = fi.CreateHoliday(
		"Itsenäisyyspäivä",
		independenceDay,
		"public",
		map[string]string{
			"fi": "Itsenäisyyspäivä",
			"sv": "Självständighetsdagen",
			"en": "Independence Day",
		},
	)

	// Christmas Eve - December 24
	christmasEve := time.Date(year, 12, 24, 0, 0, 0, 0, time.UTC)
	holidays[christmasEve] = fi.CreateHoliday(
		"Jouluaatto",
		christmasEve,
		"public",
		map[string]string{
			"fi": "Jouluaatto",
			"sv": "Julafton",
			"en": "Christmas Eve",
		},
	)

	// Christmas Day - December 25
	christmas := time.Date(year, 12, 25, 0, 0, 0, 0, time.UTC)
	holidays[christmas] = fi.CreateHoliday(
		"Joulupäivä",
		christmas,
		"public",
		map[string]string{
			"fi": "Joulupäivä",
			"sv": "Juldagen",
			"en": "Christmas Day",
		},
	)

	// Boxing Day - December 26
	boxingDay := time.Date(year, 12, 26, 0, 0, 0, 0, time.UTC)
	holidays[boxingDay] = fi.CreateHoliday(
		"Tapaninpäivä",
		boxingDay,
		"public",
		map[string]string{
			"fi": "Tapaninpäivä",
			"sv": "Annandag jul",
			"en": "Boxing Day",
		},
	)

	// Midsummer Eve (Friday between June 19-25)
	midsummerEve := NthWeekdayOfMonth(year, 6, time.Friday, 3)
	if midsummerEve.Day() < 19 {
		midsummerEve = NthWeekdayOfMonth(year, 6, time.Friday, 4)
	}
	holidays[midsummerEve] = fi.CreateHoliday(
		"Juhannusaatto",
		midsummerEve,
		"public",
		map[string]string{
			"fi": "Juhannusaatto",
			"sv": "Midsommarafton",
			"en": "Midsummer Eve",
		},
	)

	// Midsummer Day (Saturday between June 20-26)
	midsummerDay := midsummerEve.AddDate(0, 0, 1)
	holidays[midsummerDay] = fi.CreateHoliday(
		"Juhannuspäivä",
		midsummerDay,
		"public",
		map[string]string{
			"fi": "Juhannuspäivä",
			"sv": "Midsommardagen",
			"en": "Midsummer Day",
		},
	)

	// All Saints' Day (Saturday between Oct 31 and Nov 6)
	allSaints := NthWeekdayOfMonth(year, 10, time.Saturday, -1)
	if allSaints.Day() < 31 {
		allSaints = NthWeekdayOfMonth(year, 11, time.Saturday, 1)
	}
	holidays[allSaints] = fi.CreateHoliday(
		"Pyhäinpäivä",
		allSaints,
		"religious",
		map[string]string{
			"fi": "Pyhäinpäivä",
			"sv": "Alla helgons dag",
			"en": "All Saints' Day",
		},
	)

	// Easter-based holidays
	easter := EasterSunday(year)

	// Good Friday
	goodFriday := easter.AddDate(0, 0, -2)
	holidays[goodFriday] = fi.CreateHoliday(
		"Pitkäperjantai",
		goodFriday,
		"religious",
		map[string]string{
			"fi": "Pitkäperjantai",
			"sv": "Långfredag",
			"en": "Good Friday",
		},
	)

	// Easter Sunday
	holidays[easter] = fi.CreateHoliday(
		"Pääsiäispäivä",
		easter,
		"religious",
		map[string]string{
			"fi": "Pääsiäispäivä",
			"sv": "Påskdagen",
			"en": "Easter Sunday",
		},
	)

	// Easter Monday
	easterMonday := easter.AddDate(0, 0, 1)
	holidays[easterMonday] = fi.CreateHoliday(
		"Toinen pääsiäispäivä",
		easterMonday,
		"religious",
		map[string]string{
			"fi": "Toinen pääsiäispäivä",
			"sv": "Annandag påsk",
			"en": "Easter Monday",
		},
	)

	// Ascension Day (40 days after Easter)
	ascension := easter.AddDate(0, 0, 39)
	holidays[ascension] = fi.CreateHoliday(
		"Helatorstai",
		ascension,
		"religious",
		map[string]string{
			"fi": "Helatorstai",
			"sv": "Kristi himmelsfärdsdag",
			"en": "Ascension Day",
		},
	)

	// Pentecost (49 days after Easter)
	pentecost := easter.AddDate(0, 0, 49)
	holidays[pentecost] = fi.CreateHoliday(
		"Helluntaipäivä",
		pentecost,
		"religious",
		map[string]string{
			"fi": "Helluntaipäivä",
			"sv": "Pingstdagen",
			"en": "Pentecost",
		},
	)

	return holidays
}

