package service

import (
	"context"
	"nonsense/internal/store"
	"nonsense/pkg/common"
)

const (
	DeviceOffline = 0
	DeviceOnline  = 1
)

type DeviceService struct{}

func InitDeviceService()*DeviceService {
	return &DeviceService{}
}
// Register 注册设备
func (*DeviceService) Register(ctx context.Context, device store.Device) (int64, error) {
	app, err := AppServiceInst.Get(ctx, device.AppId)
	if err != nil {
		common.Sugar.Error(err)
		return 0, err
	}

	if app == nil {
		return 0, common.ErrBadRequest
	}
	var deviceId int64

	deviceId,err = store.Storage.AddDevice(device)
	if err != nil {
		return 0, err
	}

	return deviceId, nil
}

// ListOnlineByUserId 获取用户的所有在线设备
func (*DeviceService) ListOnlineByUserId(ctx context.Context, appId, userId int64) ([]store.Device, error) {
	devices, err := store.CacheInst.GetDevice(appId, userId)
	if err != nil {
		return nil, err
	}

	if devices != nil {
		return devices, nil
	}

	devices, err = store.Storage.ListOnlineByUserId(appId, userId)
	if err != nil {
		return nil, err
	}

	err = store.CacheInst.SetDevice(appId, userId, devices)
	if err != nil {
		return nil, err
	}

	return devices, nil
}

// 设备上线
func (self *DeviceService) Online(ctx context.Context, appId, deviceId, userId int64, connId string, UserIp string) error {
	err := store.Storage.UpdateDeviceStatus(deviceId, DeviceOnline)
	if err != nil {
		return err
	}

	err = store.CacheInst.DelDevice(appId, userId)
	if err != nil {
		return err
	}
	return nil
}

// Offline 设备离线
func (self *DeviceService) Offline(ctx context.Context, appId, userId, deviceId int64) error {
	err := store.Storage.UpdateDeviceStatus(deviceId, DeviceOffline)
	if err != nil {
		return err
	}

	err = store.CacheInst.DelDevice(appId, userId)
	if err != nil {
		return err
	}
	return nil
}
