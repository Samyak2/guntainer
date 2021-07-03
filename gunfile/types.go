package gunfile

import (
	"github.com/alecthomas/participle/v2/lexer"
)

type Gunfile struct {
	Pos lexer.Position

	// the base container
	Base string `parser:"'Using' @String"`
	// a bunch of commands that are run sequentially inside the container
	Commands []*Exec `parser:"@@*"`
	// MainCmd *Main `parser:"@@"`
}

// Exec commands are run during build time
// Example:
// Exec "echo 'hello world'"
type Exec struct {
	Command string `parser:"'Exec' @String"`
}

// the main command is run when the container starts
// TODO: figure out a way to pack extra information with the image
// Example:
// Main "/bin/bash"
// type Main struct {
// 	Command string `parser:"'Main' @String"`
// }
