package countries

import (
	"time"
)

// ILProvider implements holiday logic for Israel
type ILProvider struct {
	BaseProvider
}

// NewILProvider creates a new Israel holiday provider
func NewILProvider() *ILProvider {
	return &ILProvider{
		BaseProvider: *NewBaseProvider("IL"),
	}
}

// GetCountryCode returns the country code
func (il *ILProvider) GetCountryCode() string {
	return il.BaseProvider.GetCountryCode()
}

// GetCountryName returns the country name
func (il *ILProvider) GetCountryName() string {
	return "Israel"
}

// GetSubdivisions returns Israeli districts
func (il *ILProvider) GetSubdivisions() []string {
	return []string{
		"D",  // Southern District
		"HA", // Haifa District
		"JM", // Jerusalem District
		"M",  // Central District
		"TA", // Tel Aviv District
		"Z",  // Northern District
	}
}

// GetCategories returns holiday categories used in Israel
func (il *ILProvider) GetCategories() []string {
	return []string{"public", "religious", "memorial", "national", "jewish"}
}

// LoadHolidays loads Israeli holidays for the specified year
func (il *ILProvider) LoadHolidays(year int) map[time.Time]*Holiday {
	holidays := make(map[time.Time]*Holiday)

	// Hebrew calendar holidays (using known dates)
	il.addHebrewCalendarHolidays(holidays, year)

	// Memorial and national days
	il.addMemorialDays(holidays, year)

	// Independence Day (depends on Hebrew calendar)
	il.addIndependenceDay(holidays, year)

	return holidays
}

// addHebrewCalendarHolidays adds holidays based on the Hebrew calendar
func (il *ILProvider) addHebrewCalendarHolidays(holidays map[time.Time]*Holiday, year int) {
	// Use known accurate dates for Hebrew calendar holidays
	// These dates are calculated based on the Hebrew calendar and vary each year

	hebrewHolidays := il.getHebrewHolidayDates(year)

	for date, holidayInfo := range hebrewHolidays {
		holidays[date] = il.CreateHoliday(
			holidayInfo.name,
			date,
			holidayInfo.category,
			map[string]string{
				"en": holidayInfo.name,
				"he": holidayInfo.nameHE,
			},
		)
	}
}

// addMemorialDays adds memorial and remembrance days
func (il *ILProvider) addMemorialDays(holidays map[time.Time]*Holiday, year int) {
	// Holocaust Remembrance Day (Yom HaShoah) - 27th of Nisan
	// Falls in April/May, dates vary by Hebrew calendar
	yomHaShoah := il.getYomHaShoahDate(year)
	if !yomHaShoah.IsZero() {
		holidays[yomHaShoah] = il.CreateHoliday(
			"Holocaust Remembrance Day",
			yomHaShoah,
			"memorial",
			map[string]string{
				"en": "Holocaust Remembrance Day",
				"he": "יום השואה",
			},
		)
	}

	// Memorial Day (Yom HaZikaron) - 4th of Iyar (day before Independence Day)
	yomHaZikaron := il.getYomHaZikaronDate(year)
	if !yomHaZikaron.IsZero() {
		holidays[yomHaZikaron] = il.CreateHoliday(
			"Memorial Day",
			yomHaZikaron,
			"memorial",
			map[string]string{
				"en": "Memorial Day",
				"he": "יום הזיכרון",
			},
		)
	}
}

// addIndependenceDay adds Independence Day (Yom Ha'atzmaut)
func (il *ILProvider) addIndependenceDay(holidays map[time.Time]*Holiday, year int) {
	// Independence Day - 5th of Iyar (day after Memorial Day)
	independenceDay := il.getIndependenceDayDate(year)
	if !independenceDay.IsZero() {
		holidays[independenceDay] = il.CreateHoliday(
			"Independence Day",
			independenceDay,
			"national",
			map[string]string{
				"en": "Independence Day",
				"he": "יום העצמאות",
			},
		)
	}
}

// HebrewHolidayInfo holds information about a Hebrew calendar holiday
type HebrewHolidayInfo struct {
	name     string
	nameHE   string
	category string
}

// getHebrewHolidayDates returns known dates for Hebrew calendar holidays
func (il *ILProvider) getHebrewHolidayDates(year int) map[time.Time]HebrewHolidayInfo {
	holidays := make(map[time.Time]HebrewHolidayInfo)

	// Known dates for major Jewish holidays
	// These are calculated based on the Hebrew calendar
	switch year {
	case 2023:
		holidays[time.Date(2023, 9, 16, 0, 0, 0, 0, time.UTC)] = HebrewHolidayInfo{"Rosh Hashanah", "ראש השנה", "religious"}
		holidays[time.Date(2023, 9, 17, 0, 0, 0, 0, time.UTC)] = HebrewHolidayInfo{"Rosh Hashanah (Day 2)", "ראש השנה יום ב'", "religious"}
		holidays[time.Date(2023, 9, 25, 0, 0, 0, 0, time.UTC)] = HebrewHolidayInfo{"Yom Kippur", "יום כיפור", "religious"}
		holidays[time.Date(2023, 9, 30, 0, 0, 0, 0, time.UTC)] = HebrewHolidayInfo{"Sukkot", "סוכות", "religious"}
		holidays[time.Date(2023, 10, 7, 0, 0, 0, 0, time.UTC)] = HebrewHolidayInfo{"Simchat Torah", "שמחת תורה", "religious"}
		holidays[time.Date(2023, 12, 8, 0, 0, 0, 0, time.UTC)] = HebrewHolidayInfo{"Hanukkah", "חנוכה", "religious"}
		holidays[time.Date(2024, 4, 23, 0, 0, 0, 0, time.UTC)] = HebrewHolidayInfo{"Passover", "פסח", "religious"}
		holidays[time.Date(2024, 5, 14, 0, 0, 0, 0, time.UTC)] = HebrewHolidayInfo{"Shavuot", "שבועות", "religious"}

	case 2024:
		holidays[time.Date(2024, 10, 3, 0, 0, 0, 0, time.UTC)] = HebrewHolidayInfo{"Rosh Hashanah", "ראש השנה", "religious"}
		holidays[time.Date(2024, 10, 4, 0, 0, 0, 0, time.UTC)] = HebrewHolidayInfo{"Rosh Hashanah (Day 2)", "ראש השנה יום ב'", "religious"}
		holidays[time.Date(2024, 10, 12, 0, 0, 0, 0, time.UTC)] = HebrewHolidayInfo{"Yom Kippur", "יום כיפור", "religious"}
		holidays[time.Date(2024, 10, 17, 0, 0, 0, 0, time.UTC)] = HebrewHolidayInfo{"Sukkot", "סוכות", "religious"}
		holidays[time.Date(2024, 10, 24, 0, 0, 0, 0, time.UTC)] = HebrewHolidayInfo{"Simchat Torah", "שמחת תורה", "religious"}
		holidays[time.Date(2024, 12, 26, 0, 0, 0, 0, time.UTC)] = HebrewHolidayInfo{"Hanukkah", "חנוכה", "religious"}
		holidays[time.Date(2024, 4, 23, 0, 0, 0, 0, time.UTC)] = HebrewHolidayInfo{"Passover", "פסח", "religious"}
		holidays[time.Date(2024, 6, 12, 0, 0, 0, 0, time.UTC)] = HebrewHolidayInfo{"Shavuot", "שבועות", "religious"}

	case 2025:
		holidays[time.Date(2025, 9, 23, 0, 0, 0, 0, time.UTC)] = HebrewHolidayInfo{"Rosh Hashanah", "ראש השנה", "religious"}
		holidays[time.Date(2025, 9, 24, 0, 0, 0, 0, time.UTC)] = HebrewHolidayInfo{"Rosh Hashanah (Day 2)", "ראש השנה יום ב'", "religious"}
		holidays[time.Date(2025, 10, 2, 0, 0, 0, 0, time.UTC)] = HebrewHolidayInfo{"Yom Kippur", "יום כיפור", "religious"}
		holidays[time.Date(2025, 10, 7, 0, 0, 0, 0, time.UTC)] = HebrewHolidayInfo{"Sukkot", "סוכות", "religious"}
		holidays[time.Date(2025, 10, 14, 0, 0, 0, 0, time.UTC)] = HebrewHolidayInfo{"Simchat Torah", "שמחת תורה", "religious"}
		holidays[time.Date(2025, 12, 15, 0, 0, 0, 0, time.UTC)] = HebrewHolidayInfo{"Hanukkah", "חנוכה", "religious"}
		holidays[time.Date(2025, 4, 13, 0, 0, 0, 0, time.UTC)] = HebrewHolidayInfo{"Passover", "פסח", "religious"}
		holidays[time.Date(2025, 6, 2, 0, 0, 0, 0, time.UTC)] = HebrewHolidayInfo{"Shavuot", "שבועות", "religious"}
	}

	return holidays
}

// getYomHaShoahDate returns Holocaust Remembrance Day date for the given year
func (il *ILProvider) getYomHaShoahDate(year int) time.Time {
	// 27th of Nisan - known dates
	switch year {
	case 2023:
		return time.Date(2023, 4, 18, 0, 0, 0, 0, time.UTC)
	case 2024:
		return time.Date(2024, 5, 6, 0, 0, 0, 0, time.UTC)
	case 2025:
		return time.Date(2025, 4, 24, 0, 0, 0, 0, time.UTC)
	case 2026:
		return time.Date(2026, 4, 14, 0, 0, 0, 0, time.UTC)
	default:
		// For years outside our known range, return zero time
		// In production, this would use a Hebrew calendar library
		return time.Time{}
	}
}

// getYomHaZikaronDate returns Memorial Day date for the given year
func (il *ILProvider) getYomHaZikaronDate(year int) time.Time {
	// 4th of Iyar - day before Independence Day
	switch year {
	case 2023:
		return time.Date(2023, 4, 25, 0, 0, 0, 0, time.UTC)
	case 2024:
		return time.Date(2024, 5, 13, 0, 0, 0, 0, time.UTC)
	case 2025:
		return time.Date(2025, 5, 1, 0, 0, 0, 0, time.UTC)
	case 2026:
		return time.Date(2026, 4, 21, 0, 0, 0, 0, time.UTC)
	default:
		return time.Time{}
	}
}

// getIndependenceDayDate returns Independence Day date for the given year
func (il *ILProvider) getIndependenceDayDate(year int) time.Time {
	// 5th of Iyar - day after Memorial Day
	switch year {
	case 2023:
		return time.Date(2023, 4, 26, 0, 0, 0, 0, time.UTC)
	case 2024:
		return time.Date(2024, 5, 14, 0, 0, 0, 0, time.UTC)
	case 2025:
		return time.Date(2025, 5, 2, 0, 0, 0, 0, time.UTC)
	case 2026:
		return time.Date(2026, 4, 22, 0, 0, 0, 0, time.UTC)
	default:
		return time.Time{}
	}
}
