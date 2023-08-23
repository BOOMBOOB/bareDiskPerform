// @Project -> File    : bare-disk-perform -> main
// @IDE    : GoLand
// @Author    : wuji
// @Date   : 2023/8/22 10:40

package main

import (
	"bare-disk-perform/module"
)

func main() {

	// 初始化日志
	module.InitMyLogger()
	logger := module.GetLogger()

	logger.Infof("load config file....")
	config, err := module.LoadConfig("./config.json")
	if err != nil {
		logger.Fatalf("load config file failed: %v", err)
	}
	logger.Infof("load config file success.")

	// 获取db 连接
	logger.Infof("get db connection....")
	db, err := module.NewDatabase(config)
	if err != nil {
		logger.Fatalf("get db connection failed: %v.", err)
	}
	logger.Infof("get db connection success.")

	defer func() {
		if err := db.Close(); err != nil {
			logger.Infof("failed to close database connection: %v.", err)
		}
		logger.Infof("close database connection successfully.")
	}()

	disks := config.Disks.Devices
	if config.Disks.Mode == "auto" {
		disks, err = module.GetAutoScanDisks()
		if err != nil {
			logger.Fatalf("auto get disks failed: %v.", err)
		}
	}

	iotypes := config.Disks.Type

	// 遍历配置文件中指定的盘符列表
	for _, disk := range disks {
		// 获取磁盘smart信息
		logger.Infof("get disk %v smart info....", disk)
		smartInfo, err := module.GetDiskSmartInfo(disk)
		if err != nil {
			logger.Fatalf("get disk smart info failed: %v.", err)
		}
		logger.Infof("get disk smart info successfully.")

		// 遍历配置文件中指定的测试类型
		for _, iotype := range iotypes {
			logger.Infof("perform disk %v load %v testing....", disk, iotype)
			fiooutput, workload, err := module.ExecuteFio(disk, iotype, config)
			if err != nil {
				logger.Fatalf("execute load testing failed: %v.", err)
			}
			logger.Infof("load testing success.")

			// 解析 fio 执行结果
			logger.Infof("analyse fio test result....")
			ioResult, err := module.ParseFIOOutput(fiooutput)
			if err != nil {
				logger.Fatalf("parse test result failed: %v.", err)
			}
			logger.Infof("fio test result analyze successfully.")

			// 记录存储到数据库中
			logger.Infof("save test result into mysql database....")
			err = db.SaveFIOResult(ioResult, workload, smartInfo)
			if err != nil {
				logger.Fatalf("save test result failed: %v.", err)
			}
			logger.Infof("save test result success.")
		}
	}
}
