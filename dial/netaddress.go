package dial

import (
	"time"
)

type NetAddress struct {
	Network string
	Address string
}

type NetAddressTimeout struct {
	NetAddress
	Timeout time.Duration
}
