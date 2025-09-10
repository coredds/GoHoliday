package countries

import (
	"time"
)

// UAProvider implements holiday calculations for Ukraine
type UAProvider struct {
	*BaseProvider
}

// NewUAProvider creates a new Ukrainian holiday provider
func NewUAProvider() *UAProvider {
	base := NewBaseProvider("UA")
	base.subdivisions = []string{
		// 24 Oblasts (regions)
		"CK", // Cherkasy Oblast
		"CH", // Chernihiv Oblast
		"CV", // Chernivtsi Oblast
		"CR", // Crimea (disputed territory)
		"DP", // Dnipropetrovsk Oblast
		"DT", // Donetsk Oblast
		"IF", // Ivano-Frankivsk Oblast
		"KK", // Kharkiv Oblast
		"KS", // Kherson Oblast
		"KM", // Khmelnytskyi Oblast
		"KV", // Kiev Oblast
		"KR", // Kirovohrad Oblast
		"LH", // Luhansk Oblast
		"LV", // Lviv Oblast
		"MY", // Mykolaiv Oblast
		"OD", // Odessa Oblast
		"PL", // Poltava Oblast
		"RV", // Rivne Oblast
		"SM", // Sumy Oblast
		"TE", // Ternopil Oblast
		"VI", // Vinnytsia Oblast
		"VO", // Volyn Oblast
		"ZK", // Zakarpattia Oblast
		"ZP", // Zaporizhzhia Oblast
		"ZT", // Zhytomyr Oblast
		// Special administrative units
		"30", // Kyiv (capital city)
		"40", // Sevastopol (disputed)
	}
	base.categories = []string{
		"national",     // National holidays
		"orthodox",     // Orthodox Christian holidays
		"memorial",     // Memorial and remembrance days
		"professional", // Professional holidays
		"regional",     // Regional celebrations
		"cultural",     // Cultural and historical holidays
	}

	return &UAProvider{BaseProvider: base}
}

// LoadHolidays loads all Ukrainian holidays for a given year
func (ua *UAProvider) LoadHolidays(year int) map[time.Time]*Holiday {
	holidays := make(map[time.Time]*Holiday)

	// Add fixed national holidays
	ua.addNationalHolidays(holidays, year)

	// Add Orthodox Christian holidays (Easter-based)
	ua.addOrthodoxHolidays(holidays, year)

	// Add memorial and remembrance days
	ua.addMemorialHolidays(holidays, year)

	// Add professional holidays
	ua.addProfessionalHolidays(holidays, year)

	// Add cultural holidays
	ua.addCulturalHolidays(holidays, year)

	return holidays
}

// addNationalHolidays adds fixed national holidays of Ukraine
func (ua *UAProvider) addNationalHolidays(holidays map[time.Time]*Holiday, year int) {
	// New Year's Day
	newYear := time.Date(year, 1, 1, 0, 0, 0, 0, time.UTC)
	holidays[newYear] = ua.CreateHoliday(
		"New Year's Day",
		newYear,
		"national",
		map[string]string{
			"uk": "Новий рік",
			"en": "New Year's Day",
			"ru": "Новый год",
		},
	)

	// Orthodox Christmas
	orthodoxChristmas := time.Date(year, 1, 7, 0, 0, 0, 0, time.UTC)
	holidays[orthodoxChristmas] = ua.CreateHoliday(
		"Orthodox Christmas",
		orthodoxChristmas,
		"orthodox",
		map[string]string{
			"uk": "Різдво Христове",
			"en": "Orthodox Christmas",
			"ru": "Рождество Христово",
		},
	)

	// International Women's Day
	womensDay := time.Date(year, 3, 8, 0, 0, 0, 0, time.UTC)
	holidays[womensDay] = ua.CreateHoliday(
		"International Women's Day",
		womensDay,
		"national",
		map[string]string{
			"uk": "Міжнародний жіночий день",
			"en": "International Women's Day",
			"ru": "Международный женский день",
		},
	)

	// Labor Day
	laborDay := time.Date(year, 5, 1, 0, 0, 0, 0, time.UTC)
	holidays[laborDay] = ua.CreateHoliday(
		"Labor Day",
		laborDay,
		"national",
		map[string]string{
			"uk": "День праці",
			"en": "Labor Day",
			"ru": "День труда",
		},
	)

	// Victory Day (over Nazism in World War II)
	victoryDay := time.Date(year, 5, 8, 0, 0, 0, 0, time.UTC)
	holidays[victoryDay] = ua.CreateHoliday(
		"Victory Day",
		victoryDay,
		"memorial",
		map[string]string{
			"uk": "День перемоги над нацизмом у Другій світовій війні",
			"en": "Victory Day",
			"ru": "День победы",
		},
	)

	// Constitution Day
	constitutionDay := time.Date(year, 6, 28, 0, 0, 0, 0, time.UTC)
	holidays[constitutionDay] = ua.CreateHoliday(
		"Constitution Day",
		constitutionDay,
		"national",
		map[string]string{
			"uk": "День Конституції України",
			"en": "Constitution Day",
			"ru": "День Конституции Украины",
		},
	)

	// Independence Day
	independenceDay := time.Date(year, 8, 24, 0, 0, 0, 0, time.UTC)
	holidays[independenceDay] = ua.CreateHoliday(
		"Independence Day",
		independenceDay,
		"national",
		map[string]string{
			"uk": "День незалежності України",
			"en": "Independence Day",
			"ru": "День независимости Украины",
		},
	)

	// Defenders Day (since 2015, replaces Defender of the Fatherland Day)
	if year >= 2015 {
		defendersDay := time.Date(year, 10, 14, 0, 0, 0, 0, time.UTC)
		holidays[defendersDay] = ua.CreateHoliday(
			"Defenders Day",
			defendersDay,
			"memorial",
			map[string]string{
				"uk": "День захисників і захисниць України",
				"en": "Defenders Day",
				"ru": "День защитников Украины",
			},
		)
	}
}

// addOrthodoxHolidays adds Orthodox Christian holidays based on Easter calculation
func (ua *UAProvider) addOrthodoxHolidays(holidays map[time.Time]*Holiday, year int) {
	// Calculate Orthodox Easter (using Julian calendar)
	easter := ua.calculateOrthodoxEaster(year)

	// Palm Sunday (1 week before Easter)
	palmSunday := easter.AddDate(0, 0, -7)
	holidays[palmSunday] = ua.CreateHoliday(
		"Palm Sunday",
		palmSunday,
		"orthodox",
		map[string]string{
			"uk": "Вербна неділя",
			"en": "Palm Sunday",
			"ru": "Вербное воскресенье",
		},
	)

	// Orthodox Easter
	holidays[easter] = ua.CreateHoliday(
		"Orthodox Easter",
		easter,
		"orthodox",
		map[string]string{
			"uk": "Великдень",
			"en": "Orthodox Easter",
			"ru": "Пасха",
		},
	)

	// Trinity Sunday (49 days after Easter)
	trinity := easter.AddDate(0, 0, 49)
	holidays[trinity] = ua.CreateHoliday(
		"Trinity Sunday",
		trinity,
		"orthodox",
		map[string]string{
			"uk": "Трійця",
			"en": "Trinity Sunday",
			"ru": "Троица",
		},
	)
}

// addMemorialHolidays adds memorial and remembrance days
func (ua *UAProvider) addMemorialHolidays(holidays map[time.Time]*Holiday, year int) {
	// Holodomor Remembrance Day
	holodomorDay := time.Date(year, 11, 25, 0, 0, 0, 0, time.UTC)
	holidays[holodomorDay] = ua.CreateHoliday(
		"Holodomor Remembrance Day",
		holodomorDay,
		"memorial",
		map[string]string{
			"uk": "День пам'яті жертв голодомору",
			"en": "Holodomor Remembrance Day",
			"ru": "День памяти жертв голодомора",
		},
	)

	// Day of Dignity and Freedom (since 2014)
	if year >= 2014 {
		dignityDay := time.Date(year, 11, 21, 0, 0, 0, 0, time.UTC)
		holidays[dignityDay] = ua.CreateHoliday(
			"Day of Dignity and Freedom",
			dignityDay,
			"memorial",
			map[string]string{
				"uk": "День Гідності та Свободи",
				"en": "Day of Dignity and Freedom",
				"ru": "День достоинства и свободы",
			},
		)
	}

	// Day of Remembrance of Victims of Political Repressions
	repressionDay := time.Date(year, 5, 19, 0, 0, 0, 0, time.UTC)
	holidays[repressionDay] = ua.CreateHoliday(
		"Day of Remembrance of Victims of Political Repressions",
		repressionDay,
		"memorial",
		map[string]string{
			"uk": "День пам'яті жертв політичних репресій",
			"en": "Day of Remembrance of Victims of Political Repressions",
			"ru": "День памяти жертв политических репрессий",
		},
	)
}

// addProfessionalHolidays adds professional and occupational holidays
func (ua *UAProvider) addProfessionalHolidays(holidays map[time.Time]*Holiday, year int) {
	// Day of Ukrainian Language (since 2019)
	if year >= 2019 {
		languageDay := time.Date(year, 11, 9, 0, 0, 0, 0, time.UTC)
		holidays[languageDay] = ua.CreateHoliday(
			"Day of Ukrainian Language",
			languageDay,
			"cultural",
			map[string]string{
				"uk": "День української писемності та мови",
				"en": "Day of Ukrainian Language",
				"ru": "День украинской письменности и языка",
			},
		)
	}

	// Day of Ukrainian Statehood (since 2021)
	if year >= 2021 {
		statehoodDay := time.Date(year, 7, 28, 0, 0, 0, 0, time.UTC)
		holidays[statehoodDay] = ua.CreateHoliday(
			"Day of Ukrainian Statehood",
			statehoodDay,
			"national",
			map[string]string{
				"uk": "День української державності",
				"en": "Day of Ukrainian Statehood",
				"ru": "День украинской государственности",
			},
		)
	}
}

// addCulturalHolidays adds cultural and traditional holidays
func (ua *UAProvider) addCulturalHolidays(holidays map[time.Time]*Holiday, year int) {
	// Old New Year (Julian calendar New Year)
	oldNewYear := time.Date(year, 1, 14, 0, 0, 0, 0, time.UTC)
	holidays[oldNewYear] = ua.CreateHoliday(
		"Old New Year",
		oldNewYear,
		"cultural",
		map[string]string{
			"uk": "Старий Новий рік",
			"en": "Old New Year",
			"ru": "Старый Новый год",
		},
	)

	// Day of Ukrainian Cossacks (since 1999)
	if year >= 1999 {
		cossacksDay := time.Date(year, 10, 14, 0, 0, 0, 0, time.UTC)
		// Note: This is the same date as Defenders Day since 2015
		if year < 2015 {
			holidays[cossacksDay] = ua.CreateHoliday(
				"Day of Ukrainian Cossacks",
				cossacksDay,
				"cultural",
				map[string]string{
					"uk": "День українського козацтва",
					"en": "Day of Ukrainian Cossacks",
					"ru": "День украинского казачества",
				},
			)
		}
	}
}

// calculateOrthodoxEaster calculates Orthodox Easter using the Julian calendar
// Based on the algorithm for calculating Orthodox Easter
func (ua *UAProvider) calculateOrthodoxEaster(year int) time.Time {
	// Orthodox Easter calculation using the Julian calendar
	// This is a simplified version - in practice, you might want to use
	// a more sophisticated astronomical calculation

	// Calculate Julian Easter first
	a := year % 4
	b := year % 7
	c := year % 19
	d := (19*c + 15) % 30
	e := (2*a + 4*b - d + 34) % 7
	month := (d + e + 114) / 31
	day := ((d + e + 114) % 31) + 1

	// Convert Julian to Gregorian calendar
	julianEaster := time.Date(year, time.Month(month), day, 0, 0, 0, 0, time.UTC)

	// Add the difference between Julian and Gregorian calendars
	// For years 1900-2099, the difference is 13 days
	var gregorianOffset int
	if year >= 1900 && year <= 2099 {
		gregorianOffset = 13
	} else if year >= 2100 && year <= 2199 {
		gregorianOffset = 14
	} else {
		// Simplified - for more accuracy, implement full Julian/Gregorian conversion
		gregorianOffset = 13
	}

	return julianEaster.AddDate(0, 0, gregorianOffset)
}

// CreateHoliday creates a holiday with Ukrainian-specific formatting
func (ua *UAProvider) CreateHoliday(name string, date time.Time, category string, languages map[string]string) *Holiday {
	return &Holiday{
		Name:         name,
		Date:         date,
		Category:     category,
		Languages:    languages,
		IsObserved:   true,
		Subdivisions: []string{}, // Can be populated for regional holidays
	}
}

// GetCountryCode returns the country code for Ukraine
func (ua *UAProvider) GetCountryCode() string {
	return "UA"
}

// GetSupportedSubdivisions returns the list of supported Ukrainian subdivisions
func (ua *UAProvider) GetSupportedSubdivisions() []string {
	return ua.subdivisions
}

// GetSupportedCategories returns the list of supported holiday categories
func (ua *UAProvider) GetSupportedCategories() []string {
	return ua.categories
}

// GetName returns the country name
func (ua *UAProvider) GetName() string {
	return "Ukraine"
}

// GetLanguages returns the supported languages
func (ua *UAProvider) GetLanguages() []string {
	return []string{"uk", "en", "ru"}
}

// IsSubdivisionSupported checks if a subdivision is supported
func (ua *UAProvider) IsSubdivisionSupported(subdivision string) bool {
	for _, s := range ua.subdivisions {
		if s == subdivision {
			return true
		}
	}
	return false
}

// IsCategorySupported checks if a category is supported
func (ua *UAProvider) IsCategorySupported(category string) bool {
	for _, c := range ua.categories {
		if c == category {
			return true
		}
	}
	return false
}
