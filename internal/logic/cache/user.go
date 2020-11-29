package cache

import (
	"nonsense/internal/global"
	"nonsense/internal/logic/dao"
	"nonsense/pkg/common"
	"nonsense/pkg/storage"
	"strconv"
	"time"

	"github.com/go-redis/redis"
)

type UserCache struct{
	UserKey string
	UserExpire time.Duration
	dbclient *storage.DBClient
}

func InitUserCache()*UserCache{
	return &UserCache{
		UserKey: "user:",
		UserExpire: 2 * time.Hour,
		dbclient: global.StorageClient,
	}
}

func (self *UserCache) Key(appId, userId int64) string {
	return self.UserKey + strconv.FormatInt(appId, 10) + ":" + strconv.FormatInt(userId, 10)
}

// 获取用户缓存
func (self *UserCache) Get(appId, userId int64) (*dao.User, error) {
	var user dao.User
	err := GetCache(self.Key(appId, userId), &user)
	if err != nil && err != redis.Nil {
		return nil, common.WrapError(err)
	}
	if err == redis.Nil {
		return nil, nil
	}
	return &user, nil
}

// 设置用户缓存
func (self *UserCache) Set(user dao.User) error {
	err := SetCache(self.Key(user.AppId, user.UserId), user, self.UserExpire)
	if err != nil {
		return common.WrapError(err)
	}
	return nil
}

// 删除用户缓存
func (self *UserCache) Del(appId, userId int64) error {
	_, err := global.StorageClient.RedisClient.Del(self.Key(appId, userId)).Result()
	if err != nil {
		return common.WrapError(err)
	}
	return nil
}
