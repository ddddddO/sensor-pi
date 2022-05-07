package main

import (
	"fmt"
	"log"
	"os"
)

func main() {
	log.Println("start")

	if err := showCO2(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	if err := showPressTempHum(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
