/*
A simple animated rainbow example.

This is a port of https://github.com/pimoroni/mote/blob/master/python/examples/rainbow.py
*/
package main

import (
	"os"
	"os/signal"
	"syscall"
	"time"

	"fmt"

	"github.com/johnmccabe/mote"
	"github.com/lucasb-eyer/go-colorful"
)

func main() {
	mote := mote.NewMote("")

	c := make(chan os.Signal, 2)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		mote.ClearAll()
		mote.Show()
		mote.Close()
		fmt.Printf("EXITING: %v", mote)
		for c := range mote.Channels {

			fmt.Printf("CHAN: %v", c)
		}
		os.Exit(1)
	}()

	mote.ConfigureChannel(1, 16, false)
	mote.ConfigureChannel(2, 16, false)
	mote.ConfigureChannel(3, 16, false)
	mote.ConfigureChannel(4, 16, false)

	for {
		t := int(float64(time.Now().UnixNano()) * 0.00000005)
		for channel := 1; channel < 5; channel++ {
			for pixel := 0; pixel < 16; pixel++ {
				hue := (t + (channel * 64) + (pixel * 4)) % 360
				hsv := colorful.Hsv(float64(hue), 1.0, 1.0)
				r := int(hsv.R * 255)
				g := int(hsv.G * 255)
				b := int(hsv.B * 255)
				mote.SetPixel(channel, pixel, r, g, b)
			}
		}
		time.Sleep(10 * time.Millisecond)
		mote.Show()
	}
}

// func cleanup(m *mote.Mote) {
// 	m.ClearAll()
// 	fmt.Println("clear show")
// 	m.Show()
// 	time.Sleep(20 * time.Millisecond)

// 	m.Close()
// }
