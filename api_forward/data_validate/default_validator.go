/**
 * 验证自动转发请求的请求数据是否合法
 * @author duhaifeng
 * @date   2021/09/27
 */
package data_validate

import (
	"cv-api-gw/common"
	"encoding/json"
	"errors"
	"fmt"
	"reflect"

	api "github.com/duhaifeng/simpleapi"
)

var DataValidators = make(map[string]func(r *api.Request) error)

func init() {
	DataValidators["DefaultValidator"] = DefaultValidator
	//自定义验证器需要在此注册
	DataValidators["CustomValidator1"] = CustomValidator1
	DataValidators["CustomValidator2"] = CustomValidator2
	DataValidators["CustomValidator3"] = CustomValidator3
}

func GetDataValidator(validatorName string) func(r *api.Request) error {
	if validatorName == "" {
		return DefaultValidator
	}
	dataValidator, ok := DataValidators[validatorName]
	if ok {
		return dataValidator
	}
	return DefaultValidator
}

func DefaultValidator(r *api.Request) error {
	return nil
}

func validateRequestBody(r *api.Request, validationType reflect.Type) error {
	validationObject := reflect.New(validationType).Interface()
	body, _ := r.GetBody()
	if len(body) != 0 {
		err := json.Unmarshal(body, &validationObject)
		if err != nil {
			common.Log.Error("validate request body err: %s", err.Error())
			return errors.New(fmt.Sprintf("can not validate request data: %s", err.Error()))
		}
	}
	return common.ValidateStruct(validationObject)
}
