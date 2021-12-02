package main

import (
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

func initTerminal() (fd int, oldState *term.State) {
	fd = int(os.Stdin.Fd())
	oldState, err := term.MakeRaw(fd)
	if err != nil {
		panic(err)
	}
	return fd, oldState
}

func runTestMode(fd int) {
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
			panic(err)
		}
	}
}

func runNormalMode(fd int) {
	port, err := serial.Open(*device, 9600, 8, serial.ParityNone, serial.StopBits1)
	if err != nil {
		panic(err)
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
			} else if res := driver.EncodeForCH9329(code); len(res) > 0 {
				if *isDebug {
					fmt.Printf("res: %s\r\n", code)
				}
				port.Write(res)
				port.Write(driver.EncodeForCH9329(driver.EmptyKeycodes))
			}
		} else if code == driver.ComboKeycodes {
			comboMode = true
		} else if res := driver.EncodeForCH9329(code); len(res) > 0 {
			if *isDebug {
				fmt.Printf("res: %s\r\n", code)
			}
			port.Write(res)
			port.Write(driver.EncodeForCH9329(driver.EmptyKeycodes))
		}

		if err == io.EOF {
			break
		} else if err != nil {
			panic(err)
		}
	}
}

func main() {
	flag.Parse()

	fd, oldState := initTerminal()
	defer term.Restore(fd, oldState)

	if *isTest {
		runTestMode(fd)
	} else {
		runNormalMode(fd)
	}
}
