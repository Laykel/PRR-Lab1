package main

import (
    "bufio"
    "bytes"
    "github.com/Laykel/PRR-Lab1/protocol"
    "log"
    "net"
    "time"
)

// Call given function every given number of seconds
func doEvery(seconds uint, f func(uint)) {
	ticker := time.NewTicker(time.Duration(seconds) * time.Second)
	defer ticker.Stop()

	var counter uint
	for _ = range ticker.C {
		f(counter)
		counter++
	}
}

func syncAndFollowUp(id uint) {
    protocol.SendSync(id)
    protocol.SendFollowUp(id)
}

func main() {
	// Periodically broadcast
	go doEvery(protocol.SyncPeriod, syncAndFollowUp)

	// Listen on the UDP port specified in protocol
	conn, err := net.ListenPacket("udp", protocol.UnicastListenAddress)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	buf := make([]byte, protocol.MaxBufferSize)
	for {
		// Receive DELAY_REQUEST
		n, clientAddress, err := conn.ReadFrom(buf)
		if err != nil {
			log.Fatal(err)
		}

        // Syscall for time
        tM := time.Now()

		// TODO: Test whether the DELAY_REQUEST is valid
		s := bufio.NewScanner(bytes.NewReader(buf[0:n]))
		for s.Scan() {
			s := s.Text() + " from " + clientAddress.String() + "\n"
			if _, err := conn.WriteTo([]byte(s), clientAddress); err != nil {
				log.Fatal(err)
			}
		}

        // TODO: If not, continue (go to next loop iteration)

		// Send DELAY_RESPONSE
		protocol.SendDelayResponse(clientAddress, tM)
	}
}
