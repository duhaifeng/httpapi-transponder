/**
 * 请求Backend时发生Error的处理器
 * @author duhaifeng
 * @date   2021/08/19
 */
package error_process

import (
	api "github.com/duhaifeng/simpleapi"
)

var ErrorProcessors = make(map[string]func(*api.Request, map[string]map[string]interface{}, error) error)

func init() {
	ErrorProcessors["DefaultErrorProcessor"] = DefaultErrorProcess
}

func GetErrorProcessor(processorName string) func(*api.Request, map[string]map[string]interface{}, error) error {
	if processorName == "" {
		return DefaultErrorProcess
	}
	errorProcessor, ok := ErrorProcessors[processorName]
	if ok {
		return errorProcessor
	}
	return DefaultErrorProcess
}

/**
 *
 */
func DefaultErrorProcess(r *api.Request, mappedFlatData map[string]map[string]interface{}, err error) error {
	return nil
}
