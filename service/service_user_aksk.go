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

type GatewayAkskService struct {
	AkskDbOperator *db.AkskDbOperator
	UserDbOperator *db.UserDbOperator
	api.BaseService
}

func (this *GatewayAkskService) CreateAksk(user *GatewayUser) (*GatewayAksk, error) {
	err := this.DeleteAksk(user.UserId)
	if err != nil {
		return nil, err
	}
	newAksk := &GatewayAksk{UserId: user.UserId, AccessKey: uuid.New().String(), SecureKey: uuid.New().String()}
	newAksk.CreateTime = time.Now()
	newAksk.UpdateTime = time.Now()
	dbGatewayAksk := new(db.GatewayAksk)
	common.CopyStructField(newAksk, dbGatewayAksk)
	err = this.AkskDbOperator.CreateAksk(dbGatewayAksk)
	if err != nil {
		return nil, err
	}
	AccessAuthManager.addTokenToUserMapping(newAksk.AccessKey, user)
	return newAksk, nil
}

func (this *GatewayAkskService) GetUserAksk(userId string) (*GatewayAksk, error) {
	dbAkskDetail, err := this.AkskDbOperator.GetAksk(userId)
	if err != nil {
		return nil, err
	}
	if dbAkskDetail == nil {
		return nil, nil
	}
	gatewayAksk := new(GatewayAksk)
	common.CopyStructField(dbAkskDetail, gatewayAksk)
	return gatewayAksk, nil
}

func (this *GatewayAkskService) DeleteAksk(userId string) error {
	aksk, err := this.GetUserAksk(userId)
	if err != nil {
		return err
	}
	if aksk == nil {
		return nil
	}
	err = this.AkskDbOperator.DeleteAksk(userId)
	if err != nil {
		return err
	}
	AccessAuthManager.deleteAkskToUserMapping(aksk.AccessKey)
	return nil
}
