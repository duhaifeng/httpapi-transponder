/**
 * 自动转发的入口HTTPHandler，所有通过配置文件自动配置的路由请求均转到AutoForwardHandler，由AutoForwardHandler根据配置文件向Backend转发
 * @author duhaifeng
 * @date   2021/04/14
 */
package handler

import (
	"cv-api-gw/api_forward"
	"cv-api-gw/common"
	"cv-api-gw/common/busierr"
	"strings"

	api "github.com/duhaifeng/simpleapi"
)

type AutoForwardHandler struct {
	GatewayHandlerBase
}

func (this *AutoForwardHandler) HandleRequest(r *api.Request) (interface{}, error) {
	body, err := r.GetBody()
	if err != nil {
		return nil, busierr.WrapGoError(err)
	}
	url := r.GetUrl().String()
	method := r.GetMethod()
	uri := strings.Split(url, "?")[0]
	forwardConf := common.GetApiForwardConf(uri, method)
	if forwardConf == nil {
		_, resp, err := common.RequestHttp(method, url, nil, string(body))
		return resp, err
	}

	forwardProxy := new(api_forward.ForwardProxy)
	backendRespMap, err := forwardProxy.ForwardRequestByConfig(forwardConf, r, body, this.GetEnvValues())
	if err != nil {
		return nil, err
	}
	if forwardConf.DoNotWrapResponse {
		this.DoNotWrapResponse()
	}
	this.SetBusinessCode("0")
	this.SetMessage("success")
	return backendRespMap, nil
}
