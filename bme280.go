package main

import (
	"errors"
	"fmt"
	"os"
	"syscall"
)

func showPressTempHum() error {
	// for I2C
	// NOTE: for bme280
	file, err := os.OpenFile(
		"/dev/i2c-1",
		os.O_RDWR,
		os.ModeDevice,
	)
	if err != nil {
		return err
	}
	defer file.Close()

	r, _, errno := syscall.Syscall(
		syscall.SYS_IOCTL,
		uintptr(file.Fd()),
		uintptr(0x0703),
		uintptr(0x76),
	)
	if errno != 0 {
		return fmt.Errorf("failed to syscall.Syscall: %w", errno)
	}
	if r != 0 {
		return errors.New("unknown error from SYS_IOCTL")
	}

	_, err = file.Write([]byte{0xF7})
	if err != nil {
		return fmt.Errorf("faile to write reg: %w", err)
	}

	buf := make([]byte, 8)
	_, err = file.Read(buf)
	if err != nil {
		return fmt.Errorf("faile to read reg: %w", err)
	}

	// TODO: 補償計算
	press := int32(buf[0])<<12 | int32(buf[1])<<4 | int32(buf[2])>>4
	temp := int32(buf[3])<<12 | int32(buf[4])<<4 | int32(buf[5])>>4
	hum := int32(buf[6])<<8 | int32(buf[7])
	fmt.Println("PRESS:", press)
	fmt.Println("TEMP:", temp)
	fmt.Println("HUM:", hum)

	return nil
}
