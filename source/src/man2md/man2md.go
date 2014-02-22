// Copyright 2014 The Gman Authors. All rights reserved.
// Use of this source code is governed by an MIT-style
// license that can be found in LICENSE file.

// man2md
// Package for converting man (*roff) pages to Mandown format.

package man2md

import (
	"bufio"
	"fmt"
	"io"
	"sort"
	"strings"

//	"unicode"
//	"unicode/utf8"
)

type man2mdState struct {
	pageName        string
	bof             bool
	unprocessedCmds []string // dsn debug
}

func man2md(reader io.Reader, writer io.Writer) (err error) {
	bufferedReader := bufio.NewReader(reader)
	bufferedWriter := bufio.NewWriter(writer)

	var state man2mdState
	state.bof = true

	eof := false
	for !eof {
		line, err := bufferedReader.ReadString('\n')
		if err == io.EOF {
			err = nil
			eof = true
		} else if err != nil {
			return err
		}

		if len(line) > 0 {
			if line[0] == '.' {
				// Dot command
				if err := processDotCommand(&state, line, bufferedWriter); err != nil {
					return err
				}
			} else {
				// Write the line, replacing trailing "\n" with a space.
				if _, err := bufferedWriter.WriteString(line[:len(line)-1]); err != nil {
					return err
				}
				if _, err := bufferedWriter.WriteString(" "); err != nil {
					return err
				}
			}
		}
	}

	// dsn debug
	// Dump the list of unprocessed commands.
	fmt.Printf("\n%d unprocessed dot commands:\n", len(state.unprocessedCmds))
	for _, cmd := range state.unprocessedCmds {
		fmt.Printf("\t%s\n", cmd)
	}
	// end dsn debug

	return err
}

// Process lines beginning with '.'.
// Initially this will handle mdoc macros. Eventually we'll need to also process mandoc.
func processDotCommand(state *man2mdState, line string, bufferedWriter *bufio.Writer) (err error) {
	// Split the line into tokens.
	cmdToks := strings.Split(line, " ")

	// Strip "\n" from the end.
	lastCmdTokenIndex := len(cmdToks) - 1
	cmdToks[lastCmdTokenIndex] = cmdToks[lastCmdTokenIndex][:len(cmdToks[lastCmdTokenIndex])-1]

	cmd := cmdToks[0]

	switch cmd {
	case ".\\\"":
		// Comment. Ignore it.

	case ".Ar":
		// Argument
		// Emphasize it (wrap in "_").
		if _, err := bufferedWriter.WriteString("_"); err != nil {
			return err
		}
		if _, err := bufferedWriter.WriteString(cmdToks[1]); err != nil {
			return err
		}
		if _, err := bufferedWriter.WriteString("_ "); err != nil {
			return err
		}

	case ".Fl":
		// Flag
		// Prefix with '-'.
		//needSpace := true
		if _, err := bufferedWriter.WriteString("-"); err != nil {
			return err
		}

		// Process the remainder of the line.
		if err := processTokens(state, cmdToks[1:], bufferedWriter); err != nil {
			return err
		}

		/*
			for token := 1; token < len(cmdToks); token++ {
				outputHyphen := false
				if token > 1 {
					r, _ := utf8.DecodeRuneInString(cmdToks[token])
					if !unicode.IsPunct(r) {
						outputHyphen = true
					}
				}

				if outputHyphen {
					if _, err := bufferedWriter.WriteString("-"); err != nil {
						return err
					}
				}
				if _, err := bufferedWriter.WriteString(cmdToks[token]); err != nil {
					return err
				}
				if _, err := bufferedWriter.WriteString(" "); err != nil {
					return err
				}
				needSpace = false
			}

			if needSpace {
				if _, err := bufferedWriter.WriteString(" "); err != nil {
					return err
				}
			}
		*/

	case ".Nm":
		// Page name macro
		// The first occurrence sets the name.
		tokensProcessed := 1
		if len(state.pageName) == 0 {
			state.pageName = cmdToks[1]
			tokensProcessed = 2
		}

		if len(state.pageName) > 0 {
			if _, err := bufferedWriter.WriteString(state.pageName); err != nil {
				return err
			}
			if _, err := bufferedWriter.WriteString(" "); err != nil {
				return err
			}
			if len(cmdToks) > tokensProcessed {
				for tokenIndex := tokensProcessed; tokenIndex < len(cmdToks); tokenIndex++ {
					if _, err := bufferedWriter.WriteString(cmdToks[tokenIndex]); err != nil {
						return err
					}
					if _, err := bufferedWriter.WriteString(" "); err != nil {
						return err
					}
				}
			}
		}

	case ".Oc":
		// End multi-line optional section.
		if _, err := bufferedWriter.WriteString("] "); err != nil {
			return err
		}

	case ".Oo":
		// Begin multi-line optional section (enclose in "[").
		// This will be termintated by ".Oc".
		if _, err := bufferedWriter.WriteString("["); err != nil {
			return err
		}

	case ".Op":
		// Optional part of command line. Enclose in "[]".
		if _, err := bufferedWriter.WriteString("["); err != nil {
			return err
		}
		if err := processTokens(state, cmdToks[1:], bufferedWriter); err != nil {
			return err
		}
		if _, err := bufferedWriter.WriteString("] "); err != nil {
			return err
		}

	case ".Pp", ".PP":
		// Paragraph break
		if _, err := bufferedWriter.WriteString("\n\n"); err != nil {
			return err
		}

	case ".Sh", ".SH":
		// Section heading
		if !state.bof {
			if _, err := bufferedWriter.WriteString("\n\n"); err != nil {
				return err
			}
		} else {
			state.bof = false
		}

		header := cmdToks[1]
		if _, err := bufferedWriter.WriteString(header); err != nil {
			return err
		}
		if _, err := bufferedWriter.WriteString("\n"); err != nil {
			return err
		}
		ul := strings.Repeat("=", len(header))
		if _, err := bufferedWriter.WriteString(ul); err != nil {
			return err
		}
		if _, err := bufferedWriter.WriteString("\n"); err != nil {
			return err
		}

	default:
		// dsn debug
		// Add to our array tracking unprocessed commands.
		index := sort.SearchStrings(state.unprocessedCmds, cmd)
		if index >= len(state.unprocessedCmds) {
			// Add cmd to the end of the list.
			state.unprocessedCmds = append(state.unprocessedCmds, cmd)
		} else if state.unprocessedCmds[index] != cmd {
			// Insert cmd into the list.
			state.unprocessedCmds = append(state.unprocessedCmds[:index],
				append([]string{cmd}, state.unprocessedCmds[index:]...)...)
		}
		// end dsn debug
	}

	return err
}

const (
	MAN2MD_AR = iota
	MAN2MD_OP
	MAN2MD_FL
)

func processTokens(state *man2mdState, cmdToks []string, bufferedWriter *bufio.Writer) (err error) {
	var styleStack []int
	for _, token := range cmdToks {
		popStyleStack := true
		needSpace := true
		switch token {
		case "Ar":
			// Argument
			// Emphasize the next token (wrap in "_").
			if _, err := bufferedWriter.WriteString("_"); err != nil {
				return err
			}
			styleStack = append(styleStack, MAN2MD_AR)
			popStyleStack = false
			needSpace = false

		case "Fl":
			// Flag
			// Prefix next token with '-'.
			if _, err := bufferedWriter.WriteString("-"); err != nil {
				return err
			}
			needSpace = false

		case "Op":
			// Optional part of command line. Enclose rest of the line in "[]".
			if _, err := bufferedWriter.WriteString("["); err != nil {
				return err
			}
			styleStack = append(styleStack, MAN2MD_OP)
			popStyleStack = false
			needSpace = false

		default:
			// Pass the token through to the output stream.
			if _, err := bufferedWriter.WriteString(token); err != nil {
				return err
			}
			needSpace = true
		}

		if popStyleStack && len(styleStack) > 0 {
			style := styleStack[len(styleStack)-1]
			if style != MAN2MD_OP {
				styleStack = styleStack[:len(styleStack)-1]
				switch style {
				case MAN2MD_AR:
					if _, err := bufferedWriter.WriteString("_"); err != nil {
						return err
					}
				}
			}
		}

		if needSpace {
			if _, err := bufferedWriter.WriteString(" "); err != nil {
				return err
			}
		}
	}

	if len(styleStack) > 0 {
		for styleIndex := len(styleStack) - 1; styleIndex >= 0; styleIndex-- {
			style := styleStack[styleIndex]
			switch style {
			case MAN2MD_OP:
				if _, err := bufferedWriter.WriteString("]"); err != nil {
					return err
				}
			}
			styleStack = styleStack[:len(styleStack)-1]
		}
	}

	if _, err := bufferedWriter.WriteString(" "); err != nil {
		return err
	}

	return err
}
