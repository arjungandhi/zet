// zet Command Line tool
package main

import (
	"github.com/arjungandhi/zet/pkg/zet"
)

// init runs immediately after the package is loaded.
// it sets up a few things that are needed for the command to run nicely
func main() {
	zet.Cmd.Run()
}
