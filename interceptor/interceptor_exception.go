/**
 *
 * @author duhaifeng
 * @date   2021/04/14
 */
package interceptor

import (
	"cv-api-gw/common"
	"cv-api-gw/common/busierr"
	"reflect"
	"runtime/debug"

	api "github.com/duhaifeng/simpleapi"
)

/**
 * 异常检测处理拦截器
 */
type ExceptionInterceptor struct {
	api.Interceptor
}

func (this *ExceptionInterceptor) HandleRequest(r *api.Request) (interface{}, error) {
	var err error
	defer func() {
		recoverErr := recover()
		body, _ := r.GetBody()
		if recoverErr != nil {
			common.Log.Error("catch exception:\n request url: %s \n body: %s\n err: %s \n stack: %s", r.GetUrl().String(), string(body), recoverErr, string(debug.Stack()))
		}
	}()
	data, err := this.CallNextProcess(r)
	errVal := reflect.ValueOf(err)
	if !errVal.IsNil() {
		//参数校验失败的错误触发会很多，不打印堆栈
		if _, ok := err.(*busierr.ValidationError); ok {
			common.Log.Error("validation error: %s ", err.Error())
		} else if stackErr, ok := err.(*busierr.StackError); ok {
			body, _ := r.GetBody()
			common.Log.Error("handler return error: \n request url: %s \n body: %s\n err: %s \n stack: %s", r.GetUrl().String(), string(body), err.Error(), stackErr.Stack().CallStack)
		} else {
			body, _ := r.GetBody()
			common.Log.Error("handler return error: \n request url: %s \n body: %s\n err: %s \n stack: %s", r.GetUrl().String(), string(body), err.Error(), string(debug.Stack()))
		}
	}
	return data, err
}
