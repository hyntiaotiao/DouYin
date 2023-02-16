package main

import (
	"DouYIn/config"
	"fmt"

	"github.com/spf13/viper"
)

func initConfig() {
	v := viper.New()
	v.SetConfigFile("./application.yml")
	v.SetConfigType("yaml")
	// 加载配置文件内容
	if err := v.ReadInConfig(); err != nil {
		if _, ok := err.(v.ConfigFileNotFoundError); ok {
			// 配置文件未找到
			panic("配置文件未找到")
		} else {
			// Config file was found but another error was produced
			panic(fmt.Errorf("读取配置文件失败: %s \n", err.Error()))
		}
	}

	// 读取文件配置项
	mysqlConfig := config.MysqlConfig{}
	if err := v.Unmarshal(&mysqlConfig); err != nil {
		panic(err)
	}
}
