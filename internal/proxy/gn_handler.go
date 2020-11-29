package proxy

import (
	"context"
	"fmt"
	"github.com/alberliu/gn"
	"github.com/golang/protobuf/proto"
	"go.uber.org/zap"
	"nonsense/internal/global"
	"nonsense/internal/logic/cache"
	"nonsense/internal/logic/dao"
	"nonsense/internal/logic/service"
	"nonsense/pkg/common"
	pb "nonsense/pkg/proto"
)

type ConnData struct {
	AppId    int64 // AppId
	DeviceId int64 // 设备id
	UserId   int64 // 用户id
}


type Handler struct{}

func (*Handler) OnConnect(c *gn.Conn) {
	common.Logger.Debug("connect:", zap.Int32("fd", c.GetFd()), zap.String("addr", c.GetAddr()))
}

func (h *Handler) OnMessage(c *gn.Conn, bytes []byte) {
	var input pb.Input
	err := proto.Unmarshal(bytes, &input)
	if err != nil {
		common.Logger.Error("unmarshal error", zap.Error(err))
		return
	}

	// 对未登录的用户进行拦截
	if input.Type != pb.PackageType_PT_SIGN_IN && c.GetData() == nil {
		// 应该告诉用户没有登录
		return
	}

	switch input.Type {
	case pb.PackageType_PT_SIGN_IN:
		h.SignIn(c, input)
	case pb.PackageType_PT_SYNC:
		h.Sync(c, input)
	case pb.PackageType_PT_HEARTBEAT:
		h.Heartbeat(c, input)
	case pb.PackageType_PT_MESSAGE_ACK:
		h.MessageACK(c, input)
	case pb.PackageType_PT_MESSAGE_SEND:
		h.MessageSend(c, input)
	default:
		common.Logger.Error("handler switch other")
	}
	return
}

func (*Handler) OnClose(c *gn.Conn, err error) {
	common.Logger.Debug("close", zap.Any("data", c.GetData()), zap.Error(err))
	data := c.GetData().(ConnData)
	service.DeviceServiceInst.Offline(context.TODO(), data.AppId, data.UserId, data.DeviceId)
	delete(global.UserFdMap[data.UserId],c.GetFd())
}

// SignIn 登录
func (*Handler) SignIn(c *gn.Conn, input pb.Input) {
	var signIn pb.SignInInput
	err := proto.Unmarshal(input.Data, &signIn)
	if err != nil {
		common.Sugar.Error(err)
		return
	}
	var token string
fmt.Println(signIn.AppId,signIn.UserId,signIn.DeviceId,signIn.Passwd,global.AppConfig.SrvDisc.ID,c.GetAddr())
	token,err = service.AuthServiceInst.SignIn(context.TODO(), signIn.AppId,signIn.UserId,signIn.DeviceId,signIn.Passwd,global.AppConfig.SrvDisc.ID,c.GetAddr())

	global.SendToClient(c, pb.PackageType_PT_SIGN_IN, input.RequestId, err, &pb.SignInResp{
		Token: token,
	})

	if err != nil{
		common.Sugar.Error(err)
	}
	//注册用户连接
	if global.UserFdMap[signIn.UserId]==nil{
		global.UserFdMap[signIn.UserId] = make(map[int32]int32)
	}
	global.UserFdMap[signIn.UserId][c.GetFd()] = 1
	data := ConnData{
		AppId:    signIn.AppId,
		DeviceId: signIn.DeviceId,
		UserId:   signIn.UserId,
	}
	c.SetData(data)
	changeInfo := &cache.UserChangeInfo{
		Event: "online",
		Uid:   signIn.UserId,
		SrvId: global.AppConfig.SrvDisc.ID,
	}
	cache.PubOnlineUserChange(changeInfo)
}

// Sync 消息同步
func (*Handler) Sync(c *gn.Conn, input pb.Input) {
	var sync pb.SyncInput
	err := proto.Unmarshal(input.Data, &sync)
	if err != nil {
		common.Sugar.Error(err)
		return
	}

	data := c.GetData().(ConnData)
	resp, err := service.MessageServiceInst.ListByUserIdAndSeq(context.TODO(), data.AppId,data.UserId,sync.Seq)
	if err != nil {
		return
	}
	var message =&pb.SyncOutput{Messages: dao.MessagesToPB(resp)}
	global.SendToClient(c, pb.PackageType_PT_SYNC, input.RequestId, err, message)
}

// Heartbeat 心跳
func (*Handler) Heartbeat(c *gn.Conn, input pb.Input) {
	data := c.GetData().(ConnData)
	global.SendToClient(c, pb.PackageType_PT_HEARTBEAT, input.RequestId, nil, nil)
	common.Sugar.Debugw("heartbeat", "device_id", data.DeviceId, "user_id", data.UserId)
}

// 消息收到回执
func (*Handler) MessageACK(c *gn.Conn, input pb.Input) {
	var messageACK pb.MessageACKReq
	err := proto.Unmarshal(input.Data, &messageACK)
	if err != nil {
		common.Sugar.Error(err)
		return
	}
	fmt.Println(messageACK)
	service.UserServiceInst.UpdateUserAckSeq(messageACK.AppId, messageACK.UserId, messageACK.GroupId,messageACK.Seq)

}

// 消息发送
func (*Handler) MessageSend(c *gn.Conn, input pb.Input) {
	var messageBody pb.SendMessageReq
	err := proto.Unmarshal(input.Data, &messageBody)
	if err != nil {
		common.Sugar.Error(err)
		return
	}

	data := c.GetData().(ConnData)
	sender := dao.Sender{
		AppId:      data.AppId,
		SenderType: pb.SenderType_ST_USER,
		SenderId:   data.UserId,
		DeviceId:   data.DeviceId,
	}
	service.MessageServiceInst.Send(context.TODO(),sender,messageBody)
}