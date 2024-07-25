/**
 * Backend请求转发入口，负责编排数据验证、body格式转换、错误处理、字段格式转换、响应体格式转换等步骤
 * @author duhaifeng
 * @date   2021/07/14
 */
package api_forward

import (
	"cv-api-gw/api_forward/data_validate"
	"cv-api-gw/api_forward/json_convert"
	"cv-api-gw/api_forward/post_process"
	"cv-api-gw/api_forward/response_process"
	"cv-api-gw/common"
	"encoding/json"
	"errors"
	"fmt"
	api "github.com/duhaifeng/simpleapi"
	"strings"
)

type ForwardProxy struct {
}

/**
 * 将URL中引用的Body中数据占位符进行替换
 */
func (this *ForwardProxy) replaceUrlPlaceholder(url string, envValues map[string]string) string {
	//先替换url相关的占位参数，由于url相关的参数较少，因此整体遍历一遍
	for k, v := range envValues {
		url = strings.Replace(url, k, v, -1)
	}
	return url
}

/**
 * 将当前请求的Response加入环境变量，留作下次请求使用
 */
func (this *ForwardProxy) updateEnvValues(envValues map[string]string, backendMappedFlatResp map[string]interface{}) {
	for flatKey, value := range backendMappedFlatResp {
		envValues[fmt.Sprintf("{response:%s}", flatKey)] = fmt.Sprintf("%v", value)
	}
}

/**
 * 打印请求后端的调试信息
 */
func (this *ForwardProxy) printRequestInfo(backendConf *common.EndpointBackend, prefix string, data []byte) {
	if backendConf.DoNotPrintDebug {
		return
	}
	if len(data) > 20 {
		//如果超过了20个字符，则进行一次格式美化
		dataMap := make(map[string]interface{})
		err := json.Unmarshal(data, &dataMap)
		if err == nil {
			beautifulData, err := json.MarshalIndent(dataMap, "    ", "    ")
			if err == nil {
				data = beautifulData
			}
		}
	}
	if len(data) > 2000 {
		//日志限制2K长度
		common.Log.Debug("\n%s : %s", prefix, data[:2000])
	} else {
		common.Log.Debug("\n%s : %s", prefix, data)
	}
}

func (this *ForwardProxy) requestBackend(backendConf *common.EndpointBackend, paramStr string, bodyMap map[string]interface{}, envValues map[string]string) (map[string]interface{}, error) {
	if backendConf.UrlPattern == "" || len(backendConf.Host) == 0 {
		//此种情况说明配置了一个纯粹的假接口，直接返回空数据即可
		return make(map[string]interface{}), nil
	}
	bodyMapToBackend := json_convert.GetMappedFields(bodyMap, envValues, backendConf.BodyFieldOperation, true)
	//将映射后的Body扁平数据重建为层级数据
	bodyMapToBackend = post_process.DefaultPostProcess(nil, map[string]map[string]interface{}{"_": bodyMapToBackend})
	bodyToBackendBytes, err := json.MarshalIndent(bodyMapToBackend, "    ", "    ")
	if err != nil {
		return nil, err
	}

	//替换请求后端URL中对body元素的引用
	backendUrl := backendConf.Host[0] + backendConf.UrlPattern
	if backendConf.AttachUrlParamPart {
		if paramStr != "" {
			backendUrl = fmt.Sprintf("%s?%s", backendUrl, paramStr)
		}
	}
	backendUrl = this.replaceUrlPlaceholder(backendUrl, envValues)
	this.printRequestInfo(backendConf, fmt.Sprintf("request body to backend <%s>%s", backendConf.Method, backendUrl), bodyToBackendBytes)
	_, backendResp, err := common.RequestHttp(backendConf.Method, backendUrl, nil, string(bodyToBackendBytes))
	if err != nil {
		return nil, err
	}
	if backendResp == nil {
		common.Log.Error("backend response nothing: %s", backendUrl)
		return nil, errors.New("backend response nothing: " + backendUrl)
	}
	backendRespData := make(map[string]interface{})
	err = json.Unmarshal(backendResp, &backendRespData)
	if err != nil {
		common.Log.Error("backend response format is not json: %s\n %s", err.Error(), string(backendResp))
		//后端返回的可能不是json，只是一个常规Message，静文的设备查询接口在deviceID不存在时就返回非json消息
		backendRespData["code"] = 500
		backendRespData["message"] = string(backendResp)
	}
	this.printRequestInfo(backendConf, fmt.Sprintf("response from backend <%s>%s", backendConf.Method, backendUrl), backendResp)
	return backendRespData, nil
}

/**
 * 当请求backend发生错误时，则触发错误处理句柄
 */
func (this *ForwardProxy) onBackendError(backendConf *common.EndpointBackend, paramStr string, bodyMap map[string]interface{}, envValues map[string]string) {
	if len(backendConf.OnErrorCall) == 0 {
		return
	}
	for _, errorHandlerConf := range backendConf.OnErrorCall {
		common.Log.Debug("error handler <%s>%s", backendConf.Method, backendConf.UrlPattern)
		//进行错误补偿处理时，则忽略错误
		this.requestBackend(errorHandlerConf, paramStr, bodyMap, envValues)
	}
}

/**
 * 根据配置对Backend接口进行请求，并对字段进行映射
 */
func (this *ForwardProxy) ForwardRequestByConfig(forwardConf *common.ForwardEndpoint, r *api.Request, body []byte, envValues map[string]string) (interface{}, error) {
	paramStr := ""
	uriItems := strings.Split(r.GetOriReq().RequestURI, "?")
	if len(uriItems) >= 2 {
		paramStr = uriItems[len(uriItems)-1]
	}
	dataValidator := data_validate.GetDataValidator(forwardConf.ValidateRule)
	err := dataValidator(r)
	if err != nil {
		return nil, err
	}
	bodyMap := make(map[string]interface{})
	if len(body) != 0 {
		err := json.Unmarshal(body, &bodyMap)
		if err != nil {
			return nil, errors.New("request body ")
		}
	}

	allBackendResp := make(map[string]map[string]interface{})
	for _, backendConf := range forwardConf.Backend {
		backendRespData, err := this.requestBackend(backendConf, paramStr, bodyMap, envValues)
		if err != nil && !forwardConf.OnBackendErrorContinue {
			this.onBackendError(backendConf, paramStr, bodyMap, envValues)
			//默认发生错误时，都直接退出，如果配置了continue，则不退出
			if forwardConf.OnBackendErrorContinue {
				continue
			} else {
				//这里应该err信息返回，否则前端会误以为成功
				backendErrMap := make(map[string]interface{})
				backendErrMap["code"] = 500
				backendErrMap["message"] = err.Error()
				allBackendResp[backendConf.UrlPattern] = backendErrMap
				return post_process.DefaultPostProcess(r, allBackendResp), nil
			}
		}

		//如果后端返回了code字段，则直接透传
		codeFromBackend := 0.0
		if codePtr, ok1 := backendRespData["code"]; ok1 {
			if code, ok2 := codePtr.(float64); ok2 {
				codeFromBackend = code
			}
		}
		if codeFromBackend != 0 {
			//如果后端响应结果code不为0，则停止后续其他请求。（由于请求结果不全，可能导致post processor中无法正常处理分页等逻辑）
			//并且不处理append动作，避免产生一些无用字段到返回结果中。
			backendMappedFlatResp := json_convert.GetMappedFields(backendRespData, envValues, backendConf.ResponseFieldOperation, false)
			if !backendConf.DropResponse {
				allBackendResp[backendConf.UrlPattern] = backendMappedFlatResp
			}
			this.onBackendError(backendConf, paramStr, bodyMap, envValues)
			//默认发生错误时，都直接退出，如果配置了continue，则不退出
			if !forwardConf.OnBackendErrorContinue {
				//如果上面配置了DropResponse，会导致当前这次backend返回的错误信息不会被返回到前端，因此这里补偿一下
				allBackendResp[backendConf.UrlPattern] = backendMappedFlatResp
				return post_process.DefaultPostProcess(r, allBackendResp), nil
			}
		} else {
			responseProcessor := response_process.GetResponseProcessor(backendConf.ResponseProcessor)
			backendMappedFlatResp := responseProcessor(r, backendRespData)
			//为了避免在组合请求中，每一次请求都执行rebuild耗费资源，这里只是先通过Mapper返回扁平数据，以便于post process
			backendMappedFlatResp = json_convert.GetMappedFields(backendMappedFlatResp, envValues, backendConf.ResponseFieldOperation, true)
			//将当前请求返回的结果更新到环境变量中，留作后续请求替换变量使用
			this.updateEnvValues(envValues, backendMappedFlatResp)
			if !backendConf.DropResponse {
				allBackendResp[backendConf.UrlPattern] = backendMappedFlatResp
			}
		}
	}
	//主要在post process中对多次请求的结果进行重建和merge
	postProcessor := post_process.GetPostProcessor(forwardConf.PostProcessor)
	return postProcessor(r, allBackendResp), nil
}
