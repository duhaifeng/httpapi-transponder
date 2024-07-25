/**
 * 计算分页的工具
 * @author duhaifeng
 * @date   2021/07/19
 */
package common

import "strconv"

/**
 * 手工计算分页信息
 */
func CalcPagiInfo(pageNoStr, pageSizeStr string, totalCount int) map[string]int {
	if pageNoStr == "" {
		pageNoStr = "1"
	}
	pageNo, err := strconv.Atoi(pageNoStr)
	if err != nil {
		pageNo = 1
	}
	if pageNo < 1 {
		pageNo = 1
	}

	if pageSizeStr == "" {
		pageSizeStr = "10"
	}
	pageSize, err := strconv.Atoi(pageSizeStr)
	if err != nil {
		pageSize = 10
	}
	if pageSize > 1000 {
		pageSize = 1000
	}

	pageCount := totalCount / pageSize
	//如果有多余不够一页的数据，则页数加1
	if totalCount%pageSize != 0 {
		pageCount += 1
	}
	//如果数据条数为0，总页数也得返回1？
	//if pageCount == 0 {
	//	pageCount = 1
	//}
	//传入的页号不能超过总页数
	if pageNo > pageCount {
		pageNo = pageCount
	}

	startIndex := (pageNo - 1) * pageSize
	endIndex := startIndex + pageSize - 1
	//对于最后一页可能数据不足，需要机遇totalCount进行修正
	if endIndex >= totalCount {
		endIndex = totalCount - 1
	}

	pagination := make(map[string]int)
	pagination["pageNo"] = pageNo
	pagination["pageSize"] = pageSize
	pagination["pageCount"] = pageCount
	pagination["totalCount"] = totalCount
	pagination["startIndex"] = startIndex
	pagination["endIndex"] = endIndex

	return pagination
}
