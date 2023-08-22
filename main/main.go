// @Project -> File    : bare-disk-perform -> main
// @IDE    : GoLand
// @Author    : wuji
// @Date   : 2023/8/22 10:40

package main

import (
	"bare-disk-perform/module"
	"fmt"
	"log"
)

func main() {
	config, err := module.LoadConfig("./config.json")
	if err != nil {
		log.Fatal(err)
	}

	// 获取db 连接
	db, err := module.NewDatabase(config)
	if err != nil {
		log.Fatal(err)
	}

	defer func() {
		if err := db.Close(); err != nil {
			fmt.Println("failed to close database connection: ", err)
		}
	}()

	disks := config.Disks.Devices
	iotypes := config.Disks.Type
	// 遍历配置文件中指定的盘符列表
	for _, disk := range disks {
		// 获取磁盘smart信息
		smartInfo, err := module.GetDiskSmartInfo(disk)
		if err != nil {
			log.Fatal(err)
		}

		// 遍历配置文件中指定的测试类型
		for _, iotype := range iotypes {
			fiooutput, workload, err := module.ExecuteFio(disk, iotype, 128)
			if err != nil {
				log.Fatal(err)
			}

			// 解析 fio 执行结果
			ioResult, err := module.ParseFIOOutput(fiooutput)
			if err != nil {
				log.Fatal(err)
			}

			// 记录存储到数据库中
			err = db.SaveFIOResult(ioResult, workload, smartInfo)
			if err != nil {
				log.Fatal(err)
			}
		}
	}
}
