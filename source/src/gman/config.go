// Copyright 2014 The Gman Authors. All rights reserved.
// Use of this source code is governed by an MIT-style
// license that can be found in LICENSE file.

package main

import (
    "encoding/json"   // for reading config file
    "io/ioutil"       // for reading files and logging
    "os/user"         // for finding user home directory
    "log"             // for debug logging
)

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
        // If we didn't find a config file in the user's home directory
        // then we'll look for the gmanrc file colocated with the executable.
        file, err = ioutil.ReadFile("gmanrc")
        if err != nil {
            log.Println("Error reading .gmanrc file: ", err)
        }
        return config
    }

    json.Unmarshal(file, &config)
    return config
}


