package zet_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/arjungandhi/zet/pkg/zet"
)

func TestCreateOrEditNote_Create(t *testing.T) {
	zetDir := ZetDir(t)
	defer Cleanup(t, zetDir)

	originalZetDir := os.Getenv("ZETDIR")
	os.Setenv("ZETDIR", zetDir)
	defer func() {
		if originalZetDir != "" {
			os.Setenv("ZETDIR", originalZetDir)
		} else {
			os.Unsetenv("ZETDIR")
		}
	}()

	title := "My First Note"

	notePath, err := zet.CreateOrEditNote(title)
	if err != nil {
		t.Fatal(err)
	}

	expectedPath := filepath.Join(zetDir, "My First Note.md")
	if notePath != expectedPath {
		t.Errorf("CreateOrEditNote() path = %q, want %q", notePath, expectedPath)
	}

	_, err = os.Stat(notePath)
	if err != nil {
		t.Errorf("Note file should exist at %s, got error: %v", notePath, err)
	}
}

func TestCreateOrEditNote_Edit(t *testing.T) {
	zetDir := ZetDir(t)
	defer Cleanup(t, zetDir)

	originalZetDir := os.Getenv("ZETDIR")
	os.Setenv("ZETDIR", zetDir)
	defer func() {
		if originalZetDir != "" {
			os.Setenv("ZETDIR", originalZetDir)
		} else {
			os.Unsetenv("ZETDIR")
		}
	}()

	title := "Existing Note"
	existingContent := "This note already exists"
	notePath := filepath.Join(zetDir, "Existing Note.md")
	err := os.WriteFile(notePath, []byte(existingContent), 0644)
	if err != nil {
		t.Fatal(err)
	}

	returnedPath, err := zet.CreateOrEditNote(title)
	if err != nil {
		t.Fatal(err)
	}

	if returnedPath != notePath {
		t.Errorf("CreateOrEditNote() path = %q, want %q", returnedPath, notePath)
	}

	content, err := os.ReadFile(notePath)
	if err != nil {
		t.Fatal(err)
	}

	if string(content) != existingContent {
		t.Error("CreateOrEditNote() on existing note should not modify content")
	}
}

func TestDeleteNote(t *testing.T) {
	zetDir := ZetDir(t)
	defer Cleanup(t, zetDir)

	notePath := filepath.Join(zetDir, "Test Note.md")
	err := os.WriteFile(notePath, []byte("content"), 0644)
	if err != nil {
		t.Fatal(err)
	}

	note := &zet.Note{
		Title: "Test Note",
		Path:  notePath,
		Body:  "content",
	}

	err = zet.DeleteNote(note)
	if err != nil {
		t.Fatal(err)
	}

	_, err = os.Stat(notePath)
	if err == nil {
		t.Errorf("Expected note at %s to be deleted", notePath)
	}
}

func TestListNotes(t *testing.T) {
	zetDir := ZetDir(t)
	defer Cleanup(t, zetDir)

	originalZetDir := os.Getenv("ZETDIR")
	os.Setenv("ZETDIR", zetDir)
	defer func() {
		if originalZetDir != "" {
			os.Setenv("ZETDIR", originalZetDir)
		} else {
			os.Unsetenv("ZETDIR")
		}
	}()

	noteTitles := []string{
		"Zebra Note",
		"Apple Note",
		"Middle Note",
	}

	for _, title := range noteTitles {
		path := filepath.Join(zetDir, title+".md")
		err := os.WriteFile(path, []byte("content"), 0644)
		if err != nil {
			t.Fatal(err)
		}
	}

	notes, err := zet.ListNotes()
	if err != nil {
		t.Fatal(err)
	}

	if len(notes) != len(noteTitles) {
		t.Fatalf("Expected %d notes, got %d", len(noteTitles), len(notes))
	}

	expectedOrder := []string{"Apple Note", "Middle Note", "Zebra Note"}
	for i, note := range notes {
		if note.Title != expectedOrder[i] {
			t.Errorf("notes[%d].Title = %q, want %q", i, note.Title, expectedOrder[i])
		}
	}
}

func TestListNotesEmpty(t *testing.T) {
	zetDir := ZetDir(t)
	defer Cleanup(t, zetDir)

	originalZetDir := os.Getenv("ZETDIR")
	os.Setenv("ZETDIR", zetDir)
	defer func() {
		if originalZetDir != "" {
			os.Setenv("ZETDIR", originalZetDir)
		} else {
			os.Unsetenv("ZETDIR")
		}
	}()

	notes, err := zet.ListNotes()
	if err != nil {
		t.Fatal(err)
	}

	if len(notes) != 0 {
		t.Errorf("Expected 0 notes in empty directory, got %d", len(notes))
	}
}

func ZetDir(t *testing.T) string {
	dir, err := os.MkdirTemp(os.TempDir(), "zet-*")
	if err != nil {
		t.Fatal(err)
	}
	return dir
}

func Cleanup(t *testing.T, dir string) {
	err := os.RemoveAll(dir)
	if err != nil {
		t.Fatal(err)
	}
}
