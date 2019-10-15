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
