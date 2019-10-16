// Lab 1 - clock synchronization
// File: protocol/udp.go
// Authors: Jael Dubey, Luc Wachter
// Go version: 1.13.1 (linux/amd64)

// The udp part of the protocol package contains functions to send and receive
// messages through UDP
package protocol

import (
    "bufio"
    "bytes"
    "fmt"
    "log"
    "net"
    "strings"
)

// Send message to multicast group
func sendMulticast(message *bytes.Buffer) {
	// Get descriptor
	conn, err := net.Dial("udp", MulticastAddress)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	// Send message
	_, err = fmt.Fprint(conn, message)
	if err != nil {
		log.Fatal(err)
	}
}

// Send message through UDP to specified ip
func sendUnicast(ip net.Addr, port string, message *bytes.Buffer) {
	tokens := strings.Split(ip.String(), ":")

	conn, err := net.Dial("udp", tokens[0]+port)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	// Send message
	_, err = fmt.Fprint(conn, message)
	if err != nil {
		log.Fatal(err)
	}
}

// Listen a UDP connection specified by an address
func ListenUDPConnection(address string) net.PacketConn {
    result, err := net.ListenPacket("udp", address)

    if err != nil {
        log.Fatal(err)
    }

    return result
}

// Take connection and put its message in a Scanner
func ConnToScanner(conn net.PacketConn, buffer []byte) (s *bufio.Scanner, addr net.Addr) {
	n, addr, err := conn.ReadFrom(buffer)
	if err != nil {
		log.Fatal(err)
	}

	s = bufio.NewScanner(bytes.NewReader(buffer[0:n]))

	return
}
