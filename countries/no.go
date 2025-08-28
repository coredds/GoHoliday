package countries

import (
	"time"
)

// NOProvider implements holiday calculations for Norway
type NOProvider struct {
	*BaseProvider
}

// NewNOProvider creates a new Norwegian holiday provider
func NewNOProvider() *NOProvider {
	base := NewBaseProvider("NO")

	// Norway has 11 counties (fylker) after 2020 reform
	base.subdivisions = []string{
		"NO-03", // Oslo
		"NO-11", // Rogaland
		"NO-15", // Møre og Romsdal
		"NO-18", // Nordland
		"NO-30", // Viken
		"NO-34", // Innlandet
		"NO-38", // Vestfold og Telemark
		"NO-42", // Agder
		"NO-46", // Vestland
		"NO-50", // Trøndelag
		"NO-54", // Troms og Finnmark
	}

	base.categories = []string{"national", "religious", "traditional", "royal"}

	return &NOProvider{BaseProvider: base}
}

// LoadHolidays loads all Norwegian holidays for a given year
func (no *NOProvider) LoadHolidays(year int) map[time.Time]*Holiday {
	holidays := make(map[time.Time]*Holiday)

	// New Year's Day - January 1
	holidays[time.Date(year, 1, 1, 0, 0, 0, 0, time.UTC)] = no.CreateHoliday(
		"Nyttårsdag", time.Date(year, 1, 1, 0, 0, 0, 0, time.UTC), "national",
		map[string]string{
			"no": "Nyttårsdag",
			"en": "New Year's Day",
		},
	)

	// Maundy Thursday (Thursday before Easter)
	easter := EasterSunday(year)
	maundyThursday := easter.AddDate(0, 0, -3)
	holidays[maundyThursday] = no.CreateHoliday(
		"Skjærtorsdag", maundyThursday, "religious",
		map[string]string{
			"no": "Skjærtorsdag",
			"en": "Maundy Thursday",
		},
	)

	// Good Friday (Friday before Easter)
	goodFriday := easter.AddDate(0, 0, -2)
	holidays[goodFriday] = no.CreateHoliday(
		"Langfredag", goodFriday, "religious",
		map[string]string{
			"no": "Langfredag",
			"en": "Good Friday",
		},
	)

	// Easter Sunday
	holidays[easter] = no.CreateHoliday(
		"Første påskedag", easter, "religious",
		map[string]string{
			"no": "Første påskedag",
			"en": "Easter Sunday",
		},
	)

	// Easter Monday (day after Easter)
	easterMonday := easter.AddDate(0, 0, 1)
	holidays[easterMonday] = no.CreateHoliday(
		"Andre påskedag", easterMonday, "religious",
		map[string]string{
			"no": "Andre påskedag",
			"en": "Easter Monday",
		},
	)

	// Labour Day - May 1
	holidays[time.Date(year, 5, 1, 0, 0, 0, 0, time.UTC)] = no.CreateHoliday(
		"Arbeidernes dag", time.Date(year, 5, 1, 0, 0, 0, 0, time.UTC), "national",
		map[string]string{
			"no": "Arbeidernes dag",
			"en": "Labour Day",
		},
	)

	// Constitution Day - May 17
	holidays[time.Date(year, 5, 17, 0, 0, 0, 0, time.UTC)] = no.CreateHoliday(
		"Grunnlovsdag", time.Date(year, 5, 17, 0, 0, 0, 0, time.UTC), "national",
		map[string]string{
			"no": "Grunnlovsdag",
			"en": "Constitution Day",
		},
	)

	// Ascension Day (39 days after Easter)
	ascensionDay := easter.AddDate(0, 0, 39)
	holidays[ascensionDay] = no.CreateHoliday(
		"Kristi himmelfartsdag", ascensionDay, "religious",
		map[string]string{
			"no": "Kristi himmelfartsdag",
			"en": "Ascension Day",
		},
	)

	// Whit Sunday (49 days after Easter)
	whitSunday := easter.AddDate(0, 0, 49)
	holidays[whitSunday] = no.CreateHoliday(
		"Første pinsedag", whitSunday, "religious",
		map[string]string{
			"no": "Første pinsedag",
			"en": "Whit Sunday",
		},
	)

	// Whit Monday (50 days after Easter)
	whitMonday := easter.AddDate(0, 0, 50)
	holidays[whitMonday] = no.CreateHoliday(
		"Andre pinsedag", whitMonday, "religious",
		map[string]string{
			"no": "Andre pinsedag",
			"en": "Whit Monday",
		},
	)

	// Christmas Eve - December 24
	holidays[time.Date(year, 12, 24, 0, 0, 0, 0, time.UTC)] = no.CreateHoliday(
		"Julaften", time.Date(year, 12, 24, 0, 0, 0, 0, time.UTC), "traditional",
		map[string]string{
			"no": "Julaften",
			"en": "Christmas Eve",
		},
	)

	// Christmas Day - December 25
	holidays[time.Date(year, 12, 25, 0, 0, 0, 0, time.UTC)] = no.CreateHoliday(
		"Første juledag", time.Date(year, 12, 25, 0, 0, 0, 0, time.UTC), "religious",
		map[string]string{
			"no": "Første juledag",
			"en": "Christmas Day",
		},
	)

	// Boxing Day - December 26
	holidays[time.Date(year, 12, 26, 0, 0, 0, 0, time.UTC)] = no.CreateHoliday(
		"Andre juledag", time.Date(year, 12, 26, 0, 0, 0, 0, time.UTC), "traditional",
		map[string]string{
			"no": "Andre juledag",
			"en": "Boxing Day",
		},
	)

	// New Year's Eve - December 31
	holidays[time.Date(year, 12, 31, 0, 0, 0, 0, time.UTC)] = no.CreateHoliday(
		"Nyttårsaften", time.Date(year, 12, 31, 0, 0, 0, 0, time.UTC), "traditional",
		map[string]string{
			"no": "Nyttårsaften",
			"en": "New Year's Eve",
		},
	)

	return holidays
}
