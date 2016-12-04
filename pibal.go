package main

import "github.com/hybridgroup/gobot"

func main() {
	gbot := gobot.NewGobot()

	propserial := NewSerialAdaptor("propeller", "/dev/ttyAMA0")
	motor := NewMotorDriver(propserial, "motor")
	//monitorudp := NewUDPAdaptor("monitor UDP", ":25045")
	//monitor := NewMonitorDriver(monitorudp, "monitor")
	joystickudp := NewUDPAdaptor("joystick UDP", ":10000")
	remotejoystick := NewPCJoystickDriver(joystickudp, "PC joystick")
	bluetooth := NewBluetoothAdapter("bluetooth", "/dev/input/js0")
	bluetoothjoystick := NewBluetoothDriver(bluetooth, "Bluetooth joystick")

	work := func() {
		motor.Stop()

		// Handler for commands from a remote joystick
		joystickhandler := func(data interface{}) {
			j := data.(JoystickData)
			if j.deadManHandle {
				gbot.Stop()
			}
			motor.Speed(int16(j.posY*200+j.posX*50), int16(j.posY*200-j.posX*50))
		}
		remotejoystick.On(remotejoystick.Event(Joystick), joystickhandler)
		bluetoothjoystick.On(bluetoothjoystick.Event(Joystick), joystickhandler)
	}

	robot := gobot.NewRobot("PiBal",
		//[]gobot.Connection{propserial, joystickudp, monitorudp},
		[]gobot.Connection{propserial, joystickudp, bluetooth},
		//[]gobot.Device{motor, remotejoystick, monitor},
		[]gobot.Device{motor, remotejoystick, bluetoothjoystick},
		work,
	)
	gbot.AddRobot(robot)

	gbot.Start()
}
