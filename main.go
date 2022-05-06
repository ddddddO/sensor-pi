package main

import (
	"fmt"
	"log"
	"os"
)

func main() {
	log.Println("start")

	if err := showCO2(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// for I2C
	// NOTE: for bme280
	// i2cFile, err := os.OpenFile(
	// 	"/dev/i2c-1",
	// 	os.O_RDWR,
	// 	os.ModeDevice,
	// )
	// if err != nil {
	// 	fmt.Println(err)
	// 	return
	// }
	// defer i2cFile.Close()

	// r2, _, errno2 := syscall.Syscall(
	// 	syscall.SYS_IOCTL,
	// 	uintptr(i2cFile.Fd()),
	// 	uintptr(0x0703),
	// 	uintptr(0x76),
	// )
	// if errno2 != 0 {
	// 	err := fmt.Errorf("failed to syscall.Syscall: %w", errno)
	// 	fmt.Println(err)
	// 	return
	// }
	// if r2 != 0 {
	// 	err := errors.New("unknown error from SYS_IOCTL")
	// 	fmt.Println(err)
	// 	return
	// }

}
