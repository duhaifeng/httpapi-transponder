/**
 * 自定义验证器，可根据自身需求随意定制数据验证器
 * @author duhaifeng
 * @date   2021/09/27
 */
package data_validate

import (
	"cv-api-gw/service"
	"errors"
	"reflect"

	api "github.com/duhaifeng/simpleapi"
)

func CustomValidator1(r *api.Request) error {
	orgID := r.GetUrlVar("orgID")
	if orgID == "0" {
		return errors.New("root organization can no be deleted")
	}
	return nil
}

func CustomValidator2(r *api.Request) error {
	return validateRequestBody(r, reflect.TypeOf(service.Device{}))
}

func CustomValidator3(r *api.Request) error {
	return validateRequestBody(r, reflect.TypeOf(service.ImageGroup{}))
}
