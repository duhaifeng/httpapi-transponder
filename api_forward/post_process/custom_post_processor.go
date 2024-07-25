/**
 * 一个自定义后处理器，如果后端service返回的内容需要复杂的格式转换，导致无法通过简单配置就实现，可以像这样硬编码。
 *（本文件代码纯粹实例，无需参考）
 * @author duhaifeng
 * @date   2021/07/01
 */
package post_process

import (
	"cv-api-gw/api_forward/json_convert"
	"cv-api-gw/common"
	"fmt"
	"strconv"
	"strings"

	api "github.com/duhaifeng/simpleapi"
)

func CustomPostProcess(r *api.Request, mappedFlatData map[string]map[string]interface{}) map[string]interface{} {
	groupFlatMap := make(map[string]interface{}) //避免没有赋值成功导致空指针错误，声明时直接实例化
	for _, flatData := range mappedFlatData {
		groupFlatMap = flatData
	}

	groupTotalCount := 0
	for i := 0; ; i++ {
		_, ok := groupFlatMap[fmt.Sprintf("imageGroups.[%d].imageGroupID", i)]
		if !ok {
			groupTotalCount = i //这里不用i+1，刚好是在i+1时break的
			break
		}
	}

	pagiInfo := common.CalcPagiInfo(r.GetUrlParam("pageNo"), r.GetUrlParam("pageSize"), groupTotalCount)
	startIndex := pagiInfo["startIndex"]
	endIndex := pagiInfo["endIndex"]

	groupPagiMap := make(map[string]interface{})
	for k, v := range groupFlatMap {
		if !strings.HasPrefix(k, "imageGroups.[") {
			//非device list中的字段，直接原封不动保留
			groupPagiMap[k] = v
			continue
		}
		deviceIndexStr := strings.Replace(k, "imageGroups.[", "", -1)
		deviceIndexPostfix := deviceIndexStr[strings.Index(deviceIndexStr, "]"):]
		deviceIndexStr = deviceIndexStr[:strings.Index(deviceIndexStr, "]")]
		groupIndex, err := strconv.Atoi(deviceIndexStr)
		if err != nil {
			common.Log.Error("group index pagi calc error: %s, %s", k, err.Error())
			continue
		}
		if groupIndex >= startIndex && groupIndex <= endIndex {
			newKey := fmt.Sprintf("imageGroups.[%d%s", groupIndex-startIndex, deviceIndexPostfix)
			groupPagiMap[newKey] = v
		}
	}
	delete(pagiInfo, "startIndex")
	delete(pagiInfo, "endIndex")

	elemRebuilder := new(json_convert.CachedElementRebuilder)
	for mappedKey, value := range groupPagiMap {
		elemRebuilder.InsertData(mappedKey, value)
	}
	for k, v := range pagiInfo {
		elemRebuilder.InsertData(fmt.Sprintf("pagination.%s", k), v)
	}
	return elemRebuilder.GetBuiltData()
}
