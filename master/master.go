package main

import (
	"bufio"
	"bytes"
	"github.com/Laykel/PRR-Lab1/protocol"
	"io"
	"log"
	"net"
	"os"
	"time"
)

// TODO: Move this in separate package (protocol?)
// ----------------------------------------------------------------------------------------
func sendMulticast(message string) {
	conn, err := net.Dial("udp", protocol.MulticastAddress)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	conn.Write([]byte(message))
}

func sendUnicast(message string, ip net.Addr) {
	conn, err := net.Dial("udp", ip.String()+protocol.UnicastPort)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()
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

// ----------------------------------------------------------------------------------------

// Call given function every given number of seconds
func doEvery(seconds uint, f func()) {
	t := time.NewTicker(time.Duration(seconds) * time.Second)

	for _ = range t.C {
		f()
	}
}

// TODO: Move in protocol package?
// Send SYNC and FOLLOW_UP messages to multicast
func syncAndFollowUp() {
	// TODO: Generate ID
	// ...

	// SYNC (message code + ID)
	sendMulticast("CODE + ID")

	// Syscall for time
	tMaster := time.Now()

	// FOLLOW_UP (message code + ID + tMaster)
	sendMulticast("CODE + ID + " + tMaster.String())
}

func main() {
	// Periodically broadcast
	go doEvery(protocol.SyncPeriod, syncAndFollowUp)

	// Listen on the UDP port specified in protocol
	conn, err := net.ListenPacket("udp", protocol.UnicastPort)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	buf := make([]byte, 1024)
	for {
		// Receive DELAY_REQUEST
		n, clientAddress, err := conn.ReadFrom(buf)
		if err != nil {
			log.Fatal(err)
		}

		// TODO: Test whether the DELAY_REQUEST is valid
		s := bufio.NewScanner(bytes.NewReader(buf[0:n]))
		for s.Scan() {
			s := s.Text() + " from " + clientAddress.String() + "\n"
			if _, err := conn.WriteTo([]byte(s), clientAddress); err != nil {
				log.Fatal(err)
			}
		}

        // TODO: If not, continue (go to next loop iteration)

		// Syscall for time
		tM := time.Now()

		// Send DELAY_RESPONSE
		sendUnicast("CODE + ID + "+tM.String(), clientAddress)
	}
}
