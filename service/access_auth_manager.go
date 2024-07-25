/**
 *
 * @author duhaifeng
 * @date   2021/06/21
 */
package service

import (
	"cv-api-gw/common"
	"cv-api-gw/db"
	"sync"
	"time"
)

var AccessAuthManager = new(AccessAuthMapping)

/**
 * 创建token或aksk与userId的反向映射，用于网关接口请求认证
 */
type AccessAuthMapping struct {
	tokenToUser map[string]*GatewayUser
	akskToUser  map[string]*GatewayUser
	lock        *sync.RWMutex
}

func (this *AccessAuthMapping) Init() {
	this.lock = new(sync.RWMutex)
	this.tokenToUser = make(map[string]*GatewayUser)
	this.akskToUser = make(map[string]*GatewayUser)
	//每5分钟全量更新一次访问授权信息，避免通过其他入口修改了Token或AKSK
	go func() {
		this.refreshAllTokenAndAksk()
		time.Sleep(time.Minute * 5)
	}()
}

func (this *AccessAuthMapping) refreshAllTokenAndAksk() {
	this.lock.Lock()
	defer this.lock.Unlock()
	dbOperator := new(db.UserDbOperator)
	err := dbOperator.OpenOrmConnSeparately(common.Configs.Mysql.Host, common.Configs.Mysql.Port, common.Configs.Mysql.User, common.Configs.Mysql.Password, common.Configs.Mysql.Schema)
	if err != nil {
		return
	}
	userList, err := dbOperator.GetUserList()
	if err != nil {
		common.Log.Error("refresh all token and aksk from db error: %s", err.Error())
		return
	}
	akskToUser := make(map[string]*GatewayUser)
	tokenToUser := make(map[string]*GatewayUser)
	for _, dbUser := range userList {
		user := new(GatewayUser)
		common.CopyStructField(dbUser, user)
		for _, aksk := range dbUser.AkskList {
			akskToUser[aksk.AccessKey] = user
		}
		for _, token := range dbUser.TokenList {
			tokenToUser[token.AccessToken] = user
		}
	}
	this.akskToUser = akskToUser
	this.tokenToUser = tokenToUser
}

func (this *AccessAuthMapping) clearTokenMappingData() {
	this.lock.Lock()
	defer this.lock.Unlock()
	this.tokenToUser = make(map[string]*GatewayUser)
}

func (this *AccessAuthMapping) clearAkskMappingData() {
	this.lock.Lock()
	defer this.lock.Unlock()
	this.akskToUser = make(map[string]*GatewayUser)
}

func (this *AccessAuthMapping) addTokenToUserMapping(token string, user *GatewayUser) {
	this.lock.Lock()
	defer this.lock.Unlock()
	this.tokenToUser[token] = user
}

func (this *AccessAuthMapping) addAkskToUserMapping(accessKey string, user *GatewayUser) {
	this.lock.Lock()
	defer this.lock.Unlock()
	this.akskToUser[accessKey] = user
}

func (this *AccessAuthMapping) deleteTokenToUserMapping(token string) {
	this.lock.Lock()
	defer this.lock.Unlock()
	delete(this.tokenToUser, token)
}

func (this *AccessAuthMapping) deleteAkskToUserMapping(accessKey string) {
	this.lock.Lock()
	defer this.lock.Unlock()
	delete(this.akskToUser, accessKey)
}

func (this *AccessAuthMapping) GetTokenUser(token string) *GatewayUser {
	this.lock.RLock()
	defer this.lock.RUnlock()
	user, ok := this.tokenToUser[token]
	if !ok {
		return nil
	}
	return user
}

func (this *AccessAuthMapping) GetAkskUser(accessKey string) *GatewayUser {
	this.lock.RLock()
	defer this.lock.RUnlock()
	user, ok := this.akskToUser[accessKey]
	if !ok {
		return nil
	}
	return user
}
