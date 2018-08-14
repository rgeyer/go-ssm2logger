package main

import (
	"encoding/hex"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"time"

	"github.com/jacobsa/go-serial/serial"
)

func main() {
	options := serial.OpenOptions{
		PortName:              "/dev/tty.usbserial-1420",
		BaudRate:              4800,
		DataBits:              8,
		StopBits:              1,
		InterCharacterTimeout: 100,
		MinimumReadSize:       0,
		ParityMode:            serial.PARITY_NONE,
	}

	f, err := serial.Open(options)

	if err != nil {
		fmt.Println("Error opening serial port: ", err)
		os.Exit(-1)
	} else {
		defer f.Close()
	}

	initPacket := NewInitPacket(0xF0, 0x10)

	count, err := f.Write(initPacket.GetBytes())
	if err != nil {
		fmt.Println("Failed to send serial:", err)
		os.Exit(-1)
	}
	fmt.Println("Wrote bytes", count)

	time.Sleep(10 * time.Millisecond)

	resp, err := ioutil.ReadAll(f)
	if err != nil {
		if err != io.EOF {
			fmt.Println("Error reading from serial port: ", err)
		}
	} else {
		fmt.Println("Tx: ", hex.EncodeToString(resp[:count]))
		fmt.Println("Rx: ", hex.EncodeToString(resp[count:]))

		resp_bytes := resp[count:]
		fmt.Println(resp_bytes[8] & (1 << 6))
	}

	//[]byte{0x80, 0x10, 0xF0, 0x08, 0xA8, 0x00, 0x00, 0x00, 0x08, 0x00, 0x00, 0x1C, 0x54}

	readPid := NewReadPacket(0xF0, 0x10, []byte{0x46, 0x3c, 0x3D, 0x1C, 0x20, 0x22, 0x29, 0x32, 0x3B, 0xA, 0xD, 0x11, 0xE, 0x9, 0x13})
	count, err = f.Write(readPid.GetBytes())
	if err != nil {
		fmt.Println("Failed to send serial:", err)
		os.Exit(-1)
	}
	fmt.Println("Wrote bytes", count)

	time.Sleep(10 * time.Millisecond)

	resp, err = ioutil.ReadAll(f)
	if err != nil {
		if err != io.EOF {
			fmt.Println("Error reading from serial port: ", err)
		}
	} else {
		fmt.Println("Tx: ", hex.EncodeToString(resp[:count]))
		fmt.Println("Rx: ", hex.EncodeToString(resp[count:]))

		for _, val := range resp {
			fmt.Println(val)
		}
	}
}
