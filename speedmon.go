package main

import "log"

var (
	posKp = 100
	posKd = 5
	sP    int
	sD    int
)

// SpeedMon holds data about the monitor driver
type SpeedMon struct {
	targetSpeedCom int
	speedDiff      int
	targetPos      int
	setPos         int
	lastPosError   int
}

// NewSpeedMon initialises the SpeedMon struct
func NewSpeedMon() *SpeedMon {
	Monitor.Control(&posKp, "posKp", 50000, -50000, nil)
	Monitor.Control(&posKd, "posKd", 50000, -50000, nil)
	Monitor.Watch(&sP, "speed P")
	Monitor.Watch(&sD, "speed D")
	return &SpeedMon{}
}

// Reset does a reset (ns!)
func (b *SpeedMon) Reset() {
	b.targetSpeedCom = 0
	b.speedDiff = 0
	b.targetPos = 0
	b.setPos = 0
}

// UpdatePosition feed in new position info
func (b *SpeedMon) UpdatePosition(pos int) (angleoffset int) {
	newpos := pos - b.setPos
	if b.setPos == 0 {
		// Initial setting on startup
		b.setPos = newpos
		return
	}
	// Do some processing to work out new angle offset
	posError := newpos - b.targetPos
	posErrorDelta := posError - b.lastPosError
	b.lastPosError = posError
	sP = (posError * posKp) >> 18
	sD = (posErrorDelta * posKd) >> 18
	angleoffset = sP + sD
	return
}

// RequestSpeed request speed for the robot
func (b *SpeedMon) RequestSpeed(speed int) (angleoffset int) {
	log.Println("Request Speed", speed)
	b.targetSpeedCom = speed

	// Do some processing to work out new motor speeds
	return
}
