package dao

import (
	"database/sql"
	"fmt"
	"nonsense/internal/global"
	"nonsense/pkg/common"
	"nonsense/pkg/storage"
)

var AppDaoInst = InitAppDao()
var GroupDaoInst = InitGroupDao()
var DeviceDaoInst = InitDeviceDao()
var GroupUserDaoInst = InitGroupUserDao()
var MessageDaoInst = InitMessageDao()
var UserDaoInst = InitUserDao()
var UserSeqDaoInst = NewUserSeqDao()


type AppDao struct{
	dbclient *storage.DBClient
}
func InitAppDao()*AppDao{
	return &AppDao{
		dbclient: global.StorageClient,
	}
}
// Get 获取APP信息
func (self *AppDao) Get(appId int64) (*App, error) {
	fmt.Println(global.StorageClient)
	var app App
	err := global.StorageClient.MysqlClient.QueryRow("select id,name,private_key,create_time,update_time from app where id = ?", appId).Scan(
		&app.Id, &app.Name, &app.PrivateKey, &app.CreateTime, &app.UpdateTime)
	if err != nil && err != sql.ErrNoRows {
		return nil, common.WrapError(err)
	}

	if err == sql.ErrNoRows {
		return nil, nil
	}
	return &app, nil
}
