package store

import (
	"github.com/go-redis/redis"
	jsoniter "github.com/json-iterator/go"
	"nonsense/internal/config"
	"nonsense/internal/global"
	"nonsense/pkg/common"
	"strconv"
	"time"
)

type Rcache struct {
	RedisClient *redis.Client
	Conf        *config.Access
}

var CacheInst *Rcache

func NewRcache(conf *config.Access) *Rcache {
	rc := &Rcache{
		Conf: conf,
	}
	addr := conf.Storage.Redis

	rc.RedisClient = redis.NewClient(&redis.Options{
		Addr:     addr,
		DB:       0,
		Password: conf.Storage.RedisPasswd,
	})

	_, err := rc.RedisClient.Ping().Result()
	if err != nil {
		common.Sugar.Error("redis err ", err)
		panic(err)
	}
	CacheInst = rc
	return rc
}

// app
func (rcc *Rcache) GetApp(appId int64) (*App, error) {
	var app App
	err := rcc.GetCache(global.APP_CACHE_KEY+strconv.FormatInt(appId, 10), &app)
	if err != nil && err != redis.Nil {
		return nil, common.ErrCacheError
	}

	if err == redis.Nil {
		return nil, nil
	}
	return &app, nil
}

func (rcc *Rcache) SetApp(app *App) error {
	err := rcc.SetCache(global.APP_CACHE_KEY+strconv.FormatInt(app.Id, 10), app, global.APP_CACHE_EXPIRE)
	if err != nil {
		return common.ErrCacheError
	}
	return nil
}

//user_seq
func (rcc *Rcache) GetUserSeqKey(appId, userId int64) string {
	return global.USER_SEQ_KEY_PREFIX + strconv.FormatInt(appId, 10) + ":" + strconv.FormatInt(userId, 10)
}

func (rcc *Rcache) GetGroupSeqKey(appId, groupId int64) string {
	return global.GROUP_SEQ_KEY_PREFIX + strconv.FormatInt(appId, 10) + ":" + strconv.FormatInt(groupId, 10)
}

func (rcc *Rcache) IsUserSeqExist(key string) (int64, error) {
	seq, err := rcc.RedisClient.Exists(key).Result()
	if err != nil {
		return 0, common.ErrCacheError
	}
	return seq, nil
}

func (rcc *Rcache) IncrUserSeq(key string) (int64, error) {
	seq, err := rcc.RedisClient.Incr(key).Result()
	if err != nil {
		return 0, common.ErrCacheError
	}
	return seq, nil
}

func (rcc *Rcache) InitUserSeq(key string, seq int64) error {
	err := rcc.RedisClient.Set(key, seq, time.Duration(-1)).Err()
	if err != nil {
		return common.ErrCacheError
	}
	return nil
}
func (rcc *Rcache) InitGroupSeq(key string, seq int64) error {
	err := rcc.RedisClient.Set(key, seq, time.Duration(-1)).Err()
	if err != nil {
		return common.ErrCacheError
	}
	return nil
}

//群组缓存
func (rcc *Rcache) GetGroupKey(appId, groupId int64) string {
	return global.GROUP_CACHE_KEY + strconv.FormatInt(appId, 10) + ":" + strconv.FormatInt(groupId, 10)
}

func (rcc *Rcache) GetGroup(appId, groupId int64) (*Group, error) {
	var groupInfo Group
	err := rcc.GetCache(rcc.GetGroupKey(appId, groupId), &groupInfo)
	if err != nil && err != redis.Nil {
		return nil, common.ErrCacheError
	}
	if err == redis.Nil {
		return nil, nil
	}
	return &groupInfo, nil
}

func (rcc *Rcache) SetGroup(group *Group) error {
	err := rcc.SetCache(rcc.GetGroupKey(group.AppId, group.GroupId), group, global.GROUP_CACHE_EXPIRE)
	if err != nil {
		return common.ErrCacheError
	}
	return nil
}

func (rcc *Rcache) DelGroup(appId, groupId int64) error {
	_, err := rcc.RedisClient.Del(rcc.GetGroupKey(appId, groupId)).Result()
	if err != nil {
		return common.ErrCacheError
	}
	return nil
}

//群组用户缓存
func (rcc *Rcache) GetGroupUserKey(appId, groupId int64) string {
	return global.GROUP_USER_CACHE_KEY + strconv.FormatInt(appId, 10) + ":" + strconv.FormatInt(groupId, 10)
}

// 获取群组成员
func (rcc *Rcache) GetGroupMembers(appId, groupId int64) (users []*GroupUserInfo, err error) {
	userMap, err1 := rcc.RedisClient.HGetAll(rcc.GetGroupUserKey(appId, groupId)).Result()
	if err1 != nil {
		return nil, common.WrapError(err1)
	}

	users = make([]*GroupUserInfo, 0, len(userMap))
	for _, v := range userMap {
		var user GroupUserInfo
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
func (rcc *Rcache) IsGroupMember(appId, groupId, userId int64) (bool, error) {
	is, err := rcc.RedisClient.HExists(rcc.GetGroupUserKey(appId, groupId), strconv.FormatInt(userId, 10)).Result()
	if err != nil {
		return false, common.ErrCacheError
	}

	return is, nil
}

// 获取群组成员数
func (rcc *Rcache) GetGroupMembersNum(appId, groupId int64) (int64, error) {
	membersNum, err := rcc.RedisClient.HLen(rcc.GetGroupUserKey(appId, groupId)).Result()
	if err != nil {
		return 0, common.ErrCacheError
	}
	return membersNum, nil
}

// 添加群组成员
func (rcc *Rcache) SetGroupMember(appId, groupId, userId int64, label, extra string) error {
	var user = GroupUserInfo{
		AppId:     appId,
		GroupId:   groupId,
		UserId:    userId,
		Label:     label,
		UserExtra: extra,
	}
	bytes, err := jsoniter.Marshal(user)
	if err != nil {
		return common.ErrCacheError
	}
	_, err = rcc.RedisClient.HSet(rcc.GetGroupUserKey(user.AppId, user.GroupId), strconv.FormatInt(user.UserId, 10), bytes).Result()
	if err != nil {
		return common.ErrCacheError
	}
	return nil
}

// 删除群组成员
func (rcc *Rcache) DelGroupMember(appId, groupId int64, userId int64) error {
	_, err := rcc.RedisClient.HDel(rcc.GetGroupUserKey(appId, groupId), strconv.FormatInt(userId, 10)).Result()
	if err != nil {
		return common.ErrCacheError
	}
	return nil
}

//用户缓存
func (rcc *Rcache) GetUserKey(appId, userId int64) string {
	return global.USER_CACHE_KEY + strconv.FormatInt(appId, 10) + ":" + strconv.FormatInt(userId, 10)
}

// 获取用户缓存
func (rcc *Rcache) GetUser(appId, userId int64) (*User, error) {
	var user User
	err := rcc.GetCache(rcc.GetUserKey(appId, userId), &user)
	if err != nil && err != redis.Nil {
		return nil, common.ErrCacheError
	}
	if err == redis.Nil {
		return nil, nil
	}
	return &user, nil
}

// 设置用户缓存
func (rcc *Rcache) SetUser(user User) error {
	err := rcc.SetCache(rcc.GetUserKey(user.AppId, user.UserId), user, global.USER_CACHE_EXPIRE)
	if err != nil {
		return common.ErrCacheError
	}
	return nil
}

// 删除用户缓存
func (rcc *Rcache) DelUser(appId, userId int64) error {
	_, err := rcc.RedisClient.Del(rcc.GetUserKey(appId, userId)).Result()
	if err != nil {
		return common.ErrCacheError
	}
	return nil
}
func (rcc *Rcache) GetUserServerFromRedis(uid int64) map[string]string {
	result := rcc.RedisClient.HGetAll(global.USER_SERVER_MAP_KEY_PREFIX + strconv.FormatInt(uid, 10))
	data, _ := result.Result()
	return data
}

//device
func (rcc *Rcache) GetDeviceKey(appId, userId int64) string {
	return global.DEVICE_CACHE_KEY + strconv.FormatInt(appId, 10) + ":" + strconv.FormatInt(userId, 10)
}

// Get 获取指定用户的所有在线设备
func (rcc *Rcache) GetDevice(appId, userId int64) ([]Device, error) {
	var devices []Device
	err := rcc.GetCache(rcc.GetDeviceKey(appId, userId), &devices)
	if err != nil && err != redis.Nil {
		return nil, common.ErrCacheError
	}

	if err == redis.Nil {
		return nil, nil
	}
	return devices, nil
}

// Set 将指定用户的所有在线设备存入缓存
func (rcc *Rcache) SetDevice(appId, userId int64, devices []Device) error {
	err := rcc.SetCache(rcc.GetDeviceKey(appId, userId), devices, global.DEVICE_CACHE_EXPIRE)
	if err != nil {
		return common.ErrCacheError
	}
	return nil
}

// Del 删除某一用户的在线设备列表
func (rcc *Rcache) DelDevice(appId, userId int64) error {
	_, err := rcc.RedisClient.Del(rcc.GetDeviceKey(appId, userId)).Result()
	if err != nil {
		return common.ErrCacheError
	}
	return nil
}

// 将值序列化为接送并设置到redis
func (rcc *Rcache) SetCache(key string, value interface{}, duration time.Duration) error {
	bytes, err := jsoniter.Marshal(value)
	if err != nil {
		common.Sugar.Error(err)
		return err
	}

	err = rcc.RedisClient.Set(key, bytes, duration).Err()
	if err != nil {
		common.Sugar.Error(err)
		return err
	}
	return nil
}

//从redis中读取并反序列化为对象
func (rcc *Rcache) GetCache(key string, value interface{}) error {
	bytes, err := rcc.RedisClient.Get(key).Bytes()
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
