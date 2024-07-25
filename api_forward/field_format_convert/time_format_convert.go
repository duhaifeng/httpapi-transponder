/**
 * 转换字段时间格式的转换器
 * @author duhaifeng
 * @date   2021/07/16
 */
package field_format_convert

import (
	"fmt"
	"strings"
	"time"
)

var RFC3339MsLayout = "2006-01-02T15:04:05.999Z07:00"
var timeZone *time.Location
var genesisTime time.Time

func init() {
	timeZone, _ = time.LoadLocation("Asia/Shanghai")
	genesisTime = time.Date(1970, 1, 1, 0, 0, 0, 0, timeZone)
}

func timeToRFC3339Ms(datetime time.Time) string {
	RFC3339Str := datetime.Format(RFC3339MsLayout)
	if strings.Contains(RFC3339Str, ".") {
		return RFC3339Str
	}
	//如果没有纳秒时间，那么默认会返回”1970-01-01T00:00:00+08:00“，需要将格式统一为”1970-01-01T00:00:00.000+08:00“形式
	return strings.Replace(RFC3339Str, "+", ".000+", -1)
}

func ConvertDatetimeToGoTime(datetime interface{}) interface{} {
	srcLayout := "2006-01-02T15:04:05Z"
	timeStr, ok := datetime.(string)
	if !ok {
		return timeToRFC3339Ms(genesisTime)
	}
	tim, err := time.ParseInLocation(srcLayout, timeStr, timeZone)
	if err != nil {
		return timeToRFC3339Ms(genesisTime)
	}
	return timeToRFC3339Ms(tim)
}

func ConvertUnixIntToGoTime(datetime interface{}) interface{} {
	timeFloat, ok := datetime.(float64)
	if !ok {
		return genesisTime.Format(RFC3339MsLayout)
	}
	timeInt := int64(timeFloat)
	for {
		if len(fmt.Sprintf("%d", timeInt)) <= 10 {
			break
		}
		//后端返回的Unix时间戳长度可能是13位长毫秒，因此这里统一转换为秒
		timeInt = timeInt / 10
	}
	return timeToRFC3339Ms(time.Unix(timeInt, 0))
}
