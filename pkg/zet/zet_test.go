package zet_test

import (
	"os"
	"testing"
	"time"

	"github.com/arjungandhi/zet/pkg/zet"
)

func TestCreateNote(t *testing.T) {
	title := "My First Note"

	zetDir := ZetDir(t)

	note, err := zet.CreateNote(zetDir, title)
	if err != nil {
		t.Fatal(err)
	}

	content, err := os.ReadFile(note.Path)
	if err != nil {
		t.Fatal(err)
	}

	mdTitle := "# " + title
	if string(content) != mdTitle {
		t.Fatalf("Expected '%s' to be '%s'", content, mdTitle)
	}

	Cleanup(t, zetDir)
}

func TestDeleteNote(t *testing.T) {
	zetDir := ZetDir(t)
	note, err := zet.CreateNote(zetDir, "My First Note")
	if err != nil {
		t.Fatal(err)
	}

	err = zet.DeleteNote(zetDir, note)
	if err != nil {
		t.Fatal(err)
	}

	_, err = os.Stat(note.Path)
	if err == nil {
		t.Fatalf("Expected '%s' to not exist", note.Path)
	}

	Cleanup(t, zetDir)
}

func TestListNotes(t *testing.T) {
	zetDir := ZetDir(t)

	notesTitles := []string{
		"My First Note",
		"My Second Note",
		"My Third Note",
	}

	for _, title := range notesTitles {
		_, err := zet.CreateNote(zetDir, title)
		if err != nil {
			t.Fatal(err)
		}
		// make sure to sleep for a second so that the notes are created in
		// different seconds
		time.Sleep(1 * time.Second)
	}

	notes, err := zet.ListNotes(zetDir)
	if err != nil {
		t.Fatal(err)
	}

	if len(notes) != len(notesTitles) {
		t.Fatalf("Expected '%d' notes, got '%d'", len(notesTitles), len(notes))
	}

	for i, note := range notes {
		if note.Title != notesTitles[i] {
			t.Fatalf("Expected '%s' to be '%s'", note.Title, notesTitles[i])
		}
	}

}

// util test functions
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
