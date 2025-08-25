package updater

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

// PythonASTParser provides advanced Python AST parsing capabilities
// for extracting holiday definitions from Python source code
type PythonASTParser struct {
	source           string
	tokens           []Token
	currentPos       int
	holidayMethods   map[string]MethodInfo
	classDefinitions map[string]ClassInfo
}

// Token represents a parsed token from Python source
type Token struct {
	Type    TokenType
	Value   string
	Line    int
	Column  int
}

// TokenType represents the type of a parsed token
type TokenType int

const (
	TokenUnknown TokenType = iota
	TokenClass
	TokenDef
	TokenSelf
	TokenMethodCall
	TokenString
	TokenNumber
	TokenIndent
	TokenDedent
	TokenNewline
	TokenOperator
	TokenIdentifier
	TokenComment
	TokenKeyword
)

// MethodInfo contains information about a method definition
type MethodInfo struct {
	Name       string
	Parameters []string
	Body       []Statement
	Line       int
}

// ClassInfo contains information about a class definition
type ClassInfo struct {
	Name       string
	BaseClasses []string
	Methods     map[string]MethodInfo
	Line       int
}

// Statement represents a Python statement
type Statement struct {
	Type       StatementType
	MethodCall *MethodCall
	Assignment *Assignment
	Line       int
}

// StatementType represents the type of a statement
type StatementType int

const (
	StatementUnknown StatementType = iota
	StatementMethodCall
	StatementAssignment
	StatementImport
	StatementReturn
)

// MethodCall represents a method call statement
type MethodCall struct {
	Object    string
	Method    string
	Arguments []Argument
}

// Assignment represents an assignment statement
type Assignment struct {
	Target string
	Value  interface{}
}

// Argument represents a method argument
type Argument struct {
	Type  ArgumentType
	Value interface{}
}

// ArgumentType represents the type of an argument
type ArgumentType int

const (
	ArgumentString ArgumentType = iota
	ArgumentNumber
	ArgumentIdentifier
	ArgumentMethodCall
)

// HolidayCall represents a parsed holiday definition call
type HolidayCall struct {
	Method      string           // _add_holiday, _add_new_years_day, etc.
	Name        string           // Holiday name
	Date        *DateExpression  // Date calculation
	Category    string           // Holiday category
	Languages   map[string]string // Multi-language names
	Conditional string           // Conditional logic (if any)
	Line        int              // Source line number
}

// DateExpression represents a date calculation
type DateExpression struct {
	Type        DateType
	Year        string
	Month       interface{} // Can be int or method call like "date(year, JAN, 1)"
	Day         interface{} // Can be int or calculation
	Calculation string      // Complex calculations like "easter(year) + rd(days=1)"
}

// DateType represents the type of date calculation
type DateType int

const (
	DateFixed DateType = iota
	DateEasterBased
	DateWeekdayBased
	DateCalculated
)

// NewPythonASTParser creates a new Python AST parser
func NewPythonASTParser(source string) *PythonASTParser {
	return &PythonASTParser{
		source:           source,
		holidayMethods:   make(map[string]MethodInfo),
		classDefinitions: make(map[string]ClassInfo),
	}
}

// Parse parses the Python source code and extracts holiday definitions
func (p *PythonASTParser) Parse() ([]HolidayCall, error) {
	// Tokenize the source
	if err := p.tokenize(); err != nil {
		return nil, fmt.Errorf("tokenization failed: %w", err)
	}
	
	// Parse class definitions
	if err := p.parseClasses(); err != nil {
		return nil, fmt.Errorf("class parsing failed: %w", err)
	}
	
	// Extract holiday calls
	holidayCalls, err := p.extractHolidayCalls()
	if err != nil {
		return nil, fmt.Errorf("holiday extraction failed: %w", err)
	}
	
	return holidayCalls, nil
}

// tokenize breaks the source code into tokens
func (p *PythonASTParser) tokenize() error {
	lines := strings.Split(p.source, "\n")
	
	for lineNum, line := range lines {
		if err := p.tokenizeLine(line, lineNum+1); err != nil {
			return fmt.Errorf("error tokenizing line %d: %w", lineNum+1, err)
		}
	}
	
	return nil
}

// tokenizeLine tokenizes a single line
func (p *PythonASTParser) tokenizeLine(line string, lineNum int) error {
	line = strings.TrimSpace(line)
	if line == "" || strings.HasPrefix(line, "#") {
		return nil
	}
	
	// Handle indentation
	originalLine := line
	indentLevel := 0
	for i, char := range originalLine {
		if char == ' ' {
			indentLevel++
		} else {
			line = originalLine[i:]
			break
		}
	}
	
	if indentLevel > 0 {
		p.tokens = append(p.tokens, Token{
			Type:   TokenIndent,
			Value:  strings.Repeat(" ", indentLevel),
			Line:   lineNum,
			Column: 0,
		})
	}
	
	// Tokenize the rest of the line
	p.tokenizeContent(line, lineNum, indentLevel)
	
	return nil
}

// tokenizeContent tokenizes the content of a line
func (p *PythonASTParser) tokenizeContent(content string, lineNum, startColumn int) {
	// Patterns for different token types
	patterns := []struct {
		regex *regexp.Regexp
		tokenType TokenType
	}{
		{regexp.MustCompile(`^class\b`), TokenClass},
		{regexp.MustCompile(`^def\b`), TokenDef},
		{regexp.MustCompile(`^self\b`), TokenSelf},
		{regexp.MustCompile(`^"([^"\\\\]|\\\\.)*"`), TokenString},
		{regexp.MustCompile(`^'([^'\\]|\\.)*'`), TokenString},
		{regexp.MustCompile(`^\d+`), TokenNumber},
		{regexp.MustCompile(`^[a-zA-Z_][a-zA-Z0-9_]*`), TokenIdentifier},
		{regexp.MustCompile(`^[+\-*/=(),.:]`), TokenOperator},
	}
	
	pos := 0
	for pos < len(content) {
		// Skip whitespace
		for pos < len(content) && content[pos] == ' ' {
			pos++
		}
		if pos >= len(content) {
			break
		}
		
		matched := false
		for _, pattern := range patterns {
			if match := pattern.regex.FindString(content[pos:]); match != "" {
				p.tokens = append(p.tokens, Token{
					Type:   pattern.tokenType,
					Value:  match,
					Line:   lineNum,
					Column: startColumn + pos,
				})
				pos += len(match)
				matched = true
				break
			}
		}
		
		if !matched {
			pos++ // Skip unknown character
		}
	}
}

// parseClasses parses class definitions from tokens
func (p *PythonASTParser) parseClasses() error {
	for i := 0; i < len(p.tokens); i++ {
		if p.tokens[i].Type == TokenClass {
			classInfo, newPos, err := p.parseClass(i)
			if err != nil {
				return err
			}
			p.classDefinitions[classInfo.Name] = *classInfo
			i = newPos
		}
	}
	return nil
}

// parseClass parses a single class definition
func (p *PythonASTParser) parseClass(startPos int) (*ClassInfo, int, error) {
	if startPos >= len(p.tokens) || p.tokens[startPos].Type != TokenClass {
		return nil, startPos, fmt.Errorf("expected class token")
	}
	
	// Look for class name
	namePos := startPos + 1
	if namePos >= len(p.tokens) || p.tokens[namePos].Type != TokenIdentifier {
		return nil, startPos, fmt.Errorf("expected class name")
	}
	
	className := p.tokens[namePos].Value
	
	// Parse base classes (simplified)
	baseClasses := []string{}
	pos := namePos + 1
	
	// Skip to end of class header (find colon)
	for pos < len(p.tokens) && p.tokens[pos].Value != ":" {
		if p.tokens[pos].Type == TokenIdentifier && pos > namePos+1 {
			baseClasses = append(baseClasses, p.tokens[pos].Value)
		}
		pos++
	}
	
	if pos >= len(p.tokens) {
		return nil, startPos, fmt.Errorf("incomplete class definition")
	}
	
	// Parse methods within the class
	methods := make(map[string]MethodInfo)
	pos++ // Skip colon
	
	for pos < len(p.tokens) {
		// Look for method definitions
		if p.tokens[pos].Type == TokenDef {
			method, newPos, err := p.parseMethod(pos)
			if err != nil {
				// Skip this method if we can't parse it
				pos++
				continue
			}
			methods[method.Name] = *method
			pos = newPos
		} else if p.tokens[pos].Type == TokenClass {
			// Hit another class, we're done with this one
			break
		} else {
			pos++
		}
	}
	
	return &ClassInfo{
		Name:        className,
		BaseClasses: baseClasses,
		Methods:     methods,
		Line:        p.tokens[startPos].Line,
	}, pos, nil
}

// parseMethod parses a method definition
func (p *PythonASTParser) parseMethod(startPos int) (*MethodInfo, int, error) {
	if startPos >= len(p.tokens) || p.tokens[startPos].Type != TokenDef {
		return nil, startPos, fmt.Errorf("expected def token")
	}
	
	// Get method name
	namePos := startPos + 1
	if namePos >= len(p.tokens) || p.tokens[namePos].Type != TokenIdentifier {
		return nil, startPos, fmt.Errorf("expected method name")
	}
	
	methodName := p.tokens[namePos].Value
	
	// Parse parameters (simplified)
	parameters := []string{}
	pos := namePos + 1
	
	// Skip to end of method header
	for pos < len(p.tokens) && p.tokens[pos].Value != ":" {
		if p.tokens[pos].Type == TokenIdentifier {
			parameters = append(parameters, p.tokens[pos].Value)
		}
		pos++
	}
	
	if pos >= len(p.tokens) {
		return nil, startPos, fmt.Errorf("incomplete method definition")
	}
	
	// Parse method body (simplified)
	body := []Statement{}
	pos++ // Skip colon
	
	// For now, we'll just capture the raw body
	// In a full implementation, we'd parse statements properly
	
	return &MethodInfo{
		Name:       methodName,
		Parameters: parameters,
		Body:       body,
		Line:       p.tokens[startPos].Line,
	}, pos, nil
}

// extractHolidayCalls extracts holiday definition calls from the parsed AST
func (p *PythonASTParser) extractHolidayCalls() ([]HolidayCall, error) {
	var holidayCalls []HolidayCall
	
	// Look for holiday method calls in the source
	holidayPatterns := []*regexp.Regexp{
		regexp.MustCompile(`self\._add_holiday\s*\(`),
		regexp.MustCompile(`self\._add_new_years_day\s*\(`),
		regexp.MustCompile(`self\._add_christmas_day\s*\(`),
		regexp.MustCompile(`self\._add_easter_based_holiday\s*\(`),
		regexp.MustCompile(`self\._add_weekday_holiday\s*\(`),
	}
	
	lines := strings.Split(p.source, "\n")
	
	for lineNum, line := range lines {
		for _, pattern := range holidayPatterns {
			if pattern.MatchString(line) {
				holidayCall, err := p.parseHolidayCall(line, lineNum+1)
				if err != nil {
					// Log error but continue parsing
					continue
				}
				if holidayCall != nil {
					holidayCalls = append(holidayCalls, *holidayCall)
				}
			}
		}
	}
	
	return holidayCalls, nil
}

// parseHolidayCall parses a single holiday call from a line
func (p *PythonASTParser) parseHolidayCall(line string, lineNum int) (*HolidayCall, error) {
	// Extract method name
	methodPattern := regexp.MustCompile(`self\.(_add_[a-z_]+)\s*\(`)
	methodMatch := methodPattern.FindStringSubmatch(line)
	if len(methodMatch) < 2 {
		return nil, fmt.Errorf("could not extract method name")
	}
	
	methodName := methodMatch[1]
	
	// Parse different types of holiday calls
	switch {
	case strings.Contains(methodName, "_add_holiday"):
		return p.parseAddHolidayCall(line, lineNum, methodName)
	case strings.Contains(methodName, "_easter_based"):
		return p.parseEasterBasedCall(line, lineNum, methodName)
	case strings.Contains(methodName, "_weekday"):
		return p.parseWeekdayCall(line, lineNum, methodName)
	default:
		return p.parseGenericHolidayCall(line, lineNum, methodName)
	}
}

// parseAddHolidayCall parses a standard _add_holiday call
func (p *PythonASTParser) parseAddHolidayCall(line string, lineNum int, methodName string) (*HolidayCall, error) {
	// Pattern: self._add_holiday(date_expr, "Holiday Name")
	// or: self._add_holiday("Holiday Name", date_expr)
	
	// Extract holiday name - try both double and single quotes separately
	var holidayName string
	
	// Try double quotes first
	doubleQuotePattern := regexp.MustCompile(`"([^"]*)"`)
	if matches := doubleQuotePattern.FindStringSubmatch(line); len(matches) >= 2 {
		holidayName = matches[1]
	} else {
		// Try single quotes
		singleQuotePattern := regexp.MustCompile(`'([^']*)'`)
		if matches := singleQuotePattern.FindStringSubmatch(line); len(matches) >= 2 {
			holidayName = matches[1]
		} else {
			return nil, fmt.Errorf("could not extract holiday name")
		}
	}
	
	// Extract date expression
	dateExpr, err := p.extractDateExpression(line)
	if err != nil {
		return nil, fmt.Errorf("could not extract date expression: %w", err)
	}
	
	return &HolidayCall{
		Method: methodName,
		Name:   holidayName,
		Date:   dateExpr,
		Line:   lineNum,
	}, nil
}

// parseEasterBasedCall parses Easter-based holiday calls
func (p *PythonASTParser) parseEasterBasedCall(line string, lineNum int, methodName string) (*HolidayCall, error) {
	// Pattern: self._add_easter_based_holiday("Name", days_offset)
	
	var holidayName string
	
	// Try double quotes first
	doubleQuotePattern := regexp.MustCompile(`"([^"]*)"`)
	if matches := doubleQuotePattern.FindStringSubmatch(line); len(matches) >= 2 {
		holidayName = matches[1]
	} else {
		// Try single quotes
		singleQuotePattern := regexp.MustCompile(`'([^']*)'`)
		if matches := singleQuotePattern.FindStringSubmatch(line); len(matches) >= 2 {
			holidayName = matches[1]
		} else {
			return nil, fmt.Errorf("could not extract holiday name")
		}
	}
	
	// Extract days offset
	offsetPattern := regexp.MustCompile(`[+-]?\d+`)
	offsetMatch := offsetPattern.FindString(line)
	
	calculation := "easter(year)"
	if offsetMatch != "" {
		offset, _ := strconv.Atoi(offsetMatch)
		if offset != 0 {
			if offset > 0 {
				calculation = fmt.Sprintf("easter(year) + timedelta(days=%d)", offset)
			} else {
				calculation = fmt.Sprintf("easter(year) - timedelta(days=%d)", -offset)
			}
		}
	}
	
	return &HolidayCall{
		Method: methodName,
		Name:   holidayName,
		Date: &DateExpression{
			Type:        DateEasterBased,
			Calculation: calculation,
		},
		Line: lineNum,
	}, nil
}

// parseWeekdayCall parses weekday-based holiday calls
func (p *PythonASTParser) parseWeekdayCall(line string, lineNum int, methodName string) (*HolidayCall, error) {
	// Pattern: self._add_weekday_holiday("Name", month, weekday, week)
	
	var holidayName string
	
	// Try double quotes first
	doubleQuotePattern := regexp.MustCompile(`"([^"]*)"`)
	if matches := doubleQuotePattern.FindStringSubmatch(line); len(matches) >= 2 {
		holidayName = matches[1]
	} else {
		// Try single quotes
		singleQuotePattern := regexp.MustCompile(`'([^']*)'`)
		if matches := singleQuotePattern.FindStringSubmatch(line); len(matches) >= 2 {
			holidayName = matches[1]
		} else {
			return nil, fmt.Errorf("could not extract holiday name")
		}
	}
	
	return &HolidayCall{
		Method: methodName,
		Name:   holidayName,
		Date: &DateExpression{
			Type:        DateWeekdayBased,
			Calculation: "weekday_based", // Simplified
		},
		Line: lineNum,
	}, nil
}

// parseGenericHolidayCall parses other types of holiday calls
func (p *PythonASTParser) parseGenericHolidayCall(line string, lineNum int, methodName string) (*HolidayCall, error) {
	var holidayName string
	
	// Try double quotes first
	doubleQuotePattern := regexp.MustCompile(`"([^"]*)"`)
	if matches := doubleQuotePattern.FindStringSubmatch(line); len(matches) >= 2 {
		holidayName = matches[1]
	} else {
		// Try single quotes
		singleQuotePattern := regexp.MustCompile(`'([^']*)'`)
		if matches := singleQuotePattern.FindStringSubmatch(line); len(matches) >= 2 {
			holidayName = matches[1]
		} else {
			return nil, fmt.Errorf("could not extract holiday name")
		}
	}
	
	return &HolidayCall{
		Method: methodName,
		Name:   holidayName,
		Date: &DateExpression{
			Type:        DateCalculated,
			Calculation: "generic",
		},
		Line: lineNum,
	}, nil
}

// extractDateExpression extracts date expressions from holiday calls
func (p *PythonASTParser) extractDateExpression(line string) (*DateExpression, error) {
	// Look for common date patterns
	
	// Fixed date: date(year, MONTH, day)
	fixedDatePattern := regexp.MustCompile(`date\s*\(\s*year\s*,\s*([A-Z]+|\d+)\s*,\s*(\d+)\s*\)`)
	if match := fixedDatePattern.FindStringSubmatch(line); len(match) >= 3 {
		month := match[1]
		day, _ := strconv.Atoi(match[2])
		
		return &DateExpression{
			Type:  DateFixed,
			Year:  "year",
			Month: month,
			Day:   day,
		}, nil
	}
	
	// Easter-based: easter(year) + rd(days=N)
	easterPattern := regexp.MustCompile(`easter\s*\(\s*year\s*\)(\s*[+-]\s*rd\s*\(\s*days\s*=\s*(\d+)\s*\))?`)
	if match := easterPattern.FindStringSubmatch(line); len(match) >= 1 {
		calculation := "easter(year)"
		if len(match) >= 3 && match[2] != "" {
			days := match[2]
			if strings.Contains(match[0], "+") {
				calculation = fmt.Sprintf("easter(year) + timedelta(days=%s)", days)
			} else {
				calculation = fmt.Sprintf("easter(year) - timedelta(days=%s)", days)
			}
		}
		
		return &DateExpression{
			Type:        DateEasterBased,
			Calculation: calculation,
		}, nil
	}
	
	// Default to calculated
	return &DateExpression{
		Type:        DateCalculated,
		Calculation: "complex",
	}, nil
}

// ConvertToHolidayDefinitions converts parsed holiday calls to HolidayDefinition format
func (p *PythonASTParser) ConvertToHolidayDefinitions(holidayCalls []HolidayCall) map[string]HolidayDefinition {
	definitions := make(map[string]HolidayDefinition)
	
	for _, call := range holidayCalls {
		key := strings.ToLower(strings.ReplaceAll(call.Name, " ", "_"))
		
		definition := HolidayDefinition{
			Name:      call.Name,
			Category:  "public", // Default category
			Languages: map[string]string{"en": call.Name},
		}
		
		// Convert date expression to definition fields
		if call.Date != nil {
			switch call.Date.Type {
			case DateFixed:
				definition.Calculation = "fixed"
				if month, ok := call.Date.Month.(string); ok {
					definition.Month = p.convertMonthName(month)
				} else if month, ok := call.Date.Month.(int); ok {
					definition.Month = month
				}
				if day, ok := call.Date.Day.(int); ok {
					definition.Day = day
				}
				
			case DateEasterBased:
				definition.Calculation = "easter_based"
				definition.EasterOffset = p.extractEasterOffset(call.Date.Calculation)
				
			case DateWeekdayBased:
				definition.Calculation = "weekday_based"
				// Would need more parsing for weekday details
				
			default:
				definition.Calculation = "complex"
			}
		}
		
		definitions[key] = definition
	}
	
	return definitions
}

// convertMonthName converts Python month constants to integers
func (p *PythonASTParser) convertMonthName(monthName string) int {
	monthMap := map[string]int{
		"JAN": 1, "FEB": 2, "MAR": 3, "APR": 4,
		"MAY": 5, "JUN": 6, "JUL": 7, "AUG": 8,
		"SEP": 9, "OCT": 10, "NOV": 11, "DEC": 12,
	}
	
	if month, exists := monthMap[monthName]; exists {
		return month
	}
	
	// Try to parse as integer
	if month, err := strconv.Atoi(monthName); err == nil {
		return month
	}
	
	return 1 // Default to January
}

// extractEasterOffset extracts the day offset from easter calculations
func (p *PythonASTParser) extractEasterOffset(calculation string) int {
	offsetPattern := regexp.MustCompile(`days\s*=\s*([+-]?\d+)`)
	if match := offsetPattern.FindStringSubmatch(calculation); len(match) >= 2 {
		offset, _ := strconv.Atoi(match[1])
		// Check if it's a subtraction in the calculation
		if strings.Contains(calculation, "- timedelta") {
			return -offset
		}
		return offset
	}
	return 0
}
