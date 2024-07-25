/**
 * 字段格式转换器，对于时间等字段，需要将Unix时间戳形式转换为YMD可读形式输出。
 * @author duhaifeng
 * @date   2021/07/16
 */
package field_format_convert

import (
	"strings"

	"github.com/google/uuid"
)

var FieldFormatConverters = make(map[string]func(interface{}) interface{})

func init() {
	FieldFormatConverters["Time_YmdHmsStr_To_RFC3339Ms"] = ConvertDatetimeToGoTime
	FieldFormatConverters["Time_UnixInt_To_RFC3339Ms"] = ConvertUnixIntToGoTime
	FieldFormatConverters["Create_ID_On_Empty"] = CreateIdOnEmpty
	FieldFormatConverters["Msg_Translate_No_Data"] = TranslateNoDataMsg
	FieldFormatConverters["Stream_Transport_Type_In"] = StreamTransportTypeIn
	FieldFormatConverters["Stream_Transport_Type_Out"] = StreamTransportTypeOut
}

func GetFormatConverter(processorName string) func(interface{}) interface{} {
	if processorName == "" {
		return DefaultFormatConvert
	}
	formatConverter, ok := FieldFormatConverters[processorName]
	if ok {
		return formatConverter
	}
	return DefaultFormatConvert
}

func DefaultFormatConvert(data interface{}) interface{} {
	return data
}

func CreateIdOnEmpty(id interface{}) interface{} {
	if id == nil {
		return strings.Replace(uuid.New().String(), "-", "", -1)
	}
	idStr, ok := id.(string)
	if ok && idStr == "" {
		return strings.Replace(uuid.New().String(), "-", "", -1)
	}
	return id
}

func TranslateNoDataMsg(msg interface{}) interface{} {
	if msg == nil {
		return msg
	}
	msgStr, ok := msg.(string)
	if ok && msgStr == "not find query ID" {
		return "no data found"
	}
	return msg
}

func StreamTransportTypeIn(msg interface{}) interface{} {
	if msg == nil {
		return "rtsp_tcp"
	}
	msgStr, ok := msg.(string)
	if ok && msgStr == "" {
		return "rtsp_tcp"
	}
	msgStr = strings.ToUpper(msgStr)
	if strings.Contains(msgStr, "UDP") {
		return "rtsp_udp"
	}
	return "rtsp_tcp" //StreamProxy模块默认不识别的都将按照TCP对待
}

func StreamTransportTypeOut(msg interface{}) interface{} {
	if msg == nil {
		return "TCP"
	}
	msgStr, ok := msg.(string)
	if ok && msgStr == "" {
		return "TCP"
	}
	msgStr = strings.ToUpper(msgStr)
	if strings.Contains(msgStr, "UDP") {
		return "UDP"
	}
	return "TCP" //StreamProxy模块默认不识别的都将按照TCP对待
}
