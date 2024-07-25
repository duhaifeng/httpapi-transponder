/**
 * 一个简单的探活接口
 * @author duhaifeng
 * @date   2021/04/14
 */
package handler

import (
	api "github.com/duhaifeng/simpleapi"
	"time"
)

type HealthCheckHandler struct {
	api.BaseHandler
}

func (this *HealthCheckHandler) HandleRequest(r *api.Request) (interface{}, error) {
	rtnMap := make(map[string]interface{})
	rtnMap["healthCheck"] = "Ok"
	rtnMap["time"] = time.Now()
	return rtnMap, nil
}
