package countries

import (
	"time"
)

// CHProvider implements holiday calculations for Switzerland
type CHProvider struct {
	*BaseProvider
}

// NewCHProvider creates a new Swiss holiday provider
func NewCHProvider() *CHProvider {
	base := NewBaseProvider("CH")
	base.subdivisions = []string{
		// 26 cantons
		"AG", "AI", "AR", "BE", "BL", "BS", "FR", "GE", "GL", "GR",
		"JU", "LU", "NE", "NW", "OW", "SG", "SH", "SO", "SZ", "TG",
		"TI", "UR", "VD", "VS", "ZG", "ZH",
	}
	base.categories = []string{"federal", "cantonal", "religious", "cultural"}

	return &CHProvider{BaseProvider: base}
}

// LoadHolidays loads all Swiss holidays for a given year
func (ch *CHProvider) LoadHolidays(year int) map[time.Time]*Holiday {
	holidays := make(map[time.Time]*Holiday)

	// Federal holidays (observed nationwide)

	// New Year's Day - January 1
	newYear := time.Date(year, 1, 1, 0, 0, 0, 0, time.UTC)
	holidays[newYear] = ch.CreateHoliday(
		"Neujahr",
		newYear,
		"federal",
		map[string]string{
			"de": "Neujahr",
			"fr": "Nouvel An",
			"it": "Capodanno",
			"rm": "Niev on",
			"en": "New Year's Day",
		},
	)

	// Berchtoldstag - January 2 (some cantons)
	berchtoldstag := time.Date(year, 1, 2, 0, 0, 0, 0, time.UTC)
	holidays[berchtoldstag] = ch.CreateHoliday(
		"Berchtoldstag",
		berchtoldstag,
		"cantonal",
		map[string]string{
			"de": "Berchtoldstag",
			"fr": "Berchtoldstag",
			"it": "Berchtoldstag",
			"en": "Berchtoldstag",
		},
	)

	// Labour Day - May 1 (some cantons)
	labourDay := time.Date(year, 5, 1, 0, 0, 0, 0, time.UTC)
	holidays[labourDay] = ch.CreateHoliday(
		"Tag der Arbeit",
		labourDay,
		"cantonal",
		map[string]string{
			"de": "Tag der Arbeit",
			"fr": "Fête du Travail",
			"it": "Festa del Lavoro",
			"en": "Labour Day",
		},
	)

	// Swiss National Day - August 1
	nationalDay := time.Date(year, 8, 1, 0, 0, 0, 0, time.UTC)
	holidays[nationalDay] = ch.CreateHoliday(
		"Schweizer Nationalfeiertag",
		nationalDay,
		"federal",
		map[string]string{
			"de": "Schweizer Nationalfeiertag",
			"fr": "Fête nationale suisse",
			"it": "Festa nazionale svizzera",
			"rm": "Festa naziunala svizra",
			"en": "Swiss National Day",
		},
	)

	// Christmas Day - December 25
	christmas := time.Date(year, 12, 25, 0, 0, 0, 0, time.UTC)
	holidays[christmas] = ch.CreateHoliday(
		"Weihnachten",
		christmas,
		"federal",
		map[string]string{
			"de": "Weihnachten",
			"fr": "Noël",
			"it": "Natale",
			"rm": "Nadal",
			"en": "Christmas Day",
		},
	)

	// St. Stephen's Day - December 26 (most cantons)
	stephens := time.Date(year, 12, 26, 0, 0, 0, 0, time.UTC)
	holidays[stephens] = ch.CreateHoliday(
		"Stephanstag",
		stephens,
		"cantonal",
		map[string]string{
			"de": "Stephanstag",
			"fr": "Saint-Étienne",
			"it": "Santo Stefano",
			"en": "St. Stephen's Day",
		},
	)

	// Easter-based holidays
	easterDate := ch.CalculateEaster(year)

	// Good Friday - federal holiday
	goodFriday := easterDate.AddDate(0, 0, -2)
	holidays[goodFriday] = ch.CreateHoliday(
		"Karfreitag",
		goodFriday,
		"federal",
		map[string]string{
			"de": "Karfreitag",
			"fr": "Vendredi saint",
			"it": "Venerdì santo",
			"rm": "Venderdi sogn",
			"en": "Good Friday",
		},
	)

	// Easter Monday - cantonal holiday
	easterMonday := easterDate.AddDate(0, 0, 1)
	holidays[easterMonday] = ch.CreateHoliday(
		"Ostermontag",
		easterMonday,
		"cantonal",
		map[string]string{
			"de": "Ostermontag",
			"fr": "Lundi de Pâques",
			"it": "Lunedì di Pasqua",
			"rm": "Glindesdi da Pasqua",
			"en": "Easter Monday",
		},
	)

	// Ascension Day - federal holiday (39 days after Easter)
	ascension := easterDate.AddDate(0, 0, 39)
	holidays[ascension] = ch.CreateHoliday(
		"Auffahrt",
		ascension,
		"federal",
		map[string]string{
			"de": "Auffahrt",
			"fr": "Ascension",
			"it": "Ascensione",
			"rm": "Ascensiun",
			"en": "Ascension Day",
		},
	)

	// Whit Monday - federal holiday (50 days after Easter)
	whitMonday := easterDate.AddDate(0, 0, 50)
	holidays[whitMonday] = ch.CreateHoliday(
		"Pfingstmontag",
		whitMonday,
		"federal",
		map[string]string{
			"de": "Pfingstmontag",
			"fr": "Lundi de Pentecôte",
			"it": "Lunedì di Pentecoste",
			"rm": "Glindesdi da Tschuncheisma",
			"en": "Whit Monday",
		},
	)

	// Corpus Christi - cantonal holiday (60 days after Easter)
	corpusChristi := easterDate.AddDate(0, 0, 60)
	holidays[corpusChristi] = ch.CreateHoliday(
		"Fronleichnam",
		corpusChristi,
		"cantonal",
		map[string]string{
			"de": "Fronleichnam",
			"fr": "Fête-Dieu",
			"it": "Corpus Domini",
			"en": "Corpus Christi",
		},
	)

	// Assumption of Mary - August 15 (cantonal)
	assumption := time.Date(year, 8, 15, 0, 0, 0, 0, time.UTC)
	holidays[assumption] = ch.CreateHoliday(
		"Mariä Himmelfahrt",
		assumption,
		"cantonal",
		map[string]string{
			"de": "Mariä Himmelfahrt",
			"fr": "Assomption",
			"it": "Assunzione di Maria",
			"en": "Assumption of Mary",
		},
	)

	// All Saints' Day - November 1 (cantonal)
	allSaints := time.Date(year, 11, 1, 0, 0, 0, 0, time.UTC)
	holidays[allSaints] = ch.CreateHoliday(
		"Allerheiligen",
		allSaints,
		"cantonal",
		map[string]string{
			"de": "Allerheiligen",
			"fr": "Toussaint",
			"it": "Ognissanti",
			"en": "All Saints' Day",
		},
	)

	// Immaculate Conception - December 8 (cantonal)
	immaculate := time.Date(year, 12, 8, 0, 0, 0, 0, time.UTC)
	holidays[immaculate] = ch.CreateHoliday(
		"Mariä Empfängnis",
		immaculate,
		"cantonal",
		map[string]string{
			"de": "Mariä Empfängnis",
			"fr": "Immaculée Conception",
			"it": "Immacolata Concezione",
			"en": "Immaculate Conception",
		},
	)

	return holidays
}

// CreateHoliday creates a new holiday with Swiss localization
func (ch *CHProvider) CreateHoliday(name string, date time.Time, category string, languages map[string]string) *Holiday {
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
func (ch *CHProvider) CalculateEaster(year int) time.Time {
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
