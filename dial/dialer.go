package dial

import "net"

type Dialer struct {
	semaphore chan struct{}
}

type Status struct {
	Ok  bool
	Err error
}

func NewDialer(concurrentConnections int) *Dialer {
	return &Dialer{
		semaphore: make(chan struct{}, concurrentConnections),
	}
}

func (d *Dialer) NewWorker() (chan<- NetAddressTimeout, <-chan Status) {
	netAddressTimeoutCh := make(chan NetAddressTimeout)
	dialerStatusCh := make(chan Status)
	d.semaphore <- struct{}{}
	go func() {
		netAddressTimeout := <-netAddressTimeoutCh
		conn, err := net.DialTimeout(netAddressTimeout.Network, netAddressTimeout.Address, netAddressTimeout.Timeout)
		dialerStatus := Status{}
		if err != nil {
			dialerStatus.Ok = false
			dialerStatus.Err = err
		} else {
			dialerStatus.Ok = true
			conn.Close()
		}
		dialerStatusCh <- dialerStatus
		<-d.semaphore
	}()
	return netAddressTimeoutCh, dialerStatusCh
}
