package service

import (
	"context"
	"nonsense/internal/global"
	"nonsense/internal/logic/cache"
	"nonsense/internal/logic/dao"
	"nonsense/pkg/common"
)
type ReqResult struct {
	code int32
	data interface{}
}

type UserService struct{}

func InitUserService()*UserService{
	return &UserService{}
}
// 添加用户
func (*UserService) Add(user dao.User) (int64,error) {
	userId, err := dao.Storage.AddUser(&user)
	if err != nil {
		return 0,err
	}
	if userId == 0 {
		return 0,common.ErrUserAlreadyExist
	}
	userSeq := dao.UserSeq{
		AppId: user.AppId,
		UserId: user.UserId,
	}
	err = dao.Storage.InitUserSeq(userSeq)
	return userId,err
}

// 获取用户信息
func (*UserService) Get(ctx context.Context, appId, userId int64) (*dao.User, error) {
	user, err := cache.CacheInst.GetUser(appId, userId)
	if err != nil {
		common.Sugar.Errorf("get user error:",err)
		return nil, err
	}
	if user != nil {
		return user, nil
	}

	user, err = dao.Storage.GetUser(appId, userId)
	if err != nil {
		common.Sugar.Errorf("get user error:",err)
		return nil, err
	}

	if user != nil {
		err = cache.CacheInst.SetUser(*user)
		if err != nil {
			return nil, err
		}
	}
	return user, err
}

func (*UserService) UpdateUserAckSeq(appId, userId,groupId,ack int64)  error{
	return dao.Storage.UpdateAck(appId,userId,groupId,ack)
}

func (*UserService) GetUserMaxACK(appId, userId,groupId int64)  *ReqResult{
	result := &ReqResult{
		code: global.REQ_RESULT_CODE_OK,
		data: 200,
	}
	ack,err := dao.Storage.GetUserAck(appId,userId,groupId)
	if err !=nil {
		result.code = global.REQ_RESULT_CODE_DB_ERR
	}
	result.data = ack
	return result
}


// 获取用户信息
func (*UserService) Update(ctx context.Context, user dao.User) error {
	err := dao.Storage.UpdateUser(&user)
	if err != nil {
		return err
	}

	err = cache.CacheInst.DelUser(user.AppId, user.UserId)
	if err != nil {
		return err
	}
	return nil
}
