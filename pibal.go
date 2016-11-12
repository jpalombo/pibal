package main

import (
	"bytes"
	"log"
	"time"

	"github.com/hybridgroup/gobot"
)

func main() {
	gbot := gobot.NewGobot()

	propeller := NewSerialAdaptor("propeller", "/dev/ttyAMA0")
	motor := NewMotorDriver(propeller, "motor")
	gbot.On(propeller.Event(Data), func(data interface{}) {
		inbytes := data.([]byte)
		if inbytes[0] == '>' {
		}
		if bytes.Equal(inbytes[0:3], []byte{'+', 'g', 's'}) {
			propeller.SerialWrite("+gs")
		}
		log.Printf("%q", data)
	})

	work := func() {
		propeller.SerialWrite("+ss 0 1000")
		propeller.SerialWrite("+ss 1 1000")
		propeller.SerialWrite("+ss 2 1000")
		propeller.SerialWrite("+ss 3 1000")
		propeller.SerialWrite("+gs")
		time.Sleep(2 * time.Second)
		propeller.SerialWrite("+s")
	}

	robot := gobot.NewRobot("PiBal",
		[]gobot.Connection{propeller},
		[]gobot.Device{motor},
		work,
	)
	gbot.AddRobot(robot)

	gbot.Start()
}
