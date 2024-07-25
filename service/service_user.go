/**
 *
 * @author duhaifeng
 * @date   2021/04/17
 */
package service

import (
	"cv-api-gw/common"
	"cv-api-gw/db"
	api "github.com/duhaifeng/simpleapi"
	"time"
)

type GatewayUserService struct {
	UserDbOperator *db.UserDbOperator
	api.BaseService
}

func (this *GatewayUserService) RegisterNewUser(newUser *GatewayUser) error {
	dbUserEntry := new(db.UserEntry)
	common.CopyStructField(newUser, dbUserEntry)
	dbUserEntry.CreateTime = time.Now()
	dbUserEntry.UpdateTime = time.Now()
	return this.UserDbOperator.AddNewUser(dbUserEntry)
}

func (this *GatewayUserService) UpdateUser(user *GatewayUser) error {
	dbUserEntry := new(db.UserEntry)
	common.CopyStructField(user, dbUserEntry)
	dbUserEntry.CreateTime = time.Now()
	dbUserEntry.UpdateTime = time.Now()
	return this.UserDbOperator.UpdateUser(dbUserEntry)
}

func (this *GatewayUserService) GetUserList() ([]*GatewayUser, error) {
	dbUserList, err := this.UserDbOperator.GetUserList()
	if err != nil {
		return nil, err
	}
	var userList []*GatewayUser
	for _, dbUser := range dbUserList {
		user := new(GatewayUser)
		common.CopyStructField(dbUser, user)
		user.UserPassword = ""
		userList = append(userList, user)
	}
	return userList, nil
}

func (this *GatewayUserService) GetUserDetailById(userId string) (*GatewayUser, error) {
	dbUserDetail, err := this.UserDbOperator.GetUserDetail(userId)
	if err != nil {
		return nil, err
	}
	if dbUserDetail == nil {
		return nil, nil
	}
	userDetail := new(GatewayUser)
	common.CopyStructField(dbUserDetail, userDetail)
	userDetail.UserPassword = ""
	return userDetail, nil
}

func (this *GatewayUserService) GetUserDetailByName(userName string) (*GatewayUser, error) {
	dbUserDetail, err := this.UserDbOperator.GetUserByName(userName)
	if err != nil {
		return nil, nil
	}
	if dbUserDetail != nil {
		userDetail := new(GatewayUser)
		common.CopyStructField(dbUserDetail, userDetail)
		return userDetail, nil
	}
	return nil, nil
}

func (this *GatewayUserService) DeleteUser(userId string) error {
	err := this.UserDbOperator.DeleteUser(userId)
	if err != nil {
		return err
	}
	return nil
}
