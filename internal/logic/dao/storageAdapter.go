package dao

import (
	"fmt"
	"nonsense/internal/config"
)

type StorageAdapter interface {
	Open(conf *config.Access) error
	Close() error
	GetName() string
	//user
	AddUser(user *User) (int64, error)
	GetUser(appId,userId int64)(*User,error)
	UpdateUser(user *User) error
	//app
	GetAppInfo(appId int64) (*App,error)
	//device
	AddDevice(device Device) (id int64,err error)
	GetDevice(deviceId int64)(*Device,error)
	UpdateDevice(device Device) error
	DeleteDevice(deviceId int64) error
	ListOnlineByUserId(appId, userId int64) ([]Device, error)
	UpdateDeviceStatus(deviceId int64, status int) error
	//user_seq
	InitUserSeq(userSeq UserSeq) error
	GetUserAck(appId, userId,groupId int64) ( ack int64,err error)
	GetUserSeq(appId, userId,groupId int64) ( seq int64, err error)
	UpdateAck(appId,userId,groupId, ack int64) error
	GetUserNextSeq(appId,userId,groupId int64) (int64,error)
	//group 群信息
	AddGroup(group Group) (int64, error)
	UpdateGroup(group Group) (int64, error)
	GetGroup(appId, groupId int64) (*Group, error)
	DeleteGroup(appId, groupId int64) error

	//group_user 群成员
	AddMember(groupId int64,user *GroupUserInfo)error
	GetMembers(appId,groupId int64)([]*GroupUserInfo,error)
	DeleteMember(appId int64, groupId int64, userId int64) error
	UpdateMember(user *GroupUserInfo) error
	ListUserJoinGroup(appId, userId int64) ([]Group, error)
	//message 聊天消息
	AddMessage(message *Message) error
	CancelMessage(msgId int64) error //撤回消息
	ListMsgBySeq(appId, receiverId, seq int64) ([]Message, error) //取未读消息

	//moments 朋友圈消息
	AddMoments()error
	GetMoments()error

}


var AvailableAdapters = make(map[string]StorageAdapter)
var Storage StorageAdapter

func RegisterAdapter(sa StorageAdapter) {
	if sa == nil {
		panic("store: Register adapter is nil")
	}

	adapterName := sa.GetName()
	if _, ok := AvailableAdapters[adapterName]; ok {
		panic("store: adapter '" + adapterName + "' is already registered")
	}
	AvailableAdapters[adapterName] = sa
}


func OpenAdapter(conf *config.Access) error {
	if Storage == nil {
		if len(conf.Storage.Name) > 0 {
			if ad, ok := AvailableAdapters[conf.Storage.Name]; ok {
				Storage = ad
			} else {
				Storage = NewMysqlAdapter()
			}
		} else {
			return fmt.Errorf("store: db adapter is not specified. Please set `store_config.use_adapter` in `tinode.conf`")
		}
	}
	Storage.Open(conf)
	return nil
}