package main

import (
	"errors"

	"github.com/hybridgroup/gobot"
)

var (
	// ErrSerialWriteUnsupported is the error resulting when a driver attempts to use
	// hardware capabilities which a connection does not support
	ErrSerialWriteUnsupported = errors.New("SerialWrite is not supported by this platform")
)

const (
	// Error event
	Error = "error"
	// Data event
	Data = "data"
	// Joystick event
	Joystick = "joystick"
	// MotorSpeed event
	MotorSpeed = "motorspeed"
)

// JoystickData struct for data sent with the Joystick event
type JoystickData struct {
	posX          float64
	posY          float64
	deadManHandle bool
}

// Interfaces

// SerialWriter interface represents an Adaptor which has SerialWrite capabilities
type SerialWriter interface {
	gobot.Adaptor
	SerialWrite(string) (err error)
}

// UDPWriter interface represents an Adaptor which has UDPWrite capabilities
type UDPWriter interface {
	gobot.Adaptor
	UDPWrite([]byte) (err error)
}
