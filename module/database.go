// @Project -> File    : bare-disk-perform -> database
// @IDE    : GoLand
// @Author    : wuji
// @Date   : 2023/8/22 10:41

package module

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"time"
)

type Database struct {
	db *sql.DB
}

func NewDatabase(config Config) (*Database, error) {
	datasource := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", config.Mysql.Username, config.Mysql.Password,
		config.Mysql.Server, config.Mysql.Port, config.Mysql.Database)
	db, err := sql.Open("mysql", datasource)
	if err != nil {
		logger.Errorf("get database connection failed: %v.", err)
		return nil, err
	}
	logger.Debugf("get database connection successfully.")
	return &Database{db: db}, nil
}

func (d *Database) Close() error {
	if d.db != nil {
		err := d.db.Close()
		if err != nil {
			logger.Errorf("failed to close database connection: %v.", err)
			return err
		}
		logger.Debugf("close database connection successfully.")
	}
	return nil
}

func (d *Database) SaveFIOResult(result Result, workload WorkLoad, disksmart DiskSmart) error {
	// 查询数据库表中是否存在相同的 SerialNumber、BlockSize、IODepth 和 IOType 组合
	row := d.db.QueryRow("SELECT COUNT(*) FROM bareTest WHERE SerialNumber = ? AND BlockSize = ? AND IODepth = ? AND IOType = ?",
		disksmart.SerialNumber, workload.BlockSize, workload.IODepth, workload.IOType)
	var count int
	err := row.Scan(&count)
	if err != nil {
		logger.Errorf("Failed to search data in database: %v.", err)
		return err
	}
	logger.Debugf("search date in database successfully. the count number is %d.", count)

	// 获取东八区本地时间
	loc, err := time.LoadLocation("Asia/Shanghai")
	if err != nil {
		logger.Fatalf("get local time error: %v.", err)
	}
	current := time.Now().In(loc)
	currentTime := current.Format("2006-01-02T15:04:05")

	if count > 0 {
		// 如果存在相同组合，则执行更新操作
		_, err = d.db.Exec("UPDATE bareTest SET IOPS = ?, BandWidth = ?, ClatAvg = ?, Clat95 = ?, Clat99 = ?, SMART = ?, timestamp = ? WHERE SerialNumber = ? AND BlockSize = ? AND IODepth = ? AND IOType = ?",
			result.Iops, result.BandWidth, result.ClatAvg, result.Clat95, result.Clat99, disksmart.SMART, currentTime, disksmart.SerialNumber, workload.BlockSize, workload.IODepth, workload.IOType)
		if err != nil {
			logger.Errorf("Failed to update count in database: %v.", err)
			return err
		}
		logger.Debugf("update count in database successfully.")
	} else {
		// 如果不存在相同组合，则执行插入操作
		_, err = d.db.Exec("INSERT INTO bareTest (SerialNumber, BlockSize, IODepth, IOType, IOPS, BandWidth, ClatAvg, Clat95, Clat99, SMART, Model, UserCapacity, RotationRate, FormFactor, timestamp) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)",
			disksmart.SerialNumber, workload.BlockSize, workload.IODepth, workload.IOType, result.Iops, result.BandWidth, result.ClatAvg, result.Clat95, result.Clat99, disksmart.SMART, disksmart.DeviceModel, disksmart.UserCapacity, disksmart.RotationRate, disksmart.FormFactor, currentTime)
		if err != nil {
			logger.Errorf("Failed to insert count into database: %v.", err)
			return err
		}
		logger.Debugf("insert count into database successfully.")
	}
	return nil
}
