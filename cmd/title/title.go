package title

import (
	"fmt"
	"os"
	"path/filepath"

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
		zetdir := Z.Vars.Get(".zet.zetdir")
		// What this function does is walk the filesystem and dispatch various files
		// to functions to get the title of the file. The title is then printed
		// to stdout.
		filepath.Walk(zetdir, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}

			if info.IsDir() {
				return nil
			}

			ext := filepath.Ext(path)
			switch ext {
			case ".md":
				n := node.MarkdownNode{}
				err := n.Load(path)
				if err == nil {
					title := n.Title()
					if title != "" {
						fmt.Println(title)
					}
				}
			}
			return nil
		})
		return nil
	},
}

func mdTitle(path string) string {
	return ""
}
