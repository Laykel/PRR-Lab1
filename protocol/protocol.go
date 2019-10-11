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
	SyncPeriod           = 4 // [s] Period between synchronizations
	MaxBufferSize        = 256
)

// Message type codes (unsigned bytes for brevity)
const (
	Sync          uint8 = 0
	FollowUp      uint8 = 1
	DelayRequest  uint8 = 2
	DelayResponse uint8 = 3
)

// Send SYNC (message code + ID) message to multicast group
func SendSync(id uint) {
    // Build message and send
    message := fmt.Sprintf("%d|%d", Sync, id)
    sendMulticast(message)
}

// Send FOLLOW_UP (message code + ID + tMaster) message to multicast group
func SendFollowUp(id uint) {
    // Syscall for time
    tMaster := time.Now()

    // Build message and send
    message := fmt.Sprintf("%d|%d|%d", FollowUp, id, tMaster.Unix())
    sendMulticast(message)
}

// Send DELAY_REQUEST (message code) message to specified ip
func SendDelayRequest(ip net.Addr) {
    // Build message and send
    message := fmt.Sprint(DelayRequest)
    sendUnicast(ip, message)
}

// Send DELAY_RESPONSE (message code, time of request's reception) message to specified ip
func SendDelayResponse(ip net.Addr, tM time.Time) {
    // Build message and send
    message := fmt.Sprintf("%d|%d", DelayResponse, tM.Unix())
    sendUnicast(ip, message)
}
