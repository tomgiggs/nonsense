package service

import (
	"context"
	"go.uber.org/zap"
	"nonsense/internal/global"
	"nonsense/internal/logic/cache"
	"nonsense/internal/logic/dao"
	"nonsense/pkg/common"
	"nonsense/pkg/grpclib"
	pb "nonsense/pkg/proto"
)

type MessageService struct{}
func InitMessageService()*MessageService{
	return &MessageService{}
}
// 未读消息查询
func (self *MessageService) ListByUserIdAndSeq(ctx context.Context, appId, userId, seq int64) (messages []dao.Message,err error) {
	if seq == 0 {
		seqInfo := UserServiceInst.GetUserMaxACK(appId,userId,0)
		seq = seqInfo.data.(int64)
		if seqInfo.code != global.REQ_RESULT_CODE_OK {
			common.Sugar.Errorf("get user seq failed")
			return nil, global.DB_ERROR
		}
	}
	messages, err = dao.Storage.ListMsgBySeq( appId, userId, seq)
	if err != nil {
		return nil, err
	}
	return messages, nil
}

// 消息发送
func (self *MessageService) Send(ctx context.Context, sender dao.Sender, req pb.SendMessageReq) (err error) {
	switch req.ReceiverType {
	case pb.ReceiverType_RT_USER:
		if sender.SenderType == pb.SenderType_ST_USER {
			err = self.SendToFriend(ctx, sender, req)
			if err != nil {
				return err
			}
		} else {
			err = self.SendToUser(ctx, sender, req.ReceiverId, 0, req)
			if err != nil {
				return err
			}
		}
	case pb.ReceiverType_RT_NORMAL_GROUP:
		err = self.SendToGroup(ctx, sender, req,false)
		if err != nil {
			return err
		}
	case pb.ReceiverType_RT_LARGE_GROUP:
		err = self.SendToGroup(ctx, sender, req,true)
		if err != nil {
			return err
		}
	}

	return nil
}

// 好友消息
func (self *MessageService) SendToFriend(ctx context.Context, sender dao.Sender, req pb.SendMessageReq) error {

	// 发给发送者
	err := self.SendToUser(ctx, sender, sender.SenderId, 0, req)
	if err != nil {
		return err
	}

	// 发给接收者
	err = self.SendToUser(ctx, sender, req.ReceiverId, 0, req)
	if err != nil {
		return err
	}

	return nil
}

// 普通群组--写扩散
func (self *MessageService) SendToGroup(ctx context.Context, sender dao.Sender, req pb.SendMessageReq,isLargeGroup bool) error {
	isMember, err := cache.CacheInst.IsGroupMember(sender.AppId, req.ReceiverId, sender.SenderId)
	if err != nil {
		return err
	}
	if sender.SenderType == pb.SenderType_ST_USER && !isMember {
		common.Logger.Error("not int group", zap.Int64("app_id", sender.AppId), zap.Int64("group_id", req.ReceiverId),
			zap.Int64("user_id", sender.AppId))
		return common.ErrNotInGroup
	}

	users, err := cache.CacheInst.GetGroupMembers(sender.AppId, req.ReceiverId)
	if err != nil {
		return err
	}
	var seq int64 = 0
	if req.IsPersist {
		seq, err = SeqServiceInst.GetGroupNextSeq(sender.AppId, req.ReceiverId)
		if err != nil {
			return err
		}
		msg := self.PB2DB(req,seq,sender)
		err =dao.Storage.AddMessage(&msg)
		if err != nil {
			return err
		}
	}
	// 根据是否大群使用读/写扩散
	if isLargeGroup{
		req.IsPersist = false
	}
	for _, user := range users {
		err = self.SendToUser(ctx, sender, user.UserId, 0, req)
		if err != nil {
			return err
		}
	}
	return nil
}

// 将消息持久化到数据库,并且消息发送至用户
func (self *MessageService) SendToUser(ctx context.Context, sender dao.Sender, toUserId int64, roomSeq int64, req pb.SendMessageReq) error {
	common.Logger.Info("message_store_send_to_user",
		zap.String("message_id", req.MessageId),
		zap.Int64("app_id", sender.AppId),
		zap.Int64("to_user_id", toUserId),
		zap.Any("msg_body",req.MessageBody.MessageContent.Content))

	var (
		seq = roomSeq
		err error
	)
	if req.IsPersist {
		seq, err = SeqServiceInst.GetUserNextSeqFromDB(sender.AppId, toUserId)
		if err != nil {
			common.Logger.Error("message route",zap.Any("get user seq error",err))
			return err
		}
		msg := self.PB2DB(req,seq,sender)
		err =dao.Storage.AddMessage(&msg)
		if err != nil {
			common.Logger.Error("message route",zap.Any("save msg to db error",err))
			return err
		}
	}

	messageItem := &pb.MessageItem{
		MessageId:      req.MessageId,
		SenderType:     sender.SenderType,
		SenderId:       sender.SenderId,
		SenderDeviceId: sender.DeviceId,
		ReceiverType:   req.ReceiverType,
		ReceiverId:     req.ReceiverId,
		ToUserIds:      req.ToUserIds,
		MessageBody:    req.MessageBody,
		Seq:            seq,
		SendTime:       req.SendTime,
		Status:         pb.MessageStatus_MS_NORMAL,
	}
	//将不在本机的用户设备经过rpc投递消息
	for _,client := range self.FindUserConnServer(req.ReceiverId){
		if client ==nil {
			common.Sugar.Errorf("find rpc client nil")
			continue
		}
		client.DeliverMessage(context.TODO(), &pb.DeliverMessageReq{
			UserId: req.ReceiverId,
			Message: messageItem,
		})
	}
	// 查询本地用户在线设备
	self.SendToDevice(ctx,req.ReceiverId,messageItem)
	return nil
}

// 将消息发送给设备
func (self *MessageService) SendToDevice(ctx context.Context, userId int64, msgItem *pb.MessageItem) error {
	userMap := global.UserFdMap[userId]
	if userMap == nil {
		return nil
	}
	for fd,_ := range userMap{
		connection, ok := global.TcpServer.GetConn(fd)	// 获取设备对应的TCP连接
		if !ok {
			common.Logger.Warn("GetConn warn", zap.Int64("user_id", userId))
			continue
		}
		global.SendToClient(connection, pb.PackageType_PT_SYNC, grpclib.GetCtxRequstId(ctx), nil, msgItem)
	}
	return nil
}

//找出用户分布在哪几台连接服务器上
func (self *MessageService) FindUserConnServer(uid int64) []pb.LogicDispatchClient {
	serverList := make([]pb.LogicDispatchClient,0)
	if oneCache,ok := cache.UserServerMap[uid];ok{
		for k,_ := range oneCache{
			srv := global.LogicDispatchMap[k]
			serverList = append(serverList,srv)
		}
	}else {
		record := cache.CacheInst.GetUserServerFromRedis(uid)
		for k,_ := range record{
			srv := global.LogicDispatchMap[k]
			serverList = append(serverList,srv)
		}
	}
	return serverList
}

func (self *MessageService) PB2DB(req pb.SendMessageReq,seq int64,sender dao.Sender) dao.Message {
	messageType, messageContent := dao.PBToJsonStr(req.MessageBody)
	selfMessage := dao.Message{
		AppId:          sender.AppId,
		ObjectType:     dao.MessageObjectTypeUser,
		ObjectId:       req.ReceiverId,
		MessageId:      req.MessageId,
		SenderType:     int32(sender.SenderType),
		SenderId:       sender.SenderId,
		SenderDeviceId: sender.DeviceId,
		ReceiverType:   int32(req.ReceiverType),
		ReceiverId:     req.ReceiverId,
		ToUserIds:      dao.FormatUserIds(req.ToUserIds),
		Type:           messageType,
		Content:        messageContent,
		Seq:            seq,
		SendTime:       common.UnunixMilliTime(req.SendTime),
		Status:         int32(pb.MessageStatus_MS_NORMAL),
	}

	return selfMessage
}