// @Project -> File    : bare-disk-perform -> target
// @IDE    : GoLand
// @Author    : wuji
// @Date   : 2023/8/23 10:30

package module

import (
	"fmt"
	"log"
	"os/exec"
	"strings"
)

func GetAutoScanDisks() ([]string, error) {
	var disks []string

	// 获取所有盘符
	args := []string{
		"-d",
		"-o",
		"NAME",
	}
	cmd := exec.Command("lsblk", args...)
	output, err := cmd.Output()
	if err != nil {
		fmt.Println("get disks failed")
		log.Fatal(err)
	}
	outputLineStr := strings.Split(strings.TrimSpace(string(output)), "\n")
	for _, line := range outputLineStr {
		// 过滤标题行
		if line == "NAME" {
			continue
		}

		// 获取盘符名称
		diskName := strings.TrimSpace(line)
		// 判断是否是机械硬盘
		smarctlCmd := exec.Command("smartctl", "-i", fmt.Sprintf("/dev/%s", diskName))
		smartCtlOutput, err := smarctlCmd.Output()
		if err != nil {
			log.Printf("执行 smartctl 命令失败: %v", err)
			continue
		}
		if strings.Contains(string(smartCtlOutput), "Rotation Rate") && !strings.Contains(string(smartCtlOutput), "Solid State Device") {
			disks = append(disks, diskName)
		}
	}

	fmt.Println("get disk list: ", disks)
	return disks, nil
}
