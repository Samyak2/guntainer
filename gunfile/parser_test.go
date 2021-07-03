package gunfile

import (
	"testing"
)

func basicParsingHelper(testString string, t *testing.T) {
	ast, err := ParseGunfileString(testString)

	if err != nil {
		t.Errorf("Error parsing: %v", err)
	}

	PreprocessAST(ast)

	PrintAST(ast)
}

func TestBasicParsing1(t *testing.T) {
	testString := `Using "somebaseimage"`
	basicParsingHelper(testString, t)
}

func TestBasicParsing2(t *testing.T) {
	testString := `
Using	"somebaseimage"
Exec	"apt update"
`
	basicParsingHelper(testString, t)
}

func TestBasicParsing3(t *testing.T) {
	testString := `
Using "somebaseimage"
Exec "apt update"
Exec "echo hi > somefile"
`
	basicParsingHelper(testString, t)
}
