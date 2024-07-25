package busierr

import (
	"errors"
	"fmt"
	"testing"
)

func TestStackError(t *testing.T) {
	stackErr := newErrFunc()
	fmt.Println(stackErr.ToErrorWithStack())
	fmt.Println(stackErr.Stack())
}

func TestWrapStackError(t *testing.T) {
	stackErr := newWrapErrFunc()
	fmt.Println(stackErr.ToErrorWithStack())
	fmt.Println(stackErr.Stack())
}

func newErrFunc() *StackError {
	return NewError("new error")
}

func newWrapErrFunc() *StackError {
	return WrapGoError(errors.New("go error"))
}
