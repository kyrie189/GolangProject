package common

import (
	"context"
	"fmt"
	"ginProjectGRPC/model"
	"strconv"
	"github.com/go-redis/redis/v8"
)

var (
	Rd *RedisDatabase
)

type RedisDatabase struct {
	rds *redis.Client
	ctx context.Context
}

// redis 线程不安全
func init() {
	Rd = &RedisDatabase{
		rds: redis.NewClient(&redis.Options{
			Addr:     "localhost:6379",
			Password: "", // no password set
			DB:       14, // use default DB
			PoolSize: 1,
		}),
		ctx: context.Background(),
	}
	
	//初始化数据库
	initGiftRedis()
}

func (Rd *RedisDatabase) HGetAll(key string, storer interface{}) {
	data, err := Rd.rds.HGetAll(Rd.ctx, key).Result()
	if err != nil {
		fmt.Println("获取全部的hash数据失败: ", err)
	}
	switch storer.(type) {
	case *model.Gift: //只获取id和leftNum
		gift := storer.(*model.Gift)
		gift.Id, err = strconv.Atoi(data["Id"])
		if err != nil {
			fmt.Println("获取Id失败(转int失败): ", err)
		}
		gift.Title = data["Title"]
		gift.PrizeNum, err = strconv.Atoi(data["PrizeNum"])
		if err != nil {
			fmt.Println("获取PrizeNum失败(转int失败): ", err)
		}
		gift.LeftNum, err = strconv.Atoi(data["LeftNum"])
		if err != nil {
			fmt.Println("获取LeftNum失败(转int失败): ", err)
		}
		gift.Gtype, err = strconv.Atoi(data["Gtype"])
		if err != nil {
			fmt.Println("获取Gtype失败(转int失败): ", err)
		}
		gift.TimeBegin, err = strconv.Atoi(data["TimeBegin"])
		if err != nil {
			fmt.Println("获取TimeBegin失败(转int失败): ", err)
		}
		gift.TimeEnd, err = strconv.Atoi(data["TimeEnd"])
		if err != nil {
			fmt.Println("获取TimeEnd失败(转int失败): ", err)
		}
	case *model.User:
		user := storer.(*model.User)
		id, err := strconv.Atoi(data["ID"])
		if err != nil {
			fmt.Println("获取Id失败(转int失败): ", err)
		}
		user.ID = uint(id)
		user.Name = data["Name"]
		user.Telephone = data["Telephone"]
		user.Password = data["Password"]
		user.Number, err = strconv.Atoi(data["Number"])
		if err != nil {
			fmt.Println("获取Number失败(转int失败): ", err)
		}
		user.GiftNumber, err = strconv.Atoi(data["GiftNumber"])
		if err != nil {
			fmt.Println("获取GiftNumber失败(转int失败): ", err)
		}
		user.GiftId, err = strconv.Atoi(data["GiftId"])
		if err != nil {
			fmt.Println("获取GiftId失败(转int失败): ", err)
		}
	}
}

func (Rd *RedisDatabase) HSet(key string, data map[string]interface{}) {
	ok, err := Rd.rds.HMSet(Rd.ctx, key, data).Result()
	if err != nil || !ok {
		fmt.Println("更新数据失败: ", err)
	}
}

func initGiftRedis() {
	gift := model.Gift{
		Id:        1,
		Title:     "iphone14proMax",
		PrizeNum:  10, //总数
		LeftNum:   10, //剩余数量
		Gtype:     3,  //大奖
		TimeBegin: 0,
		TimeEnd:   0,
	}
	ok, err := Rd.rds.HMSet(Rd.ctx, fmt.Sprintf("gift_%d", gift.Id), map[string]interface{}{
		"Id":        gift.Id,
		"Title":     gift.Title,
		"PrizeNum":  gift.PrizeNum,
		"LeftNum":   gift.LeftNum,
		"Gtype":     gift.Gtype,
		"TimeBegin": gift.TimeBegin,
		"TimeEnd":   gift.TimeEnd,
	}).Result()
	if err != nil || !ok {
		fmt.Println("redis set failed: 11111", err)
	}
	//2
	gift = model.Gift{
		Id:        2,
		Title:     "充电器",
		PrizeNum:  1000, //总数
		LeftNum:   1000, //剩余数量
		Gtype:     2,    //小奖
		TimeBegin: 0,
		TimeEnd:   0,
	}
	ok, err = Rd.rds.HMSet(Rd.ctx, fmt.Sprintf("gift_%d", gift.Id), map[string]interface{}{
		"Id":        gift.Id,
		"Title":     gift.Title,
		"PrizeNum":  gift.PrizeNum,
		"LeftNum":   gift.LeftNum,
		"Gtype":     gift.Gtype,
		"TimeBegin": gift.TimeBegin,
		"TimeEnd":   gift.TimeEnd,
	}).Result()
	if err != nil || !ok {
		fmt.Println("redis set failed: 22222", err)
	}
	//3
	gift = model.Gift{
		Id:        3,
		Title:     "优惠卷",
		PrizeNum:  2000, //总数
		LeftNum:   2000, //剩余数量
		Gtype:     1,    //虚拟卷
		TimeBegin: 0,
		TimeEnd:   0,
	}
	ok, err = Rd.rds.HMSet(Rd.ctx, fmt.Sprintf("gift_%d", gift.Id), map[string]interface{}{
		"Id":        gift.Id,
		"Title":     gift.Title,
		"PrizeNum":  gift.PrizeNum,
		"LeftNum":   gift.LeftNum,
		"Gtype":     gift.Gtype,
		"TimeBegin": gift.TimeBegin,
		"TimeEnd":   gift.TimeEnd,
	}).Result()
	if err != nil || !ok {
		fmt.Println("redis set failed: 3333", err)
	}
	//4
	gift = model.Gift{
		Id:        4,
		Title:     "虚拟币",
		PrizeNum:  2000, //总数
		LeftNum:   2000, //剩余数量
		Gtype:     0,    //小奖
		TimeBegin: 0,
		TimeEnd:   0,
	}
	ok, err = Rd.rds.HMSet(Rd.ctx, fmt.Sprintf("gift_%d", gift.Id), map[string]interface{}{
		"Id":        gift.Id,
		"Title":     gift.Title,
		"PrizeNum":  gift.PrizeNum,
		"LeftNum":   gift.LeftNum,
		"Gtype":     gift.Gtype,
		"TimeBegin": gift.TimeBegin,
		"TimeEnd":   gift.TimeEnd,
	}).Result()
	if err != nil || !ok {
		fmt.Println("redis set failed: 3333", err)
	}
}

func (Rd *RedisDatabase) RedisToMysql(storer interface{}) {
	switch storer.(type) {
	case *model.Gift: //只获取id和leftNum
		//从redis获取gift数据
		giftList := make([]*model.Gift, 4)
		for i := 1; i <= 4; i++ {
			gift := model.Gift{}
			Rd.HGetAll(fmt.Sprintf("gift_%d", i), &gift)
			giftList[i-1] = &gift
		}
		//批量插入数据
		for _, gift := range giftList {
			result := DB.Save(&gift) //create创建 ，已经创建好就不会创建字段
			if result.Error != nil {
				// 插入数据时发生错误
				fmt.Println("Error occurred while inserting data:", result.Error)
			}
		}
	case *model.User:
		//从redis获取user数据
		userList := make([]*model.User, 5)
		for i := 1; i <= 5; i++ {
			user := model.User{}
			Rd.HGetAll(fmt.Sprintf("user_%d", i), &user)
			userList[i-1] = &user
		}
		fmt.Println()
		fmt.Println("userList: ", userList)
		fmt.Println()
		//批量插入数据
		for _, user := range userList {
			//没更新成功
			DB.Model(&user).Where("id = ?", user.ID).Update("number", user.Number)
			DB.Model(&user).Where("id = ?", user.ID).Update("gift_number", user.GiftNumber)
			DB.Model(&user).Where("id = ?", user.ID).Update("gift_id", user.GiftId)

			//result := DB.Model(&user).Where("id=?", user.ID).Updates(&user)
			//result := DB.Save(&user) //create创建 ，已经创建好就不会创建字段
			/* if result.Error != nil {
				// 插入数据时发生错误
				fmt.Println("Error occurred while inserting data:", result.Error)
			} */
		}
	}
}

func (Rd *RedisDatabase) MysqlToRedis(storer interface{}) {
	switch storer.(type) {
	case *model.Gift: //只获取id和leftNum
		//从数据库获取gift数据
		giftList := []model.Gift{}
		result := DB.Find(&giftList)
		if result.Error != nil {
			// 获取数据时发生错误
			fmt.Println("MysqlToRedis error giftFind错误:", result.Error)
		}
		//批量插入数据
		for _, gift := range giftList {
			Rd.HSet(fmt.Sprintf("gift_%d", gift.Id), map[string]interface{}{
				"Id":        gift.Id,
				"Title":     gift.Title,
				"PrizeNum":  gift.PrizeNum,
				"LeftNum":   gift.LeftNum,
				"Gtype":     gift.Gtype,
				"TimeBegin": gift.TimeBegin,
				"TimeEnd":   gift.TimeEnd,
			})
		}
	case *model.User:
		//从数据库获取user数据
		fmt.Println("userList: ")
		var userList []model.User
		result := DB.Find(&userList)

		if result.Error != nil {
			// 获取数据时发生错误
			fmt.Println("MysqlToRedis error userFind错误:", result.Error)
		}
		//批量插入数据
		for _, user := range userList {
			Rd.HSet(fmt.Sprintf("user_%d", user.ID), map[string]interface{}{
				"ID":         user.ID,
				"Name":       user.Name,
				"Telephone":  user.Telephone,
				"Password":   user.Password,
				"Number":     user.Number,
				"GiftNumber": user.GiftNumber,
				"GiftId":     user.GiftId,
			})
		}
	}
}
