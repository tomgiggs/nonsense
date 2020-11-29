package service

import (
	"context"
	"nonsense/internal/logic/cache"
	"nonsense/internal/logic/dao"
)


var AppServiceInst = InitAppService()
var GroupServiceInst = InitGroupService()
var AuthServiceInst = InitAuthService()
var SeqServiceInst = InitSeqService()
var DeviceServiceInst = InitDeviceService()
var UserServiceInst = InitUserService()
var MessageServiceInst = InitMessageService()


type AppService struct{}
func InitAppService()*AppService{
	return &AppService{}
}
// Get 注册设备
func (self *AppService) Get(ctx context.Context, appId int64) (*dao.App, error) {
	app, err := cache.CacheInst.GetApp(appId)
	if err != nil {
		return nil, err
	}
	if app != nil {
		return app, nil
	}

	app, err = dao.Storage.GetAppInfo(appId)
	if err != nil {
		return nil, err
	}

	if app != nil {
		err = cache.CacheInst.SetApp(app)
		if err != nil {
			return app, nil
		}
	}

	return app, nil
}

type SeqService struct{}

func InitSeqService()*SeqService{
	return &SeqService{}
}

// 获取下一个序列号,键值在用户登录时已写入缓存
func (self *SeqService) GetUserNextSeq(appId, userId int64) (int64, error) {

	key := cache.CacheInst.GetUserSeqKey(appId, userId)
	exist,err := cache.CacheInst.IsUserSeqExist(key)
	var currentSeq int64
	if exist != 1 {
		currentSeq,err = self.GetUserCurrentSeq(appId,userId,0)
		err = cache.CacheInst.InitUserSeq(key,currentSeq+1)
		return currentSeq+1,err
	}

	return cache.CacheInst.IncrUserSeq(key)
}
func (self *SeqService) GetUserNextSeqFromDB(appId, userId int64) (int64, error) {
	return dao.Storage.GetUserNextSeq(appId,userId,0)
}

//初始化用户序列号
func (self *SeqService) SetUserSeq(appId, userId,seq int64) error {
	return cache.CacheInst.InitUserSeq(cache.CacheInst.GetUserSeqKey(appId, userId),seq)
}
func (self * SeqService)GetUserCurrentSeq(appId, userId,groupId int64) (seq int64,err error){
	seq,err = dao.Storage.GetUserSeq(appId,userId,groupId)
	return
}

// 获取下一个序列号
func (self *SeqService) GetGroupNextSeq(appId, groupId int64) (int64, error) {
	return cache.CacheInst.IncrUserSeq(cache.CacheInst.GetUserSeqKey(appId, groupId))
}
func (self *SeqService) SetGroupSeq(appId, userId,seq int64) error {
	return cache.CacheInst.InitUserSeq(cache.CacheInst.GetGroupSeqKey(appId, userId),seq)
}