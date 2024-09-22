package cloud

import (
	"log/slog"
	"sync"
	"time"

	"github.com/damonto/libeuicc-go"
)

var (
	APDUCommunicateTimeout = []byte{0x66, 0x00}
)

type APDU interface {
	libeuicc.APDU
	Receive() chan []byte
}

type apdu struct {
	lock     sync.Mutex
	conn     *Conn
	receiver chan []byte
}

func NewAPDU(conn *Conn) APDU {
	return &apdu{conn: conn, receiver: make(chan []byte, 1)}
}

func (a *apdu) Connect() error {
	return a.conn.Send(TagAPDULock, nil)
}

func (a *apdu) Disconnect() error {
	return a.conn.Send(TagAPDUUnlock, nil)
}

func (a *apdu) OpenLogicalChannel(_ []byte) (int, error) {
	return 0, nil
}

func (a *apdu) CloseLogicalChannel(channel int) error {
	return nil
}

func (a *apdu) Transmit(command []byte) ([]byte, error) {
	a.lock.Lock()
	defer a.lock.Unlock()

	if err := a.conn.Send(TagAPDU, command); err != nil {
		return nil, err
	}

	select {
	case r := <-a.receiver:
		return r, nil
	case <-time.After(2 * time.Minute): // If response is not received in 2 minutes, return a timeout error.
		slog.Warn("wait for APDU command response timeout", "conn", a.conn.Id, "command", command, "response", APDUCommunicateTimeout)
		return APDUCommunicateTimeout, nil
	}
}

func (a *apdu) Receive() chan []byte {
	return a.receiver
}
