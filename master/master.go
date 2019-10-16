// Lab 1 - clock synchronization
// File: master/master.go
// Authors: Jael Dubey, Luc Wachter
// Go version: 1.13.1 (linux/amd64)

// Main package for master program
// Synchronizes slaves regularly and responds to delay requests
package main

import (
    "github.com/Laykel/PRR-Lab1/protocol"
    "github.com/Laykel/PRR-Lab1/utils"
    "strconv"
    "time"
)

// Call given function every given number of seconds
func doEvery(seconds uint, f func(uint8)) {
	ticker := time.NewTicker(time.Duration(seconds) * time.Second)
	defer ticker.Stop()

	var counter uint8
	for _ = range ticker.C {
		f(counter)
		counter++
	}
}

func syncAndFollowUp(id uint8) {
	protocol.SyncEncode(id)
	protocol.SendFollowUp(id, time.Now())
	utils.Trace(utils.MasterFilename, "SYNC and FOLLOW_UP sent (multicast)")
}

// Main program for master clock
// Periodically syncs and responds to DELAY_REQUEST
func main() {
	// Periodically sync
	go doEvery(protocol.SyncPeriod, syncAndFollowUp)

	// Listen on the UDP port specified in protocol
	conn := protocol.ListenUDPConnection(protocol.UnicastMasterPort)
	defer conn.Close()

	buf := make([]byte, protocol.MaxBufferSize)
	for {
		// Receive DELAY_REQUEST
        s, addr := protocol.ConnToScanner(conn, buf)
        s.Scan()
        delayRequestCode, delayRequestId := protocol.DelayRequestDecode(s.Text())

        // Get time at request reception
        tM := time.Now()

        // If the message received is indeed a DELAY_REQUEST
        if delayRequestCode == protocol.DelayRequest {
            utils.Trace(utils.MasterFilename, "DelayRequest received with id: "+strconv.Itoa(int(delayRequestId)))

            // Send DELAY_RESPONSE
            protocol.SendDelayResponse(addr, delayRequestId, tM)
            utils.Trace(utils.MasterFilename, "DelayResponse sent")
		} else {
            utils.Trace(utils.MasterFilename, "No DELAYREQUEST was received!")
        }
	}
}
