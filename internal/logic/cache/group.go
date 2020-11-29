package cache

import (
	"github.com/go-redis/redis"
	"nonsense/internal/global"
	"nonsense/internal/logic/dao"
	"nonsense/pkg/common"
	"nonsense/pkg/storage"
	"strconv"
	"time"

	jsoniter "github.com/json-iterator/go"
)

type GroupCache struct{
	dbclient *storage.DBClient
	groupKey string
	groupExpire time.Duration
}


func InitGroupCache()*GroupCache{
	return &GroupCache{
		dbclient: global.StorageClient,
		groupKey : "group:",
		groupExpire : 2 * time.Hour,
	}
}

func (self *GroupCache) Key(appId, groupId int64) string {
	return self.groupKey + strconv.FormatInt(appId, 10) + ":" + strconv.FormatInt(groupId, 10)
}

// 获取群组缓存
func (self *GroupCache) Get(appId, groupId int64) (*dao.Group, error) {
	var groupInfo dao.Group
	err := GetCache(self.Key(appId, groupId), &groupInfo)
	if err != nil && err != redis.Nil {
		return nil, common.WrapError(err)
	}
	if err == redis.Nil {
		return nil, nil
	}
	return &groupInfo, nil
}

// 设置群组缓存
func (self *GroupCache) Set(group *dao.Group) error {
	err := SetCache(self.Key(group.AppId, group.GroupId), group,self.groupExpire)
	if err != nil {
		return common.WrapError(err)
	}
	return nil
}

// 删除群组缓存
func (self *GroupCache) Del(appId, groupId int64) error {
	_, err := global.StorageClient.RedisClient.Del(self.Key(appId, groupId)).Result()
	if err != nil {
		return common.WrapError(err)
	}
	return nil
}


type GroupUserCache struct{
	dbclient *storage.DBClient
	groupUserKey string
}

func InitGroupUserCache()*GroupUserCache{
	return &GroupUserCache{
		dbclient: global.StorageClient,
		groupUserKey:"group_user:",
	}
}

func (self *GroupUserCache) Key(appId, groupId int64) string {
	return self.groupUserKey + strconv.FormatInt(appId, 10) + ":" + strconv.FormatInt(groupId, 10)
}

// 获取群组成员
func (self *GroupUserCache) Members(appId, groupId int64) (users []dao.GroupUserInfo, err error) {
	userMap, err1 := global.StorageClient.RedisClient.HGetAll(self.Key(appId, groupId)).Result()
	if err1 != nil {
		return nil, common.WrapError(err1)
	}

	users = make([]dao.GroupUserInfo, 0, len(userMap))
	for _, v := range userMap {
		var user dao.GroupUserInfo
		err = jsoniter.Unmarshal(common.Str2bytes(v), &user)
		if err != nil {
			common.Sugar.Error(err)
			continue
		}
		users = append(users, user)
	}
	return users, nil
}

// 是否是群组成员
func (self *GroupUserCache) IsMember(appId, groupId, userId int64) (bool, error) {
	is, err := global.StorageClient.RedisClient.HExists(self.Key(appId, groupId), strconv.FormatInt(userId, 10)).Result()
	if err != nil {
		return false, common.WrapError(err)
	}

	return is, nil
}

// 获取群组成员数
func (self *GroupUserCache) MembersNum(appId, groupId int64) (int64, error) {
	membersNum, err := global.StorageClient.RedisClient.HLen(self.Key(appId, groupId)).Result()
	if err != nil {
		return 0, common.WrapError(err)
	}
	return membersNum, nil
}

// 添加群组成员
func (self *GroupUserCache) Set(appId, groupId, userId int64, label, extra string) error {
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
	_, err = global.StorageClient.RedisClient.HSet(self.Key(user.AppId, user.GroupId), strconv.FormatInt(user.UserId, 10), bytes).Result()
	if err != nil {
		return common.WrapError(err)
	}
	return nil
}

// 删除群组成员
func (self *GroupUserCache) Del(appId, groupId int64, userId int64) error {
	_, err := global.StorageClient.RedisClient.HDel(self.Key(appId, groupId), strconv.FormatInt(userId, 10)).Result()
	if err != nil {
		return common.WrapError(err)
	}
	return nil
}
