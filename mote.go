package mote

import (
	"fmt"

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

// VID TODO
const VID = 5840

// PID TODO
const PID = 2244

// NAME TODO
const NAME = "Mote USB Dock"

// MAX_PIXELS TODO
const MAX_PIXELS = 512

// MAX_PIXELS_PER_CHANNEL TODO
const MAX_PIXELS_PER_CHANNEL = int(MAX_PIXELS / 4)

// Mote TODO
type Mote struct {
	PortName string
	Port     *serial.Port
	channels [4]*Channel
}

// Pixel TODO
type Pixel struct {
	Red, Green, Blue int
}

// Channel TODO
type Channel struct {
	Pixels []Pixel
	Flags  ChannelFlags
}

// ChannelFlags TODO
type ChannelFlags struct {
	GammaCorrection bool
}

func findSerialPort(v, p int, n string) *string {
	ports, _ := serial.ListPorts()
	//log.Printf("found %d ports\n", len(ports))
	if len(ports) == 0 {
		return nil
	}
	for _, info := range ports {
		name := info.Name()
		//log.Printf("port name: %v\n", name)
		if vid, pid, err := info.USBVIDPID(); err == nil {
			//log.Printf("VID: %v, PID: %v\n", vid, pid)
			if vid == v && pid == p {
				return &name
			}
		}
	}
	return nil
}

// NewMote TODO
func NewMote(portName string) *Mote {
	mote := Mote{
		PortName: portName,
	}
	if mote.PortName == "" {
		mote.PortName = *findSerialPort(VID, PID, NAME)
	}
	if mote.PortName == "" {
		//log.Fatal("unable to find Mote device")
	}

	options := serial.RawOptions
	options.Mode = serial.MODE_WRITE
	options.BitRate = 115200
	options.DataBits = 8
	options.Parity = serial.PARITY_NONE
	options.StopBits = 1
	options.FlowControl = 0
	//log.Printf("Opening port %v\n", mote.PortName)
	p, err := options.Open(mote.PortName)
	if err != nil {
		panic(err)
	}
	mote.Port = p

	return &mote
}

// ConfigureChannel configures a channel and takes the following parameters.
//
// - channel: Channel, either 1, 2, 3 or 4 corresponding to numbers on Mote
//
// - numPixels: Number of pixels to configure for this channel
//
// - gammaCorrection: Whether to enable gamma correction (default False)
func (m *Mote) ConfigureChannel(channel, numPixels int, gammaCorrection bool) error {
	if channel > 4 || channel < 1 {
		return fmt.Errorf("Channel index must be between 1 and 4")
	}
	if numPixels > MAX_PIXELS_PER_CHANNEL {
		return fmt.Errorf("Number of pixels can not be more than %d", MAX_PIXELS_PER_CHANNEL)
	}

	p := []Pixel{}
	for i := 0; i < numPixels; i++ {
		p = append(p, Pixel{0, 0, 0})
	}
	c := Channel{
		Pixels: p,
		Flags:  ChannelFlags{GammaCorrection: gammaCorrection},
	}
	m.channels[channel-1] = &c

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

// SetPixel Set the RGB colour of a single pixel, on a single channel
// - param channel: Channel, either 1, 2, 3 or 4 corresponding to numbers on Mote
// - param index: Index of pixel to set, from 0 up
// - param r: Amount of red: 0-255
// - param g: Amount of green: 0-255
// - param b: Amount of blue: 0-255
func (m *Mote) SetPixel(channel, index, r, g, b int) error {
	if channel > 4 || channel < 1 {
		return fmt.Errorf("channel index must be between 1 and 4")
	}
	if m.channels[channel-1] == nil {
		return fmt.Errorf("please set up channel %d before using it", channel)
	}
	if index >= len(m.channels[channel-1].Pixels) {
		return fmt.Errorf("Pixel index must be < %d", m.channels[channel-1].Pixels)
	}
	m.channels[channel-1].Pixels[index] = Pixel{r & 0xff, g & 0xff, b & 0xff}
	return nil
}

// Show send the pixel buffer to the hardware
func (m *Mote) Show() {
	var b []byte
	b = append(b, []byte("mote")...)
	b = append(b, []byte("o")...)
	for channel, data := range m.channels {
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

// Close TODO
func (m *Mote) Close() {
	//log.Printf("Closing port %v\n", m.PortName)
	m.Port.Close()
}
