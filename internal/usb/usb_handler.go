package usb

import (
	"fmt"
	"strconv"
	"time"

	"github.com/jochenvg/go-udev"
)

func HandleUSBDevice(dev *udev.Device) {
	if dev.Action() != "add" {
		return
	}

	fmt.Println("\n=== USB EVENT ===")
	fmt.Println("Device:", dev.Devnode())
	fmt.Println("Model:", dev.PropertyValue("ID_MODEL"))
	fmt.Println("Vendor:", dev.PropertyValue("ID_VENDOR"))
	fmt.Println("Serial:", dev.PropertyValue("ID_SERIAL_SHORT"))

	fmt.Println("=== LSBLK INFO ===")
	for attempt := 1; attempt <= 10; attempt++ {
		fmt.Printf("Попытка %v...\n", attempt)

		fields := []string{"NAME", "SIZE", "FSUSED", "MOUNTPOINT"}
		info, err := LsblkInfo(dev.Devnode(), fields)
		if err != nil {
			fmt.Println("Ошибка:", err)
			time.Sleep(500 * time.Millisecond)
			continue
		}

		if info["FSUSED"] != "" && info["FSUSED"] != "0" {
			size, _ := strconv.ParseFloat(info["SIZE"], 64)
			used, _ := strconv.ParseFloat(info["FSUSED"], 64)

			sizeGB := ParseBytesToGB(size)
			usedGB := ParseBytesToGB(used)

			fmt.Printf("Размер: %.2f GB\n", sizeGB)
			fmt.Printf("Занято: %.2f GB\n", usedGB)
			return
		}

		time.Sleep(500 * time.Millisecond)
	}

	fmt.Println("Не удалось получить данные об использовании (устройство не смонтировано?)")
}
