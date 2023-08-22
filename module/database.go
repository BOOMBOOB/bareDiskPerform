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
		return nil, err
	}
	return &Database{db: db}, nil
}

func (d *Database) Close() error {
	if d.db != nil {
		err := d.db.Close()
		if err != nil {
			return err
		}
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
		return err
	}

	if count > 0 {
		// 如果存在相同组合，则执行更新操作
		_, err = d.db.Exec("UPDATE bareTest SET IOPS = ?, BandWidth = ?, ClatAvg = ?, Clat95 = ?, Clat99 = ?, SMART = ?, timestamp = ? WHERE SerialNumber = ? AND BlockSize = ? AND IODepth = ? AND IOType = ?",
			result.Iops, result.BandWidth, result.ClatAvg, result.Clat95, result.Clat99, disksmart.SMART, time.Now(), disksmart.SerialNumber, workload.BlockSize, workload.IODepth, workload.IOType)
		if err != nil {
			return err
		}
	} else {
		// 如果不存在相同组合，则执行插入操作
		_, err = d.db.Exec("INSERT INTO bareTest (SerialNumber, BlockSize, IODepth, IOType, IOPS, BandWidth, ClatAvg, Clat95, Clat99, SMART, Model, UserCapacity, RotationRate, FormFactor, timestamp) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)",
			disksmart.SerialNumber, workload.BlockSize, workload.IODepth, workload.IOType, result.Iops, result.BandWidth, result.ClatAvg, result.Clat95, result.Clat99, disksmart.SMART, disksmart.DeviceModel, disksmart.UserCapacity, disksmart.RotationRate, disksmart.FormFactor, time.Now())
		if err != nil {
			return err
		}
	}
	return nil
}
