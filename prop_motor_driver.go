package main

import (
	"github.com/hybridgroup/gobot"
)

// MotorDriver Represents a Motor
type MotorDriver struct {
	name             string
	connection       SerialWriter
	CurrentState     byte
	CurrentSpeed     byte
	CurrentDirection string
}

// NewMotorDriver return a new MotorDriver given a SerialWriter, name and pin
func NewMotorDriver(a SerialWriter, name string) *MotorDriver {
	return &MotorDriver{
		name:             name,
		connection:       a,
		CurrentState:     0,
		CurrentSpeed:     0,
		CurrentDirection: "forward",
	}
}

// Name returns the MotorDrivers name
func (m *MotorDriver) Name() string { return m.name }

// Connection returns the MotorDrivers Connection
func (m *MotorDriver) Connection() gobot.Connection { return m.connection.(gobot.Connection) }

// Start implements the Driver interface
func (m *MotorDriver) Start() (errs []error) { return }

// Halt implements the Driver interface
func (m *MotorDriver) Halt() (errs []error) { return }

// Off turns the motor off or sets the motor to a 0 speed
func (m *MotorDriver) Off() (err error) {
	err = m.Speed(0)
	return
}

// On turns the motor on or sets the motor to a maximum speed
func (m *MotorDriver) On() (err error) {
	if m.CurrentSpeed == 0 {
		m.CurrentSpeed = 255
	}
	err = m.Speed(m.CurrentSpeed)
	return
}

// Min sets the motor to the minimum speed
func (m *MotorDriver) Min() (err error) {
	return m.Off()
}

// Max sets the motor to the maximum speed
func (m *MotorDriver) Max() (err error) {
	return m.Speed(255)
}

// IsOn returns true if the motor is on
func (m *MotorDriver) IsOn() bool {
	return m.CurrentSpeed > 0
}

// IsOff returns true if the motor is off
func (m *MotorDriver) IsOff() bool {
	return !m.IsOn()
}

// Toggle sets the motor to the opposite of it's current state
func (m *MotorDriver) Toggle() (err error) {
	if m.IsOn() {
		err = m.Off()
	} else {
		err = m.On()
	}
	return
}

// Speed sets the speed of the motor
func (m *MotorDriver) Speed(value byte) (err error) {
	if writer, ok := m.connection.(SerialWriter); ok {
		m.CurrentSpeed = value
		return writer.SerialWrite("")
	}
	return ErrSerialWriteUnsupported
}

// Forward sets the forward pin to the specified speed
func (m *MotorDriver) Forward(speed byte) (err error) {
	err = m.Direction("forward")
	if err != nil {
		return
	}
	err = m.Speed(speed)
	if err != nil {
		return
	}
	return
}

// Backward sets the backward pin to the specified speed
func (m *MotorDriver) Backward(speed byte) (err error) {
	err = m.Direction("backward")
	if err != nil {
		return
	}
	err = m.Speed(speed)
	if err != nil {
		return
	}
	return
}

// Direction sets the direction pin to the specified speed
func (m *MotorDriver) Direction(direction string) (err error) {
	/*
	m.CurrentDirection = direction
	if m.DirectionPin != "" {
		var level byte
		if direction == "forward" {
			level = 1
		} else {
			level = 0
		}
		err = m.connection.SerialWrite(m.DirectionPin)
	} else {
		var forwardLevel, backwardLevel byte
		switch direction {
		case "forward":
			forwardLevel = 1
			backwardLevel = 0
		case "backward":
			forwardLevel = 0
			backwardLevel = 1
		case "none":
			forwardLevel = 0
			backwardLevel = 0
		}
		err = m.connection.SerialWrite(m.ForwardPin)
		if err != nil {
			return
		}
		err = m.connection.SerialWrite(m.BackwardPin)
		if err != nil {
			return
		}
	}*/
	return
}

func (m *MotorDriver) changeState(state byte) (err error) {
	m.CurrentState = state
	if state == 1 {
		m.CurrentSpeed = 0
	} else {
		m.CurrentSpeed = 255
	}
	err = m.connection.SerialWrite("")

	return
}
