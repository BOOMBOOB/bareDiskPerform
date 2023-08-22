// @Project -> File    : bare-disk-perform -> smart
// @IDE    : GoLand
// @Author    : wuji
// @Date   : 2023/8/22 11:11

package module

import (
	"log"
	"os/exec"
	"regexp"
	"strings"
)

type DiskSmart struct {
	DeviceModel  string
	SerialNumber string
	UserCapacity string
	RotationRate string
	FormFactor   string
	SMART        string
}

func GetDiskSmartInfo(device string) (DiskSmart, error) {
	if !strings.HasPrefix(device, "/dev/") {
		device = "/dev/" + device
	}
	cmd := exec.Command("smartctl", "-a", device)

	output, err := cmd.Output()
	if err != nil {
		log.Fatal(err)
	}

	// 将输出结果转换为字符串
	result := string(output)

	smart := DiskSmart{}

	// 使用正则表达式提取 Device Model 和 Serial Number
	deviceModelRegex := regexp.MustCompile(`Device Model:\s+(.+)`)
	serialNumberRegex := regexp.MustCompile(`Serial Number:\s+(.+)`)
	userCapacityRegex := regexp.MustCompile(`User Capacity:\s+(.+)`)
	rotationRateRegex := regexp.MustCompile(`Rotation Rate:\s+(.+)`)
	formFactorRegex := regexp.MustCompile(`Form Factor:\s+(.+)`)

	deviceModelMatches := deviceModelRegex.FindStringSubmatch(result)
	serialNumberMatches := serialNumberRegex.FindStringSubmatch(result)
	userCapacityMatches := userCapacityRegex.FindStringSubmatch(result)
	rotationRateMatches := rotationRateRegex.FindStringSubmatch(result)
	formFactorMatches := formFactorRegex.FindStringSubmatch(result)

	if len(deviceModelMatches) > 1 {
		smart.DeviceModel = deviceModelMatches[1]
	}

	if len(serialNumberMatches) > 1 {
		smart.SerialNumber = serialNumberMatches[1]
	}

	if len(userCapacityMatches) > 1 {
		smart.UserCapacity = userCapacityMatches[1]
	}

	if len(rotationRateMatches) > 1 {
		smart.RotationRate = rotationRateMatches[1]
	}

	if len(formFactorMatches) > 1 {
		smart.FormFactor = formFactorMatches[1]
	}

	smart.SMART = result

	return smart, nil
}
