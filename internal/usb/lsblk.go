package usb

import (
	"encoding/json"
	"fmt"
	"os/exec"
	"strconv"
	"strings"
)

func getChildrenInfo(children []interface{}, field string) string {
	for _, child := range children {
		childMap := child.(map[string]interface{})
		val, exists := childMap[strings.ToLower(field)]

		if exists && val != nil {
			strVal := fmt.Sprintf("%v", val)

			if strVal != "" && strVal != "null" {
				return strVal
			}
		}
	}
	return ""
}

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

	cmdStr := fmt.Sprintf("lsblk -J -b -o %s %s", fieldStr, devnode)
	cmd := exec.Command("bash", "-c", cmdStr)

	out, err := cmd.Output()
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

	children, hasChildren := dev["children"]
	var childrenSlice []interface{}
	if hasChildren && children != nil {
		childrenSlice = children.([]interface{})
	}

	for _, f := range fields {
		key := strings.ToLower(f)

		if key == "fsused" {
			if hasChildren && len(childrenSlice) > 0 {
				fsusedBytes := sumChildrenFsused(childrenSlice)
				fsusedGB := FormatBytesToGB(fsusedBytes)
				result[f] = fsusedGB
				continue
			}
			result[f] = "0 GB"
			continue
		}
		if key == "uuid" {
			if hasChildren && len(childrenSlice) > 0 {
				result[f] = getChildrenInfo(childrenSlice, f)
				continue
			}
			result[f] = ""
			continue
		}

		if key == "size" {
			if sizeVal, ok := dev["size"]; ok && sizeVal != nil {
				sizeStr := fmt.Sprintf("%v", sizeVal)
				sizeGB := FormatBytesToGB(sizeStr)
				result[f] = sizeGB
				continue
			}
			result[f] = "0 GB"
			continue
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
