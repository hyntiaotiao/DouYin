package config

import (
	"fmt"

	"github.com/spf13/viper"
	"gorm.io/gorm/logger"
)

var (
	MYSQL_CONFIG  = &MysqlConfig{Timeout: 10, MaxOpenConns: 100, MaxIdleConns: 20, SingularTable: true}
	SERVER_CONFIG = &ServerConfig{}
	LOGGER_CONFIG = &LoggerConfig{SlowThreshold: 0, LogLevel: logger.Info, Colorful: false, IgnoreRecordNotFoundError: true}
)

type MysqlConfig struct {
	Username      string `json:"username"`
	Password      string `json:"password"`
	Host          string `json:"host"`
	Port          int    `json:"port"`
	Dbname        string `json:"dbName"`
	Timeout       int    `json:"timeout"`
	MaxOpenConns  int    `json:"maxOpenConns"`
	MaxIdleConns  int    `json:"maxIdleConns"`
	SingularTable bool   `json:"singularTable"`
}

type ServerConfig struct {
	Name string `json:"name"`
	Port int    `json:"port"`
}

type LoggerConfig struct {
	SlowThreshold             int             // 慢 SQL 阈值
	LogLevel                  logger.LogLevel // 日志级别Info, Warn, Error, Silent
	IgnoreRecordNotFoundError bool            // 忽略ErrRecordNotFound（记录未找到）错误
	Colorful                  bool            // 禁用彩色打印
}

func init() {
	fmt.Println("初始化参数：")

	v := viper.New()
	configFileName := "application.yml"
	v.SetConfigFile("./" + configFileName)
	v.SetConfigType("yaml")
	// 加载配置文件内容
	if err := v.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			fmt.Println(configFileName + "配置文件没找到.")
		} else {
			fmt.Println("读取配置文件发生错误：", err)
		}
	}

	//server配置
	if err := v.UnmarshalKey("server", SERVER_CONFIG); err != nil {
		panic(err)
	}
	fmt.Println("server config: ", SERVER_CONFIG)

	// mysql配置
	if err := v.UnmarshalKey("gorm.mysql", MYSQL_CONFIG); err != nil {
		panic(err)
	}
	fmt.Println("mysql config: ", MYSQL_CONFIG)

	// gorm logger日志
	if err := v.UnmarshalKey("gorm.logger", LOGGER_CONFIG); err != nil {
		panic(err)
	}
	fmt.Println("logger config: ", LOGGER_CONFIG)

}
