package man2md

import (
	"os"
	"testing"
)

func TestMan2md(t *testing.T) {
	manPath := "/usr/share/man/man1/xargs.1"
	mdPath := "xargs.md"

	var err error
	var manFile *os.File
	var mdFile *os.File

	if manFile, err = os.Open(manPath); err != nil {
		t.Fatal("Error opening man file: ", err)
	}

	if mdFile, err = os.Create(mdPath); err != nil {
		t.Fatal("Error creating md file: ", err)
	}

	if err = man2md(manFile, mdFile); err != nil {
		t.Fatal("man2md returned error: ", err)
	}
}
