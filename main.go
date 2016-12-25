package main

import (
	"log"

	"github.com/xperimental/onion-weather/oled"
)

func main() {
	oled, err := oled.NewOled()
	if err != nil {
		log.Fatalf("Error opening OLED device: %s", err)
	}
	defer oled.Close()

	if err := oled.Init(); err != nil {
		log.Fatalf("Error initializing display: %s", err)
	}

	if err := oled.Write("Hello Golang!\n\n  - a MIPS32 CPU"); err != nil {
		log.Fatalf("Error writing text: %s", err)
	}
}
