package json_convert

import (
	"fmt"
	"regexp"
	"testing"
	"time"
)

func testRegex(str string) {
	pattern := "^\\[[0-9|\\*]*\\]$"
	match, _ := regexp.MatchString(pattern, str)
	fmt.Println(str, match)
}

func Test_GetBucketRegex(t *testing.T) {
	testRegex("[0]")
	testRegex("[10]")
	testRegex("[a]")
	testRegex("[]")
	testRegex("0]")
	testRegex("[0")
	testRegex("[ ]")
	testRegex("[]]")
	testRegex("[ 1 ]")
	testRegex("[*]")
}

func Test_buildMap(t *testing.T) {
	eb := new(CachedElementRebuilder)
	eb.InsertData("k1.k2.k3", "k1.k2.k3")
	eb.InsertData("k1.k2.k4", "k1.k2.k4")
	eb.InsertData("k1.k2-1", 123)
	eb.InsertData("k1-1", false)
	fmt.Println(eb.GetBuiltJsonData())
}

func Test_buildArray(t *testing.T) {
	eb := new(CachedElementRebuilder)
	eb.InsertData("k1-arr[0]", "k1-arr0")
	eb.InsertData("k1-arr[1]", "k1-arr1")
	eb.InsertData("k1.k2.arr1.[0]", "k1.k2.arr1-0")
	eb.InsertData("k1.k2.arr1.[2]", "k1.k2.arr1-1")
	//eb.InsertData("k1.k2.arr.[0].[0]", "k1.k2.arr0.0") //这种[0].[0]方式还处理不了
	eb.InsertData("k1.arr2.[0].arr2-1.[0]", "k1.k2.arr2.[0].arr2-1.[0]")
	eb.InsertData("k1.arr2.[2].arr2-1.[0]", "k1.k2.arr2.[2].arr2-1.[0]")
	eb.InsertData("k1.arr2.[2].arr2-1.[2]", "k1.k2.arr2.[2].arr2-1.[2]")
	fmt.Println(eb.GetBuiltJsonData())
	time.Sleep(time.Second)
}

func Test_buildData(t *testing.T) {
	eb := new(CachedElementRebuilder)
	eb.InsertData("k1.k2.k3", "k1.k2.k3")
	eb.InsertData("k1.k2.k4", "k1.k2.k4")
	eb.InsertData("k1.k2-1", 123)
	eb.InsertData("k1-1", false)

	eb.InsertData("k1-arr.[0]", "k1-arr0")
	eb.InsertData("k1-arr.[1]", "k1-arr1")
	eb.InsertData("k1.k2.arr1.[1]", "k1.k2.arr1-1")
	eb.InsertData("k1.k2.arr1.[0]", "k1.k2.arr1-0")
	eb.InsertData("arr2.[0].k5-1", "k1.k2.arr2-0.k5-1")
	eb.InsertData("arr2.[0].k5-2", "k1.k2.arr2-0.k5-2")
	eb.InsertData("arr2.[0].k5-3", "k1.k2.arr2-0.k5-3")
	eb.InsertData("arr2.[1].k5-4", "k1.k2.arr2-1.k5-4")
	eb.InsertData("arr2.[].k5-5", "k1.k2.arr2-.k5-5")
	eb.InsertData("arr2.[].k5-6", "k1.k2.arr2-.k5-6")
	fmt.Println(eb.GetBuiltJsonData())
	time.Sleep(time.Second)
}
