/**
 * User相关数据库操作
 * @author duhaifeng
 * @date   2021/04/15
 */
package db

import (
	"cv-api-gw/common"
	"cv-api-gw/common/busierr"

	api "github.com/duhaifeng/simpleapi"
)

type UserDbOperator struct {
	api.BaseDbOperator
}

func (this *UserDbOperator) AddNewUser(user *UserEntry) error {
	err := this.OrmConn().Create(user).Error
	if err != nil {
		common.Log.Error("[db error] add gateway user failed: %s", err.Error())
		return busierr.WrapGoError(err)
	}
	return nil
}

func (this *UserDbOperator) GetUserList() (UserList, error) {
	var userList UserList
	//连带aksk和token信息一起查询出来（通过Preload）
	err := this.OrmConn().Preload("AkskList").Preload("TokenList").Find(&userList).Error
	if err != nil {
		common.Log.Error("[db error] get user list error: %s", err.Error())
		return nil, busierr.WrapGoError(err)
	}
	return userList, nil
}

func (this *UserDbOperator) GetUserDetail(userId string) (*UserEntry, error) {
	userInfo := new(UserEntry)
	err := this.OrmConn().Where("user_id=?", userId).Find(&userInfo).Error
	if err != nil {
		if err.Error() == "record not found" {
			return nil, nil
		}
		common.Log.Error("[db error] get user detail error: %s", err.Error())
		return nil, busierr.WrapGoError(err)
	}
	if userInfo.UserId == "" {
		return nil, nil
	}
	return userInfo, nil
}

func (this *UserDbOperator) GetUserByName(userName string) (*UserEntry, error) {
	userInfo := new(UserEntry)
	err := this.OrmConn().Where("user_name=?", userName).Find(&userInfo).Error
	if err != nil {
		if err.Error() == "record not found" {
			return nil, nil
		}
		common.Log.Error("[db error] get user by name error: %s", err.Error())
		return nil, busierr.WrapGoError(err)
	}
	if userInfo.UserId == "" {
		return nil, nil
	}
	return userInfo, nil
}

func (this *UserDbOperator) UpdateUser(user *UserEntry) error {
	err := this.OrmConn().Save(user).Error
	if err != nil {
		common.Log.Error("[db error] update user info error: %s", err.Error())
		return busierr.WrapGoError(err)
	}
	return nil
}

func (this *UserDbOperator) DeleteUser(userId string) error {
	userInfo := &UserEntry{UserId: userId}
	err := this.OrmConn().Delete(userInfo).Error
	if err != nil {
		common.Log.Error("[db error] update user info error: %s", err.Error())
		return busierr.WrapGoError(err)
	}
	return nil
}
