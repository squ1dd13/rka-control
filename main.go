package main

import (
	"fmt"
	"os"
)

func update(mouse *KoneAIMO) {
	err := mouse.Update()
	if err != nil {
		panic(err)
	}
}

func main() {
	println("rka-control by squ1dd13")
	println("This is /experimental/ software. I am NOT responsible for any damage caused to your hardware through use of this tool.\n")

	if len(os.Args) < 2 {
		println("Usage:\n  rka-control /path/to/lights.yml\n")
		os.Exit(1)
	}

	// Product ID of the Roccat Kone AIMO mouse.
	devices := FindAll(0x2e27)

	if len(devices) == 0 {
		fmt.Println("No devices found.")
		return
	}

	deviceIndex := 0

	// On Linux (and maybe Windows, IDK) you get two devices. We want the one
	//  for interface 1. Trying to configure interface 0 does work, but makes
	//  the mouse unresponsive until you disconnect it and plug it back in.
	if len(devices) > 1 {
		for i, dev := range devices {
			if dev.info.Interface == 1 {
				deviceIndex = i
				break
			}
		}
	}

	mouse, err := NewKoneAIMO(devices[deviceIndex])
	if err != nil {
		panic(err)
	}

	mouse.Visual = loadLights(os.Args[1])
	update(&mouse)

	closeErr := mouse.Close()

	if closeErr != nil {
		panic(closeErr)
	}
}
