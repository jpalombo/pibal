package main

import (
	"errors"

	"github.com/hybridgroup/gobot"
)

const (
	// Error event
	Error = "error"
	// Data event
	Data = "data"
)

var (
	// ErrSerialWriteUnsupported is the error resulting when a driver attempts to use
	// hardware capabilities which a connection does not support
	ErrSerialWriteUnsupported = errors.New("SerialWrite is not supported by this platform")
)

// DigitalWriter interface represents an Adaptor which has DigitalWrite capabilities
type SerialWriter interface {
	gobot.Adaptor
	SerialWrite(string) (err error)
}
