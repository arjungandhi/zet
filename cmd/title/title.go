package title

import (
	"fmt"

	Z "github.com/rwxrob/bonzai/z"
	"github.com/rwxrob/help"

	"github.com/arjungandhi/zet/pkg/node"
)

func init() {
	Z.Vars.SoftInit()
}

var Cmd = &Z.Cmd{
	Name:    "titles",
	Summary: "list all the titles from the Zettelkasten",
	Commands: []*Z.Cmd{
		help.Cmd,
	},
	Call: func(x *Z.Cmd, args ...string) error {
		titles := Titles()
		for _, t := range titles {
			fmt.Println(t)
		}
		return nil
	},
}

func Titles() []string {
	zetDir := Z.Vars.Get(".zet.zetdir")
	nodes := node.LoadAll(zetDir)
	titles := []string{}
	for _, n := range nodes {
		titles = append(titles, n.Title())
	}
	return titles
}
