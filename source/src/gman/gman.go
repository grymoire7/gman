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
    "blackfriday" // markdown parser
    "fmt"         // for printing runtime errors
    "io"          // for piping through pager
    "io/ioutil"   // for reading files and logging
    "log"         // for debug logging
    "os"          // for local file access
    "os/exec"     // for piping through pager
    "strings"     // for string manipulation
)

func init() {
}

func main() {
    // Get configuration options from .gmanrc first.
    opts := readConfig()

    // Discard logging messages if not in debug mode.
    if debug, ok := opts["--debug"]; ok && debug.(bool) {
        log.Println("Debug on")
    } else {
        log.Println("Debug off")
        log.SetOutput(ioutil.Discard)
    }

    for k, v := range opts {
        if !strings.HasPrefix(k, "_") {
            log.Printf("map: %s => %v\n", k, v)
        }
    }
    // TODO: Should be more robust. Should search a path collection and allow
    // missing os and lang dirs.
    page := opts["<page>"].(string)
    gmanpath := opts["gmanpath"].(string)
    pagepath := gmanpath + "/linux/en/gman1/" + page + ".1.md"

    var input []byte
    var err error

    // TODO: support .gz
    if input, err = ioutil.ReadFile(pagepath); err != nil {
        fmt.Fprintln(os.Stderr, "Error reading from", pagepath, ":", err)
        os.Exit(-1)
    }

    extensions := 0
    extensions |= blackfriday.EXTENSION_NO_INTRA_EMPHASIS
    extensions |= blackfriday.EXTENSION_TABLES
    extensions |= blackfriday.EXTENSION_FENCED_CODE
    extensions |= blackfriday.EXTENSION_AUTOLINK
    extensions |= blackfriday.EXTENSION_STRIKETHROUGH
    extensions |= blackfriday.EXTENSION_SPACE_HEADERS

    renderer := blackfriday.TerminalRenderer(0)
    output := blackfriday.Markdown(input, renderer, extensions)

    // declare the pager
    cmd := exec.Command("less", "-R")
    // create a blocking pipe
    r, stdin := io.Pipe()
    // Set the i/o's
    cmd.Stdin = r
    cmd.Stdout = os.Stdout
    cmd.Stderr = os.Stderr

    // Create a blocking chan, Run the pager and unblock once it is finished
    c := make(chan struct{})
    go func() {
        defer close(c)
        cmd.Run()
    }()

    // Pass output to the pipe
    fmt.Fprintf(stdin, string(output))

    // Close stdin (allows pager to exit)
    stdin.Close()

    // Wait for the pager to be finished
    <-c
}
