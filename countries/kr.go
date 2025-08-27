package countries

import (
	"time"
)

// KRProvider implements holiday calculations for South Korea
type KRProvider struct {
	*BaseProvider
}

// NewKRProvider creates a new South Korean holiday provider
func NewKRProvider() *KRProvider {
	base := NewBaseProvider("KR")
	base.subdivisions = []string{
		// 9 provinces + 8 metropolitan cities/special cities
		"KW", "KA", "CN", "CB", "KN", "KB", "JJ", "JN", "JB",
		"SE", "BS", "DG", "IC", "GJ", "DJ", "UL", "SJ",
	}
	base.categories = []string{"public", "national", "traditional", "commemorative"}

	return &KRProvider{BaseProvider: base}
}

// LoadHolidays loads all South Korean holidays for a given year
func (kr *KRProvider) LoadHolidays(year int) map[time.Time]*Holiday {
	holidays := make(map[time.Time]*Holiday)

	// Fixed national holidays

	// New Year's Day - January 1
	newYear := time.Date(year, 1, 1, 0, 0, 0, 0, time.UTC)
	holidays[newYear] = kr.CreateHoliday(
		"신정",
		newYear,
		"public",
		map[string]string{
			"ko": "신정",
			"en": "New Year's Day",
		},
	)

	// Independence Movement Day - March 1
	independence := time.Date(year, 3, 1, 0, 0, 0, 0, time.UTC)
	holidays[independence] = kr.CreateHoliday(
		"삼일절",
		independence,
		"national",
		map[string]string{
			"ko": "삼일절",
			"en": "Independence Movement Day",
		},
	)

	// Children's Day - May 5
	childrens := time.Date(year, 5, 5, 0, 0, 0, 0, time.UTC)
	holidays[childrens] = kr.CreateHoliday(
		"어린이날",
		childrens,
		"public",
		map[string]string{
			"ko": "어린이날",
			"en": "Children's Day",
		},
	)

	// Memorial Day - June 6
	memorial := time.Date(year, 6, 6, 0, 0, 0, 0, time.UTC)
	holidays[memorial] = kr.CreateHoliday(
		"현충일",
		memorial,
		"commemorative",
		map[string]string{
			"ko": "현충일",
			"en": "Memorial Day",
		},
	)

	// Liberation Day - August 15
	liberation := time.Date(year, 8, 15, 0, 0, 0, 0, time.UTC)
	holidays[liberation] = kr.CreateHoliday(
		"광복절",
		liberation,
		"national",
		map[string]string{
			"ko": "광복절",
			"en": "Liberation Day",
		},
	)

	// National Foundation Day - October 3
	foundation := time.Date(year, 10, 3, 0, 0, 0, 0, time.UTC)
	holidays[foundation] = kr.CreateHoliday(
		"개천절",
		foundation,
		"national",
		map[string]string{
			"ko": "개천절",
			"en": "National Foundation Day",
		},
	)

	// Hangeul Day - October 9
	hangeul := time.Date(year, 10, 9, 0, 0, 0, 0, time.UTC)
	holidays[hangeul] = kr.CreateHoliday(
		"한글날",
		hangeul,
		"national",
		map[string]string{
			"ko": "한글날",
			"en": "Hangeul Day",
		},
	)

	// Christmas Day - December 25
	christmas := time.Date(year, 12, 25, 0, 0, 0, 0, time.UTC)
	holidays[christmas] = kr.CreateHoliday(
		"성탄절",
		christmas,
		"public",
		map[string]string{
			"ko": "성탄절",
			"en": "Christmas Day",
		},
	)

	// Buddha's Birthday - varies each year (8th day of 4th lunar month)
	// Approximation: typically in May
	// For 2024, Buddha's Birthday is May 15
	if year == 2024 {
		buddha := time.Date(2024, 5, 15, 0, 0, 0, 0, time.UTC)
		holidays[buddha] = kr.CreateHoliday(
			"부처님 오신 날",
			buddha,
			"traditional",
			map[string]string{
				"ko": "부처님 오신 날",
				"en": "Buddha's Birthday",
			},
		)
	}

	// Lunar New Year (Seollal) - varies each year
	// For 2024, Lunar New Year is February 10-12
	if year == 2024 {
		// Seollal spans 3 days
		seollal1 := time.Date(2024, 2, 9, 0, 0, 0, 0, time.UTC)
		holidays[seollal1] = kr.CreateHoliday(
			"설날 연휴",
			seollal1,
			"traditional",
			map[string]string{
				"ko": "설날 연휴",
				"en": "Lunar New Year Holiday",
			},
		)

		seollal2 := time.Date(2024, 2, 10, 0, 0, 0, 0, time.UTC)
		holidays[seollal2] = kr.CreateHoliday(
			"설날",
			seollal2,
			"traditional",
			map[string]string{
				"ko": "설날",
				"en": "Lunar New Year",
			},
		)

		seollal3 := time.Date(2024, 2, 11, 0, 0, 0, 0, time.UTC)
		holidays[seollal3] = kr.CreateHoliday(
			"설날 연휴",
			seollal3,
			"traditional",
			map[string]string{
				"ko": "설날 연휴",
				"en": "Lunar New Year Holiday",
			},
		)
	}

	// Chuseok (Korean Thanksgiving) - varies each year
	// For 2024, Chuseok is September 16-18
	if year == 2024 {
		// Chuseok spans 3 days
		chuseok1 := time.Date(2024, 9, 16, 0, 0, 0, 0, time.UTC)
		holidays[chuseok1] = kr.CreateHoliday(
			"추석 연휴",
			chuseok1,
			"traditional",
			map[string]string{
				"ko": "추석 연휴",
				"en": "Chuseok Holiday",
			},
		)

		chuseok2 := time.Date(2024, 9, 17, 0, 0, 0, 0, time.UTC)
		holidays[chuseok2] = kr.CreateHoliday(
			"추석",
			chuseok2,
			"traditional",
			map[string]string{
				"ko": "추석",
				"en": "Chuseok",
			},
		)

		chuseok3 := time.Date(2024, 9, 18, 0, 0, 0, 0, time.UTC)
		holidays[chuseok3] = kr.CreateHoliday(
			"추석 연휴",
			chuseok3,
			"traditional",
			map[string]string{
				"ko": "추석 연휴",
				"en": "Chuseok Holiday",
			},
		)
	}

	return holidays
}

// CreateHoliday creates a new holiday with Korean localization
func (kr *KRProvider) CreateHoliday(name string, date time.Time, category string, languages map[string]string) *Holiday {
	return &Holiday{
		Name:         name,
		Date:         date,
		Category:     category,
		Languages:    languages,
		IsObserved:   true,
		Subdivisions: []string{},
	}
}
