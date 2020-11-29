package service

import (
	"context"
	"nonsense/internal/logic/cache"
	"nonsense/internal/logic/dao"
	"nonsense/pkg/common"
)

const (
	DeviceOffline = 0
	DeviceOnline  = 1
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

	deviceId,err = dao.Storage.AddDevice(device)
	if err != nil {
		return 0, err
	}

	return deviceId, nil
}

// ListOnlineByUserId 获取用户的所有在线设备
func (*DeviceService) ListOnlineByUserId(ctx context.Context, appId, userId int64) ([]dao.Device, error) {
	devices, err := cache.CacheInst.GetDevice(appId, userId)
	if err != nil {
		return nil, err
	}

	if devices != nil {
		return devices, nil
	}

	devices, err = dao.Storage.ListOnlineByUserId(appId, userId)
	if err != nil {
		return nil, err
	}

	err = cache.CacheInst.SetDevice(appId, userId, devices)
	if err != nil {
		return nil, err
	}

	return devices, nil
}

// 设备上线
func (self *DeviceService) Online(ctx context.Context, appId, deviceId, userId int64, connId string, UserIp string) error {
	err := dao.Storage.UpdateDeviceStatus(deviceId, DeviceOnline)
	if err != nil {
		return err
	}

	err = cache.CacheInst.DelDevice(appId, userId)
	if err != nil {
		return err
	}
	return nil
}

// Offline 设备离线
func (self *DeviceService) Offline(ctx context.Context, appId, userId, deviceId int64) error {
	err := dao.Storage.UpdateDeviceStatus(deviceId, DeviceOffline)
	if err != nil {
		return err
	}

	err = cache.CacheInst.DelDevice(appId, userId)
	if err != nil {
		return err
	}
	return nil
}
