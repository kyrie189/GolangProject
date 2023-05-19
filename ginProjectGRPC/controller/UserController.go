package controller

import (
	"encoding/json"
	"fmt"
	"ginProjectGRPC/common"
	"ginProjectGRPC/dto"
	"ginProjectGRPC/model"
	"ginProjectGRPC/response"
	"ginProjectGRPC/util"
	"log"
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

// 用户注册
func Register(c *gin.Context) {
	db := common.GetDB()

	//获取参数方法1（不可取）
	/* name := c.PostForm("name")
	telephone := c.PostForm("telephone")
	password := c.PostForm("password")
	*/
	//使用map获取参数
	/* var requestMap = make(map[string]string)
	json.NewDecoder(c.Request.Body).Decode(&requestMap) */

	var requestUser = model.User{}
	c.ShouldBind(&requestUser)
	name := requestUser.Name
	telephone := requestUser.Telephone
	password := requestUser.Password

	fmt.Println("requestUser: ", requestUser)

	//fmt.Println("requestUser: ", telephone,"len: ", len(telephone))
	//数据验证
	if len(telephone) != 11 {
		response.Response(c,
			http.StatusUnprocessableEntity,
			422, nil, "手机号必须为11位")
		return
	}
	if len(password) < 6 {
		response.Response(c, http.StatusUnprocessableEntity,
			422, nil, "密码不能少于6位")
		return
	}
	//如果名称没有传就给一个随机字符串
	if len(name) == 0 {
		name = util.Randomstring(10)
	}

	//判断手机号是否存在
	if common.IsTelephoneExist(db, telephone) {
		response.Response(c, http.StatusUnprocessableEntity,
			422, nil, "用户已经存在")
		return
	}

	//用户密码加密
	hasedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		//系统级错误
		response.Response(c, http.StatusInternalServerError,
			500, nil, "加密错误")
		return
	}

	//创建用户
	newUser := model.User{
		Name:       name,
		Telephone:  telephone,
		Password:   string(hasedPassword),
		Expiration: 1,
	}
	db.Create(&newUser) //插入数据

	//生成token
	token, err := common.ReleaseToken(newUser)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": 500,
			"msg":  "生成token错误",
		})
		log.Printf("token generate error: %v", err)
		return
	}
	//向客户端返回结果
	response.Success(c, gin.H{"token": token}, "注册成功")
	//打印参数
	log.Println(name, telephone, password)
}

func Login(c *gin.Context) {
	fmt.Println("进入Login")
	//用户是否已经登陆
	session := sessions.Default(c)
	islogin := session.Get("is_logged_in")

	if islogin == nil {
		fmt.Println("第二次islogin为nil")
	} else {
		fmt.Println("第二次islogin不为nil:  ", islogin.(bool))
	}
	if islogin != nil && islogin.(bool) {
		response.Success(c, gin.H{"token": nil, "is_logged_in": true}, "登录成功")
		return
	}

	//登陆
	var requestMap = make(map[string]string)
	json.NewDecoder(c.Request.Body).Decode(&requestMap)
	//获取参数
	telephone := requestMap["Telephone"]
	password := requestMap["Password"]
	fmt.Println("手机号和密码", telephone, password)

	//数据验证
	if len(telephone) != 11 {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"code":         422,
			"msg":          "手机号必须为11位",
			"is_logged_in": false,
		})
		return
	}
	//用户验证
	db := common.GetDB()
	user := model.User{}
	db.Where("telephone=?", telephone).First(&user)
	if user.ID == 0 {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"code":         422,
			"msg":          "用户不存在",
			"is_logged_in": false,
		})
		return
	}
	//密码验证
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":         400,
			"msg":          "密码错误",
			"is_logged_in": false,
		})
		return
	}
	//登陆成功 设置session
	session.Set("is_logged_in", true)
	session.Save()
	islogin = session.Get("is_logged_in")
	if islogin == nil {
		fmt.Println("islogin为nil")
	} else {
		fmt.Println("islogin不为nil:  ", islogin.(bool))
	}
	//发放token
	token, err := common.ReleaseToken(user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":         500,
			"msg":          "生成token错误",
			"is_logged_in": true,
		})
		log.Printf("token generate error: %v", err)
		return
	}

	db.Model(&user).Update("expiration", 1) //更新expiration

	//返回结果
	response.Success(c, gin.H{"token": token, "is_logged_in": true}, "登录成功")

	//eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJVc2VySWQiOjUsImV4cCI6MTY4NDQyMTE4NywiaWF0IjoxNjgzODE2Mzg3LCJpc3MiOiJ3ZWlsaWZlbmciLCJzdWIiOiJ1c2VyIHRva2VuIn0.c2v9qiXQDpmUz4g0K2w5Jh2Q-Gl2RAXdL_UUI07M7Q4
	/*
		Header = {"alg":"HS256","typ":"JWT"}"}   //头部包含算法和类型
		payload = {UserID:3,"exp":1583476999,"iat":1583476999,"iss":"weilifeng","sub":"user token"}  //有效载荷包含用户ID，过期时间，发行时间，发行人，主题	}
	*/
}

func Info(c *gin.Context) {
	user, _ := c.Get("user")                //从上下文中获取user
	user = dto.ToUserDto(user.(model.User)) //转换为UserDto
	c.JSON(http.StatusOK, gin.H{"code": 200,
		"data": gin.H{"user": user}})
}
