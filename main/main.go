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
	fmt.Println("加载配置文件....")
	config, err := module.LoadConfig("./config.json")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("配置文件加载成功.")

	// 获取db 连接
	fmt.Println("获取 db 连接....")
	db, err := module.NewDatabase(config)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("获取 db 连接成功.")

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
		fmt.Printf("获取磁盘 %v smart信息....", disk)
		smartInfo, err := module.GetDiskSmartInfo(disk)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println("获取磁盘smart信息成功")

		// 遍历配置文件中指定的测试类型
		for _, iotype := range iotypes {
			fmt.Printf("对磁盘 %v 进行 %v 负载测试....", disk, iotype)
			fiooutput, workload, err := module.ExecuteFio(disk, iotype, 128)
			if err != nil {
				log.Fatal(err)
			}
			fmt.Println("负载测试完毕.")

			// 解析 fio 执行结果
			fmt.Println("分析 fio 测试结果....")
			ioResult, err := module.ParseFIOOutput(fiooutput)
			if err != nil {
				log.Fatal(err)
			}
			fmt.Println("fio测试结果分析完毕.")

			// 记录存储到数据库中
			fmt.Println("记录测试结果到数据库中....")
			err = db.SaveFIOResult(ioResult, workload, smartInfo)
			if err != nil {
				log.Fatal(err)
			}
			fmt.Println("记录完毕.")
		}
	}
}
