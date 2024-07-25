/**
 *
 * @author duhaifeng
 * @date   2021/04/14
 */
package interceptor

import (
	"cv-api-gw/common"

	api "github.com/duhaifeng/simpleapi"
)

/**
 * 请求ID缓存拦截器，一次请求对应产生一个Goroutine，将该Goroutine对应的RequestID进行全局缓存（通过Goroutine No关联），
 * 以方便日志输出中包含该Request ID
 */
type RequestIdInterceptor struct {
	api.Interceptor
}

func (this *RequestIdInterceptor) HandleRequest(r *api.Request) (interface{}, error) {
	common.RequestIdHolder.PutRoutineReqId(this.GetContext().GetRequestId())
	data, err := this.CallNextProcess(r)
	common.RequestIdHolder.DelRoutineReqId()
	return data, err
}
