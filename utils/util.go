package utils

import (
	"github.com/Laykel/PRR-Lab1/protocol"
	"log"
	"strconv"
	"strings"
)

// Return message in a given position from a string separated by a character
func ParseUdpMessage(s string, position uint) uint64 {
	var result uint64

	tokens := strings.FieldsFunc(s, func(r rune) bool {
		return r == protocol.Separator
	})

	// Get the message code
	result, err := strconv.ParseUint(tokens[position], 10, 64)
	if err != nil {
		log.Fatal(err)
	}

	return result
}