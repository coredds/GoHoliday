package countries

import (
	"time"
)

// USProvider implements holiday calculations for the United States
type USProvider struct {
	*BaseProvider
}

// NewUSProvider creates a new US holiday provider
func NewUSProvider() *USProvider {
	base := NewBaseProvider("US")
	base.subdivisions = []string{
		"AL", "AK", "AZ", "AR", "CA", "CO", "CT", "DE", "FL", "GA",
		"HI", "ID", "IL", "IN", "IA", "KS", "KY", "LA", "ME", "MD",
		"MA", "MI", "MN", "MS", "MO", "MT", "NE", "NV", "NH", "NJ",
		"NM", "NY", "NC", "ND", "OH", "OK", "OR", "PA", "RI", "SC",
		"SD", "TN", "TX", "UT", "VT", "VA", "WA", "WV", "WI", "WY",
		"DC", "AS", "GU", "MP", "PR", "VI",
	}
	base.categories = []string{"federal", "state", "religious", "observance"}
	
	return &USProvider{BaseProvider: base}
}

// LoadHolidays loads all US holidays for a given year
func (us *USProvider) LoadHolidays(year int) map[time.Time]*Holiday {
	holidays := make(map[time.Time]*Holiday)
	
	// Fixed date holidays
	holidays[time.Date(year, 1, 1, 0, 0, 0, 0, time.UTC)] = us.CreateHoliday(
		"New Year's Day",
		time.Date(year, 1, 1, 0, 0, 0, 0, time.UTC),
		"federal",
		map[string]string{
			"en": "New Year's Day",
			"es": "Año Nuevo",
		},
	)
	
	// Juneteenth - June 19 (federal holiday since 2021)
	if year >= 2021 {
		juneteenth := time.Date(year, 6, 19, 0, 0, 0, 0, time.UTC)
		holidays[juneteenth] = us.CreateHoliday(
			"Juneteenth",
			juneteenth,
			"federal",
			map[string]string{
				"en": "Juneteenth",
				"es": "Juneteenth",
			},
		)
	}
	
	holidays[time.Date(year, 7, 4, 0, 0, 0, 0, time.UTC)] = us.CreateHoliday(
		"Independence Day",
		time.Date(year, 7, 4, 0, 0, 0, 0, time.UTC),
		"federal",
		map[string]string{
			"en": "Independence Day",
			"es": "Día de la Independencia",
		},
	)
	
	holidays[time.Date(year, 11, 11, 0, 0, 0, 0, time.UTC)] = us.CreateHoliday(
		"Veterans Day",
		time.Date(year, 11, 11, 0, 0, 0, 0, time.UTC),
		"federal",
		map[string]string{
			"en": "Veterans Day",
			"es": "Día de los Veteranos",
		},
	)
	
	holidays[time.Date(year, 12, 25, 0, 0, 0, 0, time.UTC)] = us.CreateHoliday(
		"Christmas Day",
		time.Date(year, 12, 25, 0, 0, 0, 0, time.UTC),
		"federal",
		map[string]string{
			"en": "Christmas Day",
			"es": "Navidad",
		},
	)
	
	// Variable date holidays
	
	// Martin Luther King Jr. Day - 3rd Monday in January (since 1983)
	if year >= 1983 {
		mlkDay := NthWeekdayOfMonth(year, 1, time.Monday, 3)
		holidays[mlkDay] = us.CreateHoliday(
			"Martin Luther King Jr. Day",
			mlkDay,
			"federal",
			map[string]string{
				"en": "Martin Luther King Jr. Day",
				"es": "Día de Martin Luther King Jr.",
			},
		)
	}
	
	// Presidents' Day - 3rd Monday in February
	presidentsDay := NthWeekdayOfMonth(year, 2, time.Monday, 3)
	holidays[presidentsDay] = us.CreateHoliday(
		"Presidents' Day",
		presidentsDay,
		"federal",
		map[string]string{
			"en": "Presidents' Day",
			"es": "Día de los Presidentes",
		},
	)
	
	// Memorial Day - Last Monday in May
	memorialDay := NthWeekdayOfMonth(year, 5, time.Monday, -1)
	holidays[memorialDay] = us.CreateHoliday(
		"Memorial Day",
		memorialDay,
		"federal",
		map[string]string{
			"en": "Memorial Day",
			"es": "Día de los Caídos",
		},
	)
	
	// Labor Day - 1st Monday in September
	laborDay := NthWeekdayOfMonth(year, 9, time.Monday, 1)
	holidays[laborDay] = us.CreateHoliday(
		"Labor Day",
		laborDay,
		"federal",
		map[string]string{
			"en": "Labor Day",
			"es": "Día del Trabajo",
		},
	)
	
	// Columbus Day - 2nd Monday in October
	columbusDay := NthWeekdayOfMonth(year, 10, time.Monday, 2)
	holidays[columbusDay] = us.CreateHoliday(
		"Columbus Day",
		columbusDay,
		"federal",
		map[string]string{
			"en": "Columbus Day",
			"es": "Día de Colón",
		},
	)
	
	// Thanksgiving Day - 4th Thursday in November
	thanksgiving := NthWeekdayOfMonth(year, 11, time.Thursday, 4)
	holidays[thanksgiving] = us.CreateHoliday(
		"Thanksgiving Day",
		thanksgiving,
		"federal",
		map[string]string{
			"en": "Thanksgiving Day",
			"es": "Día de Acción de Gracias",
		},
	)
	
	return holidays
}

// GetStateHolidays returns state-specific holidays for given subdivisions
func (us *USProvider) GetStateHolidays(year int, subdivisions []string) map[time.Time]*Holiday {
	holidays := make(map[time.Time]*Holiday)
	
	for _, state := range subdivisions {
		switch state {
		case "CA":
			// California-specific holidays
			// Cesar Chavez Day - March 31
			chavezDay := time.Date(year, 3, 31, 0, 0, 0, 0, time.UTC)
			holidays[chavezDay] = us.CreateHoliday(
				"Cesar Chavez Day",
				chavezDay,
				"public",
				map[string]string{
					"en": "Cesar Chavez Day",
					"es": "Día de César Chávez",
				},
			)
			
		case "TX":
			// Texas-specific holidays
			// Texas Independence Day - March 2
			texasIndependence := time.Date(year, 3, 2, 0, 0, 0, 0, time.UTC)
			holidays[texasIndependence] = us.CreateHoliday(
				"Texas Independence Day",
				texasIndependence,
				"public",
				map[string]string{
					"en": "Texas Independence Day",
					"es": "Día de la Independencia de Texas",
				},
			)
			
		case "MA":
			// Massachusetts-specific holidays
			// Patriots' Day - 3rd Monday in April
			patriotsDay := NthWeekdayOfMonth(year, 4, time.Monday, 3)
			holidays[patriotsDay] = us.CreateHoliday(
				"Patriots' Day",
				patriotsDay,
				"public",
				map[string]string{
					"en": "Patriots' Day",
				},
			)
		}
	}
	
	return holidays
}
