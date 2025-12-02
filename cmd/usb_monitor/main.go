package main

import (
	"log"

	"usb_monitoring/internal/usb"
)

func main() {
	if err := usb.MonitorUSBInsertions(); err != nil {
		log.Fatal(err)
	}
}
