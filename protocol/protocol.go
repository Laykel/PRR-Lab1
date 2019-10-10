package protocol

import (
    "fmt"
    "time"
)

// Networking values
const (
	MulticastAddress     = "224.0.0.1:2204"
	UnicastListenAddress = ":2205"
	SyncPeriod           = 4 // [s] Period between to synchronizations
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
    syncCode := fmt.Sprint(Sync)
    syncId := fmt.Sprint(id)
    sendMulticast(syncCode+"|"+syncId)
}

// FOLLOW_UP (message code + ID + tMaster)
func SendFollowUp(id uint) {
    followUpCode := fmt.Sprint(FollowUp)
    followUpId := fmt.Sprint(id)

    // Syscall for time
    tMaster := time.Now()

    sendMulticast(followUpCode+"|"+followUpId+"|"+tMaster.String())
}
