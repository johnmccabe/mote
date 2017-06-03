/*
A simple example that reads its colour data from the Cheerlights feed.

Tweet a colour with the hashtag `#Cheerlights` to have it change the colour of Motes all over the world.

This is a port of https://github.com/pimoroni/mote/blob/master/python/examples/cheerlights.py
*/
package main

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/johnmccabe/mote"
	"github.com/johnmccabe/mote/example/cheerlights/feed"
)

// CheerlightsURL is the Cheerlights feed
const CheerlightsURL = "http://api.thingspeak.com/channels/1417/feed.json"

func main() {
	mote := mote.NewMote("")

	c := make(chan os.Signal, 2)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		mote.ClearAll()
		mote.Show()
		mote.Close()
		os.Exit(1)
	}()

	mote.ConfigureChannel(1, 16, false)
	mote.ConfigureChannel(2, 16, false)
	mote.ConfigureChannel(3, 16, false)
	mote.ConfigureChannel(4, 16, false)

	for {
		res, _ := http.Get(CheerlightsURL)
		body, _ := ioutil.ReadAll(res.Body)
		var data feed.CheerLights
		err := json.Unmarshal(body, &data)
		if err != nil {
			panic("can't unmarshall")
		}

		fmt.Printf("%s", time.Now().Format("2006/01/02 15:04:05 "))

		for channel := 1; channel < 5; channel++ {
			f := data.Feeds[len(data.Feeds)-(channel*2)]
			fmt.Printf("%v", f)
			b, _ := hex.DecodeString(strings.Trim(f.Field2, "#"))
			for pixel := 0; pixel < 16; pixel++ {
				r := int(b[0])
				g := int(b[1])
				b := int(b[2])
				mote.SetPixel(channel, pixel, r, g, b)
			}
		}
		fmt.Printf("\n")

		mote.Show()
		time.Sleep(5 * time.Second)
	}
}
