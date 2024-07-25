/**
 * 为每个go routine生成和管理id的工具（为了便于打日志）
 * @author duhaifeng
 * @date   2021/04/15
 */
package common

import (
	"bytes"
	"runtime"
	"sync"
)

var RequestIdHolder RoutineRequestIdHolder

/**
 * 为每次请求的Go Routine保持一个Request ID
 */
type RoutineRequestIdHolder struct {
	routineReqIdMap map[string]string
	rwlock          sync.RWMutex
}

func (this *RoutineRequestIdHolder) PutRoutineReqId(reqId string) {
	defer this.rwlock.Unlock()
	this.rwlock.Lock()
	if this.routineReqIdMap == nil {
		this.routineReqIdMap = make(map[string]string)
	}
	this.routineReqIdMap[GetGoroutineNo()] = reqId
}

func (this *RoutineRequestIdHolder) GetRoutineReqId() string {
	defer this.rwlock.RUnlock()
	this.rwlock.RLock()
	if this.routineReqIdMap == nil {
		return "<no-request-id>"
	}
	reqId, ok := this.routineReqIdMap[GetGoroutineNo()]
	if !ok {
		return "<no-request-id>"
	}
	return reqId
}

func (this *RoutineRequestIdHolder) DelRoutineReqId() {
	defer this.rwlock.Unlock()
	this.rwlock.Lock()
	if this.routineReqIdMap == nil {
		return
	}
	delete(this.routineReqIdMap, GetGoroutineNo())
}

/**
 * 获取groutine no，由于go语言默认不提供，因此采用第三方实现
 */
func GetGoroutineNo() string {
	/**
	runtime.Stack()返回格式：
	goroutine 18 [running]:
	runtime/debug.Stack(0x0, 0x0, 0x0)
		/usr/local/go/src/runtime/debug/stack.go:24 +0xbe
	*/
	//查找当前的goroutine号（位于调用栈的的第一行中）
	stackBuf := make([]byte, 64)
	bufSize := runtime.Stack(stackBuf, false)
	stackBuf = stackBuf[:bufSize]
	stackBuf = bytes.TrimPrefix(stackBuf, []byte("goroutine "))
	stackBuf = stackBuf[:bytes.IndexByte(stackBuf, ' ')]
	routineNo := string(stackBuf)
	return routineNo
}
