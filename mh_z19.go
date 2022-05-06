package main

import (
	"errors"
	"fmt"
	"os"
	"syscall"
	"unsafe"
)

// for Termios
type tcflag_t uint
type speed_t uint
type cc_t byte

const NCCS = 19

type termios struct {
	c_iflag  tcflag_t
	c_oflag  tcflag_t
	c_cflag  tcflag_t
	c_lflag  tcflag_t
	c_cc     [NCCS]cc_t
	c_ispeed speed_t
	c_ospeed speed_t
}

func showCO2() error {
	// for UART
	file, err := os.OpenFile(
		"/dev/serial0",
		syscall.O_RDWR|syscall.O_NOCTTY|syscall.O_NONBLOCK,
		0600,
	)
	if err != nil {
		return err
	}
	defer file.Close()

	fmt.Println("debug 1")

	const MINIMAMREADSIZE = 4
	c := [NCCS]cc_t{}
	c[syscall.VTIME] = cc_t(0)
	c[syscall.VMIN] = cc_t(MINIMAMREADSIZE)
	t := &termios{
		c_cflag:  syscall.CLOCAL | syscall.CREAD | syscall.CS8,
		c_cc:     c,
		c_ispeed: speed_t(9600),
		c_ospeed: speed_t(9600),
	}

	_ = syscall.SetNonblock(int(file.Fd()), false)
	r, _, errno := syscall.Syscall(
		syscall.SYS_IOCTL,
		uintptr(file.Fd()),
		uintptr(0x402C542B),
		uintptr(unsafe.Pointer(t)),
	)
	if errno != 0 {
		err := fmt.Errorf("failed to syscall.Syscall: %w", errno)
		return err
	}
	if r != 0 {
		err := errors.New("unknown error from SYS_IOCTL")
		return err
	}

	fmt.Println("debug 2")

	var co2 int64
	writeN, err := file.Write([]byte{0xff, 0x01, 0x86, 0x00, 0x00, 0x00, 0x00, 0x00, 0x79})
	if err != nil {
		return err
	}

	fmt.Println("debug 3")

	for i := 0; i < writeN; i++ {
		buf := make([]byte, 1)
		_, err = file.Read(buf)
		if err != nil {
			return err
		}

		// TODO: checksum
		switch i {
		case 2: // high level concentration
			co2 += int64(buf[0]) * 256
		case 3: // low level concentration
			co2 += int64(buf[0])
		}
	}

	fmt.Printf("CO2: %d\n", co2)
	return nil
}
