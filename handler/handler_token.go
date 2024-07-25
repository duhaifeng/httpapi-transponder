/**
 * 管理用户登录Token（网关自有Handler）
 * @author duhaifeng
 * @date   2021/04/14
 */
package handler

import (
	"cv-api-gw/common/busierr"
	"cv-api-gw/service"
	"encoding/json"
	"errors"

	api "github.com/duhaifeng/simpleapi"
)

type CreateTokenHandler struct {
	UserService  *service.GatewayUserService
	TokenService *service.GatewayTokenService
	GatewayToken *service.GatewayToken
	GatewayHandlerBase
}

func (this *CreateTokenHandler) HandleRequest(r *api.Request) (interface{}, error) {
	//TODO:这里是否应该使用userName？
	userId := r.GetUrlVar("userID")
	if userId == "" {
		return nil, errors.New("user id is necessary")
	}
	if !isAdminUser(this.GetUser()) && !isSelfUser(userId, this.GetUser()) {
		return nil, errors.New("can not manage other user's token")
	}
	user, err := this.UserService.GetUserDetailById(userId)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, errors.New("user does not exist")
	}
	newToken, err := this.TokenService.CreateToken(user)
	if err != nil {
		return nil, err
	}
	rtnMap := make(map[string]interface{})
	rtnMap["userToken"] = newToken
	return rtnMap, nil
	return newToken, nil
}

type GetBoxTokenHandler struct {
	UserService         *service.GatewayUserService
	GatewayTokenService *service.GatewayTokenService
	GatewayHandlerBase
}

/**
 * 兼容Box的接口，body传入{username:xxx, password:xxx}，换取一个访问token。
 * 返回：
{
  "exp": 1735689600,
  "access_token": "eyJhbGciOiJIUzI1NiIsImtpZCI6InNpbTIifQ.eyJhdWQiOiJodHRwOi8vYXBpLmV4YW1wbGUuY29tIiwiZXhwIjoxNzM1Njg5NjAwLCJpc3MiOiJodHRwczovL2tyYWtlbmQuaW8iLCJqdGkiOiJtbmIyM3Zjc3J0NzU2eXVpb21uYnZjeDk4ZXJ0eXVpb3AiLCJyb2xlcyI6WyJyb2xlX2EiLCJyb2xlX2IiXSwic3ViIjoiMTIzNDU2Nzg5MHF3ZXJ0eXVpbyJ9.htgbhantGcv6zrN1i43Rl58q1sokh3lzuFgzfenI0Rk",
  "refresh_token": "eyJhbGciOiJIUzI1NiIsImtpZCI6InNpbTIifQ.eyJhdWQiOiJodHRwOi8vYXBpLmV4YW1wbGUuY29tIiwiZXhwIjoxNzM1Njg5NjAwLCJpc3MiOiJodHRwczovL2tyYWtlbmQuaW8iLCJqdGkiOiJtbmIyM3Zjc3J0NzU2eXVpb21uMTI4NzZidmN4OThlcnR5dWlvcCIsInN1YiI6IjEyMzQ1Njc4OTBxd2VydHl1aW8ifQ.4v36tuYHe4E9gCVO-_asuXfzSzoJdoR0NJfVQdVKidw"
}
*/
func (this *GetBoxTokenHandler) HandleRequest(r *api.Request) (interface{}, error) {
	rtnMap := make(map[string]interface{})
	bodyBytes, err := r.GetBody()
	if err != nil {
		return nil, busierr.WrapGoError(err)
	}
	bodyMap := make(map[string]string)
	err = json.Unmarshal(bodyBytes, &bodyMap)
	if err != nil {
		return nil, busierr.WrapGoError(err)
	}
	username, ok := bodyMap["username"]
	if !ok || username == "" {
		return nil, errors.New("username can not empty")
	}
	password, ok := bodyMap["password"]
	if !ok || password == "" {
		return nil, errors.New("password can not empty")
	}
	//userDetail, err := this.UserService.GetUserDetailByName(username)
	//if err != nil {
	//	return nil, err
	//}
	//if userDetail == nil {
	//	return nil, errors.New("user does not exist")
	//}
	//
	//passCipher := fmt.Sprintf("%x", md5.Sum([]byte(password)))
	//if passCipher != userDetail.UserPassword {
	//	return nil, errors.New("password is incorrect")
	//}
	if username == "admin" && password == "admin" {
		rtnMap["exp"] = 173568960
		rtnMap["access_token"] = "eyJhbGciOiJIUzI1NiIsImtpZCI6InNpbTIifQ.eyJhdWQiOiJodHRwOi8vYXBpLmV4YW1wbGUuY29tIiwiZXhwIjoxNzM1Njg5NjAwLCJpc3MiOiJodHRwczovL2tyYWtlbmQuaW8iLCJqdGkiOiJtbmIyM3Zjc3J0NzU2eXVpb21uYnZjeDk4ZXJ0eXVpb3AiLCJyb2xlcyI6WyJyb2xlX2EiLCJyb2xlX2IiXSwic3ViIjoiMTIzNDU2Nzg5MHF3ZXJ0eXVpbyJ9.htgbhantGcv6zrN1i43Rl58q1sokh3lzuFgzfenI0Rk"
		rtnMap["refresh_token"] = "eyJhbGciOiJIUzI1NiIsImtpZCI6InNpbTIifQ.eyJhdWQiOiJodHRwOi8vYXBpLmV4YW1wbGUuY29tIiwiZXhwIjoxNzM1Njg5NjAwLCJpc3MiOiJodHRwczovL2tyYWtlbmQuaW8iLCJqdGkiOiJtbmIyM3Zjc3J0NzU2eXVpb21uMTI4NzZidmN4OThlcnR5dWlvcCIsInN1YiI6IjEyMzQ1Njc4OTBxd2VydHl1aW8ifQ.4v36tuYHe4E9gCVO-_asuXfzSzoJdoR0NJfVQdVKidw"
		this.DoNotWrapResponse()
	} else {
		return nil, errors.New("username or password incorrect")
	}

	return rtnMap, nil
}

type TokenDetailHandler struct {
	GatewayTokenService *service.GatewayTokenService
	GatewayHandlerBase
}

func (this *TokenDetailHandler) HandleRequest(r *api.Request) (interface{}, error) {
	userId := r.GetUrlVar("userID")
	if userId == "" {
		return nil, errors.New("user id is necessary")
	}
	if !isAdminUser(this.GetUser()) && !isSelfUser(userId, this.GetUser()) {
		return nil, errors.New("can not manage other user's token")
	}
	token, err := this.GatewayTokenService.GetUserToken(userId)
	if err != nil {
		return nil, err
	}
	rtnMap := make(map[string]interface{})
	rtnMap["userToken"] = token
	return rtnMap, nil
}

type TokenDeleteHandler struct {
	GatewayTokenService *service.GatewayTokenService
	GatewayHandlerBase
}

func (this *TokenDeleteHandler) HandleRequest(r *api.Request) (interface{}, error) {
	userId := r.GetUrlVar("userID")
	if userId == "" {
		return nil, errors.New("user id is necessary")
	}
	if !isAdminUser(this.GetUser()) && !isSelfUser(userId, this.GetUser()) {
		return nil, errors.New("can not manage other user's token")
	}
	return nil, this.GatewayTokenService.DeleteToken(userId)
}
