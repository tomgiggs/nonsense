package service

import (
	"context"
	"nonsense/internal/logic/cache"
	"nonsense/internal/logic/dao"
	"nonsense/pkg/common"
)

const (
	DeviceOnline  = 1
	DeviceOffline = 0
)

type DeviceService struct{}

func InitDeviceService()*DeviceService{
	return &DeviceService{}
}
// Register 注册设备
func (*DeviceService) Register(ctx context.Context, device dao.Device) (int64, error) {
	app, err := AppServiceInst.Get(ctx, device.AppId)
	if err != nil {
		common.Sugar.Error(err)
		return 0, err
	}

	if app == nil {
		return 0, common.ErrBadRequest
	}
	var deviceId int64

	deviceId,err = dao.DeviceDaoInst.Add(device)
	if err != nil {
		return 0, err
	}

	return deviceId, nil
}

// ListOnlineByUserId 获取用户的所有在线设备
func (*DeviceService) ListOnlineByUserId(ctx context.Context, appId, userId int64) ([]dao.Device, error) {
	devices, err := cache.UserDeviceCacheInst.Get(appId, userId)
	if err != nil {
		return nil, err
	}

	if devices != nil {
		return devices, nil
	}

	devices, err = dao.DeviceDaoInst.ListOnlineByUserId(appId, userId)
	if err != nil {
		return nil, err
	}

	err = cache.UserDeviceCacheInst.Set(appId, userId, devices)
	if err != nil {
		return nil, err
	}

	return devices, nil
}

// 设备上线
func (self *DeviceService) Online(ctx context.Context, appId, deviceId, userId int64, connId string, UserIp string) error {
	err := dao.DeviceDaoInst.UpdateUserIdAndStatus(deviceId, userId, DeviceOnline, connId, UserIp)
	if err != nil {
		return err
	}

	err = cache.UserDeviceCacheInst.Del(appId, userId)
	if err != nil {
		return err
	}
	return nil
}

// Offline 设备离线
func (self *DeviceService) Offline(ctx context.Context, appId, userId, deviceId int64) error {
	err := dao.DeviceDaoInst.UpdateStatus(deviceId, DeviceOffline)
	if err != nil {
		return err
	}

	err = cache.UserDeviceCacheInst.Del(appId, userId)
	if err != nil {
		return err
	}

	err = cache.DeviceIPCacheInst.Del(deviceId)
	if err != nil {
		return err
	}
	return nil
}
