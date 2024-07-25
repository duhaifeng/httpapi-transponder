/**
 * 基于validator的数据验证工具
 * @author duhaifeng
 * @date   2021/06/30
 */
package common

import (
	"cv-api-gw/common/busierr"

	"github.com/go-playground/validator/v10"
)

func ValidateStruct(s interface{}) error {
	dataValidator := validator.New()
	err := dataValidator.Struct(s)
	if err == nil {
		return nil
	}
	if _, ok := err.(*validator.InvalidValidationError); ok {
		return busierr.WrapGoError(err)
	}
	//多重校验错误时，返回第一个错误
	for _, err := range err.(validator.ValidationErrors) {
		// err.Field() 只返回 field，使用 err.Namespace() 返回完整的路径
		return &busierr.ValidationError{err.Namespace(), err.Tag()}
	}
	return nil
}
