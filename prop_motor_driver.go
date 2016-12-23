package main

import (
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/hybridgroup/gobot"
	"golang.org/x/exp/io/i2c"
)

// MotorDriver Represents a Motor
type MotorDriver struct {
	name       string
	connection SerialWriter
	//CurrentSpeed [4]int16
	//interlock    *Interlock
	gobot.Eventer
}

// NewMotorDriver return a new MotorDriver given a SerialWriter, name and pin
func NewMotorDriver(a SerialWriter, name string) *MotorDriver {

	m := &MotorDriver{
		name:       name,
		connection: a,
		//interlock:  NewInterlock(2, 100),
		Eventer: gobot.NewEventer(),
	}
	m.AddEvent(MotorSpeed)

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
	var newSpeed [4]int16
	switch len(value) {
	case 0: // nothing to do, newSpeed already 0
	case 1:
		for i := range newSpeed {
			newSpeed[i] = value[0]
		}
	case 2:
		newSpeed[0] = value[0]
		newSpeed[1] = 0 //value[0]
		newSpeed[2] = 0 //value[1]
		newSpeed[3] = value[1]
	case 4:
		for i := range newSpeed {
			newSpeed[i] = value[i]
		}
	}

	/*outstring := fmt.Sprintf("+sa %d %d %d %d",
		newSpeed[0],
		newSpeed[1],
		newSpeed[2],
		newSpeed[3])

	//	if writer, ok := m.connection.(SerialWriter); ok {
	//		return m.interlock.Write(outstring, writer.SerialWrite)
	//	}
	//	return ErrSerialWriteUnsupported */

	// Write directly using I2C rather than serial interface for increased speed
	d, err := i2c.Open(&i2c.Devfs{Dev: "/dev/i2c-1"}, 0x42)
	if err != nil {
		return err
	}
	var regs [8]byte
	for i := 0; i < 4; i++ {
		regs[i*2] = byte(Abs(int(newSpeed[i])))
		if newSpeed[i] > 0 {
			regs[i*2+1] = 0
		} else {
			regs[i*2+1] = 1
		}
	}
	if err = d.WriteReg(0, regs[:]); err != nil {
		return err
	}
	d.Close()
	return
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

// SetPid sets the motor pid parameters
func (m *MotorDriver) SetPid(kp, ki, kd int) (err error) {
	if writer, ok := m.connection.(SerialWriter); ok {
		return writer.SerialWrite(fmt.Sprintf("+sp %d %d %d", kp, ki, kd))
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
		mdata := MotorSpeedData{}
		for i := 0; i < 4; i++ {
			mdata.speed[i], _ = strconv.Atoi(split[i+1])
		}
		mdata.millis, _ = strconv.Atoi(split[5])
		m.Publish(MotorSpeed, mdata)
	case "+gp:":
		log.Printf("position: %s %s %s %s %s",
			split[1],
			split[2],
			split[3],
			split[4],
			split[5])
	case "+sa:":
		//m.interlock.ResponseRcvd()
	default:
		log.Printf("%s", data)
	}
}
