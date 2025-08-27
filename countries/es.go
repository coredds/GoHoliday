package countries

import (
	"time"
)

// ESProvider implements holiday calculations for Spain
type ESProvider struct {
	*BaseProvider
}

// NewESProvider creates a new Spanish holiday provider
func NewESProvider() *ESProvider {
	base := NewBaseProvider("ES")
	base.subdivisions = []string{
		// 17 autonomous communities + 2 autonomous cities
		"AN", "AR", "AS", "IB", "IC", "CB", "CL", "CM", "CT", "VC",
		"EX", "GA", "MD", "MU", "NA", "PV", "RI", "CE", "ML",
	}
	base.categories = []string{"public", "national", "religious", "regional"}

	return &ESProvider{BaseProvider: base}
}

// LoadHolidays loads all Spanish holidays for a given year
func (es *ESProvider) LoadHolidays(year int) map[time.Time]*Holiday {
	holidays := make(map[time.Time]*Holiday)

	// Fixed national holidays

	// New Year's Day - January 1
	newYear := time.Date(year, 1, 1, 0, 0, 0, 0, time.UTC)
	holidays[newYear] = es.CreateHoliday(
		"Año Nuevo",
		newYear,
		"public",
		map[string]string{
			"es": "Año Nuevo",
			"en": "New Year's Day",
		},
	)

	// Epiphany - January 6
	epiphany := time.Date(year, 1, 6, 0, 0, 0, 0, time.UTC)
	holidays[epiphany] = es.CreateHoliday(
		"Día de los Reyes Magos",
		epiphany,
		"religious",
		map[string]string{
			"es": "Día de los Reyes Magos",
			"en": "Epiphany",
		},
	)

	// Labour Day - May 1
	labourDay := time.Date(year, 5, 1, 0, 0, 0, 0, time.UTC)
	holidays[labourDay] = es.CreateHoliday(
		"Día del Trabajador",
		labourDay,
		"public",
		map[string]string{
			"es": "Día del Trabajador",
			"en": "Labour Day",
		},
	)

	// Assumption of Mary - August 15
	assumption := time.Date(year, 8, 15, 0, 0, 0, 0, time.UTC)
	holidays[assumption] = es.CreateHoliday(
		"Asunción de la Virgen",
		assumption,
		"religious",
		map[string]string{
			"es": "Asunción de la Virgen",
			"en": "Assumption of Mary",
		},
	)

	// National Day of Spain - October 12
	nationalDay := time.Date(year, 10, 12, 0, 0, 0, 0, time.UTC)
	holidays[nationalDay] = es.CreateHoliday(
		"Fiesta Nacional de España",
		nationalDay,
		"public",
		map[string]string{
			"es": "Fiesta Nacional de España",
			"en": "National Day of Spain",
		},
	)

	// All Saints' Day - November 1
	allSaints := time.Date(year, 11, 1, 0, 0, 0, 0, time.UTC)
	holidays[allSaints] = es.CreateHoliday(
		"Día de Todos los Santos",
		allSaints,
		"religious",
		map[string]string{
			"es": "Día de Todos los Santos",
			"en": "All Saints' Day",
		},
	)

	// Constitution Day - December 6
	constitution := time.Date(year, 12, 6, 0, 0, 0, 0, time.UTC)
	holidays[constitution] = es.CreateHoliday(
		"Día de la Constitución",
		constitution,
		"public",
		map[string]string{
			"es": "Día de la Constitución",
			"en": "Constitution Day",
		},
	)

	// Immaculate Conception - December 8
	immaculate := time.Date(year, 12, 8, 0, 0, 0, 0, time.UTC)
	holidays[immaculate] = es.CreateHoliday(
		"Inmaculada Concepción",
		immaculate,
		"religious",
		map[string]string{
			"es": "Inmaculada Concepción",
			"en": "Immaculate Conception",
		},
	)

	// Christmas Day - December 25
	christmas := time.Date(year, 12, 25, 0, 0, 0, 0, time.UTC)
	holidays[christmas] = es.CreateHoliday(
		"Navidad",
		christmas,
		"religious",
		map[string]string{
			"es": "Navidad",
			"en": "Christmas Day",
		},
	)

	// Easter-based holidays
	easterDate := es.CalculateEaster(year)

	// Maundy Thursday (Jueves Santo)
	maundyThursday := easterDate.AddDate(0, 0, -3)
	holidays[maundyThursday] = es.CreateHoliday(
		"Jueves Santo",
		maundyThursday,
		"religious",
		map[string]string{
			"es": "Jueves Santo",
			"en": "Maundy Thursday",
		},
	)

	// Good Friday (Viernes Santo)
	goodFriday := easterDate.AddDate(0, 0, -2)
	holidays[goodFriday] = es.CreateHoliday(
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

// CreateHoliday creates a new holiday with Spanish localization
func (es *ESProvider) CreateHoliday(name string, date time.Time, category string, languages map[string]string) *Holiday {
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
func (es *ESProvider) CalculateEaster(year int) time.Time {
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
