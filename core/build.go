package core

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/Samyak2/guntainer/gunfile"
	"github.com/Samyak2/guntainer/guntar"
)

func Build(gunfileName string, outputPath string) {
	ast, err := gunfile.ParseGunfile(gunfileName)
	if err != nil {
		log.Fatalln("Error in reading Gunfile: ", err)
	}
	gunfile.PreprocessAST(ast)

	gunfile.PrintAST(ast)

	divider, err := strconv.Unquote("'\u0000'")
	if err != nil {
		log.Fatalf("Divider could not be unquoted. If you're seeing this error, uhhh. You shoudn't be.")
	}

	cmdWithArgs := make([]string, len(ast.Commands))
	for i, cmd := range ast.Commands {
		newCommand := strings.Join([]string{"/bin/sh", "-c", fmt.Sprintf("%v", cmd.Command)}, divider)

		cmdWithArgs[i] = newCommand
	}

	dname, err := RunMultipleNoClean(ast.Base, cmdWithArgs)
	if dname != "" {
		defer os.RemoveAll(dname)
	}

	if err != nil {
		log.Fatalln(err)
	}

	err = guntar.ArchiveDirectory(dname, outputPath)
	if err != nil {
		log.Fatalln(err)
	}
}
