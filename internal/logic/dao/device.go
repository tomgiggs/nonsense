package dao

import (
	"database/sql"
	"fmt"
	"nonsense/internal/global"
	"nonsense/pkg/common"
	"nonsense/pkg/storage"
)

type DeviceDao struct{
	daoClient *storage.DBClient
}

func InitDeviceDao()*DeviceDao{
	return &DeviceDao{
		daoClient: global.StorageClient,
	}
}

// 插入一条设备信息
func (self *DeviceDao) Add(device Device)(id int64,err error) {
	res, err := global.StorageClient.MysqlClient.Exec(`insert into device(device_id,app_id,type,brand,model,system_version,sdk_version,status,conn_id,user_ip) values(?,?,?,?,?,?,?,?,?,?)`,device.DeviceId, device.AppId, device.Type, device.Brand, device.Model, device.SystemVersion, device.SDKVersion, device.Status, "", 0)
	if err != nil {
		return 0,common.WrapError(err)
	}
	id,err =res.LastInsertId()
	if err != nil{
		common.Sugar.Error("get last insert id error: ",err)
	}
	common.Sugar.Info("register device :",id)

	return
}

// 获取设备
func (self *DeviceDao) Get(deviceId int64) (*Device, error) {
	device := Device{
		DeviceId: deviceId,
	}
	row := global.StorageClient.MysqlClient.QueryRow(`
		select app_id,user_id,type,brand,model,system_version,sdk_version,status,conn_id,user_ip,create_time,update_time
		from device where device_id = ?`, deviceId)
	err := row.Scan(&device.AppId, &device.UserId, &device.Type, &device.Brand, &device.Model, &device.SystemVersion, &device.SDKVersion,
		&device.Status, &device.ConnId, &device.UserIp, &device.CreateTime, &device.UpdateTime)
	if err != nil && err != sql.ErrNoRows {
		return nil, common.WrapError(err)
	}

	if err == sql.ErrNoRows {
		return nil, nil
	}

	return &device, err
}

// 查询用户所有的在线设备
func (self *DeviceDao) ListOnlineByUserId(appId, userId int64) ([]Device, error) {
	rows, err := global.StorageClient.MysqlClient.Query(
		`select device_id,type,brand,model,system_version,sdk_version,status,conn_id,user_ip,create_time,update_time from device where app_id = ? and user_id = ? and status = ?`,
		appId, userId, DeviceOnLine)
	if err != nil {
		return nil, common.WrapError(err)
	}

	devices := make([]Device, 0, 5)
	for rows.Next() {
		device := new(Device)
		err = rows.Scan(&device.DeviceId, &device.Type, &device.Brand, &device.Model, &device.SystemVersion, &device.SDKVersion,
			&device.Status, &device.ConnId, &device.UserIp, &device.CreateTime, &device.UpdateTime)
		if err != nil {
			common.Sugar.Error(err)
			return nil, err
		}
		devices = append(devices, *device)
	}
	return devices, nil
}

// 更新设备绑定用户和设备在线状态
func (self *DeviceDao) UpdateUserIdAndStatus(deviceId, userId int64, status int, connId string, UserIp string) error {
	_, err := global.StorageClient.MysqlClient.Exec("update device  set user_id = ?,status = ?,conn_id = ?,user_ip = ? where device_id = ? ",
		userId, status, connId, UserIp, deviceId)
	if err != nil {
		return common.WrapError(err)
	}
	return nil
}

// 更新设备的在线状态
func (self *DeviceDao) UpdateStatus(deviceId int64, status int) error {
	_, err := global.StorageClient.MysqlClient.Exec("update device set status = ? where device_id = ?", status, deviceId)
	if err != nil {
		return common.WrapError(err)
	}
	return nil
}

// 升级设备
func (self *DeviceDao) Upgrade(deviceId int64, systemVersion, sdkVersion string) error {
	_, err := global.StorageClient.MysqlClient.Exec("update device set system_version = ?,sdk_version = ? where device_id = ? ",
		systemVersion, sdkVersion, deviceId)
	if err != nil {
		return common.WrapError(err)
	}
	return nil
}


type UserSeqDao struct {

}
func NewUserSeqDao()*UserSeqDao{
	return &UserSeqDao{}
}

func(self *UserSeqDao)InitUserSeq(appId, userId int64)(err error){
	_, err = global.StorageClient.MysqlClient.Exec("insert into user_seq(app_id,user_id) values(?,?)", appId,userId)
	if err != nil {
		return common.WrapError(err)
	}
	return nil
}

func (self *UserSeqDao) GetUserSeq(appId, userId,groupId int64) ( ack int64,seq int64, err error) {
	row := global.StorageClient.MysqlClient.QueryRow(`select read_seq,receive_seq from user_seq where app_id = ? and group_id = ? and user_id = ?`, appId, groupId,userId)

	err = row.Scan(&ack,&seq)
	if err != nil {
		return 0, 0,common.WrapError(err)
	}
	return ack,seq, nil
}
func (self *UserSeqDao) UpdateAck(appId,userId,groupId, ack int64) error {
	fmt.Println(ack,groupId,userId,appId)
	_, err := global.StorageClient.MysqlClient.Exec("update user_seq set read_seq = ? where group_id=? and app_id=? and  user_id = ?",ack,groupId,appId, userId)
	if err != nil {
		common.Sugar.Errorf("update userseq error:%v",err)
		return common.WrapError(err)
	}
	return nil
}

func (self *UserSeqDao) Update(appId,userId,groupId, ack int64) error {
	fmt.Println(ack,groupId,userId,appId)
	_, err := global.StorageClient.MysqlClient.Exec("update user_seq set receive_seq = ? where group_id=? and app_id=? and  user_id = ?",ack,groupId,appId, userId)
	if err != nil {
		common.Sugar.Errorf("update userseq error:%v",err)
		return common.WrapError(err)
	}
	return nil
}

func (self *UserSeqDao) GetUserNextSeq(appId,userId,groupId int64) (int64,error) {
	tx,err := global.StorageClient.MysqlClient.Begin()
	rsp, err := global.StorageClient.MysqlClient.Exec("update user_seq set receive_seq = receive_seq+1 where group_id=? and app_id=? and  user_id = ?",groupId,appId, userId)
	if err != nil {
		tx.Rollback()
		return 0,common.WrapError(err)
	}
	tx.Commit()
	updateNum,err := rsp.RowsAffected()
	if err != nil {
		return 0,err
	}
	if updateNum == 0 {
		return 0,fmt.Errorf("record not exist")
	}
	_,seq,err := self.GetUserSeq(appId,userId,groupId)

	return seq,err
}