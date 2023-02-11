package del

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/arjungandhi/zet/cmd/search"
	"github.com/manifoldco/promptui"
	Z "github.com/rwxrob/bonzai/z"
	"github.com/rwxrob/help"
)

func init() {
	Z.Vars.SoftInit()
}

var Cmd = &Z.Cmd{
	Name:    "del",
	Summary: "delete a zettelkasten note",
	Usage:   "del [File Num]",
	Description: `
	The del command allows you to delete a note in the zettelkasten. It uses
	the file numbers from the search and title commands to identify the note
	to delete.
	`,
	MinArgs: 0,
	MaxArgs: 1,
	Commands: []*Z.Cmd{
		help.Cmd,
	},
	Call: func(x *Z.Cmd, args ...string) error {
		query := ""
		if len(args) > 0 {
			query = args[0]
		}

		n, err := search.Search(query)
		if err != nil {
			return err
		}

		// check if the user prompt for confirmation
		prompt := promptui.Prompt{
			Label:     fmt.Sprintf("Are you sure you want to delete %s", n.Title()),
			IsConfirm: true,
		}
		_, err = prompt.Run()
		if err != nil {
			return nil
		}

		// remove the file from the path
		dirPath := filepath.Dir(n.Path())
		os.RemoveAll(dirPath)

		return nil
	},
}
