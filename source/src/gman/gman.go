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
    "encoding/json"   // for reading config file
    "flag"            // for command line parameters
    // "fmt"             // for printing runtime errors
    "io/ioutil"       // for reading files and logging
    // "net/http"        // for creating an http web service
    // "net/url"         // for url escaping
    // "open"            // for launching the web browser
    // "os/exec"         // for executing conversions
    "os"              // for local file access
    "os/user"            // for finding user home directory
    "log"             // for debug logging
    // "path"            // for local file system access
    // "path/filepath"   // for local file system access
    // "regexp"          // for file parsing
    // "sort"            // for sorting the file list
    // "strings"         // for string manipulation
)

// TODO: Configuration reading is probably complex enough to be in it's
// own source file.

// TODO: use a better type for GManPathMap
type ConfigObject struct {
    Help  bool                   // Show help (gman help page)
    Debug bool
    Man   bool
    Section  string
    HttpPort string
    SectionSearchOrder string
    GManPath    []string
    GManPathMap []string
}

// Define options with some defaults
var config = ConfigObject {
    Help:     false,
    Debug:    false,
    Man:      false,
    Section:  "",
    HttpPort: "8088",
    SectionSearchOrder: "1 n l 8 3 2 3posix 3pm 3perl 5 4 9 6 7",
    GManPath: []string{"", ""},
    GManPathMap: []string{"", ""},
}

// TODO: We probably need to return extra command line arguments
// as an array of strings.
func readConfig() ConfigObject {

    // get user directory
    usr, err := user.Current()
    if err != nil {
        log.Println("Error file user home directory: ", err)
        return config
    }
    log.Println("User homedir: ", usr.HomeDir)

    // TODO: Check for config file in /etc/gmanrc
    // DONE: Check for config file in $HOME/.gmanrc
    // TODO: If config.Man then read /etc/manpath.config
    file, err := ioutil.ReadFile(usr.HomeDir + "/.gmanrc")
    if err != nil {
        log.Println("Error reading .gmanrc file: ", err)
        return config
    }

    json.Unmarshal(file, &config)
    return config
}

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

