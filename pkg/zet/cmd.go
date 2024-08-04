package zet

import (
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"

	bonzai "github.com/rwxrob/bonzai/z"
	"github.com/rwxrob/help"
)

func init() {
}

// exported main command
var Cmd = &bonzai.Cmd{
	Name:     "zet",
	Commands: []*bonzai.Cmd{help.Cmd, listCmd, deleteCmd, newCmd, renderCmd},
	Call: func(cmd *bonzai.Cmd, args ...string) error {
		zetDir, err := getZetDir()
		if err != nil {
			return err
		}

		// args is an optional search term to feed into fzf
		search := ""
		if len(args) > 0 {
			search = strings.Join(args, " ")
		}

		// list all the notes in our zet directory
		notes, err := ListNotes(zetDir)
		if err != nil {
			return err
		}

		// use fzf to find the note we want
		note, err := findNote(notes, search)
		if err != nil {
			return err
		}

		// open the note in our editor
		editNote(note)
		return nil
	},
}

var deleteCmd = &bonzai.Cmd{
	Name: "delete",
	Call: func(cmd *bonzai.Cmd, args ...string) error {
		zetDir, err := getZetDir()
		if err != nil {
			return err
		}

		// argsis an optional search term to feed into fzf
		search := ""
		if len(args) > 0 {
			search = strings.Join(args, " ")
		}

		// list all the notes in our zet directory
		notes, err := ListNotes(zetDir)
		if err != nil {
			return err
		}

		// use fzf to find the note we want
		note, err := findNote(notes, search)
		if err != nil {
			return err
		}

		// delete the note
		err = DeleteNote(zetDir, note)
		if err != nil {
			return err
		}

		fmt.Printf("Deleted %s @ %s\n", note.Title, note.Path)
		return nil
	},
}

var newCmd = &bonzai.Cmd{
	Name: "new",
	Call: func(cmd *bonzai.Cmd, args ...string) error {
		zetDir, err := getZetDir()
		if err != nil {
			return err
		}

		// args is the title of the note
		title := ""
		if len(args) > 0 {
			title = strings.Join(args, " ")
		}

		// create the note
		note, err := CreateNote(zetDir, title)
		if err != nil {
			return err
		}

		// open the note in our editor
		editNote(note)
		return nil
	},
}

var listCmd = &bonzai.Cmd{
	Name: "list",
	Call: func(cmd *bonzai.Cmd, args ...string) error {
		zetDir, err := getZetDir()
		if err != nil {
			return err
		}

		// list all the notes in our zet directory
		notes, err := ListNotes(zetDir)
		if err != nil {
			return err
		}

		// print the notes
		for _, note := range notes {
			fmt.Println(note.Title)
		}

		return nil
	},
}

var renderCmd = &bonzai.Cmd{
	Name: "render",
	Call: func(cmd *bonzai.Cmd, args ...string) error {
		zetDir, err := getZetDir()
		if err != nil {
			return err
		}

		// args is an optional search term to feed into fzf
		search := ""
		if len(args) > 0 {
			search = strings.Join(args, " ")
		}

		// list all the notes in our zet directory
		notes, err := ListNotes(zetDir)
		if err != nil {
			return err
		}

		// use fzf to find the note we want
		note, err := findNote(notes, search)
		if err != nil {
			return err
		}

		// render the note
		err = renderNote(note)
		if err != nil {
			return err
		}

		return nil
	},
}

// render uses glow to render a note
func renderNote(n *Note) error {
	cmd := exec.Command("glow", n.Path)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func editNote(n *Note) {
	editor, exists := os.LookupEnv("EDITOR")
	if !exists {
		editor = "vi"
	}

	bonzai.SysExec(editor, n.Path)
}

// get zet dir gets the zet directory from the ZETDIR environment variable
func getZetDir() (string, error) {
	if value, ok := os.LookupEnv("ZETDIR"); ok {
		return value, nil
	}
	return "", fmt.Errorf("ZETDIR environment variable not set")
}

// findNote takes a list of notes and a search term and returns the note that
// matches the search term
// TODO: I wanna write a more sane search wrapper around fzf at some point
func findNote(notes []*Note, search string) (*Note, error) {
	fzf_args := []string{
		fmt.Sprintf("--query=%s", search),
		"--delimiter=\t",
		"--with-nth=2",
		"--layout=reverse",
		"-1",
		"--preview=cat {3}",
	}

	// generate our search list
	var searchOpts []string
	for i, note := range notes {
		searchOpts = append(searchOpts, fmt.Sprintf("%d\t%s\t%s", i, note.Title, note.Path))
	}

	// use fzf to find the note we want
	cmd := exec.Command("fzf", fzf_args...)
	cmd.Stdin = strings.NewReader(strings.Join(searchOpts, "\n"))
	cmd.Stderr = os.Stderr

	selected, err := cmd.Output()
	if err != nil {
		return nil, err
	}
	selectedIndex := strings.Split(string(selected), "\t")[0]

	// convert the selected string to an int
	i, err := strconv.Atoi(selectedIndex)
	if err != nil {
		return nil, err
	}

	return notes[i], nil
}
