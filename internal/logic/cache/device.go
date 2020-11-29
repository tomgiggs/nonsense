package cache

import (
	"github.com/go-redis/redis"
	"nonsense/internal/global"
	"nonsense/internal/logic/dao"
	"nonsense/pkg/common"
	"nonsense/pkg/storage"
	"strconv"
	"time"
)



type DeviceIPCache struct{
	DeviceIPKey string
	dbclient *storage.DBClient
}

func InitDeviceIpCache()*DeviceIPCache{
	return &DeviceIPCache{
		dbclient: global.StorageClient,
	}
}

func (self *DeviceIPCache) Key(deviceId int64) string {
	return self.DeviceIPKey + strconv.FormatInt(deviceId, 10)
}

// Get 获取设备所建立长连接的主机IP
func (self *DeviceIPCache) Get(deviceId int64) (string, error) {
	ip, err := global.StorageClient.RedisClient.Get(self.Key(deviceId)).Result()
	if err != nil && err != redis.Nil {
		return "", common.WrapError(err)
	}
	if err == redis.Nil {
		return "", nil
	}
	return ip, nil
}

// Set 设置设备所建立长连接的主机IP
func (self *DeviceIPCache) Set(deviceId int64, ip string) error {
	_, err := global.StorageClient.RedisClient.Set(self.Key(deviceId), ip, 0).Result()
	return common.WrapError(err)
}

// Del 删除设备所建立长连接的主机IP
func (self *DeviceIPCache) Del(deviceId int64) error {
	_, err := global.StorageClient.RedisClient.Del(self.Key(deviceId)).Result()
	if err != nil {
		return common.WrapError(err)
	}

	return nil
}

type UserDeviceCache struct{
	DeviceKey string
	DeviceExpire time.Duration
	dbclient *storage.DBClient
}

func InitUserDeviceCache()*UserDeviceCache{
	return &UserDeviceCache{
		DeviceKey  :"user:device:",
		DeviceExpire:2 * time.Hour,
		dbclient: global.StorageClient,
	}
}
func (self *UserDeviceCache) Key(appId, userId int64) string {
	return self.DeviceKey + strconv.FormatInt(appId, 10) + ":" + strconv.FormatInt(userId, 10)
}

// Get 获取指定用户的所有在线设备
func (self *UserDeviceCache) Get(appId, userId int64) ([]dao.Device, error) {
	var devices []dao.Device
	err := GetCache(self.Key(appId, userId), &devices)
	if err != nil && err != redis.Nil {
		return nil, common.WrapError(err)
	}

	if err == redis.Nil {
		return nil, nil
	}
	return devices, nil
}

// Set 将指定用户的所有在线设备存入缓存
func (self *UserDeviceCache) Set(appId, userId int64, devices []dao.Device) error {
	err := SetCache(self.Key(appId, userId), devices, self.DeviceExpire)
	return common.WrapError(err)
}

// Del 删除某一用户的在线设备列表
func (self *UserDeviceCache) Del(appId, userId int64) error {
	_, err := global.StorageClient.RedisClient.Del(self.Key(appId, userId)).Result()
	return common.WrapError(err)
}
