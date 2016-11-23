package main

import (
	"time"

	"github.com/hybridgroup/gobot"
)

func main() {
	gbot := gobot.NewGobot()

	var t int
	inc := 1
	MonitorInit()
	Watch(&t, "Count")
	Control(&inc, "Increment", 5, -5)

	propserial := NewSerialAdaptor("propeller", "/dev/ttyAMA0")
	motor := NewMotorDriver(propserial, "motor")
	joystickudp := NewUDPAdaptor("joystick UDP", ":10000")
	remotejoystick := NewPCJoystickDriver(joystickudp, "PC joystick")

	work := func() {
		//gobot.Every(time.Millisecond*10, func() { motor.GetSpeed() })
		gobot.Every(time.Millisecond*300, func() { t += inc })

		remotejoystick.On(remotejoystick.Event(Joystick), func(data interface{}) {
			j := data.(JoystickData)
			if j.deadManHandle {
				gbot.Stop()
			}
			motor.Speed(int16(j.posY) * 200)
		})

	}

	robot := gobot.NewRobot("PiBal",
		[]gobot.Connection{propserial, joystickudp},
		[]gobot.Device{motor, remotejoystick},
		work,
	)
	gbot.AddRobot(robot)

	gbot.Start()
}
