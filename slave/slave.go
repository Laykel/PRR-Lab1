package main

import (
    "fmt"
    "github.com/Laykel/PRR-Lab1/protocol"
	"github.com/Laykel/PRR-Lab1/utils"
	"golang.org/x/net/ipv4"
	"log"
	"math/rand"
	"net"
	"runtime"
	"strconv"
	"time"
)

func main() {
	connMulticast := protocol.ListenUDPConnection(protocol.MulticastAddress)
	defer connMulticast.Close()
	connUnicast := protocol.ListenUDPConnection(protocol.UnicastSlavePort)
	defer connUnicast.Close()

	// Get server's ipv4
	p := ipv4.NewPacketConn(connMulticast)
	addr, err := net.ResolveUDPAddr("udp", protocol.MulticastAddress)
	if err != nil {
		log.Fatal(err)
	}

	var interf *net.Interface
	if runtime.GOOS == "darwin" {
		interf, _ = net.InterfaceByName("en0")
	}

	// Listen for message on multicast group
	if err = p.JoinGroup(interf, addr); err != nil {
		log.Fatal(err)
	}

	buf := make([]byte, protocol.MaxBufferSize)
	var tI, offsetI, tES, shiftI int64
	var delayRequestId uint8

	for {
		// SYNC
		s, addr := protocol.ConnToScanner(connMulticast, buf)
		s.Scan()
		syncCode, syncId := protocol.SyncDecode(s.Text())

		if syncCode == protocol.Sync {
			utils.Trace(utils.SlaveFilename, "SYNC received with id: "+strconv.Itoa(int(syncId)))

			// Record slave time
			tI = time.Now().UnixNano() / int64(time.Microsecond)
		} else {
			utils.Trace(utils.SlaveFilename, "First message received was not a SYNC!")
			continue
		}

		// FOLLOW_UP
		s, addr = protocol.ConnToScanner(connMulticast, buf)
		s.Scan()
		followUpCode, followUpId, tMaster := protocol.FollowUpDecode(s.Text())

		if followUpCode == protocol.FollowUp {
			utils.Trace(utils.SlaveFilename, "FOLLOWUP received with id: "+strconv.Itoa(int(followUpId)))

			if followUpId == syncId {
				// Calculate offset
				offsetI = tMaster - tI
			} else {
				utils.Trace(utils.SlaveFilename, "FOLLOWUP id is not equal to previous SYNC id!")
				continue
			}
		} else {
			utils.Trace(utils.SlaveFilename, "No FOLLOWUP message was received!")
			continue
		}

		// DELAY_REQUEST
		rand.Seed(time.Now().UnixNano())
		// TODO: uncomment that
		// Wait between 4 and 60 times the sync period
		//timeToWait := (rand.Intn(56) + 4) * protocol.SyncPeriod
		timeToWait := 2
		time.Sleep(time.Duration(timeToWait) * time.Second)

		// Record time
		tES = time.Now().UnixNano() / int64(time.Microsecond)

		protocol.SendDelayRequest(addr, delayRequestId)
		utils.Trace(utils.SlaveFilename, "DelayRequest sent")

		// DELAY_RESPONSE
		s, addr = protocol.ConnToScanner(connUnicast, buf)
		s.Scan()
        delayResponseCode, delayResponseId, tM := protocol.DelayResponseDecode(s.Text())

        if delayResponseCode == protocol.DelayResponse {
            utils.Trace(utils.SlaveFilename, "DelayResponse received with id: "+strconv.Itoa(int(delayResponseId)))

            if delayResponseId == delayRequestId {
                // Calculate delay
                delayI := (tM - tES) / 2
                // Calculate shift
                shiftI = offsetI + delayI

                utils.Trace(utils.SlaveFilename, "Shift_i determined: "+strconv.Itoa(int(shiftI))+" [μs]")
                delayRequestId++
            } else {
                utils.Trace(utils.SlaveFilename, "DELAYRESPONSE id is not equal to DELAYREQUEST id!")
            }
        } else {
            utils.Trace(utils.SlaveFilename, "No DELAYRESPONSE was received!")
        }

        fmt.Println("------------------------------------")
	}
}
