package main

import (
	"log"
	"time"

	"github.com/johnmccabe/mote"
	"github.com/lucasb-eyer/go-colorful"
)

func main() {
	mote := mote.NewMote("")
	mote.ConfigureChannel(1, 16, false)
	mote.ConfigureChannel(2, 16, false)
	mote.ConfigureChannel(3, 16, false)
	mote.ConfigureChannel(4, 16, false)

	i := 0

	for {
		// t := int(time.Now().Unix())
		t := int(float64(time.Now().UnixNano()) * 0.00000005)
		// h := t * 50
		for channel := 1; channel < 5; channel++ {
			// log.Printf("Channel: %d\n", channel)
			for pixel := 0; pixel < 16; pixel++ {
				hue := (t + (channel * 64) + (pixel * 4)) % 360
				hsv := colorful.Hsv(float64(hue), 1.0, 1.0)
				// log.Printf("Hsv: %v\n", hsv)
				r := int(hsv.R * 255)
				g := int(hsv.G * 255)
				b := int(hsv.B * 255)
				if pixel == 0 {
					log.Printf("Pixel: 0, Hue: %v, R: %d, G: %d, B: %d\n", hue, r, g, b)
				}
				mote.SetPixel(channel, pixel, r, g, b)

				// log.Printf("i: %d, t: %d\n", i, t)
			}
		}

		time.Sleep(2 * time.Millisecond)
		log.Printf("i: %d, t: %d\n", i, t)
		i++
		mote.Show()
	}
	// for channel in range(4):
	// for pixel in range(16):
	//     mote.set_pixel(channel + 1, pixel, r, g, b)
	// time.Sleep(100 * time.Millisecond)
	// mote.Close()
}
