package mote

import (
	"bytes"
	"fmt"
	"testing"
)

var configureChannelTests = []struct {
	title           string
	channel         int
	numPixels       int
	gammaCorrection bool
	expected        []byte
	expectClose     bool
	expectError     bool
}{
	{
		title:           "valid channel all pixels and no gamma",
		channel:         1,
		numPixels:       16,
		gammaCorrection: false,
		expected:        []byte{'m', 'o', 't', 'e', 'c', 1, 16, 0},
		expectError:     false,
	},
	{
		title:           "valid channel partial pixels and gamma",
		channel:         1,
		numPixels:       8,
		gammaCorrection: true,
		expected:        []byte{'m', 'o', 't', 'e', 'c', 1, 8, 1},
		expectError:     false,
	},
	{
		title:           "invalid channel less than",
		channel:         0,
		numPixels:       8,
		gammaCorrection: false,
		expected:        []byte{},
		expectError:     true,
	},
	{
		title:           "invalid channel greater than",
		channel:         12,
		numPixels:       8,
		gammaCorrection: false,
		expected:        []byte{},
		expectError:     true,
	},
	{
		title:           "invalid numPixels less than",
		channel:         1,
		numPixels:       -1,
		gammaCorrection: false,
		expected:        []byte{},
		expectError:     true,
	},
	{
		title:           "invalid numPixels greater than",
		channel:         1,
		numPixels:       1000,
		gammaCorrection: false,
		expected:        []byte{},
		expectError:     true,
	},
}

type MockPort struct {
	expected    []byte
	expectClose bool
}

func (mp MockPort) Write(p []byte) (n int, err error) {
	if !bytes.Equal(p, mp.expected) {
		return 0, fmt.Errorf("written data [%v] does not match expected data [%v]", p, mp.expected)
	}
	return len(p), nil
}

func (mp MockPort) Close() error {
	if mp.expectClose {
		return nil
	}
	return fmt.Errorf("unexpected close invocation")
}

func Test_ConfigureChannel(t *testing.T) {
	for _, test := range configureChannelTests {
		t.Run(test.title, func(t *testing.T) {
			m := Mote{
				PortName: "/a/dummy/portname",
				Port: MockPort{
					expected:    test.expected,
					expectClose: test.expectClose,
				},
			}
			err := m.ConfigureChannel(test.channel, test.numPixels, test.gammaCorrection)
			if err != nil {
				t.Logf("[%s] error thrown: %v\n", test.title, err)
				if !test.expectError {
					t.Errorf("[%s] error was unexpected: %v", test.title, err)
				}
			}
			t.Logf("[%s] no error thrown", test.title)
		})
	}
}
