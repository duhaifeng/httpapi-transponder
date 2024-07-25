package json_convert

import (
	"fmt"
	"testing"
)

func Test_recurseArray(t *testing.T) {
	er := new(ElementRecursor)
	er.recurseArray(nil, []string{"a1", "a2", "a3"})
}

func Test_recurseData(t *testing.T) {
	er := new(ElementRecursor)
	m1 := make(map[string]interface{})
	//m2 := make(map[string]interface{})
	//m3 := make(map[string]interface{})
	//m4 := make(map[string]interface{})
	//m4["k4"] = 4
	//m2["k2-4"] = m4
	//m3["k3"] = "v3"
	//m2["k2-3"] = m3
	//m1["k1-2"] = m2
	//m1["k1-a"] = "a"
	var arr1 []string
	arr1 = append(arr1, "a1")
	arr1 = append(arr1, "a2")
	m1["k1-a[]"] = arr1
	er.recurseData(nil, m1)
}

func Test_convertData(t *testing.T) {
	er := new(ElementRecursor)
	elemMapper := new(ElementMapper)
	elemMapper.Init(nil, nil)
	er.SetElementProcessor(elemMapper)
	m1 := make(map[string]interface{})
	m2 := make(map[string]interface{})
	m3 := make(map[string]interface{})
	m4 := make(map[string]interface{})
	m4["k4"] = 4
	m2["k2-4"] = m4
	m3["k3"] = "v3"
	m2["k2-3"] = m3
	m1["k1-2"] = m2
	m1["k1-a"] = "a"

	var arr1 []string
	arr1 = append(arr1, "a1")
	arr1 = append(arr1, "a2")
	m1["k1-a[]"] = arr1
	er.RecurseElement(m1)
	fmt.Println(elemMapper.GetMappedData())
}
