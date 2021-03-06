package service

import (
	"context"
	"fmt"
	"nonsense/internal/global"
	"nonsense/internal/store"
	"nonsense/pkg/common"
)

type AuthService struct{}

func InitAuthService()*AuthService {
	return &AuthService{}
}

// 登录
func (self *AuthService) SignIn(ctx context.Context, appId, userId int64, deviceId int64, passwd string, connId string, UserIp string)(string, error) {
	user,err := UserServiceInst.Get(ctx,appId,userId)
	if user == nil{
		return "",fmt.Errorf("user not exist")
	}
	if err != nil{
		return "", err
	}
	if user.Passwd != passwd{
		return "",fmt.Errorf("wrong passwd")
	}
	// 生成token
	tokenInfo := common.JwtTokenInfo{
		UserId: userId,
		Passwd: passwd,
		AppId: appId,
	}
	global.UserTokenInfos[userId] = &tokenInfo

	tokenStr := common.JwtEncry(tokenInfo)

	// 标记用户在设备上登录
	err = DeviceServiceInst.Online(ctx, appId, deviceId, userId, connId, UserIp)
	if err != nil {
		return "",err
	}
	//初始化用户seq到缓存
	var seq int64
	seq,err = store.Storage.GetUserSeq(appId,userId,0)
	SeqServiceInst.SetUserSeq(appId,userId,seq)

	return tokenStr,nil
}

// Auth 验证用户是否登录
func (self *AuthService) Auth(ctx context.Context, appId, userId int64, passwd string, token string) error {
	return self.VerifyToken(ctx, appId, userId, passwd, token)
}

// 对用户秘钥进行校验
func (self *AuthService) VerifyToken(ctx context.Context, appId, userId int64, passwd string, token string) error {

	//info, err := common.DecryptToken(token, app.PrivateKey)
	//if err != nil {
	//	return common.ErrUnauthorized
	//}

	info,succ := common.JwtDecry(token)
	if !succ {
		return common.ErrUnauthorized
	}

	if !(info.AppId == appId && info.UserId == userId && info.Passwd == passwd) {
		return common.ErrUnauthorized
	}

	return nil
}

func (self *AuthService) IsTokenExpire(ctx context.Context, token string) error {

	_,succ := common.JwtDecry(token)
	if !succ {
		return common.ErrUnauthorized
	}

	return nil
}