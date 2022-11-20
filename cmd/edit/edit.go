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
	Usage:   "edit file",
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
			return errors.New("index must be a number")
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
			return errors.New("index not found")
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

func mdTitle(path string) string {
	return ""
}
