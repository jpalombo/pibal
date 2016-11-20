package main

import (
	"time"

	"github.com/hybridgroup/gobot"
)

func main() {
	gbot := gobot.NewGobot()

	propserial := NewSerialAdaptor("propeller", "/dev/ttyAMA0")
	motor := NewMotorDriver(propserial, "motor")

	work := func() {
		gobot.Every(time.Millisecond*10, func() { motor.GetSpeed() })

		motor.Speed(200)
		time.Sleep(1 * time.Second)
		motor.Speed(-200)
		time.Sleep(1 * time.Second)
		motor.Stop()
		gbot.Stop()
	}

	robot := gobot.NewRobot("PiBal",
		[]gobot.Connection{propserial},
		[]gobot.Device{motor},
		work,
	)
	gbot.AddRobot(robot)

	gbot.Start()
}
