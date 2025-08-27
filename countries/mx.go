package countries

import (
	"time"
)

// MXProvider implements holiday calculations for Mexico
type MXProvider struct {
	*BaseProvider
}

// NewMXProvider creates a new Mexican holiday provider
func NewMXProvider() *MXProvider {
	base := NewBaseProvider("MX")
	base.subdivisions = []string{
		// 32 federal entities (31 states + 1 federal district)
		"AGU", "BCN", "BCS", "CAM", "CHP", "CHH", "COA", "COL", "CMX",
		"DUR", "GUA", "GRO", "HID", "JAL", "MEX", "MIC", "MOR", "NAY",
		"NLE", "OAX", "PUE", "QUE", "ROO", "SLP", "SIN", "SON", "TAB",
		"TAM", "TLA", "VER", "YUC", "ZAC",
	}
	base.categories = []string{"public", "national", "religious", "regional", "civic"}

	return &MXProvider{BaseProvider: base}
}

// LoadHolidays loads all Mexican holidays for a given year
func (mx *MXProvider) LoadHolidays(year int) map[time.Time]*Holiday {
	holidays := make(map[time.Time]*Holiday)

	// Fixed national holidays

	// New Year's Day - January 1
	newYear := time.Date(year, 1, 1, 0, 0, 0, 0, time.UTC)
	holidays[newYear] = mx.CreateHoliday(
		"Año Nuevo",
		newYear,
		"public",
		map[string]string{
			"es": "Año Nuevo",
			"en": "New Year's Day",
		},
	)

	// Constitution Day - First Monday of February (since 2006)
	constitutionDay := mx.getFirstMondayOfFebruary(year)
	holidays[constitutionDay] = mx.CreateHoliday(
		"Día de la Constitución",
		constitutionDay,
		"civic",
		map[string]string{
			"es": "Día de la Constitución",
			"en": "Constitution Day",
		},
	)

	// Benito Juárez's Birthday - Third Monday of March (since 2006)
	juarezBirthday := mx.getThirdMondayOfMarch(year)
	holidays[juarezBirthday] = mx.CreateHoliday(
		"Natalicio de Benito Juárez",
		juarezBirthday,
		"civic",
		map[string]string{
			"es": "Natalicio de Benito Juárez",
			"en": "Benito Juárez's Birthday",
		},
	)

	// Labour Day - May 1
	labourDay := time.Date(year, 5, 1, 0, 0, 0, 0, time.UTC)
	holidays[labourDay] = mx.CreateHoliday(
		"Día del Trabajo",
		labourDay,
		"public",
		map[string]string{
			"es": "Día del Trabajo",
			"en": "Labour Day",
		},
	)

	// Independence Day - September 16
	independence := time.Date(year, 9, 16, 0, 0, 0, 0, time.UTC)
	holidays[independence] = mx.CreateHoliday(
		"Día de la Independencia",
		independence,
		"national",
		map[string]string{
			"es": "Día de la Independencia",
			"en": "Independence Day",
		},
	)

	// Revolution Day - Third Monday of November (since 2006)
	revolutionDay := mx.getThirdMondayOfNovember(year)
	holidays[revolutionDay] = mx.CreateHoliday(
		"Día de la Revolución",
		revolutionDay,
		"civic",
		map[string]string{
			"es": "Día de la Revolución",
			"en": "Revolution Day",
		},
	)

	// Christmas Day - December 25
	christmas := time.Date(year, 12, 25, 0, 0, 0, 0, time.UTC)
	holidays[christmas] = mx.CreateHoliday(
		"Navidad",
		christmas,
		"religious",
		map[string]string{
			"es": "Navidad",
			"en": "Christmas Day",
		},
	)

	// Variable holidays based on Easter
	easter := EasterSunday(year)

	// Maundy Thursday (3 days before Easter)
	maundyThursday := easter.AddDate(0, 0, -3)
	holidays[maundyThursday] = mx.CreateHoliday(
		"Jueves Santo",
		maundyThursday,
		"religious",
		map[string]string{
			"es": "Jueves Santo",
			"en": "Maundy Thursday",
		},
	)

	// Good Friday (2 days before Easter)
	goodFriday := easter.AddDate(0, 0, -2)
	holidays[goodFriday] = mx.CreateHoliday(
		"Viernes Santo",
		goodFriday,
		"religious",
		map[string]string{
			"es": "Viernes Santo",
			"en": "Good Friday",
		},
	)

	// Easter Saturday (1 day before Easter)
	easterSaturday := easter.AddDate(0, 0, -1)
	holidays[easterSaturday] = mx.CreateHoliday(
		"Sábado de Gloria",
		easterSaturday,
		"religious",
		map[string]string{
			"es": "Sábado de Gloria",
			"en": "Easter Saturday",
		},
	)

	// Additional observances (not official holidays but widely observed)

	// Day of the Dead - November 2
	dayOfDead := time.Date(year, 11, 2, 0, 0, 0, 0, time.UTC)
	holidays[dayOfDead] = mx.CreateHoliday(
		"Día de los Muertos",
		dayOfDead,
		"religious",
		map[string]string{
			"es": "Día de los Muertos",
			"en": "Day of the Dead",
		},
	)

	// Our Lady of Guadalupe - December 12
	guadalupe := time.Date(year, 12, 12, 0, 0, 0, 0, time.UTC)
	holidays[guadalupe] = mx.CreateHoliday(
		"Día de la Virgen de Guadalupe",
		guadalupe,
		"religious",
		map[string]string{
			"es": "Día de la Virgen de Guadalupe",
			"en": "Our Lady of Guadalupe",
		},
	)

	return holidays
}

// Helper methods for calculating variable holidays

// getFirstMondayOfFebruary returns the first Monday of February
func (mx *MXProvider) getFirstMondayOfFebruary(year int) time.Time {
	return mx.getNthWeekdayOfMonth(year, time.February, time.Monday, 1)
}

// getThirdMondayOfMarch returns the third Monday of March
func (mx *MXProvider) getThirdMondayOfMarch(year int) time.Time {
	return mx.getNthWeekdayOfMonth(year, time.March, time.Monday, 3)
}

// getThirdMondayOfNovember returns the third Monday of November
func (mx *MXProvider) getThirdMondayOfNovember(year int) time.Time {
	return mx.getNthWeekdayOfMonth(year, time.November, time.Monday, 3)
}

// getNthWeekdayOfMonth calculates the nth occurrence of a weekday in a month
func (mx *MXProvider) getNthWeekdayOfMonth(year int, month time.Month, weekday time.Weekday, n int) time.Time {
	// Find the first occurrence of the weekday in the month
	firstDay := time.Date(year, month, 1, 0, 0, 0, 0, time.UTC)
	daysToWeekday := (int(weekday) - int(firstDay.Weekday()) + 7) % 7
	firstOccurrence := firstDay.AddDate(0, 0, daysToWeekday)

	// Add weeks to get the nth occurrence
	return firstOccurrence.AddDate(0, 0, (n-1)*7)
}

// GetCountryCode returns the country code for Mexico
func (mx *MXProvider) GetCountryCode() string {
	return "MX"
}

// GetSupportedSubdivisions returns the list of supported Mexican states
func (mx *MXProvider) GetSupportedSubdivisions() []string {
	return mx.subdivisions
}

// GetSupportedCategories returns the list of supported holiday categories
func (mx *MXProvider) GetSupportedCategories() []string {
	return mx.categories
}

// GetName returns the country name
func (mx *MXProvider) GetName() string {
	return "Mexico"
}

// GetLanguages returns the supported languages
func (mx *MXProvider) GetLanguages() []string {
	return []string{"es", "en"}
}

// IsSubdivisionSupported checks if a subdivision is supported
func (mx *MXProvider) IsSubdivisionSupported(subdivision string) bool {
	for _, s := range mx.subdivisions {
		if s == subdivision {
			return true
		}
	}
	return false
}

// IsCategorySupported checks if a category is supported
func (mx *MXProvider) IsCategorySupported(category string) bool {
	for _, c := range mx.categories {
		if c == category {
			return true
		}
	}
	return false
}
