package usb

const bytesInGB = 1024.0 * 1024.0 * 1024.0

func ParseBytesToGB(val float64) float64 {
	if val == 0 {
		return 0
	}

	return val / bytesInGB
}
