package protocol

import (
	"fmt"
	"io"
	"log"
	"net"
	"os"
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
	conn, err := net.Dial("udp", ip.String()+UnicastListenAddress)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	// Write message
	go func() {
		mustCopy(os.Stdout, conn)
	}()
	mustCopy(conn, os.Stdin)
}

func mustCopy(dst io.Writer, src io.Reader) {
	if _, err := io.Copy(dst, src); err != nil {
		log.Fatal(err)
	}
}
