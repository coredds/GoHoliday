package countries

import (
	"time"
)

// ITProvider implements holiday calculations for Italy
type ITProvider struct {
	*BaseProvider
}

// NewITProvider creates a new Italian holiday provider
func NewITProvider() *ITProvider {
	base := NewBaseProvider("IT")
	base.subdivisions = []string{
		// 20 regions
		"ABR", "BAS", "CAL", "CAM", "EMR", "FVG", "LAZ", "LIG", "LOM", "MAR",
		"MOL", "PIE", "PUG", "SAR", "SIC", "TOS", "TAA", "UMB", "VDA", "VEN",
	}
	base.categories = []string{"public", "national", "religious", "regional"}

	return &ITProvider{BaseProvider: base}
}

// LoadHolidays loads all Italian holidays for a given year
func (it *ITProvider) LoadHolidays(year int) map[time.Time]*Holiday {
	holidays := make(map[time.Time]*Holiday)

	// Fixed national holidays

	// New Year's Day - January 1
	newYear := time.Date(year, 1, 1, 0, 0, 0, 0, time.UTC)
	holidays[newYear] = it.CreateHoliday(
		"Capodanno",
		newYear,
		"public",
		map[string]string{
			"it": "Capodanno",
			"en": "New Year's Day",
		},
	)

	// Epiphany - January 6
	epiphany := time.Date(year, 1, 6, 0, 0, 0, 0, time.UTC)
	holidays[epiphany] = it.CreateHoliday(
		"Epifania",
		epiphany,
		"religious",
		map[string]string{
			"it": "Epifania",
			"en": "Epiphany",
		},
	)

	// Liberation Day - April 25
	liberation := time.Date(year, 4, 25, 0, 0, 0, 0, time.UTC)
	holidays[liberation] = it.CreateHoliday(
		"Festa della Liberazione",
		liberation,
		"public",
		map[string]string{
			"it": "Festa della Liberazione",
			"en": "Liberation Day",
		},
	)

	// Labour Day - May 1
	labourDay := time.Date(year, 5, 1, 0, 0, 0, 0, time.UTC)
	holidays[labourDay] = it.CreateHoliday(
		"Festa del Lavoro",
		labourDay,
		"public",
		map[string]string{
			"it": "Festa del Lavoro",
			"en": "Labour Day",
		},
	)

	// Republic Day - June 2
	republic := time.Date(year, 6, 2, 0, 0, 0, 0, time.UTC)
	holidays[republic] = it.CreateHoliday(
		"Festa della Repubblica",
		republic,
		"public",
		map[string]string{
			"it": "Festa della Repubblica",
			"en": "Republic Day",
		},
	)

	// Assumption of Mary - August 15
	assumption := time.Date(year, 8, 15, 0, 0, 0, 0, time.UTC)
	holidays[assumption] = it.CreateHoliday(
		"Assunzione di Maria",
		assumption,
		"religious",
		map[string]string{
			"it": "Assunzione di Maria",
			"en": "Assumption of Mary",
		},
	)

	// All Saints' Day - November 1
	allSaints := time.Date(year, 11, 1, 0, 0, 0, 0, time.UTC)
	holidays[allSaints] = it.CreateHoliday(
		"Ognissanti",
		allSaints,
		"religious",
		map[string]string{
			"it": "Ognissanti",
			"en": "All Saints' Day",
		},
	)

	// Immaculate Conception - December 8
	immaculate := time.Date(year, 12, 8, 0, 0, 0, 0, time.UTC)
	holidays[immaculate] = it.CreateHoliday(
		"Immacolata Concezione",
		immaculate,
		"religious",
		map[string]string{
			"it": "Immacolata Concezione",
			"en": "Immaculate Conception",
		},
	)

	// Christmas Day - December 25
	christmas := time.Date(year, 12, 25, 0, 0, 0, 0, time.UTC)
	holidays[christmas] = it.CreateHoliday(
		"Natale",
		christmas,
		"religious",
		map[string]string{
			"it": "Natale",
			"en": "Christmas Day",
		},
	)

	// St. Stephen's Day - December 26
	stephens := time.Date(year, 12, 26, 0, 0, 0, 0, time.UTC)
	holidays[stephens] = it.CreateHoliday(
		"Santo Stefano",
		stephens,
		"religious",
		map[string]string{
			"it": "Santo Stefano",
			"en": "St. Stephen's Day",
		},
	)

	// Easter-based holidays
	easterDate := it.CalculateEaster(year)

	// Easter Monday (Lunedì di Pasqua)
	easterMonday := easterDate.AddDate(0, 0, 1)
	holidays[easterMonday] = it.CreateHoliday(
		"Lunedì di Pasqua",
		easterMonday,
		"religious",
		map[string]string{
			"it": "Lunedì di Pasqua",
			"en": "Easter Monday",
		},
	)

	return holidays
}

// CreateHoliday creates a new holiday with Italian localization
func (it *ITProvider) CreateHoliday(name string, date time.Time, category string, languages map[string]string) *Holiday {
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
func (it *ITProvider) CalculateEaster(year int) time.Time {
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
