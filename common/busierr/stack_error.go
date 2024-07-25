/**
 * 带有堆栈信息的异常定义，之所以定义StackError是为了记录error发生点的stack信息以打印到日志中，方便快速定位error源头
 * @author duhaifeng
 * @date   2021/04/15
 */
package busierr

import (
	"errors"
	"fmt"
	"runtime"
	"runtime/debug"
)

type ErrorStack struct {
	File      string
	Line      int
	FuncName  string
	CallStack string
}

type StackError struct {
	msg   string
	stack ErrorStack
}

/**
 * 创建一个带有堆栈信息的error对象
 */
func NewError(msg string) *StackError {
	var err = &StackError{
		msg: msg,
	}
	if funcName, file, line, ok := runtime.Caller(1); ok {
		err.stack = ErrorStack{
			File:      file,
			Line:      line,
			FuncName:  runtime.FuncForPC(funcName).Name(),
			CallStack: string(debug.Stack()),
		}
	}

	return err
}

/**
 * 将go原生err包装为带有堆栈信息的error
 */
func WrapGoError(goErr error) *StackError {
	if goErr == nil {
		return nil
	}
	//避免二次保证
	if goStackError, ok := goErr.(*StackError); ok {
		return goStackError
	}
	err := &StackError{
		msg: goErr.Error(),
	}
	if funcName, file, line, ok := runtime.Caller(1); ok {
		err.stack = ErrorStack{
			File:      file,
			Line:      line,
			FuncName:  runtime.FuncForPC(funcName).Name(),
			CallStack: string(debug.Stack()),
		}
	}

	return err
}

/**
 * 获取go原生错误对象
 */
func (e StackError) ToError() error {
	return errors.New(e.msg)
}

/**
 * 获取带有堆栈基本信息的go原生错误对象
 */
func (e StackError) ToErrorWithStack() error {
	return errors.New(fmt.Sprintf("%s; %s:%d:%s", e.msg, e.stack.File, e.stack.Line, e.stack.FuncName))
}

/**
 * 获取go原生错误对象的message
 */
func (e StackError) Error() string {
	return e.msg
}

/**
 * 获取错误对象的堆栈信息
 */
func (e StackError) Stack() ErrorStack {
	return e.stack
}
