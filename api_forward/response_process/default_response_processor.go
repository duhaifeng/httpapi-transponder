/**
 * 针对一次请求的最后环节处理（只是预留，没有真实参与逻辑）
 * @author duhaifeng
 * @date   2021/07/14
 */
package response_process

import api "github.com/duhaifeng/simpleapi"

var ResponseProcessors = make(map[string]func(*api.Request, map[string]interface{}) map[string]interface{})

func init() {
	ResponseProcessors["DefaultPostProcessor"] = DefaultResponseProcess
}

func GetResponseProcessor(processorName string) func(*api.Request, map[string]interface{}) map[string]interface{} {
	if processorName == "" {
		return DefaultResponseProcess
	}
	postProcessor, ok := ResponseProcessors[processorName]
	if ok {
		return postProcessor
	}
	return DefaultResponseProcess
}

func DefaultResponseProcess(r *api.Request, respData map[string]interface{}) map[string]interface{} {
	return respData
}
