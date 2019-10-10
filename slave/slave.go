package main

import (
    "bufio"
    "bytes"
    "fmt"
    "github.com/Laykel/PRR-Lab1/protocol"
    "golang.org/x/net/ipv4"
    "log"
    "net"
    "runtime"
)

func main() {
    // TODO: Fill in the blanks xD
    // Listen for messages on multicast group

    // When SYNC message arrives, record time
    // tI := time.Now()

    // When FOLLOW_UP message arrives, parse master time

    // And calculate offset
    // offset := tMaster.Sub(tI)

    // Record time and send DELAY_REQUEST to extracted ip address
    // tEs := time.Now()
    // protocol.SendDelayRequest(serverAddress)

    // Wait for DELAY_RESPONSE and parse master reception time

    // Then calculate delay
    // delay := tM.Sub(tEs)






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
	for {
		n, addr, err := conn.ReadFrom(buf)
		if err != nil {
			log.Fatal(err)
		}

		s := bufio.NewScanner(bytes.NewReader(buf[0:n]))
		for s.Scan() {
			fmt.Printf("%s from %v\n", s.Text(), addr)
		}
	}
}
