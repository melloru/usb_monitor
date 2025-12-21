package usb

import (
	"fmt"
	"strconv"
)

const bytesInGB = 1024.0 * 1024.0 * 1024.0

func FormatBytesToGB(bytesStr string) string {
	if bytesStr == "" {
		return "0"
	}

	bytes, err := strconv.ParseFloat(bytesStr, 64)

	if err != nil {
		return "0"
	}

	gb := bytes / bytesInGB

	return fmt.Sprintf("%.0f GB", gb)
}
