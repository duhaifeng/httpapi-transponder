/**
 * 存放扁平Key中的数组信息，辅助转换过程中的数组索引左对齐或右对齐判断
 * @author duhaifeng
 * @date   2021/07/16
 */
package json_convert

import "strings"

func (this *ArrIndexHolder) Init() {
	this.BOTTOM = -100
	this.cursor = this.BOTTOM
}

func (this *ArrIndexHolder) AppendIndex(idx string) {
	this.arrIndexes = append(this.arrIndexes, idx)
}

func (this *ArrIndexHolder) PopIndex(alignMode string) string {
	if strings.ToLower(alignMode) == "right" {
		return this.popFromRight()
	} else {
		return this.popFromLeft()
	}
}

func (this *ArrIndexHolder) popFromLeft() string {
	if len(this.arrIndexes) == 0 {
		return ""
	}
	if this.cursor == this.BOTTOM {
		this.cursor = 0
	} else {
		this.cursor++
	}
	if this.cursor >= len(this.arrIndexes) {
		return ""
	}
	return this.arrIndexes[this.cursor]
}

func (this *ArrIndexHolder) popFromRight() string {
	if len(this.arrIndexes) == 0 {
		return ""
	}
	if this.cursor == this.BOTTOM {
		this.cursor = len(this.arrIndexes) - 1
	} else {
		this.cursor--
	}
	if this.cursor < 0 {
		return ""
	}
	return this.arrIndexes[this.cursor]
}
