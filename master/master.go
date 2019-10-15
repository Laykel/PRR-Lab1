// TODO: headers
package main

import (
	"bufio"
	"bytes"
	"github.com/Laykel/PRR-Lab1/protocol"
	"github.com/Laykel/PRR-Lab1/utils"
	"log"
	"net"
	"time"
)

// TODO: Put that in another package???
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
	utils.Trace(utils.MasterFilename, "SYNC and FOLLOWUP sent (multicast)")
	protocol.SendSync(id)
	protocol.SendFollowUp(id, time.Now())
}

func main() {
	// Periodically broadcast
	go doEvery(protocol.SyncPeriod, syncAndFollowUp)

	// Listen on the UDP port specified in protocol
	conn, err := net.ListenPacket("udp", protocol.UnicastMasterPort)
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
			messageCode := utils.ParseUdpMessage(s.Text(), 0, protocol.Separator)

			// If the message received is indeed a DELAY_REQUEST
			if uint8(messageCode) == protocol.DelayRequest {

				idDelayRequest := utils.ParseUdpMessage(s.Text(), 1, protocol.Separator)

				utils.Trace(utils.MasterFilename, "DelayRequest received with message : " + s.Text())

				message := s.Text() + " from " + clientAddress.String() + "\n"
				if _, err := conn.WriteTo([]byte(message), clientAddress); err != nil {
					log.Fatal(err)
				}

				utils.Trace(utils.MasterFilename, "DelayResponse sent")
				// Send DELAY_RESPONSE
				protocol.SendDelayResponse(clientAddress, tM, uint(idDelayRequest))
			}
		}
	}
}
