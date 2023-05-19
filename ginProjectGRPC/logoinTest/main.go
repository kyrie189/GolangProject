package main

import (
	"fmt"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"net/http"
)

func main() {
	// 创建 Gin 引擎
	router := gin.Default()

	// 初始化 session 中间件
	store := cookie.NewStore([]byte("secret_key")) 
	//创建一个会话管理中间件，命名为session,该中间件将会话存储（store）传入
	router.Use(sessions.Sessions("session", store)) 

	// 登录路由处理函数
	router.GET("/login", func(c *gin.Context) {
		// 重定向到登录页面
		c.Redirect(http.StatusTemporaryRedirect, "/login.html")
	})

	// 处理登录表单提交
	router.POST("/login", func(c *gin.Context) {
		/* session := sessions.Default(c)
		islogin := session.Get("is_logged_in")
		if islogin != nil && islogin.(bool) {
			c.Redirect(http.StatusTemporaryRedirect, "/choujiang")
			return
		} */

		username := c.PostForm("username")
		password := c.PostForm("password")
		fmt.Println(username)
		fmt.Println(password)
		// 验证用户名和密码，这里简化为检查是否匹配
		
		if username == "admin" && password == "123456" {
			// 登录成功，设置会话标记
			session := sessions.Default(c)
			session.Set("is_logged_in", true)
			session.Save()

			// 重定向到抽奖页面
			c.Redirect(http.StatusTemporaryRedirect, "/choujiang")
			return
		}

		// 登录失败，重定向回登录页面
		c.Redirect(http.StatusTemporaryRedirect, "/login")
	})

	// 抽奖页面
	router.GET("/choujiang", func(c *gin.Context) {
		// 检查用户登录状态
		session := sessions.Default(c)
		isLoggedIn := session.Get("is_logged_in")
		if isLoggedIn != nil && isLoggedIn.(bool) {
			// 用户已登录，显示抽奖页面
			c.HTML(http.StatusOK, "choujiang.html", nil)
		} else {
			// 用户未登录，重定向回登录页面
			c.Redirect(http.StatusTemporaryRedirect, "/login")
		}
	})

	// 运行服务器
	fmt.Println("服务器已启动，监听端口 8081")
	router.Run(":8081")
}
