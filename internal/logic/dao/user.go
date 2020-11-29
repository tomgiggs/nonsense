package dao

import (
	"database/sql"
	"nonsense/internal/global"
	"nonsense/pkg/common"
	"nonsense/pkg/storage"
)

type UserDao struct{
	daoClient *storage.DBClient
}

func InitUserDao()*UserDao{
	return &UserDao{
		daoClient: global.StorageClient,
	}
}
// 创建新用户
func (self *UserDao) Add(user User) (int64, error) {
	result, err := global.StorageClient.MysqlClient.Exec("insert ignore into user(app_id,user_id,nickname,sex,avatar_url,extra) values(?,?,?,?,?,?)",
		user.AppId, user.UserId, user.Nickname, user.Sex, user.AvatarUrl, user.Extra)
	if err != nil {
		return 0, common.WrapError(err)
	}

	lastId, err := result.LastInsertId()
	if err != nil {
		return 0, common.WrapError(err)
	}
	return lastId, nil
}

// 获取用户信息
func (self *UserDao) Get(appId, userId int64) (*User, error) {
	row := global.StorageClient.MysqlClient.QueryRow(
		"select nickname,sex,birthday,passwd,email,weixin_openid,avatar_url,extra,create_time,last_login_time from user where app_id = ? and user_id = ?",
		appId, userId)
	user := User{
		AppId:  appId,
		UserId: userId,
	}

	err := row.Scan(&user.Nickname, &user.Sex,&user.Birthday,&user.Passwd,&user.Email,&user.WeixinOpenid, &user.AvatarUrl, &user.Extra, &user.CreateTime, &user.LastLoginTime)
	if err != nil && err != sql.ErrNoRows {
		return nil, common.WrapError(err)
	}

	if err == sql.ErrNoRows {
		return nil, nil
	}

	return &user, err
}


// 更新用户信息
func (self *UserDao) Update(user User) error {
	_, err := global.StorageClient.MysqlClient.Exec("update user set nickname = ?,sex = ?,avatar_url = ?,extra = ? where app_id = ? and user_id = ?",
		user.Nickname, user.Sex, user.AvatarUrl, user.Extra, user.AppId, user.UserId)
	if err != nil {
		return common.WrapError(err)
	}

	return nil
}
