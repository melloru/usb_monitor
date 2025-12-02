package usb

import (
	"encoding/json"
	"fmt"
	"os/exec"
	"strconv"
	"strings"
)

func sumChildrenFsused(children []interface{}) string {
	var total int64 = 0

	for _, child := range children {
		childMap := child.(map[string]interface{})

		fsusedRaw, exists := childMap["fsused"]
		if !exists || fsusedRaw == nil {
			continue
		}

		fsusedStr := fsusedRaw.(string)
		if fsusedStr == "" {
			continue
		}

		parsed, err := strconv.ParseInt(fsusedStr, 10, 64)
		if err != nil {
			continue
		}

		total += parsed
	}

	return strconv.FormatInt(total, 10)
}

func LsblkInfo(devnode string, fields []string) (map[string]string, error) {
	fieldStr := strings.Join(fields, ",")
	out, err := exec.Command("lsblk", "-J", "-b", "-o", fieldStr, devnode).Output()
	if err != nil {
		return nil, err
	}

	var parsed struct {
		Blockdevices []map[string]interface{} `json:"blockdevices"`
	}
	err = json.Unmarshal(out, &parsed)
	if err != nil {
		return nil, err
	}

	if len(parsed.Blockdevices) == 0 {
		return nil, fmt.Errorf("блочные устройства не найдены")
	}

	dev := parsed.Blockdevices[0]
	result := make(map[string]string)

	for _, f := range fields {
		key := strings.ToLower(f)

		if key == "fsused" {
			if children, exists := dev["children"]; exists && children != nil {
				childrenSlice := children.([]interface{})
				if len(childrenSlice) > 0 {
					result[f] = sumChildrenFsused(childrenSlice)
					continue
				}
			}
		}

		val, ok := dev[key]
		if !ok || val == nil {
			result[f] = ""
			continue
		}
		result[f] = fmt.Sprintf("%v", val)
	}

	return result, nil
}
