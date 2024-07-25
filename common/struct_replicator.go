/**
 * 将一个结构体字段数据复制到另一个结构体的工具（同名字段会复制）
 * @author duhaifeng
 * @date   2021/04/15
 */
package common

import (
	"reflect"
	"strings"
	"time"
)

/**
 * 将src中的数据段的数据赋值到des中
 * 只赋值数据字段名字和类型完全一样的字段，其余的字段不会赋值
 * 如果一旦出现desField.CanSet()==false的情况，则需要检测结构体中字段是否私有的(即小写的)
 */
func CopyStructField(src, des interface{}) {
	//如果src为nil，此处不能直接通过src==nil来判空，因为interface{}会把空指针也转型为一个对象
	if reflect.ValueOf(src).IsNil() {
		return
	}
	srcType := reflect.TypeOf(src)
	desType := reflect.TypeOf(des)
	srcRef := reflect.ValueOf(src).Elem()
	desRef := reflect.ValueOf(des).Elem()
	for i := 0; i < srcType.Elem().NumField(); i++ {
		if srcType.Elem().Field(i).Type.Kind() == reflect.Struct && srcType.Elem().Field(i).Type != reflect.TypeOf(time.Time{}) {
			for j := 0; j < desType.Elem().NumField(); j++ {
				if desType.Elem().Field(j).Type.Kind() == reflect.Struct {
					CopyStructField(srcRef.Field(i).Addr().Interface(), desRef.Field(j).Addr().Interface())
				}
			}
		} else if srcType.Elem().Field(i).Type.Kind() == reflect.Ptr {
			//指针类型
			for j := 0; j < desType.Elem().NumField(); j++ {
				if desType.Elem().Field(j).Type.Kind() == reflect.Ptr {
					CopyStructField(srcRef.Field(i).Interface(), desRef.Field(j).Interface())
				}
			}
		} else if srcType.Elem().Field(i).Type.Kind() == reflect.Slice {
			//切片类型
			//源结构体中数组属性的字段数组长度大于0才复制
			if srcRef.Field(i).Len() != 0 {
				//循环数组属性的字段每一项进行复制,要求目的结构体已经初始化数组进行接收
				//是根据源结构体的字段进行遍历的，因此需要判断目的结构体的字段是否存在，否目的结构体与源结构体不一样时可能出现数组越界
				for j := 0; j < srcRef.Field(i).Len(); j++ {
					if srcRef.Field(i).Index(j).Type().Kind() == reflect.Ptr {
						//添加判断目的结构体的字段是否合法或存在
						desRef.FieldByName(srcType.Elem().Field(i).Name)
						if desRef.IsValid() && desRef.FieldByName(srcType.Elem().Field(i).Name).IsValid() && srcRef.Field(i).Len() == desRef.FieldByName(srcType.Elem().Field(i).Name).Len() {
							CopyStructField(srcRef.Field(i).Index(j).Interface(), desRef.FieldByName(srcType.Elem().Field(i).Name).Index(j).Interface())
						}
					} else {
						//不是指针类型，赋值一次
						desField := desRef.FieldByName(srcType.Elem().Field(i).Name)
						if desField.IsValid() && desField.Type() == srcRef.Field(i).Type() && desField.CanSet() {
							desField.Set(srcRef.Field(i))
						}
						break
					}
				}
			}
		} else {
			desField := desRef.FieldByName(srcType.Elem().Field(i).Name)
			if desField.IsValid() && desField.Type() == srcRef.Field(i).Type() && desField.CanSet() {
				desField.Set(srcRef.Field(i))
			}
		}
	}
}

func ReplaceFieldValueRecursively(pointer interface{}, find, replaceValue string) {
	replaceFieldValueRecursively(reflect.ValueOf(pointer), find, replaceValue)
}

func replaceFieldValueRecursively(value reflect.Value, find, replaceValue string) {
	switch value.Kind() {
	case reflect.Slice, reflect.Array:
		for i := 0; i < value.Len(); i++ {
			replaceFieldValueRecursively(value.Index(i), find, replaceValue)
		}
	case reflect.Struct:
		for i := 0; i < value.NumField(); i++ {
			replaceFieldValueRecursively(value.Field(i), find, replaceValue)
		}
	case reflect.Map:
		for _, key := range value.MapKeys() {
			v := value.MapIndex(key)
			if v.Kind() == reflect.String {
				// map 的 key/value 不能原地更新
				// https://stackoverflow.com/questions/49147706/reflect-accessed-map-does-not-provide-a-modifiable-value
				str := v.String()
				if strings.Contains(str, find) {
					str = strings.ReplaceAll(str, find, replaceValue)
					value.SetMapIndex(key, reflect.ValueOf(str))
				}
			} else {
				replaceFieldValueRecursively(v, find, replaceValue)
			}
		}
	case reflect.Ptr, reflect.Interface:
		if !value.IsNil() {
			replaceFieldValueRecursively(value.Elem(), find, replaceValue)
		}
	case reflect.String:
		if value.IsValid() && value.CanSet() {
			str := value.String()
			if strings.Contains(str, find) {
				str = strings.ReplaceAll(str, find, replaceValue)
				value.SetString(str)
			}
		}
	}
}
