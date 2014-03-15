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
	"regexp"
	"sort"
	"strings"

//	"unicode"
//	"unicode/utf8"
)

type Man2mdParser struct {
	pageName        string
	bof             bool
	reader          *bufio.Reader
	writer          *bufio.Writer
	unprocessedCmds []string // dsn debug
}

func Convert(reader io.Reader, writer io.Writer) (err error) {
	parser := NewParser(reader, writer)
	err = parser.Parse()
	return err
}

func NewParser(reader io.Reader, writer io.Writer) (parser *Man2mdParser) {
	parser = new(Man2mdParser)
	parser.bof = true
	parser.reader = bufio.NewReader(reader)
	parser.writer = bufio.NewWriter(writer)

	return parser
}

func (parser *Man2mdParser) Parse() (err error) {
	eof := false
	for !eof {
		line, err := parser.reader.ReadString('\n')
		if err == io.EOF {
			err = nil
			eof = true
		} else if err != nil {
			return err
		}

		if len(line) > 0 {
			if line[0] == '.' {
				// Dot command
				// Eat the initial '.' and send to the parser method.
				if err = parser.parseMacroLine(line[1:]); err != nil {
					return err
				}
			} else {
				// Check for start of options.
				// TODO: localize this
				if regexp.MustCompile(".*options.*:[ \t\n]*$").MatchString(line) {
					header := "OPTIONS"
					s := fmt.Sprintf("%s\n%s\n", header, strings.Repeat("-", len(header)))
					if _, err = parser.writer.WriteString(s); err != nil {
						return err
					}
				} else {
					// Write the line, replacing trailing "\n" with a space.
					s := fmt.Sprintf("%s ", line[:len(line)-1])
					if _, err = parser.writer.WriteString(s); err != nil {
						return err
					}
				}
			}
		}
	}

	parser.writer.Flush()

	// dsn debug
	// Dump the list of unprocessed commands.
	fmt.Printf("\n%d unprocessed dot commands:\n", len(parser.unprocessedCmds))
	for _, cmd := range parser.unprocessedCmds {
		fmt.Printf("\t%s\n", cmd)
	}
	// end dsn debug

	return err
}

// Process a line beginning with '.' (macro).
func (parser *Man2mdParser) parseMacroLine(line string) (err error) {
	// Split the line into tokens delimited by white space.
	// Get rid of the terminating '\n' while we're at it.
	tokens := regexp.MustCompile("[ \t]").Split(line[:len(line)-1], -1)

	var parsedTokens []string
	if parsedTokens, err = parser.parseTokens(tokens, true); err != nil {
		return err
	}

	// Write the parsed tokens to the output stream.
	for _, token := range parsedTokens {
		if _, err = parser.writer.WriteString(fmt.Sprintf("%s ", token)); err != nil {
			return err
		}
	}

	return err
}

// Parse the macros in the given slice of tokens.
// Returns an array of parsed tokens.
// bol (beginning of line) is true if the "tokens" slice represents the entire line (i.e. the first
// token is a dot command).
func (parser *Man2mdParser) parseTokens(tokens []string, bol bool) (parsedTokens []string, err error) {
	var parsedRemainder []string
	token := tokens[0]
	switch token {
	case "\\\"":
		// Comment. Ignore it.

	case "Ar":
		// Argument
		// Emphasize the next token (wrap in "_").
		if len(tokens) > 1 {
			if parsedRemainder, err = parser.parseTokens(tokens[1:], false); err != nil {
				return nil, err
			}

			parsedTokens = append(parsedTokens, fmt.Sprintf("_%s_", parsedRemainder[0]))

			if len(parsedRemainder) > 1 {
				parsedTokens = append(parsedTokens, parsedRemainder[1:]...)
			}
		}

	case "Bl":
		// Begin list. Ignore

	case "Dd":
		// Document date. Ignore.

	case "Dl":
		// Literal text. Present like a code block (indented by 4 spaces).
		if len(tokens) > 1 {
			s := fmt.Sprintf("\n\n        %s\n", strings.Join(tokens[1:], " "))
			parsedTokens = append(parsedTokens, s)
		}

	case "Dt":
		// Document title.
		// We're only interested in the title and section for now.
		if len(tokens) > 1 {
			title := tokens[1]
			if len(tokens) > 2 {
				title = fmt.Sprintf("%s(%s)", title, tokens[2])
			}

			s := fmt.Sprintf("%s\n%s\n\n", title, strings.Repeat("=", len(title)))
			parsedTokens = append(parsedTokens, s)
		}

	case "Dv":
		// Defined variable.
		// Emphasize (strong)
		if len(tokens) > 1 {
			if parsedRemainder, err = parser.parseTokens(tokens[1:], false); err != nil {
				return nil, err
			}
			parsedTokens = append(parsedTokens, fmt.Sprintf("__%s__ ", parsedRemainder[0]))

			if len(parsedRemainder) > 1 {
				parsedTokens = append(parsedTokens, parsedRemainder[1:]...)
			}
		}

	case "El":
		// End list.
		// Ignore for now.

	case "Er":
		// Error code. Pass it through with "strong" emphasis.
		if len(tokens) > 1 {
			if parsedRemainder, err = parser.parseTokens(tokens[1:], false); err != nil {
				return nil, err
			}

			parsedTokens = append(parsedTokens, fmt.Sprintf("__%s__", parsedRemainder[0]))

			if len(parsedRemainder) > 1 {
				parsedTokens = append(parsedTokens, parsedRemainder[1:]...)
			}
		}

	case "Fx":
		// FreeBSD version
		// For now, just substitute "FreeBSD".
		// TODO: handle default version.
		parsedTokens = append(parsedTokens, "FreeBSD ")

		if len(tokens) > 1 {
			if parsedRemainder, err = parser.parseTokens(tokens[1:], false); err != nil {
				return nil, err
			}
			parsedTokens = append(parsedTokens, parsedRemainder...)
		}

	case "Fl":
		// Flag
		// Prefix subsequent token(s) with '-'.
		// Affects multiple tokens if they are comma-separated.
		if len(tokens) > 1 {
			if parsedRemainder, err = parser.parseTokens(tokens[1:], false); err != nil {
				return nil, err
			}

			flagNext := true
			flagged := false
			for _, token := range parsedRemainder {
				if flagNext && !flagged {
					parsedTokens = append(parsedTokens, fmt.Sprintf("-%s", token))
					flagged = true
				} else if token == "," {
					parsedTokens = append(parsedTokens, ", ")
					flagNext = true
					flagged = false
				} else {
					parsedTokens = append(parsedTokens, token)
					flagNext = false
					flagged = false
				}
			}
		}

	case "It":
		// List item. Precede with "*".
		if len(tokens) > 1 {
			parsedTokens = append(parsedTokens, "\n\n* ")

			if parsedRemainder, err = parser.parseTokens(tokens[1:], false); err != nil {
				return nil, err
			}
			parsedTokens = append(parsedTokens, parsedRemainder...)
			parsedTokens = append(parsedTokens, "\n")
		}

	case "Nd":
		// Description
		parsedTokens = append(parsedTokens, fmt.Sprintf("-- %s", strings.Join(tokens[1:], " ")))

	case "Nm":
		// Page name macro
		// The first occurrence sets the name.
		nextToken := 1
		if len(parser.pageName) == 0 {
			if len(tokens) > 1 {
				parser.pageName = tokens[1]
				nextToken = 2
			}
		}

		parsedTokens = append(parsedTokens, parser.pageName)

		if len(tokens) > nextToken {
			parsedTokens = append(parsedTokens, " ")
			parsedTokens = append(parsedTokens, strings.Join(tokens[1:], " "))
		}

	case "Oc":
		// End multi-line optional section.
		parsedTokens = append(parsedTokens, "] ")

	case "Oo":
		// Begin multi-line optional section (enclose in "[").
		// This will be termintated by ".Oc".
		parsedTokens = append(parsedTokens, "[")

	case "Op":
		// Optional part of command line. Enclose rest of the line in "[]".
		if len(tokens) > 1 {
			if parsedRemainder, err = parser.parseTokens(tokens[1:], false); err != nil {
				return nil, err
			}

			parsedTokens = append(parsedTokens, "[")
			parsedTokens = append(parsedTokens, parsedRemainder...)
			parsedTokens = append(parsedTokens, "]")

		}

	case "Os":
		// Operating system. Ignore.

	case "Pa":
		// Path. Mark next token for emphasis.
		if len(tokens) > 1 {
			parsedTokens = append(parsedTokens, fmt.Sprintf("_%s_", tokens[1]))

			if len(tokens) > 2 {
				if parsedRemainder, err = parser.parseTokens(tokens[2:], false); err != nil {
					return nil, err
				}

				parsedTokens = append(parsedTokens, parsedRemainder...)
			}
		}

	case "Pp", "PP":
		// Paragraph break
		parsedTokens = append(parsedTokens, "\n\n")

	case "Pq":
		// Enclose rest of the line in parens.
		if len(tokens) > 1 {
			parsedTokens = append(parsedTokens, "(")
			if parsedRemainder, err = parser.parseTokens(tokens[1:], false); err != nil {
				return nil, err
			}
			parsedTokens = append(parsedTokens, parsedRemainder...)
			parsedTokens = append(parsedTokens, ") ")
		}

	case "Ql":
		// Single-quoted literal
		if len(tokens) > 1 {
			parsedTokens = append(parsedTokens, fmt.Sprintf("'%s' ", tokens[1]))
			if len(tokens) > 2 {
				if parsedRemainder, err = parser.parseTokens(tokens[2:], false); err != nil {
					return nil, err
				}
				parsedTokens = append(parsedTokens, parsedRemainder...)
			}
			s := fmt.Sprintf("'%s' ", tokens[1])
			parsedTokens = append(parsedTokens, s)
		}

	case "Sh", "SH":
		// Section heading
		if !parser.bof {
			parsedTokens = append(parsedTokens, "\n\n")
		} else {
			parser.bof = false
		}

		if len(tokens) > 1 {
			header := strings.Join(tokens[1:], " ")

			// Rename the "Synopsis" section to "Usage"
			// TODO: Localize
			if header == "SYNOPSIS" {
				header = "USAGE"
			}

			parsedTokens = append(parsedTokens, fmt.Sprintf("%s\n%s\n", header, strings.Repeat("-", len(header))))
		}

	case "Xr":
		// Cross-reference (link to another man page).
		// TODO: determine how we're handling links to other gman pages.
		if len(tokens) > 2 {
			page := tokens[1]
			section := fmt.Sprintf("(%s)", tokens[2])
			sectionExt := fmt.Sprintf(".%s", tokens[2])
			s := fmt.Sprintf("[%s%s](gman://%s%s) ", page, section, page, sectionExt)
			parsedTokens = append(parsedTokens, s)
		}

	default:
		if bol {
			// Beginning of line
			// This must be an unhandled macro.
			// dsn debug
			// Add to our array tracking unprocessed commands.
			index := sort.SearchStrings(parser.unprocessedCmds, token)
			if index >= len(parser.unprocessedCmds) {
				// Add token to the end of the list.
				parser.unprocessedCmds = append(parser.unprocessedCmds, token)
			} else if parser.unprocessedCmds[index] != token {
				// Insert token into the list.
				parser.unprocessedCmds = append(parser.unprocessedCmds[:index],
					append([]string{token}, parser.unprocessedCmds[index:]...)...)
			}
		} else {
			// Not a macro
			// Pass the token through unaltered.
			parsedTokens = append(parsedTokens, token)
			if len(tokens) > 1 {
				if parsedRemainder, err = parser.parseTokens(tokens[1:], false); err != nil {
					return nil, err
				}

				parsedTokens = append(parsedTokens, parsedRemainder...)
			}
		}
	}

	return parsedTokens, err
}
