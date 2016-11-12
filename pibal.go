package main

import (
	"log"

	"github.com/hybridgroup/gobot"
)

func main() {
	gbot := gobot.NewGobot()

	propeller := NewSerialAdaptor("propeller", "/dev/ttyAMA0")
	motor := NewMotorDriver(propeller, "motor")
	gbot.On(propeller.Event(Data), func(data interface{}) {
		log.Printf("%q", data)
	})

	work := func() {
		propeller.SerialWrite("+gs")
		propeller.SerialWrite("+gp")
	}

	robot := gobot.NewRobot("PiBal",
		[]gobot.Connection{propeller},
		[]gobot.Device{motor},
		work,
	)
	gbot.AddRobot(robot)

	gbot.Start()
}
