/**
 * 将json中的元素按照字段操作配置，映射为新的字段。大概过程是先将json反序列化到一个多层Map中（其中可能还套着数组），然后逐层遍历得到基于扁平key描述的二维数据结构。
 * 然后基于二维数据结构进行变换后，反向再重建为一个多层Map，最终将Map序列化成json。
 * @author duhaifeng
 * @date   2021/05/11
 */
package json_convert

import (
	"cv-api-gw/api_forward/field_format_convert"
	"cv-api-gw/common"
	"fmt"
	"regexp"
	"strings"
)

/**
 * 获取映射后的扁平数据
 */
func GetMappedFields(bodyMap map[string]interface{}, envValues map[string]string, fieldOperations []*common.FieldOperation, processAppend bool) map[string]interface{} {
	//if len(fieldOperations) == 0 {
	//	return bodyMap
	//}
	//定义json元素循环器
	elemRecursor := new(ElementRecursor)
	elemMapper := new(ElementMapper)
	elemMapper.Init(fieldOperations, envValues)
	elemRecursor.SetElementProcessor(elemMapper)
	elemRecursor.RecurseElement(bodyMap)
	//由于append的新节点无法通过上面的遍历添加，因此单独进行添加
	if processAppend {
		elemMapper.processAppendOperation()
	}
	return elemMapper.GetMappedData()
}

type ElementMapper struct {
	envValues       map[string]string
	fieldOperations []*common.FieldOperation
	mappedData      map[string]interface{}
}

func (this *ElementMapper) Init(fieldOperations []*common.FieldOperation, envValues map[string]string) {
	this.fieldOperations = fieldOperations
	this.envValues = envValues
	this.mappedData = make(map[string]interface{})
}

func (this *ElementMapper) saveMappedData(hierarchicKeys []string, data interface{}) {
	joinedKey := strings.Join(hierarchicKeys, ".")
	oldValue, ok := this.mappedData[joinedKey]
	if ok {
		//如果旧值是空值或空字符串才进行更新，避免通过append补偿生成一些带有默认值的字段时，覆盖了真正的值
		if oldValue == nil {
			this.mappedData[joinedKey] = data
		}
		oldValueStr, strOk := oldValue.(string)
		if strOk && oldValueStr == "" {
			this.mappedData[joinedKey] = data
		}
	} else {
		this.mappedData[joinedKey] = data
	}
}

func (this *ElementMapper) GetMappedData() map[string]interface{} {
	return this.mappedData
}

/**
 * 获取字段对应的操作定义
 */
func (this *ElementMapper) getFieldOperation(hierarchicKeys []string) *common.FieldOperation {
	//由于操作定义的key可以只定义高层key，因此可能存在当前元素命中多个key的情况，
	//处理原则是：
	//	1、del优先级高于move（元素都要被删掉了，其它命中的操作也没意义），这里实际上没有考虑move操作匹配的key也许更长
	//	2、相同的del操作返回最短匹配的key（要删就多删点）
	//	3、相同的move操作返回最长匹配的key（不删就精确操作）
	var delOperation *common.FieldOperation
	for _, fieldOperation := range this.fieldOperations {
		if fieldOperation.Operation == "del" {
			if this.isFieldNameMatch(fieldOperation.FieldName, hierarchicKeys) {
				if delOperation == nil {
					delOperation = fieldOperation
				} else {
					if len(fieldOperation.FieldName) < len(delOperation.FieldName) {
						delOperation = fieldOperation
					}
				}
			}
		}
	}
	if delOperation != nil {
		return delOperation
	}

	var moveOperation *common.FieldOperation
	for _, fieldOperation := range this.fieldOperations {
		if fieldOperation.Operation == "move" {
			if this.isFieldNameMatch(fieldOperation.FieldName, hierarchicKeys) {
				if moveOperation == nil {
					moveOperation = fieldOperation
				} else {
					if len(moveOperation.FieldName) < len(fieldOperation.FieldName) {
						moveOperation = fieldOperation
					}
				}
			}
		}
	}
	return moveOperation
}

/**
 * 判断某一层级的key是否为数组索引：[0]、[*]
 */
func (this *ElementMapper) isArrayIndex(str string) bool {
	pattern := "^\\[[0-9|\\*]*\\]$"
	match, _ := regexp.MatchString(pattern, str)
	return match
}

/**
 * 判断某一个字段操作是否就是当前遍历到字段对应的配置，之所以要单独判断是由于配置里只配了高层的key
 */
func (this *ElementMapper) isFieldNameMatch(operationDefinedFieldName string, hierarchicKeys []string) bool {
	fieldNameItems := strings.Split(operationDefinedFieldName, ".")
	//如果定义的字段长度大于当前遍历字段的长度，肯定不match
	//之所以这里不是=，而是>判断，是由于配置可能只配置顶层key，所以前缀匹配了也算匹配
	//例如当前遍历的层级是a.b[*].c.d，而配置中可以只配置a.b，降低配置工作量
	if len(fieldNameItems) > len(hierarchicKeys) {
		return false
	}

	for i, fieldNameItem := range fieldNameItems {
		hierarchicKey := hierarchicKeys[i]
		if fieldNameItem == hierarchicKey {
			continue
		}
		//如果当前层级key代表数组索引，且配置项为[*]、[]这种通用匹配形式，则认为对应有效
		if this.isArrayIndex(fieldNameItem) && this.isArrayIndex(hierarchicKey) {
			if fieldNameItem == "[*]" || fieldNameItem == "[]" {
				continue
			}
		}
		return false
	}
	return true
}

type ArrIndexHolder struct {
	BOTTOM     int
	cursor     int
	arrIndexes []string
}

/**
 * 获取移动映射后的层级键名，如果配置的目的key是泛化的[*]数组索引，则需要展开为具体的索引数字
 */
func (this *ElementMapper) getMovedKey(fieldOperation *common.FieldOperation, hierarchicKeys []string) []string {
	//从hierarchicKeys中把数组索引抽取出来，留作下面替换使用
	arrIndexHolder := new(ArrIndexHolder)
	arrIndexHolder.Init()
	for _, hierarchicKeyItem := range hierarchicKeys {
		if this.isArrayIndex(hierarchicKeyItem) {
			arrIndexHolder.AppendIndex(hierarchicKeyItem)
		}
	}

	destNameItems := strings.Split(fieldOperation.DestName, ".")
	var movedKeyItems []string
	if fieldOperation.ArrayIndexAlignMode == "right" {
		var invertedMovedKeyItems []string
		for i := len(destNameItems) - 1; i >= 0; i-- {
			destNameItem := destNameItems[i]
			arrIndex := ""
			if strings.HasPrefix(destNameItem, "[") && strings.HasSuffix(destNameItem, "]") {
				//即便destName中不明确地配置了带数字的索引（非[*]形式），也要从arrIndexHolder中弹出一个索引，目的是为了保证双方索引的对齐
				arrIndex = arrIndexHolder.popFromRight()
			}
			//配置定义的数组格式可能是[]或[*]，需要替代为实际的数组序号，避免重建数据时错误
			if destNameItem == "[*]" || destNameItem == "[]" {
				if arrIndex != "" {
					invertedMovedKeyItems = append(invertedMovedKeyItems, arrIndex)
				} else {
					common.Log.Error("move json node error: %s → %s: %s", fieldOperation.FieldName, fieldOperation.DestName, strings.Join(hierarchicKeys, "."))
					return nil
				}
			} else {
				invertedMovedKeyItems = append(invertedMovedKeyItems, destNameItem)
			}
		}
		for i := len(invertedMovedKeyItems) - 1; i >= 0; i-- {
			movedKeyItems = append(movedKeyItems, invertedMovedKeyItems[i])
		}
	} else {
		for _, destNameItem := range destNameItems {
			arrIndex := ""
			if strings.HasPrefix(destNameItem, "[") && strings.HasSuffix(destNameItem, "]") {
				//即便destName中不明确地配置了带数字的索引（非[*]形式），也要从arrIndexHolder中弹出一个索引，目的是为了保证双方索引的对齐
				arrIndex = arrIndexHolder.popFromLeft()
			}
			//配置定义的数组格式可能是[]或[*]，需要替代为实际的数组序号，避免重建数据时错误
			if destNameItem == "[*]" || destNameItem == "[]" {
				if arrIndex != "" {
					movedKeyItems = append(movedKeyItems, arrIndex)
				} else {
					common.Log.Error("move json node error: %s → %s: %s", fieldOperation.FieldName, fieldOperation.DestName, strings.Join(hierarchicKeys, "."))
					return nil
				}
			} else {
				movedKeyItems = append(movedKeyItems, destNameItem)
			}
		}
	}

	//由于move的配置可能只配置了高层key，因此需要将剩余层级的key也一并加上。
	//例如：如果只配置了devices→deviceList，那么应该满足devices.[*].deviceID → deviceList.[*].deviceID
	fieldNameItems := strings.Split(fieldOperation.FieldName, ".")
	if len(fieldNameItems) != len(hierarchicKeys) { //说明此时是部分定义
		for i := len(fieldNameItems); i < len(hierarchicKeys); i++ {
			movedKeyItems = append(movedKeyItems, hierarchicKeys[i])
		}
	}
	return movedKeyItems
}

/**
 * 处理append的元素
 */
func (this *ElementMapper) processAppendOperation() {
	needAppendData := make(map[string]interface{})
	for _, fieldOperation := range this.fieldOperations {
		if fieldOperation.Operation == "append" {
			expandedKeys := this.expandAppendKey(fieldOperation.FieldName)
			for _, expandedKey := range expandedKeys {
				//如果相应的key已经有了对应的值，则不进行覆盖
				if this.ifLongerKeyExist(expandedKey) {
					continue
				}
				value := this.replacePlaceholder(fieldOperation.Value)
				//如果定义了数据格式转换器，则进行格式转换
				if fieldOperation.FormatConverter != "" {
					formatConverter := field_format_convert.GetFormatConverter(fieldOperation.FormatConverter)
					if formatConverter != nil {
						value = formatConverter(value)
						//如果页面没有输入ID那么将会自动生成一个，这里是为了将自动生成的ID更新到环境变量中，以方便在组合请求中后续请求可使用相同的ID。
						//在这里更新envValues是利用Go中map是引用赋值的性质，不过在这里更新不太优雅（由于环境变量是在拦截器中最早生成的，所以跨度有点大）
						//有些ID在RequestBody连字段都不会传入，会导致无法借助move操作生成ID，因此借助配置一个append来补偿生成ID
						if fieldOperation.FormatConverter == "Create_ID_On_Empty" {
							this.envValues[fmt.Sprintf("{body:%s}", expandedKey)] = fmt.Sprintf("%v", value)
						}
					}
				}
				//先临时缓存要添加的数据，避免直接加入后影响ifLongerKeyExist方法的判断
				needAppendData[expandedKey] = value
			}
		}
	}
	for expandedKey, value := range needAppendData {
		this.saveMappedData(strings.Split(expandedKey, "."), value)
	}
}

/**
 * 判断一个具有相同前缀的更长的key是否已经存在。对于append操作，已经有了对应的值或者如果更长Key存在，则取消本次append，避免覆盖更详细的数据。
 * （由于append操作一般只用于补空）
 */
func (this *ElementMapper) ifLongerKeyExist(prefixKey string) bool {
	for mappedKey := range this.mappedData {
		if prefixKey == mappedKey || strings.HasPrefix(mappedKey, prefixKey) {
			return true
		}
	}
	return false
}

/**
 * 获取append后的key，获取move key一样，需要把[*]展开为实际的数组索引
 */
func (this *ElementMapper) expandAppendKey(tobeExpandKey string) []string {
	tobeExpandKey = strings.Replace(tobeExpandKey, "[]", "[*]", -1)
	if !strings.Contains(tobeExpandKey, "[*]") {
		return []string{tobeExpandKey}
	}
	//找到第一个[*]出现的位置，并进行展开。展开的关键是根据既有的数据展开，要与既有的数据match上
	//a.b.[*].c.[*].d →
	//		a.b.[0].c.[*].d →
	//				a.b.[0].c.[0].d
	//				a.b.[0].c.[1].d
	//		a.b.[1].c.[*].d、→
	//				a.b.[1].c.[0].d
	//		a.b.[2].c.[*].d→
	//				a.b.[2].c.[0].d
	//				a.b.[2].c.[1].d
	//				a.b.[2].c.[2].d

	//注意prefix要保留一个判断用的“.[”，这样做的目的是确保能一直数组匹配，不允许a.[*].b.c与a.[*].b.[*]发生命中
	keyPrefix := tobeExpandKey[:strings.Index(tobeExpandKey, "*]")]
	expandedKeys := make(map[string]bool)
	needNestedExpandKeys := make(map[string]bool)
	for mappedKey := range this.mappedData {
		if strings.HasPrefix(mappedKey, keyPrefix) {
			//找到第一个[*]在mappedData中对应的Keys，并依次从这些Keys中扣出其中的数组索引数，并替换[*]
			mappedKeyArrayStartIndex := strings.Index(mappedKey, keyPrefix) + len(keyPrefix)
			mappedKeyArrayIndexStr := mappedKey[mappedKeyArrayStartIndex:]
			mappedKeyArrayEndIndex := strings.Index(mappedKeyArrayIndexStr, "]")
			mappedKeyArrayIndexStr = mappedKeyArrayIndexStr[:mappedKeyArrayEndIndex]
			expandedKey := strings.Replace(tobeExpandKey, "[*]", "["+mappedKeyArrayIndexStr+"]", 1)

			if strings.Contains(expandedKey, "[*]") {
				//将需要二次展开的key收集起来，这里之所以不直接嵌套调用，而是先通过map保存是为了去重
				needNestedExpandKeys[expandedKey] = true
			} else {
				//如果是泛化展开的key，并且key最后一级就是数组，则该key非法。程序能执行到这一步，说明该key以前缀方式命中了mappedData中更长的key，
				//这种情况下append一个前缀进来，无法和既有的更长key的数据类型一致起来。
				if !strings.HasSuffix(expandedKey, "]") {
					//先通过map进行去重，发生重复的情况一般是由于append的key比较短，与mappedData中前缀重复的长key匹配了多次导致
					expandedKeys[expandedKey] = true
				}
			}
		}
	}
	for needNestedExpandKey := range needNestedExpandKeys {
		nestedExpendKeys := this.expandAppendKey(needNestedExpandKey)
		for _, nestedExpendKey := range nestedExpendKeys {
			//fmt.Println("nestedExpendKey++", expandedKey)
			expandedKeys[nestedExpendKey] = true
		}
	}
	var returnKeys []string
	for expandedKey := range expandedKeys {
		returnKeys = append(returnKeys, expandedKey)
	}
	return returnKeys
}

/**
 * 替换数据中的占位符
 */
func (this *ElementMapper) replacePlaceholder(data interface{}) interface{} {
	//TODO：这里没有考虑非字符串的情况（如果要跨类型替换，则必须要整个字符串对应一个占位符名字）
	if dataStr, ok := data.(string); ok {
		for k, v := range this.envValues {
			dataStr = strings.Replace(dataStr, k, v, -1)
		}
		return dataStr
	}
	return data
}

/**
 * 对接ElementRecursor的处理
 */
func (this *ElementMapper) Process(hierarchicKeys []string, data interface{}) {
	//置换数据中配置的占位符
	data = this.replacePlaceholder(data)
	fieldOperation := this.getFieldOperation(hierarchicKeys)
	if fieldOperation == nil {
		//如果未做映射配置，则保留原数据
		this.saveMappedData(hierarchicKeys, data)
		return
	}
	//如果定义了数据格式转换器，则进行格式转换
	if fieldOperation.FormatConverter != "" {
		formatConverter := field_format_convert.GetFormatConverter(fieldOperation.FormatConverter)
		if formatConverter != nil {
			data = formatConverter(data)
			//如果页面没有输入ID那么将会自动生成一个，这里是为了将自动生成的ID更新到环境变量中，以方便在组合请求中后续请求可使用相同的ID。
			//在这里更新envValues是利用Go中map是引用赋值的性质，不过在这里更新不太优雅（由于环境变量是在拦截器中最早生成的，所以跨度有点大）
			if fieldOperation.FormatConverter == "Create_ID_On_Empty" {
				this.envValues[fmt.Sprintf("{body:%s}", strings.Join(hierarchicKeys, "."))] = fmt.Sprintf("%v", data)
			}
		}
	}

	switch fieldOperation.Operation {
	case "del":
		//hiKey := strings.Join(hierarchicKeys, ".")
		//fmt.Println("del", hiKey)
		return
	case "move":
		movedKeyItems := this.getMovedKey(fieldOperation, hierarchicKeys)
		//oriKey := strings.Join(hierarchicKeys, ".")
		//mvKey := strings.Join(movedKeyItems, ".")
		//fmt.Println(oriKey, "->", mvKey)
		this.saveMappedData(movedKeyItems, data)
	default:
		//如果配置无法解析，默认保留原数据
		this.saveMappedData(hierarchicKeys, data)
	}
}
