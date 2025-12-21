package usb

import (
	"fmt"
	"os/exec"
	"strings"
)

func extractContainerName(fullPath string) string {
	parts := strings.Split(fullPath, "\\")
	if len(parts) == 0 {
		return ""
	}
	return parts[len(parts)-1]
}

func getKeyAndCertInfo(containerName string) string {
	cmdStr := fmt.Sprintf(
		"/opt/cprocsp/bin/amd64/csptest -keyset -cont %s -info | grep -A999 'Key pair info:'",
		containerName,
	)
	cmd := exec.Command("bash", "-c", cmdStr)

	output, err := cmd.Output()
	if err != nil {
		if len(output) == 0 {
			return "Информация о ключевом контейеере не найдена"
		}
	}

	lines := strings.Split(string(output), "\n")
	var resultLines []string

	for _, line := range lines {
		line = strings.TrimRight(line, "\r")

		if strings.HasPrefix(line, "Container version:") {
			break
		}

		resultLines = append(resultLines, line)
	}

	return strings.Join(resultLines, "\n")
}

func FindCSPContainersByUUID(uuid string) ([]string, error) {
	cmdStr := "/opt/cprocsp/bin/amd64/csptest -keyset -enum_containers -fqcn"
	cmd := exec.Command("bash", "-c", cmdStr)

	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("ошибка CSP: %w, вывод: %s", err, output)
	}

	var containers []string

	uuidUpper := strings.ToUpper(uuid)

	lines := strings.Split(string(output), "\n")

	for _, line := range lines {
		line = strings.TrimSpace(line)

		if strings.Contains(strings.ToUpper(line), uuidUpper) {
			containerName := extractContainerName(line)
			if containerName != "" {
				containers = append(containers, containerName)
			}
		}
	}

	return containers, nil
}
