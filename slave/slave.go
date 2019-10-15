package main

import (
	"bufio"
	"bytes"
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
	// Listen for multicast
	connMulticast, err := net.ListenPacket("udp", protocol.MulticastAddress)
	if err != nil {
		log.Fatal(err)
	}
	defer connMulticast.Close()

	// Listen for unicast
	connUnicast, err := net.ListenPacket("udp", ":2206")
	if err != nil {
		log.Fatal(err)
	}
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
		n, addr, err := connMulticast.ReadFrom(buf)
		if err != nil {
			log.Fatal(err)
		}

		s := bufio.NewScanner(bytes.NewReader(buf[0:n]))

		// Sync loop
		for s.Scan() {
			utils.Trace(utils.SlaveFilename, "SYNC received with message : " + s.Text())

			messageType := utils.ParseUdpMessage(s.Text(), 0, protocol.Separator)

			if uint8(messageType) == protocol.Sync {
				tI = time.Now().UnixNano() / int64(time.Microsecond)
			}
		}
		n, addr, err = connMulticast.ReadFrom(buf)
		if err != nil {
			log.Fatal(err)
		}

		s = bufio.NewScanner(bytes.NewReader(buf[0:n]))

		// FollowUp loop
		for s.Scan() {
			utils.Trace(utils.SlaveFilename, "FOLLOWUP received with message : " + s.Text())

			messageType := utils.ParseUdpMessage(s.Text(), 0, protocol.Separator)

			if uint8(messageType) == protocol.FollowUp {
				tMaster := utils.ParseUdpMessage(s.Text(), 2, protocol.Separator)

				offsetI = int64(tMaster) - tI

				rand.Seed(time.Now().UnixNano())
				//timeToWait := rand.Intn(56) + 4
				timeToWait := 2

				time.Sleep(time.Duration(timeToWait) * time.Second)

				tES = time.Now().UnixNano() / int64(time.Microsecond)

				utils.Trace(utils.SlaveFilename, "DelayRequest sent")
				protocol.SendDelayRequest(addr, idDelayRequest)
			}
		}

		n, addr, err = connUnicast.ReadFrom(buf)
		if err != nil {
			log.Fatal(err)
		}

		s = bufio.NewScanner(bytes.NewReader(buf[0:n]))

		// DelayResponse loop
		for s.Scan() {

			utils.Trace(utils.SlaveFilename, "DelayResponse received with message : " + s.Text())

			messageType := utils.ParseUdpMessage(s.Text(), 0, protocol.Separator)

			if uint8(messageType) == protocol.DelayResponse {

				tM := utils.ParseUdpMessage(s.Text(), 1, protocol.Separator)
				idDelayResponse := utils.ParseUdpMessage(s.Text(), 2, protocol.Separator)

				if uint64(idDelayRequest) != idDelayResponse {
					log.Fatal("id delayRequest and delayResponse not the same")
				}

				delayI := (int64(tM) - tES) / 2

				shiftI = offsetI + delayI

				utils.Trace(utils.SlaveFilename, "Shift_i determined : " + strconv.Itoa(int(shiftI)) + " [μs]\n------------------------------------")
				idDelayRequest++
			}
		}
	}
}
