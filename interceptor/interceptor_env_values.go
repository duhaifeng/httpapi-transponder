/**
 *
 * @author duhaifeng
 * @date   2021/04/14
 */
package interceptor

import (
	"cv-api-gw/api_forward/json_convert"
	"cv-api-gw/common"
	"cv-api-gw/common/busierr"
	"encoding/json"
	"fmt"
	"strings"

	api "github.com/duhaifeng/simpleapi"
)

/**
 * 将url中的变量和参数、以及userID的元素抽取出来，作为环境变量使用
 * @author duhaifeng
 * @date   2021/06/22
 */
type EnvValuesInterceptor struct {
	api.Interceptor
}

func (this *EnvValuesInterceptor) extractUrlVars(r *api.Request, envValues map[string]string) {
	urlItems := strings.Split(r.GetUrl().String(), "/")
	for i, urlItem := range urlItems {
		if urlItem == "" {
			continue
		}
		envValues[fmt.Sprintf("{url:path.[%d]}", i)] = urlItem
	}
}

func (this *EnvValuesInterceptor) extractUrlParams(r *api.Request, envValues map[string]string) {
	for k, v := range r.GetOriReq().URL.Query() {
		envValues[fmt.Sprintf("{url:param.%s}", k)] = v[0]
	}
}

type ElementRecorder struct {
	envValues map[string]string
}

func (this *ElementRecorder) SetEnvValues(envValues map[string]string) {
	this.envValues = envValues
}
func (this *ElementRecorder) Process(hierarchicKeys []string, data interface{}) {
	this.envValues[fmt.Sprintf("{body:%s}", strings.Join(hierarchicKeys, "."))] = fmt.Sprintf("%v", data)
}

func (this *EnvValuesInterceptor) extractBodyItems(r *api.Request, envValues map[string]string) {
	body, err := r.GetBody()
	if err != nil {
		common.Log.Error("extract env values from body error: %s", err.Error())
		return
	}
	bodyMap := make(map[string]interface{})
	if len(body) != 0 {
		err := json.Unmarshal(body, &bodyMap)
		if err != nil {
			common.Log.Error("extract env values from body error: %s", err.Error())
			return
		}
	}
	elemRecursor := new(json_convert.ElementRecursor)
	elemRecorder := new(ElementRecorder)
	elemRecorder.SetEnvValues(envValues)
	elemRecursor.SetElementProcessor(elemRecorder)
	elemRecursor.RecurseElement(bodyMap)
}

func (this *EnvValuesInterceptor) HandleRequest(r *api.Request) (interface{}, error) {
	envValues := make(map[string]string)
	this.extractUrlParams(r, envValues)
	this.extractUrlVars(r, envValues)
	this.extractBodyItems(r, envValues)
	envValues["{env:userID}"] = r.GetHeader("userID")
	this.GetContext().SetAttachment("envValues", envValues)
	data, err := this.CallNextProcess(r)
	return data, busierr.WrapGoError(err)
}
