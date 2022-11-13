# 🌳 Zet a cli for managing my zettelkasten 

[![Go Report Card](https://goreportcard.com/badge/github.com/arjungandhi/zet?style=flat-square)](https://goreportcard.com/report/github.com/arjungandhi/zet)
[![Go Doc](https://img.shields.io/badge/godoc-reference-blue.svg?style=flat-square)](http://godoc.org/github.com/arjungandhi/zet)
[![PkgGoDev](https://pkg.go.dev/badge/github.com/arjungandhi/zet)](https://pkg.go.dev/github.com/arjungandhi/zet)
[![Release](https://img.shields.io/github/release/arjungandhi/zet.svg?style=flat-square)](https://github.com/arjungandhi/zet/releases/latest)

## Install

This command can be installed as a standalone program or composed into a
Bonzai command tree.

Standalone

```
go install github.com/arjungandi/zet/cmd/zet@latest
```

Composed

```go
package z

import (
	Z "github.com/rwxrob/bonzai/z"
	zet "github.com/arjungandhi/zet"
)

var Cmd = &Z.Cmd{
	Name:     `z`,
	Commands: []*Z.Cmd{help.Cmd, zet.Cmd},
}
```

## Tab Completion

To activate bash completion just use the `complete -C` option from your
`.bashrc` or command line. There is no messy sourcing required. All the
completion is done by the program itself.

```
complete -C zet zet
```

If you don't have bash or tab completion check use the shortcut
commands instead.

## Embedded Documentation

All documentation (like manual pages) has been embedded into the source
code of the application. See the source or run the program with help to
access it.

