package main

import (
	"io"
	"os"

	"github.com/whoisnian/observer/driver"
	"github.com/whoisnian/observer/serial"
	"golang.org/x/sys/unix"
	"golang.org/x/term"
)

const device = "/dev/ttyUSB0"

func main() {
	port, err := serial.Open(device, 9600, 8, serial.ParityNone, serial.StopBits1)
	if err != nil {
		panic(err)
	}
	defer port.Close()

	// wg := new(sync.WaitGroup)
	// wg.Add(1)
	// go func() {
	// 	defer wg.Done()
	// 	var buf [128]byte
	// 	for {
	// 		n, err := port.Read(buf[:])
	// 		fmt.Println("read ", n, err)
	// 		if err == io.EOF {
	// 			break
	// 		} else if err != nil {
	// 			panic(err)
	// 		}
	// 		fmt.Println(string(buf[:n]))
	// 	}
	// }()

	fd := int(os.Stdin.Fd())
	oldState, err := term.MakeRaw(fd)
	if err != nil {
		panic(err)
	}
	defer term.Restore(fd, oldState)

	var buf [8]byte
	for {
		n, err := unix.Read(fd, buf[:])
		if buf[0] == 1 { // input Ctrl-A to exit
			break
		}
		if res := driver.EncodeForCH9329_M(buf[:n]); len(res) > 0 {
			port.Write(res)
		}
		if err == io.EOF {
			break
		} else if err != nil {
			panic(err)
		}
	}

	// wg.Wait()
}
