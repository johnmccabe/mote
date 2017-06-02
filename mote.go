/*
Package mote represents a connected Pimoroni Mote device, communicating over USB serial.

It allows you to configure the 4 channels and set individual pixels.

It is a port of the Pimoroni Mote Python library (https://github.com/pimoroni/mote).

The Mote device can be obtained directly from Pimoroni (https://shop.pimoroni.com/products/mote).
*/
package mote

import (
	"fmt"
	"log"

	serial "github.com/johnmccabe/go-serial-native"
)

// __init__
// find_serial_port
// configure_channel
// get_pixel_count
// get_pixel
// set_pixel
// clear
// show

// VID is the USB Vendor ID of the Pimoroni Mote
const VID = 5840

// PID is the USB Product ID of the Pimoroni Mote
const PID = 2244

// ProductName is the USB Product Name of the Pimoroni Mote
const ProductName = "Mote USB Dock"

// MaxPixels is the maximum addressable number of pixels across all channels
const MaxPixels = 512

// MaxPixelsPerChannel is the maximum addressable number of pixels across a single channel
const MaxPixelsPerChannel = int(MaxPixels / 4)

// Mote represents a connected Pimoroni Mote device
type Mote struct {
	PortName string
	Port     *serial.Port
	Channels [4]*Channel
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
		mote.PortName = *findSerialPort(VID, PID, ProductName)
	}
	if mote.PortName == "" {
		log.Fatal("unable to detect connected Mote")
	}

	options := serial.RawOptions
	options.Mode = serial.MODE_WRITE
	options.BitRate = 115200
	options.DataBits = 8
	options.Parity = serial.PARITY_NONE
	options.StopBits = 1
	options.FlowControl = 0
	p, err := options.Open(mote.PortName)
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
	if channel > 4 || channel < 1 {
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

	m.Port.Write(b)
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
	if channel > 4 || channel < 1 {
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
	if channel > 4 || channel < 0 {
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
	for channel := 1; channel < 5; channel++ {
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

func findSerialPort(v, p int, n string) *string {
	ports, _ := serial.ListPorts()
	if len(ports) == 0 {
		return nil
	}
	for _, info := range ports {
		portName := info.Name()
		if vid, pid, err := info.USBVIDPID(); err == nil {
			if vid == v && pid == p {
				log.Printf("found Mote connected to port: %s\n", portName)
				return &portName
			}
		}
		product := info.USBProduct()
		if product == ProductName {
			log.Printf("found Mote connected to port: %s\n", portName)
			return &portName
		}
	}
	return nil
}
