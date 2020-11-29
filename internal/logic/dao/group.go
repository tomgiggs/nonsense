package dao

import (
	"database/sql"
	"nonsense/internal/global"
	"nonsense/pkg/common"
	"nonsense/pkg/storage"
)

type GroupDao struct{
	daoClient *storage.DBClient
}

func InitGroupDao()*GroupDao{
	return &GroupDao{
		daoClient: global.StorageClient,
	}
}

// Get 获取群组信息
func (self *GroupDao) Get(appId, groupId int64) (*Group, error) {
	row := global.StorageClient.MysqlClient.QueryRow("select name,introduction,user_num,type,extra,create_time,update_time from `group` where app_id = ? and group_id = ?",
		appId, groupId)
	group := Group{
		AppId:   appId,
		GroupId: groupId,
	}
	err := row.Scan(&group.Name, &group.Introduction, &group.UserNum, &group.Type, &group.Extra, &group.CreateTime, &group.UpdateTime)
	if err != nil && err != sql.ErrNoRows {
		return nil, common.WrapError(err)
	}

	if err == sql.ErrNoRows {
		return nil, nil
	}
	return &group, nil
}

// Insert 插入一条群组
func (self *GroupDao) Add(group Group) (int64, error) {
	result, err := global.StorageClient.MysqlClient.Exec("insert ignore into `group`(app_id,group_id,name,introduction,type,extra) value(?,?,?,?,?,?)",
		group.AppId, group.GroupId, group.Name, group.Introduction, group.Type, group.Extra)
	if err != nil {
		common.Sugar.Error(err)
		return 0, err
	}
	num, err := result.RowsAffected()
	if err != nil {
		return 0, common.WrapError(err)
	}
	return num, nil
}

// Update 更新群组信息
func (self *GroupDao) Update(appId, groupId int64, name, introduction, extra string) error {
	_, err := global.StorageClient.MysqlClient.Exec("update `group` set name = ?,introduction = ?,extra = ? where app_id = ? and group_id = ?",
		name, introduction, extra, appId, groupId)
	if err != nil {
		return common.WrapError(err)
	}

	return nil
}

// AddUserNum 更新群组信息
func (self *GroupDao) AddUserNum(appId, groupId int64, userNum int) error {
	_, err := global.StorageClient.MysqlClient.Exec("update `group` set user_num = user_num + ? where app_id = ? and group_id = ?",
		userNum, appId, groupId)
	if err != nil {
		return common.WrapError(err)
	}

	return nil
}

// UpdateUserNum 更新群组群成员人数
func (self *GroupDao) UpdateUserNum(appId, groupId, userNum int64) error {
	_, err := global.StorageClient.MysqlClient.Exec("update `group` set user_num = user_num + ? where app_id = ? and group_id = ?",
		userNum, appId, groupId)
	if err != nil {
		return common.WrapError(err)
	}

	return nil
}

//--------------

type GroupUserDao struct{
	daoClient *storage.DBClient
}

func InitGroupUserDao()*GroupUserDao{
	return &GroupUserDao{
		daoClient:global.StorageClient,
	}
}
// ListByUser 获取用户加入的群组信息
func (self *GroupUserDao) ListByUserId(appId, userId int64) ([]Group, error) {
	rows, err := global.StorageClient.MysqlClient.Query(
		"select g.group_id,g.name,g.introduction,g.user_num,g.type,g.extra,g.create_time,g.update_time "+
			"from group_user u left join `group` g on u.app_id = g.app_id and u.group_id = g.group_id "+
			"where u.app_id = ? and u.user_id = ?",
		appId, userId)
	if err != nil {
		return nil, common.WrapError(err)
	}
	var groups []Group
	var group Group
	for rows.Next() {
		err := rows.Scan(&group.GroupId, &group.Name, &group.Introduction, &group.UserNum, &group.Type, &group.Extra, &group.CreateTime, &group.UpdateTime)
		if err != nil {
			return nil, common.WrapError(err)
		}
		groups = append(groups, group)
	}
	return groups, nil
}

// ListGroupUser 获取群组用户信息
func (self *GroupUserDao) ListUser(appId, groupId int64) ([]GroupUserInfo, error) {
	rows, err := global.StorageClient.MysqlClient.Query(`
		select user_id,label,extra,create_time,update_time 
		from group_user
		where app_id = ? and group_id = ?`, appId, groupId)
	if err != nil {
		return nil, common.WrapError(err)
	}
	groupUsers := make([]GroupUserInfo, 0, 5)
	for rows.Next() {
		var groupUser = GroupUserInfo{
		}
		err := rows.Scan(&groupUser.UserId, &groupUser.Label, &groupUser.UserExtra, &groupUser.CreateTime, &groupUser.UpdateTime)
		if err != nil {
			return nil, common.WrapError(err)
		}
		groupUsers = append(groupUsers, groupUser)
	}
	return groupUsers, nil
}

// GetGroupUser 获取群组用户信息,用户不存在返回nil
func (self *GroupUserDao) Get(appId, groupId, userId int64) (*GroupUserInfo, error) {
	var groupUser = GroupUserInfo{
		AppId:   appId,
		GroupId: groupId,
		UserId:  userId,
	}
	err := global.StorageClient.MysqlClient.QueryRow("select label,extra from group_user where app_id = ? and group_id = ? and user_id = ?",
		appId, groupId, userId).
		Scan(&groupUser.Label, &groupUser.UserExtra)
	if err != nil && err != sql.ErrNoRows {
		return nil, common.WrapError(err)
	}

	if err == sql.ErrNoRows {
		return nil, nil
	}

	return &groupUser, nil
}

// Add 将用户添加到群组
func (self *GroupUserDao) Add(appId, groupId, userId int64, label, extra string) error {
	_, err := global.StorageClient.MysqlClient.Exec("insert ignore into group_user(app_id,group_id,user_id,label,extra) values(?,?,?,?,?)",
		appId, groupId, userId, label, extra)
	if err != nil {
		return common.WrapError(err)
	}
	return nil
}

// Delete 将用户从群组删除
func (self *GroupUserDao) Delete(appId int64, groupId int64, userId int64) error {
	_, err := global.StorageClient.MysqlClient.Exec("delete from group_user where app_id = ? and group_id = ? and user_id = ?",
		appId, groupId, userId)
	if err != nil {
		return common.WrapError(err)
	}
	return nil
}

// Update 更新用户群组信息
func (self *GroupUserDao) Update(appId, groupId, userId int64, label string, extra string) error {
	_, err := global.StorageClient.MysqlClient.Exec("update group_user set label = ?,extra = ? where app_id = ? and group_id = ? and user_id = ?",
		label, extra, appId, groupId, userId)
	if err != nil {
		return common.WrapError(err)
	}
	return nil
}

