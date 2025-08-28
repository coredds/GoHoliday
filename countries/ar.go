package countries

import (
	"time"
)

// ARProvider implements holiday calculations for Argentina
type ARProvider struct {
	*BaseProvider
}

// NewARProvider creates a new Argentine holiday provider
func NewARProvider() *ARProvider {
	base := NewBaseProvider("AR")
	base.subdivisions = []string{
		// 23 provinces + 1 autonomous city
		"C", "B", "K", "H", "U", "X", "W", "E", "P", "Y", "L",
		"F", "M", "N", "Q", "R", "A", "D", "Z", "S", "G", "V", "T", "J",
	}
	base.categories = []string{"national", "religious", "provincial", "commemorative"}

	return &ARProvider{BaseProvider: base}
}

// LoadHolidays loads all Argentine holidays for a given year
func (ar *ARProvider) LoadHolidays(year int) map[time.Time]*Holiday {
	holidays := make(map[time.Time]*Holiday)

	// Fixed national holidays

	// New Year's Day - January 1
	newYear := time.Date(year, 1, 1, 0, 0, 0, 0, time.UTC)
	holidays[newYear] = ar.CreateHoliday(
		"Año Nuevo",
		newYear,
		"national",
		map[string]string{
			"es": "Año Nuevo",
			"en": "New Year's Day",
		},
	)

	// Truth and Justice Day - March 24
	truthDay := time.Date(year, 3, 24, 0, 0, 0, 0, time.UTC)
	holidays[truthDay] = ar.CreateHoliday(
		"Día Nacional de la Memoria por la Verdad y la Justicia",
		truthDay,
		"commemorative",
		map[string]string{
			"es": "Día Nacional de la Memoria por la Verdad y la Justicia",
			"en": "Day of Remembrance for Truth and Justice",
		},
	)

	// Veterans Day and Fallen in Malvinas War - April 2
	malvinasDay := time.Date(year, 4, 2, 0, 0, 0, 0, time.UTC)
	holidays[malvinasDay] = ar.CreateHoliday(
		"Día del Veterano y de los Caídos en la Guerra de Malvinas",
		malvinasDay,
		"commemorative",
		map[string]string{
			"es": "Día del Veterano y de los Caídos en la Guerra de Malvinas",
			"en": "Veterans Day and Fallen in Malvinas War",
		},
	)

	// Labour Day - May 1
	labourDay := time.Date(year, 5, 1, 0, 0, 0, 0, time.UTC)
	holidays[labourDay] = ar.CreateHoliday(
		"Día del Trabajador",
		labourDay,
		"national",
		map[string]string{
			"es": "Día del Trabajador",
			"en": "Labour Day",
		},
	)

	// May Revolution Day - May 25
	mayRevolution := time.Date(year, 5, 25, 0, 0, 0, 0, time.UTC)
	holidays[mayRevolution] = ar.CreateHoliday(
		"Día de la Revolución de Mayo",
		mayRevolution,
		"national",
		map[string]string{
			"es": "Día de la Revolución de Mayo",
			"en": "May Revolution Day",
		},
	)

	// Flag Day - June 20 (moved to Monday if weekend)
	flagDay := ar.calculateMovableHoliday(year, 6, 20)
	holidays[flagDay] = ar.CreateHoliday(
		"Día de la Bandera",
		flagDay,
		"national",
		map[string]string{
			"es": "Día de la Bandera",
			"en": "Flag Day",
		},
	)

	// Independence Day - July 9
	independenceDay := time.Date(year, 7, 9, 0, 0, 0, 0, time.UTC)
	holidays[independenceDay] = ar.CreateHoliday(
		"Día de la Independencia",
		independenceDay,
		"national",
		map[string]string{
			"es": "Día de la Independencia",
			"en": "Independence Day",
		},
	)

	// San Martín Day - August 17 (moved to Monday if weekend)
	sanMartinDay := ar.calculateMovableHoliday(year, 8, 17)
	holidays[sanMartinDay] = ar.CreateHoliday(
		"Paso a la Inmortalidad del General José de San Martín",
		sanMartinDay,
		"national",
		map[string]string{
			"es": "Paso a la Inmortalidad del General José de San Martín",
			"en": "San Martín Day",
		},
	)

	// Columbus Day - October 12 (moved to Monday if weekend)
	columbusDay := ar.calculateMovableHoliday(year, 10, 12)
	holidays[columbusDay] = ar.CreateHoliday(
		"Día del Respeto a la Diversidad Cultural",
		columbusDay,
		"national",
		map[string]string{
			"es": "Día del Respeto a la Diversidad Cultural",
			"en": "Day of Respect for Cultural Diversity",
		},
	)

	// National Sovereignty Day - November 20 (moved to Monday if weekend)
	sovereigntyDay := ar.calculateMovableHoliday(year, 11, 20)
	holidays[sovereigntyDay] = ar.CreateHoliday(
		"Día de la Soberanía Nacional",
		sovereigntyDay,
		"national",
		map[string]string{
			"es": "Día de la Soberanía Nacional",
			"en": "National Sovereignty Day",
		},
	)

	// Immaculate Conception - December 8
	immaculate := time.Date(year, 12, 8, 0, 0, 0, 0, time.UTC)
	holidays[immaculate] = ar.CreateHoliday(
		"Inmaculada Concepción de María",
		immaculate,
		"religious",
		map[string]string{
			"es": "Inmaculada Concepción de María",
			"en": "Immaculate Conception",
		},
	)

	// Christmas Day - December 25
	christmas := time.Date(year, 12, 25, 0, 0, 0, 0, time.UTC)
	holidays[christmas] = ar.CreateHoliday(
		"Navidad",
		christmas,
		"religious",
		map[string]string{
			"es": "Navidad",
			"en": "Christmas Day",
		},
	)

	// Easter-based holidays
	easterDate := ar.CalculateEaster(year)

	// Carnival Monday (48 days before Easter)
	carnivalMonday := easterDate.AddDate(0, 0, -48)
	holidays[carnivalMonday] = ar.CreateHoliday(
		"Lunes de Carnaval",
		carnivalMonday,
		"national",
		map[string]string{
			"es": "Lunes de Carnaval",
			"en": "Carnival Monday",
		},
	)

	// Carnival Tuesday (47 days before Easter)
	carnivalTuesday := easterDate.AddDate(0, 0, -47)
	holidays[carnivalTuesday] = ar.CreateHoliday(
		"Martes de Carnaval",
		carnivalTuesday,
		"national",
		map[string]string{
			"es": "Martes de Carnaval",
			"en": "Carnival Tuesday",
		},
	)

	// Maundy Thursday (3 days before Easter)
	maundyThursday := easterDate.AddDate(0, 0, -3)
	holidays[maundyThursday] = ar.CreateHoliday(
		"Jueves Santo",
		maundyThursday,
		"religious",
		map[string]string{
			"es": "Jueves Santo",
			"en": "Maundy Thursday",
		},
	)

	// Good Friday (2 days before Easter)
	goodFriday := easterDate.AddDate(0, 0, -2)
	holidays[goodFriday] = ar.CreateHoliday(
		"Viernes Santo",
		goodFriday,
		"religious",
		map[string]string{
			"es": "Viernes Santo",
			"en": "Good Friday",
		},
	)

	return holidays
}

// calculateMovableHoliday calculates holidays that are moved to Monday if they fall on weekends
func (ar *ARProvider) calculateMovableHoliday(year int, month int, day int) time.Time {
	date := time.Date(year, time.Month(month), day, 0, 0, 0, 0, time.UTC)

	// If the holiday falls on Saturday or Sunday, move to the following Monday
	switch date.Weekday() {
	case time.Saturday:
		return date.AddDate(0, 0, 2) // Move to Monday
	case time.Sunday:
		return date.AddDate(0, 0, 1) // Move to Monday
	default:
		return date // Keep original date
	}
}

// CreateHoliday creates a new holiday with Argentine localization
func (ar *ARProvider) CreateHoliday(name string, date time.Time, category string, languages map[string]string) *Holiday {
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
func (ar *ARProvider) CalculateEaster(year int) time.Time {
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
