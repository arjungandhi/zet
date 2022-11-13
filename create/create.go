package create

import (
	"errors"
	"os"
	"strings"
	"text/template"

	"github.com/manifoldco/promptui"
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
		urlCmd,
		help.Cmd,
	},
}

var noteCmd = &Z.Cmd{
	Name:    "note",
	Summary: "Create a new zettelkasten note",
	Usage:   "create note [title]",
	Commands: []*Z.Cmd{
		help.Cmd,
	},
	Call: func(x *Z.Cmd, args ...string) error {
		type Note struct {
			Title string
		}
		note := &Note{}

		if len(args) > 0 {
			note.Title = strings.Join(args, " ")
		} else {
			// Prompt for title
			prompt := promptui.Prompt{
				Label: "Title",
			}
			result, err := prompt.Run()
			if err != nil {
				return err
			}

			note.Title = result
		}

		zetdir := Z.Vars.Get("zetdir")

		notePath := zetdir + "/notes/" + uniq.IsoSecond()

		// create a new folder for the note
		err := os.MkdirAll(notePath, 0755)
		if err != nil {
			return err
		}

		// create a new file for the note
		noteFile, err := os.Create(notePath + "/README.md")

		snippets := Z.Vars.Get("snippets")
		// use the Read the contents of the snippet if it exists
		snippet, err := os.ReadFile(snippets + "/zet/note.md")
		if err == nil {
			// Run the snippet through the template engine
			template, err := template.New("note").Parse(string(snippet))
			if err != nil {
				return errors.New("Error parsing template")
			}
			template.Execute(noteFile, note)
		}

		// use SysExec to open the note in the default editor
		editor, exists := os.LookupEnv("EDITOR")
		if !exists {
			editor = "vi"
		}
		Z.SysExec(editor, notePath+"/README.md")

		return nil
	},
}

var urlCmd = &Z.Cmd{
	Name:    "url",
	Summary: "Create a new zettelkasten url",
	Usage:   "[url]",
	Commands: []*Z.Cmd{
		help.Cmd,
	},
	MaxArgs: 1,
	Call: func(x *Z.Cmd, args ...string) error {
		type UrlNote struct {
			Title string
			Url   string
		}

		urlFields := &UrlNote{}

		// check if the url is passed as an argument
		if len(args) > 0 {
			urlFields.Url = args[0]
		} else {
			// Prompt for the url
			prompt := promptui.Prompt{
				Label: "Url",
			}
			result, err := prompt.Run()
			if err != nil {
				return err
			}
			urlFields.Url = result
		}

		// Prompt for the title
		prompt := promptui.Prompt{
			Label: "Title",
		}
		result, err := prompt.Run()
		if err != nil {
			return err
		}
		urlFields.Title = result

		zetdir := Z.Vars.Get("zetdir")

		urlPath := zetdir + "/notes/" + uniq.IsoSecond()

		// create a new folder for the note
		err = os.MkdirAll(urlPath, 0755)
		if err != nil {
			return err
		}

		// create a new file for the url
		urlNote, err := os.Create(urlPath + "/README.md")

		snippets := Z.Vars.Get("snippets")
		// use the Read the contents of the snippet if it exists
		snippet, err := os.ReadFile(snippets + "/zet/url.md")
		if err == nil {
			// Load the snippet into the template engine
			template, err := template.New("url").Parse(string(snippet))
			if err != nil {
				return errors.New("Error parsing template")
			}
			template.Execute(urlNote, urlFields)
		}

		// use SysExec to open the note in the default editor
		editor, exists := os.LookupEnv("EDITOR")
		if !exists {
			editor = "vi"
		}
		Z.SysExec(editor, urlPath+"/README.md")

		return nil
	},
}
