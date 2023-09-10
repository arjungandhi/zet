// the zet package contians functions and structs for working with a zettelcasten.
// on the users system
package zet

import (
	"io/fs"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/rwxrob/uniq-go"
)

// Note is a struct that represents a note in the zettelkasten.
type Note struct {
	Title string
	Path  string
	Body  string
}

// createNote creates a new note in the zettelkasten.
func CreateNote(dir string, title string) (*Note, error) {
	// get a the current time
	isoSec := uniq.IsoSecond()

	// create a new directory for the note
	noteDir := filepath.Join(dir, isoSec)
	err := os.Mkdir(noteDir, 0755)
	if err != nil {
		return nil, err
	}

	// create a new file for the note
	noteFile := filepath.Join(noteDir, "README.md")
	err = ioutil.WriteFile(noteFile, []byte("# "+title), 0644)
	if err != nil {
		return nil, err
	}

	// create a new note struct
	return &Note{
		Title: title,
		Path:  noteFile,
		Body:  "# " + title,
	}, nil

}

// DeleteNote deletes a note from the zettelkasten.
func DeleteNote(dir string, note *Note) error {
	// get the note directory
	noteDir := filepath.Dir(note.Path)

	// remove the note directory
	err := os.RemoveAll(noteDir)
	if err != nil {
		return err
	}

	return nil
}

// ListNotes lists all the notes in the zettelkasten.
func ListNotes(dir string) ([]*Note, error) {
	// Walk the zettelkasten directory
	var notes []*Note

	err := filepath.WalkDir(dir,
		func(path string, d fs.DirEntry, err error) error {
			if err != nil {
				return err
			}

			// if the path is a directory, skip it
			if d.IsDir() {
				return nil
			}

			// if the path is a READEME.md file and the parent directory name is an int
			// then add it to the notes slice

			// get the parent directory
			parentDir := filepath.Base(filepath.Dir(path))

			if _, err = strconv.Atoi(parentDir); err == nil && d.Name() == "README.md" {
				// Read the Note, the title is the first line with out the '# ' prefix
				content, err := os.ReadFile(path)
				if err != nil {
					return err
				}
				title := strings.TrimPrefix(strings.Split(string(content), "\n")[0], "# ")

				// create a new note struct
				note := &Note{
					Title: title,
					Path:  path,
					Body:  string(content),
				}

				// add the note to the notes slice
				notes = append(notes, note)

			}

			return nil
		},
	)
	if err != nil {
		return nil, err
	}

	return notes, nil
}
