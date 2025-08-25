package updater

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/coredds/GoHoliday"
)

// GitHubSyncer handles real GitHub API integration for Python holidays sync
type GitHubSyncer struct {
	client      *http.Client
	baseURL     string
	repoOwner   string
	repoName    string
	branch      string
	rateLimiter chan struct{}
}

// NewGitHubSyncer creates a new GitHub API syncer
func NewGitHubSyncer() *GitHubSyncer {
	// Rate limiter: GitHub allows 60 requests/hour for unauthenticated requests
	// We'll be conservative and limit to 1 request per second
	rateLimiter := make(chan struct{}, 1)
	go func() {
		for {
			rateLimiter <- struct{}{}
			time.Sleep(1 * time.Second)
		}
	}()

	return &GitHubSyncer{
		client:      &http.Client{Timeout: 30 * time.Second},
		baseURL:     "https://api.github.com",
		repoOwner:   "vacanza",
		repoName:    "holidays",
		branch:      "dev", // Python holidays uses 'dev' as main branch
		rateLimiter: rateLimiter,
	}
}

// GitHubFile represents a file from GitHub API
type GitHubFile struct {
	Name        string `json:"name"`
	Path        string `json:"path"`
	SHA         string `json:"sha"`
	Size        int    `json:"size"`
	URL         string `json:"url"`
	HTMLURL     string `json:"html_url"`
	GitURL      string `json:"git_url"`
	DownloadURL string `json:"download_url"`
	Type        string `json:"type"`
}

// GitHubContent represents file content from GitHub API
type GitHubContent struct {
	Name        string `json:"name"`
	Path        string `json:"path"`
	SHA         string `json:"sha"`
	Size        int    `json:"size"`
	URL         string `json:"url"`
	HTMLURL     string `json:"html_url"`
	GitURL      string `json:"git_url"`
	DownloadURL string `json:"download_url"`
	Type        string `json:"type"`
	Content     string `json:"content"`
	Encoding    string `json:"encoding"`
}

// FetchCountryList retrieves the list of available country modules
func (gs *GitHubSyncer) FetchCountryList(ctx context.Context) ([]string, error) {
	<-gs.rateLimiter // Rate limiting

	url := fmt.Sprintf("%s/repos/%s/%s/contents/holidays/countries",
		gs.baseURL, gs.repoOwner, gs.repoName)

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Accept", "application/vnd.github.v3+json")
	req.Header.Set("User-Agent", fmt.Sprintf("GoHolidays/%s", goholidays.Version))

	resp, err := gs.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch country list: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("GitHub API error %d: %s", resp.StatusCode, string(body))
	}

	var files []GitHubFile
	if err := json.NewDecoder(resp.Body).Decode(&files); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	var countries []string
	for _, file := range files {
		if file.Type == "file" && strings.HasSuffix(file.Name, ".py") && file.Name != "__init__.py" {
			// Extract country code from filename (e.g., "united_states.py" -> "US")
			countryCode := gs.extractCountryCode(file.Name)
			if countryCode != "" {
				countries = append(countries, countryCode)
			}
		}
	}

	return countries, nil
}

// FetchCountryFile retrieves the Python source file for a specific country
func (gs *GitHubSyncer) FetchCountryFile(ctx context.Context, countryCode string) (string, error) {
	<-gs.rateLimiter // Rate limiting

	filename := gs.getCountryFilename(countryCode)
	url := fmt.Sprintf("%s/repos/%s/%s/contents/holidays/countries/%s",
		gs.baseURL, gs.repoOwner, gs.repoName, filename)

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Accept", "application/vnd.github.v3+json")
	req.Header.Set("User-Agent", fmt.Sprintf("GoHolidays/%s", goholidays.Version))

	resp, err := gs.client.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to fetch country file: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("GitHub API error %d: %s", resp.StatusCode, string(body))
	}

	var content GitHubContent
	if err := json.NewDecoder(resp.Body).Decode(&content); err != nil {
		return "", fmt.Errorf("failed to decode response: %w", err)
	}

	if content.Encoding != "base64" {
		return "", fmt.Errorf("unexpected encoding: %s", content.Encoding)
	}

	// Decode base64 content
	decoded, err := decodeBase64Content(content.Content)
	if err != nil {
		return "", fmt.Errorf("failed to decode content: %w", err)
	}

	return decoded, nil
}

// ParseHolidayDefinitions extracts holiday definitions from Python source code
func (gs *GitHubSyncer) ParseHolidayDefinitions(pythonSource string) (*CountryData, error) {
	countryData := &CountryData{
		Holidays:   make(map[string]HolidayDefinition),
		Categories: []string{"public"},
		Languages:  []string{"en"},
		UpdatedAt:  time.Now(),
	}

	// Extract country information from class definition
	classRegex := regexp.MustCompile(`class\s+(\w+)\s*\([^)]*\):`)
	if matches := classRegex.FindStringSubmatch(pythonSource); len(matches) > 1 {
		countryData.Name = gs.convertClassName(matches[1])
		countryData.CountryCode = gs.extractCountryCodeFromClass(matches[1])
	}

	// Extract subdivisions
	countryData.Subdivisions = gs.extractSubdivisions(pythonSource)

	// Use the new Python AST parser for better accuracy
	astParser := NewPythonASTParser(pythonSource)
	holidayCalls, err := astParser.Parse()
	if err != nil {
		// Fallback to old regex method if AST parsing fails
		holidays := gs.extractHolidays(pythonSource)
		for name, holiday := range holidays {
			countryData.Holidays[name] = holiday
		}
	} else {
		// Convert AST results to holiday definitions
		astHolidays := astParser.ConvertToHolidayDefinitions(holidayCalls)
		for name, holiday := range astHolidays {
			countryData.Holidays[name] = holiday
		}
	}

	return countryData, nil
}

// extractCountryCode converts filename to country code
func (gs *GitHubSyncer) extractCountryCode(filename string) string {
	// Remove .py extension
	name := strings.TrimSuffix(filename, ".py")

	// Skip special files
	if name == "__init__" || !strings.HasSuffix(filename, ".py") {
		return ""
	}

	// Map common filename patterns to ISO country codes
	countryMap := map[string]string{
		"united_states":        "US",
		"united_kingdom":       "GB",
		"canada":               "CA",
		"australia":            "AU",
		"germany":              "DE",
		"france":               "FR",
		"italy":                "IT",
		"spain":                "ES",
		"netherlands":          "NL",
		"belgium":              "BE",
		"austria":              "AT",
		"switzerland":          "CH",
		"sweden":               "SE",
		"norway":               "NO",
		"denmark":              "DK",
		"finland":              "FI",
		"poland":               "PL",
		"czechia":              "CZ",
		"hungary":              "HU",
		"portugal":             "PT",
		"greece":               "GR",
		"ireland":              "IE",
		"luxembourg":           "LU",
		"slovakia":             "SK",
		"slovenia":             "SI",
		"estonia":              "EE",
		"latvia":               "LV",
		"lithuania":            "LT",
		"romania":              "RO",
		"bulgaria":             "BG",
		"croatia":              "HR",
		"malta":                "MT",
		"cyprus":               "CY",
		"japan":                "JP",
		"south_korea":          "KR",
		"china":                "CN",
		"india":                "IN",
		"singapore":            "SG",
		"hong_kong":            "HK",
		"taiwan":               "TW",
		"thailand":             "TH",
		"malaysia":             "MY",
		"indonesia":            "ID",
		"philippines":          "PH",
		"vietnam":              "VN",
		"new_zealand":          "NZ",
		"south_africa":         "ZA",
		"brazil":               "BR",
		"argentina":            "AR",
		"chile":                "CL",
		"mexico":               "MX",
		"colombia":             "CO",
		"peru":                 "PE",
		"venezuela":            "VE",
		"ecuador":              "EC",
		"bolivia":              "BO",
		"paraguay":             "PY",
		"uruguay":              "UY",
		"russia":               "RU",
		"ukraine":              "UA",
		"belarus":              "BY",
		"israel":               "IL",
		"turkey":               "TR",
		"egypt":                "EG",
		"saudi_arabia":         "SA",
		"united_arab_emirates": "AE",
		"qatar":                "QA",
		"kuwait":               "KW",
		"bahrain":              "BH",
		"oman":                 "OM",
		"jordan":               "JO",
		"lebanon":              "LB",
		"iran":                 "IR",
		"iraq":                 "IQ",
		"pakistan":             "PK",
		"bangladesh":           "BD",
		"sri_lanka":            "LK",
		"nepal":                "NP",
		"myanmar":              "MM",
		"cambodia":             "KH",
		"laos":                 "LA",
		"mongolia":             "MN",
	}

	if code, exists := countryMap[name]; exists {
		return code
	}

	// For unknown countries, return empty string to indicate we can't map them
	return ""
}

// getCountryFilename converts country code to expected filename
func (gs *GitHubSyncer) getCountryFilename(countryCode string) string {
	filenameMap := map[string]string{
		"US": "united_states.py",
		"GB": "united_kingdom.py",
		"CA": "canada.py",
		"AU": "australia.py",
		"DE": "germany.py",
		"FR": "france.py",
		"NZ": "new_zealand.py",
		"ZA": "south_africa.py",
		// Add more mappings as needed
	}

	if filename, exists := filenameMap[countryCode]; exists {
		return filename
	}

	// Fallback: lowercase country code + .py
	return strings.ToLower(countryCode) + ".py"
}

// Helper functions for parsing Python source code
func (gs *GitHubSyncer) convertClassName(className string) string {
	// Convert "UnitedStates" to "United States"
	re := regexp.MustCompile(`([a-z])([A-Z])`)
	return re.ReplaceAllString(className, "$1 $2")
}

func (gs *GitHubSyncer) extractCountryCodeFromClass(className string) string {
	// This is a simplified extraction - in practice you'd need more logic
	countryMap := map[string]string{
		"UnitedStates":  "US",
		"UnitedKingdom": "GB",
		"Canada":        "CA",
		"Australia":     "AU",
		"Germany":       "DE",
		"France":        "FR",
	}

	if code, exists := countryMap[className]; exists {
		return code
	}

	return ""
}

func (gs *GitHubSyncer) extractSubdivisions(source string) map[string]string {
	subdivisions := make(map[string]string)

	// Look for subdivision definitions in the Python source
	// This is a simplified parser - real implementation would need AST parsing
	subdivisionRegex := regexp.MustCompile(`['"']([A-Z]{2,3})['"]:\s*['"]([^'"]+)['"]`)
	matches := subdivisionRegex.FindAllStringSubmatch(source, -1)

	for _, match := range matches {
		if len(match) >= 3 {
			subdivisions[match[1]] = match[2]
		}
	}

	return subdivisions
}

func (gs *GitHubSyncer) extractHolidays(source string) map[string]HolidayDefinition {
	holidays := make(map[string]HolidayDefinition)

	// Extract fixed date holidays (simplified pattern)
	fixedDateRegex := regexp.MustCompile(`self\._add_holiday[^(]*\([^,]+,\s*(\w+)\s*,\s*(\d+)\s*,\s*["']([^"']+)["']`)
	matches := fixedDateRegex.FindAllStringSubmatch(source, -1)

	for _, match := range matches {
		if len(match) >= 4 {
			month, _ := strconv.Atoi(match[2])
			day, _ := strconv.Atoi(match[1])
			name := match[3]

			holidayKey := strings.ToLower(strings.ReplaceAll(name, " ", "_"))
			holidays[holidayKey] = HolidayDefinition{
				Name:        name,
				Category:    "public",
				Calculation: "fixed",
				Month:       month,
				Day:         day,
				Languages:   map[string]string{"en": name},
			}
		}
	}

	// Extract more complex patterns (Easter-based, weekday-based) would go here
	// This is a simplified implementation for demonstration

	return holidays
}

// decodeBase64Content decodes base64 content from GitHub API
func decodeBase64Content(content string) (string, error) {
	// Remove whitespace and newlines
	content = strings.ReplaceAll(content, "\n", "")
	content = strings.ReplaceAll(content, " ", "")
	content = strings.ReplaceAll(content, "\t", "")

	// Use Go's base64 package
	decoded, err := base64.StdEncoding.DecodeString(content)
	if err != nil {
		return "", fmt.Errorf("failed to decode base64: %w", err)
	}

	return string(decoded), nil
}

// ParseWithComparison compares AST parsing with regex parsing for validation
func (gs *GitHubSyncer) ParseWithComparison(pythonSource string) (*CountryData, *ParsingComparison, error) {
	countryData := &CountryData{
		CountryCode:  "XX", // Will be set from context
		Name:         "Unknown",
		Subdivisions: make(map[string]string),
		Holidays:     make(map[string]HolidayDefinition),
	}

	// Parse with AST parser
	astParser := NewPythonASTParser(pythonSource)
	holidayCalls, astErr := astParser.Parse()
	astHolidays := make(map[string]HolidayDefinition)
	if astErr == nil {
		astHolidays = astParser.ConvertToHolidayDefinitions(holidayCalls)
	}

	// Parse with regex method
	regexHolidays := gs.extractHolidays(pythonSource)

	// Create comparison
	comparison := &ParsingComparison{
		Source:        "Python AST vs Regex",
		ASTHolidays:   len(astHolidays),
		RegexHolidays: len(regexHolidays),
		ASTError:      astErr,
		Differences:   make([]string, 0),
	}

	// Compare results
	astKeys := make(map[string]bool)
	for key := range astHolidays {
		astKeys[key] = true
	}

	regexKeys := make(map[string]bool)
	for key := range regexHolidays {
		regexKeys[key] = true
	}

	// Find holidays only in AST
	for key := range astKeys {
		if !regexKeys[key] {
			comparison.Differences = append(comparison.Differences, fmt.Sprintf("AST only: %s", key))
		}
	}

	// Find holidays only in regex
	for key := range regexKeys {
		if !astKeys[key] {
			comparison.Differences = append(comparison.Differences, fmt.Sprintf("Regex only: %s", key))
		}
	}

	// Use AST results if available, otherwise fallback to regex
	if astErr == nil && len(astHolidays) > 0 {
		countryData.Holidays = astHolidays
		comparison.PreferredMethod = "AST"
	} else {
		countryData.Holidays = regexHolidays
		comparison.PreferredMethod = "Regex"
	}

	return countryData, comparison, nil
}

// ParsingComparison provides comparison data between parsing methods
type ParsingComparison struct {
	Source          string
	ASTHolidays     int
	RegexHolidays   int
	ASTError        error
	PreferredMethod string
	Differences     []string
}

// Validation helper to check if decoded content is valid Python
func (gs *GitHubSyncer) ValidatePythonContent(content string) error {
	// Basic validation checks
	if !strings.Contains(content, "class") {
		return fmt.Errorf("no class definition found")
	}

	if !strings.Contains(content, "_add_holiday") {
		return fmt.Errorf("no holiday definitions found")
	}

	return nil
}
