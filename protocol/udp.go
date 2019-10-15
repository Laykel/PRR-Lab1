package protocol

import (
	"bufio"
	"bytes"
	"fmt"
	"github.com/Laykel/PRR-Lab1/utils"
	"io"
	"log"
	"net"
	"strconv"
	"strings"
	"time"
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
func sendUnicast(ip net.Addr, port string, message string) {
	// Get descriptor
	tokens := strings.FieldsFunc(ip.String(), func(r rune) bool {
		return r == ':'
	})

	conn, err := net.Dial("udp", tokens[0]+port)

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

// Receive message through UDP
func ReceiveUnicast(message string, messageType uint8) (int64) {
	var result int64

	_messageType := utils.ParseUdpMessage(message, 0, Separator)

	if uint8(_messageType) == messageType {
		result = time.Now().UnixNano() / int64(time.Microsecond)
	} else {
		log.Fatal("Message type isn't " + strconv.Itoa(int(messageType)))
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

// Listen an UDP connection specified by an address
func ListenUDPConnection(address string) net.PacketConn {
	result, err := net.ListenPacket("udp", address)

	if err != nil {
		log.Fatal(err)
	}


	return result
}
