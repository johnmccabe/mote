/*
Package mote represents a connected Pimoroni Mote device, communicating over USB serial.

It allows you to configure the 4 channels and set individual pixels, see the `examples` subdirectory for soem demo applications using the library.

It is a port of the Pimoroni Mote Python library (https://github.com/pimoroni/mote).

The Mote device can be obtained directly from Pimoroni (https://shop.pimoroni.com/products/mote).
*/
package mote

import (
	"fmt"
	"io"
	"log"
	"strings"

	"go.bug.st/serial"
	"go.bug.st/serial/enumerator"
)

const (
	// Pimoroni Mote USB Vendor ID
	VID = "16d0"
	// Pimoroni Mote USB Product ID
	PID = "08c4"
	// Pimoroni Mote USB Product Name
	ProductName = "Mote USB Dock"
	// NumChannels is the number of available channel connections on the Mote device
	NumChannels = 4
	// MaxPixels is the maximum addressable number of pixels across all channels
	MaxPixels = 512
	// MaxPixelsPerChannel is the maximum addressable number of pixels across a single channel
	MaxPixelsPerChannel = int(MaxPixels / NumChannels)
)

// Mote represents a connected Pimoroni Mote device
type Mote struct {
	PortName string
	Port     io.WriteCloser
	Channels [NumChannels]*Channel
}

// Pixel represents a single RGB pixel
type Pixel struct {
	Red, Green, Blue int
}

// Channel represents a single channel on the Mote board
type Channel struct {
	Pixels          []Pixel
	GammaCorrection bool
}

// NewMote creates a connection to a Mote device, communicating over USB serial.
//
// It will attach to the first available Mote device unless a non-empty string `port_name` is specified at init.
//
//   - portName: override auto-detect and specify an explicit port to use. Must be a complete path ie: /dev/tty.usbmodem1234
func NewMote(portName string) *Mote {
	mote := Mote{
		PortName: portName,
	}
	if mote.PortName == "" {
		mote.PortName = *findSerialPort(VID, PID)
	}
	if mote.PortName == "" {
		log.Fatal("unable to detect connected Mote")
	}

	mode := serial.Mode{
		BaudRate: 115200,
		DataBits: 8,
		Parity:   serial.NoParity,
		StopBits: serial.OneStopBit,
	}

	p, err := serial.Open(mote.PortName, &mode)
	if err != nil {
		panic(err)
	}
	mote.Port = p

	return &mote
}

// ConfigureChannel configures a channel, taking the following parameters.
//
//   - channel: Channel, either 1, 2, 3 or 4 corresponding to numbers on Mote
//   - numPixels: Number of pixels to configure for this channel
//   - gammaCorrection: Whether to enable gamma correction
func (m *Mote) ConfigureChannel(channel, numPixels int, gammaCorrection bool) error {
	if channel > NumChannels || channel < 1 {
		return fmt.Errorf("channel index must be between 1 and 4")
	}
	if numPixels > MaxPixelsPerChannel {
		return fmt.Errorf("number of pixels can not be more than %d", MaxPixelsPerChannel)
	}

	p := []Pixel{}
	for i := 0; i < numPixels; i++ {
		p = append(p, Pixel{0, 0, 0})
	}
	c := Channel{
		Pixels:          p,
		GammaCorrection: gammaCorrection,
	}
	m.Channels[channel-1] = &c

	var b []byte
	b = append(b, []byte("mote")...)
	b = append(b, []byte("c")...)
	b = append(b, byte(channel))
	b = append(b, byte(numPixels))
	gammaCorrectionVar := 0
	if gammaCorrection {
		gammaCorrectionVar = 1
	}
	b = append(b, byte(gammaCorrectionVar))

	_, err := m.Port.Write(b)
	if err != nil {
		return err
	}
	return nil
}

// SetPixel Set the RGB colour of a single pixel, on a single channel, taking the following parameters.
//
//   - channel: Channel, either 1, 2, 3 or 4 corresponding to numbers on Mote
//   - index: Index of pixel to set, from 0 up
//   - r: Amount of red: 0-255
//   - g: Amount of green: 0-255
//   - b: Amount of blue: 0-255
func (m *Mote) SetPixel(channel, index, r, g, b int) error {
	if channel > NumChannels || channel < 1 {
		return fmt.Errorf("channel index must be between 1 and 4")
	}
	if m.Channels[channel-1] == nil {
		return fmt.Errorf("please set up channel %d before using it", channel)
	}
	if index >= len(m.Channels[channel-1].Pixels) {
		return fmt.Errorf("Pixel index must be < %d", m.Channels[channel-1].Pixels)
	}
	m.Channels[channel-1].Pixels[index] = Pixel{r & 0xff, g & 0xff, b & 0xff}
	return nil
}

// Show sends the pixel buffer to the hardware.
func (m *Mote) Show() {
	var b []byte
	b = append(b, []byte("mote")...)
	b = append(b, []byte("o")...)
	for channel, data := range m.Channels {
		if data == nil {
			fmt.Printf("skipping empty channel %d\n", channel+1)
			continue
		}
		for _, pixel := range data.Pixels {
			b = append(b, byte(pixel.Blue))
			b = append(b, byte(pixel.Green))
			b = append(b, byte(pixel.Red))
		}
	}
	m.Port.Write(b)
}

// Clear the buffer of a specific channels, taking the following parameter.
//
//   - channel: Channel, either 1, 2, 3 or 4 corresponding to numbers on Mote
func (m *Mote) Clear(channel int) error {
	if channel > NumChannels || channel < 0 {
		return fmt.Errorf("channel index must be between 1 and 4")
	}
	if m.Channels[channel-1] == nil {
		return fmt.Errorf("please set up channel %d before using it", channel)
	}
	for pixel := 0; pixel < 16; pixel++ {
		err := m.SetPixel(channel, pixel, 0, 0, 0)
		if err != nil {
			return err
		}
	}
	return nil
}

// ClearAll clears the buffers of all configured channels.
func (m *Mote) ClearAll() error {
	for channel := 1; channel < NumChannels+1; channel++ {
		if m.Channels[channel-1] != nil {
			err := m.Clear(channel)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

// Close closes the serial port connected to the Mote device.
func (m *Mote) Close() {
	log.Printf("Closing port %v\n", m.PortName)
	m.Port.Close()
}

// GetPixelCount gets the number of pixels a channel is configured to use, taking the following parameters.
//   - channel: Channel, either 1, 2 3 or 4 corresponding to numbers on Mote
func (m *Mote) GetPixelCount(channel int) (int, error) {
	if channel > NumChannels || channel < 1 {
		return 0, fmt.Errorf("channel index must be between 1 and 4")
	}
	if m.Channels[channel-1] == nil {
		return 0, fmt.Errorf("channel %d has not been set up", channel)
	}
	return len(m.Channels[channel-1].Pixels), nil
}

// GetPixel gets the RGB colour of a single pixel, on a single channel, taking the following parameters.
//   - channel: Channel, either 1, 2, 3 or 4 corresponding to numbers on Mote
//   - index: Index of pixel to set, from 0 up
func (m *Mote) GetPixel(channel, index int) (Pixel, error) {
	if channel > NumChannels || channel < 1 {
		return Pixel{}, fmt.Errorf("channel index must be between 1 and 4")
	}
	if m.Channels[channel-1] == nil {
		return Pixel{}, fmt.Errorf("channel %d has not been set up", channel)
	}
	if index >= len(m.Channels[channel-1].Pixels) {
		return Pixel{}, fmt.Errorf("pixel index must be < %d", len(m.Channels[channel-1].Pixels))
	}
	return m.Channels[channel-1].Pixels[index], nil
}

func findSerialPort(v, p string) *string {
	ports, _ := enumerator.GetDetailedPortsList()
	if len(ports) == 0 {
		return nil
	}
	for _, port := range ports {
		if strings.EqualFold(port.VID, v) && strings.EqualFold(port.PID, p) {
			log.Printf("found Mote connected to port: %s\n", port.Name)
			return &port.Name
		}
	}
	return nil
}
