/**
 *
 * @author duhaifeng
 * @date   2021/04/14
 */
package interceptor

import (
	"cv-api-gw/common"
	"encoding/json"
	"net/http"
	"reflect"
	"time"

	api "github.com/duhaifeng/simpleapi"
)

func isInterfaceNil(i interface{}) bool {
	if i == nil {
		return true
	}
	//go中直接对interface{}判空会失效，因此需要借助反射
	switch reflect.TypeOf(i).Kind() {
	case reflect.Ptr, reflect.Map, reflect.Array, reflect.Chan, reflect.Slice:
		return reflect.ValueOf(i).IsNil()
	}
	return false
}

/**
 * 将响应结果格式化后输出的拦截器
 */
type FormatResponseInterceptor struct {
	api.Interceptor
}

func (this *FormatResponseInterceptor) HandleRequest(r *api.Request) (interface{}, error) {
	w := this.GetResponse()
	w.AlreadyResponsed()
	w.SetHeader("Content-Type", "Application/json")

	respData := make(map[string]interface{})
	startTime := time.Now()
	handlerRtnData, err := this.CallNextProcess(r)

	//增加跨域访问支持
	this.GetResponse().SetHeader("Access-Control-Allow-Origin", "*")
	this.GetResponse().SetHeader("Access-Control-Request-Method", "POST,GET,OPTIONS,PUT,DELETE")
	this.GetResponse().SetHeader("Access-Control-Allow-Methods", "POST,GET,OPTIONS,PUT,DELETE")
	this.GetResponse().SetHeader("Access-Control-Allow-Headers", "*")

	//如果Handler指示了不要包装响应内容，则直接json化返回
	wrapRespInter := this.GetResponse().GetResponseData("_wrap_response_")
	if wrapRespInter != nil {
		wrapResp, ok := wrapRespInter.(bool)
		if ok && !wrapResp {
			if handlerRtnData != nil && !isInterfaceNil(handlerRtnData) {
				jsonBytes, _ := json.Marshal(handlerRtnData)
				return w.Write(jsonBytes)
			} else {
				return nil, nil
			}
		}
	}

	useTime := time.Since(startTime)
	//如果Handler返回了错误，则直接返回错误信息
	errVal := reflect.ValueOf(err)
	if !errVal.IsNil() {
		respData["code"] = 500
		respData["requestID"] = this.GetContext().GetRequestId()
		respData["message"] = err.Error()
		respData["processTime"] = useTime.Milliseconds()
		//error很多时候是业务抛出来的，并不是真的系统错误，所以不用500，避免Web前端不解析Response
		//TODO：这里还需要补充判断一下panic的异常，这类异常应该返回500
		w.WriteHeader(http.StatusOK)
		jsonBytes, _ := json.Marshal(respData)
		return w.Write(jsonBytes)
	}
	if handlerRtnData == nil || isInterfaceNil(handlerRtnData) {
		respData["code"] = 0
		respData["requestID"] = this.GetContext().GetRequestId()
		respData["message"] = "success"
		respData["processTime"] = useTime.Milliseconds()
		respData["data"] = map[string]interface{}{}
		w.WriteHeader(http.StatusOK)
		jsonBytes, _ := json.Marshal(respData)
		return w.Write(jsonBytes)
	}

	//通过将Handler返回的响应数据序列化再反序列化，将其中的结构体类型信息抹掉，转为Map，并在其中注入requestId
	handlerRtnDataBytes, err := json.Marshal(handlerRtnData)
	if err != nil {
		common.Log.Error("format response error when marshal handler response %s %v", err, handlerRtnData)
		w.WriteHeader(http.StatusOK)
		//如果无法序列化分析结果就原封不动返回
		jsonBytes, _ := json.Marshal(handlerRtnData)
		return w.Write(jsonBytes)
	}

	unmarshalRtnData := make(map[string]interface{})
	err = json.Unmarshal(handlerRtnDataBytes, &unmarshalRtnData)
	if err != nil {
		common.Log.Error("format response error when unmarshal handler response %s %v \n %s", err, handlerRtnData, string(handlerRtnDataBytes))
		respData = make(map[string]interface{})
		respData["code"] = 0
		respData["requestID"] = this.GetContext().GetRequestId()
		respData["message"] = "success"
		respData["processTime"] = useTime.Milliseconds()
		respData["data"] = handlerRtnData
		w.WriteHeader(http.StatusOK)
		jsonBytes, _ := json.Marshal(respData)
		return w.Write(jsonBytes)
	}
	w.WriteHeader(http.StatusOK)
	jsonBytes, _ := json.Marshal(unmarshalRtnData)
	return w.Write(jsonBytes)
}
