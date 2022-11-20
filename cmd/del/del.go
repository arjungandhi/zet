package del

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strconv"

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
			return errors.New("FileNum must be a number")
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

		// check if the user prompt for confirmation
		prompt := promptui.Prompt{
			Label:     fmt.Sprintf("Are you sure you want to delete %s", path),
			IsConfirm: true,
		}
		_, err = prompt.Run()
		if err != nil {
			return nil
		}

		// remove the file from the path
		dirPath := filepath.Dir(path)
		os.RemoveAll(dirPath)

		// delete the index from the zetlist
		delete(zetlistMap, indexNum)
		// marshal the zetlist
		b, err := json.Marshal(zetlistMap)
		if err != nil {
			return err
		}

		// set the zetlist
		Z.Vars.Set(".zet.list", string(b))

		return nil
	},
}

func mdTitle(path string) string {
	return ""
}
