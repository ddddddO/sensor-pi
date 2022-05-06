package main

import (
	"errors"
	"fmt"
	"log"
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

func main() {
	log.Println("start")

	// for UART
	file, err := os.OpenFile(
		"/dev/serial0",
		syscall.O_RDWR|syscall.O_NOCTTY|syscall.O_NONBLOCK,
		0600,
	)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer file.Close()

	c := [NCCS]cc_t{}
	c[syscall.VTIME] = cc_t(0)
	c[syscall.VMIN] = cc_t()

	_ = syscall.SetNonblock(int(file.Fd()), false)
	r, _, errno := syscall.Syscall(
		syscall.SYS_IOCTL,
		uintptr(file.Fd()),
		uintptr(0x402C542B),
		uintptr(unsafe.Pointer(t)),
	)
	if errno != 0 {
		err := fmt.Errorf("failed to syscall.Syscall: %w", errno)
		fmt.Println(err)
		return
	}
	if r != 0 {
		err := errors.New("unknown error from SYS_IOCTL")
		fmt.Println(err)
		return
	}

}
