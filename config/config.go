package config

import (
	"fmt"

	"github.com/spf13/viper"
)

var (
	MYSQL_CONFIG  = &MysqlConfig{Timeout: 10, MaxOpenConns: 100, MaxIdleConns: 20}
	SERVER_CONFIG = &ServerConfig{}
)

type MysqlConfig struct {
	Username     string `json:"username"`
	Password     string `json:"password"`
	Host         string `json:"host"`
	Port         int    `json:"port"`
	Dbname       string `json:"dbName"`
	Timeout      int    `json:"timeout"`
	MaxOpenConns int    `json:"maxOpenConns"`
	MaxIdleConns int    `json:"maxIdleConns"`
}

type ServerConfig struct {
	Name string `json:"name"`
	Host string `json:"host"`
	Port int    `json:"port"`
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

	if err := v.UnmarshalKey("mysql", MYSQL_CONFIG); err != nil {
		panic(err)
	}
	fmt.Println("mysql config: ", MYSQL_CONFIG)

	if err := v.UnmarshalKey("server", SERVER_CONFIG); err != nil {
		panic(err)
	}
	fmt.Println("server config: ", SERVER_CONFIG)
}
