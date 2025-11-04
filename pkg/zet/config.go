package zet

import (
	"fmt"
	"os"
)

func GetZetDir() (string, error) {
	dir := os.Getenv("ZETDIR")
	if dir == "" {
		return "", fmt.Errorf("ZETDIR environment variable not set")
	}
	return dir, nil
}

func GetEditor() string {
	editor := os.Getenv("EDITOR")
	if editor == "" {
		return "vi"
	}
	return editor
}

func GetRenderer() string {
	return "glow"
}
