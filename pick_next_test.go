package main

import (
	"fmt"
	"os"
	"strings"
	"testing"
)

func TestGetTeamFile(t *testing.T) {
	tests := []struct {
		name        string
		flagValue   string
		envValue    string
		expected    string
		description string
	}{
		{
			name:        "flag takes precedence over env",
			flagValue:   "flag-team.txt",
			envValue:    "env-team.txt",
			expected:    "flag-team.txt",
			description: "Command-line flag should override environment variable",
		},
		{
			name:        "env fallback when no flag",
			flagValue:   "",
			envValue:    "env-team.txt",
			expected:    "env-team.txt",
			description: "Environment variable should be used when no flag is provided",
		},
		{
			name:        "default when neither flag nor env",
			flagValue:   "",
			envValue:    "",
			expected:    "team.txt",
			description: "Should use default team.txt when neither flag nor env is set",
		},
		{
			name:        "stdin flag takes precedence",
			flagValue:   "-",
			envValue:    "env-team.txt",
			expected:    "-",
			description: "Stdin flag should override environment variable",
		},
		{
			name:        "empty flag falls back to env",
			flagValue:   "",
			envValue:    "backup-team.txt",
			expected:    "backup-team.txt",
			description: "Empty flag should fall back to environment variable",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Set up environment variable
			if tt.envValue != "" {
				t.Setenv("TEAM_FILE", tt.envValue)
			} else {
				// Explicitly unset the environment variable for default behavior tests
				t.Setenv("TEAM_FILE", "")
			}

			result := getTeamFile(tt.flagValue)
			if result != tt.expected {
				t.Errorf("getTeamFile(%q) with TEAM_FILE=%q = %q; want %q\n%s",
					tt.flagValue, tt.envValue, result, tt.expected, tt.description)
			}
		})
	}
}

func TestLoadTeamMembers_FromFile(t *testing.T) {
	// Create a temporary test file
	testContent := `# Team members for testing
Alice
Bob
# This is a comment
Charlie

Diana
# Another comment
Eve`

	tmpFile, err := os.CreateTemp("", "test-team-*.txt")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer func() {
		if err := os.Remove(tmpFile.Name()); err != nil {
			t.Logf("Warning: failed to remove temp file: %v", err)
		}
	}()

	if _, err := tmpFile.WriteString(testContent); err != nil {
		t.Fatalf("Failed to write to temp file: %v", err)
	}
	if err := tmpFile.Close(); err != nil {
		t.Fatalf("Failed to close temp file: %v", err)
	}

	members, err := loadTeamMembers(tmpFile.Name())
	if err != nil {
		t.Fatalf("loadTeamMembers failed: %v", err)
	}

	expected := []string{"Alice", "Bob", "Charlie", "Diana", "Eve"}
	if len(members) != len(expected) {
		t.Fatalf("Expected %d members, got %d", len(expected), len(members))
	}

	for i, member := range members {
		if member != expected[i] {
			t.Errorf("Expected member %d to be %q, got %q", i, expected[i], member)
		}
	}
}

func TestLoadTeamMembers_FromStdin(t *testing.T) {
	testInput := "Alice\nBob\n# Comment line\nCharlie\n\nDiana\n"

	// Save original stdin
	oldStdin := os.Stdin
	defer func() { os.Stdin = oldStdin }()

	// Create a pipe with test data
	r, w, err := os.Pipe()
	if err != nil {
		t.Fatalf("Failed to create pipe: %v", err)
	}

	// Write test data to the pipe
	go func() {
		defer func() {
			if err := w.Close(); err != nil {
				t.Errorf("Failed to close writer: %v", err)
			}
		}()
		if _, err := w.WriteString(testInput); err != nil {
			t.Errorf("Failed to write test input: %v", err)
		}
	}()

	// Replace stdin with our pipe
	os.Stdin = r
	defer func() {
		if err := r.Close(); err != nil {
			t.Logf("Warning: failed to close reader: %v", err)
		}
	}()

	members, err := loadTeamMembers("-")
	if err != nil {
		t.Fatalf("loadTeamMembers from stdin failed: %v", err)
	}

	expected := []string{"Alice", "Bob", "Charlie", "Diana"}
	if len(members) != len(expected) {
		t.Fatalf("Expected %d members, got %d", len(expected), len(members))
	}

	for i, member := range members {
		if member != expected[i] {
			t.Errorf("Expected member %d to be %q, got %q", i, expected[i], member)
		}
	}
}

func TestLoadTeamMembers_FileNotFound(t *testing.T) {
	_, err := loadTeamMembers("nonexistent-file.txt")
	if err == nil {
		t.Error("Expected error for nonexistent file, got nil")
	}
}

func TestLoadTeamMembers_EmptyFile(t *testing.T) {
	// Create empty temp file
	tmpFile, err := os.CreateTemp("", "empty-team-*.txt")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer func() {
		if err := os.Remove(tmpFile.Name()); err != nil {
			t.Logf("Warning: failed to remove temp file: %v", err)
		}
	}()
	if err := tmpFile.Close(); err != nil {
		t.Fatalf("Failed to close temp file: %v", err)
	}

	members, err := loadTeamMembers(tmpFile.Name())
	if err != nil {
		t.Fatalf("loadTeamMembers failed: %v", err)
	}

	if len(members) != 0 {
		t.Errorf("Expected empty slice for empty file, got %d members", len(members))
	}
}

func TestLoadTeamMembers_OnlyComments(t *testing.T) {
	testContent := `# Only comments here
# Another comment
# Yet another comment`

	tmpFile, err := os.CreateTemp("", "comments-team-*.txt")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer func() {
		if err := os.Remove(tmpFile.Name()); err != nil {
			t.Logf("Warning: failed to remove temp file: %v", err)
		}
	}()

	if _, err := tmpFile.WriteString(testContent); err != nil {
		t.Fatalf("Failed to write to temp file: %v", err)
	}
	if err := tmpFile.Close(); err != nil {
		t.Fatalf("Failed to close temp file: %v", err)
	}

	members, err := loadTeamMembers(tmpFile.Name())
	if err != nil {
		t.Fatalf("loadTeamMembers failed: %v", err)
	}

	if len(members) != 0 {
		t.Errorf("Expected empty slice for comments-only file, got %d members", len(members))
	}
}

func TestLoadTeamMembers_WhitespaceHandling(t *testing.T) {
	testContent := `  Alice  
	Bob	
Charlie
  Diana  `

	tmpFile, err := os.CreateTemp("", "whitespace-team-*.txt")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer func() {
		if err := os.Remove(tmpFile.Name()); err != nil {
			t.Logf("Warning: failed to remove temp file: %v", err)
		}
	}()

	if _, err := tmpFile.WriteString(testContent); err != nil {
		t.Fatalf("Failed to write to temp file: %v", err)
	}
	if err := tmpFile.Close(); err != nil {
		t.Fatalf("Failed to close temp file: %v", err)
	}

	members, err := loadTeamMembers(tmpFile.Name())
	if err != nil {
		t.Fatalf("loadTeamMembers failed: %v", err)
	}

	expected := []string{"Alice", "Bob", "Charlie", "Diana"}
	if len(members) != len(expected) {
		t.Fatalf("Expected %d members, got %d", len(expected), len(members))
	}

	for i, member := range members {
		if member != expected[i] {
			t.Errorf("Expected member %d to be %q, got %q", i, expected[i], member)
		}
	}
}

func TestGetStateFile(t *testing.T) {
	tests := []struct {
		name       string
		envValue   string
		shouldTest string // "default" or "exact"
	}{
		{
			name:       "default state file",
			envValue:   "",
			shouldTest: "default",
		},
		{
			name:       "custom state file from env",
			envValue:   "/custom/path/state.txt",
			shouldTest: "exact",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.envValue != "" {
				t.Setenv("STATE_FILE", tt.envValue)
			} else {
				// Explicitly unset the environment variable for default behavior tests
				t.Setenv("STATE_FILE", "")
			}

			result := getStateFile()

			if tt.shouldTest == "default" {
				// For default, should contain temp directory and expected filename
				if !strings.Contains(result, "daily-scrum-picker-remaining.txt") {
					t.Errorf("Default state file should contain 'daily-scrum-picker-remaining.txt', got %q", result)
				}
			} else {
				// For custom, should be exact match
				if result != tt.envValue {
					t.Errorf("getStateFile() = %q; want %q", result, tt.envValue)
				}
			}
		})
	}
}

// Benchmark tests
func BenchmarkGetTeamFile(b *testing.B) {
	b.Setenv("TEAM_FILE", "bench-team.txt")

	for i := 0; i < b.N; i++ {
		getTeamFile("flag-team.txt")
	}
}

func BenchmarkLoadTeamMembers(b *testing.B) {
	// Create a test file with many members
	tmpFile, err := os.CreateTemp("", "bench-team-*.txt")
	if err != nil {
		b.Fatalf("Failed to create temp file: %v", err)
	}
	defer func() {
		if err := os.Remove(tmpFile.Name()); err != nil {
			b.Logf("Warning: failed to remove temp file: %v", err)
		}
	}()

	// Write 1000 team members
	for i := 0; i < 1000; i++ {
		if _, err := fmt.Fprintf(tmpFile, "Member%d\n", i); err != nil {
			b.Fatalf("Failed to write to temp file: %v", err)
		}
	}
	if err := tmpFile.Close(); err != nil {
		b.Fatalf("Failed to close temp file: %v", err)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if _, err := loadTeamMembers(tmpFile.Name()); err != nil {
			b.Fatalf("loadTeamMembers failed: %v", err)
		}
	}
}
