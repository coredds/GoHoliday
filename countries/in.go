package countries

import (
	"time"
)

// INProvider implements holiday calculations for India
type INProvider struct {
	*BaseProvider
}

// NewINProvider creates a new India holiday provider
func NewINProvider() *INProvider {
	base := NewBaseProvider("IN")
	base.subdivisions = []string{
		"AN", "AP", "AR", "AS", "BR", "CH", "CT", "DH", "DL", "GA",
		"GJ", "HR", "HP", "JK", "JH", "KA", "KL", "LA", "LD", "MP",
		"MH", "MN", "ML", "MZ", "NL", "OR", "PY", "PB", "RJ", "SK",
		"TN", "TG", "TR", "UP", "UT", "WB",
	}
	base.categories = []string{"public", "national", "religious", "hindu", "islamic", "christian"}

	return &INProvider{BaseProvider: base}
}

// LoadHolidays loads all India holidays for a given year
func (in *INProvider) LoadHolidays(year int) map[time.Time]*Holiday {
	holidays := make(map[time.Time]*Holiday)

	// National holidays observed across India

	// Republic Day - January 26
	holidays[time.Date(year, 1, 26, 0, 0, 0, 0, time.UTC)] = &Holiday{
		Name:     "Republic Day",
		Date:     time.Date(year, 1, 26, 0, 0, 0, 0, time.UTC),
		Category: "national",
		Languages: map[string]string{
			"en": "Republic Day",
			"hi": "गणतंत्र दिवस",
		},
		IsObserved: false,
	}

	// Independence Day - August 15
	holidays[time.Date(year, 8, 15, 0, 0, 0, 0, time.UTC)] = &Holiday{
		Name:     "Independence Day",
		Date:     time.Date(year, 8, 15, 0, 0, 0, 0, time.UTC),
		Category: "national",
		Languages: map[string]string{
			"en": "Independence Day",
			"hi": "स्वतंत्रता दिवस",
		},
		IsObserved: false,
	}

	// Gandhi Jayanti - October 2
	holidays[time.Date(year, 10, 2, 0, 0, 0, 0, time.UTC)] = &Holiday{
		Name:     "Gandhi Jayanti",
		Date:     time.Date(year, 10, 2, 0, 0, 0, 0, time.UTC),
		Category: "national",
		Languages: map[string]string{
			"en": "Gandhi Jayanti",
			"hi": "गांधी जयंती",
		},
		IsObserved: false,
	}

	// Add major Hindu festivals (dates vary by lunar calendar)
	in.addHinduFestivals(holidays, year)

	// Add major Islamic festivals (dates vary by lunar calendar)
	in.addIslamicFestivals(holidays, year)

	// Add major Christian festivals
	in.addChristianFestivals(holidays, year)

	return holidays
}

// addHinduFestivals adds Hindu festivals for the given year
// Note: These are approximate dates - in production, use proper lunar calendar calculations
func (in *INProvider) addHinduFestivals(holidays map[time.Time]*Holiday, year int) {
	// Diwali (approximate - varies by lunar calendar)
	// This is a simplified implementation - real dates require lunar calendar calculations
	var diwaliDate time.Time
	switch year {
	case 2024:
		diwaliDate = time.Date(2024, 11, 1, 0, 0, 0, 0, time.UTC)
	case 2025:
		diwaliDate = time.Date(2025, 10, 20, 0, 0, 0, 0, time.UTC)
	case 2026:
		diwaliDate = time.Date(2026, 11, 8, 0, 0, 0, 0, time.UTC)
	default:
		// Fallback calculation needed for other years
		return
	}

	holidays[diwaliDate] = &Holiday{
		Name:     "Diwali",
		Date:     diwaliDate,
		Category: "religious",
		Languages: map[string]string{
			"en": "Diwali",
			"hi": "दीपावली",
		},
		IsObserved: false,
	}

	// Holi (approximate dates)
	var holiDate time.Time
	switch year {
	case 2024:
		holiDate = time.Date(2024, 3, 25, 0, 0, 0, 0, time.UTC)
	case 2025:
		holiDate = time.Date(2025, 3, 14, 0, 0, 0, 0, time.UTC)
	case 2026:
		holiDate = time.Date(2026, 3, 3, 0, 0, 0, 0, time.UTC)
	default:
		return
	}

	holidays[holiDate] = &Holiday{
		Name:     "Holi",
		Date:     holiDate,
		Category: "religious",
		Languages: map[string]string{
			"en": "Holi",
			"hi": "होली",
		},
		IsObserved: false,
	}
}

// addIslamicFestivals adds Islamic festivals for the given year
func (in *INProvider) addIslamicFestivals(holidays map[time.Time]*Holiday, year int) {
	// Eid al-Fitr (approximate dates)
	var eidFitrDate time.Time
	switch year {
	case 2024:
		eidFitrDate = time.Date(2024, 4, 10, 0, 0, 0, 0, time.UTC)
	case 2025:
		eidFitrDate = time.Date(2025, 3, 30, 0, 0, 0, 0, time.UTC)
	case 2026:
		eidFitrDate = time.Date(2026, 3, 20, 0, 0, 0, 0, time.UTC)
	default:
		return
	}

	holidays[eidFitrDate] = &Holiday{
		Name:     "Eid al-Fitr",
		Date:     eidFitrDate,
		Category: "religious",
		Languages: map[string]string{
			"en": "Eid al-Fitr",
			"hi": "ईद उल-फितर",
			"ar": "عيد الفطر",
		},
		IsObserved: false,
	}
}

// addChristianFestivals adds Christian festivals for the given year
func (in *INProvider) addChristianFestivals(holidays map[time.Time]*Holiday, year int) {
	// Christmas Day - December 25
	holidays[time.Date(year, 12, 25, 0, 0, 0, 0, time.UTC)] = &Holiday{
		Name:     "Christmas Day",
		Date:     time.Date(year, 12, 25, 0, 0, 0, 0, time.UTC),
		Category: "religious",
		Languages: map[string]string{
			"en": "Christmas Day",
			"hi": "क्रिसमस",
		},
		IsObserved: false,
	}

	// Good Friday (varies by Easter calculation)
	easterDate := in.calculateEaster(year)
	goodFridayDate := easterDate.AddDate(0, 0, -2)

	holidays[goodFridayDate] = &Holiday{
		Name:     "Good Friday",
		Date:     goodFridayDate,
		Category: "religious",
		Languages: map[string]string{
			"en": "Good Friday",
			"hi": "गुड फ्राइडे",
		},
		IsObserved: false,
	}
}

// calculateEaster calculates Easter date using the algorithm
func (in *INProvider) calculateEaster(year int) time.Time {
	// Simplified Easter calculation (Gregorian calendar)
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
	n := (h + l - 7*m + 114) / 31
	p := (h + l - 7*m + 114) % 31

	return time.Date(year, time.Month(n), p+1, 0, 0, 0, 0, time.UTC)
}
