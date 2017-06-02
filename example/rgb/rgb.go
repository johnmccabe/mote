/*
A simple static RGB colour setting example.

This is a port of https://github.com/pimoroni/mote/blob/master/python/examples/rgb.py
*/
package main

import (
	"time"

	"github.com/johnmccabe/mote"
)

func main() {
	mote := mote.NewMote("")
	mote.ConfigureChannel(1, 16, false)
	mote.ConfigureChannel(2, 16, false)
	mote.ConfigureChannel(3, 16, false)
	mote.ConfigureChannel(4, 16, false)

	for channel := 1; channel < 5; channel++ {
		for pixel := 0; pixel < 16; pixel++ {
			mote.SetPixel(channel, pixel, 0, 0, 0)
		}
		time.Sleep(10 * time.Millisecond)
	}

	mote.Show()
	time.Sleep(100 * time.Millisecond)
	mote.Close()
}
