package countries

import (
	"time"
)

// CLProvider implements holiday logic for Chile
type CLProvider struct {
	BaseProvider
}

// NewCLProvider creates a new Chile holiday provider
func NewCLProvider() *CLProvider {
	return &CLProvider{
		BaseProvider: *NewBaseProvider("CL"),
	}
}

// GetCountryCode returns the country code
func (cl *CLProvider) GetCountryCode() string {
	return cl.BaseProvider.GetCountryCode()
}

// GetCountryName returns the country name
func (cl *CLProvider) GetCountryName() string {
	return "Chile"
}

// GetSubdivisions returns Chilean regions
func (cl *CLProvider) GetSubdivisions() []string {
	return []string{
		"AI", // Aisén
		"AN", // Antofagasta
		"AP", // Arica y Parinacota
		"AR", // Araucanía
		"AT", // Atacama
		"BI", // Biobío
		"CO", // Coquimbo
		"LI", // Libertador General Bernardo O'Higgins
		"LL", // Los Lagos
		"LR", // Los Ríos
		"MA", // Magallanes y Antártica Chilena
		"ML", // Maule
		"NB", // Ñuble
		"RM", // Región Metropolitana de Santiago
		"TA", // Tarapacá
		"VS", // Valparaíso
	}
}

// GetCategories returns holiday categories used in Chile
func (cl *CLProvider) GetCategories() []string {
	return []string{"public", "religious", "civic", "regional"}
}

// LoadHolidays loads Chilean holidays for the specified year
func (cl *CLProvider) LoadHolidays(year int) map[time.Time]*Holiday {
	holidays := make(map[time.Time]*Holiday)

	// Fixed national holidays
	cl.addFixedHolidays(holidays, year)

	// Easter-based holidays
	cl.addEasterBasedHolidays(holidays, year)

	// Variable holidays (Monday holidays law)
	cl.addVariableHolidays(holidays, year)

	// Regional holidays
	cl.addRegionalHolidays(holidays, year)

	return holidays
}

// addFixedHolidays adds fixed date holidays
func (cl *CLProvider) addFixedHolidays(holidays map[time.Time]*Holiday, year int) {
	fixedHolidays := []struct {
		month    int
		day      int
		name     string
		nameES   string
		category string
	}{
		{1, 1, "New Year's Day", "Año Nuevo", "public"},
		{5, 1, "Labour Day", "Día del Trabajador", "public"},
		{5, 21, "Navy Day", "Día de las Glorias Navales", "civic"},
		{9, 18, "Independence Day", "Día de la Independencia", "public"},
		{9, 19, "Army Day", "Día de las Glorias del Ejército", "civic"},
		{10, 12, "Columbus Day", "Día del Encuentro de Dos Mundos", "public"},
		{11, 1, "All Saints' Day", "Día de Todos los Santos", "religious"},
		{12, 8, "Immaculate Conception", "Inmaculada Concepción", "religious"},
		{12, 25, "Christmas Day", "Navidad", "religious"},
	}

	for _, h := range fixedHolidays {
		date := time.Date(year, time.Month(h.month), h.day, 0, 0, 0, 0, time.UTC)
		holidays[date] = cl.CreateHoliday(
			h.name,
			date,
			h.category,
			map[string]string{
				"en": h.name,
				"es": h.nameES,
			},
		)
	}
}

// addEasterBasedHolidays adds holidays based on Easter calculation
func (cl *CLProvider) addEasterBasedHolidays(holidays map[time.Time]*Holiday, year int) {
	easter := EasterSunday(year)

	// Good Friday (Viernes Santo)
	goodFriday := easter.AddDate(0, 0, -2)
	holidays[goodFriday] = cl.CreateHoliday(
		"Good Friday",
		goodFriday,
		"religious",
		map[string]string{
			"en": "Good Friday",
			"es": "Viernes Santo",
		},
	)

	// Holy Saturday (Sábado Santo)
	holySaturday := easter.AddDate(0, 0, -1)
	holidays[holySaturday] = cl.CreateHoliday(
		"Holy Saturday",
		holySaturday,
		"religious",
		map[string]string{
			"en": "Holy Saturday",
			"es": "Sábado Santo",
		},
	)
}

// addVariableHolidays adds holidays that can be moved to Monday (Ley de Feriados)
func (cl *CLProvider) addVariableHolidays(holidays map[time.Time]*Holiday, year int) {
	// These holidays are moved to Monday if they fall on Tuesday-Friday
	// This is part of Chile's "Ley de Feriados" (Holidays Law)

	variableHolidays := []struct {
		month    int
		day      int
		name     string
		nameES   string
		category string
	}{
		{6, 29, "Saint Peter and Saint Paul", "San Pedro y San Pablo", "religious"},
		{8, 15, "Assumption of Mary", "Asunción de la Virgen", "religious"},
		{10, 31, "Reformation Day", "Día de las Iglesias Evangélicas y Protestantes", "religious"},
	}

	for _, h := range variableHolidays {
		originalDate := time.Date(year, time.Month(h.month), h.day, 0, 0, 0, 0, time.UTC)
		observedDate := cl.getObservedDate(originalDate)

		holidays[observedDate] = cl.CreateHoliday(
			h.name,
			observedDate,
			h.category,
			map[string]string{
				"en": h.name,
				"es": h.nameES,
			},
		)
	}
}

// addRegionalHolidays adds region-specific holidays
func (cl *CLProvider) addRegionalHolidays(holidays map[time.Time]*Holiday, year int) {
	// Arica y Parinacota region
	if year >= 2020 { // Law 21.394
		date := time.Date(year, 6, 7, 0, 0, 0, 0, time.UTC)
		holidays[date] = cl.CreateHoliday(
			"Battle of Arica",
			date,
			"regional",
			map[string]string{
				"en": "Battle of Arica",
				"es": "Asalto y Toma del Morro de Arica",
			},
		)
	}

	// Ñuble region
	if year >= 2019 { // Law 21.148
		date := time.Date(year, 8, 20, 0, 0, 0, 0, time.UTC)
		holidays[date] = cl.CreateHoliday(
			"Chillán Foundation Day",
			date,
			"regional",
			map[string]string{
				"en": "Chillán Foundation Day",
				"es": "Nacimiento de Chillán",
			},
		)
	}
}

// getObservedDate returns the observed date for holidays that can be moved to Monday
func (cl *CLProvider) getObservedDate(date time.Time) time.Time {
	weekday := date.Weekday()

	switch weekday {
	case time.Tuesday:
		// Move to previous Monday
		return date.AddDate(0, 0, -1)
	case time.Wednesday:
		// Move to previous Monday
		return date.AddDate(0, 0, -2)
	case time.Thursday:
		// Move to previous Monday
		return date.AddDate(0, 0, -3)
	case time.Friday:
		// Move to previous Monday
		return date.AddDate(0, 0, -4)
	case time.Saturday:
		// Move to next Monday
		return date.AddDate(0, 0, 2)
	case time.Sunday:
		// Move to next Monday
		return date.AddDate(0, 0, 1)
	default:
		// Already Monday, keep as is
		return date
	}
}
