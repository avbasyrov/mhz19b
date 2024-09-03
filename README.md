mhz19b
========
A Go package to allow you to communicate with MH-Z19B CO2 sensor.

MH-Z19B has UART and PWM output modes. This driver communicates in UART mode.

Sensor can be connected to a small single-board computers (like Orange Pi, Raspberry Pi and so on) or even to desktop PC. All that is needed are GPIO pins that support UART (and this should be supported by OS).

MH-Z19B documentation
-------
[mh-z19b-co2-ver1_0.pdf](mh-z19b-co2-ver1_0.pdf)

Usage
-------
```go
package main

import (
	"fmt"
	"log"

	"github.com/avbasyrov/mhz19b"
)

func main() {
	sensor := mhz19b.New("/dev/ttyS5") // Choose correct serial device, on Windows systems it will look like "COM45"

	const set2kDetectionRange = true // Sets detection range 2000ppm
	const disableABC = true          // Disables Automatic Baseline Correction

	err := sensor.Connect(set2kDetectionRange, disableABC)
	if err != nil {
		log.Fatal(err)
	}

	co2level, err := sensor.ReadCO2()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("CO2 level: %d ppm\n", co2level)
}
```

Methods
-------
```go
Connect(set2kDetectionRange, disableABC bool) error
ReadCO2() (uint16, error)
Set2kDetectionRange() error
DisableABC() error
```
