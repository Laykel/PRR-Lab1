package utils

import (
	"bufio"
	"github.com/Laykel/PRR-Lab1/protocol"
	"log"
	"strconv"
	"strings"
)

// Parse message separate by a character
func ParseUdpMessage(s bufio.Scanner) uint64 {
	var result uint64

	tokens := strings.FieldsFunc(s.Text(), func(r rune) bool {
		return r == protocol.Separator
	})

	// Get the message code
	result, err := strconv.ParseUint(tokens[1], 10, 64)
	if err != nil {
		log.Fatal(err)
	}

	return result
}