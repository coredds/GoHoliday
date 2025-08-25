package updater

import (
	"fmt"
	"strings"
	"testing"
	"time"
)

func TestPythonASTParser_BasicParsing(t *testing.T) {
	source := `
class TestCountry(HolidayBase):
    """Test country holidays."""
    
    def __init__(self, **kwargs):
        super().__init__(**kwargs)
    
    def _populate(self, year):
        # New Year's Day
        self._add_holiday("New Year's Day", date(year, JAN, 1))
        
        # Easter-based holidays
        self._add_holiday("Good Friday", easter(year) - rd(days=2))
        self._add_holiday("Easter Monday", easter(year) + rd(days=1))
        
        # Fixed date with month constant
        self._add_holiday("Independence Day", date(year, JUL, 4))
`

	parser := NewPythonASTParser(source)
	holidayCalls, err := parser.Parse()
	if err != nil {
		t.Fatalf("Parse() failed: %v", err)
	}
	
	if len(holidayCalls) == 0 {
		t.Error("Expected to find holiday calls, got none")
	}
	
	// Check for specific holidays
	found := make(map[string]bool)
	for _, call := range holidayCalls {
		found[call.Name] = true
	}
	
	expectedHolidays := []string{"New Year's Day", "Good Friday", "Easter Monday", "Independence Day"}
	for _, expected := range expectedHolidays {
		if !found[expected] {
			t.Errorf("Expected to find holiday '%s', but it was not parsed", expected)
		}
	}
}

func TestPythonASTParser_TokenizeBasicContent(t *testing.T) {
	parser := NewPythonASTParser("")
	
	testCases := []struct {
		name     string
		content  string
		expected []TokenType
	}{
		{
			name:     "class definition",
			content:  "class TestCountry:",
			expected: []TokenType{TokenClass, TokenIdentifier, TokenOperator},
		},
		{
			name:     "method definition",
			content:  "def _populate(self, year):",
			expected: []TokenType{TokenDef, TokenIdentifier, TokenOperator, TokenSelf, TokenOperator, TokenIdentifier, TokenOperator, TokenOperator},
		},
		{
			name:     "holiday call",
			content:  `self._add_holiday("New Year", date(year, 1, 1))`,
			expected: []TokenType{TokenSelf, TokenOperator, TokenIdentifier, TokenOperator, TokenString, TokenOperator, TokenIdentifier, TokenOperator, TokenIdentifier, TokenOperator, TokenNumber, TokenOperator, TokenNumber, TokenOperator, TokenOperator},
		},
		{
			name:     "string literal",
			content:  `"Independence Day"`,
			expected: []TokenType{TokenString},
		},
		{
			name:     "number literal",
			content:  "42",
			expected: []TokenType{TokenNumber},
		},
	}
	
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			parser.tokens = []Token{} // Reset tokens
			parser.tokenizeContent(tc.content, 1, 0)
			
			if len(parser.tokens) != len(tc.expected) {
				t.Errorf("Expected %d tokens, got %d", len(tc.expected), len(parser.tokens))
				for i, token := range parser.tokens {
					t.Logf("Token %d: %s (%v)", i, token.Value, token.Type)
				}
				return
			}
			
			for i, expected := range tc.expected {
				if parser.tokens[i].Type != expected {
					t.Errorf("Token %d: expected type %v, got %v (value: %s)", 
						i, expected, parser.tokens[i].Type, parser.tokens[i].Value)
				}
			}
		})
	}
}

func TestPythonASTParser_ExtractDateExpression(t *testing.T) {
	parser := NewPythonASTParser("")
	
	testCases := []struct {
		name           string
		line           string
		expectedType   DateType
		expectedMonth  interface{}
		expectedDay    interface{}
		expectedCalc   string
	}{
		{
			name:          "fixed date with month constant",
			line:          `self._add_holiday("Test", date(year, JAN, 1))`,
			expectedType:  DateFixed,
			expectedMonth: "JAN",
			expectedDay:   1,
		},
		{
			name:          "fixed date with number",
			line:          `self._add_holiday("Test", date(year, 7, 4))`,
			expectedType:  DateFixed,
			expectedMonth: "7",
			expectedDay:   4,
		},
		{
			name:         "easter basic",
			line:         `self._add_holiday("Easter", easter(year))`,
			expectedType: DateEasterBased,
			expectedCalc: "easter(year)",
		},
		{
			name:         "easter plus days",
			line:         `self._add_holiday("Easter Monday", easter(year) + rd(days=1))`,
			expectedType: DateEasterBased,
			expectedCalc: "easter(year) + timedelta(days=1)",
		},
		{
			name:         "easter minus days",
			line:         `self._add_holiday("Good Friday", easter(year) - rd(days=2))`,
			expectedType: DateEasterBased,
			expectedCalc: "easter(year) - timedelta(days=2)",
		},
	}
	
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			dateExpr, err := parser.extractDateExpression(tc.line)
			if err != nil {
				t.Fatalf("extractDateExpression() failed: %v", err)
			}
			
			if dateExpr.Type != tc.expectedType {
				t.Errorf("Expected type %v, got %v", tc.expectedType, dateExpr.Type)
			}
			
			if tc.expectedMonth != nil && dateExpr.Month != tc.expectedMonth {
				t.Errorf("Expected month %v, got %v", tc.expectedMonth, dateExpr.Month)
			}
			
			if tc.expectedDay != nil && dateExpr.Day != tc.expectedDay {
				t.Errorf("Expected day %v, got %v", tc.expectedDay, dateExpr.Day)
			}
			
			if tc.expectedCalc != "" && dateExpr.Calculation != tc.expectedCalc {
				t.Errorf("Expected calculation %s, got %s", tc.expectedCalc, dateExpr.Calculation)
			}
		})
	}
}

func TestPythonASTParser_ParseHolidayCall(t *testing.T) {
	parser := NewPythonASTParser("")
	
	testCases := []struct {
		name           string
		line           string
		expectedMethod string
		expectedName   string
		expectedType   DateType
	}{
		{
			name:           "basic add_holiday",
			line:           `        self._add_holiday("New Year's Day", date(year, JAN, 1))`,
			expectedMethod: "_add_holiday",
			expectedName:   "New Year's Day",
			expectedType:   DateFixed,
		},
		{
			name:           "easter based holiday",
			line:           `        self._add_easter_based_holiday("Good Friday", -2)`,
			expectedMethod: "_add_easter_based_holiday",
			expectedName:   "Good Friday",
			expectedType:   DateEasterBased,
		},
		{
			name:           "new years method",
			line:           `        self._add_new_years_day("New Year's Day")`,
			expectedMethod: "_add_new_years_day",
			expectedName:   "New Year's Day",
			expectedType:   DateCalculated,
		},
		{
			name:           "christmas method",
			line:           `        self._add_christmas_day("Christmas Day")`,
			expectedMethod: "_add_christmas_day",
			expectedName:   "Christmas Day",
			expectedType:   DateCalculated,
		},
	}
	
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			holidayCall, err := parser.parseHolidayCall(tc.line, 1)
			if err != nil {
				t.Fatalf("parseHolidayCall() failed: %v", err)
			}
			
			if holidayCall.Method != tc.expectedMethod {
				t.Errorf("Expected method %s, got %s", tc.expectedMethod, holidayCall.Method)
			}
			
			if holidayCall.Name != tc.expectedName {
				t.Errorf("Expected name %s, got %s", tc.expectedName, holidayCall.Name)
			}
			
			if holidayCall.Date != nil && holidayCall.Date.Type != tc.expectedType {
				t.Errorf("Expected date type %v, got %v", tc.expectedType, holidayCall.Date.Type)
			}
		})
	}
}

func TestPythonASTParser_ConvertToHolidayDefinitions(t *testing.T) {
	parser := NewPythonASTParser("")
	
	holidayCalls := []HolidayCall{
		{
			Method: "_add_holiday",
			Name:   "New Year's Day",
			Date: &DateExpression{
				Type:  DateFixed,
				Year:  "year",
				Month: "JAN",
				Day:   1,
			},
			Line: 10,
		},
		{
			Method: "_add_easter_based_holiday",
			Name:   "Good Friday",
			Date: &DateExpression{
				Type:        DateEasterBased,
				Calculation: "easter(year) - timedelta(days=2)",
			},
			Line: 15,
		},
		{
			Method: "_add_weekday_holiday",
			Name:   "Labor Day",
			Date: &DateExpression{
				Type:        DateWeekdayBased,
				Calculation: "weekday_based",
			},
			Line: 20,
		},
	}
	
	definitions := parser.ConvertToHolidayDefinitions(holidayCalls)
	
	// Test New Year's Day
	if def, exists := definitions["new_year's_day"]; !exists {
		t.Error("Expected 'new_year's_day' definition not found")
	} else {
		if def.Name != "New Year's Day" {
			t.Errorf("Expected name 'New Year's Day', got '%s'", def.Name)
		}
		if def.Calculation != "fixed" {
			t.Errorf("Expected calculation 'fixed', got '%s'", def.Calculation)
		}
		if def.Month != 1 {
			t.Errorf("Expected month 1, got %d", def.Month)
		}
		if def.Day != 1 {
			t.Errorf("Expected day 1, got %d", def.Day)
		}
	}
	
	// Test Good Friday
	if def, exists := definitions["good_friday"]; !exists {
		t.Error("Expected 'good_friday' definition not found")
	} else {
		if def.Name != "Good Friday" {
			t.Errorf("Expected name 'Good Friday', got '%s'", def.Name)
		}
		if def.Calculation != "easter_based" {
			t.Errorf("Expected calculation 'easter_based', got '%s'", def.Calculation)
		}
		if def.EasterOffset != -2 {
			t.Errorf("Expected easter offset -2, got %d", def.EasterOffset)
		}
	}
	
	// Test Labor Day
	if def, exists := definitions["labor_day"]; !exists {
		t.Error("Expected 'labor_day' definition not found")
	} else {
		if def.Name != "Labor Day" {
			t.Errorf("Expected name 'Labor Day', got '%s'", def.Name)
		}
		if def.Calculation != "weekday_based" {
			t.Errorf("Expected calculation 'weekday_based', got '%s'", def.Calculation)
		}
	}
}

func TestPythonASTParser_ConvertMonthName(t *testing.T) {
	parser := NewPythonASTParser("")
	
	testCases := []struct {
		input    string
		expected int
	}{
		{"JAN", 1},
		{"FEB", 2},
		{"MAR", 3},
		{"APR", 4},
		{"MAY", 5},
		{"JUN", 6},
		{"JUL", 7},
		{"AUG", 8},
		{"SEP", 9},
		{"OCT", 10},
		{"NOV", 11},
		{"DEC", 12},
		{"1", 1},
		{"12", 12},
		{"INVALID", 1}, // Should default to 1
	}
	
	for _, tc := range testCases {
		result := parser.convertMonthName(tc.input)
		if result != tc.expected {
			t.Errorf("convertMonthName(%s): expected %d, got %d", tc.input, tc.expected, result)
		}
	}
}

func TestPythonASTParser_ExtractEasterOffset(t *testing.T) {
	parser := NewPythonASTParser("")
	
	testCases := []struct {
		calculation string
		expected    int
	}{
		{"easter(year)", 0},
		{"easter(year) + timedelta(days=1)", 1},
		{"easter(year) - timedelta(days=2)", -2},
		{"easter(year) + timedelta(days=50)", 50},
		{"easter(year) - timedelta(days=100)", -100},
		{"complex calculation", 0}, // Should default to 0
	}
	
	for _, tc := range testCases {
		result := parser.extractEasterOffset(tc.calculation)
		if result != tc.expected {
			t.Errorf("extractEasterOffset(%s): expected %d, got %d", tc.calculation, tc.expected, result)
		}
	}
}

func TestPythonASTParser_RealWorldExample(t *testing.T) {
	// Example from a real Python holidays file
	source := `
class UnitedStates(HolidayBase):
    """Official public holidays for the United States."""
    
    country = "US"
    subdivisions = ["AL", "AK", "AZ", "AR", "CA"]
    
    def __init__(self, **kwargs):
        super().__init__(**kwargs)
    
    def _populate(self, year):
        super()._populate(year)
        
        # New Year's Day
        self._add_holiday("New Year's Day", date(year, JAN, 1))
        
        # Martin Luther King Jr. Day
        if year >= 1983:
            self._add_holiday(
                "Martin Luther King Jr. Day",
                date(year, JAN, 1) + rd(weekday=MO(3))
            )
        
        # Presidents' Day
        self._add_holiday(
            "Washington's Birthday",
            date(year, FEB, 1) + rd(weekday=MO(3))
        )
        
        # Good Friday
        self._add_holiday("Good Friday", easter(year) - rd(days=2))
        
        # Easter Monday  
        self._add_holiday("Easter Monday", easter(year) + rd(days=1))
        
        # Memorial Day
        self._add_holiday(
            "Memorial Day",
            date(year, MAY, 31) + rd(weekday=MO(-1))
        )
        
        # Independence Day
        self._add_holiday("Independence Day", date(year, JUL, 4))
        
        # Labor Day
        self._add_holiday(
            "Labor Day",
            date(year, SEP, 1) + rd(weekday=MO(1))
        )
        
        # Christmas Day
        self._add_holiday("Christmas Day", date(year, DEC, 25))
`

	parser := NewPythonASTParser(source)
	holidayCalls, err := parser.Parse()
	if err != nil {
		t.Fatalf("Parse() failed: %v", err)
	}
	
	if len(holidayCalls) < 5 {
		t.Errorf("Expected at least 5 holiday calls, got %d", len(holidayCalls))
	}
	
	// Convert to definitions
	definitions := parser.ConvertToHolidayDefinitions(holidayCalls)
	
	// Check specific holidays
	expectedHolidays := map[string]struct {
		calculation string
		month       int
		day         int
	}{
		"new_year's_day":    {"fixed", 1, 1},
		"independence_day":  {"fixed", 7, 4},
		"christmas_day":     {"fixed", 12, 25},
		"good_friday":       {"easter_based", 0, 0},
		"easter_monday":     {"easter_based", 0, 0},
	}
	
	for key, expected := range expectedHolidays {
		if def, exists := definitions[key]; !exists {
			t.Errorf("Expected holiday '%s' not found", key)
		} else {
			if def.Calculation != expected.calculation {
				t.Errorf("Holiday '%s': expected calculation '%s', got '%s'", 
					key, expected.calculation, def.Calculation)
			}
			if expected.month > 0 && def.Month != expected.month {
				t.Errorf("Holiday '%s': expected month %d, got %d", 
					key, expected.month, def.Month)
			}
			if expected.day > 0 && def.Day != expected.day {
				t.Errorf("Holiday '%s': expected day %d, got %d", 
					key, expected.day, def.Day)
			}
		}
	}
}

func TestPythonASTParser_PerformanceComparison(t *testing.T) {
	// Create a large Python source for performance testing
	var builder strings.Builder
	builder.WriteString("class TestCountry(HolidayBase):\n")
	builder.WriteString("    def _populate(self, year):\n")
	
	// Add many holiday definitions
	months := []string{"JAN", "FEB", "MAR", "APR", "MAY", "JUN", "JUL", "AUG", "SEP", "OCT", "NOV", "DEC"}
	for i := 0; i < 100; i++ {
		month := months[i%12]
		day := (i % 28) + 1
		builder.WriteString(fmt.Sprintf(`        self._add_holiday("Holiday %d", date(year, %s, %d))`+"\n", i, month, day))
	}
	
	source := builder.String()
	
	// Test our AST parser
	parser := NewPythonASTParser(source)
	start := time.Now()
	holidayCalls, err := parser.Parse()
	astDuration := time.Since(start)
	
	if err != nil {
		t.Fatalf("AST Parser failed: %v", err)
	}
	
	t.Logf("AST Parser: Found %d holidays in %v", len(holidayCalls), astDuration)
	
	// The AST parser should find significantly more holidays than the old regex method
	if len(holidayCalls) < 50 {
		t.Errorf("Expected to find at least 50 holidays, got %d", len(holidayCalls))
	}
}

// Benchmark tests
func BenchmarkPythonASTParser_Parse(b *testing.B) {
	source := `
class TestCountry(HolidayBase):
    def _populate(self, year):
        self._add_holiday("New Year's Day", date(year, JAN, 1))
        self._add_holiday("Good Friday", easter(year) - rd(days=2))
        self._add_holiday("Easter Monday", easter(year) + rd(days=1))
        self._add_holiday("Independence Day", date(year, JUL, 4))
        self._add_holiday("Christmas Day", date(year, DEC, 25))
`

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		parser := NewPythonASTParser(source)
		_, _ = parser.Parse()
	}
}

func BenchmarkPythonASTParser_ConvertToDefinitions(b *testing.B) {
	parser := NewPythonASTParser("")
	
	holidayCalls := []HolidayCall{
		{
			Method: "_add_holiday",
			Name:   "New Year's Day",
			Date: &DateExpression{
				Type:  DateFixed,
				Month: "JAN",
				Day:   1,
			},
		},
		{
			Method: "_add_easter_based_holiday",
			Name:   "Good Friday",
			Date: &DateExpression{
				Type:        DateEasterBased,
				Calculation: "easter(year) - timedelta(days=2)",
			},
		},
	}
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = parser.ConvertToHolidayDefinitions(holidayCalls)
	}
}
