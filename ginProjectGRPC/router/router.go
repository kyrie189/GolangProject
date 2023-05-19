package router

import (
	"ginProjectGRPC/controller"
	"ginProjectGRPC/middleware"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
)

func CollectRoute(r *gin.Engine) *gin.Engine {
	// 初始化 session 中间件
	store := cookie.NewStore([]byte("miyao")) // 替换为实际的密钥
	store.Options(sessions.Options{
		MaxAge: 86400, // 设置持续时间为 24 小时
		Path:   "/",
	})
	//使用中间件
	r.Use(middleware.CORSMiddleware())
	//创建一个会话管理中间件，命名为session,该中间件将会话存储（store）传入
	r.Use(sessions.Sessions("session", store))

	//
	r.POST("/api/auth/register", controller.Register)
	r.POST("/api/auth/login", controller.Login)
	//登陆成功获取数据
	r.GET("/api/data", controller.GetPrizeInfo)
	//抽奖
	r.POST("api/lucky", controller.Lucky)

	//使用中间件保护用户登陆的接口
	r.GET("/api/auth/info", middleware.AuthMiddleware(), controller.Info)

	return r
}
