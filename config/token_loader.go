package config

import (
	"os"
	"path/filepath"
	"strings"
)

// LoadGitHubToken loads the GitHub token from various sources in order of priority:
// 1. Environment variable GITHUB_TOKEN
// 2. config/github_token.txt file
// 3. Returns empty string if not found
func LoadGitHubToken() string {
	// First, try environment variable
	if token := os.Getenv("GITHUB_TOKEN"); token != "" {
		return strings.TrimSpace(token)
	}

	// Second, try loading from config file
	tokenFile := filepath.Join("config", "github_token.txt")
	if data, err := os.ReadFile(tokenFile); err == nil {
		token := strings.TrimSpace(string(data))
		if token != "" {
			return token
		}
	}

	// No token found
	return ""
}

// HasGitHubToken checks if a GitHub token is available from any source
func HasGitHubToken() bool {
	return LoadGitHubToken() != ""
}
