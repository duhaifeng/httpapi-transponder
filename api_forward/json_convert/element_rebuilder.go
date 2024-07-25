/**
 * 基于扁平二维数据描述，重建得到多层级Map数据结构
 * @author duhaifeng
 * @date   2021/05/11
 */
package json_convert

import (
	"cv-api-gw/common"
	"encoding/json"
	"reflect"
	"regexp"
	"strconv"
	"strings"
)

type ElementRebuilder interface {
	InsertData(hierarchicKeys []string, data interface{})
	GetBuiltData() map[string]interface{}
	GetBuiltJsonData() string
}

type CachedElementRebuilder struct {
	isDataChanged bool
	flatMap       map[string]interface{}
	builtData     map[string]interface{}
}

/**
 * 先以扁平方式记录要重建的数据，降低资源消耗成本
 */
func (this *CachedElementRebuilder) InsertData(mappedFlatKey string, data interface{}) {
	if this.flatMap == nil {
		this.flatMap = make(map[string]interface{})
	}
	this.isDataChanged = true
	mappedFlatKey = strings.Replace(mappedFlatKey, " ", "", -1)
	//如果旧值是空值或空字符串才进行更新的逻辑在mapper中已经有了，这里其实可以去掉
	oldValue, ok := this.flatMap[mappedFlatKey]
	if ok {
		common.Log.Warn("build key conflict <%s>", mappedFlatKey)
		if oldValue == nil {
			this.flatMap[mappedFlatKey] = data
		}
		oldValueStr, strOk := oldValue.(string)
		if strOk && oldValueStr == "" {
			this.flatMap[mappedFlatKey] = data
		}
	} else {
		this.flatMap[mappedFlatKey] = data
	}
}

/**
 * 获取重建数据的JSON形式，便于预览
 */
func (this *CachedElementRebuilder) GetBuiltJsonData() string {
	jsonBytes, _ := json.MarshalIndent(this.GetBuiltData(), "", "  ") //使输出的json直接格式化为易读模式
	return string(jsonBytes)
}

/**
 * 基于扁平数据定义，将数据重建为层级形式
 */
func (this *CachedElementRebuilder) GetBuiltData() map[string]interface{} {
	//如果距离上次构建未发生任何改变，就不反复重建
	if this.builtData == nil || this.isDataChanged {
		builtMap := make(map[string]interface{})
		for flatKey, data := range this.flatMap {
			hierarchicKeys := strings.Split(flatKey, ".")
			//为了简化实现，这里限定顶层数据结构必须为对象，不能为数组
			if this.isArrayKey(hierarchicKeys[0]) {
				common.Log.Error("top level data is array %s", flatKey)
				continue
			}
			builtMap = this.buildData(builtMap, hierarchicKeys, data).(map[string]interface{})
		}
		this.builtData = builtMap
	}
	return this.builtData
}

/**
 * 判断指定的key是否对应为数组
 */
func (this *CachedElementRebuilder) isArrayKey(key string) bool {
	match, _ := regexp.MatchString("^\\[[0-9|\\*]*\\]$", key)
	return match
}

/**
 * 基于单个（注意是单个）扁平key重建数据层级
 */
func (this *CachedElementRebuilder) buildData(curLevelDataset interface{}, hierarchicKeys []string, data interface{}) interface{} {
	//fmt.Println(curLevelDataset, hierarchicKeys, data)
	curLevelKey := hierarchicKeys[0] //最顶层key即为当前要处理这一层的key
	if this.isArrayKey(curLevelKey) {
		return this.buildArrayData(curLevelDataset, hierarchicKeys, data)
	} else {
		return this.buildMapData(curLevelDataset, hierarchicKeys, data)
	}
}

/**
 * 重建map数据
 */
func (this *CachedElementRebuilder) buildMapData(curLevelDataset interface{}, hierarchicKeys []string, data interface{}) map[string]interface{} {
	//common.Log.Debug("build map data: %s %v %v", strings.Join(hierarchicKeys, "."), curLevelDataset, data)
	var curLevelMap map[string]interface{}
	//如果当前层级对应的数据集合已经存在就复用，否则说明是开始插入第一个元素，就新建一个新集合
	if curLevelDataset == nil {
		curLevelMap = make(map[string]interface{})
	} else {
		ok := true
		curLevelMap, ok = curLevelDataset.(map[string]interface{})
		if !ok {
			common.Log.Error("rebuild error: %s <%s> is not map[string]interface{} %v", strings.Join(hierarchicKeys, "."), reflect.TypeOf(curLevelDataset).String(), curLevelDataset)
			return curLevelMap
		}
		//bugfix:如果curLevelDataset为nil，上面转型ok的同时还会导致curLevelMap也为nil，导致下面出现赋值空指针错误
		if curLevelMap == nil {
			curLevelMap = make(map[string]interface{})
		}
	}
	curLevelKey := hierarchicKeys[0]
	if len(hierarchicKeys) > 1 {
		//借助下层递归构建当前一级数据
		curLevelData := this.buildData(curLevelMap[curLevelKey], hierarchicKeys[1:], data)
		//fmt.Println("put1: ", strings.Join(hierarchicKeys, "."), curLevelData)
		curLevelMap[curLevelKey] = curLevelData
	} else {
		//fmt.Println("put2: ", strings.Join(hierarchicKeys, "."), data)
		curLevelMap[curLevelKey] = data
	}
	return curLevelMap
}

/**
 * 重建数组数据，注意：对于b.[N]而言，b仍然对应一个map，[N]才对应一个数组
 */
func (this *CachedElementRebuilder) buildArrayData(curLevelDataset interface{}, hierarchicKeys []string, data interface{}) []interface{} {
	//common.Log.Debug("build array data: %s %v %v", strings.Join(hierarchicKeys, "."), curLevelDataset, data)
	curLevelIndex := this.extractArrayIndex(hierarchicKeys[0])
	var curLevelArray []interface{}
	//如果当前层级对应的数据集合已经存在就复用，否则说明是开始插入第一个元素，就新建一个新集合
	//对于数组比较特殊一点，需要根据元素的插入落位，动态扩充数组长度
	if curLevelDataset == nil {
		if curLevelIndex == -1 {
			curLevelIndex = 0
		}
		curLevelArray = make([]interface{}, curLevelIndex+1)
	} else {
		ok := true
		curLevelArray, ok = curLevelDataset.([]interface{})
		if !ok {
			common.Log.Error("rebuild error: %s <%s> is not []interface{}", strings.Join(hierarchicKeys, "."), reflect.TypeOf(curLevelDataset).String())
			return curLevelArray
		}
		//对[]或[*]形式，无法拿到索引，默认就加到当前数组的最后面
		if curLevelIndex == -1 {
			curLevelIndex = len(curLevelArray)
		}
		//如果当前传入的数组元素个数少于要放入数据的索引位，则用nil补位
		if len(curLevelArray) <= curLevelIndex {
			for i := len(curLevelArray); i <= curLevelIndex; i++ {
				curLevelArray = append(curLevelArray, nil)
			}
		}
	}

	//由于提前进行了空元素补位，因此这里直接赋值即可，而不是append
	if len(hierarchicKeys) > 1 {
		curLevelData := this.buildData(curLevelArray[curLevelIndex], hierarchicKeys[1:], data)
		//fmt.Println("append1: ", strings.Join(hierarchicKeys, "."), curLevelData)
		curLevelArray[curLevelIndex] = curLevelData
	} else {
		//fmt.Println("append2: ", strings.Join(hierarchicKeys, "."), data)
		curLevelArray[curLevelIndex] = data
	}
	return curLevelArray
}

/**
 * 将key中指定索引对应的数组元素取出来
 */
func (this *CachedElementRebuilder) extractArrayIndex(curArrayLevelKey string) int {
	arrIndexStr := curArrayLevelKey[strings.LastIndex(curArrayLevelKey, "[")+1:]
	arrIndexStr = strings.Replace(arrIndexStr, "]", "", -1)
	arrIndex, err := strconv.ParseInt(arrIndexStr, 10, 64)
	if err != nil {
		return -1
	}
	if arrIndex < 0 {
		return -1
	}
	return int(arrIndex)
}
