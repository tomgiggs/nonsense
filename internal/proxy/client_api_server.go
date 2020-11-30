package proxy

import (
	"context"
	"nonsense/internal/global"
	"nonsense/internal/service"
	"nonsense/internal/store"
	"nonsense/pkg/common"
	pb "nonsense/pkg/proto"
)
type ClientApiServer struct{
	pb.UnimplementedLogicClientExtServer
}

// 设备登录
func (m *ClientApiServer) SignIn(ctx context.Context, req *pb.SignInReq) (*pb.SignInResp, error) {
	tokenStr,err:= service.AuthServiceInst.SignIn(ctx, req.AppId, req.UserId, req.DeviceId, req.Passwd, req.ConnId, req.UserIp)

	return &pb.SignInResp{
		Token: tokenStr,
	}, err
}

// 设备同步消息
func (m *ClientApiServer) Sync(ctx context.Context, req *pb.SyncReq) (*pb.SyncResp, error) {
	messages, err := service.MessageServiceInst.ListByUserIdAndSeq(ctx, req.AppId, req.UserId, req.Seq)
	if err != nil {
		return nil, err
	}
	return &pb.SyncResp{
		Messages: store.MessagesToPB(messages),
	}, nil
}

// 收到消息ack
func (m *ClientApiServer) MessageACK(ctx context.Context, req *pb.MessageACKReq) (*pb.MessageACKResp, error) {
	return &pb.MessageACKResp{}, service.UserServiceInst.UpdateUserAckSeq(req.AppId, req.UserId, req.GroupId,req.Seq)
}

// 设备离线
func (m *ClientApiServer) Offline(ctx context.Context, req *pb.OfflineReq) (*pb.OfflineResp, error) {
	return &pb.OfflineResp{}, service.DeviceServiceInst.Offline(ctx, req.AppId, req.UserId, req.DeviceId)
}

// 注册设备
func (m *ClientApiServer) RegisterDevice(ctx context.Context, in *pb.RegisterDeviceReq) (*pb.RegisterDeviceResp, error) {
	device := store.Device{
		AppId:         in.AppId,
		Type:          in.Type,
		Brand:         in.Brand,
		Model:         in.Model,
		SystemVersion: in.SystemVersion,
		SDKVersion:    in.SdkVersion,
	}

	if device.AppId == 0 || device.Type == 0 || device.Brand == "" || device.Model == "" ||
		device.SystemVersion == "" || device.SDKVersion == "" {
		return nil, common.ErrBadRequest
	}

	id, err := service.DeviceServiceInst.Register(ctx, device)
	if err != nil {
		return nil, err
	}
	return &pb.RegisterDeviceResp{DeviceId: id}, nil
}

// 添加用户
func (m *ClientApiServer) AddUser(ctx context.Context, in *pb.AddUserReq) (*pb.AddUserResp, error) {
	appId, _, _, err := common.GetCtxData(ctx)
	if err != nil {
		return &pb.AddUserResp{
			UserId: int64(0),
		}, err
	}

	user := store.User{
		AppId:     appId,
		Nickname:  in.User.Nickname,
		Sex:       in.User.Sex,
		AvatarUrl: in.User.AvatarUrl,
		Extra:     in.User.Extra,
	}
	var userId int64
	userId,err = service.UserServiceInst.Add(user)

	return &pb.AddUserResp{
		UserId: userId,
	},err
}

// 获取用户信息
func (m *ClientApiServer) GetUser(ctx context.Context, in *pb.GetUserReq) (*pb.GetUserResp, error) {
	appId, _, _, err := common.GetCtxData(ctx)
	if err != nil {
		return &pb.GetUserResp{}, err
	}

	user, err := service.UserServiceInst.Get(ctx, appId, in.UserId)
	if err != nil {
		return &pb.GetUserResp{}, nil
	}

	if user == nil {
		return nil, common.ErrUserNotExist
	}

	pbUser := &pb.User{
		UserId:     user.UserId,
		Nickname:   user.Nickname,
		Sex:        user.Sex,
		AvatarUrl:  user.AvatarUrl,
		Extra:      user.Extra,
		CreateTime: user.CreateTime.Unix(),
		UpdateTime: user.LastLoginTime.Unix(),
	}
	return &pb.GetUserResp{User: pbUser}, nil
}

// 发送消息
func (m *ClientApiServer) SendMessage(ctx context.Context, in *pb.SendMessageReq) (*pb.SendMessageResp,error) {
	appId, userId, deviceId, err := common.GetCtxData(ctx)
	rsp := &pb.SendMessageResp{
		ResultCode: global.REQ_RESULT_CODE_OK,
	}
	if err != nil {
		return nil, err
	}

	sender := store.Sender{
		AppId:      appId,
		SenderType: pb.SenderType_ST_USER,
		SenderId:   userId,
		DeviceId:   deviceId,
	}
	err = service.MessageServiceInst.Send(ctx, sender, *in)
	if err != nil {
		return nil, err
	}
	return rsp, nil
}

// 创建群组
func (m *ClientApiServer) CreateGroup(ctx context.Context, in *pb.CreateGroupReq) (*pb.CreateGroupResp, error) {
	rsp := &pb.CreateGroupResp{
		ResultCode: global.REQ_RESULT_CODE_FAIL,
	}
	appId, _, _, err := common.GetCtxData(ctx)
	if err != nil {
		return rsp, err
	}

	var group = store.Group{
		AppId:        appId,
		GroupId:      in.Group.GroupId,
		Name:         in.Group.Name,
		Introduction: in.Group.Introduction,
		Type:         in.Group.Type,
		Extra:        in.Group.Extra,
	}
	err = service.GroupServiceInst.Create(ctx, group)
	if err != nil {
		return rsp, err
	}
	rsp.ResultCode = global.REQ_RESULT_CODE_OK
	return rsp, nil
}

// 更新群组
func (m *ClientApiServer) UpdateGroup(ctx context.Context, in *pb.UpdateGroupReq) (*pb.UpdateGroupResp, error) {
	appId, _, _, err := common.GetCtxData(ctx)
	if err != nil {
		return nil, err
	}

	var group = store.Group{
		AppId:        appId,
		GroupId:      in.Group.GroupId,
		Name:         in.Group.Name,
		Introduction: in.Group.Introduction,
		Type:         in.Group.Type,
		Extra:        in.Group.Extra,
	}
	err = service.GroupServiceInst.Update(ctx, group)
	if err != nil {
		return nil, err
	}
	return &pb.UpdateGroupResp{}, nil
}

// 获取群组信息
func (m *ClientApiServer) GetGroup(ctx context.Context, in *pb.GetGroupReq) (*pb.GetGroupResp, error) {
	appId, _, _, err := common.GetCtxData(ctx)
	if err != nil {
		return nil, err
	}

	group, err := service.GroupServiceInst.Get(ctx, appId, in.GroupId)
	if err != nil {
		return nil, err
	}

	if group == nil {
		return nil, common.ErrGroupNotExist
	}

	return &pb.GetGroupResp{
		Group: &pb.Group{
			GroupId:      group.GroupId,
			Name:         group.Name,
			Introduction: group.Introduction,
			UserMum:      group.UserNum,
			Type:         group.Type,
			Extra:        group.Extra,
			CreateTime:   group.CreateTime.Unix(),
			UpdateTime:   group.UpdateTime.Unix(),
		},
	}, nil
}

// 获取用户加入的所有群组
func (m *ClientApiServer) GetUserGroups(ctx context.Context, in *pb.GetUserGroupsReq) (*pb.GetUserGroupsResp, error) {
	appId, userId, _, err := common.GetCtxData(ctx)
	if err != nil {
		return nil, err
	}

	groups, err := service.GroupServiceInst.ListUserJoinGroup(ctx, appId, userId)
	if err != nil {
		common.Sugar.Error(err)
		return nil, err
	}
	pbGroups := make([]*pb.Group, 0, len(groups))
	for i := range groups {
		pbGroups = append(pbGroups, &pb.Group{
			GroupId:      groups[i].GroupId,
			Name:         groups[i].Name,
			Introduction: groups[i].Introduction,
			UserMum:      groups[i].UserNum,
			Type:         groups[i].Type,
			Extra:        groups[i].Extra,
			CreateTime:   groups[i].CreateTime.Unix(),
			UpdateTime:   groups[i].UpdateTime.Unix(),
		})
	}
	return &pb.GetUserGroupsResp{Groups: pbGroups}, err
}

// 添加群组成员
func (m *ClientApiServer) AddGroupMember(ctx context.Context, in *pb.AddGroupMemberReq) (*pb.AddGroupMemberResp, error) {
	appId, _, _, err := common.GetCtxData(ctx)
	if err != nil {
		common.Sugar.Error(err)
		return nil, err
	}

	err = service.GroupServiceInst.AddUser(ctx, appId, in.GroupUser.GroupId, in.GroupUser.UserId, in.GroupUser.Label, in.GroupUser.Extra)
	if err != nil {
		common.Sugar.Error(err)
		return nil, err
	}

	return &pb.AddGroupMemberResp{}, nil
}

// 更新群组成员信息
func (m *ClientApiServer) UpdateGroupMember(ctx context.Context, in *pb.UpdateGroupMemberReq) (*pb.UpdateGroupMemberResp, error) {
	appId, _, _, err := common.GetCtxData(ctx)
	if err != nil {
		return nil, err
	}

	err = service.GroupServiceInst.UpdateUser(ctx, appId, in.GroupUser.GroupId, in.GroupUser.UserId, in.GroupUser.Label, in.GroupUser.Extra)
	if err != nil {
		return nil, err
	}

	return &pb.UpdateGroupMemberResp{}, nil
}

// 添加群组成员
func (m *ClientApiServer) DeleteGroupMember(ctx context.Context, in *pb.DeleteGroupMemberReq) (*pb.DeleteGroupMemberResp, error) {
	appId, _, _, err := common.GetCtxData(ctx)
	if err != nil {
		return nil, err
	}

	err = service.GroupServiceInst.DeleteUser(ctx, appId, in.GroupId, in.UserId)
	if err != nil {
		return nil, err
	}

	return &pb.DeleteGroupMemberResp{}, nil
}
