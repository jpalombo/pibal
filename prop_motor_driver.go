package main

import (
	"fmt"
	"log"
	"strings"

	"github.com/hybridgroup/gobot"
)

// MotorDriver Represents a Motor
type MotorDriver struct {
	name         string
	connection   SerialWriter
	CurrentSpeed [4]int16
}

// NewMotorDriver return a new MotorDriver given a SerialWriter, name and pin
func NewMotorDriver(a SerialWriter, name string) *MotorDriver {

	m := &MotorDriver{
		name:       name,
		connection: a,
	}

	if eventer, ok := a.(gobot.Eventer); ok {
		eventer.On(eventer.Event(Data), func(data interface{}) {
			if inbytes, ok := data.([]byte); ok {
				if inbytes[0] == '>' {
					inbytes = inbytes[1:]
				}
				m.parseReadData(string(inbytes))
			}
		})
	}

	return m
}

// Name returns the MotorDrivers name
func (m *MotorDriver) Name() string { return m.name }

// Connection returns the MotorDrivers Connection
func (m *MotorDriver) Connection() gobot.Connection { return m.connection.(gobot.Connection) }

// Start implements the Driver interface
func (m *MotorDriver) Start() (errs []error) { return }

// Halt implements the Driver interface
func (m *MotorDriver) Halt() (errs []error) { return }

// Speed sets the motor speeds
func (m *MotorDriver) Speed(value ...int16) (err error) {
	if value == nil { // called with zero parameters
		for i := range m.CurrentSpeed {
			m.CurrentSpeed[i] = 0
		}
	}
	switch len(value) {
	case 1:
		for i := range m.CurrentSpeed {
			m.CurrentSpeed[i] = value[0]
		}
	case 2:
		m.CurrentSpeed[0] = value[0]
		m.CurrentSpeed[1] = value[0]
		m.CurrentSpeed[2] = value[1]
		m.CurrentSpeed[3] = value[1]
	case 4:
		for i := range m.CurrentSpeed {
			m.CurrentSpeed[i] = value[i]
		}
	}

	outstring := fmt.Sprintf("+sa %d %d %d %d",
		m.CurrentSpeed[0],
		m.CurrentSpeed[1],
		m.CurrentSpeed[2],
		m.CurrentSpeed[3])
	if writer, ok := m.connection.(SerialWriter); ok {
		return writer.SerialWrite(outstring)
	}

	return ErrSerialWriteUnsupported
}

// GetSpeed gets the motor speeds
func (m *MotorDriver) GetSpeed() (err error) {
	if writer, ok := m.connection.(SerialWriter); ok {
		return writer.SerialWrite("+gs")
	}
	return ErrSerialWriteUnsupported
}

// GetPosition gets the motor speeds
func (m *MotorDriver) GetPosition() (err error) {
	if writer, ok := m.connection.(SerialWriter); ok {
		return writer.SerialWrite("+gp")
	}
	return ErrSerialWriteUnsupported
}

// Stop stops the motor
func (m *MotorDriver) Stop() (err error) {
	if writer, ok := m.connection.(SerialWriter); ok {
		return writer.SerialWrite("+st")
	}
	return ErrSerialWriteUnsupported
}

func (m *MotorDriver) parseReadData(data string) {
	split := strings.Split(data, " ")
	switch split[0] {
	case "+gs:":
		log.Printf("speed: %s %s %s %s %s",
			split[1],
			split[2],
			split[3],
			split[4],
			split[5])
	case "+gp:":
		log.Printf("position: %s %s %s %s %s",
			split[1],
			split[2],
			split[3],
			split[4],
			split[5])
	default:
		log.Printf("%s", data)
	}
}
