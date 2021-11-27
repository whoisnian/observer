// from: https://github.com/tarm/serial
package serial

import (
	"errors"
	"os"

	"golang.org/x/sys/unix"
)

type Port struct {
	f *os.File
}

// Example:
// device: /dev/ttyUSB0
// baudrate: 9600 (9600|19200|38400|57600|115200)
// databits: 8 (5|6|7|8)
// parity: ParityNone (ParityNone|ParityOdd|ParityEven|ParityMark|ParitySpace)
// stopbits: StopBits1 (StopBits1|StopBits2)
func Open(device string, baudrate int, databits int, parity uint32, stopbits uint32) (p *Port, err error) {
	baudrateValid, ok := baudrateMap[baudrate]
	if !ok {
		return nil, errors.New("invalid baudrate")
	}
	databitsValid, ok := databitsMap[databits]
	if !ok {
		return nil, errors.New("invalid databits")
	}

	f, err := os.OpenFile(device, unix.O_RDWR|unix.O_NOCTTY|unix.O_NONBLOCK, 0666)
	if err != nil {
		return nil, err
	}

	fd := f.Fd()
	state := unix.Termios{
		Iflag:  unix.IGNPAR,
		Cflag:  unix.CREAD | unix.CLOCAL | baudrateValid | databitsValid | parity | stopbits,
		Ispeed: baudrateValid,
		Ospeed: baudrateValid,
	}
	state.Cc[unix.VMIN] = 1
	state.Cc[unix.VTIME] = 0
	if err = unix.IoctlSetTermios(int(fd), unix.TCSETS, &state); err != nil {
		return nil, err
	}

	if err = unix.SetNonblock(int(fd), false); err != nil {
		return nil, err
	}

	return &Port{f}, nil
}

func (p *Port) Read(b []byte) (n int, err error) {
	return p.f.Read(b)
}

func (p *Port) Write(b []byte) (n int, err error) {
	return p.f.Write(b)
}

func (p *Port) Close() (err error) {
	return p.f.Close()
}
