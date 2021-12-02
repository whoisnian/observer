package main

import (
	"errors"
	"flag"
	"fmt"
	"io"
	"os"

	"github.com/whoisnian/observer/driver"
	"github.com/whoisnian/observer/serial"
	"golang.org/x/sys/unix"
	"golang.org/x/term"
)

var isDebug = flag.Bool("d", false, "Output debug message")
var isTest = flag.Bool("t", false, "Keycodes test mode")
var device = flag.String("dev", "/dev/ttyUSB0", "UART device name")
var encode = flag.String("enc", "ch9329", "Encode driver (ch9329, kcom3)")

func initTerminal() (fd int, oldState *term.State, err error) {
	fd = int(os.Stdin.Fd())
	oldState, err = term.MakeRaw(fd)
	if err != nil {
		return 0, nil, err
	}
	return fd, oldState, nil
}

func runTestMode(fd int) error {
	var buf [8]byte
	comboMode := false
	for {
		n, err := unix.Read(fd, buf[:])
		fmt.Printf("ori: %x\r\n", buf[:n])
		code := driver.VT100Decode(buf[:n])
		if comboMode {
			comboMode = false
			code = driver.CalcCombo(code)
			fmt.Printf("res: %s\r\n", code)
			if code == driver.ComboKeycodesExit {
				break
			}
		} else {
			comboMode = code == driver.ComboKeycodes
			fmt.Printf("res: %s\r\n", code)
		}

		if err == io.EOF {
			break
		} else if err != nil {
			return err
		}
	}
	return nil
}

func openPort() (port *serial.Port, encodeFunc driver.EncodeFunc, err error) {
	if *encode == "ch9329" {
		port, err = serial.Open(*device, 9600, 8, serial.ParityNone, serial.StopBits1)
		encodeFunc = driver.EncodeForCH9329
	} else if *encode == "kcom3" {
		port, err = serial.Open(*device, 57600, 8, serial.ParityNone, serial.StopBits1)
		encodeFunc = driver.EncodeForKCOM3
	} else {
		return nil, nil, errors.New("unknown encode driver")
	}
	if err != nil {
		return nil, nil, err
	}
	return port, encodeFunc, nil
}

func runNormalMode(fd int) error {
	port, encodeFunc, err := openPort()
	if err != nil {
		return err
	}
	defer port.Close()

	var buf [8]byte
	comboMode := false
	for {
		n, err := unix.Read(fd, buf[:])
		if *isDebug {
			fmt.Printf("ori: %x\r\n", buf[:n])
		}
		code := driver.VT100Decode(buf[:n])
		if comboMode {
			comboMode = false
			code = driver.CalcCombo(code)
			if code == driver.ComboKeycodesExit {
				break
			} else if res := encodeFunc(code); len(res) > 0 {
				if *isDebug {
					fmt.Printf("res: %s\r\n", code)
				}
				port.Write(res)
				port.Write(encodeFunc(driver.EmptyKeycodes))
			}
		} else if code == driver.ComboKeycodes {
			comboMode = true
		} else if res := encodeFunc(code); len(res) > 0 {
			if *isDebug {
				fmt.Printf("res: %s\r\n", code)
			}
			port.Write(res)
			port.Write(encodeFunc(driver.EmptyKeycodes))
		}

		if err == io.EOF {
			break
		} else if err != nil {
			return err
		}
	}
	return nil
}

func main() {
	flag.Parse()

	fd, oldState, err := initTerminal()
	if err != nil {
		panic(err)
	}
	defer term.Restore(fd, oldState)

	if *isTest {
		err = runTestMode(fd)
	} else {
		err = runNormalMode(fd)
	}
	if err != nil {
		panic(err)
	}
}
