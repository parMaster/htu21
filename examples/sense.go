package main

import (
	"log"

	"github.com/parMaster/htu21"
	"periph.io/x/conn/v3/i2c/i2creg"
	"periph.io/x/conn/v3/physic"
	"periph.io/x/host/v3"
)

var (
	htu21Data   physic.Env
	htu21Device *htu21.Dev
)

func main() {

	// Preparing to read sensor
	if _, err := host.Init(); err != nil {
		log.Fatal(err)
	}

	// Use i2creg I²C bus registry to find the first available I²C bus.
	b, err := i2creg.Open("4")
	if err != nil {
		log.Fatalf("failed to open I²C: %v", err)
	}
	defer b.Close()

	htu21Device, err = htu21.NewI2C(b, 0x40)
	if err != nil {
		log.Fatalf("failed to initialize htu21: %v", err)
	}

	if err := htu21Device.Sense(&htu21Data); err != nil {
		log.Fatal(err)
	}

	log.Println(htu21Data.Temperature)
	log.Println(htu21Data.Humidity)
}
