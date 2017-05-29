package main

import (
	"fmt"
	"strings"
	"unicode/utf8"

	serial "github.com/facchinm/go-serial-native"
)

// serial "github.com/mikepb/go-serial"

func main() {
	underline("Mote Scratch", "=")

	underline("ListPorts", "-")
	ports, _ := serial.ListPorts()

	fmt.Printf(" - found %d ports\n", len(ports))
	for _, info := range ports {
		fmt.Printf("   - Name: %v\n", info.Name())
		if vid, pid, err := info.USBVIDPID(); err == nil {
			fmt.Printf("     VID:  %v\n", vid)
			fmt.Printf("     PID:  %v\n", pid)
		}
	}

}

func underline(s, u string) {
	fmt.Println(s)
	underline := strings.Repeat(u, utf8.RuneCountInString(s))
	fmt.Println(underline)
}
