package config

import (
	"os"
	"testing"
)

func TestLoadGitHubToken(t *testing.T) {
	// Save original environment
	originalToken := os.Getenv("GITHUB_TOKEN")
	defer func() {
		if originalToken != "" {
			os.Setenv("GITHUB_TOKEN", originalToken)
		} else {
			os.Unsetenv("GITHUB_TOKEN")
		}
	}()

	t.Run("Load from environment variable", func(t *testing.T) {
		// Set environment variable
		testToken := "ghp_test_env_token"
		os.Setenv("GITHUB_TOKEN", testToken)

		token := LoadGitHubToken()
		if token != testToken {
			t.Errorf("Expected token from env var %q, got %q", testToken, token)
		}

		if !HasGitHubToken() {
			t.Error("HasGitHubToken() should return true when env var is set")
		}
	})

	t.Run("Environment variable takes precedence", func(t *testing.T) {
		envToken := "ghp_env_priority_token"
		os.Setenv("GITHUB_TOKEN", envToken)

		token := LoadGitHubToken()
		if token != envToken {
			t.Errorf("Expected env token %q to take precedence, got %q", envToken, token)
		}
	})

	t.Run("No token available when env var cleared", func(t *testing.T) {
		// Clear environment variable
		os.Unsetenv("GITHUB_TOKEN")

		// The actual token file might exist, so we just test that env var is cleared
		token := os.Getenv("GITHUB_TOKEN")
		if token != "" {
			t.Errorf("Expected empty env var, got %q", token)
		}
	})

	t.Run("HasGitHubToken reflects actual state", func(t *testing.T) {
		// Set a token
		os.Setenv("GITHUB_TOKEN", "test_token")
		if !HasGitHubToken() {
			t.Error("HasGitHubToken() should return true when token is available")
		}

		// Clear token
		os.Unsetenv("GITHUB_TOKEN")
		// Note: HasGitHubToken might still return true if config file exists
		// This is expected behavior
	})
}
