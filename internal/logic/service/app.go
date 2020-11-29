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
	app, err := cache.AppCacheInst.Get(appId)
	if err != nil {
		return nil, err
	}
	if app != nil {
		return app, nil
	}

	app, err = dao.AppDaoInst.Get(appId)
	if err != nil {
		return nil, err
	}

	if app != nil {
		err = cache.AppCacheInst.Set(app)
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

	key := cache.SeqCacheInst.UserKey(appId, userId)
	exist,err := cache.SeqCacheInst.Exist(key)
	var currentSeq int64
	if exist != 1 {
		currentSeq,err = self.GetUserCurrentSeq(appId,userId,0)
		err = cache.SeqCacheInst.InitUserSeq(key,currentSeq+1)
		return currentSeq+1,err
	}

	return cache.SeqCacheInst.Incr(key)
}
func (self *SeqService) GetUserNextSeqFromDB(appId, userId int64) (int64, error) {
	return dao.UserSeqDaoInst.GetUserNextSeq(appId,userId,0)
}

//初始化用户序列号
func (self *SeqService) SetUserSeq(appId, userId,seq int64) error {
	return cache.SeqCacheInst.InitUserSeq(cache.SeqCacheInst.GetUserSeqKey(appId, userId),seq)
}
func (self * SeqService)GetUserCurrentSeq(appId, userId,groupId int64) (seq int64,err error){
	_,seq,err = dao.UserSeqDaoInst.GetUserSeq(appId,userId,groupId)
	return
}

// 获取下一个序列号
func (self *SeqService) GetGroupNextSeq(appId, groupId int64) (int64, error) {
	return cache.SeqCacheInst.Incr(cache.SeqCacheInst.GetUserSeqKey(appId, groupId))
}
func (self *SeqService) SetGroupSeq(appId, userId,seq int64) error {
	return cache.SeqCacheInst.InitUserSeq(cache.SeqCacheInst.GetGroupSeqKey(appId, userId),seq)
}