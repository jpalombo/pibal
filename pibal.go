package main

import (
	"time"

	"github.com/hybridgroup/gobot"
)

func main() {
	gbot := gobot.NewGobot()

	t := 10
	inc := 1

	propserial := NewSerialAdaptor("propeller", "/dev/ttyAMA0")
	motor := NewMotorDriver(propserial, "motor")
	monitorudp := NewUDPAdaptor("monitor UDP", ":25045")
	monitor := NewMonitorDriver(monitorudp, "monitor")
	joystickudp := NewUDPAdaptor("joystick UDP", ":10000")
	remotejoystick := NewPCJoystickDriver(joystickudp, "PC joystick")

	monitor.Watch(&t, "Count")
	monitor.Control(&inc, "Increment", 5, -5)

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
		[]gobot.Connection{propserial, joystickudp, monitorudp},
		[]gobot.Device{motor, remotejoystick, monitor},
		work,
	)
	gbot.AddRobot(robot)

	gbot.Start()
}
