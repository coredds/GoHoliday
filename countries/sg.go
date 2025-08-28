package countries

import (
	"time"
)

// SGProvider implements holiday calculations for Singapore
type SGProvider struct {
	*BaseProvider
}

// NewSGProvider creates a new Singaporean holiday provider
func NewSGProvider() *SGProvider {
	base := NewBaseProvider("SG")
	base.subdivisions = []string{
		// 5 regions (not administrative divisions but planning areas)
		"CR", "ER", "NR", "NER", "WR", // Central, East, North, Northeast, West
	}
	base.categories = []string{"public", "religious", "cultural", "national"}

	return &SGProvider{BaseProvider: base}
}

// LoadHolidays loads all Singaporean holidays for a given year
func (sg *SGProvider) LoadHolidays(year int) map[time.Time]*Holiday {
	holidays := make(map[time.Time]*Holiday)

	// Fixed national holidays

	// New Year's Day - January 1
	newYear := time.Date(year, 1, 1, 0, 0, 0, 0, time.UTC)
	holidays[newYear] = sg.CreateHoliday(
		"New Year's Day",
		newYear,
		"public",
		map[string]string{
			"en": "New Year's Day",
			"zh": "元旦",
			"ms": "Hari Tahun Baru",
			"ta": "புத்தாண்டு",
		},
	)

	// Labour Day - May 1
	labourDay := time.Date(year, 5, 1, 0, 0, 0, 0, time.UTC)
	holidays[labourDay] = sg.CreateHoliday(
		"Labour Day",
		labourDay,
		"public",
		map[string]string{
			"en": "Labour Day",
			"zh": "劳动节",
			"ms": "Hari Pekerja",
			"ta": "தொழிலாளர் தினம்",
		},
	)

	// National Day - August 9
	nationalDay := time.Date(year, 8, 9, 0, 0, 0, 0, time.UTC)
	holidays[nationalDay] = sg.CreateHoliday(
		"National Day",
		nationalDay,
		"national",
		map[string]string{
			"en": "National Day",
			"zh": "国庆节",
			"ms": "Hari Kebangsaan",
			"ta": "தேசிய தினம்",
		},
	)

	// Christmas Day - December 25
	christmas := time.Date(year, 12, 25, 0, 0, 0, 0, time.UTC)
	holidays[christmas] = sg.CreateHoliday(
		"Christmas Day",
		christmas,
		"religious",
		map[string]string{
			"en": "Christmas Day",
			"zh": "圣诞节",
			"ms": "Hari Krismas",
			"ta": "கிறிஸ்துமஸ்",
		},
	)

	// Variable holidays - approximations for demonstration
	// In production, these would use proper lunar calendar calculations

	// Chinese New Year - varies each year (lunar calendar)
	if year == 2024 {
		// Chinese New Year 2024: February 10-11
		cny1 := time.Date(2024, 2, 10, 0, 0, 0, 0, time.UTC)
		holidays[cny1] = sg.CreateHoliday(
			"Chinese New Year",
			cny1,
			"cultural",
			map[string]string{
				"en": "Chinese New Year",
				"zh": "农历新年",
				"ms": "Tahun Baru Cina",
				"ta": "சீன புத்தாண்டு",
			},
		)

		cny2 := time.Date(2024, 2, 11, 0, 0, 0, 0, time.UTC)
		holidays[cny2] = sg.CreateHoliday(
			"Chinese New Year (Day 2)",
			cny2,
			"cultural",
			map[string]string{
				"en": "Chinese New Year (Day 2)",
				"zh": "农历新年第二天",
				"ms": "Tahun Baru Cina (Hari 2)",
				"ta": "சீன புத்தாண்டு (2ம் நாள்)",
			},
		)
	}

	// Good Friday - Easter-based
	easterDate := sg.CalculateEaster(year)
	goodFriday := easterDate.AddDate(0, 0, -2)
	holidays[goodFriday] = sg.CreateHoliday(
		"Good Friday",
		goodFriday,
		"religious",
		map[string]string{
			"en": "Good Friday",
			"zh": "耶稣受难节",
			"ms": "Jumaat Agung",
			"ta": "புனித வெள்ளி",
		},
	)

	// Vesak Day - Buddhist holiday (varies each year)
	if year == 2024 {
		vesak := time.Date(2024, 5, 22, 0, 0, 0, 0, time.UTC) // Approximation
		holidays[vesak] = sg.CreateHoliday(
			"Vesak Day",
			vesak,
			"religious",
			map[string]string{
				"en": "Vesak Day",
				"zh": "卫塞节",
				"ms": "Hari Vesak",
				"ta": "வேசாக் தினம்",
			},
		)
	}

	// Hari Raya Puasa (Eid al-Fitr) - Islamic holiday (varies each year)
	if year == 2024 {
		hariRaya := time.Date(2024, 4, 10, 0, 0, 0, 0, time.UTC) // Approximation
		holidays[hariRaya] = sg.CreateHoliday(
			"Hari Raya Puasa",
			hariRaya,
			"religious",
			map[string]string{
				"en": "Hari Raya Puasa",
				"zh": "开斋节",
				"ms": "Hari Raya Puasa",
				"ta": "ஹரி ராயா புவாசா",
			},
		)
	}

	// Hari Raya Haji (Eid al-Adha) - Islamic holiday (varies each year)
	if year == 2024 {
		hariRayaHaji := time.Date(2024, 6, 17, 0, 0, 0, 0, time.UTC) // Approximation
		holidays[hariRayaHaji] = sg.CreateHoliday(
			"Hari Raya Haji",
			hariRayaHaji,
			"religious",
			map[string]string{
				"en": "Hari Raya Haji",
				"zh": "哈芝节",
				"ms": "Hari Raya Haji",
				"ta": "ஹரி ராயா ஹாஜி",
			},
		)
	}

	// Deepavali (Diwali) - Hindu festival (varies each year)
	if year == 2024 {
		deepavali := time.Date(2024, 10, 31, 0, 0, 0, 0, time.UTC) // Approximation
		holidays[deepavali] = sg.CreateHoliday(
			"Deepavali",
			deepavali,
			"religious",
			map[string]string{
				"en": "Deepavali",
				"zh": "屠妖节",
				"ms": "Deepavali",
				"ta": "தீபாவளி",
			},
		)
	}

	return holidays
}

// CreateHoliday creates a new holiday with Singaporean localization
func (sg *SGProvider) CreateHoliday(name string, date time.Time, category string, languages map[string]string) *Holiday {
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
func (sg *SGProvider) CalculateEaster(year int) time.Time {
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
