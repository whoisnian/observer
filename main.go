package main

import (
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"time"

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
	code := driver.EmptyKeycodes
	isCombo := false
	isExit := false
	for {
		n, err := unix.Read(fd, buf[:])
		fmt.Printf("ori: %x\r\n", buf[:n])
		code, isCombo, isExit = driver.Decode(buf[:n], isCombo)
		if isCombo {
			fmt.Printf("enter comboMode\r\n")
		} else {
			fmt.Printf("res: %s\r\n", code)
			if isExit {
				break
			}
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
	port, err = serial.Open(*device, 9600, 8, serial.ParityNone, serial.StopBits1)
	if err != nil {
		return nil, nil, err
	}

	if *encode == "ch9329" {
		encodeFunc = driver.EncodeForCH9329
	} else if *encode == "kcom3" {
		encodeFunc = driver.EncodeForKCOM3
		port.SetInterval(time.Millisecond * 16)
	} else {
		port.Close()
		return nil, nil, errors.New("unknown encode driver")
	}

	return port, encodeFunc, nil
}

func runNormalMode(fd int) error {
	port, encodeFunc, err := openPort()
	if err != nil {
		return err
	}
	defer port.Close()
	stop := port.GoWaitAndSend()

	var buf [8]byte
	code := driver.EmptyKeycodes
	isCombo := false
	isExit := false
	for {
		n, err := unix.Read(fd, buf[:])
		if *isDebug {
			fmt.Printf("ori: %x\r\n", buf[:n])
		}
		code, isCombo, isExit = driver.Decode(buf[:n], isCombo)
		if isCombo {
			if *isDebug {
				fmt.Printf("enter comboMode\r\n")
			}
		} else {
			if *isDebug {
				fmt.Printf("res: %s\r\n", code)
			}
			if isExit {
				break
			}
			if res := encodeFunc(code); len(res) > 0 {
				port.Push(res)
				port.Push(encodeFunc(driver.EmptyKeycodes))
			}
		}

		if err == io.EOF {
			break
		} else if err != nil {
			return err
		}
	}

	stop()
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
