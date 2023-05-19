package middleware

import (
	"ginProjectGRPC/common"
	"ginProjectGRPC/model"
	"log"
	"net/http"
	"strings"
	"github.com/gin-gonic/gin"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		//获取authorization header
		tokenString := c.GetHeader("Authorization")

		//验证token格式
		//oauth2.0规定，所有的token都应该以Bearer开头
		log.Println(tokenString)
		if tokenString == "" || !strings.HasPrefix(tokenString, "Bearer ") {
			c.JSON(http.StatusUnauthorized,
				gin.H{"code": 401, "msg": "(无Bearer)权限不足"})
			c.Abort()
			return
		}
		tokenString = tokenString[7:] //把Bearer 截取掉

		//解析token
		token, claims, err := common.ParseToken(tokenString)
		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized,
				gin.H{"code": 401, "msg": "(解析token)权限不足"})
			c.Abort()
			return
		}

		//验证通过后获取claim中的userId
		userId := claims.UserId //通过id查找mysql
		DB := common.GetDB()
		var user model.User
		DB.First(&user, userId)

		//用户
		if user.ID == 0 {
			c.JSON(http.StatusUnauthorized,
				gin.H{"code": 401, "msg": "(用户不存在)权限不足"})
			c.Abort()
			return
		}

		//用户存在 将user信息写入上下文
		c.Set("user", user)
		c.Next()
	}
}
