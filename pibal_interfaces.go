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
	// Balance event
	Balance = "balance"
)

// BluetoothData struct for data sent with the Joystick event
// see https://www.kernel.org/doc/Documentation/input/joystick-api.txt for
// meaning of fields
type BluetoothData struct {
	value   int16
	bType   byte
	bNumber byte
}

// JoystickData struct for data sent with the Joystick event
type JoystickData struct {
	posX          float64
	posY          float64
	deadManHandle bool
}

// MotorSpeedData struct for data sent with the MotorSpeed event
type MotorSpeedData struct {
	speed  [4]int
	millis int
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

// MPU9250Sender interface
type MPU9250Sender interface {
	gobot.Adaptor
	SensorAccel(int) (int, error)
	SensorGyro(int) (int, error)
}
