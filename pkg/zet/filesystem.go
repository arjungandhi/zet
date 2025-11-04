package zet

import (
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
)

func SanitizeFilename(title string) string {
	re := regexp.MustCompile(`[^a-zA-Z0-9 ]`)
	sanitized := re.ReplaceAllString(title, "")
	sanitized = strings.Join(strings.Fields(sanitized), " ")
	return strings.TrimSpace(sanitized)
}

func NoteExists(dir, title string) bool {
	sanitized := SanitizeFilename(title)
	path := filepath.Join(dir, sanitized+".md")
	_, err := os.Stat(path)
	return err == nil
}

func ReadNote(path string) (*Note, error) {
	content, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	filename := filepath.Base(path)
	title := strings.TrimSuffix(filename, ".md")

	return &Note{
		Title: title,
		Path:  path,
		Body:  string(content),
	}, nil
}

func WriteNote(dir, title, content string) error {
	sanitized := SanitizeFilename(title)
	path := filepath.Join(dir, sanitized+".md")
	return os.WriteFile(path, []byte(content), 0644)
}

func ListNoteFiles(dir string) ([]string, error) {
	pattern := filepath.Join(dir, "*.md")
	files, err := filepath.Glob(pattern)
	if err != nil {
		return nil, err
	}

	sort.Strings(files)
	return files, nil
}
