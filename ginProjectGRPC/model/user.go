package model

import (
	"gorm.io/gorm"
)

type User struct {
	gorm.Model        //含有id
	Name       string `gorm:"type:varchar(20);not null" json:"Name"`
	Password   string `gorm:"size:255;not null" json:"Password"`
	Telephone  string `gorm:"varchar(11);not null;unique" json:"Telephone"`
	Expiration int8   `gorm:"TINYINT;default:0" json:"Expiration"`
	Number     int   `gorm:"INT;default:0" json:"Number"`
	GiftNumber int   `gorm:"INT;default:0" json:"GiftNumber"`
	GiftId     int    `gorm:"INT;default:5" json:"GiftId"`
}

type Result struct {
	Id         int    `gorm:"not null pk autoincr INT(11)"`
	GiftId     int    `gorm:"comment('奖品ID，关联gift表') INT(11) "`
	GiftName   string `gorm:"comment('奖品名称') VARCHAR(255)"`
	GiftType   int    `gorm:"comment('奖品类型，同gift. gtype') INT(11)"`
	Uid        int    `gorm:"comment('用户ID') INT(11)"`
	UserName   string `gorm:"comment('用户名') VARCHAR(50)"`
	Telephone  string `gorm:"comment('联系方式') VARCHAR(50)"`
	PrizeCode  int    `gorm:"comment('抽奖编号（4位的随机数）') INT(11)"`
	GiftData   string `gorm:"comment('获奖信息') TEXT"`
	SysCreated int    `gorm:"comment('创建时间') INT(11)"`
	SysIp      string `gorm:"comment('用户抽奖的IP') VARCHAR(50)"`
	SysStatus  int    `gorm:"comment('状态，0 正常，1删除，2作弊') INT(11)"`
}

type Gift struct {
	Id           int    `gorm:"not null INT(11)"`
	Title        string `gorm:"not null comment('奖品名称') VARCHAR(255)"`
	PrizeNum     int    `gorm:"not null default -1 comment('奖品数量，0 无限量，>0限量，<0无奖品') INT(10)"`
	LeftNum      int    `gorm:"not null default 0 comment('剩余数量') INT(10)"`
	PrizeCode    string `gorm:"not null comment('0-9999表示100%，0-0表示万分之一的中奖概率') VARCHAR(255)"`
	PrizeTime    int    `gorm:"not null default 0 comment('发奖周期，D天') INT(10)"`
	Img          string `gorm:"comment('奖品图片') TEXT(0)"`
	DisplayOrder int    `gorm:"comment('位置序号，小的排在前面') INT(255)"`
	Gtype        int    `gorm:"not null comment('奖品类型，0 虚拟币，1 虚拟券，2 实物-小奖，3 实物-大奖'') INT(10)"`
	Gdata        string `gorm:"comment('扩展数据，如：虚拟币数量') TEXT"`
	TimeBegin    int    `gorm:"not null comment('开始时间') INT(13)"`
	TimeEnd      int    `gorm:"not null comment('结束时间') INT(13)"`
	PrizeData    string `gorm:"comment('发奖计划，[[时间1,数量1],[时间2,数量2]]') TEXT"`
	PrizeBegin   int    `gorm:"comment('发奖计划周期的开始') INT(13)"`
	PrizeEnd     int    `gorm:"comment('发奖计划周期的结束') INT(13)"`
	SysStatus    int    `gorm:"not null default 0 comment('状态，0 正常，1 删除') INT(13)"`
	SysCreated   int    `gorm:"comment('创建时间') INT(13)"`
	SysUpdated   int    `gorm:"comment('修改时间') INT(13)"`
	SysIp        string `gorm:"comment('操作人IP') VARCHAR(50)"`
	Result       Result `gorm:"foreignKey:GiftId"`
}
