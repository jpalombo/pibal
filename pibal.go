package main

import (
  "github.com/hybridgroup/gobot"
)

func main() {
  gbot := gobot.NewGobot()

  serialAdaptor := NewSerialAdaptor("propeller", "/dev/ttyAMA0")
  motor := NewMotorDriver(serialAdaptor, "motor")

  work := func() {
    serialAdaptor.SerialWrite("+gs")
    serialAdaptor.SerialWrite("+gp")
    serialAdaptor.SerialWrite("")
    serialAdaptor.SerialWrite("")
  }

  robot := gobot.NewRobot("PiBal",
    []gobot.Connection{serialAdaptor},
    []gobot.Device{motor},
    work,
  )
  gbot.AddRobot(robot)

  gbot.Start()
}