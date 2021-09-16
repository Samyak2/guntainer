package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/Samyak2/guntainer/gunfile"
	"github.com/Samyak2/guntainer/guntar"
	// "github.com/mholt/archiver/v3"
)

func build() {
	if len(os.Args) <= 2 {
		log.Fatalln("Gib Gunfile")
	}
	if len(os.Args) <= 3 {
		log.Fatalln("Gib output image path")
	}

	filename := os.Args[2]
	outputPath := os.Args[3]

	ast, err := gunfile.ParseGunfile(filename)
	if err != nil {
		log.Fatalln("Error in reading Gunfile:", err)
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

	// err = archiver.Archive([]string{dname}, outputPath)
	err = guntar.ArchiveDirectory(dname, outputPath)
	if err != nil {
		log.Fatalln(err)
	}
}
