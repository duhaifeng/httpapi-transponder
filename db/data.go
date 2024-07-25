/**
 * 数据库所用数据结构的声明
 * @author duhaifeng
 * @date   2021/04/15
 */
package db

import "time"

type DbDataOrmBase struct {
	CreateUser int
	CreateTime time.Time
	UpdateUser int
	UpdateTime time.Time
}

type UserList []*UserEntry

/**
 * 网关用户
 */
type UserEntry struct {
	UserId       string `gorm:"primary_key:yes"`
	UserName     string
	UserPassword string
	UserType     int
	TokenList    []GatewayToken `gorm:"foreignKey:UserId;references:UserId"`
	AkskList     []GatewayAksk  `gorm:"foreignKey:UserId;references:UserId"`
	Remark       string
	DbDataOrmBase
}

func (*UserEntry) TableName() string {
	return "user"
}

type TokenList []*GatewayToken

type GatewayToken struct {
	Id          int `gorm:"primary_key:yes"`
	UserId      string
	AccessToken string
	FreshToken  string
	DbDataOrmBase
}

func (*GatewayToken) TableName() string {
	return "user_token"
}

type AkskList []*GatewayAksk

type GatewayAksk struct {
	Id        int `gorm:"primary_key:yes"`
	UserId    string
	AccessKey string
	SecureKey string
	DbDataOrmBase
}

func (*GatewayAksk) TableName() string {
	return "user_aksk"
}
