// Lab 1 - clock synchronization
// File: slave/slave.go
// Authors: Jael Dubey, Luc Wachter
// Go version: 1.13.1 (linux/amd64)

// Main package for slave program
// Listens for the master clock's messages and sends delay requests to synchronize their clocks
package main

import (
	"fmt"
	"github.com/Laykel/PRR-Lab1/protocol"
	"golang.org/x/net/ipv4"
	"log"
	"math/rand"
	"net"
	"runtime"
	"strconv"
	"time"
)

func main() {
    log.SetFlags(log.LstdFlags | log.Lshortfile)

	var tI, offsetI, tES, shiftI int64
	var delayRequestId uint8

	for {
		//Create unicast and multicast connection
		connMulticast := protocol.ListenUDPConnection(protocol.MulticastAddress)
		connUnicast := protocol.ListenUDPConnection(protocol.UnicastSlavePort)

		p := ipv4.NewPacketConn(connMulticast)
		masterAddr, err := net.ResolveUDPAddr("udp", protocol.MulticastAddress)
		if err != nil {
			log.Fatal(err)
		}

		var interf *net.Interface
		// For MacOS
		if runtime.GOOS == "darwin" {
			interf, _ = net.InterfaceByName("en0")
		}

		// Join multicast group
		if err = p.JoinGroup(interf, masterAddr); err != nil {
			log.Fatal(err)
		}

		buf := make([]byte, protocol.MaxBufferSize)

		// SYNC
		s, addr := protocol.ConnToScanner(connMulticast, buf)
		s.Scan()
		syncCode, syncId := protocol.SyncDecode(s.Text())

		if syncCode == protocol.Sync {
			log.Printf("SYNC received with id: "+strconv.Itoa(int(syncId)))

			// Record slave time
			tI = time.Now().UnixNano() / int64(time.Microsecond)
		} else {
			log.Printf("First message received was not a SYNC!")
			continue
		}

		// FOLLOW_UP
		s, addr = protocol.ConnToScanner(connMulticast, buf)
		s.Scan()
		followUpCode, followUpId, tMaster := protocol.FollowUpDecode(s.Text())

		if followUpCode == protocol.FollowUp {
			if followUpId == syncId {
				// Calculate offset
				offsetI = tMaster - tI

				log.Printf("FOLLOWUP received, offset determined: "+strconv.Itoa(int(offsetI))+" [μs]")
			} else {
				log.Printf("FOLLOWUP id is not equal to previous SYNC id!")
				continue
			}
		} else {
			log.Printf("No FOLLOWUP message was received!")
			continue
		}

		// DELAY_REQUEST
		rand.Seed(time.Now().UnixNano())
		// Wait between 4 and 60 times the sync period
		timeToWait := (rand.Intn(56) + 4) * protocol.SyncPeriod
		log.Printf("Waiting "+strconv.Itoa(timeToWait)+" [s] before DELAYREQUEST")
		time.Sleep(time.Duration(timeToWait) * time.Second)

		// Record time and add offset
		tES = time.Now().UnixNano() / int64(time.Microsecond) + offsetI

		protocol.SendDelayRequest(addr, delayRequestId)
		log.Printf("DelayRequest sent")

		// DELAY_RESPONSE
		// Would be good to add a timeout...
		s, addr = protocol.ConnToScanner(connUnicast, buf)
		s.Scan()
		delayResponseCode, delayResponseId, tM := protocol.DelayResponseDecode(s.Text())

		if delayResponseCode == protocol.DelayResponse {
			log.Printf("DelayResponse received with id: "+strconv.Itoa(int(delayResponseId)))

			if delayResponseId == delayRequestId {
				// Calculate delay
				delayI := (tM - tES) / 2
				// Calculate shift
				shiftI = offsetI + delayI

				log.Printf("Delay determined: "+strconv.Itoa(int(delayI))+" [μs]")
				log.Printf("Shift determined: "+strconv.Itoa(int(shiftI))+" [μs]")
				delayRequestId++
			} else {
				log.Printf("DELAYRESPONSE id is not equal to DELAYREQUEST id!")
			}
		} else {
			log.Printf("No DELAYRESPONSE was received!")
		}

		fmt.Println("------------------------------------")

		connMulticast.Close()
		connUnicast.Close()
	}
}
