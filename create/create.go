package create

import (
	"os"

	Z "github.com/rwxrob/bonzai/z"
	"github.com/rwxrob/help"
	"github.com/rwxrob/uniq-go"
)

func init() {
	Z.Vars.SoftInit()
}

var Cmd = &Z.Cmd{
	Name:    "create",
	Summary: "Create a new zettelkasten note",
	Commands: []*Z.Cmd{
		noteCmd,
		help.Cmd,
	},
}

var noteCmd = &Z.Cmd{
	Name:    "note",
	Summary: "Create a new zettelkasten note",
	Call: func(x *Z.Cmd, _ ...string) error {
		zetdir := Z.Vars.Get("zetdir")

		note_path := zetdir + "/notes/" + uniq.IsoSecond()

		// create a new folder for the note
		err := os.MkdirAll(note_path, 0755)
		if err != nil {
			return err
		}

		// create a new file for the note
		note_file, err := os.Create(note_path + "/README.md")

		snippets := Z.Vars.Get("snippets")
		// use the Read the contents of the snippet if it exists
		snippet, err := os.ReadFile(snippets + "/zet/note.md")
		if err == nil {
			// write the contents of the snippet to the note
			note_file.Write(snippet)
		}

		// use SysExec to open the note in the default editor
		editor, exists := os.LookupEnv("EDITOR")
		if !exists {
			editor = "vi"
		}
		Z.SysExec(editor, note_path+"/README.md")

		return nil
	},
}
