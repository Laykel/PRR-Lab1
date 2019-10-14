package main

import (
	"bufio"
	"bytes"
	"fmt"
	"github.com/Laykel/PRR-Lab1/protocol"
	"github.com/Laykel/PRR-Lab1/utils"
	"golang.org/x/net/ipv4"
	"log"
	"math/rand"
	"net"
	"runtime"
	"time"
)

func main() {
	// Listen for multicast
	conn, err := net.ListenPacket("udp", protocol.MulticastAddress)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	// Get server's ipv4
	p := ipv4.NewPacketConn(conn)
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
	var tI, gapI, tES, shiftI int64
	var idDelayRequest uint

	for {
		n, addr, err := conn.ReadFrom(buf)
		if err != nil {
			log.Fatal(err)
		}

		s := bufio.NewScanner(bytes.NewReader(buf[0:n]))

		for s.Scan() {
			fmt.Printf("%s from %v\n", s.Text(), addr)

			messageType := utils.ParseUdpMessage(s.Text(), 0)

			switch uint8(messageType) {
				case protocol.Sync:
					tI = time.Now().Unix()

				case protocol.FollowUp:
					tMaster := utils.ParseUdpMessage(s.Text(), 2)

					gapI = int64(tMaster) - tI

					rand.Seed(time.Now().UnixNano())
					timeToWait := rand.Intn(56) + 4

					fmt.Printf("%d secondes\n", timeToWait)

					time.Sleep(time.Duration(timeToWait) * time.Second)

					tES = time.Now().Unix()
					protocol.SendDelayRequest(addr, idDelayRequest)
					idDelayRequest++

				case protocol.DelayResponse:
					tM := utils.ParseUdpMessage(s.Text(), 1)
					idDelayResponse := utils.ParseUdpMessage(s.Text(), 2)

					if uint64(idDelayRequest) != idDelayResponse {
						log.Fatal("id delayRequest and delayResponse not the same")
					}

					noticeI := (int64(tM) - tES) / 2

					shiftI = gapI + noticeI

					fmt.Printf("The shift of this slave is %d\n", shiftI)
			}
		}
	}
}
