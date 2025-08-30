package countries

import (
	"time"
)

// INProvider implements holiday calculations for India
type INProvider struct {
	*BaseProvider
}

// NewINProvider creates a new Indian holiday provider
func NewINProvider() *INProvider {
	base := NewBaseProvider("IN")
	base.subdivisions = []string{
		// 28 States
		"AP", "AR", "AS", "BR", "CT", "GA", "GJ", "HR", "HP", "JH", "KA", "KL", "MP", "MH",
		"MN", "ML", "MZ", "NL", "OR", "PB", "RJ", "SK", "TN", "TG", "TR", "UP", "UT", "WB",
		// 8 Union Territories
		"AN", "CH", "DN", "DL", "JK", "LA", "LD", "PY",
		// Andhra Pradesh, Arunachal Pradesh, Assam, Bihar, Chhattisgarh, Goa, Gujarat,
		// Haryana, Himachal Pradesh, Jharkhand, Karnataka, Kerala, Madhya Pradesh, Maharashtra,
		// Manipur, Meghalaya, Mizoram, Nagaland, Odisha, Punjab, Rajasthan, Sikkim,
		// Tamil Nadu, Telangana, Tripura, Uttar Pradesh, Uttarakhand, West Bengal,
		// Andaman and Nicobar Islands, Chandigarh, Dadra and Nagar Haveli, Delhi,
		// Jammu and Kashmir, Ladakh, Lakshadweep, Puducherry
	}
	base.categories = []string{"national", "gazetted", "restricted", "hindu", "muslim", "christian", "sikh", "buddhist", "jain", "regional"}

	return &INProvider{BaseProvider: base}
}

// LoadHolidays loads all Indian holidays for a given year
func (in *INProvider) LoadHolidays(year int) map[time.Time]*Holiday {
	holidays := make(map[time.Time]*Holiday)

	// Fixed national holidays
	in.addNationalHolidays(holidays, year)

	// Christian holidays (Easter-based)
	in.addChristianHolidays(holidays, year)

	// Note: Hindu, Muslim, and other religious holidays are based on lunar calendars
	// and would require complex astronomical calculations. For now, we include
	// the major fixed ones and note that lunar-based holidays need special handling.

	return holidays
}

// addNationalHolidays adds fixed national holidays of India
func (in *INProvider) addNationalHolidays(holidays map[time.Time]*Holiday, year int) {
	nationalHolidays := []struct {
		month    int
		day      int
		name     string
		nameHi   string
		category string
	}{
		{1, 26, "Republic Day", "गणतंत्र दिवस", "national"},
		{8, 15, "Independence Day", "स्वतंत्रता दिवस", "national"},
		{10, 2, "Gandhi Jayanti", "गांधी जयंती", "national"},
		{12, 25, "Christmas Day", "क्रिसमस", "christian"},
	}

	for _, h := range nationalHolidays {
		date := time.Date(year, time.Month(h.month), h.day, 0, 0, 0, 0, time.UTC)
		holidays[date] = &Holiday{
			Name:     h.name,
			Date:     date,
			Category: h.category,
			Languages: map[string]string{
				"en": h.name,
				"hi": h.nameHi,
			},
			IsObserved: true,
		}
	}
}

// addChristianHolidays adds Christian holidays (Easter-based)
func (in *INProvider) addChristianHolidays(holidays map[time.Time]*Holiday, year int) {
	easter := in.calculateEaster(year)

	christianHolidays := []struct {
		offset   int
		name     string
		category string
	}{
		{-2, "Good Friday", "christian"},
		{0, "Easter Sunday", "christian"},
	}

	for _, h := range christianHolidays {
		date := easter.AddDate(0, 0, h.offset)
		holidays[date] = &Holiday{
			Name:     h.name,
			Date:     date,
			Category: h.category,
			Languages: map[string]string{
				"en": h.name,
			},
			IsObserved: true,
		}
	}
}

// calculateEaster calculates Easter Sunday for a given year using the Western (Gregorian) calendar
func (in *INProvider) calculateEaster(year int) time.Time {
	// Using the algorithm for Western Easter (Gregorian calendar)
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

// GetStateHolidays returns state-specific holidays for India
// Note: This is a simplified implementation. Real-world usage would require
// more comprehensive state-specific holiday data.
func (in *INProvider) GetStateHolidays(year int, state string) map[time.Time]*Holiday {
	holidays := make(map[time.Time]*Holiday)

	// Some examples of state-specific holidays
	stateHolidays := map[string][]struct {
		month    int
		day      int
		name     string
		category string
	}{
		"MH": { // Maharashtra
			{5, 1, "Maharashtra Day", "regional"},
		},
		"GJ": { // Gujarat
			{5, 1, "Gujarat Day", "regional"},
		},
		"WB": { // West Bengal
			{8, 16, "Poila Boishakh", "regional"}, // Bengali New Year (approximate)
		},
		"TN": { // Tamil Nadu
			{4, 14, "Tamil New Year", "regional"},
		},
		"KL": { // Kerala
			{8, 15, "Onam", "regional"}, // Approximate date
		},
		"PB": { // Punjab
			{4, 13, "Baisakhi", "regional"},
		},
	}

	if stateHols, exists := stateHolidays[state]; exists {
		for _, h := range stateHols {
			date := time.Date(year, time.Month(h.month), h.day, 0, 0, 0, 0, time.UTC)
			holidays[date] = &Holiday{
				Name:     h.name,
				Date:     date,
				Category: h.category,
				Languages: map[string]string{
					"en": h.name,
				},
				IsObserved:   true,
				Subdivisions: []string{state},
			}
		}
	}

	return holidays
}

// Note: For a production implementation of Indian holidays, you would need:
// 1. Hindu calendar calculations for festivals like Diwali, Holi, Dussehra
// 2. Islamic calendar calculations for Eid al-Fitr, Eid al-Adha, Muharram
// 3. Sikh calendar for Guru Nanak Jayanti, etc.
// 4. Buddhist calendar for Buddha Purnima
// 5. Jain calendar for Mahavir Jayanti
// 6. Regional lunar calendar calculations for state-specific festivals
//
// These require complex astronomical calculations and calendar conversions
// that are beyond the scope of this basic implementation.

// GetMajorFestivals returns approximate dates for major Indian festivals
// Note: These are approximate dates as actual dates depend on lunar calculations
func (in *INProvider) GetMajorFestivals(year int) map[time.Time]*Holiday {
	holidays := make(map[time.Time]*Holiday)

	// These are very approximate dates - real implementation would need lunar calendar
	approximateFestivals := []struct {
		month    int
		day      int
		name     string
		category string
	}{
		{3, 8, "Holi", "hindu"},                // Approximate - varies by lunar calendar
		{10, 24, "Dussehra", "hindu"},          // Approximate
		{11, 12, "Diwali", "hindu"},            // Approximate
		{8, 19, "Janmashtami", "hindu"},        // Approximate
		{4, 6, "Ram Navami", "hindu"},          // Approximate
		{5, 23, "Buddha Purnima", "buddhist"},  // Approximate
		{4, 14, "Mahavir Jayanti", "jain"},     // Approximate
		{11, 15, "Guru Nanak Jayanti", "sikh"}, // Approximate
	}

	for _, h := range approximateFestivals {
		date := time.Date(year, time.Month(h.month), h.day, 0, 0, 0, 0, time.UTC)
		holidays[date] = &Holiday{
			Name:     h.name,
			Date:     date,
			Category: h.category,
			Languages: map[string]string{
				"en": h.name,
			},
			IsObserved: true,
		}
	}

	return holidays
}
