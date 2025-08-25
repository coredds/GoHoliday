package countries

import (
	"time"
)

// NZProvider implements holiday calculations for New Zealand
type NZProvider struct {
	*BaseProvider
}

// NewNZProvider creates a new New Zealand holiday provider
func NewNZProvider() *NZProvider {
	base := NewBaseProvider("NZ")
	base.subdivisions = []string{
		"AUK", // Auckland
		"BOP", // Bay of Plenty
		"CAN", // Canterbury
		"GIS", // Gisborne
		"HKB", // Hawke's Bay
		"MWT", // Manawatu-Wanganui
		"MBH", // Marlborough
		"NSN", // Nelson
		"NTL", // Northland
		"OTA", // Otago
		"STL", // Southland
		"TKI", // Taranaki
		"TAS", // Tasman
		"WKO", // Waikato
		"WGN", // Wellington
		"WTC", // West Coast
		"CIT", // Chatham Islands Territory
	}
	base.categories = []string{"public", "regional"}
	
	return &NZProvider{BaseProvider: base}
}

// LoadHolidays loads all New Zealand holidays for a given year
func (nz *NZProvider) LoadHolidays(year int) map[time.Time]*Holiday {
	holidays := make(map[time.Time]*Holiday)
	
	// Fixed date holidays
	holidays[time.Date(year, 1, 1, 0, 0, 0, 0, time.UTC)] = nz.CreateHoliday(
		"New Year's Day",
		time.Date(year, 1, 1, 0, 0, 0, 0, time.UTC),
		"public",
		map[string]string{
			"en": "New Year's Day",
			"mi": "Te Rā Tau Hou", // Māori
		},
	)
	
	holidays[time.Date(year, 1, 2, 0, 0, 0, 0, time.UTC)] = nz.CreateHoliday(
		"Day after New Year's Day",
		time.Date(year, 1, 2, 0, 0, 0, 0, time.UTC),
		"public",
		map[string]string{
			"en": "Day after New Year's Day",
			"mi": "Te Rā i muri i te Rā Tau Hou",
		},
	)
	
	holidays[time.Date(year, 2, 6, 0, 0, 0, 0, time.UTC)] = nz.CreateHoliday(
		"Waitangi Day",
		time.Date(year, 2, 6, 0, 0, 0, 0, time.UTC),
		"public",
		map[string]string{
			"en": "Waitangi Day",
			"mi": "Te Rā o Waitangi",
		},
	)
	
	holidays[time.Date(year, 4, 25, 0, 0, 0, 0, time.UTC)] = nz.CreateHoliday(
		"ANZAC Day",
		time.Date(year, 4, 25, 0, 0, 0, 0, time.UTC),
		"public",
		map[string]string{
			"en": "ANZAC Day",
			"mi": "Te Rā ANZAC",
		},
	)
	
	holidays[time.Date(year, 12, 25, 0, 0, 0, 0, time.UTC)] = nz.CreateHoliday(
		"Christmas Day",
		time.Date(year, 12, 25, 0, 0, 0, 0, time.UTC),
		"public",
		map[string]string{
			"en": "Christmas Day",
			"mi": "Te Rā Kirihimete",
		},
	)
	
	holidays[time.Date(year, 12, 26, 0, 0, 0, 0, time.UTC)] = nz.CreateHoliday(
		"Boxing Day",
		time.Date(year, 12, 26, 0, 0, 0, 0, time.UTC),
		"public",
		map[string]string{
			"en": "Boxing Day",
			"mi": "Te Rā Pākete",
		},
	)
	
	// Easter-based holidays
	easter := EasterSunday(year)
	
	// Good Friday
	goodFriday := easter.AddDate(0, 0, -2)
	holidays[goodFriday] = nz.CreateHoliday(
		"Good Friday",
		goodFriday,
		"public",
		map[string]string{
			"en": "Good Friday",
			"mi": "Paraire Pai",
		},
	)
	
	// Easter Monday
	easterMonday := easter.AddDate(0, 0, 1)
	holidays[easterMonday] = nz.CreateHoliday(
		"Easter Monday",
		easterMonday,
		"public",
		map[string]string{
			"en": "Easter Monday",
			"mi": "Mane Aranga",
		},
	)
	
	// Variable date holidays
	
	// Queen's Birthday - first Monday in June
	queensBirthday := NthWeekdayOfMonth(year, 6, time.Monday, 1)
	holidays[queensBirthday] = nz.CreateHoliday(
		"Queen's Birthday",
		queensBirthday,
		"public",
		map[string]string{
			"en": "Queen's Birthday",
			"mi": "Te Rā Whānau o te Kuini",
		},
	)
	
	// Labour Day - fourth Monday in October
	labourDay := NthWeekdayOfMonth(year, 10, time.Monday, 4)
	holidays[labourDay] = nz.CreateHoliday(
		"Labour Day",
		labourDay,
		"public",
		map[string]string{
			"en": "Labour Day",
			"mi": "Te Rā Whakanui i nga Kaimahi",
		},
	)
	
	// Matariki - Māori New Year (varies each year)
	matariki := nz.getMatarikiDate(year)
	if !matariki.IsZero() {
		holidays[matariki] = nz.CreateHoliday(
			"Matariki",
			matariki,
			"public",
			map[string]string{
				"en": "Matariki",
				"mi": "Matariki",
			},
		)
	}
	
	return holidays
}

// GetRegionalHolidays returns region-specific holidays (provincial anniversaries)
func (nz *NZProvider) GetRegionalHolidays(year int, regions []string) map[time.Time]*Holiday {
	holidays := make(map[time.Time]*Holiday)
	
	for _, region := range regions {
		switch region {
		case "AUK": // Auckland Anniversary
			aucklandDay := nz.getAucklandAnniversary(year)
			holidays[aucklandDay] = nz.CreateHoliday(
				"Auckland Anniversary Day",
				aucklandDay,
				"regional",
				map[string]string{
					"en": "Auckland Anniversary Day",
					"mi": "Te Rā Whakamaumahara o Tāmaki-makau-rau",
				},
			)
			
		case "WGN": // Wellington Anniversary
			wellingtonDay := nz.getWellingtonAnniversary(year)
			holidays[wellingtonDay] = nz.CreateHoliday(
				"Wellington Anniversary Day",
				wellingtonDay,
				"regional",
				map[string]string{
					"en": "Wellington Anniversary Day",
					"mi": "Te Rā Whakamaumahara o Te Whanganui-a-Tara",
				},
			)
			
		case "CAN": // Canterbury Anniversary
			canterburyDay := nz.getCanterburyAnniversary(year)
			holidays[canterburyDay] = nz.CreateHoliday(
				"Canterbury Anniversary Day",
				canterburyDay,
				"regional",
				map[string]string{
					"en": "Canterbury Anniversary Day",
					"mi": "Te Rā Whakamaumahara o Waitaha",
				},
			)
			
		case "OTA": // Otago Anniversary
			otagoDay := nz.getOtagoAnniversary(year)
			holidays[otagoDay] = nz.CreateHoliday(
				"Otago Anniversary Day",
				otagoDay,
				"regional",
				map[string]string{
					"en": "Otago Anniversary Day",
					"mi": "Te Rā Whakamaumahara o Ōtākou",
				},
			)
			
		case "STL": // Southland Anniversary
			southlandDay := nz.getSouthlandAnniversary(year)
			holidays[southlandDay] = nz.CreateHoliday(
				"Southland Anniversary Day",
				southlandDay,
				"regional",
				map[string]string{
					"en": "Southland Anniversary Day",
					"mi": "Te Rā Whakamaumahara o Murihiku",
				},
			)
			
		case "WKO": // Waikato Anniversary
			waikatoDay := nz.getWaikatoAnniversary(year)
			holidays[waikatoDay] = nz.CreateHoliday(
				"Waikato Anniversary Day",
				waikatoDay,
				"regional",
				map[string]string{
					"en": "Waikato Anniversary Day",
					"mi": "Te Rā Whakamaumahara o Waikato",
				},
			)
			
		case "HKB": // Hawke's Bay Anniversary
			hawkesBayDay := nz.getHawkesBayAnniversary(year)
			holidays[hawkesBayDay] = nz.CreateHoliday(
				"Hawke's Bay Anniversary Day",
				hawkesBayDay,
				"regional",
				map[string]string{
					"en": "Hawke's Bay Anniversary Day",
					"mi": "Te Rā Whakamaumahara o Te Matau-a-Māui",
				},
			)
			
		case "TKI": // Taranaki Anniversary
			taranakiDay := nz.getTaranakiAnniversary(year)
			holidays[taranakiDay] = nz.CreateHoliday(
				"Taranaki Anniversary Day",
				taranakiDay,
				"regional",
				map[string]string{
					"en": "Taranaki Anniversary Day",
					"mi": "Te Rā Whakamaumahara o Taranaki",
				},
			)
			
		case "NSN", "MBH", "TAS": // Nelson, Marlborough, Tasman Anniversary
			nelsonDay := nz.getNelsonAnniversary(year)
			holidays[nelsonDay] = nz.CreateHoliday(
				"Nelson Anniversary Day",
				nelsonDay,
				"regional",
				map[string]string{
					"en": "Nelson Anniversary Day",
					"mi": "Te Rā Whakamaumahara o Whakatū",
				},
			)
			
		case "WTC": // West Coast Anniversary
			westCoastDay := nz.getWestCoastAnniversary(year)
			holidays[westCoastDay] = nz.CreateHoliday(
				"West Coast Anniversary Day",
				westCoastDay,
				"regional",
				map[string]string{
					"en": "West Coast Anniversary Day",
					"mi": "Te Rā Whakamaumahara o Te Tai Poutini",
				},
			)
			
		case "CIT": // Chatham Islands Anniversary
			chathamDay := nz.getChathamAnniversary(year)
			holidays[chathamDay] = nz.CreateHoliday(
				"Chatham Islands Anniversary Day",
				chathamDay,
				"regional",
				map[string]string{
					"en": "Chatham Islands Anniversary Day",
					"mi": "Te Rā Whakamaumahara o Rēkohu",
				},
			)
		}
	}
	
	return holidays
}

// Anniversary calculation methods

// getAucklandAnniversary - Monday closest to January 29
func (nz *NZProvider) getAucklandAnniversary(year int) time.Time {
	jan29 := time.Date(year, 1, 29, 0, 0, 0, 0, time.UTC)
	return nz.getClosestMonday(jan29)
}

// getWellingtonAnniversary - Monday closest to January 22
func (nz *NZProvider) getWellingtonAnniversary(year int) time.Time {
	jan22 := time.Date(year, 1, 22, 0, 0, 0, 0, time.UTC)
	return nz.getClosestMonday(jan22)
}

// getCanterburyAnniversary - Friday after the first Tuesday in November (Canterbury Show Day)
func (nz *NZProvider) getCanterburyAnniversary(year int) time.Time {
	firstTuesday := NthWeekdayOfMonth(year, 11, time.Tuesday, 1)
	return firstTuesday.AddDate(0, 0, 3) // Friday after
}

// getOtagoAnniversary - Monday closest to March 23
func (nz *NZProvider) getOtagoAnniversary(year int) time.Time {
	mar23 := time.Date(year, 3, 23, 0, 0, 0, 0, time.UTC)
	return nz.getClosestMonday(mar23)
}

// getSouthlandAnniversary - Easter Tuesday (day after Easter Monday)
func (nz *NZProvider) getSouthlandAnniversary(year int) time.Time {
	easter := EasterSunday(year)
	return easter.AddDate(0, 0, 2) // Tuesday after Easter
}

// getWaikatoAnniversary - Monday closest to December 1
func (nz *NZProvider) getWaikatoAnniversary(year int) time.Time {
	dec1 := time.Date(year, 12, 1, 0, 0, 0, 0, time.UTC)
	return nz.getClosestMonday(dec1)
}

// getHawkesBayAnniversary - Friday before Labour Day weekend
func (nz *NZProvider) getHawkesBayAnniversary(year int) time.Time {
	labourDay := NthWeekdayOfMonth(year, 10, time.Monday, 4)
	return labourDay.AddDate(0, 0, -3) // Friday before
}

// getTaranakiAnniversary - Second Monday in March
func (nz *NZProvider) getTaranakiAnniversary(year int) time.Time {
	return NthWeekdayOfMonth(year, 3, time.Monday, 2)
}

// getNelsonAnniversary - Monday closest to February 1
func (nz *NZProvider) getNelsonAnniversary(year int) time.Time {
	feb1 := time.Date(year, 2, 1, 0, 0, 0, 0, time.UTC)
	return nz.getClosestMonday(feb1)
}

// getWestCoastAnniversary - Monday closest to December 1
func (nz *NZProvider) getWestCoastAnniversary(year int) time.Time {
	dec1 := time.Date(year, 12, 1, 0, 0, 0, 0, time.UTC)
	return nz.getClosestMonday(dec1)
}

// getChathamAnniversary - Monday closest to November 30
func (nz *NZProvider) getChathamAnniversary(year int) time.Time {
	nov30 := time.Date(year, 11, 30, 0, 0, 0, 0, time.UTC)
	return nz.getClosestMonday(nov30)
}

// getClosestMonday returns the Monday closest to the given date
func (nz *NZProvider) getClosestMonday(date time.Time) time.Time {
	weekday := date.Weekday()
	switch weekday {
	case time.Monday:
		return date
	case time.Tuesday:
		return date.AddDate(0, 0, -1) // Previous Monday
	case time.Wednesday:
		return date.AddDate(0, 0, -2) // Previous Monday
	case time.Thursday:
		return date.AddDate(0, 0, -3) // Previous Monday
	case time.Friday:
		return date.AddDate(0, 0, 3) // Next Monday
	case time.Saturday:
		return date.AddDate(0, 0, 2) // Next Monday
	case time.Sunday:
		return date.AddDate(0, 0, 1) // Next Monday
	}
	return date
}

// getMatarikiDate calculates Matariki date (simplified - actual calculation is astronomical)
func (nz *NZProvider) getMatarikiDate(year int) time.Time {
	// Matariki dates are determined astronomically. Here are known dates:
	matarikiDates := map[int]time.Time{
		2022: time.Date(2022, 6, 24, 0, 0, 0, 0, time.UTC),
		2023: time.Date(2023, 7, 14, 0, 0, 0, 0, time.UTC),
		2024: time.Date(2024, 6, 28, 0, 0, 0, 0, time.UTC),
		2025: time.Date(2025, 6, 20, 0, 0, 0, 0, time.UTC),
		2026: time.Date(2026, 7, 10, 0, 0, 0, 0, time.UTC),
		2027: time.Date(2027, 6, 25, 0, 0, 0, 0, time.UTC),
		2028: time.Date(2028, 7, 14, 0, 0, 0, 0, time.UTC),
		2029: time.Date(2029, 7, 6, 0, 0, 0, 0, time.UTC),
		2030: time.Date(2030, 6, 21, 0, 0, 0, 0, time.UTC),
	}
	
	if date, exists := matarikiDates[year]; exists {
		return date
	}
	
	// For years not in our table, return zero time (no Matariki holiday)
	return time.Time{}
}

// GetSeasons returns New Zealand seasons (Southern Hemisphere)
func (nz *NZProvider) GetSeasons(year int) map[string][]time.Time {
	return map[string][]time.Time{
		"Summer": {
			time.Date(year-1, 12, 1, 0, 0, 0, 0, time.UTC), // Dec 1 (previous year)
			time.Date(year, 2, 28, 0, 0, 0, 0, time.UTC),   // Feb 28/29
		},
		"Autumn": {
			time.Date(year, 3, 1, 0, 0, 0, 0, time.UTC),  // Mar 1
			time.Date(year, 5, 31, 0, 0, 0, 0, time.UTC), // May 31
		},
		"Winter": {
			time.Date(year, 6, 1, 0, 0, 0, 0, time.UTC),  // Jun 1
			time.Date(year, 8, 31, 0, 0, 0, 0, time.UTC), // Aug 31
		},
		"Spring": {
			time.Date(year, 9, 1, 0, 0, 0, 0, time.UTC),   // Sep 1
			time.Date(year, 11, 30, 0, 0, 0, 0, time.UTC), // Nov 30
		},
	}
}
