/**
 * 用于遍历多层Map和数组嵌套的遍历器，对于每个节点会触发一个回调，这是Recursor与Mapper结合的地方
 * @author duhaifeng
 * @date   2021/05/11
 */
package json_convert

import (
	"fmt"
	"reflect"
)

type ElementProcessor interface {
	Process([]string, interface{})
}

type ElementRecursor struct {
	elemProcessor ElementProcessor
}

func (this *ElementRecursor) SetElementProcessor(elementProcessor ElementProcessor) {
	this.elemProcessor = elementProcessor
}

func (this *ElementRecursor) RecurseElement(jsonMap map[string]interface{}) {
	this.recurseData(nil, jsonMap)
}

func (this *ElementRecursor) recurseData(hierarchicKeys []string, srcData interface{}) {
	if hierarchicKeys == nil {
		hierarchicKeys = make([]string, 0)
	}
	v := reflect.ValueOf(srcData)
	if v.Kind() == reflect.Map {
		this.recurseMap(hierarchicKeys, srcData)
	} else if v.Kind() == reflect.Slice || v.Kind() == reflect.Array {
		if v.Len() != 0 {
			this.recurseArray(hierarchicKeys, srcData)
		} else {
			//如果当前是个空数组，则没必要对其进行深入递归，注意不能再用索引形式对其映射
			this.elemProcessor.Process(hierarchicKeys, srcData)
		}
	} else {
		//fmt.Printf("%s=%v \n", strings.Join(hierarchicKeys, "."), srcData)
		if this.elemProcessor != nil {
			this.elemProcessor.Process(hierarchicKeys, srcData)
		}
	}
}

func (this *ElementRecursor) recurseMap(hierarchicKeys []string, srcMap interface{}) {
	v := reflect.ValueOf(srcMap)
	for _, key := range v.MapKeys() {
		nextLevelHierarchicKeys := hierarchicKeys
		nextLevelHierarchicKeys = append(nextLevelHierarchicKeys, key.String())
		//fmt.Printf("%s = %v \n", key, val)
		this.recurseData(nextLevelHierarchicKeys, v.MapIndex(key).Interface())
	}
}

func (this *ElementRecursor) recurseArray(hierarchicKeys []string, srcArray interface{}) {
	//arrKey := hierarchicKeys[len(hierarchicKeys)-1]
	v := reflect.ValueOf(srcArray)
	for i := 0; i < v.Len(); i++ {
		elementHierarchicKeys := hierarchicKeys
		elementHierarchicKeys = append(elementHierarchicKeys, fmt.Sprintf("[%d]", i))
		this.recurseData(elementHierarchicKeys, v.Index(i).Interface())
	}
}
