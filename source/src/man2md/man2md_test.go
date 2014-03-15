package man2md

import (
	"os"
	"testing"
)

func TestMan2md_ParseUsingConvertFunc(t *testing.T) {
	manPath := "testing/xargs.1"
	mdPath := "testing/xargs.md"

	var err error
	var manFile *os.File
	var mdFile *os.File

	if manFile, err = os.Open(manPath); err != nil {
		t.Fatal("Error opening man file: ", err)
	}
	defer manFile.Close()

	if mdFile, err = os.Create(mdPath); err != nil {
		t.Fatal("Error creating md file: ", err)
	}
	defer mdFile.Close()

	if err = Convert(manFile, mdFile); err != nil {
		t.Fatal("man2md returned error: ", err)
	}
}

func TestMan2md_ParseUsingParserObject(t *testing.T) {
	manPath := "testing/xargs.1"
	mdPath := "testing/xargs.md"

	var err error
	var manFile *os.File
	var mdFile *os.File

	if manFile, err = os.Open(manPath); err != nil {
		t.Fatal("Error opening man file: ", err)
	}
	defer manFile.Close()

	if mdFile, err = os.Create(mdPath); err != nil {
		t.Fatal("Error creating md file: ", err)
	}
	defer mdFile.Close()

	parser := NewParser(manFile, mdFile)
	if err = parser.Parse(); err != nil {
		t.Fatal("man2md returned error: ", err)
	}
}
