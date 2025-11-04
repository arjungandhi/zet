package zet

import (
	"os"
	"os/exec"
	"path/filepath"

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
