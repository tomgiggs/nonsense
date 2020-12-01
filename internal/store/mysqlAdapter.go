package store

import (
	"database/sql"
	"fmt"
	"nonsense/internal/config"
	"nonsense/pkg/common"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

type MysqlAdapter struct {
	Name        string
	MysqlClient *sql.DB
	Conf        *config.Access
	DBclient    *gorm.DB
}

func NewMysqlAdapter() *MysqlAdapter {
	return &MysqlAdapter{
		Name: "mysql",
	}
}

func (mya *MysqlAdapter) Open(conf *config.Access) error {
	var err error
	mya.DBclient, err = gorm.Open("mysql", conf.Storage.MySQL)
	if err != nil {
		common.Sugar.Error("open db error:",err)
		panic(err)
	}
	mya.DBclient.SingularTable(true)// 全局禁用表名复数
	mya.DBclient.DB().SetMaxIdleConns(10)//数据库连接池
	mya.DBclient.DB().SetMaxOpenConns(100)
	//更改默认表名
	//gorm.DefaultTableNameHandler = func (db *gorm.DB, defaultTableName string) string  {
	//	return "prefix_" + defaultTableName;
	//}
	return err
}

func (mya *MysqlAdapter) Close() (err error) {
	if mya.DBclient != nil {
		err = mya.DBclient.Close()
	}
	if mya.MysqlClient != nil {
		err = mya.MysqlClient.Close()
	}
	return err

}

func (mya *MysqlAdapter) GetName() string {
	return mya.Name
}

func (mya *MysqlAdapter) AddUser(user *User) (int64, error) {
	result := mya.DBclient.Create(user)
	if result.Error != nil {
		common.Sugar.Error("insert user error:",result.Error)
		return 0, common.ErrDBError
	}
	return user.Id, nil
}

func (mya *MysqlAdapter) GetUser(appId, userId int64) (*User, error) {
	user := User{
		AppId:  appId,
		UserId: userId,
	}
	result := mya.DBclient.Where(&User{
		AppId: appId,
		UserId: userId,
	}).First(&user)
	err := result.Error
	if err != nil {
		return nil, common.ErrDBError
	}
	return &user, err
}

func (mya *MysqlAdapter) UpdateUser(user *User) error {
	result := mya.DBclient.Model(&user).Where("app_id = ?", user.AppId,"user_id",user.UserId).
		Updates(map[string]interface{}{"nickname": user.Nickname, "sex": user.Sex, "avatar_url": user.AvatarUrl,"extra":user.Extra})
	err := result.Error
	if err != nil {
		return common.ErrDBError
	}

	return nil
}

func (mya *MysqlAdapter) GetAppInfo(appId int64) (*App, error) {
	var app App
	result := mya.DBclient.Where(&App{
		Id: appId,
	}).First(&app)
	err := result.Error
	if err != nil {
		return nil, common.ErrDBError
	}

	return &app, nil
}

func (mya *MysqlAdapter) AddDevice(device Device) (id int64, err error) {
	result := mya.DBclient.Create(&device)
	err = result.Error
	if err != nil {
		common.Sugar.Error("insert device error:",err)
		return 0, common.ErrDBError
	}
	return device.Id, nil

}

func (mya *MysqlAdapter) GetDevice(deviceId int64) (*Device, error) {
	device := Device{
		DeviceId: deviceId,
	}
	result := mya.DBclient.Where(&Device{
		Id: deviceId,
	}).First(&device)
	err := result.Error
	if err != nil {
		return nil, common.ErrDBError
	}

	return &device, nil
}

func (mya *MysqlAdapter) UpdateDevice(device Device) error {
	return fmt.Errorf("unreachable")
}

func (mya *MysqlAdapter) DeleteDevice(deviceId int64) error {
	return fmt.Errorf("unreachable")
}
func (mya *MysqlAdapter) ListOnlineByUserId(appId, userId int64) ([]Device, error) {
	msgs := []Device{}
	result := mya.DBclient.Where(&Device{
		AppId: appId,
		UserId: userId,
	}).Find(&msgs)

	err := result.Error
	if err != nil {
		return nil, common.ErrDBError
	}
	return msgs, err
}
func (mya *MysqlAdapter) UpdateDeviceStatus(deviceId int64, status int) error {
	result := mya.DBclient.Model(&Device{}).Where("device_id = ?", deviceId).
		Updates(map[string]interface{}{"status": status})
	err := result.Error
	if err != nil {
		return common.ErrDBError
	}

	return nil
}

func (mya *MysqlAdapter) InitUserSeq(userSeq UserSeq) error {
	result := mya.DBclient.Create(&userSeq)
	err := result.Error
	if  err!= nil {
		common.Sugar.Error("insert userseq error:",result.Error)
		return common.ErrDBError
	}
	return nil
}

func (mya *MysqlAdapter) GetUserAck(appId, userId, groupId int64) (ack int64, err error) {
	userSeq := &UserSeq{}
	result := mya.DBclient.Select("read_seq").Where("app_id = ?", appId,"user_id",userId,"group_id",groupId).First(userSeq)
	err = result.Error
	if  err!= nil {
		common.Sugar.Error("get userseq error:",result.Error)
		return 0,common.ErrDBError
	}
	return userSeq.ReceiveSeq,nil

}

func (mya *MysqlAdapter) GetUserSeq(appId, userId, groupId int64) (seq int64, err error) {
	userSeq := &UserSeq{}
	result := mya.DBclient.Select("receive_seq").Where("app_id = ? and user_id=? and group_id=?", appId,userId,groupId).First(userSeq)
	err = result.Error
	if  err!= nil {
		common.Sugar.Error("get userseq error:",result.Error)
		return 0,common.ErrDBError
	}
	return userSeq.ReadSeq,nil

}

func (mya *MysqlAdapter) UpdateAck(appId, userId, groupId, ack int64) error {
	result := mya.DBclient.Model(&UserSeq{}).Where("app_id = ? and user_id=? and group_id=?", appId,userId,groupId).
		Updates(map[string]interface{}{"read_seq": ack})
	err := result.Error
	if err != nil {
		return common.ErrDBError
	}

	return nil
}

func (mya *MysqlAdapter) GetUserNextSeq(appId, userId, groupId int64) (int64, error) {
	tx := mya.DBclient.Begin()
	err := mya.DBclient.Exec("update user_seq set receive_seq = receive_seq+1 where group_id=? and app_id=? and  user_id = ?", groupId, appId, userId).Error
	if err != nil {
		common.Sugar.Error("db error: ", err)
		tx.Rollback()
		return 0, common.ErrDBError
	}
	result := tx.Commit()
	err = result.Error
	if err != nil {
		common.Sugar.Error("db error: ", err)
		return 0, err
	}
	seq, err := mya.GetUserSeq(appId, userId, groupId)

	return seq, err
}

func (mya *MysqlAdapter) AddGroup(group Group) (int64, error) {
	result := mya.DBclient.Create(&group)
	err := result.Error
	if  err!= nil {
		common.Sugar.Error("insert group error:",result.Error)
		return 0,common.ErrDBError
	}
	return group.Id,nil
}

func (mya *MysqlAdapter) UpdateGroup(group Group) (int64, error) {
	result := mya.DBclient.Model(&group).Where("app_id = ?", group.AppId,"group_id",group.GroupId).
		Updates(map[string]interface{}{"name": group.Name, "introduction": group.Introduction, "extra": group.Extra})
	err := result.Error
	if err != nil {
		return 0,common.ErrDBError
	}

	return 1,nil
}

func (mya *MysqlAdapter) GetGroup(appId, groupId int64) (*Group, error) {
	var group Group
	result := mya.DBclient.Where(&Group{
		AppId: appId,
		GroupId: groupId,
	}).First(&group)
	err := result.Error
	if err != nil {
		return nil, common.ErrDBError
	}
	return &group,nil
}

func (mya *MysqlAdapter) DeleteGroup(appId, groupId int64) error {
	return fmt.Errorf("unreachable")
}

func (mya *MysqlAdapter) AddMember(groupId int64, groupUser *GroupUserInfo) error {
	result := mya.DBclient.Create(groupUser)
	err := result.Error
	if  err!= nil {
		common.Sugar.Error("insert group_user error:",result.Error)
		return common.ErrDBError
	}
	return nil
}

func (mya *MysqlAdapter) GetMembers(appId, groupId int64) ([]GroupUserInfo, error) {
	msgs := []GroupUserInfo{}
	result := mya.DBclient.Where(&GroupUserInfo{
		AppId: appId,
		GroupId: groupId,
	}).Find(&msgs)

	err := result.Error
	if err != nil {
		return nil, common.ErrDBError
	}
	return msgs, err
}

func (mya *MysqlAdapter) DeleteMember(appId int64, groupId int64, userId int64) error {
	result := mya.DBclient.Delete(&GroupUserInfo{
		AppId: appId,
		GroupId: groupId,
		UserId: userId,
	})
	err := result.Error
	if err != nil {
		common.Sugar.Error("db error: ", err)
		return common.ErrDBError
	}
	return nil
}

func (mya *MysqlAdapter) UpdateMember(user *GroupUserInfo) error {
	result := mya.DBclient.Model(&user).Where("app_id = ? and group_id = ? and user_id = ?", user.AppId,user.GroupId,user.UserId).
		Updates(map[string]interface{}{"name": user.Nickname, "label": user.Label, "extra": user.UserExtra})
	err := result.Error
	if err != nil {
		return common.ErrDBError
	}

	return nil
}

// ListByUser 获取用户加入的群组信息
func (mya *MysqlAdapter) ListUserJoinGroup(appId, userId int64) ([]Group, error) {
	//mya.DBclient.Table("users").Select("users.name, emails.email").Joins("left join emails on emails.user_id = users.id").Scan(&results)

	result := mya.DBclient.Exec("select g.group_id,g.name,g.introduction,g.user_num,g.type,g.extra,g.create_time,g.update_time "+
			"from group_user u left join `group` g on u.app_id = g.app_id and u.group_id = g.group_id "+
			"where u.app_id = ? and u.user_id = ?", appId, userId)
	err := result.Error
	if err != nil {
		common.Sugar.Error("db error: ", err)
		return nil, common.ErrDBError
	}

	var groups []Group
	result.Scan(&groups)
	return groups, nil
}

func (mya *MysqlAdapter) AddMessage(message *Message) error {
	result := mya.DBclient.Create(message)
	err := result.Error
	if  err!= nil {
		common.Sugar.Error("insert message error:",result.Error)
		return common.ErrDBError
	}
	return nil
}

func (mya *MysqlAdapter) CancelMessage(msgId int64) error {
	return fmt.Errorf("unreachable")
}

func (mya *MysqlAdapter) ListMsgBySeq(appId, receiverId, seq int64) ([]Message, error) {
	msgs := []Message{}
	result := mya.DBclient.Where(&Message{
		AppId: appId,
		ReceiverId: receiverId,
	}).Where("seq >?",seq).Find(&msgs)

	err := result.Error
	if err != nil {
		return nil, common.ErrDBError
	}
	return msgs, err

}

func (mya *MysqlAdapter) AddMoments() error {
	return fmt.Errorf("unreachable")
}

func (mya *MysqlAdapter) GetMoments() error {
	return fmt.Errorf("unreachable")
}
