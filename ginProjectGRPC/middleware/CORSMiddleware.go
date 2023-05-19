package middleware

import (
	_ "fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		origin := c.GetHeader("Origin")
		c.Writer.Header().Set("Access-Control-Allow-Origin", origin) //更改为前端域名//不能使用*
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Max-Age", "86400")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST,GET,OPTIONS,PUT,DELETE,UPDATE")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")  //不能使用*
		if c.Request.Method == http.MethodOptions {
			c.AbortWithStatus(200)
		} else {
			c.Next() //将请求交给以一个处理函数
		}
	}
}
