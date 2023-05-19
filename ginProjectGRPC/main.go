package main

import (
	"fmt"
	"ginProjectGRPC/common"
	"ginProjectGRPC/router"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

func main() {
	InitConfig() //初始化配置文件

	//初始化数据库
	common.InitDB()

	//创建路由
	r := gin.Default()
	router.CollectRoute(r)
	port := viper.GetString("server.port")
	if port != "" {
		panic(r.Run(":" + port)) // :号是灵魂
	} else {
		r.Run(port)
	}

}

// yml和json基本一致，yml可以注释
func InitConfig() {
	workDir, _ := os.Getwd() //获取当前工作目录
	fmt.Println(workDir)

	//设置读取的配置文件名
	viper.SetConfigName("application")
	//设置读取配置文件的类型
	viper.SetConfigType("yml")
	//设置读取配置文件的路径
	viper.AddConfigPath(workDir + "/config")
	err := viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("Fatal error config file: %s \n", err))
	}
}
