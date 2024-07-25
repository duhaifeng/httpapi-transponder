/**
 * handler基类，提供一些网关业务相关的共通函数
 * @author duhaifeng
 * @date   2021/04/14
 */
package handler

import (
	"cv-api-gw/service"
	"fmt"
	api "github.com/duhaifeng/simpleapi"
)

type GatewayHandlerBase struct {
	api.BaseHandler
}

func (this *GatewayHandlerBase) SetBusinessCode(data interface{}) {
	this.GetResponse().SetResponseData("code", data)
}

func (this *GatewayHandlerBase) SetMessage(message string) {
	this.GetResponse().SetResponseData("message", message)
}

func (this *GatewayHandlerBase) DoNotWrapResponse() {
	this.GetResponse().SetResponseData("_wrap_response_", false)
}

func (this *GatewayHandlerBase) GetUserId() string {
	userId, ok := this.GetContext().GetAttachment("userID")
	if !ok {
		return ""
	}
	return fmt.Sprintf("%v", userId)
}

func (this *GatewayHandlerBase) GetUser() *service.GatewayUser {
	user, ok := this.GetContext().GetAttachment("userInfo")
	if !ok {
		return nil
	}
	return user.(*service.GatewayUser)
}

func (this *GatewayHandlerBase) GetEnvValues() map[string]string {
	envValues, _ := this.GetContext().GetAttachment("envValues")
	return envValues.(map[string]string)
}
