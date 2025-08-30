package updater

import (
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

// MockHTTPTransport is a mock HTTP transport for testing
type MockHTTPTransport struct {
	responses    map[string]*http.Response
	shouldError  bool
	errorMessage string
}

// NewMockHTTPTransport creates a new mock HTTP transport
func NewMockHTTPTransport() *MockHTTPTransport {
	return &MockHTTPTransport{
		responses:   make(map[string]*http.Response),
		shouldError: false,
	}
}

// SetError configures the mock to return an error
func (m *MockHTTPTransport) SetError(shouldError bool, message string) {
	m.shouldError = shouldError
	m.errorMessage = message
}

// AddResponse adds a mock response for a specific URL
func (m *MockHTTPTransport) AddResponse(url string, statusCode int, body string) {
	m.responses[url] = &http.Response{
		StatusCode: statusCode,
		Status:     fmt.Sprintf("%d %s", statusCode, http.StatusText(statusCode)),
		Body:       io.NopCloser(strings.NewReader(body)),
		Header:     make(http.Header),
	}
}

// RoundTrip implements the http.RoundTripper interface
func (m *MockHTTPTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	if m.shouldError {
		return nil, fmt.Errorf("mock HTTP error: %s", m.errorMessage)
	}

	url := req.URL.String()
	if resp, exists := m.responses[url]; exists {
		return resp, nil
	}

	// Default response for unmatched URLs
	return &http.Response{
		StatusCode: 404,
		Status:     "404 Not Found",
		Body:       io.NopCloser(strings.NewReader("Not Found")),
		Header:     make(http.Header),
	}, nil
}

// NewMockPythonHolidaysSync creates a PythonHolidaysSync with a mock HTTP client
func NewMockPythonHolidaysSync(dataDir string) *PythonHolidaysSync {
	transport := NewMockHTTPTransport()

	// Add mock response for commits endpoint
	commitsResponse := `[
		{
			"sha": "abc123",
			"commit": {
				"committer": {
					"date": "` + time.Now().Format(time.RFC3339) + `"
				}
			}
		}
	]`
	transport.AddResponse("https://api.github.com/repos/vacanza/holidays/commits", 200, commitsResponse)

	client := &http.Client{
		Transport: transport,
		Timeout:   30 * time.Second,
	}

	sync := NewPythonHolidaysSync(dataDir)
	sync.httpClient = client
	return sync
}
