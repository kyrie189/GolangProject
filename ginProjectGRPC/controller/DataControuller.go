package controller

import (
	"fmt"
	"ginProjectGRPC/common"
	"ginProjectGRPC/model"
	"ginProjectGRPC/response"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

func GetPrizeInfo(c *gin.Context) {
	//用户是否已经登陆
	session := sessions.Default(c)
	islogin := session.Get("is_logged_in")
	if islogin == nil {
		fmt.Println("islogin为nil")
	} else {
		fmt.Println("islogin不为nil:  ", islogin.(bool))
	}
	if islogin == nil || !islogin.(bool) {
		response.Success(c, gin.H{"is_logged_in": false}, "注册成功")
		return
	}

	db := common.GetDB()

	// 从数据库获取gift数据
	giftList := []model.Gift{}
	result := db.Find(&giftList)
	if result.Error != nil {
		// 获取数据时发生错误
		fmt.Println("giftList Error occurred while retrieving data:", result.Error)
	}
	// fmt.Println("giftList: ", giftList)

	//获取user数据
	userList := []model.User{}
	result = db.Find(&userList)
	if result.Error != nil {
		// 获取数据时发生错误
		fmt.Println("userlist Error occurred while retrieving data:", result.Error)
	}
	// fmt.Println("giftList: ", userList)

	//向客户端返回结果, "user": userList
	response.Success(c, gin.H{"gift": giftList, "user": userList, "is_logged_in": true}, "注册成功")
}
