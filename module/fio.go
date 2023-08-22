// @Project -> File    : bare-disk-perform -> fio
// @IDE    : GoLand
// @Author    : wuji
// @Date   : 2023/8/22 10:56

package module

import (
	"fmt"
	"log"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
)

type Result struct {
	Iops      int
	BandWidth int
	ClatAvg   string
	Clat95    string
	Clat99    string
}

type WorkLoad struct {
	BlockSize string
	IODepth   int
	IOType    string
}

func ExecuteFio(device string, iotype string, iodepth int) ([]byte, WorkLoad, error) {
	if !strings.HasSuffix(device, "/dev/") {
		device = "/dev/" + device
	}
	var bs string
	if iotype == "read" || iotype == "write" {
		bs = "1m"
	} else if iotype == "randread" || iotype == "randwrite" {
		bs = "4k"
	} else {
		log.Fatal("unknown iotype")
	}

	workload := WorkLoad{}
	workload.BlockSize = bs
	workload.IODepth = iodepth
	workload.IOType = iotype

	args := []string{
		"-filename=" + device,
		"-rw=" + iotype,
		"-iodepth=" + strconv.Itoa(iodepth),
		"-bs=" + bs,
		"-ioengine=libaio",
		"-sync=1",
		"-direct=1",
		"-time_based",
		"-runtime=20",
		"-ramp_time=10",
		"-group_reporting",
		"-name=job",
	}
	fmt.Println("args: ", args)
	cmd := exec.Command("fio", args...)
	output, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Println("output: ", output)
		fmt.Println("err: ", err)
		//log.Fatal(err)
	}
	return output, workload, nil
}

func ParseFIOOutput(output []byte) (Result, error) {
	result := Result{}

	// 转换输出为字符串
	outputStr := string(output)

	// 提取 IOPS
	iopsPattern := `IOPS=\s*(\d+)`
	iopsRegex := regexp.MustCompile(iopsPattern)
	iopsMatches := iopsRegex.FindStringSubmatch(outputStr)
	if len(iopsMatches) >= 2 {
		iops, err := strconv.Atoi(iopsMatches[1])
		if err == nil {
			result.Iops = iops
		}
	}

	// 提取带宽
	bandwidthPattern := `BW=\s*(\d+)`
	bandwidthRegex := regexp.MustCompile(bandwidthPattern)
	bandwidthMatches := bandwidthRegex.FindStringSubmatch(outputStr)
	if len(bandwidthMatches) >= 2 {
		bandwidth, err := strconv.Atoi(bandwidthMatches[1])
		if err == nil {
			result.BandWidth = bandwidth
		}
	}

	// 提取 Clat 平均值单位
	clatAvgUnitPattern := `clat\s*\((\w+)\)`
	clatAvgUnitRegex := regexp.MustCompile(clatAvgUnitPattern)
	clatAvgUnitMatches := clatAvgUnitRegex.FindStringSubmatch(outputStr)

	// 提取 Clat 平均值
	clatAvgPattern := `clat\s+avg=\s*(\S+)`
	clatAvgRegex := regexp.MustCompile(clatAvgPattern)
	clatAvgMatches := clatAvgRegex.FindStringSubmatch(outputStr)
	if len(clatAvgMatches) > 1 && len(clatAvgUnitMatches) > 1 {
		result.ClatAvg = clatAvgMatches[1] + clatAvgUnitMatches[1]
	}

	// 提取 Clat percentiles 单位
	clatPercentilesUnitPattern := `clat percentiles \((.+)\)`
	clatPercentilesUnitRegex := regexp.MustCompile(clatPercentilesUnitPattern)
	clatPercentilesUnitMatches := clatPercentilesUnitRegex.FindStringSubmatch(outputStr)
	if len(clatPercentilesUnitMatches) > 1 {

		// 提取 Clat 95th 百分位
		clat95Pattern := `95\.00th=\[\s*(\d+)\]`
		clat95Regex := regexp.MustCompile(clat95Pattern)
		clat95Matches := clat95Regex.FindStringSubmatch(outputStr)
		if len(clat95Matches) > 1 {
			result.Clat95 = clat95Matches[1] + clatPercentilesUnitMatches[1]
		}

		// 提取 Clat 99th 百分位
		clat99Pattern := `99\.00th=\[\s*(\d+)\]`
		clat99Regex := regexp.MustCompile(clat99Pattern)
		clat99Matches := clat99Regex.FindStringSubmatch(outputStr)
		if len(clat99Matches) > 1 {
			result.Clat99 = clat99Matches[1] + clatPercentilesUnitMatches[1]
		}
	}

	return result, nil
}
