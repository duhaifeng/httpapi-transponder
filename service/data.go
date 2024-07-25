/**
 *
 * @author duhaifeng
 * @date   2021/04/14
 */
package service

import "time"

type ServiceDataBase struct {
	CreateUser int       `json:"-"`
	CreateTime time.Time `json:"createTime"`
	UpdateUser int       `json:"-"`
	UpdateTime time.Time `json:"-"`
}

/**
 * 网关用户
 */
type GatewayUser struct {
	UserId       string `json:"userID"`
	UserName     string `validate:"required" json:"userName"`
	UserPassword string `validate:"required" json:"password"`
	UserType     string `json:"-"`
	Remark       string `json:"remark"`
	ServiceDataBase
}

/**
 * 用户Token（为了兼容盒子版本API）
 */
type GatewayToken struct {
	UserId      string `validate:"required" json:"userID"`
	AccessToken string
	FreshToken  string
	ServiceDataBase
}

/**
 * API调用AKSK（主要用于API调用使用）
 */
type GatewayAksk struct {
	UserId    string `validate:"required" json:"userID"`
	AccessKey string
	SecureKey string
	ServiceDataBase
}

type Device struct {
	DeviceID   string
	DeviceName string `validate:"required"`
}

type ImageGroup struct {
	SerialNo       string
	ImageGroupID   string
	ImageGroupName string `validate:"required" `
	Remark         string
}
