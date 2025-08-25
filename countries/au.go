package countries

import (
	"time"
)

// AUProvider implements holiday calculations for Australia
type AUProvider struct {
	*BaseProvider
}

// NewAUProvider creates a new Australia holiday provider
func NewAUProvider() *AUProvider {
	base := NewBaseProvider("AU")
	base.subdivisions = []string{
		"NSW", // New South Wales
		"VIC", // Victoria
		"QLD", // Queensland
		"SA",  // South Australia
		"WA",  // Western Australia
		"TAS", // Tasmania
		"NT",  // Northern Territory
		"ACT", // Australian Capital Territory
	}
	base.categories = []string{"public", "bank", "government"}
	
	return &AUProvider{BaseProvider: base}
}

// LoadHolidays loads all Australian holidays for a given year
func (au *AUProvider) LoadHolidays(year int) map[time.Time]*Holiday {
	holidays := make(map[time.Time]*Holiday)
	
	// Fixed date holidays
	holidays[time.Date(year, 1, 1, 0, 0, 0, 0, time.UTC)] = au.CreateHoliday(
		"New Year's Day",
		time.Date(year, 1, 1, 0, 0, 0, 0, time.UTC),
		"public",
		map[string]string{
			"en": "New Year's Day",
		},
	)
	
	holidays[time.Date(year, 1, 26, 0, 0, 0, 0, time.UTC)] = au.CreateHoliday(
		"Australia Day",
		time.Date(year, 1, 26, 0, 0, 0, 0, time.UTC),
		"public",
		map[string]string{
			"en": "Australia Day",
		},
	)
	
	holidays[time.Date(year, 4, 25, 0, 0, 0, 0, time.UTC)] = au.CreateHoliday(
		"ANZAC Day",
		time.Date(year, 4, 25, 0, 0, 0, 0, time.UTC),
		"public",
		map[string]string{
			"en": "ANZAC Day",
		},
	)
	
	holidays[time.Date(year, 12, 25, 0, 0, 0, 0, time.UTC)] = au.CreateHoliday(
		"Christmas Day",
		time.Date(year, 12, 25, 0, 0, 0, 0, time.UTC),
		"public",
		map[string]string{
			"en": "Christmas Day",
		},
	)
	
	holidays[time.Date(year, 12, 26, 0, 0, 0, 0, time.UTC)] = au.CreateHoliday(
		"Boxing Day",
		time.Date(year, 12, 26, 0, 0, 0, 0, time.UTC),
		"public",
		map[string]string{
			"en": "Boxing Day",
		},
	)
	
	// Easter-based holidays
	easter := EasterSunday(year)
	
	// Good Friday
	goodFriday := easter.AddDate(0, 0, -2)
	holidays[goodFriday] = au.CreateHoliday(
		"Good Friday",
		goodFriday,
		"public",
		map[string]string{
			"en": "Good Friday",
		},
	)
	
	// Easter Saturday
	easterSaturday := easter.AddDate(0, 0, -1)
	holidays[easterSaturday] = au.CreateHoliday(
		"Easter Saturday",
		easterSaturday,
		"public",
		map[string]string{
			"en": "Easter Saturday",
		},
	)
	
	// Easter Monday
	easterMonday := easter.AddDate(0, 0, 1)
	holidays[easterMonday] = au.CreateHoliday(
		"Easter Monday",
		easterMonday,
		"public",
		map[string]string{
			"en": "Easter Monday",
		},
	)
	
	// Variable date holidays (most states)
	
	// Queen's Birthday - 2nd Monday in June (most states)
	queensBirthday := NthWeekdayOfMonth(year, 6, time.Monday, 2)
	holidays[queensBirthday] = au.CreateHoliday(
		"Queen's Birthday",
		queensBirthday,
		"public",
		map[string]string{
			"en": "Queen's Birthday",
		},
	)
	
	// Labour Day - 1st Monday in October (most states)
	labourDay := NthWeekdayOfMonth(year, 10, time.Monday, 1)
	holidays[labourDay] = au.CreateHoliday(
		"Labour Day",
		labourDay,
		"public",
		map[string]string{
			"en": "Labour Day",
		},
	)
	
	return holidays
}

// GetStateHolidays returns state-specific holidays
func (au *AUProvider) GetStateHolidays(year int, states []string) map[time.Time]*Holiday {
	holidays := make(map[time.Time]*Holiday)
	
	for _, state := range states {
		switch state {
		case "VIC": // Victoria
			// Melbourne Cup Day - 1st Tuesday in November
			melbourneCup := au.getMelbourneCupDay(year)
			holidays[melbourneCup] = au.CreateHoliday(
				"Melbourne Cup Day",
				melbourneCup,
				"public",
				map[string]string{
					"en": "Melbourne Cup Day",
				},
			)
			
			// Labour Day - 2nd Monday in March (different from other states)
			labourDayVic := NthWeekdayOfMonth(year, 3, time.Monday, 2)
			holidays[labourDayVic] = au.CreateHoliday(
				"Labour Day",
				labourDayVic,
				"public",
				map[string]string{
					"en": "Labour Day",
				},
			)
			
		case "WA": // Western Australia
			// Western Australia Day - 1st Monday in June
			waDay := NthWeekdayOfMonth(year, 6, time.Monday, 1)
			holidays[waDay] = au.CreateHoliday(
				"Western Australia Day",
				waDay,
				"public",
				map[string]string{
					"en": "Western Australia Day",
				},
			)
			
			// Labour Day - 1st Monday in March
			labourDayWA := NthWeekdayOfMonth(year, 3, time.Monday, 1)
			holidays[labourDayWA] = au.CreateHoliday(
				"Labour Day",
				labourDayWA,
				"public",
				map[string]string{
					"en": "Labour Day",
				},
			)
			
		case "QLD": // Queensland
			// Labour Day - 1st Monday in May
			labourDayQLD := NthWeekdayOfMonth(year, 5, time.Monday, 1)
			holidays[labourDayQLD] = au.CreateHoliday(
				"Labour Day",
				labourDayQLD,
				"public",
				map[string]string{
					"en": "Labour Day",
				},
			)
			
			// Queen's Birthday - 1st Monday in October (different from other states)
			queensBirthdayQLD := NthWeekdayOfMonth(year, 10, time.Monday, 1)
			holidays[queensBirthdayQLD] = au.CreateHoliday(
				"Queen's Birthday",
				queensBirthdayQLD,
				"public",
				map[string]string{
					"en": "Queen's Birthday",
				},
			)
			
		case "SA": // South Australia
			// Adelaide Cup Day - 2nd Monday in March
			adelaideCup := NthWeekdayOfMonth(year, 3, time.Monday, 2)
			holidays[adelaideCup] = au.CreateHoliday(
				"Adelaide Cup Day",
				adelaideCup,
				"public",
				map[string]string{
					"en": "Adelaide Cup Day",
				},
			)
			
			// Labour Day - 1st Monday in October
			labourDaySA := NthWeekdayOfMonth(year, 10, time.Monday, 1)
			holidays[labourDaySA] = au.CreateHoliday(
				"Labour Day",
				labourDaySA,
				"public",
				map[string]string{
					"en": "Labour Day",
				},
			)
			
		case "TAS": // Tasmania
			// Eight Hours Day - 2nd Monday in March
			eightHoursDay := NthWeekdayOfMonth(year, 3, time.Monday, 2)
			holidays[eightHoursDay] = au.CreateHoliday(
				"Eight Hours Day",
				eightHoursDay,
				"public",
				map[string]string{
					"en": "Eight Hours Day",
				},
			)
			
		case "NT": // Northern Territory
			// May Day - 1st Monday in May
			mayDay := NthWeekdayOfMonth(year, 5, time.Monday, 1)
			holidays[mayDay] = au.CreateHoliday(
				"May Day",
				mayDay,
				"public",
				map[string]string{
					"en": "May Day",
				},
			)
			
			// Picnic Day - 1st Monday in August
			picnicDay := NthWeekdayOfMonth(year, 8, time.Monday, 1)
			holidays[picnicDay] = au.CreateHoliday(
				"Picnic Day",
				picnicDay,
				"public",
				map[string]string{
					"en": "Picnic Day",
				},
			)
		}
	}
	
	return holidays
}

// getMelbourneCupDay calculates Melbourne Cup Day (1st Tuesday in November)
func (au *AUProvider) getMelbourneCupDay(year int) time.Time {
	// Find the first Tuesday in November
	firstNov := time.Date(year, 11, 1, 0, 0, 0, 0, time.UTC)
	daysToTuesday := (int(time.Tuesday) - int(firstNov.Weekday()) + 7) % 7
	return firstNov.AddDate(0, 0, daysToTuesday)
}

// GetSeasons returns Australian seasons (opposite to Northern Hemisphere)
func (au *AUProvider) GetSeasons(year int) map[string][]time.Time {
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
