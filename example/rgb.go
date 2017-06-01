package main

import (
	"time"

	"github.com/johnmccabe/mote"
)

func main() {
	mote := mote.NewMote("")
	mote.ConfigureChannel(1, 8, false)
	// mote.ConfigureChannel(2, 16, false)
	// mote.ConfigureChannel(3, 16, false)
	// mote.ConfigureChannel(4, 16, false)

	for channel := 1; channel < 5; channel++ {
		for pixel := 0; pixel < 16; pixel++ {
			mote.SetPixel(channel, pixel, 0, 0, 0)
		}
		time.Sleep(10 * time.Millisecond)
	}

	mote.Show()
	// for channel in range(4):
	// for pixel in range(16):
	//     mote.set_pixel(channel + 1, pixel, r, g, b)
	time.Sleep(100 * time.Millisecond)
	mote.Close()
}
