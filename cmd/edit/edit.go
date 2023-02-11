package edit

import (
	"os"

	"github.com/arjungandhi/zet/cmd/search"
	Z "github.com/rwxrob/bonzai/z"
	"github.com/rwxrob/help"
)

func init() {
	Z.Vars.SoftInit()
}

var Cmd = &Z.Cmd{
	Name:    "edit",
	Summary: "edit a zettelkasten note",
	Usage:   "edit [File Num]",
	Description: `
	The edit command allows you to edit a note in the zettelkasten. 
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

		// use SysExec to open the note in the default editor
		editor, exists := os.LookupEnv("EDITOR")
		if !exists {
			editor = "vi"
		}
		Z.SysExec(editor, n.Path())

		return nil
	},
}
