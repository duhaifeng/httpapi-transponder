/**
 * 管理网关用户的Handler（网关自有Handler）
 * @author duhaifeng
 * @date   2021/04/14
 */
package handler

import (
	"crypto/md5"
	"cv-api-gw/common"
	"cv-api-gw/service"
	"errors"
	"fmt"
	api "github.com/duhaifeng/simpleapi"
	"github.com/google/uuid"
	"strings"
)

func isAdminUser(user *service.GatewayUser) bool {
	if user == nil {
		return false
	}
	if user.UserName != "admin" {
		return false
	}
	return true
}

func isSelfUser(userId string, user *service.GatewayUser) bool {
	if user == nil {
		return false
	}
	if user.UserId != userId {
		return false
	}
	return true
}

type AddUserHandler struct {
	UserService *service.GatewayUserService
	GatewayUser *service.GatewayUser
	GatewayHandlerBase
}

func (this *AddUserHandler) HandleRequest(r *api.Request) (interface{}, error) {
	//只有admin用户才能调用本方法
	if !isAdminUser(this.GetUser()) {
		return nil, errors.New("only admin can manage user")
	}
	err := common.ValidateStruct(this.GatewayUser)
	if err != nil {
		return nil, err
	}
	//用户如果没有指定ID，则创建一个
	if this.GatewayUser.UserId == "" {
		this.GatewayUser.UserId = strings.Replace(uuid.New().String(), "-", "", -1)
	}
	this.GatewayUser.UserPassword = fmt.Sprintf("%x", md5.Sum([]byte(this.GatewayUser.UserPassword)))
	userInfo, err := this.UserService.GetUserDetailById(this.GatewayUser.UserId)
	if err != nil {
		return nil, err
	}
	if userInfo != nil {
		return nil, errors.New("user id already exist: " + this.GatewayUser.UserId)
	}
	userDetail, err := this.UserService.GetUserDetailByName(this.GatewayUser.UserName)
	if err != nil {
		return nil, err
	}
	if userDetail != nil {
		return nil, errors.New("user name already exist: " + this.GatewayUser.UserName)
	}
	err = this.UserService.RegisterNewUser(this.GatewayUser)
	if err != nil {
		//TODO:定义User应用的异常返回Code
		return nil, err
	}
	//TODO：在VPaaS中创建一个对应brand，同时DB中增加一个vpaas_brand字段
	userInfo, err = this.UserService.GetUserDetailById(this.GatewayUser.UserId)
	if err != nil {
		return nil, err
	}
	rtnMap := make(map[string]interface{})
	rtnMap["user"] = userInfo
	return rtnMap, nil
}

type UpdateUserHandler struct {
	UserService *service.GatewayUserService
	GatewayUser *service.GatewayUser
	GatewayHandlerBase
}

func (this *UpdateUserHandler) HandleRequest(r *api.Request) (interface{}, error) {
	if !isAdminUser(this.GetUser()) {
		return nil, errors.New("only admin can manage user")
	}
	err := common.ValidateStruct(this.GatewayUser)
	if err != nil {
		return nil, err
	}
	err = this.UserService.UpdateUser(this.GatewayUser)
	if err != nil {
		return nil, err
	}
	userInfo, err := this.UserService.GetUserDetailById(this.GatewayUser.UserId)
	if err != nil {
		return nil, err
	}
	rtnMap := make(map[string]interface{})
	rtnMap["user"] = userInfo
	return rtnMap, nil
}

type UserListHandler struct {
	UserService *service.GatewayUserService
	GatewayHandlerBase
}

func (this *UserListHandler) getUserPagiList(r *api.Request, userList []*service.GatewayUser) ([]*service.GatewayUser, map[string]int) {
	pagiInfo := common.CalcPagiInfo(r.GetUrlParam("pageNo"), r.GetUrlParam("pageSize"), len(userList))
	startIndex := pagiInfo["startIndex"]
	endIndex := pagiInfo["endIndex"]
	var userPagiList []*service.GatewayUser
	for i, userInfo := range userList {
		if i >= startIndex && i <= endIndex {
			userPagiList = append(userPagiList, userInfo)
		}
	}
	delete(pagiInfo, "startIndex")
	delete(pagiInfo, "endIndex")
	return userPagiList, pagiInfo
}

func (this *UserListHandler) HandleRequest(r *api.Request) (interface{}, error) {
	if !isAdminUser(this.GetUser()) {
		return nil, errors.New("only admin can manage user")
	}
	userList, err := this.UserService.GetUserList()
	if err != nil {
		return nil, err
	}
	rtnMap := make(map[string]interface{})
	userPagiList, pagiInfo := this.getUserPagiList(r, userList)
	rtnMap["userList"] = userPagiList
	rtnMap["pagination"] = pagiInfo
	return rtnMap, nil
}

type UserDetailHandler struct {
	UserService *service.GatewayUserService
	GatewayHandlerBase
}

func (this *UserDetailHandler) HandleRequest(r *api.Request) (interface{}, error) {
	if !isAdminUser(this.GetUser()) {
		return nil, errors.New("only admin can manage user")
	}
	userId := r.GetUrlVar("userID")
	if userId == "" {
		return nil, errors.New("user id is necessary")
	}
	userInfo, err := this.UserService.GetUserDetailById(userId)
	if err != nil {
		return nil, err
	}
	rtnMap := make(map[string]interface{})
	rtnMap["user"] = userInfo
	return rtnMap, nil
}

type UserDeleteHandler struct {
	UserService *service.GatewayUserService
	GatewayHandlerBase
}

func (this *UserDeleteHandler) HandleRequest(r *api.Request) (interface{}, error) {
	if !isAdminUser(this.GetUser()) {
		return nil, errors.New("only admin can manage user")
	}
	userId := r.GetUrlVar("userID")
	if userId == "" {
		return nil, errors.New("user id is necessary")
	}
	return nil, this.UserService.DeleteUser(userId)
}
