package countries

import (
	"time"
)

// JPProvider implements HolidayProvider for Japan
type JPProvider struct {
	BaseProvider
}

// NewJPProvider creates a new Japan holiday provider
func NewJPProvider() *JPProvider {
	return &JPProvider{
		BaseProvider: BaseProvider{
			countryCode: "JP",
		},
	}
}

// LoadHolidays returns holidays for Japan for the given year
func (p *JPProvider) LoadHolidays(year int) map[time.Time]*Holiday {
	holidays := make(map[time.Time]*Holiday)

	// New Year's Day (元日, Ganjitsu)
	newYear := time.Date(year, 1, 1, 0, 0, 0, 0, time.UTC)
	holidays[newYear] = &Holiday{
		Name:     "New Year's Day",
		Date:     newYear,
		Category: "public",
		Languages: map[string]string{
			"en": "New Year's Day",
			"ja": "元日",
		},
	}

	// Coming of Age Day (成人の日, Seijin no Hi) - Second Monday of January
	comingOfAge := NthWeekdayOfMonth(year, 1, time.Monday, 2)
	holidays[comingOfAge] = &Holiday{
		Name:     "Coming of Age Day",
		Date:     comingOfAge,
		Category: "public",
		Languages: map[string]string{
			"en": "Coming of Age Day",
			"ja": "成人の日",
		},
	}

	// National Foundation Day (建国記念の日, Kenkoku Kinen no Hi)
	foundation := time.Date(year, 2, 11, 0, 0, 0, 0, time.UTC)
	holidays[foundation] = &Holiday{
		Name:     "National Foundation Day",
		Date:     foundation,
		Category: "public",
		Languages: map[string]string{
			"en": "National Foundation Day",
			"ja": "建国記念の日",
		},
	}

	// Emperor's Birthday (天皇誕生日, Tennō Tanjōbi)
	var emperorBirthday time.Time
	if year >= 2020 {
		emperorBirthday = time.Date(year, 2, 23, 0, 0, 0, 0, time.UTC) // Emperor Naruhito
	} else {
		emperorBirthday = time.Date(year, 12, 23, 0, 0, 0, 0, time.UTC) // Emperor Akihito
	}
	holidays[emperorBirthday] = &Holiday{
		Name:     "Emperor's Birthday",
		Date:     emperorBirthday,
		Category: "public",
		Languages: map[string]string{
			"en": "Emperor's Birthday",
			"ja": "天皇誕生日",
		},
	}

	// Vernal Equinox Day (春分の日, Shunbun no Hi)
	vernalEquinox := getVernalEquinox(year)
	holidays[vernalEquinox] = &Holiday{
		Name:     "Vernal Equinox Day",
		Date:     vernalEquinox,
		Category: "public",
		Languages: map[string]string{
			"en": "Vernal Equinox Day",
			"ja": "春分の日",
		},
	}

	// Showa Day (昭和の日, Shōwa no Hi)
	showa := time.Date(year, 4, 29, 0, 0, 0, 0, time.UTC)
	holidays[showa] = &Holiday{
		Name:     "Showa Day",
		Date:     showa,
		Category: "public",
		Languages: map[string]string{
			"en": "Showa Day",
			"ja": "昭和の日",
		},
	}

	// Constitution Memorial Day (憲法記念日, Kenpō Kinenbi)
	constitution := time.Date(year, 5, 3, 0, 0, 0, 0, time.UTC)
	holidays[constitution] = &Holiday{
		Name:     "Constitution Memorial Day",
		Date:     constitution,
		Category: "public",
		Languages: map[string]string{
			"en": "Constitution Memorial Day",
			"ja": "憲法記念日",
		},
	}

	// Greenery Day (みどりの日, Midori no Hi)
	greenery := time.Date(year, 5, 4, 0, 0, 0, 0, time.UTC)
	holidays[greenery] = &Holiday{
		Name:     "Greenery Day",
		Date:     greenery,
		Category: "public",
		Languages: map[string]string{
			"en": "Greenery Day",
			"ja": "みどりの日",
		},
	}

	// Children's Day (こどもの日, Kodomo no Hi)
	children := time.Date(year, 5, 5, 0, 0, 0, 0, time.UTC)
	holidays[children] = &Holiday{
		Name:     "Children's Day",
		Date:     children,
		Category: "public",
		Languages: map[string]string{
			"en": "Children's Day",
			"ja": "こどもの日",
		},
	}

	// Marine Day (海の日, Umi no Hi) - Third Monday of July
	marine := NthWeekdayOfMonth(year, 7, time.Monday, 3)
	holidays[marine] = &Holiday{
		Name:     "Marine Day",
		Date:     marine,
		Category: "public",
		Languages: map[string]string{
			"en": "Marine Day",
			"ja": "海の日",
		},
	}

	// Mountain Day (山の日, Yama no Hi)
	mountain := time.Date(year, 8, 11, 0, 0, 0, 0, time.UTC)
	holidays[mountain] = &Holiday{
		Name:     "Mountain Day",
		Date:     mountain,
		Category: "public",
		Languages: map[string]string{
			"en": "Mountain Day",
			"ja": "山の日",
		},
	}

	// Respect for the Aged Day (敬老の日, Keirō no Hi) - Third Monday of September
	aged := NthWeekdayOfMonth(year, 9, time.Monday, 3)
	holidays[aged] = &Holiday{
		Name:     "Respect for the Aged Day",
		Date:     aged,
		Category: "public",
		Languages: map[string]string{
			"en": "Respect for the Aged Day",
			"ja": "敬老の日",
		},
	}

	// Autumnal Equinox Day (秋分の日, Shūbun no Hi)
	autumnalEquinox := getAutumnalEquinox(year)
	holidays[autumnalEquinox] = &Holiday{
		Name:     "Autumnal Equinox Day",
		Date:     autumnalEquinox,
		Category: "public",
		Languages: map[string]string{
			"en": "Autumnal Equinox Day",
			"ja": "秋分の日",
		},
	}

	// Sports Day (スポーツの日, Supōtsu no Hi) - Second Monday of October
	sports := NthWeekdayOfMonth(year, 10, time.Monday, 2)
	sportsName := "Sports Day"
	if year < 2020 {
		sportsName = "Health and Sports Day"
	}
	holidays[sports] = &Holiday{
		Name:     sportsName,
		Date:     sports,
		Category: "public",
		Languages: map[string]string{
			"en": sportsName,
			"ja": "スポーツの日",
		},
	}

	// Culture Day (文化の日, Bunka no Hi)
	culture := time.Date(year, 11, 3, 0, 0, 0, 0, time.UTC)
	holidays[culture] = &Holiday{
		Name:     "Culture Day",
		Date:     culture,
		Category: "public",
		Languages: map[string]string{
			"en": "Culture Day",
			"ja": "文化の日",
		},
	}

	// Labor Thanksgiving Day (勤労感謝の日, Kinrō Kansha no Hi)
	laborThanksgiving := time.Date(year, 11, 23, 0, 0, 0, 0, time.UTC)
	holidays[laborThanksgiving] = &Holiday{
		Name:     "Labor Thanksgiving Day",
		Date:     laborThanksgiving,
		Category: "public",
		Languages: map[string]string{
			"en": "Labor Thanksgiving Day",
			"ja": "勤労感謝の日",
		},
	}

	return holidays
}

// GetCountryCode returns the country code
func (p *JPProvider) GetCountryCode() string {
	return "JP"
}

// GetSupportedSubdivisions returns supported subdivisions (none for Japan currently)
func (p *JPProvider) GetSupportedSubdivisions() []string {
	return []string{}
}

// GetSupportedCategories returns supported holiday categories
func (p *JPProvider) GetSupportedCategories() []string {
	return []string{"public"}
}

// Helper functions for Japanese holidays

// getVernalEquinox calculates the vernal equinox date for Japan
func getVernalEquinox(year int) time.Time {
	// Simplified calculation - the vernal equinox typically falls around March 20-21
	day := 20
	if year >= 1851 && year <= 1899 {
		day = 19
	} else if year >= 1900 && year <= 1979 {
		day = 21
	} else if year >= 1980 && year <= 2099 {
		day = 20
		if year%4 == 0 && year > 2092 {
			day = 21
		}
	}
	return time.Date(year, 3, day, 0, 0, 0, 0, time.UTC)
}

// getAutumnalEquinox calculates the autumnal equinox date for Japan
func getAutumnalEquinox(year int) time.Time {
	// Simplified calculation - the autumnal equinox typically falls around September 22-23
	day := 23
	if year >= 1851 && year <= 1899 {
		day = 22
	} else if year >= 1900 && year <= 1979 {
		day = 23
	} else if year >= 1980 && year <= 2099 {
		day = 23
		if year%4 == 0 && year > 2092 {
			day = 22
		}
	}
	return time.Date(year, 9, day, 0, 0, 0, 0, time.UTC)
}
