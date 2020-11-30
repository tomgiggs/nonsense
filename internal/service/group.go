package service

import (
	"context"
	"nonsense/internal/store"
	"nonsense/pkg/common"
)

type GroupService struct{
}
func InitGroupService() *GroupService {
	return &GroupService{
	}
}

// 获取群组信息
func (self *GroupService) Get(ctx context.Context, appId, groupId int64) (*store.Group, error) {
	group, err := store.CacheInst.GetGroup(appId, groupId)
	if err != nil {
		return nil, err
	}
	if group != nil {
		return group, nil
	}
	group, err = store.Storage.GetGroup(appId, groupId)
	if err != nil {
		return nil, err
	}

	if group == nil {
		return nil, nil
	}

	err = store.CacheInst.SetGroup(group)
	if err != nil {
		return nil, err
	}
	return group, nil
}

// 创建群组
func (self *GroupService) Create(ctx context.Context, group store.Group) error {
	affected, err := store.Storage.AddGroup(group)
	if err != nil {
		return err
	}

	if affected == 0 {
		return common.ErrGroupAlreadyExist
	}
	return nil
}

// 更新群组
func (self *GroupService) Update(ctx context.Context, group store.Group) error {
	_,err := store.Storage.UpdateGroup(group)
	if err != nil {
		return err
	}
	err = store.CacheInst.DelGroup(group.AppId, group.GroupId)
	if err != nil {
		return err
	}
	return nil
}

// 给群组添加用户
func (self *GroupService) AddUser(ctx context.Context, appId, groupId, userId int64, label, extra string) error {
	group, err := self.Get(ctx, appId, groupId)
	if err != nil {
		return err
	}
	if group == nil {
		return common.ErrGroupNotExist
	}

	user, err := UserServiceInst.Get(ctx, appId, userId)
	if err != nil {
		return err
	}
	if user == nil {
		return common.ErrUserNotExist
	}

	if group.Type == store.GroupTypeChatRoom {
		err = store.CacheInst.SetGroupMember(appId, groupId, userId, label, extra)
		if err != nil {
			return err
		}
	}
	return nil
}

// 更新群组用户
func (self *GroupService) UpdateUser(ctx context.Context, appId, groupId, userId int64, label, extra string) error {
	group, err := self.Get(ctx, appId, groupId)
	if err != nil {
		return err
	}

	if group == nil {
		return common.ErrGroupNotExist
	}

	err = store.CacheInst.SetGroupMember(appId, groupId, userId, label, extra)
	if err != nil {
		return err
	}

	return nil
}

// 删除用户群组
func (self *GroupService) DeleteUser(ctx context.Context, appId, groupId, userId int64) error {
	group, err := self.Get(ctx, appId, groupId)
	if err != nil {
		return err
	}

	if group == nil {
		return common.ErrGroupNotExist
	}

	err = store.CacheInst.DelGroupMember(appId, groupId, userId)
	if err != nil {
		return err
	}
	return nil
}

func (self *GroupService) IsMember(appId, groupId, userId int64) (bool, error) {
	is, err := store.CacheInst.IsGroupMember(appId,groupId,userId)
	return is, err
}
// GetUsers 获取群组的所有用户信息
func (self *GroupService) GetUsers(appId, groupId int64) ([]*store.GroupUserInfo, error) {
	users, err := store.CacheInst.GetGroupMembers(appId, groupId)
	if err != nil {
		return nil, err
	}

	if users != nil {
		return users, nil
	}

	users, err = store.Storage.GetMembers(appId, groupId)
	if err != nil {
		return nil, err
	}

	for _,user :=range users{
		err = store.CacheInst.SetGroupMember(appId, groupId,user.UserId,user.Label, user.UserExtra)
	}
	if err != nil {
		return nil, err
	}
	return users, err
}

// 获取用户所加入的群组
func (self *GroupService) ListUserJoinGroup(ctx context.Context, appId, userId int64) ([]store.Group, error) {
	groups, err := store.Storage.ListUserJoinGroup(appId, userId)
	if err != nil {
		return nil, err
	}
	return groups, nil
}

