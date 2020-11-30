package store

import (
	"database/sql"
	"fmt"
	"nonsense/internal/config"
	"nonsense/pkg/common"
)

type MysqlAdapter struct {
	Name        string
	MysqlClient *sql.DB
	Conf        *config.Access
}

func NewMysqlAdapter() *MysqlAdapter {
	return &MysqlAdapter{
		Name: "mysql",
	}
}

func (mya *MysqlAdapter) Open(conf *config.Access) error {
	var err error
	mya.MysqlClient, err = sql.Open("mysql", conf.Storage.MySQL)
	if err != nil {
		panic(err)
	}
	return err
}

func (mya *MysqlAdapter) Close() error {
	return mya.MysqlClient.Close()
}

func (mya *MysqlAdapter) GetName() string {
	return mya.Name
}

func (mya *MysqlAdapter) AddUser(user *User) (int64, error) {
	result, err := mya.MysqlClient.Exec("insert ignore into user(app_id,user_id,nickname,sex,avatar_url,extra) values(?,?,?,?,?,?)",
		user.AppId, user.UserId, user.Nickname, user.Sex, user.AvatarUrl, user.Extra)
	if err != nil {
		return 0, common.ErrDBError
	}

	lastId, err := result.LastInsertId()
	if err != nil {
		return 0, common.ErrDBError
	}
	return lastId, nil
}

func (mya *MysqlAdapter) GetUser(appId, userId int64) (*User, error) {
	row := mya.MysqlClient.QueryRow(
		"select nickname,sex,birthday,passwd,email,weixin_openid,avatar_url,extra,create_time,last_login_time from user where app_id = ? and user_id = ?",
		appId, userId)
	user := User{
		AppId:  appId,
		UserId: userId,
	}

	err := row.Scan(&user.Nickname, &user.Sex, &user.Birthday, &user.Passwd, &user.Email, &user.WeixinOpenid, &user.AvatarUrl, &user.Extra, &user.CreateTime, &user.LastLoginTime)
	if err != nil && err != sql.ErrNoRows {
		return nil, common.ErrDBError
	}

	if err == sql.ErrNoRows {
		return nil, nil
	}

	return &user, err
}

func (mya *MysqlAdapter) UpdateUser(user *User) error {
	_, err := mya.MysqlClient.Exec("update user set nickname = ?,sex = ?,avatar_url = ?,extra = ? where app_id = ? and user_id = ?",
		user.Nickname, user.Sex, user.AvatarUrl, user.Extra, user.AppId, user.UserId)
	if err != nil {
		return common.ErrDBError
	}

	return nil
}

func (mya *MysqlAdapter) GetAppInfo(appId int64) (*App, error) {
	var app App
	err := mya.MysqlClient.QueryRow("select id,name,private_key,create_time,update_time from app where id = ?", appId).Scan(
		&app.Id, &app.Name, &app.PrivateKey, &app.CreateTime, &app.UpdateTime)
	if err != nil && err != sql.ErrNoRows {
		return nil, common.ErrDBError
	}

	if err == sql.ErrNoRows {
		return nil, nil
	}
	return &app, nil
}

func (mya *MysqlAdapter) AddDevice(device Device) (id int64, err error) {
	res, err := mya.MysqlClient.Exec(`insert into device(device_id,app_id,type,brand,model,system_version,sdk_version,status,conn_id,user_ip) values(?,?,?,?,?,?,?,?,?,?)`, device.DeviceId, device.AppId, device.Type, device.Brand, device.Model, device.SystemVersion, device.SDKVersion, device.Status, "", 0)
	if err != nil {
		return 0, common.ErrDBError
	}
	id, err = res.LastInsertId()
	if err != nil {
		common.Sugar.Error("get last insert id error: ", err)
	}
	common.Sugar.Info("register device :", id)

	return
}

func (mya *MysqlAdapter) GetDevice(deviceId int64) (*Device, error) {
	device := Device{
		DeviceId: deviceId,
	}
	row := mya.MysqlClient.QueryRow(`
		select app_id,user_id,type,brand,model,system_version,sdk_version,status,conn_id,user_ip,create_time,update_time
		from device where device_id = ?`, deviceId)
	err := row.Scan(&device.AppId, &device.UserId, &device.Type, &device.Brand, &device.Model, &device.SystemVersion, &device.SDKVersion,
		&device.Status, &device.ConnId, &device.UserIp, &device.CreateTime, &device.UpdateTime)
	if err != nil && err != sql.ErrNoRows {
		return nil, common.ErrDBError
	}

	if err == sql.ErrNoRows {
		return nil, nil
	}

	return &device, err
}

func (mya *MysqlAdapter) UpdateDevice(device Device) error {
	return fmt.Errorf("unreachable")
}

func (mya *MysqlAdapter) DeleteDevice(deviceId int64) error {
	return fmt.Errorf("unreachable")
}
func (mya *MysqlAdapter) ListOnlineByUserId(appId, userId int64) ([]Device, error) {
	rows, err := mya.MysqlClient.Query(
		`select device_id,type,brand,model,system_version,sdk_version,status,conn_id,user_ip,create_time,update_time from device where app_id = ? and user_id = ? and status = ?`,
		appId, userId, DeviceOnLine)
	if err != nil {
		return nil, common.ErrDBError
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
func (mya *MysqlAdapter) UpdateDeviceStatus(deviceId int64, status int) error {
	_, err := mya.MysqlClient.Exec("update device set status = ? where device_id = ?", status, deviceId)
	if err != nil {
		return common.ErrDBError
	}
	return nil
}

func (mya *MysqlAdapter) InitUserSeq(userSeq UserSeq) error {
	_, err := mya.MysqlClient.Exec("insert into user_seq(app_id,user_id) values(?,?)", userSeq.AppId, userSeq.UserId)
	if err != nil {
		return common.ErrDBError
	}
	return nil
}

func (mya *MysqlAdapter) GetUserAck(appId, userId, groupId int64) (ack int64, err error) {
	row := mya.MysqlClient.QueryRow(`select read_seq from user_seq where app_id = ? and group_id = ? and user_id = ?`, appId, groupId, userId)

	err = row.Scan(&ack)
	if err != nil {
		common.Sugar.Errorf("query db err:", err)
		return 0, common.ErrDBError
	}
	return ack, nil
}

func (mya *MysqlAdapter) GetUserSeq(appId, userId, groupId int64) (seq int64, err error) {
	row := mya.MysqlClient.QueryRow(`select receive_seq from user_seq where app_id = ? and group_id = ? and user_id = ?`, appId, groupId, userId)

	err = row.Scan(&seq)
	if err != nil {
		common.Sugar.Errorf("query db err:", err)
		return 0, common.ErrDBError
	}
	return seq, nil
}

func (mya *MysqlAdapter) UpdateAck(appId, userId, groupId, ack int64) error {
	_, err := mya.MysqlClient.Exec("update user_seq set read_seq = ? where group_id=? and app_id=? and  user_id = ?", ack, groupId, appId, userId)
	if err != nil {
		common.Sugar.Errorf("update user_seq error:%v", err)
		return common.ErrDBError
	}
	return nil
}

func (mya *MysqlAdapter) GetUserNextSeq(appId, userId, groupId int64) (int64, error) {
	tx, err := mya.MysqlClient.Begin()
	rsp, err := mya.MysqlClient.Exec("update user_seq set receive_seq = receive_seq+1 where group_id=? and app_id=? and  user_id = ?", groupId, appId, userId)
	if err != nil {
		err = tx.Rollback()
		return 0, common.ErrDBError
	}
	err = tx.Commit()
	updateNum, err := rsp.RowsAffected()
	if err != nil {
		return 0, err
	}
	if updateNum == 0 {
		return 0, fmt.Errorf("record not exist")
	}
	seq, err := mya.GetUserSeq(appId, userId, groupId)

	return seq, err
}

func (mya *MysqlAdapter) AddGroup(group Group) (int64, error) {
	result, err := mya.MysqlClient.Exec("insert ignore into `group`(app_id,group_id,name,introduction,type,extra) value(?,?,?,?,?,?)",
		group.AppId, group.GroupId, group.Name, group.Introduction, group.Type, group.Extra)
	if err != nil {
		common.Sugar.Error(err)
		return 0, err
	}
	num, err := result.RowsAffected()
	if err != nil {
		common.Sugar.Errorf("query db err:", err)
		return 0, common.ErrDBError
	}
	return num, nil
}

func (mya *MysqlAdapter) UpdateGroup(group Group) (int64, error) {
	_, err := mya.MysqlClient.Exec("update `group` set name = ?,introduction = ?,extra = ? where app_id = ? and group_id = ?",
		group.Name, group.Introduction, group.Extra, group.AppId, group.GroupId)
	if err != nil {
		common.Sugar.Errorf("query db err:", err)
		return 0, common.ErrDBError
	}

	return 1, nil
}

func (mya *MysqlAdapter) GetGroup(appId, groupId int64) (*Group, error) {
	row := mya.MysqlClient.QueryRow("select name,introduction,user_num,type,extra,create_time,update_time from `group` where app_id = ? and group_id = ?",
		appId, groupId)
	group := Group{
		AppId:   appId,
		GroupId: groupId,
	}
	err := row.Scan(&group.Name, &group.Introduction, &group.UserNum, &group.Type, &group.Extra, &group.CreateTime, &group.UpdateTime)
	if err != nil && err != sql.ErrNoRows {
		return nil, common.ErrDBError
	}

	if err == sql.ErrNoRows {
		return nil, nil
	}
	return &group, nil
}

func (mya *MysqlAdapter) DeleteGroup(appId, groupId int64) error {
	return fmt.Errorf("unreachable")
}

func (mya *MysqlAdapter) AddMember(groupId int64, user *GroupUserInfo) error {
	_, err := mya.MysqlClient.Exec("insert ignore into group_user(app_id,group_id,user_id,label,extra) values(?,?,?,?,?)",
		user.AppId, user.GroupId, user.UserId, user.Label, user.UserExtra)
	if err != nil {
		return common.ErrDBError
	}
	return nil
}

func (mya *MysqlAdapter) GetMembers(appId, groupId int64) ([]*GroupUserInfo, error) {
	rows, err := mya.MysqlClient.Query(`
		select user_id,label,extra,create_time,update_time 
		from group_user
		where app_id = ? and group_id = ?`, appId, groupId)
	if err != nil {
		return nil, common.ErrDBError
	}
	groupUsers := make([]*GroupUserInfo, 0, 5)
	for rows.Next() {
		var groupUser = GroupUserInfo{
		}
		err := rows.Scan(&groupUser.UserId, &groupUser.Label, &groupUser.UserExtra, &groupUser.CreateTime, &groupUser.UpdateTime)
		if err != nil {
			return nil, common.ErrDBError
		}
		groupUsers = append(groupUsers, &groupUser)
	}
	return groupUsers, nil
}

func (mya *MysqlAdapter) DeleteMember(appId int64, groupId int64, userId int64) error {
	_, err := mya.MysqlClient.Exec("delete from group_user where app_id = ? and group_id = ? and user_id = ?",
		appId, groupId, userId)
	if err != nil {
		return common.ErrDBError
	}
	return nil
}

func (mya *MysqlAdapter) UpdateMember(user *GroupUserInfo) error {
	_, err := mya.MysqlClient.Exec("update group_user set label = ?,extra = ? where app_id = ? and group_id = ? and user_id = ?",
		user.Label, user.UserExtra, user.AppId, user.GroupId, user.UserId)
	if err != nil {
		return common.ErrDBError
	}
	return nil
}

// ListByUser 获取用户加入的群组信息
func (mya *MysqlAdapter) ListUserJoinGroup(appId, userId int64) ([]Group, error) {
	rows, err := mya.MysqlClient.Query(
		"select g.group_id,g.name,g.introduction,g.user_num,g.type,g.extra,g.create_time,g.update_time "+
			"from group_user u left join `group` g on u.app_id = g.app_id and u.group_id = g.group_id "+
			"where u.app_id = ? and u.user_id = ?",
		appId, userId)
	if err != nil {
		return nil, common.ErrDBError
	}
	var groups []Group
	var group Group
	for rows.Next() {
		err := rows.Scan(&group.GroupId, &group.Name, &group.Introduction, &group.UserNum, &group.Type, &group.Extra, &group.CreateTime, &group.UpdateTime)
		if err != nil {
			return nil, common.ErrDBError
		}
		groups = append(groups, group)
	}
	return groups, nil
}

func (mya *MysqlAdapter) AddMessage(message *Message) error {
	sqlraw := fmt.Sprintf(`insert into message (app_id,object_type,object_id,message_id,sender_type,sender_id,sender_device_id,receiver_type,receiver_id,
			to_user_ids,type,content,seq,send_time) values(?,?,?,?,?,?,?,?,?,?,?,?,?,?)`)
	_, err := mya.MysqlClient.Exec(sqlraw, message.AppId, message.ObjectType, message.ObjectId, message.MessageId, message.SenderType, message.SenderId,
		message.SenderDeviceId, message.ReceiverType, message.ReceiverId, message.ToUserIds, message.Type, message.Content, message.Seq, message.SendTime)
	if err != nil {
		return common.ErrDBError
	}
	return nil
}

func (mya *MysqlAdapter) CancelMessage(msgId int64) error {
	return fmt.Errorf("unreachable")
}

func (mya *MysqlAdapter) ListMsgBySeq(appId, receiverId, seq int64) ([]Message, error) {
	sqlraw := `select app_id,object_type,object_id,message_id,sender_type,sender_id,sender_device_id,receiver_type,receiver_id,
		to_user_ids,type,content,seq,send_time from message where app_id = ? and receiver_id = ? and seq > ? order by seq limit  300`
	rows, err := mya.MysqlClient.Query(sqlraw, appId, receiverId, seq)
	if err != nil {
		return nil, common.ErrDBError
	}

	messages := make([]Message, 0, 10)
	for rows.Next() {
		message := new(Message)
		err := rows.Scan(&message.AppId, &message.ObjectType, &message.ObjectId, &message.MessageId, &message.SenderType, &message.SenderId,
			&message.SenderDeviceId, &message.ReceiverType, &message.ReceiverId, &message.ToUserIds, &message.Type, &message.Content, &message.Seq, &message.SendTime)
		if err != nil {
			return nil, common.ErrDBError
		}
		messages = append(messages, *message)
	}
	return messages, nil
}

func (mya *MysqlAdapter) AddMoments() error {
	return fmt.Errorf("unreachable")
}

func (mya *MysqlAdapter) GetMoments() error {
	return fmt.Errorf("unreachable")
}
