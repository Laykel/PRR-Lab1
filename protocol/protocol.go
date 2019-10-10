package protocol

import (
    "fmt"
    "net"
    "time"
)

// Networking values
const (
	MulticastAddress     = "224.0.0.1:2204"
	UnicastListenAddress = ":2205"
	SyncPeriod           = 4 // [s] Period between to synchronizations
	MaxBufferSize        = 256
)

// Message type codes (unsigned bytes for brevity)
const (
	Sync          uint8 = 0
	FollowUp      uint8 = 1
	DelayRequest  uint8 = 2
	DelayResponse uint8 = 3
)

// SYNC (message code + ID)
func SendSync(id uint) {
    // Build message and send
    message := fmt.Sprintf("%d|%d", Sync, id)
    sendMulticast(message)
}

// FOLLOW_UP (message code + ID + tMaster)
func SendFollowUp(id uint) {
    // Syscall for time
    tMaster := time.Now()

    // Build message and send
    message := fmt.Sprintf("%d|%d|%s", FollowUp, id, tMaster)
    sendMulticast(message)
}

// DELAY_REQUEST (message code)
func SendDelayRequest(ip net.Addr) {
    // Build message and send
    message := fmt.Sprint(DelayRequest)
    sendUnicast(ip, message)
}

// DELAY_RESPONSE (message code, time of request's reception)
func SendDelayResponse(ip net.Addr, tM time.Time) {
    // Build message and send
    message := fmt.Sprintf("%d|%s", DelayResponse, tM)
    sendUnicast(ip, message)
}
