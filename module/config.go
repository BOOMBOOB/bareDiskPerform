// @Project -> File    : bare-disk-perform -> config
// @IDE    : GoLand
// @Author    : wuji
// @Date   : 2023/8/22 10:42

package module

import (
	"encoding/json"
	"log"
	"os"
)

type Config struct {
	Mysql    MysqlConfig `json:"mysql"`
	Disks    DisksConfig `json:"disks"`
	Level    string      `json:"level"`
	RampTime string      `json:"ramp_time"`
	Runtime  string      `json:"runtime"`
	Iodepth  string      `json:"iodepth"`
}

type MysqlConfig struct {
	Server   string `json:"server"`
	Port     string `json:"port"`
	Database string `json:"database"`
	Username string `json:"username"`
	Password string `json:"password"`
}

type DisksConfig struct {
	Mode    string   `json:"mode"`
	Type    []string `json:"type"`
	Devices []string `json:"devices"`
}

func LoadConfig(filepath string) (Config, error) {
	var config Config

	configData, err := os.ReadFile(filepath)
	if err != nil {
		log.Println("Error reading.")
		return config, err
	}

	err = json.Unmarshal(configData, &config)
	if err != nil {
		log.Println("Error decoding config.")
		return config, err
	}
	log.Println("Config loaded successfully.")
	return config, nil
}
