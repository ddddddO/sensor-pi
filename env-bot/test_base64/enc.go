package main

import (
	"encoding/base64"
	"fmt"
	"io"
	"os"
)

func main() {
	f, err := os.Open("../plotter/pressure.png")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer f.Close()

	body, err := io.ReadAll(f)
	if err != nil {
		fmt.Println(err)
		return
	}

	encoded := base64.StdEncoding.EncodeToString(body)
	fmt.Print(encoded)
}
