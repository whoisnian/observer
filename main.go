package main

import (
	"errors"
	"fmt"
	"io"
	"os"
	"time"

	"github.com/whoisnian/glb/ansi"
	"github.com/whoisnian/glb/config"
	"github.com/whoisnian/glb/logger"
	"github.com/whoisnian/observer/driver"
	"github.com/whoisnian/observer/serial"
	"github.com/whoisnian/observer/server"
	"golang.org/x/sys/unix"
	"golang.org/x/term"
)

var CFG struct {
	Debug      bool   `flag:"d,false,Enable debug output"`
	Mode       string `flag:"m,cli,Running mode (cli,test,server)"`
	ListenAddr string `flag:"l,127.0.0.1:8080,Listen address (for server mode)"`
	UpStream   string `flag:"u,127.0.0.1:8081,Upstream µStreamer server (for server mode)"`
	Device     string `flag:"dev,/dev/ttyUSB0,UART device name"`
	Encode     string `flag:"enc,ch9329,Encode driver (ch9329, kcom3)"`
	Baudrate   int    `flag:"baud,9600,UART device baudrate"`
}

func main() {
	err := config.FromCommandLine(&CFG)
	if err != nil {
		logger.Fatal(err)
	}
	logger.SetDebug(CFG.Debug)

	if CFG.Mode == "server" {
		if err := runServerMode(); err != nil {
			logger.Fatal(err)
		}
		return
	}

	fd := int(os.Stdin.Fd())
	oldState, err := term.MakeRaw(fd)
	if err != nil {
		logger.Fatal(err)
	}
	defer term.Restore(fd, oldState)

	if CFG.Mode == "test" {
		err = runTestMode(fd)
	} else {
		err = runCliMode(fd)
	}
	if err != nil {
		logger.Fatal(err)
	}
}

func runTestMode(fd int) error {
	var buf [8]byte
	var code driver.Keycodes
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
	port, err = serial.Open(CFG.Device, CFG.Baudrate, 8, serial.ParityNone, serial.StopBits1)
	if err != nil {
		return nil, nil, err
	}

	if CFG.Encode == "ch9329" {
		encodeFunc = driver.EncodeForCH9329
		port.SetInterval(time.Millisecond * 50)
	} else if CFG.Encode == "kcom3" {
		encodeFunc = driver.EncodeForKCOM3
		port.SetInterval(time.Millisecond * 50)
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
	var code driver.Keycodes
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

	return server.Start(CFG.ListenAddr, CFG.UpStream, port, encodeFunc)
}
