package main

import (
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
	connUnicast   := protocol.ListenUDPConnection(protocol.UnicastSlavePort)
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

	// Listen on multicast
	if err = p.JoinGroup(interf, addr); err != nil {
		log.Fatal(err)
	}

	buf := make([]byte, protocol.MaxBufferSize)
	var tI, offsetI, tES, shiftI int64
	var idDelayRequest uint

	for {
		// SYNC
		s, addr := protocol.ConnToScanner(connMulticast, buf)
		s.Scan()
		utils.Trace(utils.SlaveFilename, "SYNC received with message : "+s.Text())
		tI = protocol.ReceiveUnicast(s.Text(), protocol.Sync)


		// FOLLOW_UP
		s, addr = protocol.ConnToScanner(connMulticast, buf)
		s.Scan()
		utils.Trace(utils.SlaveFilename, "FOLLOWUP received with message : "+s.Text())
		tMaster := protocol.ReceiveUnicast(s.Text(), protocol.FollowUp)
		offsetI = tMaster - tI

		// DELAY_REQUEST
		rand.Seed(time.Now().UnixNano())
		//timeToWait := (rand.Intn(56) + 4) * protocol.SyncPeriod
		timeToWait := 2
		time.Sleep(time.Duration(timeToWait) * time.Second)

		tES = time.Now().UnixNano() / int64(time.Microsecond)

		utils.Trace(utils.SlaveFilename, "DelayRequest sent")
		protocol.SendDelayRequest(addr, idDelayRequest)

		// DELAY_RESPONSE
		s, addr = protocol.ConnToScanner(connUnicast, buf)
		s.Scan()
		utils.Trace(utils.SlaveFilename, "DelayResponse received with message : "+s.Text())
		tM := protocol.ReceiveUnicast(s.Text(), protocol.DelayResponse)

		idDelayResponse := utils.ParseUdpMessage(s.Text(), 2, protocol.Separator)
		if uint64(idDelayRequest) != idDelayResponse {
			log.Fatal("id delayRequest and delayResponse not the same")
		}


		delayI := (tM - tES) / 2
		shiftI = offsetI + delayI

		utils.Trace(utils.SlaveFilename, "Shift_i determined : "+strconv.Itoa(int(shiftI))+" [μs]\n------------------------------------")
		idDelayRequest++

	}
}
