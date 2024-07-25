/**
 *
 * @author duhaifeng
 * @date   2021/04/14
 */
package interceptor

import (
	"cv-api-gw/service"
	api "github.com/duhaifeng/simpleapi"
)

type ValidateAccessAuthInterceptor struct {
	api.Interceptor
}

func (this *ValidateAccessAuthInterceptor) HandleRequest(r *api.Request) (interface{}, error) {
	accessToken := r.GetHeader("accessToken")
	accessKey := r.GetHeader("accessKey")
	//secureKey := r.GetHeader("secureKey")
	var userInfo *service.GatewayUser
	if accessToken != "" {
		//TODO:如果User模块从网关剥离，那么这里要改为从user api获取
		userInfo = service.AccessAuthManager.GetTokenUser(accessToken)
	} else if accessKey != "" {
		userInfo = service.AccessAuthManager.GetAkskUser(accessKey)
	}
	if userInfo == nil {
		//return nil, errors.New("gateway authentication failed: token or aksk corresponding user doesn't exist")
		//TODO：为了方便测试，临时设定userID为admin
		userInfo = new(service.GatewayUser)
		userInfo.UserId = "admin"
		userInfo.UserName = "admin"
		r.SetHeader("userID", "admin")
		this.GetContext().SetAttachment("userID", "admin")
		this.GetContext().SetAttachment("userInfo", userInfo)
	} else {
		this.GetContext().SetAttachment("userID", userInfo.UserId)
		this.GetContext().SetAttachment("userInfo", userInfo)
	}
	data, err := this.CallNextProcess(r)
	return data, err
}
