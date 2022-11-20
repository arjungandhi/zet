// zet Command Line tool
package main

import (
	"fmt"
	"os"

	Z "github.com/rwxrob/bonzai/z"
	"github.com/rwxrob/help"
	"github.com/rwxrob/vars"

	"github.com/arjungandhi/zet/cmd/create"
	"github.com/arjungandhi/zet/cmd/del"
	"github.com/arjungandhi/zet/cmd/edit"
	"github.com/arjungandhi/zet/cmd/search"
	"github.com/arjungandhi/zet/cmd/title"
)

func main() {
	zetCmd.Run()
}

// init runs immediately after the package is loaded.
// it sets up a few things that are needed for the command to run nicely
func init() {
	// check for the existence of $ZETDIR
	Z.Vars.SoftInit()
	// set the zet directory
	zetdir, exists := os.LookupEnv("ZETDIR")
	if !exists {
		fmt.Println("ZETDIR environment variable not set")
		os.Exit(1)
	}

	snipdir, exists := os.LookupEnv("SNIPPETS")
	if !exists {
		fmt.Println("SNIPPETS environment variable not set")
		os.Exit(1)
	}

	// set the zetdir var
	Z.Vars.Set(".zet.zetdir", zetdir)
	Z.Vars.Set(".zet.snipdir", snipdir)
}

// rootCmd is the main command for the zet command line tool
// its just holds all the other useful commands
var zetCmd = &Z.Cmd{
	Name:    "zet",
	Summary: "zet is a command line tool for managing zettelkasten notes",
	Commands: []*Z.Cmd{
		help.Cmd,
		vars.Cmd,
		create.Cmd,
		title.Cmd,
		search.Cmd,
		edit.Cmd,
		del.Cmd,
	},
}
