package common

import (
	"fmt"
	"ginProjectGRPC/model"
	"github.com/spf13/viper"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB = nil

// 连接数据库 并创建一张表
func InitDB() *gorm.DB {
	host := viper.GetString("datasource.host")
	port := viper.GetString("datasource.port")
	username := viper.GetString("datasource.username")
	password := viper.GetString("datasource.password")
	database := viper.GetString("datasource.database")
	charset := viper.GetString("datasource.charset")
	args := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=%s&parseTime=True&loc=Local",
		username, password, host, port, database, charset)
	db, err := gorm.Open(mysql.Open(args), &gorm.Config{})
	if err != nil {
		panic("failed to connect database,err:" + err.Error())
	}
	DB = db
	//创建mysql的表
	DB.AutoMigrate(&model.User{})
	MigrateTable() //创建gift和result表

	initGift() //初始化gift表
	initUser() //初始化user表
	Rd.MysqlToRedis(&model.User{})
	fmt.Println("数据库连接成功")
	return DB
}

func IsTelephoneExist(db *gorm.DB, telephone string) bool {
	var user model.User
	db.Where("telephone = ?", telephone).First(&user) //没找打就返回空
	return user.ID != 0
}

func GetDB() *gorm.DB {
	return DB
}

func MigrateTable() {
	gift := model.Gift{}
	result := model.Result{}
	DB.AutoMigrate(&gift)
	DB.AutoMigrate(&result)
}

func initGift() {
	giftList := make([]*model.Gift, 4)
	g1 := model.Gift{
		Id:        1,
		Title:     "iphone14proMax",
		PrizeNum:  10, //总数
		LeftNum:   10, //剩余数量
		Gtype:     3,  //大奖
		TimeBegin: 0,
		TimeEnd:   0,
	}
	giftList[0] = &g1
	g2 := model.Gift{
		Id:        2,
		Title:     "充电器",
		PrizeNum:  1000, //总数
		LeftNum:   1000, //剩余数量
		Gtype:     2,    //小奖
		TimeBegin: 0,
		TimeEnd:   0,
	}
	giftList[1] = &g2
	g3 := model.Gift{
		Id:        3,
		Title:     "优惠卷",
		PrizeNum:  2000, //总数
		LeftNum:   2000, //剩余数量
		Gtype:     1,    //虚拟卷
		TimeBegin: 0,
		TimeEnd:   0,
	}
	giftList[2] = &g3
	g4 := model.Gift{
		Id:        4,
		Title:     "虚拟币",
		PrizeNum:  2000, //总数
		LeftNum:   2000, //剩余数量
		Gtype:     0,    //小奖
		TimeBegin: 0,
		TimeEnd:   0,
	}
	giftList[3] = &g4

	//抽奖方法
	//在1-10000取一个随机数，看看这个数落在那个区间

	//批量插入数据
	for _, user := range giftList {
		result := DB.Save(&user) //create创建 ，已经创建好就不会创建字段
		if result.Error != nil {
			// 插入数据时发生错误
			fmt.Println("Error occurred while inserting data:", result.Error)
		}
	}
}

func initUser() {
	user := model.User{}
	for i := 1; i <= 5; i++ {
		DB.Where("id = ?", i).Find(&user)
		user.Number = 0
		user.GiftNumber = 0
		user.GiftId = 5
		DB.Save(&user)
		user = model.User{}
	}
}
