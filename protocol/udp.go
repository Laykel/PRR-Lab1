package protocol

import (
	"fmt"
	"io"
	"log"
	"net"
	"strings"
)

// Send message to multicast group
func sendMulticast(message string) {
    // Get descriptor
	conn, err := net.Dial("udp", MulticastAddress)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	// Write message
    _, err = fmt.Fprintln(conn, message)
	if err != nil {
	    log.Fatal(err)
    }
}

// Send message through UDP to specified ip
func sendUnicast(ip net.Addr, message string) {
    // Get descriptor
	tokens := strings.FieldsFunc(ip.String(), func(r rune) bool {
		return r == ':'
	})

	conn, err := net.Dial("udp", tokens[0]+UnicastPort)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	reader := strings.NewReader(message)

	// Write message
	if _, err := io.Copy(conn, reader); err != nil {
		log.Fatal(err)
	}
}
