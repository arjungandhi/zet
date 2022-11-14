package zet

import (
	"fmt"
	"os"

	Z "github.com/rwxrob/bonzai/z"
	"github.com/rwxrob/help"
	"github.com/rwxrob/vars"

	"github.com/arjungandhi/zet/create"
	"github.com/arjungandhi/zet/title"
)

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

	snippets, exists := os.LookupEnv("SNIPPETS")
	if !exists {
		fmt.Println("SNIPPETS environment variable not set")
		os.Exit(1)
	}

	// set the zetdir var
	Z.Vars.Set("zetdir", zetdir)
	Z.Vars.Set("snippets", snippets)
}

// rootCmd is the main command for the zet command line tool
// its just holds all the other useful commands
var Cmd = &Z.Cmd{
	Name:    "zet",
	Summary: "zet is a command line tool for managing zettelkasten notes",
	Commands: []*Z.Cmd{
		help.Cmd,
		vars.Cmd,
		create.Cmd,
		title.Cmd,
	},
}
