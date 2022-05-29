package main

import (
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"time"

	"github.com/whoisnian/glb/ansi"
	"github.com/whoisnian/glb/logger"
	"github.com/whoisnian/glb/serial"
	"github.com/whoisnian/observer/driver"
	"github.com/whoisnian/observer/server"
	"golang.org/x/sys/unix"
	"golang.org/x/term"
)

var isDebug = flag.Bool("d", false, "Output debug message")
var mode = flag.String("m", "cli", "Running mode (cli,test,server)")
var listenAt = flag.String("l", "127.0.0.1:8080", "Listen address (for server mode)")
var upStream = flag.String("u", "127.0.0.1:8081", "Upstream µStreamer server (for server mode)")
var device = flag.String("dev", "/dev/ttyUSB0", "UART device name")
var encode = flag.String("enc", "ch9329", "Encode driver (ch9329, kcom3)")

func main() {
	flag.Parse()
	logger.SetDebug(*isDebug)

	if *mode == "server" {
		if err := runServerMode(); err != nil {
			panic(err)
		}
		return
	}

	fd := int(os.Stdin.Fd())
	oldState, err := term.MakeRaw(fd)
	if err != nil {
		panic(err)
	}
	defer term.Restore(fd, oldState)

	if *mode == "test" {
		err = runTestMode(fd)
	} else {
		err = runCliMode(fd)
	}
	if err != nil {
		panic(err)
	}
}

func runTestMode(fd int) error {
	var buf [8]byte
	code := driver.EmptyKeycodes
	isCombo := false
	isExit := false
	for {
		n, err := unix.Read(fd, buf[:])
		fmt.Printf("ori: %s%x%s\r\n", ansi.Blue, buf[:n], ansi.Reset)
		code, isCombo, isExit = driver.DecodeFromCli(buf[:n], isCombo)
		if isCombo {
			fmt.Printf("res: %sComboMode%s\r\n", ansi.Green, ansi.Reset)
		} else {
			fmt.Printf("res: %s%s%s\r\n", ansi.Green, code, ansi.Reset)
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

func runCliMode(fd int) error {
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
		if logger.IsDebug() {
			fmt.Printf("ori: %s%x%s\r\n", ansi.Blue, buf[:n], ansi.Reset)
		}
		code, isCombo, isExit = driver.DecodeFromCli(buf[:n], isCombo)
		if isCombo {
			if logger.IsDebug() {
				fmt.Printf("res: %sComboMode%s\r\n", ansi.Green, ansi.Reset)
			}
		} else {
			if logger.IsDebug() {
				fmt.Printf("res: %s%s%s\r\n", ansi.Green, code, ansi.Reset)
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

func runServerMode() error {
	port, encodeFunc, err := openPort()
	if err != nil {
		return err
	}
	defer port.Close()

	return server.Start(*listenAt, *upStream, port, encodeFunc)
}
