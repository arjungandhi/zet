package zet_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/arjungandhi/zet/pkg/zet"
)

func TestSanitizeFilename(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "simple title",
			input:    "My Note",
			expected: "My Note",
		},
		{
			name:     "title with special characters",
			input:    "My New Note: Ideas & Thoughts!",
			expected: "My New Note Ideas Thoughts",
		},
		{
			name:     "title with question mark",
			input:    "What's the plan?",
			expected: "Whats the plan",
		},
		{
			name:     "title with multiple spaces",
			input:    "My    Note   Title",
			expected: "My Note Title",
		},
		{
			name:     "title with leading/trailing spaces",
			input:    "  My Note  ",
			expected: "My Note",
		},
		{
			name:     "title with slashes",
			input:    "TCP/IP & HTTP/HTTPS",
			expected: "TCPIP HTTPHTTPS",
		},
		{
			name:     "empty title",
			input:    "",
			expected: "",
		},
		{
			name:     "only special characters",
			input:    "!@#$%^&*()",
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := zet.SanitizeFilename(tt.input)
			if result != tt.expected {
				t.Errorf("SanitizeFilename(%q) = %q, want %q", tt.input, result, tt.expected)
			}
		})
	}
}

func TestNoteExists(t *testing.T) {
	zetDir := ZetDir(t)
	defer Cleanup(t, zetDir)

	testPath := filepath.Join(zetDir, "Test Note.md")
	err := os.WriteFile(testPath, []byte("test content"), 0644)
	if err != nil {
		t.Fatal(err)
	}

	tests := []struct {
		name     string
		title    string
		expected bool
	}{
		{
			name:     "existing note",
			title:    "Test Note",
			expected: true,
		},
		{
			name:     "non-existing note",
			title:    "Non Existent Note",
			expected: false,
		},
		{
			name:     "title with special characters (should sanitize and match)",
			title:    "Test: Note!",
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := zet.NoteExists(zetDir, tt.title)
			if result != tt.expected {
				t.Errorf("NoteExists(%q) = %v, want %v", tt.title, result, tt.expected)
			}
		})
	}
}

func TestWriteNote(t *testing.T) {
	zetDir := ZetDir(t)
	defer Cleanup(t, zetDir)

	title := "My Test Note"
	content := "This is my test note content"

	err := zet.WriteNote(zetDir, title, content)
	if err != nil {
		t.Fatalf("WriteNote() error = %v", err)
	}

	// Verify file was created with correct name
	expectedPath := filepath.Join(zetDir, "My Test Note.md")
	fileContent, err := os.ReadFile(expectedPath)
	if err != nil {
		t.Fatalf("Failed to read created note: %v", err)
	}

	if string(fileContent) != content {
		t.Errorf("File content = %q, want %q", string(fileContent), content)
	}
}

func TestWriteNoteWithSpecialCharacters(t *testing.T) {
	zetDir := ZetDir(t)
	defer Cleanup(t, zetDir)

	title := "My Note: Ideas & Thoughts!"
	content := "Test content"

	err := zet.WriteNote(zetDir, title, content)
	if err != nil {
		t.Fatalf("WriteNote() error = %v", err)
	}

	// Verify file was created with sanitized name
	expectedPath := filepath.Join(zetDir, "My Note Ideas Thoughts.md")
	_, err = os.Stat(expectedPath)
	if err != nil {
		t.Errorf("Expected file at %s to exist, but got error: %v", expectedPath, err)
	}
}

func TestReadNote(t *testing.T) {
	zetDir := ZetDir(t)
	defer Cleanup(t, zetDir)

	// Create a test note
	title := "Test Note"
	content := "This is test content"
	path := filepath.Join(zetDir, "Test Note.md")
	err := os.WriteFile(path, []byte(content), 0644)
	if err != nil {
		t.Fatal(err)
	}

	// Read the note
	note, err := zet.ReadNote(path)
	if err != nil {
		t.Fatalf("ReadNote() error = %v", err)
	}

	if note.Title != title {
		t.Errorf("note.Title = %q, want %q", note.Title, title)
	}

	if note.Path != path {
		t.Errorf("note.Path = %q, want %q", note.Path, path)
	}

	if note.Body != content {
		t.Errorf("note.Body = %q, want %q", note.Body, content)
	}
}

func TestListNoteFiles(t *testing.T) {
	zetDir := ZetDir(t)
	defer Cleanup(t, zetDir)

	// Create multiple test notes
	noteFiles := []string{
		"Zebra Note.md",
		"Apple Note.md",
		"Middle Note.md",
	}

	for _, filename := range noteFiles {
		path := filepath.Join(zetDir, filename)
		err := os.WriteFile(path, []byte("content"), 0644)
		if err != nil {
			t.Fatal(err)
		}
	}

	// Create a non-markdown file (should be ignored)
	txtPath := filepath.Join(zetDir, "ignore.txt")
	err := os.WriteFile(txtPath, []byte("ignore"), 0644)
	if err != nil {
		t.Fatal(err)
	}

	// List notes
	files, err := zet.ListNoteFiles(zetDir)
	if err != nil {
		t.Fatalf("ListNoteFiles() error = %v", err)
	}

	// Should only return .md files
	if len(files) != len(noteFiles) {
		t.Errorf("ListNoteFiles() returned %d files, want %d", len(files), len(noteFiles))
	}

	// Verify files are sorted alphabetically
	expectedOrder := []string{
		filepath.Join(zetDir, "Apple Note.md"),
		filepath.Join(zetDir, "Middle Note.md"),
		filepath.Join(zetDir, "Zebra Note.md"),
	}

	for i, file := range files {
		if file != expectedOrder[i] {
			t.Errorf("files[%d] = %q, want %q", i, file, expectedOrder[i])
		}
	}
}

func TestListNoteFilesEmptyDirectory(t *testing.T) {
	zetDir := ZetDir(t)
	defer Cleanup(t, zetDir)

	files, err := zet.ListNoteFiles(zetDir)
	if err != nil {
		t.Fatalf("ListNoteFiles() error = %v", err)
	}

	if len(files) != 0 {
		t.Errorf("ListNoteFiles() on empty dir returned %d files, want 0", len(files))
	}
}
