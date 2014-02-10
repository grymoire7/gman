// Copyright 2014 The Gman Authors. All rights reserved.
// Use of this source code is governed by an MIT-style
// license that can be found in LICENSE file.

/*
 * gman is a man-like help facility with the following features:
 *     - ability to show just certain pieces/sections of a gman page.
 *     - TLDR sections for the user on the go.
 *     - ability to leverage existing man pages, if needed
 *     - ability to start http server for interactive browsing
 *     - other stuff.
 *
 * See https://github.com/grymoire7/gman for details.
 *
 * Developed with Go 1.1.2
 */

package main

import (
    "flag"            // for command line parameters
    "io/ioutil"       // for reading files and logging
    "os"              // for local file access
    "log"             // for debug logging
    "blackfriday"     // markdown parser
    "fmt"             // for printing runtime errors
    // "net/http"        // for creating an http web service
    // "net/url"         // for url escaping
    // "open"            // for launching the web browser
    // "os/exec"         // for executing conversions
    // "path"            // for local file system access
    // "path/filepath"   // for local file system access
    // "regexp"          // for file parsing
    // "sort"            // for sorting the file list
    // "strings"         // for string manipulation
)

func init() {
}

func main() {
    // Get configuration options from .gmanrc first.
    config := readConfig()

    // TODO: Command line options handling may soon deserve it's own source file.
    // Override config values with options from the command line.
    // We do flags this way so we can supplement os.Args with
    // arguments from the .gmanrc file (config.Clopts).
    f := flag.NewFlagSet("gman", flag.ExitOnError)
    f.BoolVar(  &config.Debug,   "d",   config.Debug,   "turn on debug info")
    f.BoolVar(  &config.Man,     "man", config.Man,     "use man pages as input")
    f.StringVar(&config.Section, "s",   config.Section, "which help sections to print")

    // TODO: prepend config.Clopts  to os.Args[1:] for parsing
    f.Parse(os.Args[1:])
    // f.Parse([]string{"-d", "-s", "bob"})

    // Discard logging messages if not in debug mode.
    if config.Debug {
        log.Println("Debug on");
    } else {
        log.Println("Debug off");
        log.SetOutput(ioutil.Discard)
    }

    var input []byte
    var err error
    var args = f.Args()
    // TODO: look for gman file in gmanpath
    // TODO: support .gz

    switch len(args) {
    case 0:
        if input, err = ioutil.ReadAll(os.Stdin); err != nil {
            fmt.Fprintln(os.Stderr, "Error reading from Stdin:", err)
            os.Exit(-1)
        }
    case 1,2:
        if input, err = ioutil.ReadFile(args[0]); err != nil {
            fmt.Fprintln(os.Stderr, "Error reading from", args[0], ":", err)
            os.Exit(-1)
        }
    default:
        fmt.Println("Too many arguments (", len(args), ").")
        os.Exit(-1);
    }

    if len(args) > 0 {
        log.Println("In from:", args[0])
    }

	extensions := 0
    extensions |= blackfriday.EXTENSION_NO_INTRA_EMPHASIS
    extensions |= blackfriday.EXTENSION_TABLES
    extensions |= blackfriday.EXTENSION_FENCED_CODE
    extensions |= blackfriday.EXTENSION_AUTOLINK
    extensions |= blackfriday.EXTENSION_STRIKETHROUGH
    extensions |= blackfriday.EXTENSION_SPACE_HEADERS

    renderer := blackfriday.TerminalRenderer(0)
    // output := blackfriday.MarkdownCommon(input)
    output := blackfriday.Markdown(input, renderer, extensions)

   	// output the result
    var out *os.File
    if len(args) == 2 {
        if out, err = os.Create(args[1]); err != nil {
            fmt.Fprintf(os.Stderr, "Error creating %s: %v", args[1], err)
            os.Exit(-1)
        }
        defer out.Close()
    } else {
        out = os.Stdout
    }

    if _, err = out.Write(output); err != nil {
        fmt.Fprintln(os.Stderr, "Error writing output:", err)
        os.Exit(-1)
    }

    log.Printf("Section: %s\n", config.Section)
}

