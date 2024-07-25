/**
 *
 * @author duhaifeng
 * @date   2021/04/17
 */
package service

import (
	"cv-api-gw/common"
	"cv-api-gw/db"
	api "github.com/duhaifeng/simpleapi"
	"github.com/google/uuid"
	"time"
)

type GatewayTokenService struct {
	TokenDbOperator *db.TokenDbOperator
	api.BaseService
}

func (this *GatewayTokenService) CreateToken(user *GatewayUser) (*GatewayToken, error) {
	err := this.DeleteToken(user.UserId)
	if err != nil {
		return nil, err
	}
	newToken := &GatewayToken{UserId: user.UserId, AccessToken: uuid.New().String(), FreshToken: uuid.New().String()}
	newToken.CreateTime = time.Now()
	newToken.UpdateTime = time.Now()
	dbGatewayToken := new(db.GatewayToken)
	common.CopyStructField(newToken, dbGatewayToken)
	err = this.TokenDbOperator.CreateToken(dbGatewayToken)
	if err != nil {
		return nil, err
	}
	AccessAuthManager.addTokenToUserMapping(newToken.AccessToken, user)
	return newToken, nil
}

func (this *GatewayTokenService) GetUserToken(userId string) (*GatewayToken, error) {
	dbTokenDetail, err := this.TokenDbOperator.GetToken(userId)
	if err != nil {
		return nil, err
	}
	if dbTokenDetail == nil {
		return nil, nil
	}
	gatewayToken := new(GatewayToken)
	common.CopyStructField(dbTokenDetail, gatewayToken)
	return gatewayToken, nil
}

func (this *GatewayTokenService) DeleteToken(userId string) error {
	token, err := this.GetUserToken(userId)
	if err != nil {
		return err
	}
	if token == nil {
		return nil
	}
	err = this.TokenDbOperator.DeleteToken(userId)
	if err != nil {
		return err
	}
	AccessAuthManager.deleteTokenToUserMapping(token.AccessToken)
	return nil
}
