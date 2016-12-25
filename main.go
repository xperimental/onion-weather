package main

import (
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	netatmo "github.com/exzz/netatmo-api-go"
	"github.com/xperimental/onion-weather/oled"
)

type textDisplay interface {
	Clear() error
	Write(text string) error
}

func main() {
	cfg, err := parseConfig()
	if err != nil {
		log.Fatalf("Error in configuration: %s", err)
	}

	log.Printf("Username: %s", cfg.Netatmo.Username)
	log.Printf("Update interval: %s", cfg.UpdateInterval)

	client, err := netatmo.NewClient(cfg.Netatmo)
	if err != nil {
		log.Fatalf("Error creating client: %s", err)
	}

	var display textDisplay
	if cfg.Dummy {
		display = NewDummyDisplay(6, 21, os.Stdout)
	} else {
		display, err = initOledDisplay()
		if err != nil {
			log.Fatalf("Error initializing display: %s", err)
		}
	}

	log.Println("Starting update loop...")
	for {
		if err := updateDisplay(client, display); err != nil {
			log.Printf("Error updating display: %s", err)
		}
		<-time.After(cfg.UpdateInterval)
	}
}

func initOledDisplay() (textDisplay, error) {
	oled, err := oled.NewOled()
	if err != nil {
		return nil, err
	}

	if err := oled.Init(); err != nil {
		return nil, err
	}

	return oled, nil
}

func updateDisplay(client *netatmo.Client, display textDisplay) error {
	dc, err := client.Read()
	if err != nil {
		return err
	}

	display.Clear()
	display.Write("--- N e t A t m o ---")
	for _, device := range dc.Devices() {
		displayDevice(display, *device)

		for _, module := range device.LinkedModules {
			displayModule(display, *module)
		}
	}
	return nil
}

func displayDevice(display textDisplay, device netatmo.Device) {
	writeLine(display, fmt.Sprintf("%s @ %s", device.ModuleName, dateString(device.DashboardData)))
	writeLine(display, fmt.Sprintf("  %.1f C - %d %%", *device.DashboardData.Temperature, *device.DashboardData.Humidity))
	writeLine(display, fmt.Sprintf("  %d ppm", *device.DashboardData.CO2))
}

func displayModule(display textDisplay, module netatmo.Device) {
	writeLine(display, fmt.Sprintf("%s @ %s", module.ModuleName, dateString(module.DashboardData)))
	writeLine(display, fmt.Sprintf("  %.1f C - %d %%", *module.DashboardData.Temperature, *module.DashboardData.Humidity))
}

func dateString(data netatmo.DashboardData) string {
	date := time.Unix(*data.LastMesure, 0)
	if time.Since(date) > 2*time.Hour {
		return "inactive"
	}

	return date.Format("15:04:05")
}

func writeLine(display textDisplay, line string) {
	if len(line) < 21 {
		fill := strings.Repeat(" ", 21-len(line))
		line += fill
	}

	if err := display.Write(line); err != nil {
		log.Printf("Error during output: %s", err)
	}
}
