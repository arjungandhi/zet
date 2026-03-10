package zet

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	bonzai "github.com/rwxrob/bonzai/z"
)

type Note struct {
	Title string
	Path  string
	Body  string
}

func CreateOrEditNote(title string) (string, error) {
	dir, err := GetZetDir()
	if err != nil {
		return "", err
	}

	sanitized := SanitizeFilename(title)
	path := filepath.Join(dir, sanitized+".md")

	if !NoteExists(dir, title) {
		err = WriteNote(dir, title, "")
		if err != nil {
			return "", err
		}
	}

	return path, nil
}

func JournalNote(dateStr string) (string, error) {
	dir, err := GetZetDir()
	if err != nil {
		return "", err
	}

	var date string
	if dateStr == "" {
		date = time.Now().Format("2006-01-02")
	} else {
		t, err := time.Parse("2006-01-02", dateStr)
		if err != nil {
			return "", fmt.Errorf("invalid date %q, expected YYYY-MM-DD format", dateStr)
		}
		date = t.Format("2006-01-02")
	}

	path := filepath.Join(dir, "journal.md")

	// Create file with top-level header if it doesn't exist
	if _, err := os.Stat(path); os.IsNotExist(err) {
		if err := os.WriteFile(path, []byte("# Journal\n\n"), 0644); err != nil {
			return "", err
		}
	}

	// Add a dated entry header near the top if one doesn't already exist
	content, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}

	header := "## " + date
	text := string(content)
	if !strings.Contains(text, header) {
		// Insert after the "# Journal" line so newest entries are at the top
		entry := header + "\n\n"
		idx := strings.Index(text, "\n")
		if idx != -1 {
			text = text[:idx+1] + "\n" + entry + text[idx+1:]
		} else {
			text = text + "\n" + entry
		}
		if err := os.WriteFile(path, []byte(text), 0644); err != nil {
			return "", err
		}
	}

	return path, nil
}

func DeleteNote(note *Note) error {
	return os.Remove(note.Path)
}

func ListNotes() ([]*Note, error) {
	dir, err := GetZetDir()
	if err != nil {
		return nil, err
	}

	files, err := ListNoteFiles(dir)
	if err != nil {
		return nil, err
	}

	var notes []*Note
	for _, path := range files {
		note, err := ReadNote(path)
		if err != nil {
			continue
		}
		notes = append(notes, note)
	}

	return notes, nil
}

func OpenNote(searchTerm string) error {
	notes, err := ListNotes()
	if err != nil {
		return err
	}

	note, err := FindNote(notes, searchTerm)
	if err != nil {
		return err
	}

	editor := GetEditor()
	return bonzai.SysExec(editor, note.Path)
}

func RenderNote(searchTerm string) error {
	notes, err := ListNotes()
	if err != nil {
		return err
	}

	note, err := FindNote(notes, searchTerm)
	if err != nil {
		return err
	}

	renderer := GetRenderer()
	cmd := exec.Command(renderer, note.Path)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}
