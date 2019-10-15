package utils

import (
    "log"
    "strconv"
    "strings"
)

const (
	MasterFilename = "master.go"
	SlaveFilename  = "slave.go"
)

func Trace(filename string, message string) {
	log.Printf("From %s : %s\n", filename, message)
}

// Return message in a given position from a string separated by a character
func ParseUdpMessage(s string, position uint, separator string) uint64 {
    var result uint64

    tokens := strings.Split(s, separator)

    // Get the message code
    result, err := strconv.ParseUint(tokens[position], 10, 64)
    if err != nil {
        log.Fatal(err)
    }

    return result
}
