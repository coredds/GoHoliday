package countries

import (
	"time"
)

// IEProvider implements holiday logic for Ireland
type IEProvider struct {
	BaseProvider
}

// NewIEProvider creates a new Ireland holiday provider
func NewIEProvider() *IEProvider {
	return &IEProvider{
		BaseProvider: *NewBaseProvider("IE"),
	}
}

// GetCountryCode returns the country code
func (ie *IEProvider) GetCountryCode() string {
	return ie.BaseProvider.GetCountryCode()
}

// GetCountryName returns the country name
func (ie *IEProvider) GetCountryName() string {
	return "Ireland"
}

// GetSubdivisions returns Irish counties
func (ie *IEProvider) GetSubdivisions() []string {
	return []string{
		"C",  // Connacht
		"L",  // Leinster
		"M",  // Munster
		"U",  // Ulster (part of)
		"CW", // Carlow
		"CN", // Cavan
		"CE", // Clare
		"CO", // Cork
		"DL", // Donegal
		"D",  // Dublin
		"G",  // Galway
		"KY", // Kerry
		"KE", // Kildare
		"KK", // Kilkenny
		"LS", // Laois
		"LM", // Leitrim
		"LK", // Limerick
		"LD", // Longford
		"LH", // Louth
		"MO", // Mayo
		"MH", // Meath
		"MN", // Monaghan
		"OY", // Offaly
		"RN", // Roscommon
		"SO", // Sligo
		"TA", // Tipperary
		"WD", // Waterford
		"WH", // Westmeath
		"WX", // Wexford
		"WW", // Wicklow
	}
}

// GetCategories returns holiday categories used in Ireland
func (ie *IEProvider) GetCategories() []string {
	return []string{"public", "bank", "religious", "national", "cultural"}
}

// LoadHolidays loads Irish holidays for the specified year
func (ie *IEProvider) LoadHolidays(year int) map[time.Time]*Holiday {
	holidays := make(map[time.Time]*Holiday)

	// Fixed public holidays
	ie.addFixedHolidays(holidays, year)

	// Easter-based holidays
	ie.addEasterBasedHolidays(holidays, year)

	// Bank holidays (first Monday in months)
	ie.addBankHolidays(holidays, year)

	// Cultural and national holidays
	ie.addCulturalHolidays(holidays, year)

	return holidays
}

// addFixedHolidays adds fixed date holidays
func (ie *IEProvider) addFixedHolidays(holidays map[time.Time]*Holiday, year int) {
	fixedHolidays := []struct {
		month    int
		day      int
		name     string
		nameGA   string // Irish Gaelic
		category string
	}{
		{1, 1, "New Year's Day", "Lá na Bliana Nua", "public"},
		{3, 17, "Saint Patrick's Day", "Lá Fhéile Pádraig", "national"},
		{12, 25, "Christmas Day", "Lá na Nollag", "religious"},
		{12, 26, "Saint Stephen's Day", "Lá Fhéile Stiofáin", "religious"},
	}

	for _, h := range fixedHolidays {
		date := time.Date(year, time.Month(h.month), h.day, 0, 0, 0, 0, time.UTC)
		holidays[date] = ie.CreateHoliday(
			h.name,
			date,
			h.category,
			map[string]string{
				"en": h.name,
				"ga": h.nameGA,
			},
		)
	}
}

// addEasterBasedHolidays adds holidays based on Easter calculation
func (ie *IEProvider) addEasterBasedHolidays(holidays map[time.Time]*Holiday, year int) {
	easter := EasterSunday(year)

	// Good Friday (not a public holiday but widely observed)
	goodFriday := easter.AddDate(0, 0, -2)
	holidays[goodFriday] = ie.CreateHoliday(
		"Good Friday",
		goodFriday,
		"religious",
		map[string]string{
			"en": "Good Friday",
			"ga": "Aoine an Chéasta",
		},
	)

	// Easter Monday (public holiday)
	easterMonday := easter.AddDate(0, 0, 1)
	holidays[easterMonday] = ie.CreateHoliday(
		"Easter Monday",
		easterMonday,
		"public",
		map[string]string{
			"en": "Easter Monday",
			"ga": "Luan Cásca",
		},
	)
}

// addBankHolidays adds bank holidays (first Monday of certain months)
func (ie *IEProvider) addBankHolidays(holidays map[time.Time]*Holiday, year int) {
	// First Monday in June - June Bank Holiday
	juneHoliday := ie.getFirstMondayOfMonth(year, 6)
	holidays[juneHoliday] = ie.CreateHoliday(
		"June Bank Holiday",
		juneHoliday,
		"bank",
		map[string]string{
			"en": "June Bank Holiday",
			"ga": "Lá Saoire an Bhainc Meitheamh",
		},
	)

	// First Monday in August - August Bank Holiday
	augustHoliday := ie.getFirstMondayOfMonth(year, 8)
	holidays[augustHoliday] = ie.CreateHoliday(
		"August Bank Holiday",
		augustHoliday,
		"bank",
		map[string]string{
			"en": "August Bank Holiday",
			"ga": "Lá Saoire an Bhainc Lúnasa",
		},
	)

	// Last Monday in October - October Bank Holiday
	octoberHoliday := ie.getLastMondayOfMonth(year, 10)
	holidays[octoberHoliday] = ie.CreateHoliday(
		"October Bank Holiday",
		octoberHoliday,
		"bank",
		map[string]string{
			"en": "October Bank Holiday",
			"ga": "Lá Saoire an Bhainc Deireadh Fómhair",
		},
	)
}

// addCulturalHolidays adds cultural and historical holidays
func (ie *IEProvider) addCulturalHolidays(holidays map[time.Time]*Holiday, year int) {
	// May Day (May 1st) - not a public holiday but culturally significant
	mayDay := time.Date(year, 5, 1, 0, 0, 0, 0, time.UTC)
	holidays[mayDay] = ie.CreateHoliday(
		"May Day",
		mayDay,
		"cultural",
		map[string]string{
			"en": "May Day",
			"ga": "Lá Bealtaine",
		},
	)

	// Brigid's Day (February 1st) - traditional Celtic festival
	brigidsDay := time.Date(year, 2, 1, 0, 0, 0, 0, time.UTC)
	holidays[brigidsDay] = ie.CreateHoliday(
		"Saint Brigid's Day",
		brigidsDay,
		"cultural",
		map[string]string{
			"en": "Saint Brigid's Day",
			"ga": "Lá Fhéile Bríde",
		},
	)

	// Lughnasadh (August 1st) - ancient Celtic harvest festival
	lughnasadh := time.Date(year, 8, 1, 0, 0, 0, 0, time.UTC)
	holidays[lughnasadh] = ie.CreateHoliday(
		"Lughnasadh",
		lughnasadh,
		"cultural",
		map[string]string{
			"en": "Lughnasadh",
			"ga": "Lúnasa",
		},
	)

	// Samhain (October 31st) - ancient Celtic festival, origin of Halloween
	samhain := time.Date(year, 10, 31, 0, 0, 0, 0, time.UTC)
	holidays[samhain] = ie.CreateHoliday(
		"Samhain",
		samhain,
		"cultural",
		map[string]string{
			"en": "Samhain",
			"ga": "Samhain",
		},
	)

	// Add modern public holidays introduced recently
	if year >= 2023 {
		// Saint Brigid's Day became a public holiday in 2023
		// It's observed on the first Monday in February if February 1st is not a Monday
		brigidsPublicHoliday := ie.getBrigidsPublicHoliday(year)
		holidays[brigidsPublicHoliday] = ie.CreateHoliday(
			"Saint Brigid's Day (Public Holiday)",
			brigidsPublicHoliday,
			"public",
			map[string]string{
				"en": "Saint Brigid's Day",
				"ga": "Lá Fhéile Bríde",
			},
		)
	}
}

// getFirstMondayOfMonth returns the first Monday of the specified month
func (ie *IEProvider) getFirstMondayOfMonth(year, month int) time.Time {
	firstDay := time.Date(year, time.Month(month), 1, 0, 0, 0, 0, time.UTC)
	
	// Find the first Monday
	daysUntilMonday := (int(time.Monday) - int(firstDay.Weekday()) + 7) % 7
	return firstDay.AddDate(0, 0, daysUntilMonday)
}

// getLastMondayOfMonth returns the last Monday of the specified month
func (ie *IEProvider) getLastMondayOfMonth(year, month int) time.Time {
	// Get the last day of the month
	lastDay := time.Date(year, time.Month(month+1), 0, 0, 0, 0, 0, time.UTC)
	
	// Find the last Monday
	daysBackToMonday := (int(lastDay.Weekday()) - int(time.Monday) + 7) % 7
	return lastDay.AddDate(0, 0, -daysBackToMonday)
}

// getBrigidsPublicHoliday returns the observed date for Saint Brigid's Day public holiday
func (ie *IEProvider) getBrigidsPublicHoliday(year int) time.Time {
	brigidsDay := time.Date(year, 2, 1, 0, 0, 0, 0, time.UTC)
	
	// If February 1st is a Monday, it's observed on that day
	if brigidsDay.Weekday() == time.Monday {
		return brigidsDay
	}
	
	// Otherwise, it's observed on the first Monday in February
	return ie.getFirstMondayOfMonth(year, 2)
}
