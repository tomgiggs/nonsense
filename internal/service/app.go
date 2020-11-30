package service

import (
	"context"
	"nonsense/internal/store"
)


var AppServiceInst = InitAppService()
var GroupServiceInst = InitGroupService()
var AuthServiceInst = InitAuthService()
var SeqServiceInst = InitSeqService()
var DeviceServiceInst = InitDeviceService()
var UserServiceInst = InitUserService()
var MessageServiceInst = InitMessageService()


type AppService struct{}
func InitAppService()*AppService {
	return &AppService{}
}
// Get 注册设备
func (self *AppService) Get(ctx context.Context, appId int64) (*store.App, error) {
	app, err := store.CacheInst.GetApp(appId)
	if err != nil {
		return nil, err
	}
	if app != nil {
		return app, nil
	}

	app, err = store.Storage.GetAppInfo(appId)
	if err != nil {
		return nil, err
	}

	if app != nil {
		err = store.CacheInst.SetApp(app)
		if err != nil {
			return app, nil
		}
	}

	return app, nil
}

type SeqService struct{}

func InitSeqService()*SeqService {
	return &SeqService{}
}

// 获取下一个序列号,键值在用户登录时已写入缓存
func (self *SeqService) GetUserNextSeq(appId, userId int64) (int64, error) {

	key := store.CacheInst.GetUserSeqKey(appId, userId)
	exist,err := store.CacheInst.IsUserSeqExist(key)
	var currentSeq int64
	if exist != 1 {
		currentSeq,err = self.GetUserCurrentSeq(appId,userId,0)
		err = store.CacheInst.InitUserSeq(key,currentSeq+1)
		return currentSeq+1,err
	}

	return store.CacheInst.IncrUserSeq(key)
}
func (self *SeqService) GetUserNextSeqFromDB(appId, userId int64) (int64, error) {
	return store.Storage.GetUserNextSeq(appId,userId,0)
}

//初始化用户序列号
func (self *SeqService) SetUserSeq(appId, userId,seq int64) error {
	return store.CacheInst.InitUserSeq(store.CacheInst.GetUserSeqKey(appId, userId),seq)
}
func (self *SeqService)GetUserCurrentSeq(appId, userId,groupId int64) (seq int64,err error){
	seq,err = store.Storage.GetUserSeq(appId,userId,groupId)
	return
}

// 获取下一个序列号
func (self *SeqService) GetGroupNextSeq(appId, groupId int64) (int64, error) {
	return store.CacheInst.IncrUserSeq(store.CacheInst.GetUserSeqKey(appId, groupId))
}
func (self *SeqService) SetGroupSeq(appId, userId,seq int64) error {
	return store.CacheInst.InitUserSeq(store.CacheInst.GetGroupSeqKey(appId, userId),seq)
}