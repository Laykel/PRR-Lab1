package protocol

// Networking values
const (
	MulticastAddress = "224.0.0.1:2204"
	//MulticastPort = "2204"
	UnicastPort = "2205"
	SyncPeriod  = 2 // [s] Period between to synchronizations
	// IEEE 1588 suggests 2 seconds
)

// Message type codes
const (
	Sync          uint8 = 0
	FollowUp      uint8 = 1
	DelayRequest  uint8 = 2
	DelayResponse uint8 = 3
)
