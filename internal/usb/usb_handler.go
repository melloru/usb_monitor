package usb

import (
	"encoding/json"
	"fmt"
	"os/exec"
	"time"

	"github.com/jochenvg/go-udev"
)

func PollForMountpoint(devnode string, attempts int) (string, error) {

	for i := 0; i < attempts; i++ {
		cmdStr := fmt.Sprintf("lsblk -J -o MOUNTPOINT %s", devnode)
		cmd := exec.Command("bash", "-c", cmdStr)

		output, err := cmd.Output()
		if err != nil {
			return "", err
		}

		var result struct {
			Blockdevices []struct {
				Mountpoint string `json:"mountpoint"`
			} `json:"blockdevices"`
		}

		if err := json.Unmarshal(output, &result); err != nil {
			return "", err
		}

		if len(result.Blockdevices) > 1 {
			for i := 1; i < len(result.Blockdevices); i++ {
				mp := result.Blockdevices[i].Mountpoint
				if mp != "" && mp != "null" {
					return mp, nil
				}
			}
		}

		time.Sleep(500 * time.Millisecond)
	}

	return "", fmt.Errorf("таймаут ожидания монтирования")
}

func HandleUSBDevice(dev *udev.Device) {
	fmt.Println("\n=== USB EVENT ===")
	fmt.Println("Device:", dev.Devnode())
	fmt.Println("Model:", dev.PropertyValue("ID_MODEL"))
	fmt.Println("Vendor:", dev.PropertyValue("ID_VENDOR"))
	fmt.Println("Serial:", dev.PropertyValue("ID_SERIAL_SHORT"))

	fmt.Println("=== LSBLK INFO ===")

	mountpoint, err := PollForMountpoint(dev.Devnode(), 10)

	if err != nil {
		fmt.Printf("Mount timeout: %v\n", err)
	} else {
		fmt.Printf("Mounted at: %s\n", mountpoint)
	}

	fields := []string{"NAME", "SIZE", "FSUSED", "UUID"}

	info, err := LsblkInfo(dev.Devnode(), fields)
	if err != nil {
		fmt.Println("Ошибка:", err)
	}

	for _, f := range fields {
		val := info[f]
		fmt.Printf("%s: %s\n", f, val)
	}

	if uuid := info["UUID"]; uuid != "" {
		fmt.Println("=== CSP CONTAINERS ===")

		containers, _ := FindCSPContainersByUUID(uuid)
		if len(containers) == 0 {
			fmt.Println("Нет контейнеров CSP на этой флешке")
		} else {
			fmt.Printf("Найдено контейнеров: %d\n", len(containers))
			for i, container := range containers {
				fmt.Printf("%d. %s\n", i+1, container)
			}
		}

		fmt.Println("=== CSP CONTAINERS INFO ===")
		if len(containers) > 0 {
			for _, container := range containers {
				containerInfo := getKeyAndCertInfo(container)
				fmt.Printf("Контейнер: %s\n%s", container, containerInfo)
			}
		}
	}

}
