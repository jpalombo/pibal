package main

import "github.com/hybridgroup/gobot"

var (
	monitorudp *UDPAdaptor
	// Monitor global used to monitor variables
	Monitor *MonitorDriver
)

func main() {
	gbot := gobot.NewGobot()
	monitorudp = NewUDPAdaptor("monitor UDP", ":25045")
	Monitor = NewMonitorDriver(monitorudp, "monitor")

	propserial := NewSerialAdaptor("propeller", "/dev/ttyAMA0")
	motor := NewMotorDriver(propserial, "motor")
	joystickudp := NewUDPAdaptor("joystick UDP", ":10000")
	remotejoystick := NewPCJoystickDriver(joystickudp, "PC joystick")
	bluetooth := NewBluetoothAdapter("bluetooth", "/dev/input/js0")
	bluetoothjoystick := NewBluetoothDriver(bluetooth, "Bluetooth joystick")
	mpu9250 := NewMPU9250Driver("MPU9250", "I2C")
	balance := NewBalanceDriver(mpu9250, "Balance")

	work := func() {
		motor.Stop()

		deadzone := func(i int16) int16 {
			if i < 15 && i > -15 {
				i = 0
			}
			return i
		}
		// Handler for commands from a remote joystick
		joystickhandler := func(data interface{}) {
			j := data.(JoystickData)
			if j.deadManHandle {
				gbot.Stop()
			}
			motor.Speed(
				deadzone(int16(j.posY*200+j.posX*50)),
				deadzone(int16(j.posY*200-j.posX*50)))
		}

		remotejoystick.On(remotejoystick.Event(Joystick), joystickhandler)
		bluetoothjoystick.On(bluetoothjoystick.Event(Joystick), joystickhandler)
		balance.On(balance.Event(Balance), func(data interface{}) {
			b := data.(int)
			motor.Speed(int16(b), int16(b))
		})
	}

	robot := gobot.NewRobot("PiBal",
		[]gobot.Connection{monitorudp, propserial, joystickudp, bluetooth, mpu9250},
		[]gobot.Device{Monitor, motor, remotejoystick, bluetoothjoystick, balance},
		work,
	)
	gbot.AddRobot(robot)

	gbot.Start()
}
