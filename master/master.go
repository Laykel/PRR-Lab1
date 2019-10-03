package main

import (
    "github.com/Laykel/PRR-Lab1/protocol"
    "log"
    "net"
)

func main() {
	conn, err := net.Dial("udp", protocol.MulticastAddress)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	for {
        conn.Write([]byte("testst"))
    }
}
