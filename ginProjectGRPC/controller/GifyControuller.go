package controller

import (
	"context"
	"fmt"
	"ginProjectGRPC/proto/service"
	"ginProjectGRPC/response"

	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func Lucky(c *gin.Context) {
	//1、连接服务器
	grpcClient, err := grpc.Dial("localhost:8080", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		fmt.Println("grpc.Dial err:", err)
		return
	}
	//2、注册客户端
	Client := service.NewGreeterClient(grpcClient)
	//contenxt.background()是一个空的context

	//获取参数
	var id int = 0
	c.ShouldBindJSON(&id)
	fmt.Println("获取id: ", id)
	//发送给服务端
	Client.Lucky(context.Background(), &service.LuckyReq{Id: int32(id)})

	response.Success(c, nil, "抽奖成功")
}
