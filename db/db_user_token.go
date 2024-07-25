/**
 * 用户Token相关数据库操作，Token目前主要用于用户登录（为了兼容AIBox）
 * @author duhaifeng
 * @date   2021/04/15
 */
package db

import (
	"cv-api-gw/common"
	"cv-api-gw/common/busierr"

	api "github.com/duhaifeng/simpleapi"
)

type TokenDbOperator struct {
	api.BaseDbOperator
}

func (this *TokenDbOperator) CreateToken(gatewayToken *GatewayToken) error {
	err := this.OrmConn().Create(gatewayToken).Error
	if err != nil {
		common.Log.Error("[db error] create gateway token error: %s", err.Error())
		return busierr.WrapGoError(err)
	}
	return nil
}

func (this *TokenDbOperator) GetAllToken() ([]*GatewayToken, error) {
	var tokenList TokenList
	err := this.OrmConn().Find(&tokenList).Error
	if err != nil {
		common.Log.Error("[db error] get token list error: %s", err.Error())
		return nil, busierr.WrapGoError(err)
	}
	return tokenList, nil
}

func (this *TokenDbOperator) GetToken(userId string) (*GatewayToken, error) {
	gatewayToken := new(GatewayToken)
	err := this.OrmConn().Where("user_id=?", userId).Find(&gatewayToken).Error
	if err != nil {
		if err.Error() == "record not found" {
			return nil, nil
		}
		common.Log.Error("[db error] get gateway token error: %s", err.Error())
		return nil, busierr.WrapGoError(err)
	}
	if gatewayToken.UserId == "" {
		return nil, nil
	}
	return gatewayToken, nil
}

func (this *TokenDbOperator) DeleteToken(userId string) error {
	err := this.OrmConn().Where("user_id=?", userId).Delete(GatewayToken{}).Error
	if err != nil {
		common.Log.Error("[db error] delete gateway token error: %s", err.Error())
		return busierr.WrapGoError(err)
	}
	return nil
}
