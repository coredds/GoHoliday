package countries

import (
	"time"
)

// CNProvider implements holiday calculations for China
type CNProvider struct {
	*BaseProvider
}

// NewCNProvider creates a new China holiday provider
func NewCNProvider() *CNProvider {
	base := NewBaseProvider("CN")
	base.subdivisions = []string{
		// Provinces
		"11", "12", "13", "14", "15", "21", "22", "23", "31", "32",
		"33", "34", "35", "36", "37", "41", "42", "43", "44", "45",
		"46", "50", "51", "52", "53", "54", "61", "62", "63", "64",
		"65", "71", "91", "92",
		// Beijing, Tianjin, Hebei, Shanxi, Inner Mongolia, Liaoning, Jilin, Heilongjiang,
		// Shanghai, Jiangsu, Zhejiang, Anhui, Fujian, Jiangxi, Shandong, Henan, Hubei,
		// Hunan, Guangdong, Guangxi, Hainan, Chongqing, Sichuan, Guizhou, Yunnan, Tibet,
		// Shaanxi, Gansu, Qinghai, Ningxia, Xinjiang, Taiwan, Hong Kong, Macau
	}
	base.categories = []string{"public", "traditional", "lunar"}

	return &CNProvider{BaseProvider: base}
}

// LoadHolidays loads all Chinese holidays for a given year
func (cn *CNProvider) LoadHolidays(year int) map[time.Time]*Holiday {
	holidays := make(map[time.Time]*Holiday)

	// Fixed date holidays

	// New Year's Day - January 1
	newYear := time.Date(year, 1, 1, 0, 0, 0, 0, time.UTC)
	holidays[newYear] = cn.CreateHoliday(
		"元旦",
		newYear,
		"public",
		map[string]string{
			"zh": "元旦",
			"en": "New Year's Day",
		},
	)

	// Women's Day - March 8 (half-day for women)
	womensDay := time.Date(year, 3, 8, 0, 0, 0, 0, time.UTC)
	holidays[womensDay] = cn.CreateHoliday(
		"妇女节",
		womensDay,
		"public",
		map[string]string{
			"zh": "妇女节",
			"en": "Women's Day",
		},
	)

	// Labor Day - May 1
	laborDay := time.Date(year, 5, 1, 0, 0, 0, 0, time.UTC)
	holidays[laborDay] = cn.CreateHoliday(
		"劳动节",
		laborDay,
		"public",
		map[string]string{
			"zh": "劳动节",
			"en": "Labor Day",
		},
	)

	// Youth Day - May 4 (half-day for youth)
	youthDay := time.Date(year, 5, 4, 0, 0, 0, 0, time.UTC)
	holidays[youthDay] = cn.CreateHoliday(
		"青年节",
		youthDay,
		"public",
		map[string]string{
			"zh": "青年节",
			"en": "Youth Day",
		},
	)

	// Children's Day - June 1 (for children under 14)
	childrensDay := time.Date(year, 6, 1, 0, 0, 0, 0, time.UTC)
	holidays[childrensDay] = cn.CreateHoliday(
		"儿童节",
		childrensDay,
		"public",
		map[string]string{
			"zh": "儿童节",
			"en": "Children's Day",
		},
	)

	// Army Day - August 1
	armyDay := time.Date(year, 8, 1, 0, 0, 0, 0, time.UTC)
	holidays[armyDay] = cn.CreateHoliday(
		"建军节",
		armyDay,
		"public",
		map[string]string{
			"zh": "建军节",
			"en": "Army Day",
		},
	)

	// National Day - October 1-3
	nationalDay := time.Date(year, 10, 1, 0, 0, 0, 0, time.UTC)
	holidays[nationalDay] = cn.CreateHoliday(
		"国庆节",
		nationalDay,
		"public",
		map[string]string{
			"zh": "国庆节",
			"en": "National Day",
		},
	)

	nationalDay2 := time.Date(year, 10, 2, 0, 0, 0, 0, time.UTC)
	holidays[nationalDay2] = cn.CreateHoliday(
		"国庆节第二天",
		nationalDay2,
		"public",
		map[string]string{
			"zh": "国庆节第二天",
			"en": "National Day (Day 2)",
		},
	)

	nationalDay3 := time.Date(year, 10, 3, 0, 0, 0, 0, time.UTC)
	holidays[nationalDay3] = cn.CreateHoliday(
		"国庆节第三天",
		nationalDay3,
		"public",
		map[string]string{
			"zh": "国庆节第三天",
			"en": "National Day (Day 3)",
		},
	)

	// Lunar calendar-based holidays
	// These dates change every year based on the lunar calendar

	// Spring Festival (Chinese New Year) - usually late January or February
	springFestivalDates := cn.getSpringFestivalDates(year)
	for i, date := range springFestivalDates {
		var name, enName string
		if i == 0 {
			name = "除夕"
			enName = "Chinese New Year's Eve"
		} else {
			name = "春节第" + cn.getChineseNumber(i) + "天"
			enName = "Chinese New Year (Day " + cn.getEnglishNumber(i) + ")"
		}

		holidays[date] = cn.CreateHoliday(
			name,
			date,
			"lunar",
			map[string]string{
				"zh": name,
				"en": enName,
			},
		)
	}

	// Qingming Festival (Tomb Sweeping Day) - usually early April
	qingmingDate := cn.getQingmingDate(year)
	holidays[qingmingDate] = cn.CreateHoliday(
		"清明节",
		qingmingDate,
		"traditional",
		map[string]string{
			"zh": "清明节",
			"en": "Qingming Festival",
		},
	)

	// Dragon Boat Festival - usually in June
	dragonBoatDate := cn.getDragonBoatDate(year)
	holidays[dragonBoatDate] = cn.CreateHoliday(
		"端午节",
		dragonBoatDate,
		"traditional",
		map[string]string{
			"zh": "端午节",
			"en": "Dragon Boat Festival",
		},
	)

	// Mid-Autumn Festival - usually in September or October
	midAutumnDate := cn.getMidAutumnDate(year)
	holidays[midAutumnDate] = cn.CreateHoliday(
		"中秋节",
		midAutumnDate,
		"traditional",
		map[string]string{
			"zh": "中秋节",
			"en": "Mid-Autumn Festival",
		},
	)

	return holidays
}

// getChineseNumber returns the Chinese character for a number
func (cn *CNProvider) getChineseNumber(num int) string {
	numbers := []string{"一", "二", "三", "四", "五", "六", "七", "八", "九", "十"}
	if num >= 1 && num <= 10 {
		return numbers[num-1]
	}
	return ""
}

// getEnglishNumber returns the ordinal English number as string
func (cn *CNProvider) getEnglishNumber(num int) string {
	switch num {
	case 1:
		return "1st"
	case 2:
		return "2nd"
	case 3:
		return "3rd"
	default:
		return string(rune('0'+num)) + "th"
	}
}

// getSpringFestivalDates returns the dates for Spring Festival (Chinese New Year)
// including Eve and the following days (total 7 days)
func (cn *CNProvider) getSpringFestivalDates(year int) []time.Time {
	var dates []time.Time
	var eve time.Time

	// Lookup table for Spring Festival Eve (Chinese New Year's Eve)
	switch year {
	case 2023:
		eve = time.Date(2023, 1, 21, 0, 0, 0, 0, time.UTC)
	case 2024:
		eve = time.Date(2024, 2, 9, 0, 0, 0, 0, time.UTC)
	case 2025:
		eve = time.Date(2025, 1, 28, 0, 0, 0, 0, time.UTC)
	case 2026:
		eve = time.Date(2026, 2, 16, 0, 0, 0, 0, time.UTC)
	case 2027:
		eve = time.Date(2027, 2, 5, 0, 0, 0, 0, time.UTC)
	case 2028:
		eve = time.Date(2028, 1, 25, 0, 0, 0, 0, time.UTC)
	case 2029:
		eve = time.Date(2029, 2, 12, 0, 0, 0, 0, time.UTC)
	case 2030:
		eve = time.Date(2030, 2, 2, 0, 0, 0, 0, time.UTC)
	default:
		// Fallback to an approximation for years not in the lookup table
		// This is just a rough estimate and should be updated with actual dates
		eve = time.Date(year, 2, 1, 0, 0, 0, 0, time.UTC)
	}

	// Add Eve and the following 6 days (total 7 days celebration)
	dates = append(dates, eve)
	for i := 1; i <= 6; i++ {
		dates = append(dates, eve.AddDate(0, 0, i))
	}

	return dates
}

// getQingmingDate returns the date for Qingming Festival
func (cn *CNProvider) getQingmingDate(year int) time.Time {
	// Lookup table for Qingming Festival
	switch year {
	case 2023:
		return time.Date(2023, 4, 5, 0, 0, 0, 0, time.UTC)
	case 2024:
		return time.Date(2024, 4, 4, 0, 0, 0, 0, time.UTC)
	case 2025:
		return time.Date(2025, 4, 4, 0, 0, 0, 0, time.UTC)
	case 2026:
		return time.Date(2026, 4, 4, 0, 0, 0, 0, time.UTC)
	case 2027:
		return time.Date(2027, 4, 5, 0, 0, 0, 0, time.UTC)
	case 2028:
		return time.Date(2028, 4, 4, 0, 0, 0, 0, time.UTC)
	case 2029:
		return time.Date(2029, 4, 4, 0, 0, 0, 0, time.UTC)
	case 2030:
		return time.Date(2030, 4, 5, 0, 0, 0, 0, time.UTC)
	default:
		// Qingming is usually April 4-6
		return time.Date(year, 4, 4, 0, 0, 0, 0, time.UTC)
	}
}

// getDragonBoatDate returns the date for Dragon Boat Festival
func (cn *CNProvider) getDragonBoatDate(year int) time.Time {
	// Lookup table for Dragon Boat Festival
	switch year {
	case 2023:
		return time.Date(2023, 6, 22, 0, 0, 0, 0, time.UTC)
	case 2024:
		return time.Date(2024, 6, 10, 0, 0, 0, 0, time.UTC)
	case 2025:
		return time.Date(2025, 5, 31, 0, 0, 0, 0, time.UTC)
	case 2026:
		return time.Date(2026, 6, 19, 0, 0, 0, 0, time.UTC)
	case 2027:
		return time.Date(2027, 6, 9, 0, 0, 0, 0, time.UTC)
	case 2028:
		return time.Date(2028, 5, 28, 0, 0, 0, 0, time.UTC)
	case 2029:
		return time.Date(2029, 6, 16, 0, 0, 0, 0, time.UTC)
	case 2030:
		return time.Date(2030, 6, 5, 0, 0, 0, 0, time.UTC)
	default:
		// Fallback to an approximation
		return time.Date(year, 6, 10, 0, 0, 0, 0, time.UTC)
	}
}

// getMidAutumnDate returns the date for Mid-Autumn Festival
func (cn *CNProvider) getMidAutumnDate(year int) time.Time {
	// Lookup table for Mid-Autumn Festival
	switch year {
	case 2023:
		return time.Date(2023, 9, 29, 0, 0, 0, 0, time.UTC)
	case 2024:
		return time.Date(2024, 9, 17, 0, 0, 0, 0, time.UTC)
	case 2025:
		return time.Date(2025, 10, 6, 0, 0, 0, 0, time.UTC)
	case 2026:
		return time.Date(2026, 9, 25, 0, 0, 0, 0, time.UTC)
	case 2027:
		return time.Date(2027, 9, 15, 0, 0, 0, 0, time.UTC)
	case 2028:
		return time.Date(2028, 10, 3, 0, 0, 0, 0, time.UTC)
	case 2029:
		return time.Date(2029, 9, 22, 0, 0, 0, 0, time.UTC)
	case 2030:
		return time.Date(2030, 9, 12, 0, 0, 0, 0, time.UTC)
	default:
		// Fallback to an approximation
		return time.Date(year, 9, 20, 0, 0, 0, 0, time.UTC)
	}
}

// GetSpecialObservances returns non-public observances
func (cn *CNProvider) GetSpecialObservances(year int) map[time.Time]*Holiday {
	holidays := make(map[time.Time]*Holiday)

	// Teachers' Day - September 10
	teachersDay := time.Date(year, 9, 10, 0, 0, 0, 0, time.UTC)
	holidays[teachersDay] = cn.CreateHoliday(
		"教师节",
		teachersDay,
		"traditional",
		map[string]string{
			"zh": "教师节",
			"en": "Teachers' Day",
		},
	)

	// Journalists' Day - November 8
	journalistsDay := time.Date(year, 11, 8, 0, 0, 0, 0, time.UTC)
	holidays[journalistsDay] = cn.CreateHoliday(
		"记者节",
		journalistsDay,
		"traditional",
		map[string]string{
			"zh": "记者节",
			"en": "Journalists' Day",
		},
	)

	// Constitution Day - December 4
	constitutionDay := time.Date(year, 12, 4, 0, 0, 0, 0, time.UTC)
	holidays[constitutionDay] = cn.CreateHoliday(
		"宪法日",
		constitutionDay,
		"traditional",
		map[string]string{
			"zh": "宪法日",
			"en": "Constitution Day",
		},
	)

	return holidays
}

// GetRegionalHolidays returns region-specific holidays
func (cn *CNProvider) GetRegionalHolidays(year int, regions []string) map[time.Time]*Holiday {
	holidays := make(map[time.Time]*Holiday)

	for _, region := range regions {
		switch region {
		case "91": // Hong Kong
			// Hong Kong SAR Establishment Day - July 1
			hkSARDay := time.Date(year, 7, 1, 0, 0, 0, 0, time.UTC)
			holidays[hkSARDay] = cn.CreateHoliday(
				"香港特别行政区成立纪念日",
				hkSARDay,
				"regional",
				map[string]string{
					"zh": "香港特别行政区成立纪念日",
					"en": "Hong Kong SAR Establishment Day",
				},
			)

		case "92": // Macau
			// Macau SAR Establishment Day - December 20
			macauSARDay := time.Date(year, 12, 20, 0, 0, 0, 0, time.UTC)
			holidays[macauSARDay] = cn.CreateHoliday(
				"澳门特别行政区成立纪念日",
				macauSARDay,
				"regional",
				map[string]string{
					"zh": "澳门特别行政区成立纪念日",
					"en": "Macau SAR Establishment Day",
				},
			)
		}
	}

	return holidays
}
