/**
 * 针对Backend响应数据的后处理，对于设备等接口需要对数据进行大尺度变换，需要写这类数据后处理器
 * @author duhaifeng
 * @date   2021/07/01
 */
package post_process

import (
	"cv-api-gw/api_forward/json_convert"

	api "github.com/duhaifeng/simpleapi"
)

var PostProcessors = make(map[string]func(*api.Request, map[string]map[string]interface{}) map[string]interface{})

func init() {
	PostProcessors["DefaultPostProcessor"] = DefaultPostProcess
	PostProcessors["GroupListPostProcessor"] = CustomPostProcess
}

func GetPostProcessor(processorName string) func(*api.Request, map[string]map[string]interface{}) map[string]interface{} {
	if processorName == "" {
		return DefaultPostProcess
	}
	postProcessor, ok := PostProcessors[processorName]
	if ok {
		return postProcessor
	}
	return DefaultPostProcess
}

/**
 * 默认动作只是简单地把各请求数据合并到一起（即注册到同一个Rebuilder中进行重建）
 */
func DefaultPostProcess(r *api.Request, mappedFlatData map[string]map[string]interface{}) map[string]interface{} {
	//基于遍历后的扁平数据映射重建为层级结构
	elemRebuilder := new(json_convert.CachedElementRebuilder)
	for _, flatData := range mappedFlatData {
		for mappedKey, value := range flatData {
			elemRebuilder.InsertData(mappedKey, value)
		}
	}
	return elemRebuilder.GetBuiltData()
}
