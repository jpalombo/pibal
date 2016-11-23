package main

import (
	"encoding/json"
	"log"

	"github.com/hybridgroup/gobot"
)

// PCJoystickDriver Represents a Joystick
type PCJoystickDriver struct {
	name       string
	connection UDPWriter
	lastPosX   float64
	lastPosY   float64
	lastDMH    bool
	gobot.Eventer
}

// NewPCJoystickDriver return a new PCJoystickDriver given a UDPWriter and name
func NewPCJoystickDriver(u UDPWriter, name string) *PCJoystickDriver {
	j := &PCJoystickDriver{
		name:       name,
		connection: u,
		Eventer:    gobot.NewEventer(),
	}
	j.AddEvent(Joystick)

	if eventer, ok := u.(gobot.Eventer); ok {
		eventer.On(eventer.Event(Data), func(data interface{}) {
			if inbytes, ok := data.([]byte); ok {
				j.parseJoystickData(inbytes)
			}
		})
	}
	return j
}

// Name returns the PCJoystickDrivers name
func (j *PCJoystickDriver) Name() string { return j.name }

// Connection returns the PCJoystickDrivers Connection
func (j *PCJoystickDriver) Connection() gobot.Connection { return j.connection.(gobot.Connection) }

// Start implements the Driver interface
func (j *PCJoystickDriver) Start() (errs []error) { return }

// Halt implements the Driver interface
func (j *PCJoystickDriver) Halt() (errs []error) { return }

// Speed sets the motor speeds
func (j *PCJoystickDriver) Speed(value ...int16) (err error) {
	return
}

func (j *PCJoystickDriver) parseJoystickData(data []byte) {

	var jdat map[string]interface{}
	if err := json.Unmarshal(data, &jdat); err != nil {
		log.Panic(err)
	}

	var newPosX, newPosY float64
	var newDMH bool

	switch jdat["controller"] {
	case "keypad":
		//Extract keypad data
		newPosX = jdat["K_RIGHT"].(float64) - jdat["K_LEFT"].(float64)
		newPosY = jdat["K_UP"].(float64) - jdat["K_DOWN"].(float64)
		newDMH = jdat["K_SPACE"].(float64) != 0.0

	case "Wireless Controller":
		//Extract Playstation Controller data
		log.Println(jdat)
		/*
			newPosX = float32(jdat["sticks"][2])
			newPosY = float32(jdat["sticks"][3])
			newDMH = (jdat["'buttons"][6] + jdat["buttons"][7]) > 0
		*/

	case "Controller (XBOX 360 For Windows)":
		// Extract XBOX 360 controller data
		log.Println(jdat)
		/*
			newPosX = float32(jdat["sticks"][4])
			newPosY = float32(jdat["sticks"][3])
			newDMH = (jdat["'buttons"][6] + jdat["buttons"][7]) > 0
		*/
	}

	if newPosX != j.lastPosX || newPosY != j.lastPosY || newDMH != j.lastDMH {
		j.lastPosX = newPosX
		j.lastPosY = newPosY
		j.lastDMH = newDMH
		j.Publish(Joystick, JoystickData{newPosX, newPosY, newDMH})
	}
}
