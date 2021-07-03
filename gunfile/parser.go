package gunfile

import (
	"fmt"
	"log"
	"os"

	"github.com/alecthomas/participle/v2"
)

func ParseGunfileString(s string) (*Gunfile, error) {
	parser, err := participle.Build(&Gunfile{})
	if err != nil {
		return nil, err
	}

	ast := &Gunfile{}
	err = parser.ParseString("", s, ast)
	if err != nil {
		return nil, err
	}

	return ast, nil
}

func ParseGunfile(filename string) (*Gunfile, error) {
	parser, err := participle.Build(&Gunfile{})
	if err != nil {
		return nil, err
	}

	ast := &Gunfile{}
	r, err := os.Open(filename)
	defer r.Close()
	if err != nil {
		return nil, err
	}
	if err := parser.Parse(filename, r, ast); err != nil {
		return nil, err
	}

	return ast, nil
}

func must(err error, s string) {
	if err != nil {
		log.Fatalf("%v: %v", s, err)
	}
}

func PrintAST(ast *Gunfile) {
	if ast == nil {
		fmt.Println("Empty AST")
		return
	}

	fmt.Printf("%v:%v -> Base: %s\n", ast.Pos.Line, ast.Pos.Column, ast.Base)

	fmt.Printf("Commands:\n")
	for _, cmd := range ast.Commands {
		fmt.Printf("    -> %v\n", cmd.Command)
	}
}

func PreprocessAST(ast *Gunfile) {
	// remove quotes
	ast.Base = ast.Base[1:len(ast.Base)-1]

	for i, cmd := range ast.Commands {
		// remove quotes
		ast.Commands[i].Command = cmd.Command[1:len(cmd.Command)-1]
	}
}
