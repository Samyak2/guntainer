package main

import (
	"os"

	"github.com/Samyak2/guntainer/cmd"
	"github.com/Samyak2/guntainer/core"
)

func main() {
	// these need to be handled specially because:
	//  1. they are internal subcommands, not to be used in the CLI
	//  2. environment variables may not be set when these are called,
	//     leading to viper erroring out because it cannot find $HOME
	if len(os.Args) > 1 {
		switch os.Args[1] {
			case "child":
				core.Child(os.Args[2], os.Args[3:])
				return
			case "childMultiple":
				core.ChildMultiple(os.Args[2], os.Args[3:])
				return
		}
	}

	cmd.Execute()
}
