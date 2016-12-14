package main

import (
	"log"
	"sync"
	"time"
)

// NewInterlock returns a new interlock
func NewInterlock(size int, timeout int) *Interlock {
	return &Interlock{size: size, timeout: timeout}
}

// Interlock is a basic LIFO stack that resizes as needed.
type Interlock struct {
	size     int
	count    int
	str      string
	writer   func(string) error
	timeout  int
	lastResp time.Time
	m        sync.Mutex
}

func (i *Interlock) checkTimeout() {
	if i.count > 0 && int(time.Since(i.lastResp)/time.Millisecond) > i.timeout {
		// assume that we've dropped some responses somewhere
		log.Println("Prop rcv has timed out, resetting")
		i.count = 0
	}
}

// ResponseRcvd - free up a space in the interlock
func (i *Interlock) ResponseRcvd() (err error) {
	i.m.Lock()
	i.lastResp = time.Now()
	if i.writer != nil {
		err = i.writer(i.str)
		i.writer = nil
	} else if i.count > 0 {
		i.count--
	}
	i.m.Unlock()
	return
}

// Write a string or store it if not enough space
func (i *Interlock) Write(str string, writer func(string) error) (err error) {
	//log.Println("Write", i.count, str)
	i.m.Lock()
	if i.count == 0 {
		i.lastResp = time.Now()
	} else {
		i.checkTimeout()
	}
	if i.count < i.size {
		i.count++
		err = writer(str)
	} else {
		// Not enough space to send, so hold on for later
		// This overwrites any previously stored string as we only want to send the latest
		i.str = str
		i.writer = writer
	}
	i.m.Unlock()
	return
}
