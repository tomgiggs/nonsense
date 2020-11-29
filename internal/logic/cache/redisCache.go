package cache

import (
	"encoding/json"
	"github.com/go-redis/redis"
	jsoniter "github.com/json-iterator/go"
	"nonsense/internal/config"
	"nonsense/internal/global"
	"nonsense/internal/logic/dao"
	"nonsense/pkg/common"
	"strconv"
	"time"
)

type Rcache struct {
	RedisClient *redis.Client
	Conf  *config.Access
}

var  CacheInst Rcache

func NewRcache(conf *config.Access)*Rcache{
	rc :=  &Rcache{
		Conf: conf,
	}
	addr := conf.Storage.Redis

	rc.RedisClient = redis.NewClient(&redis.Options{
		Addr: addr,
		DB:   0,
		Password: conf.Storage.RedisPasswd,
	})

	_, err := rc.RedisClient.Ping().Result()
	if err != nil {
		common.Sugar.Error("redis err ")
		panic(err)
	}
	return rc
}

// app
func (self *Rcache) GetApp(appId int64) (*dao.App, error) {
	var app dao.App
	err := GetCache(global.APP_CACHE_KEY+strconv.FormatInt(appId, 10), &app)
	if err != nil && err != redis.Nil {
		return nil, common.WrapError(err)
	}

	if err == redis.Nil {
		return nil, nil
	}
	return &app, nil
}

func (self *Rcache) SetApp(app *dao.App) error {
	err := SetCache(global.APP_CACHE_KEY+strconv.FormatInt(app.Id, 10), app, global.APP_CACHE_EXPIRE)
	if err != nil {
		return common.WrapError(err)
	}
	return nil
}


//user_seq
func (self *Rcache) GetUserSeqKey(appId, userId int64) string {
	return global.USER_SEQ_KEY_PREFIX + strconv.FormatInt(appId, 10) + ":" + strconv.FormatInt(userId, 10)
}

func (self *Rcache) GetGroupSeqKey(appId, groupId int64) string {
	return global.GROUP_SEQ_KEY_PREFIX + strconv.FormatInt(appId, 10) + ":" + strconv.FormatInt(groupId, 10)
}

func (self *Rcache) IsUserSeqExist(key string) (int64, error) {
	seq, err := self.RedisClient.Exists(key).Result()
	if err != nil {
		return 0, common.WrapError(err)
	}
	return seq, nil
}

func (self *Rcache) IncrUserSeq(key string) (int64, error) {
	seq, err := self.RedisClient.Incr(key).Result()
	if err != nil {
		return 0, common.WrapError(err)
	}
	return seq, nil
}

func (self *Rcache) InitUserSeq(key string,seq int64) error {
	err := self.RedisClient.Set(key,seq,time.Duration(-1)).Err()
	if err != nil {
		return common.WrapError(err)
	}
	return nil
}
func (self *Rcache) InitGroupSeq(key string,seq int64) error {
	err := self.RedisClient.Set(key,seq,time.Duration(-1)).Err()
	if err != nil {
		return common.WrapError(err)
	}
	return nil
}

//user discovery
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

func (self *Rcache) GetUserServerFromRedis(uid int64) map[string]string{
	result := global.StorageClient.RedisClient.HGetAll(global.USER_SERVER_MAP_KEY_PREFIX+strconv.FormatInt(uid,10))
	data,_ := result.Result()
	return data
}

//群组缓存
func (self *Rcache) GetGroupKey(appId, groupId int64) string {
	return global.GROUP_CACHE_KEY + strconv.FormatInt(appId, 10) + ":" + strconv.FormatInt(groupId, 10)
}

func (self *Rcache) GetGroup(appId, groupId int64) (*dao.Group, error) {
	var groupInfo dao.Group
	err := GetCache(self.GetGroupKey(appId, groupId), &groupInfo)
	if err != nil && err != redis.Nil {
		return nil, common.WrapError(err)
	}
	if err == redis.Nil {
		return nil, nil
	}
	return &groupInfo, nil
}

func (self *Rcache) SetGroup(group *dao.Group) error {
	err := SetCache(self.GetGroupKey(group.AppId, group.GroupId), group,global.GROUP_CACHE_EXPIRE)
	if err != nil {
		return common.WrapError(err)
	}
	return nil
}

func (self *Rcache) DelGroup(appId, groupId int64) error {
	_, err := global.StorageClient.RedisClient.Del(self.GetGroupKey(appId, groupId)).Result()
	if err != nil {
		return common.WrapError(err)
	}
	return nil
}

//群组用户缓存
func (self *Rcache) GetGroupUserKey(appId, groupId int64) string {
	return global.GROUP_USER_CACHE_KEY + strconv.FormatInt(appId, 10) + ":" + strconv.FormatInt(groupId, 10)
}

// 获取群组成员
func (self *Rcache) GetGroupMembers(appId, groupId int64) (users []*dao.GroupUserInfo, err error) {
	userMap, err1 := global.StorageClient.RedisClient.HGetAll(self.GetGroupUserKey(appId, groupId)).Result()
	if err1 != nil {
		return nil, common.WrapError(err1)
	}

	users = make([]*dao.GroupUserInfo, 0, len(userMap))
	for _, v := range userMap {
		var user dao.GroupUserInfo
		err = jsoniter.Unmarshal(common.Str2bytes(v), &user)
		if err != nil {
			common.Sugar.Error(err)
			continue
		}
		users = append(users, &user)
	}
	return users, nil
}

// 是否是群组成员
func (self *Rcache) IsGroupMember(appId, groupId, userId int64) (bool, error) {
	is, err := global.StorageClient.RedisClient.HExists(self.GetGroupUserKey(appId, groupId), strconv.FormatInt(userId, 10)).Result()
	if err != nil {
		return false, common.WrapError(err)
	}

	return is, nil
}

// 获取群组成员数
func (self *Rcache) GetGroupMembersNum(appId, groupId int64) (int64, error) {
	membersNum, err := global.StorageClient.RedisClient.HLen(self.GetGroupUserKey(appId, groupId)).Result()
	if err != nil {
		return 0, common.WrapError(err)
	}
	return membersNum, nil
}

// 添加群组成员
func (self *Rcache) SetGroupMember(appId, groupId, userId int64, label, extra string) error {
	var user = dao.GroupUserInfo{
		AppId :appId,
		GroupId: groupId,
		UserId:  userId,
		Label:   label,
		UserExtra:   extra,
	}
	bytes, err := jsoniter.Marshal(user)
	if err != nil {
		return common.WrapError(err)
	}
	_, err = global.StorageClient.RedisClient.HSet(self.GetGroupUserKey(user.AppId, user.GroupId), strconv.FormatInt(user.UserId, 10), bytes).Result()
	if err != nil {
		return common.WrapError(err)
	}
	return nil
}

// 删除群组成员
func (self *Rcache) DelGroupMember(appId, groupId int64, userId int64) error {
	_, err := global.StorageClient.RedisClient.HDel(self.GetGroupUserKey(appId, groupId), strconv.FormatInt(userId, 10)).Result()
	if err != nil {
		return common.WrapError(err)
	}
	return nil
}

//用户缓存
func (self *Rcache) GetUserKey(appId, userId int64) string {
	return global.USER_CACHE_KEY + strconv.FormatInt(appId, 10) + ":" + strconv.FormatInt(userId, 10)
}

// 获取用户缓存
func (self *Rcache) GetUser(appId, userId int64) (*dao.User, error) {
	var user dao.User
	err := GetCache(self.GetUserKey(appId, userId), &user)
	if err != nil && err != redis.Nil {
		return nil, common.WrapError(err)
	}
	if err == redis.Nil {
		return nil, nil
	}
	return &user, nil
}

// 设置用户缓存
func (self *Rcache) SetUser(user dao.User) error {
	err := SetCache(self.GetUserKey(user.AppId, user.UserId), user, global.USER_CACHE_EXPIRE)
	if err != nil {
		return common.WrapError(err)
	}
	return nil
}

// 删除用户缓存
func (self *Rcache) DelUser(appId, userId int64) error {
	_, err := global.StorageClient.RedisClient.Del(self.GetUserKey(appId, userId)).Result()
	if err != nil {
		return common.WrapError(err)
	}
	return nil
}

//device
func (self *Rcache) GetDeviceKey(appId, userId int64) string {
	return global.DEVICE_CACHE_KEY + strconv.FormatInt(appId, 10) + ":" + strconv.FormatInt(userId, 10)
}

// Get 获取指定用户的所有在线设备
func (self *Rcache) GetDevice(appId, userId int64) ([]dao.Device, error) {
	var devices []dao.Device
	err := GetCache(self.GetDeviceKey(appId, userId), &devices)
	if err != nil && err != redis.Nil {
		return nil, common.WrapError(err)
	}

	if err == redis.Nil {
		return nil, nil
	}
	return devices, nil
}

// Set 将指定用户的所有在线设备存入缓存
func (self *Rcache) SetDevice(appId, userId int64, devices []dao.Device) error {
	err := SetCache(self.GetDeviceKey(appId, userId), devices, global.DEVICE_CACHE_EXPIRE)
	return common.WrapError(err)
}

// Del 删除某一用户的在线设备列表
func (self *Rcache) DelDevice(appId, userId int64) error {
	_, err := global.StorageClient.RedisClient.Del(self.GetDeviceKey(appId, userId)).Result()
	return common.WrapError(err)
}

//------------------

type UserServerMapCache map[int64]map[string]bool

var UserServerMap = make(UserServerMapCache,0)

type UserChangeInfo struct {
	Event string
	Uid  int64
	SrvId string
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
