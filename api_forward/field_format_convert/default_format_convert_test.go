package field_format_convert

import (
	"fmt"
	"testing"
	"time"
)

func TestConvertDatetimeToGoTime(t *testing.T) {
	layout := "2006-01-02T15:04:05Z"
	timeStr := "2021-05-21T10:37:21Z"
	ti, err := time.Parse(layout, timeStr)
	fmt.Println(ti, err)
	fmt.Println(time.Date(1970, 1, 1, 0, 0, 0, 0, time.UTC))
	timeZone, _ = time.LoadLocation("Asia/Shanghai")
	fmt.Println(timeZone)
}

func TestConvertUnixIntToGoTime(t *testing.T) {
	fmt.Println(time.Now())
	fmt.Println(time.Now().Unix())
	fmt.Println(time.Now().UnixNano())
	unixTime := time.Now().Unix()
	fmt.Println(unixTime)
	fmt.Println(time.Unix(unixTime, 0))
	fmt.Println(time.Unix(1624442043, 0))
	fmt.Println(time.Unix(1626092919755, 0))
}

func TestGoTimeFormat(t *testing.T) {
	now := time.Now()
	RFC3339Ms := "2006-01-02T15:04:05.999Z07:00"
	fmt.Println(now.Format(time.RFC3339Nano))
	fmt.Println(now.Format(RFC3339Ms))
}

func TestMapCopy(t *testing.T) {
	//测试map是否为引用赋值
	m1 := make(map[string]string)
	m1["k1"] = "v1"
	m1["k2"] = "v2"
	m2 := m1
	m2["k2"] = "v2-a"
	m2["k3"] = "v3"
	fmt.Println(m1, m2)
}
