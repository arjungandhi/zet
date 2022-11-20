package edit

import (
	"encoding/json"
	"errors"
	"os"
	"strconv"

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
	The edit command allows you to edit a note in the zettelkasten. It uses
	the file numbers from the search and title commands to identify the note
	to edit.
	`,
	MinArgs: 1,
	MaxArgs: 1,
	Commands: []*Z.Cmd{
		help.Cmd,
	},
	Call: func(x *Z.Cmd, args ...string) error {
		index := args[0]
		// convert the index to number
		indexNum, err := strconv.Atoi(index)
		if err != nil {
			return errors.New("File Num must be a number")
		}

		zetlist := Z.Vars.Get(".zet.list")
		// unmarshal the zetlist
		var zetlistMap map[int]string
		err = json.Unmarshal([]byte(zetlist), &zetlistMap)
		if err != nil {
			return err
		}

		// get the path of the file
		path, ok := zetlistMap[indexNum]
		if !ok {
			return errors.New("No file found with that index, make sure you've run 'zet search' or 'zet titles' first")
		}

		// open the file
		// use SysExec to open the note in the default editor
		editor, exists := os.LookupEnv("EDITOR")
		if !exists {
			editor = "vi"
		}
		Z.SysExec(editor, path)

		return nil
	},
}
