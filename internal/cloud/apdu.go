package cloud

import (
	"encoding/hex"
	"log/slog"
	"sync"
	"time"

	"github.com/damonto/estkme-cloud/internal/driver"
)

const (
	APDUCommunicateTimeout = "6600"
)

type apdu struct {
	lock     sync.Mutex
	conn     *Conn
	receiver chan []byte
}

func NewAPDU(conn *Conn) driver.APDU {
	return &apdu{conn: conn, receiver: make(chan []byte, 1)}
}

func (a *apdu) Lock() error {
	a.lock.Lock()
	return a.conn.Send(TagAPDULock, nil)
}

func (a *apdu) Unlock() error {
	defer a.lock.Unlock()
	return a.conn.Send(TagAPDUUnlock, nil)
}

func (a *apdu) Connect() error {
	return nil
}

func (a *apdu) Disconnect() error {
	return nil
}

func (a *apdu) OpenLogicalChannel(_ []byte) (int, error) {
	return 0, nil
}

func (a *apdu) CloseLogicalChannel(channel int) error {
	return nil
}

func (a *apdu) Transmit(command []byte) ([]byte, error) {
	b, _ := hex.DecodeString(string(command))
	if err := a.conn.Send(TagAPDU, b); err != nil {
		return nil, err
	}

	select {
	case r := <-a.receiver:
		return r, nil
	case <-time.After(2 * time.Minute): // If response is not received in 2 minutes, return a timeout error.
		slog.Warn("wait for APDU command response timeout", "conn", a.conn.Id, "command", command, "response", APDUCommunicateTimeout)
		return []byte(APDUCommunicateTimeout), nil
	}
}

func (a *apdu) Receive() chan []byte {
	return a.receiver
}
