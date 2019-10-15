package utils

import (
	"log"
	"strconv"
	"strings"
)

const (
	MasterFilename = "master.go"
	SlaveFilename = "slave.go"
)

// Return message in a given position from a string separated by a character
func ParseUdpMessage(s string, position uint, separator rune) uint64 {
	var result uint64

	tokens := strings.FieldsFunc(s, func(r rune) bool {
		return r == separator
	})

	// Get the message code
	result, err := strconv.ParseUint(tokens[position], 10, 64)
	if err != nil {
		log.Fatal(err)
	}

	return result
}

func Trace(filename string, message string) {
	log.Printf("From %s : %s\n", filename, message)
}
