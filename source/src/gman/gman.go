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
    // "fmt"             // for printing runtime errors
    "io/ioutil"       // for reading files and logging
    // "net/http"        // for creating an http web service
    // "net/url"         // for url escaping
    // "open"            // for launching the web browser
    // "os/exec"         // for executing conversions
    "os"              // for local file access
    "log"             // for debug logging
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

    // log.Printf("Config: %v\n", config)
    log.Printf("Section: %s\n", config.Section)
}

