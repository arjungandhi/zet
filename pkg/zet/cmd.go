package zet

import (
	"fmt"
	"strings"

	bonzai "github.com/rwxrob/bonzai/z"
	"github.com/rwxrob/help"
)

var Cmd = &bonzai.Cmd{
	Name:     "zet",
	Commands: []*bonzai.Cmd{help.Cmd, listCmd, deleteCmd, newCmd, renderCmd},
	Call: func(cmd *bonzai.Cmd, args ...string) error {
		search := strings.Join(args, " ")
		return OpenNote(search)
	},
}

var deleteCmd = &bonzai.Cmd{
	Name: "delete",
	Call: func(cmd *bonzai.Cmd, args ...string) error {
		search := strings.Join(args, " ")

		notes, err := ListNotes()
		if err != nil {
			return err
		}

		note, err := FindNote(notes, search)
		if err != nil {
			return err
		}

		err = DeleteNote(note)
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
		title := strings.Join(args, " ")
		if title == "" {
			return fmt.Errorf("title required")
		}

		path, err := CreateOrEditNote(title)
		if err != nil {
			return err
		}

		editor := GetEditor()
		return bonzai.SysExec(editor, path)
	},
}

var listCmd = &bonzai.Cmd{
	Name: "list",
	Call: func(cmd *bonzai.Cmd, args ...string) error {
		notes, err := ListNotes()
		if err != nil {
			return err
		}

		for _, note := range notes {
			fmt.Println(note.Title)
		}

		return nil
	},
}

var renderCmd = &bonzai.Cmd{
	Name: "render",
	Call: func(cmd *bonzai.Cmd, args ...string) error {
		search := strings.Join(args, " ")
		return RenderNote(search)
	},
}
