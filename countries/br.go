package countries

import (
	"time"
)

// BRProvider implements holiday calculations for Brazil
type BRProvider struct {
	*BaseProvider
}

// NewBRProvider creates a new Brazilian holiday provider
func NewBRProvider() *BRProvider {
	base := NewBaseProvider("BR")
	base.subdivisions = []string{
		// 26 states + 1 federal district
		"AC", "AL", "AP", "AM", "BA", "CE", "DF", "ES", "GO", "MA",
		"MT", "MS", "MG", "PA", "PB", "PR", "PE", "PI", "RJ", "RN",
		"RS", "RO", "RR", "SC", "SP", "SE", "TO",
	}
	base.categories = []string{"public", "national", "religious", "regional", "carnival"}

	return &BRProvider{BaseProvider: base}
}

// LoadHolidays loads all Brazilian holidays for a given year
func (br *BRProvider) LoadHolidays(year int) map[time.Time]*Holiday {
	holidays := make(map[time.Time]*Holiday)

	// Fixed national holidays

	// New Year's Day - January 1
	newYear := time.Date(year, 1, 1, 0, 0, 0, 0, time.UTC)
	holidays[newYear] = br.CreateHoliday(
		"Confraternização Universal",
		newYear,
		"public",
		map[string]string{
			"pt": "Confraternização Universal",
			"en": "New Year's Day",
		},
	)

	// Tiradentes - April 21
	tiradentes := time.Date(year, 4, 21, 0, 0, 0, 0, time.UTC)
	holidays[tiradentes] = br.CreateHoliday(
		"Tiradentes",
		tiradentes,
		"public",
		map[string]string{
			"pt": "Tiradentes",
			"en": "Tiradentes",
		},
	)

	// Labour Day - May 1
	labourDay := time.Date(year, 5, 1, 0, 0, 0, 0, time.UTC)
	holidays[labourDay] = br.CreateHoliday(
		"Dia do Trabalhador",
		labourDay,
		"public",
		map[string]string{
			"pt": "Dia do Trabalhador",
			"en": "Labour Day",
		},
	)

	// Independence Day - September 7
	independence := time.Date(year, 9, 7, 0, 0, 0, 0, time.UTC)
	holidays[independence] = br.CreateHoliday(
		"Independência do Brasil",
		independence,
		"public",
		map[string]string{
			"pt": "Independência do Brasil",
			"en": "Independence of Brazil",
		},
	)

	// Our Lady of Aparecida - October 12
	aparecida := time.Date(year, 10, 12, 0, 0, 0, 0, time.UTC)
	holidays[aparecida] = br.CreateHoliday(
		"Nossa Senhora Aparecida",
		aparecida,
		"religious",
		map[string]string{
			"pt": "Nossa Senhora Aparecida",
			"en": "Our Lady of Aparecida",
		},
	)

	// All Souls' Day - November 2
	allSouls := time.Date(year, 11, 2, 0, 0, 0, 0, time.UTC)
	holidays[allSouls] = br.CreateHoliday(
		"Finados",
		allSouls,
		"religious",
		map[string]string{
			"pt": "Finados",
			"en": "All Souls' Day",
		},
	)

	// Proclamation of the Republic - November 15
	republic := time.Date(year, 11, 15, 0, 0, 0, 0, time.UTC)
	holidays[republic] = br.CreateHoliday(
		"Proclamação da República",
		republic,
		"public",
		map[string]string{
			"pt": "Proclamação da República",
			"en": "Proclamation of the Republic",
		},
	)

	// Christmas Day - December 25
	christmas := time.Date(year, 12, 25, 0, 0, 0, 0, time.UTC)
	holidays[christmas] = br.CreateHoliday(
		"Natal",
		christmas,
		"religious",
		map[string]string{
			"pt": "Natal",
			"en": "Christmas Day",
		},
	)

	// Variable holidays based on Easter
	easter := EasterSunday(year)

	// Carnival Monday (48 days before Easter)
	carnivalMonday := easter.AddDate(0, 0, -48)
	holidays[carnivalMonday] = br.CreateHoliday(
		"Segunda-feira de Carnaval",
		carnivalMonday,
		"carnival",
		map[string]string{
			"pt": "Segunda-feira de Carnaval",
			"en": "Carnival Monday",
		},
	)

	// Carnival Tuesday (47 days before Easter)
	carnivalTuesday := easter.AddDate(0, 0, -47)
	holidays[carnivalTuesday] = br.CreateHoliday(
		"Terça-feira de Carnaval",
		carnivalTuesday,
		"carnival",
		map[string]string{
			"pt": "Terça-feira de Carnaval",
			"en": "Carnival Tuesday",
		},
	)

	// Good Friday (2 days before Easter)
	goodFriday := easter.AddDate(0, 0, -2)
	holidays[goodFriday] = br.CreateHoliday(
		"Sexta-feira Santa",
		goodFriday,
		"religious",
		map[string]string{
			"pt": "Sexta-feira Santa",
			"en": "Good Friday",
		},
	)

	// Corpus Christi (60 days after Easter)
	corpusChristi := easter.AddDate(0, 0, 60)
	holidays[corpusChristi] = br.CreateHoliday(
		"Corpus Christi",
		corpusChristi,
		"religious",
		map[string]string{
			"pt": "Corpus Christi",
			"en": "Corpus Christi",
		},
	)

	return holidays
}

// GetCountryCode returns the country code for Brazil
func (br *BRProvider) GetCountryCode() string {
	return "BR"
}

// GetSupportedSubdivisions returns the list of supported Brazilian states
func (br *BRProvider) GetSupportedSubdivisions() []string {
	return br.subdivisions
}

// GetSupportedCategories returns the list of supported holiday categories
func (br *BRProvider) GetSupportedCategories() []string {
	return br.categories
}

// GetName returns the country name
func (br *BRProvider) GetName() string {
	return "Brazil"
}

// GetLanguages returns the supported languages
func (br *BRProvider) GetLanguages() []string {
	return []string{"pt", "en"}
}

// IsSubdivisionSupported checks if a subdivision is supported
func (br *BRProvider) IsSubdivisionSupported(subdivision string) bool {
	for _, s := range br.subdivisions {
		if s == subdivision {
			return true
		}
	}
	return false
}

// IsCategorySupported checks if a category is supported
func (br *BRProvider) IsCategorySupported(category string) bool {
	for _, c := range br.categories {
		if c == category {
			return true
		}
	}
	return false
}
