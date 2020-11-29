package cache

import (
	"encoding/json"
	jsoniter "github.com/json-iterator/go"
	"nonsense/internal/global"
	"nonsense/internal/logic/dao"
	"nonsense/pkg/common"
	"strconv"
	"time"

	"github.com/go-redis/redis"
)

var DeviceIPCacheInst = InitDeviceIpCache()
var UserCacheInst = InitUserCache()
var UserDeviceCacheInst = InitUserDeviceCache()
var GroupCacheInst = InitGroupCache()
var GroupUserCacheInst = InitGroupUserCache()
var SeqCacheInst = InitSeqCache()
var AppCacheInst = InitAppCache()

type AppCache struct{
	AppKey string
	AppExpire time.Duration
}

func InitAppCache()*AppCache{
	return &AppCache{
		AppKey    : "app:",
		AppExpire : 24 * time.Hour,
	}

}

func (self *AppCache) Get(appId int64) (*dao.App, error) {
	var app dao.App
	err := GetCache(self.AppKey+strconv.FormatInt(appId, 10), &app)
	if err != nil && err != redis.Nil {
		return nil, common.WrapError(err)
	}

	if err == redis.Nil {
		return nil, nil
	}
	return &app, nil
}

func (self *AppCache) Set(app *dao.App) error {
	err := SetCache(self.AppKey+strconv.FormatInt(app.Id, 10), app, self.AppExpire)
	if err != nil {
		return common.WrapError(err)
	}
	return nil
}


type SeqCache struct{
}
func InitSeqCache()*SeqCache{
	return &SeqCache{}
}

func (self *SeqCache) UserKey(appId, userId int64) string {
	return global.USER_SEQ_KEY_PREFIX + strconv.FormatInt(appId, 10) + ":" + strconv.FormatInt(userId, 10)
}
func (self *SeqCache) GetUserSeqKey(appId, userId int64) string {
	return global.USER_SEQ_KEY_PREFIX + strconv.FormatInt(appId, 10) + ":" + strconv.FormatInt(userId, 10)
}

func (self *SeqCache) GetGroupSeqKey(appId, groupId int64) string {
	return global.GROUP_SEQ_KEY_PREFIX + strconv.FormatInt(appId, 10) + ":" + strconv.FormatInt(groupId, 10)
}

//缓存是否存在
func (self *SeqCache) Exist(key string) (int64, error) {
	seq, err := global.StorageClient.RedisClient.Exists(key).Result()
	if err != nil {
		return 0, common.WrapError(err)
	}
	return seq, nil
}

// 将序列号增加1
func (self *SeqCache) Incr(key string) (int64, error) {
	seq, err := global.StorageClient.RedisClient.Incr(key).Result()
	if err != nil {
		return 0, common.WrapError(err)
	}
	return seq, nil
}


func (self *SeqCache) InitUserSeq(key string,seq int64) error {
	err := global.StorageClient.RedisClient.Set(key,seq,time.Duration(-1)).Err()
	if err != nil {
		return common.WrapError(err)
	}
	return nil
}
func (self *SeqCache) InitGroupSeq(key string,seq int64) error {
	err := global.StorageClient.RedisClient.Set(key,seq,time.Duration(-1)).Err()
	if err != nil {
		return common.WrapError(err)
	}
	return nil
}

// 将值序列化为接送并设置到redis
func SetCache(key string, value interface{}, duration time.Duration) error {
	bytes, err := jsoniter.Marshal(value)
	if err != nil {
		common.Sugar.Error(err)
		return err
	}

	err = global.StorageClient.RedisClient.Set(key, bytes, duration).Err()
	if err != nil {
		common.Sugar.Error(err)
		return err
	}
	return nil
}

//从redis中读取并反序列化为对象
func GetCache(key string, value interface{}) error {
	bytes, err := global.StorageClient.RedisClient.Get(key).Bytes()
	if err != nil {
		return err
	}
	err = jsoniter.Unmarshal(bytes, value)
	if err != nil {
		common.Sugar.Error(err)
		return err
	}
	return nil
}

type UserServerMapCache map[int64]map[string]bool

var UserServerMap = make(UserServerMapCache,0)

type UserChangeInfo struct {
	Event string
	Uid  int64
	SrvId string
}
func RefreshOnlineUserServer(){
	SubOnlineUserChange(global.USER_LOGIN_SERVER_KEY, func(info string) {
		changeInfo := new(UserChangeInfo)
		json.Unmarshal([]byte(info),&changeInfo)
		if changeInfo.Event == global.USER_CHANGE_EVENT_OFFLINE{
			if _,ok := UserServerMap[changeInfo.Uid][changeInfo.SrvId];ok{
				delete(UserServerMap[changeInfo.Uid],changeInfo.SrvId)
			}
			return
		}

		UserServerMap[changeInfo.Uid][changeInfo.SrvId] = true
	})
}
func PubOnlineUserChange(data *UserChangeInfo){

	jsonStr,_ := json.Marshal(data)
	err := global.StorageClient.RedisClient.HSet(global.USER_SERVER_MAP_KEY_PREFIX+strconv.FormatInt(data.Uid,10),data.SrvId,1).Err()
	if err != nil {
		common.Sugar.Error("更新用户所在服务器信息失败:",err)
	}

	sub := global.StorageClient.RedisClient.Subscribe(global.USER_LOGIN_SERVER_KEY)
	defer sub.Close()
	err = global.StorageClient.RedisClient.Publish("gim-user-login-server",string(jsonStr)).Err()
	if err != nil {
		common.Sugar.Error("发布用户更新消息失败:",err)
	}
}

func SubOnlineUserChange(channel string, callback func(info string)){

	sub := global.StorageClient.RedisClient.Subscribe(channel)
	defer sub.Close()
	for msg := range sub.Channel(){
		callback(msg.Payload)
	}
}
func GetUserServerFromRedis(uid int64) map[string]string{
	result := global.StorageClient.RedisClient.HGetAll(global.USER_SERVER_MAP_KEY_PREFIX+strconv.FormatInt(uid,10))
	data,_ := result.Result()
	return data
}
