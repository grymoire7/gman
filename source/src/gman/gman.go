// Copyright 2014 The Gman Authors. All rights reserved.
// Use of this source code is governed by an BSD-style
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
    "bufio"       // for section extraction
    "bytes"       // for section extraction
    "errors"      // for reporting errors
    "fmt"         // for printing runtime errors
    "io"          // for piping through pager
    "io/ioutil"   // for reading files and logging
    "log"         // for debug logging
    "os"          // for local file access
    "os/exec"     // for piping through pager
    "strings"     // for string manipulation
    "compress/gzip"
)

func init() {
}

func main() {
    // Get configuration options from rc files and command line.
    opts := readConfig()

    // Discard logging messages if not in debug mode.
    if debug, ok := opts["--debug"]; ok && debug.(bool) {
        log.Println("Debug on")
    } else {
        log.SetOutput(ioutil.Discard)
    }

    // log configuration options
    for k, v := range opts {
        if !strings.HasPrefix(k, "_") {
            log.Printf("opt: %s => %v\n", k, v)
        }
    }

    // TODO: Should be more robust. Should search a path collection and allow
    // missing os and lang dirs.
    page := opts["<page>"].(string)
    gmanpath := opts["gmanpath"].(string)
    pagepath := gmanpath + "/linux/en/gman1/" + page + ".1.md"
    gzpagepath := gmanpath + "/linux/en/gman1/" + page + ".1.gz"

    var input []byte
    var err error

    // try to open the compressed version first
    f, err := os.Open(gzpagepath)
    if err != nil {
        // else try to open the non-compressed version
        if input, err = ioutil.ReadFile(pagepath); err != nil {
            fmt.Fprintln(os.Stderr, "gman: help page not found")
            log.Println("Error reading from", pagepath, ":", err)
            os.Exit(-1)
        }
    } else {
        defer f.Close()
        gz, err := gzip.NewReader(f)
        if err != nil {
            fmt.Fprintln(os.Stderr, "gman: help page not found")
            log.Println("Error reading from", gzpagepath, ":", err)
            os.Exit(-1)
        } else {
            // is there a better way to do this?
            buf := bytes.NewBuffer(input)
            _, err := io.Copy(buf, gz)
            if err != nil {
                fmt.Fprintln(os.Stderr, "gman: help page not found")
                log.Println("Error reading from", gzpagepath, ":", err)
                os.Exit(-1)
            }
            input = buf.Bytes()
        }
    }


    // handle section extraction option
    if opts["-s"] != nil {
        log.Println("Exracting", opts["-s"], "...")
        c, err := extractDocSection(input, opts["-s"].(string))
        if err == nil {
            input = c
        } else {
            fmt.Fprintln(os.Stderr, err)
            os.Exit(-1)
        }
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

    // declare the pager in raw mode
    cmd := exec.Command("less", "-R")
    // create a blocking pipe
    r, stdin := io.Pipe()
    // set the ins and outs
    cmd.Stdin = r
    cmd.Stdout = os.Stdout
    cmd.Stderr = os.Stderr

    // Create a blocking chan, run the pager and unblock once it is finished
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

func extractDocSection(input []byte, sectionPattern string) ([]byte, error) {
    var lines []string
    var inSection = false
    var foundSection = false
    var sectionLevel = 0
    scanner := bufio.NewScanner(bytes.NewReader(input))
    for scanner.Scan() {
        line := scanner.Text()
        if strings.HasPrefix(line, "#") {
            var level = 0
            for _, c := range line {
                if c != '#' {
                    break
                }
                level++
            }
            if inSection {
                if level <= sectionLevel {
                    inSection = false
                }
            } else {
                if strings.Contains(line, sectionPattern) {
                    inSection = true
                    foundSection = true
                    sectionLevel = level
                }
            }
        }
        if inSection {
            lines = append(lines, scanner.Text())
        }
    }
    s := strings.Join(lines, "\n")
    if !foundSection {
        return nil, errors.New("gman: document section not found")
    }
    return []byte(s), scanner.Err()
}
