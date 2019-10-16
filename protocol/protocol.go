// Lab 1 - clock synchronization
// File: protocol/protocol.go
// Authors: Jael Dubey, Luc Wachter
// Go version: 1.13.1 (linux/amd64)

// The protocol package contains the constants and types that define the synchronization protocol
// It also provides functions to encode and send messages, and to decode them
package protocol

import (
	"bytes"
	"encoding/binary"
	"log"
	"net"
	"strings"
	"time"
)

// Networking values
const (
	MulticastAddress  = "224.97.6.27:2204"
	UnicastMasterPort = ":2205"
	UnicastSlavePort  = ":2206"
	SyncPeriod        = 2 // [s] Period between synchronizations
	MaxBufferSize     = 32
)

// Message type codes (unsigned bytes for brevity)
const (
	Sync          uint8 = 0
	FollowUp      uint8 = 1
	DelayRequest  uint8 = 2
	DelayResponse uint8 = 3
)

type SyncMessage struct {
	MessageCode uint8
	Id          uint8
}

type FollowUpMessage struct {
	MessageCode uint8
	Id          uint8
	Time        int64
}

type DelayRequestMessage struct {
	MessageCode uint8
	Id          uint8
}

type DelayResponseMessage struct {
	MessageCode uint8
	Id          uint8
	Time        int64
}

// Encode given struct as bytes and return bytes buffer
func encode(message interface{}) *bytes.Buffer {
	buffer := &bytes.Buffer{}
	// Write struct's data as bytes
	err := binary.Write(buffer, binary.BigEndian, message)
	if err != nil {
		log.Fatal(err)
	}

	return buffer
}

// Send SYNC message to multicast group
func SendSync(id uint8) {
	// Build message, encode and send
	message := SyncMessage{
		MessageCode: Sync,
		Id:          id,
	}

	encoded := encode(message)
	sendMulticast(encoded)
}

// Decode sync bytes buffer and return code and id
func SyncDecode(buffer string) (uint8, uint8) {
	message := SyncMessage{}
	err := binary.Read(strings.NewReader(buffer), binary.BigEndian, &message)
	if err != nil {
		log.Fatal(err)
	}

	return message.MessageCode, message.Id
}

// Send FOLLOW_UP message to multicast group
func SendFollowUp(id uint8, tMaster time.Time) {
	// Build message, encode and send
	message := FollowUpMessage{
		MessageCode: FollowUp,
		Id:          id,
		Time:        tMaster.UnixNano() / int64(time.Microsecond),
	}

	encoded := encode(message)
	sendMulticast(encoded)
}

// Decode follow up bytes buffer and return code, id and master time
func FollowUpDecode(buffer string) (uint8, uint8, int64) {
	message := FollowUpMessage{}
	err := binary.Read(strings.NewReader(buffer), binary.BigEndian, &message)
	if err != nil {
		log.Fatal(err)
	}

	return message.MessageCode, message.Id, message.Time
}

// Send DELAY_REQUEST message to specified ip
func SendDelayRequest(ip net.Addr, id uint8) {
	// Build message, encode and send
	message := DelayRequestMessage{
		MessageCode: DelayRequest,
		Id:          id,
	}

	encoded := encode(message)
	sendUnicast(ip, UnicastMasterPort, encoded)
}

// Decode delay request bytes buffer and return code and id
func DelayRequestDecode(buffer string) (uint8, uint8) {
	message := DelayRequestMessage{}
	err := binary.Read(strings.NewReader(buffer), binary.BigEndian, &message)
	if err != nil {
		log.Fatal(err)
	}

	return message.MessageCode, message.Id
}

// Send DELAY_RESPONSE message to specified ip
func SendDelayResponse(ip net.Addr, id uint8, tM time.Time) {
	// Build message, encode and send
	message := DelayResponseMessage{
		MessageCode: DelayResponse,
		Id:          id,
		Time:        tM.UnixNano() / int64(time.Microsecond),
	}

	encoded := encode(message)
	sendUnicast(ip, UnicastSlavePort, encoded)
}

// Decode delay response bytes buffer and return code, id and master time
func DelayResponseDecode(buffer string) (uint8, uint8, int64) {
	message := DelayResponseMessage{}
	err := binary.Read(strings.NewReader(buffer), binary.BigEndian, &message)
	if err != nil {
		log.Fatal(err)
	}

	return message.MessageCode, message.Id, message.Time
}
