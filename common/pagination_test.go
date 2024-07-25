package common

import (
	"fmt"
	"testing"
)

func TestCalcPagination(t *testing.T) {
	fmt.Println(CalcPagiInfo("2", "20", 100))
	fmt.Println(CalcPagiInfo("0", "20", 95))
	fmt.Println(CalcPagiInfo("2", "20", 95))
	fmt.Println(CalcPagiInfo("15", "20", 95))
	fmt.Println(CalcPagiInfo("15", "aaa", 95))
	fmt.Println(CalcPagiInfo("2", "1", 15))
}
