package main

import (
	"context"
	"fmt"
	"ginProjectGRPC/common"
	"ginProjectGRPC/model"
	"ginProjectGRPC/proto/service"
	"math/rand"
	"net"
	"os"
	"sync"
	"time"

	"github.com/spf13/viper"
	"google.golang.org/grpc"
)

var mu sync.Mutex = sync.Mutex{}

type luckyService struct {
}

func (l *luckyService) Lucky(ctx context.Context, req *service.LuckyReq) (*service.LuckyRes, error) {
	//获取参数
	var id int = int(req.Id) //从客户端获取参数

	//获取随机数并抽奖
	giftNumber := 5
	code := luckyCode(int64(id)) // 获取随机数
	fmt.Println("获取随机数: ", code)
	if code >= 1 && code <= 10 {
		giftNumber = 1
	} else if code > 10 && code <= 1010 {
		giftNumber = 2
	} else if code > 1010 && code <= 3010 {
		giftNumber = 3
	} else if code > 3010 && code <= 5010 {
		giftNumber = 4
	}

	mu.Lock()
	defer mu.Unlock()
	//更新gift数据
	if giftNumber != 5 {
		GiftData := model.Gift{}
		common.Rd.HGetAll("gift_"+fmt.Sprintf("%d", giftNumber), &GiftData)
		fmt.Println("GiftData: ", GiftData)
		if GiftData.LeftNum <= 0 {
			giftNumber = 5 + giftNumber //代表已经发完了
		} else { // 更新数据库
			GiftData.LeftNum = GiftData.LeftNum - 1
			common.Rd.HSet("gift_"+fmt.Sprintf("%d", giftNumber), map[string]interface{}{
				"Id":      GiftData.Id,
				"LeftNum": GiftData.LeftNum,
			})
		}
		common.Rd.RedisToMysql(&GiftData)
	}
	//fmt.Println("giftNumber: ", giftNumber)
	//更新user数据
	//抽奖次数，中奖次数，中奖id，
	userdata := model.User{}
	common.Rd.HGetAll("user_"+fmt.Sprintf("%d", id), &userdata)
	userdata.Number = userdata.Number + 1
	if giftNumber < 5 {
		userdata.GiftNumber = userdata.GiftNumber + 1
	}
	userdata.GiftId = giftNumber
	common.Rd.HSet("user_"+fmt.Sprintf("%d", id), map[string]interface{}{
		"ID":         userdata.ID,
		"Name":       userdata.Name,
		"Telephone":  userdata.Telephone,
		"Password":   userdata.Password,
		"Number":     userdata.Number,
		"GiftNumber": userdata.GiftNumber,
		"GiftId":     userdata.GiftId,
	})
	//fmt.Println("userdata: ", userdata)
	common.Rd.RedisToMysql(&userdata) //更新数据库
	fmt.Println("用户id: ", id, "抽奖成功，中奖id: ", giftNumber, "剩余数量: ", userdata.GiftNumber, "抽奖数量: ", userdata.Number)
	res := &service.LuckyRes{}
	return res, nil
}

func main() {
	InitConfig() //初始化配置文件
	//初始化数据库
	common.InitDB()

	//1、初始化一个grpc的对象
	grpcServe := grpc.NewServer()
	//2、注册服务
	service.RegisterGreeterServer(grpcServe, &luckyService{})
	//3、监听端口,指定IP、port
	listener, err := net.Listen("tcp", ":8080")
	if err != nil {
		fmt.Println("net.Listen err:", err)
		return
	}
	defer listener.Close()
	//4、启动服务
	fmt.Println("服务端已启动")
	grpcServe.Serve(listener)

}

func luckyCode(id int64) int32 {
	rateMax := 10000
	seed := time.Now().UnixNano()
	seed = seed + id
	r := rand.New(rand.NewSource(seed))
	code := r.Int31n(int32(rateMax)) + 1
	return code //返回一个随机数[1,10000]
}

// yml和json基本一致，yml可以注释
func InitConfig() {
	workDir, _ := os.Getwd() //获取当前工作目录
	fmt.Println(workDir)
	workDir = "F:/go/gocode/ginProjectGRPC"
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
