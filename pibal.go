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

	work := func() {
		propeller.On(propeller.Event(Data), func(data interface{}) {
			log.Printf("%q", data)
			inbytes := data.([]byte)
			if inbytes[0] == '>' {
				inbytes = inbytes[1:]
			}
			if bytes.Equal(inbytes[0:3], []byte{'+', 'g', 's'}) {
				propeller.SerialWrite("+gs")
			}
			log.Printf("%q", inbytes)
		})

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
