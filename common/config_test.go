package common

import (
	"encoding/json"
	"fmt"
	"testing"
)

func Test_getCurrentDirectory(t *testing.T) {
	fmt.Println(GetCurrentDirectory())
}

func Test_ParseConfigFile(t *testing.T) {
	err := ParseConfigFile()
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(GetApiForwardConf("/v1/a/b/c/d", "GET"))
	fmt.Println(GetApiForwardConf("/v2/a/b/c/d", "GET"))
}

func Test_MultiMarshal(t *testing.T) {
	m := make(map[string]string)
	m1 := map[string]string{"k1": "v1"}
	jsonBytes, _ := json.Marshal(m1)
	json.Unmarshal(jsonBytes, &m)
	fmt.Println(m)
	m2 := map[string]string{"k2": "v2"}
	jsonBytes, _ = json.Marshal(m2)
	json.Unmarshal(jsonBytes, &m)
	fmt.Println(m)

	var arr []string
	arr1 := []string{"a1", "a2"}
	jsonBytes, _ = json.Marshal(arr1)
	json.Unmarshal(jsonBytes, &arr)
	fmt.Println(arr)
	arr2 := []string{"b1", "b2"}
	jsonBytes, _ = json.Marshal(arr2)
	json.Unmarshal(jsonBytes, &arr)
	fmt.Println(arr)
}
