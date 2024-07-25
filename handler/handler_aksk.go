/**
 * 负责处理登录AKSK的handler（网关自有Handler）
 * @author duhaifeng
 * @date   2021/06/21
 */
package handler

import (
	"cv-api-gw/service"
	"errors"
	api "github.com/duhaifeng/simpleapi"
)

type CreateAkskHandler struct {
	UserService *service.GatewayUserService
	AkskService *service.GatewayAkskService
	GatewayAksk *service.GatewayAksk
	GatewayHandlerBase
}

func (this *CreateAkskHandler) HandleRequest(r *api.Request) (interface{}, error) {
	userId := r.GetUrlVar("userID")
	if userId == "" {
		return nil, errors.New("user id is necessary")
	}
	if !isAdminUser(this.GetUser()) && !isSelfUser(userId, this.GetUser()) {
		return nil, errors.New("can not manage other user's aksk")
	}
	user, err := this.UserService.GetUserDetailById(userId)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, errors.New("user does not exist")
	}
	newAksk, err := this.AkskService.CreateAksk(user)
	if err != nil {
		return nil, err
	}
	rtnMap := make(map[string]interface{})
	rtnMap["userAksk"] = newAksk
	return rtnMap, nil
}

type AkskDetailHandler struct {
	AkskService *service.GatewayAkskService
	GatewayHandlerBase
}

func (this *AkskDetailHandler) HandleRequest(r *api.Request) (interface{}, error) {
	userId := r.GetUrlVar("userID")
	if userId == "" {
		return nil, errors.New("user id is necessary")
	}
	if !isAdminUser(this.GetUser()) && !isSelfUser(userId, this.GetUser()) {
		return nil, errors.New("can not manage other user's aksk")
	}
	aksk, err := this.AkskService.GetUserAksk(userId)
	if err != nil {
		return nil, err
	}
	rtnMap := make(map[string]interface{})
	rtnMap["userAksk"] = aksk
	return rtnMap, nil
}

type AkskDeleteHandler struct {
	UserService *service.GatewayUserService
	AkskService *service.GatewayAkskService
	GatewayHandlerBase
}

func (this *AkskDeleteHandler) HandleRequest(r *api.Request) (interface{}, error) {
	userId := r.GetUrlVar("userID")
	if userId == "" {
		return nil, errors.New("user id is necessary")
	}
	if !isAdminUser(this.GetUser()) && !isSelfUser(userId, this.GetUser()) {
		return nil, errors.New("can not manage other user's aksk")
	}
	user, err := this.UserService.GetUserDetailById(userId)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, errors.New("user does not exist")
	}
	return nil, this.AkskService.DeleteAksk(userId)
}
