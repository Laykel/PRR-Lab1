// Lab 1 - clock synchronization
// File: utils/utils.go
// Authors: Jael Dubey, Luc Wachter
// Go version: 1.13.1 (linux/amd64)

// The utils package contains the needed values and functions to write
// traces to stdout
package utils

import (
	"log"
)

const (
	MasterFilename = "master.go"
	SlaveFilename  = "slave.go"
)

func Trace(filename string, message string) {
	log.Printf("from %s: %s\n", filename, message)
}
