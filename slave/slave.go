package main

import (
	"bufio"
	"bytes"
	"fmt"
	"github.com/Laykel/PRR-Lab1/protocol"
	"golang.org/x/net/ipv4"
	"log"
	"math/rand"
	"net"
	"runtime"
	"strconv"
	"strings"
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

			// Separate message with the separator
			tokens := strings.FieldsFunc(s.Text(), func(r rune) bool {
				return r == protocol.Separator
			})

			// Get the message code
			messageType, err := strconv.ParseUint(tokens[0], 10, 8)
			if err != nil {
				log.Fatal(err)
			}

			switch uint8(messageType) {
				case protocol.Sync:
					tI = time.Now().Unix()

				case protocol.FollowUp:
					// Separate message with the separator
					tokens := strings.FieldsFunc(s.Text(), func(r rune) bool {
						return r == protocol.Separator
					})

					// Get the message code
					tMaster, err := strconv.ParseUint(tokens[2], 10, 32)
					if err != nil {
						log.Fatal(err)
					}

					gapI = int64(tMaster) - tI

					rand.Seed(time.Now().UnixNano())
					timeToWait := rand.Intn(56) + 4

					fmt.Printf("%d secondes\n", timeToWait)

					time.Sleep(time.Duration(timeToWait) * time.Second)

					tES = time.Now().Unix()
					protocol.SendDelayRequest(addr, idDelayRequest)
					idDelayRequest++

				case protocol.DelayResponse:
					// Separate message with the separator
					tokens := strings.FieldsFunc(s.Text(), func(r rune) bool {
						return r == protocol.Separator
					})

					// Get the message code
					tM, err := strconv.ParseUint(tokens[1], 10, 32)
					if err != nil {
						log.Fatal(err)
					}

					idDelayResponse, err := strconv.ParseUint(tokens[2],10,32)
					if err != nil {
						log.Fatal(err)
					}

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
