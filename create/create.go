package create

import (
	"errors"
	"os"
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
		recipeCmd,
		help.Cmd,
	},
}

type cmdParams struct {
	Args    []string // positional args
	Snippet string   // snippet to use
}

var noteCmdArgs = cmdParams{
	Args:    []string{"title"},
	Snippet: "note",
}
var noteCmd = &Z.Cmd{
	Name:    "note",
	Summary: "Create a note in the zettelkasten",
	Commands: []*Z.Cmd{
		help.Cmd,
	},
	Call: noteCmdArgs.Call,
}

var urlCmdArgs = cmdParams{
	Args:    []string{"title", "url"},
	Snippet: "url",
}

var urlCmd = &Z.Cmd{
	Name:    "url",
	Summary: "Create a url in the zettelkasten",
	Commands: []*Z.Cmd{
		help.Cmd,
	},
	Call: urlCmdArgs.Call,
}

var recipeCmdArgs = cmdParams{
	Args:    []string{"title"},
	Snippet: "recipe",
}

var recipeCmd = &Z.Cmd{
	Name:    "recipe",
	Summary: "Create a recipe in the zettelkasten",
	Commands: []*Z.Cmd{
		help.Cmd,
	},
	Call: recipeCmdArgs.Call,
}

func (cmdParams *cmdParams) Call(x *Z.Cmd, _ ...string) error {
	// create a map to hold the values for the template
	varMap := map[string]interface{}{}

	for _, arg := range cmdParams.Args {
		// prompt for the value of the arg
		prompt := promptui.Prompt{
			Label: arg,
		}
		result, err := prompt.Run()
		if err != nil {
			return err
		}
		varMap[arg] = result
	}

	// fetch the path of useful vairables from Z.Vars
	zetdir := Z.Vars.Get(".zet.zetdir")
	snippets := Z.Vars.Get(".zet.snipdir")

	// create a unique filename
	notePath := zetdir + "/" + uniq.IsoSecond()

	// create a new folder for the note
	err := os.MkdirAll(notePath, 0755)
	if err != nil {
		return err
	}

	// create a new file for the note
	noteFile, err := os.Create(notePath + "/README.md")

	// use the Read the contents of the snippet if it exists
	snippet, err := os.ReadFile(snippets + "/zet/" + cmdParams.Snippet + ".md")
	if err == nil {
		// Run the snippet through the template engine
		template, err := template.New("note").Parse(string(snippet))
		if err != nil {
			return errors.New("Error parsing template")
		}
		template.Execute(noteFile, varMap)
	}

	// use SysExec to open the note in the default editor
	editor, exists := os.LookupEnv("EDITOR")
	if !exists {
		editor = "vi"
	}
	Z.SysExec(editor, notePath+"/README.md")

	return nil
}
