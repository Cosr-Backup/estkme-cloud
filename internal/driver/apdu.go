package driver

import "github.com/damonto/libeuicc-go"

type APDU interface {
	libeuicc.APDU
	Lock() error
	Unlock() error
	Receive() chan []byte
}
