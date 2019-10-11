// TODO: headers
package main

import (
    "bufio"
    "bytes"
    "github.com/Laykel/PRR-Lab1/protocol"
    "log"
    "net"
    "strconv"
    "strings"
    "time"
)

// TODO: Put that in another package???
// Call given function every given number of seconds
func doEvery(seconds uint, f func(uint)) {
	ticker := time.NewTicker(time.Duration(seconds) * time.Second)
	defer ticker.Stop()

	// TODO: make the id go back to 0 when capacity attained
	var counter uint
	for _ = range ticker.C {
		f(counter)
		counter++
	}
}

func syncAndFollowUp(id uint) {
	protocol.SendSync(id)
	protocol.SendFollowUp(id, time.Now())
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

		// Get time at request reception
		tM := time.Now()

		// Read message
		s := bufio.NewScanner(bytes.NewReader(buf[0:n]))
		for s.Scan() {
			// TODO: Put tokens parsing in protocol package?
		// Separate message with the separator
			tokens := strings.FieldsFunc(s.Text(), func(r rune) bool {
				return r == protocol.Separator
			})

			// Get the message code
			messageCode, err := strconv.ParseUint(tokens[0], 10, 8)
			if err != nil {
				log.Fatal(err)
			}

			// If the message received is indeed a DELAY_REQUEST
			if uint8(messageCode) == protocol.DelayRequest {
				s := s.Text() + " from " + clientAddress.String() + "\n"
				if _, err := conn.WriteTo([]byte(s), clientAddress); err != nil {
					log.Fatal(err)
				}

				// Send DELAY_RESPONSE
				protocol.SendDelayResponse(clientAddress, tM)
			}
		}
	}
}
