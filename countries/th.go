package countries

import (
	"time"
)

// THProvider implements holiday calculations for Thailand
type THProvider struct {
	*BaseProvider
}

// NewTHProvider creates a new Thailand holiday provider
func NewTHProvider() *THProvider {
	base := NewBaseProvider("TH")

	// Thailand has 77 provinces
	base.subdivisions = []string{
		"TH-10", "TH-11", "TH-12", "TH-13", "TH-14", "TH-15", "TH-16", "TH-17", "TH-18", "TH-19",
		"TH-20", "TH-21", "TH-22", "TH-23", "TH-24", "TH-25", "TH-26", "TH-27", "TH-30", "TH-31",
		"TH-32", "TH-33", "TH-34", "TH-35", "TH-36", "TH-37", "TH-38", "TH-39", "TH-40", "TH-41",
		"TH-42", "TH-43", "TH-44", "TH-45", "TH-46", "TH-47", "TH-48", "TH-49", "TH-50", "TH-51",
		"TH-52", "TH-53", "TH-54", "TH-55", "TH-56", "TH-57", "TH-58", "TH-60", "TH-61", "TH-62",
		"TH-63", "TH-64", "TH-65", "TH-66", "TH-67", "TH-70", "TH-71", "TH-72", "TH-73", "TH-74",
		"TH-75", "TH-76", "TH-77", "TH-80", "TH-81", "TH-82", "TH-83", "TH-84", "TH-85", "TH-86",
		"TH-90", "TH-91", "TH-92", "TH-93", "TH-94", "TH-95", "TH-96",
	}

	base.categories = []string{"national", "religious", "royal", "buddhist", "cultural"}

	return &THProvider{BaseProvider: base}
}

// LoadHolidays loads all holidays for Thailand for the given year
func (p *THProvider) LoadHolidays(year int) map[time.Time]*Holiday {
	holidays := make(map[time.Time]*Holiday)

	// Fixed national holidays
	holidays[time.Date(year, 1, 1, 0, 0, 0, 0, time.UTC)] = p.CreateHoliday(
		"วันขึ้นปีใหม่", time.Date(year, 1, 1, 0, 0, 0, 0, time.UTC), "national",
		map[string]string{
			"th": "วันขึ้นปีใหม่",
			"en": "New Year's Day",
		},
	)

	// Magha Puja Day (full moon of 3rd lunar month) - February/March
	maghaPuja := p.calculateMaghaPuja(year)
	holidays[maghaPuja] = p.CreateHoliday(
		"วันมาฆบูชา", maghaPuja, "buddhist",
		map[string]string{
			"th": "วันมาฆบูชา",
			"en": "Magha Puja Day",
		},
	)

	// Chakri Day - April 6
	holidays[time.Date(year, 4, 6, 0, 0, 0, 0, time.UTC)] = p.CreateHoliday(
		"วันจักรี", time.Date(year, 4, 6, 0, 0, 0, 0, time.UTC), "royal",
		map[string]string{
			"th": "วันจักรี",
			"en": "Chakri Day",
		},
	)

	// Songkran Festival - April 13-15
	holidays[time.Date(year, 4, 13, 0, 0, 0, 0, time.UTC)] = p.CreateHoliday(
		"วันสงกรานต์", time.Date(year, 4, 13, 0, 0, 0, 0, time.UTC), "cultural",
		map[string]string{
			"th": "วันสงกรานต์",
			"en": "Songkran Festival",
		},
	)

	holidays[time.Date(year, 4, 14, 0, 0, 0, 0, time.UTC)] = p.CreateHoliday(
		"วันสงกรานต์ (วันที่ 2)", time.Date(year, 4, 14, 0, 0, 0, 0, time.UTC), "cultural",
		map[string]string{
			"th": "วันสงกรานต์ (วันที่ 2)",
			"en": "Songkran Festival (Day 2)",
		},
	)

	holidays[time.Date(year, 4, 15, 0, 0, 0, 0, time.UTC)] = p.CreateHoliday(
		"วันสงกรานต์ (วันที่ 3)", time.Date(year, 4, 15, 0, 0, 0, 0, time.UTC), "cultural",
		map[string]string{
			"th": "วันสงกรานต์ (วันที่ 3)",
			"en": "Songkran Festival (Day 3)",
		},
	)

	// National Labour Day - May 1
	holidays[time.Date(year, 5, 1, 0, 0, 0, 0, time.UTC)] = p.CreateHoliday(
		"วันแรงงานแห่งชาติ", time.Date(year, 5, 1, 0, 0, 0, 0, time.UTC), "national",
		map[string]string{
			"th": "วันแรงงานแห่งชาติ",
			"en": "National Labour Day",
		},
	)

	// Coronation Day - May 4
	holidays[time.Date(year, 5, 4, 0, 0, 0, 0, time.UTC)] = p.CreateHoliday(
		"วันฉัตรมงคล", time.Date(year, 5, 4, 0, 0, 0, 0, time.UTC), "royal",
		map[string]string{
			"th": "วันฉัตรมงคล",
			"en": "Coronation Day",
		},
	)

	// Royal Ploughing Ceremony (variable date in May)
	royalPloughing := p.calculateRoyalPloughing(year)
	holidays[royalPloughing] = p.CreateHoliday(
		"วันพืชมงคล", royalPloughing, "royal",
		map[string]string{
			"th": "วันพืชมงคล",
			"en": "Royal Ploughing Ceremony",
		},
	)

	// Visakha Puja Day (full moon of 6th lunar month) - May/June
	visakhaPuja := p.calculateVisakhaPuja(year)
	holidays[visakhaPuja] = p.CreateHoliday(
		"วันวิสาขบูชา", visakhaPuja, "buddhist",
		map[string]string{
			"th": "วันวิสาขบูชา",
			"en": "Visakha Puja Day",
		},
	)

	// Queen Suthida's Birthday - June 3
	holidays[time.Date(year, 6, 3, 0, 0, 0, 0, time.UTC)] = p.CreateHoliday(
		"วันเฉลิมพระชนมพรรษาสมเด็จพระนางเจ้าสุทิดา", time.Date(year, 6, 3, 0, 0, 0, 0, time.UTC), "royal",
		map[string]string{
			"th": "วันเฉลิมพระชนมพรรษาสมเด็จพระนางเจ้าสุทิดา",
			"en": "Queen Suthida's Birthday",
		},
	)

	// Asalha Puja Day (full moon of 8th lunar month) - July/August
	asalhaPuja := p.calculateAsalhaPuja(year)
	holidays[asalhaPuja] = p.CreateHoliday(
		"วันอาสาฬหบูชา", asalhaPuja, "buddhist",
		map[string]string{
			"th": "วันอาสาฬหบูชา",
			"en": "Asalha Puja Day",
		},
	)

	// Khao Phansa (Buddhist Lent begins) - day after Asalha Puja
	khaoPhansal := asalhaPuja.AddDate(0, 0, 1)
	holidays[khaoPhansal] = p.CreateHoliday(
		"วันเข้าพรรษา", khaoPhansal, "buddhist",
		map[string]string{
			"th": "วันเข้าพรรษา",
			"en": "Khao Phansa",
		},
	)

	// HM King Maha Vajiralongkorn's Birthday - July 28
	holidays[time.Date(year, 7, 28, 0, 0, 0, 0, time.UTC)] = p.CreateHoliday(
		"วันเฉลิมพระชนมพรรษาพระบาทสมเด็จพระปรเมนทรรามาธิบดีศรีสินทรมหาวชิราลงกรณ พระวชิรเกล้าเจ้าอยู่หัว", time.Date(year, 7, 28, 0, 0, 0, 0, time.UTC), "royal",
		map[string]string{
			"th": "วันเฉลิมพระชนมพรรษาพระบาทสมเด็จพระปรเมนทรรามาธิบดีศรีสินทรมหาวชิราลงกรณ พระวชิรเกล้าเจ้าอยู่หัว",
			"en": "HM King Maha Vajiralongkorn's Birthday",
		},
	)

	// HM Queen Sirikit The Queen Mother's Birthday - August 12
	holidays[time.Date(year, 8, 12, 0, 0, 0, 0, time.UTC)] = p.CreateHoliday(
		"วันเฉลิมพระชนมพรรษาสมเด็จพระนางเจ้าสิริกิติ์ พระบรมราชินีนาถ พระบรมราชชนนีพันปีหลวง", time.Date(year, 8, 12, 0, 0, 0, 0, time.UTC), "royal",
		map[string]string{
			"th": "วันเฉลิมพระชนมพรรษาสมเด็จพระนางเจ้าสิริกิติ์ พระบรมราชินีนาถ พระบรมราชชนนีพันปีหลวง",
			"en": "HM Queen Sirikit The Queen Mother's Birthday",
		},
	)

	// HM King Bhumibol Adulyadej Memorial Day - October 13
	holidays[time.Date(year, 10, 13, 0, 0, 0, 0, time.UTC)] = p.CreateHoliday(
		"วันคล้ายวันสวรรคตพระบาทสมเด็จพระปรมินทรมหาภูมิพลอดุลยเดช บรมนาถบพิตร", time.Date(year, 10, 13, 0, 0, 0, 0, time.UTC), "royal",
		map[string]string{
			"th": "วันคล้ายวันสวรรคตพระบาทสมเด็จพระปรมินทรมหาภูมิพลอดุลยเดช บรมนาถบพิตร",
			"en": "HM King Bhumibol Adulyadej Memorial Day",
		},
	)

	// Chulalongkorn Day - October 23
	holidays[time.Date(year, 10, 23, 0, 0, 0, 0, time.UTC)] = p.CreateHoliday(
		"วันปิยมหาราช", time.Date(year, 10, 23, 0, 0, 0, 0, time.UTC), "royal",
		map[string]string{
			"th": "วันปิยมหาราช",
			"en": "Chulalongkorn Day",
		},
	)

	// HM King Bhumibol Adulyadej's Birthday - December 5
	holidays[time.Date(year, 12, 5, 0, 0, 0, 0, time.UTC)] = p.CreateHoliday(
		"วันเฉลิมพระชนมพรรษาพระบาทสมเด็จพระปรมินทรมหาภูมิพลอดุลยเดช บรมนาถบพิตร", time.Date(year, 12, 5, 0, 0, 0, 0, time.UTC), "royal",
		map[string]string{
			"th": "วันเฉลิมพระชนมพรรษาพระบาทสมเด็จพระปรมินทรมหาภูมิพลอดุลยเดช บรมนาถบพิตร",
			"en": "HM King Bhumibol Adulyadej's Birthday",
		},
	)

	// Constitution Day - December 10
	holidays[time.Date(year, 12, 10, 0, 0, 0, 0, time.UTC)] = p.CreateHoliday(
		"วันรัฐธรรมนูญ", time.Date(year, 12, 10, 0, 0, 0, 0, time.UTC), "national",
		map[string]string{
			"th": "วันรัฐธรรมนูญ",
			"en": "Constitution Day",
		},
	)

	// New Year's Eve - December 31
	holidays[time.Date(year, 12, 31, 0, 0, 0, 0, time.UTC)] = p.CreateHoliday(
		"วันสิ้นปี", time.Date(year, 12, 31, 0, 0, 0, 0, time.UTC), "national",
		map[string]string{
			"th": "วันสิ้นปี",
			"en": "New Year's Eve",
		},
	)

	return holidays
}

// calculateMaghaPuja calculates Magha Puja Day (full moon of 3rd lunar month)
// This is an approximation - in practice, Thai authorities announce the exact date
func (p *THProvider) calculateMaghaPuja(year int) time.Time {
	// Magha Puja typically falls in February or early March
	// Using approximation based on historical patterns
	switch year {
	case 2024:
		return time.Date(2024, 2, 24, 0, 0, 0, 0, time.UTC)
	case 2025:
		return time.Date(2025, 2, 12, 0, 0, 0, 0, time.UTC)
	case 2026:
		return time.Date(2026, 3, 4, 0, 0, 0, 0, time.UTC)
	case 2027:
		return time.Date(2027, 2, 21, 0, 0, 0, 0, time.UTC)
	default:
		// Fallback calculation - approximately 57-59 days after Western New Year
		return time.Date(year, 2, 24, 0, 0, 0, 0, time.UTC)
	}
}

// calculateVisakhaPuja calculates Visakha Puja Day (full moon of 6th lunar month)
func (p *THProvider) calculateVisakhaPuja(year int) time.Time {
	// Visakha Puja typically falls in May or early June
	switch year {
	case 2024:
		return time.Date(2024, 5, 22, 0, 0, 0, 0, time.UTC)
	case 2025:
		return time.Date(2025, 5, 12, 0, 0, 0, 0, time.UTC)
	case 2026:
		return time.Date(2026, 5, 31, 0, 0, 0, 0, time.UTC)
	case 2027:
		return time.Date(2027, 5, 21, 0, 0, 0, 0, time.UTC)
	default:
		// Fallback calculation
		return time.Date(year, 5, 22, 0, 0, 0, 0, time.UTC)
	}
}

// calculateAsalhaPuja calculates Asalha Puja Day (full moon of 8th lunar month)
func (p *THProvider) calculateAsalhaPuja(year int) time.Time {
	// Asalha Puja typically falls in July or early August
	switch year {
	case 2024:
		return time.Date(2024, 7, 21, 0, 0, 0, 0, time.UTC)
	case 2025:
		return time.Date(2025, 7, 11, 0, 0, 0, 0, time.UTC)
	case 2026:
		return time.Date(2026, 7, 30, 0, 0, 0, 0, time.UTC)
	case 2027:
		return time.Date(2027, 7, 19, 0, 0, 0, 0, time.UTC)
	default:
		// Fallback calculation
		return time.Date(year, 7, 21, 0, 0, 0, 0, time.UTC)
	}
}

// calculateRoyalPloughing calculates the Royal Ploughing Ceremony date
// This varies each year and is officially announced
func (p *THProvider) calculateRoyalPloughing(year int) time.Time {
	// Royal Ploughing Ceremony typically falls in May
	switch year {
	case 2024:
		return time.Date(2024, 5, 9, 0, 0, 0, 0, time.UTC)
	case 2025:
		return time.Date(2025, 5, 8, 0, 0, 0, 0, time.UTC)
	case 2026:
		return time.Date(2026, 5, 14, 0, 0, 0, 0, time.UTC)
	case 2027:
		return time.Date(2027, 5, 13, 0, 0, 0, 0, time.UTC)
	default:
		// Fallback to second Thursday of May
		return p.getNthWeekdayOfMonth(year, 5, time.Thursday, 2)
	}
}

// getNthWeekdayOfMonth returns the nth occurrence of a weekday in a given month
func (p *THProvider) getNthWeekdayOfMonth(year int, month time.Month, weekday time.Weekday, n int) time.Time {
	firstDay := time.Date(year, month, 1, 0, 0, 0, 0, time.UTC)

	// Find the first occurrence of the target weekday
	daysToAdd := int((weekday - firstDay.Weekday() + 7) % 7)
	firstOccurrence := firstDay.AddDate(0, 0, daysToAdd)

	// Add weeks to get the nth occurrence
	return firstOccurrence.AddDate(0, 0, (n-1)*7)
}
