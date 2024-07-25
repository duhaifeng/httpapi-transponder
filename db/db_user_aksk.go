/**
 * AKSK相关数据库操作，一个用户可以创建多个AKSK用于外部程序使用
 * @author duhaifeng
 * @date   2021/04/15
 */
package db

import (
	"cv-api-gw/common"
	"cv-api-gw/common/busierr"

	api "github.com/duhaifeng/simpleapi"
)

type AkskDbOperator struct {
	api.BaseDbOperator
}

func (this *AkskDbOperator) CreateAksk(gatewayAksk *GatewayAksk) error {
	err := this.OrmConn().Create(gatewayAksk).Error
	if err != nil {
		common.Log.Error("[db error] create gateway aksk error: %s", err.Error())
		return busierr.WrapGoError(err)
	}
	return nil
}

func (this *AkskDbOperator) GetAllAksk() ([]*GatewayAksk, error) {
	var akskList AkskList
	err := this.OrmConn().Find(&akskList).Error
	if err != nil {
		common.Log.Error("[db error] get aksk list error: %s", err.Error())
		return nil, busierr.WrapGoError(err)
	}
	return akskList, nil
}

func (this *AkskDbOperator) GetAksk(userId string) (*GatewayAksk, error) {
	gatewayAksk := new(GatewayAksk)
	err := this.OrmConn().Where("user_id=?", userId).Find(&gatewayAksk).Error
	if err != nil {
		if err.Error() == "record not found" {
			return nil, nil
		}
		common.Log.Error("[db error] get gateway aksk error: %s", err.Error())
		return nil, busierr.WrapGoError(err)
	}
	if gatewayAksk.UserId == "" {
		return nil, nil
	}
	return gatewayAksk, nil
}

func (this *AkskDbOperator) DeleteAksk(userId string) error {
	err := this.OrmConn().Where("user_id=?", userId).Delete(GatewayAksk{}).Error
	if err != nil {
		common.Log.Error("[db error] delete gateway aksk error: %s", err.Error())
		return busierr.WrapGoError(err)
	}
	return nil
}
