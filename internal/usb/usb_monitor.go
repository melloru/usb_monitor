package usb

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/jochenvg/go-udev"
)

func MonitorUSBInsertions() error {
	u := udev.Udev{}
	monitor := u.NewMonitorFromNetlink("udev")
	err := monitor.FilterAddMatchSubsystemDevtype("block", "disk")
	if err != nil {
		return fmt.Errorf("ошибка при установке фильтра udev: %w", err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	devCh, errCh, _ := monitor.DeviceChan(ctx)
	fmt.Println("Жду подключения флешки...")

	for {
		select {
		case dev := <-devCh:
			if dev == nil || strings.ToLower(dev.PropertyValue("ID_BUS")) != "usb" {
				continue
			}
			HandleUSBDevice(dev)

		case err := <-errCh:
			if err != nil {
				log.Println("udev error:", err)
			}
		}
	}
}
