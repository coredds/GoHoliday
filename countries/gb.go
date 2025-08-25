package countries

import (
	"time"
)

// GBProvider implements holiday calculations for the United Kingdom
type GBProvider struct {
	*BaseProvider
}

// NewGBProvider creates a new UK holiday provider
func NewGBProvider() *GBProvider {
	base := NewBaseProvider("GB")
	base.subdivisions = []string{
		"ENG", "SCT", "WLS", "NIR", // England, Scotland, Wales, Northern Ireland
	}
	base.categories = []string{"public", "bank", "government"}
	
	return &GBProvider{BaseProvider: base}
}

// LoadHolidays loads all UK holidays for a given year
func (gb *GBProvider) LoadHolidays(year int) map[time.Time]*Holiday {
	holidays := make(map[time.Time]*Holiday)
	
	// Fixed date holidays
	holidays[time.Date(year, 1, 1, 0, 0, 0, 0, time.UTC)] = gb.CreateHoliday(
		"New Year's Day",
		time.Date(year, 1, 1, 0, 0, 0, 0, time.UTC),
		"public",
		map[string]string{
			"en": "New Year's Day",
		},
	)
	
	holidays[time.Date(year, 12, 25, 0, 0, 0, 0, time.UTC)] = gb.CreateHoliday(
		"Christmas Day",
		time.Date(year, 12, 25, 0, 0, 0, 0, time.UTC),
		"public",
		map[string]string{
			"en": "Christmas Day",
		},
	)
	
	holidays[time.Date(year, 12, 26, 0, 0, 0, 0, time.UTC)] = gb.CreateHoliday(
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
	holidays[goodFriday] = gb.CreateHoliday(
		"Good Friday",
		goodFriday,
		"public",
		map[string]string{
			"en": "Good Friday",
		},
	)
	
	// Easter Monday
	easterMonday := easter.AddDate(0, 0, 1)
	holidays[easterMonday] = gb.CreateHoliday(
		"Easter Monday",
		easterMonday,
		"public",
		map[string]string{
			"en": "Easter Monday",
		},
	)
	
	// Variable date holidays
	
	// Early May Bank Holiday - 1st Monday in May
	earlyMayBankHoliday := NthWeekdayOfMonth(year, 5, time.Monday, 1)
	holidays[earlyMayBankHoliday] = gb.CreateHoliday(
		"Early May Bank Holiday",
		earlyMayBankHoliday,
		"bank",
		map[string]string{
			"en": "Early May Bank Holiday",
		},
	)
	
	// Spring Bank Holiday - Last Monday in May
	springBankHoliday := NthWeekdayOfMonth(year, 5, time.Monday, -1)
	holidays[springBankHoliday] = gb.CreateHoliday(
		"Spring Bank Holiday",
		springBankHoliday,
		"bank",
		map[string]string{
			"en": "Spring Bank Holiday",
		},
	)
	
	// Summer Bank Holiday - Last Monday in August
	summerBankHoliday := NthWeekdayOfMonth(year, 8, time.Monday, -1)
	holidays[summerBankHoliday] = gb.CreateHoliday(
		"Summer Bank Holiday",
		summerBankHoliday,
		"bank",
		map[string]string{
			"en": "Summer Bank Holiday",
		},
	)
	
	// Special holidays for specific years
	gb.addSpecialHolidays(year, holidays)
	
	return holidays
}

// addSpecialHolidays adds one-off special holidays for specific years
func (gb *GBProvider) addSpecialHolidays(year int, holidays map[time.Time]*Holiday) {
	switch year {
	case 2022:
		// Platinum Jubilee
		platinumJubilee := time.Date(2022, 6, 3, 0, 0, 0, 0, time.UTC)
		holidays[platinumJubilee] = gb.CreateHoliday(
			"Platinum Jubilee",
			platinumJubilee,
			"public",
			map[string]string{
				"en": "Platinum Jubilee",
			},
		)
	case 2023:
		// Coronation of King Charles III
		coronation := time.Date(2023, 5, 8, 0, 0, 0, 0, time.UTC)
		holidays[coronation] = gb.CreateHoliday(
			"Coronation of King Charles III",
			coronation,
			"public",
			map[string]string{
				"en": "Coronation of King Charles III",
			},
		)
	}
}

// GetRegionalHolidays returns region-specific holidays
func (gb *GBProvider) GetRegionalHolidays(year int, subdivisions []string) map[time.Time]*Holiday {
	holidays := make(map[time.Time]*Holiday)
	
	for _, region := range subdivisions {
		switch region {
		case "SCT": // Scotland
			// St. Andrew's Day - November 30
			stAndrewsDay := time.Date(year, 11, 30, 0, 0, 0, 0, time.UTC)
			holidays[stAndrewsDay] = gb.CreateHoliday(
				"St. Andrew's Day",
				stAndrewsDay,
				"public",
				map[string]string{
					"en": "St. Andrew's Day",
				},
			)
			
		case "WLS": // Wales
			// St. David's Day - March 1
			stDavidsDay := time.Date(year, 3, 1, 0, 0, 0, 0, time.UTC)
			holidays[stDavidsDay] = gb.CreateHoliday(
				"St. David's Day",
				stDavidsDay,
				"public",
				map[string]string{
					"en": "St. David's Day",
				},
			)
			
		case "NIR": // Northern Ireland
			// St. Patrick's Day - March 17
			stPatricksDay := time.Date(year, 3, 17, 0, 0, 0, 0, time.UTC)
			holidays[stPatricksDay] = gb.CreateHoliday(
				"St. Patrick's Day",
				stPatricksDay,
				"public",
				map[string]string{
					"en": "St. Patrick's Day",
				},
			)
			
			// Battle of the Boyne - July 12
			battleOfBoyne := time.Date(year, 7, 12, 0, 0, 0, 0, time.UTC)
			holidays[battleOfBoyne] = gb.CreateHoliday(
				"Battle of the Boyne",
				battleOfBoyne,
				"public",
				map[string]string{
					"en": "Battle of the Boyne",
				},
			)
		}
	}
	
	return holidays
}
