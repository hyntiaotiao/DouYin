package main

import (
	"DouYIn/config"
	"DouYIn/controller"
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

var (
	serverConfig config.ServerConfig
)

func main() {
	log.Println("开始！")

	run_server := initConfig()

	r := gin.Default()

	//r.Use(utils.JwtVerify) // Jwt验证中间件

	initRouter(r)

	r.Run(run_server) // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}

func initConfig() (run_server string) {
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

	if err := v.UnmarshalKey("mysql", config.MYSQL_CONFIG); err != nil {
		panic(err)
	}
	fmt.Println("mysql config: ", config.MYSQL_CONFIG)

	if err := v.UnmarshalKey("server", config.SERVER_CONFIG); err != nil {
		panic(err)
	}
	fmt.Println("server config: ", config.SERVER_CONFIG)
	run_server = config.SERVER_CONFIG.Host + ":" + fmt.Sprint(config.SERVER_CONFIG.Port)
	fmt.Println(run_server)
	return run_server
}

func initRouter(r *gin.Engine) {
	// public directory is used to serve static resources
	//r.Static("/static", "./public")

	r.MaxMultipartMemory = 128 << 20 //设置视频最大上传容量

	apiRouter := r.Group("/douyin")

	// basic apis
	apiRouter.GET("/feed/", controller.Feed)
	apiRouter.GET("/user/", controller.UserInfo)
	apiRouter.POST("/user/register/", controller.Register)
	apiRouter.POST("/user/login/", controller.Login)

	apiRouter.POST("/publish/action/", controller.Publish)
	apiRouter.GET("/publish/list/", controller.PublishList)
	//
	//// extra apis - I
	apiRouter.POST("/favorite/action/", controller.FavoriteAction)
	apiRouter.GET("/favorite/list/", controller.FavoriteList)
	apiRouter.POST("/comment/action/", controller.CommentAction)
	apiRouter.GET("/comment/list/", controller.CommentList)
	//
	//// extra apis - II
	apiRouter.POST("/relation/action/", controller.RelationAction)
	apiRouter.GET("/relation/follow/list/", controller.FollowList)
	apiRouter.GET("/relation/follower/list/", controller.FollowerList)
	apiRouter.GET("/relation/friend/list/", controller.FriendList)
	apiRouter.GET("/message/chat/", controller.MessageChat)
	apiRouter.POST("/message/action/", controller.MessageAction)
}
