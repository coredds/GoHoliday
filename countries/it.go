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
		// 20 regions of Italy
		"ABR", "BAS", "CAL", "CAM", "EMR", "FVG", "LAZ", "LIG", "LOM", "MAR",
		"MOL", "PIE", "PUG", "SAR", "SIC", "TOS", "TAA", "UMB", "VDA", "VEN",
		// Abruzzo, Basilicata, Calabria, Campania, Emilia-Romagna, Friuli-Venezia Giulia,
		// Lazio, Liguria, Lombardy, Marche, Molise, Piedmont, Apulia, Sardinia,
		// Sicily, Tuscany, Trentino-Alto Adige, Umbria, Aosta Valley, Veneto
	}
	base.categories = []string{"public", "religious", "regional", "patron"}

	return &ITProvider{BaseProvider: base}
}

// LoadHolidays loads all Italian holidays for a given year
func (it *ITProvider) LoadHolidays(year int) map[time.Time]*Holiday {
	holidays := make(map[time.Time]*Holiday)

	// Fixed date holidays
	it.addFixedHolidays(holidays, year)

	// Easter-based holidays
	it.addEasterHolidays(holidays, year)

	return holidays
}

// addFixedHolidays adds fixed-date Italian holidays
func (it *ITProvider) addFixedHolidays(holidays map[time.Time]*Holiday, year int) {
	fixedHolidays := []struct {
		month    int
		day      int
		name     string
		nameIt   string
		category string
	}{
		{1, 1, "New Year's Day", "Capodanno", "public"},
		{1, 6, "Epiphany", "Epifania", "religious"},
		{4, 25, "Liberation Day", "Festa della Liberazione", "public"},
		{5, 1, "Labour Day", "Festa del Lavoro", "public"},
		{6, 2, "Republic Day", "Festa della Repubblica", "public"},
		{8, 15, "Assumption of Mary", "Assunzione di Maria", "religious"},
		{11, 1, "All Saints' Day", "Ognissanti", "religious"},
		{12, 8, "Immaculate Conception", "Immacolata Concezione", "religious"},
		{12, 25, "Christmas Day", "Natale", "religious"},
		{12, 26, "St. Stephen's Day", "Santo Stefano", "religious"},
	}

	for _, h := range fixedHolidays {
		date := time.Date(year, time.Month(h.month), h.day, 0, 0, 0, 0, time.UTC)
		holidays[date] = &Holiday{
			Name:     h.name,
			Date:     date,
			Category: h.category,
			Languages: map[string]string{
				"en": h.name,
				"it": h.nameIt,
			},
			IsObserved: true,
		}
	}
}

// addEasterHolidays adds Easter-based Italian holidays
func (it *ITProvider) addEasterHolidays(holidays map[time.Time]*Holiday, year int) {
	easter := it.calculateEaster(year)

	easterHolidays := []struct {
		offset   int
		name     string
		nameIt   string
		category string
	}{
		{1, "Easter Monday", "Lunedì dell'Angelo", "religious"}, // Easter Monday (Pasquetta)
	}

	for _, h := range easterHolidays {
		date := easter.AddDate(0, 0, h.offset)
		holidays[date] = &Holiday{
			Name:     h.name,
			Date:     date,
			Category: h.category,
			Languages: map[string]string{
				"en": h.name,
				"it": h.nameIt,
			},
			IsObserved: true,
		}
	}
}

// calculateEaster calculates Easter Sunday for a given year using the Western (Gregorian) calendar
func (it *ITProvider) calculateEaster(year int) time.Time {
	// Using the algorithm for Western Easter (Gregorian calendar)
	// This is the standard algorithm used in Western Christianity

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

// GetRegionalHolidays returns region-specific holidays for Italy
func (it *ITProvider) GetRegionalHolidays(year int, region string) map[time.Time]*Holiday {
	holidays := make(map[time.Time]*Holiday)

	// Some examples of regional patron saint days
	regionalHolidays := map[string][]struct {
		month  int
		day    int
		name   string
		nameIt string
	}{
		"LOM": { // Lombardy
			{12, 7, "St. Ambrose Day", "Sant'Ambrogio"}, // Milan patron saint
		},
		"VEN": { // Veneto
			{4, 25, "St. Mark's Day", "San Marco"}, // Venice patron saint
		},
		"SIC": { // Sicily
			{7, 15, "St. Rosalia Day", "Santa Rosalia"}, // Palermo patron saint
		},
		"LAZ": { // Lazio (Rome)
			{6, 29, "St. Peter and Paul Day", "Santi Pietro e Paolo"}, // Rome patron saints
		},
		"CAM": { // Campania (Naples)
			{9, 19, "St. Januarius Day", "San Gennaro"}, // Naples patron saint
		},
	}

	if regionHolidays, exists := regionalHolidays[region]; exists {
		for _, h := range regionHolidays {
			date := time.Date(year, time.Month(h.month), h.day, 0, 0, 0, 0, time.UTC)
			holidays[date] = &Holiday{
				Name:     h.name,
				Date:     date,
				Category: "patron",
				Languages: map[string]string{
					"en": h.name,
					"it": h.nameIt,
				},
				IsObserved:   true,
				Subdivisions: []string{region},
			}
		}
	}

	return holidays
}
