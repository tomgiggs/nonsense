package service

import (
	"context"
	"fmt"
	"nonsense/internal/store"
	"nonsense/pkg/common"
)

type UserService struct{}

func InitUserService()*UserService {
	return &UserService{}
}
// 添加用户
func (*UserService) Add(user store.User) (int64,error) {
	userId, err := store.Storage.AddUser(&user)
	if err != nil {
		return 0,common.ErrDBError
	}
	if userId == 0 {
		return 0,common.ErrUserAlreadyExist
	}
	userSeq := store.UserSeq{
		AppId: user.AppId,
		UserId: user.UserId,
	}
	err = store.Storage.InitUserSeq(userSeq)
	return userId,common.ErrDBError
}

// 获取用户信息
func (*UserService) Get(ctx context.Context, appId, userId int64) (*store.User, error) {
	user, err := store.CacheInst.GetUser(appId, userId)
	if err != nil {
		common.Sugar.Errorf("get user cache error:",err)
		return nil, common.ErrCacheError
	}
	if user != nil {
		return user, nil
	}

	user, err = store.Storage.GetUser(appId, userId)
	if err != nil {
		common.Sugar.Errorf("get user db error:",err)
		return nil, common.ErrDBError
	}
	fmt.Println(user,err)

	if user != nil {
		err = store.CacheInst.SetUser(*user)
		if err != nil {
			return nil, common.ErrCacheError
		}
	}
	return user, err
}

func (*UserService) UpdateUserAckSeq(appId, userId,groupId,ack int64)  error{
	return store.Storage.UpdateAck(appId,userId,groupId,ack)
}

func (*UserService) GetUserMaxACK(appId, userId,groupId int64)  (int64,error){
	return store.Storage.GetUserAck(appId,userId,groupId)
}


// 获取用户信息
func (*UserService) Update(ctx context.Context, user store.User) error {
	err := store.Storage.UpdateUser(&user)
	if err != nil {
		return common.ErrDBError
	}

	err = store.CacheInst.DelUser(user.AppId, user.UserId)
	if err != nil {
		return common.ErrCacheError
	}
	return nil
}
