package zet

import (
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

func FindNote(notes []*Note, searchTerm string) (*Note, error) {
	if len(notes) == 0 {
		return nil, fmt.Errorf("no notes to search")
	}

	input := BuildFzfInput(notes)

	args := []string{
		"--delimiter=\t",
		"--with-nth=2",
		"--layout=reverse",
		"-1",
		"--preview=cat {3}",
	}

	if searchTerm != "" {
		args = append([]string{fmt.Sprintf("--query=%s", searchTerm)}, args...)
	}

	cmd := exec.Command("fzf", args...)
	cmd.Stdin = strings.NewReader(input)
	cmd.Stderr = os.Stderr

	output, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	return ParseFzfOutput(notes, string(output))
}

func BuildFzfInput(notes []*Note) string {
	var lines []string
	for i, note := range notes {
		line := fmt.Sprintf("%d\t%s\t%s", i, note.Title, note.Path)
		lines = append(lines, line)
	}
	return strings.Join(lines, "\n")
}

func ParseFzfOutput(notes []*Note, output string) (*Note, error) {
	output = strings.TrimSpace(output)
	if output == "" {
		return nil, fmt.Errorf("empty fzf output")
	}

	parts := strings.Split(output, "\t")
	if len(parts) < 1 {
		return nil, fmt.Errorf("malformed fzf output")
	}

	index, err := strconv.Atoi(parts[0])
	if err != nil {
		return nil, fmt.Errorf("invalid index in fzf output: %w", err)
	}

	if index < 0 || index >= len(notes) {
		return nil, fmt.Errorf("index out of range: %d", index)
	}

	return notes[index], nil
}
