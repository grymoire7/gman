// Copyright 2014 The Gman Authors. All rights reserved.
// Use of this source code is governed by an MIT-style
// license that can be found in LICENSE file.

package main

import (
    "docopt"        // for command line options parsing
    "encoding/json" // for reading config file
    "io/ioutil"     // for reading files and logging
    "log"           // for debug logging
    "os/user"       // for finding user home directory
)

var GMAN_VERSION = "0.1rc"

func loadJSONConfig(filename string) (map[string]interface{}, error) {
    var result map[string]interface{}
    file, err := ioutil.ReadFile(filename)
    if err != nil {
        return result, err
    }
    err = json.Unmarshal(file, &result)
    return result, err
}

// merge combines two maps.
// truthiness takes priority over falsiness
// mapA takes priority over mapB
func merge(mapA, mapB map[string]interface{}) map[string]interface{} {
    result := make(map[string]interface{})
    for k, v := range mapA {
        result[k] = v
    }
    for k, v := range mapB {
        if _, ok := result[k]; !ok || result[k] == nil || result[k] == false {
            result[k] = v
        }
    }
    return result
}

func readConfig() map[string]interface{} {
    usage := `GMan

Usage:
  gman [-d | --debug] [--color] [-s <docsection>]
       [-P pager | --pager=pager] <page>
  gman (-b | --browse) [(-p <port> | --port <port>)]
  gman (-h | --help | -V | --version )

Options:
  -h --help                Show this help.
  -d --debug               Print debug information.
  --color                  Use color text in terminal.
  -s <docsection>          Print document section.
  -p <port> --port <port>  Specifiy port for web server.
  -V --version             Show version.`

    // get user directory
    usr, err := user.Current()
    if err != nil {
        log.Println("Error file user home directory: ", err)
        return nil
    }
    // log.Println("User homedir: ", usr.HomeDir)

    binConfig, _ := loadJSONConfig("gmanrc")
    etcConfig, _ := loadJSONConfig("/etc/gman.conf")
    homeConfig, _ := loadJSONConfig(usr.HomeDir + "/.gmanrc")
    arguments, _ := docopt.Parse(usage, nil, true, GMAN_VERSION, false)

    result := merge(arguments, merge(homeConfig, merge(etcConfig, binConfig)))

    // config.Section = result["this"].(string)
    // json.Unmarshal(file, &config)
    // return config
    return result
}
