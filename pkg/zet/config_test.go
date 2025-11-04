package zet_test

import (
	"os"
	"testing"

	"github.com/arjungandhi/zet/pkg/zet"
)

func TestGetZetDir(t *testing.T) {
	// Save original env var and restore after test
	originalZetDir := os.Getenv("ZETDIR")
	defer func() {
		if originalZetDir != "" {
			os.Setenv("ZETDIR", originalZetDir)
		} else {
			os.Unsetenv("ZETDIR")
		}
	}()

	tests := []struct {
		name      string
		envValue  string
		shouldSet bool
		wantErr   bool
		expected  string
	}{
		{
			name:      "ZETDIR is set",
			envValue:  "/home/user/notes",
			shouldSet: true,
			wantErr:   false,
			expected:  "/home/user/notes",
		},
		{
			name:      "ZETDIR is not set",
			envValue:  "",
			shouldSet: false,
			wantErr:   true,
			expected:  "",
		},
		{
			name:      "ZETDIR is empty string",
			envValue:  "",
			shouldSet: true,
			wantErr:   true,
			expected:  "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.shouldSet {
				os.Setenv("ZETDIR", tt.envValue)
			} else {
				os.Unsetenv("ZETDIR")
			}

			result, err := zet.GetZetDir()

			if tt.wantErr && err == nil {
				t.Error("GetZetDir() expected error, got nil")
			}

			if !tt.wantErr && err != nil {
				t.Errorf("GetZetDir() unexpected error: %v", err)
			}

			if result != tt.expected {
				t.Errorf("GetZetDir() = %q, want %q", result, tt.expected)
			}
		})
	}
}

func TestGetEditor(t *testing.T) {
	// Save original env var and restore after test
	originalEditor := os.Getenv("EDITOR")
	defer func() {
		if originalEditor != "" {
			os.Setenv("EDITOR", originalEditor)
		} else {
			os.Unsetenv("EDITOR")
		}
	}()

	tests := []struct {
		name      string
		envValue  string
		shouldSet bool
		expected  string
	}{
		{
			name:      "EDITOR is set to vim",
			envValue:  "vim",
			shouldSet: true,
			expected:  "vim",
		},
		{
			name:      "EDITOR is set to nano",
			envValue:  "nano",
			shouldSet: true,
			expected:  "nano",
		},
		{
			name:      "EDITOR is not set (should default to vi)",
			envValue:  "",
			shouldSet: false,
			expected:  "vi",
		},
		{
			name:      "EDITOR is empty string (should default to vi)",
			envValue:  "",
			shouldSet: true,
			expected:  "vi",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.shouldSet {
				os.Setenv("EDITOR", tt.envValue)
			} else {
				os.Unsetenv("EDITOR")
			}

			result := zet.GetEditor()

			if result != tt.expected {
				t.Errorf("GetEditor() = %q, want %q", result, tt.expected)
			}
		})
	}
}

func TestGetRenderer(t *testing.T) {
	result := zet.GetRenderer()
	expected := "glow"

	if result != expected {
		t.Errorf("GetRenderer() = %q, want %q", result, expected)
	}
}
