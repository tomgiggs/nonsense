package dao

import (
	"fmt"
	"nonsense/internal/global"
	"nonsense/pkg/common"
	pb "nonsense/pkg/proto"
	"nonsense/pkg/storage"
)

type MessageDao struct{
	daoClient *storage.DBClient
}

func InitMessageDao()*MessageDao{
	return &MessageDao{
		daoClient: global.StorageClient,
	}
}
// Add 插入一条消息
func (self *MessageDao) Add(tableName string, message Message) error {
	sql := fmt.Sprintf(`insert into %s(app_id,object_type,object_id,message_id,sender_type,sender_id,sender_device_id,receiver_type,receiver_id,
			to_user_ids,type,content,seq,send_time) values(?,?,?,?,?,?,?,?,?,?,?,?,?,?)`, tableName)
	_, err := global.StorageClient.MysqlClient.Exec(sql, message.AppId, message.ObjectType, message.ObjectId, message.MessageId, message.SenderType, message.SenderId,
		message.SenderDeviceId, message.ReceiverType, message.ReceiverId, message.ToUserIds, message.Type, message.Content, message.Seq, message.SendTime)
	if err != nil {
		return common.WrapError(err)
	}
	return nil
}

// 根据appid,userid,查询序号大于seq的消息,传入表名是为了后续做分表
func (self *MessageDao) ListBySeq(appId, userId, seq int64) ([]Message, error) {
	sql := `select app_id,object_type,object_id,message_id,sender_type,sender_id,sender_device_id,receiver_type,receiver_id,
		to_user_ids,type,content,seq,send_time from message where app_id = ? and receiver_id = ? and seq > ? order by seq limit  300`
	rows, err := global.StorageClient.MysqlClient.Query(sql, appId, userId, seq)
	if err != nil {
		return nil, common.WrapError(err)
	}

	messages := make([]Message, 0, 10)
	for rows.Next() {
		message := new(Message)
		err := rows.Scan(&message.AppId, &message.ObjectType, &message.ObjectId, &message.MessageId, &message.SenderType, &message.SenderId,
			&message.SenderDeviceId, &message.ReceiverType, &message.ReceiverId, &message.ToUserIds, &message.Type, &message.Content, &message.Seq, &message.SendTime)
		if err != nil {
			return nil, common.WrapError(err)
		}
		messages = append(messages, *message)
	}
	return messages, nil
}
func (self *MessageDao) ListHugeGroupMessageBySeq(appId, objectType, objectId, seq int64) ([]Message, error) {
	sql := `select app_id,object_type,object_id,message_id,sender_type,sender_id,sender_device_id,receiver_type,receiver_id,
		to_user_ids,type,content,seq,send_time from message where app_id = ? and object_type = ? and object_id = ? and seq > ?`
	rows, err := global.StorageClient.MysqlClient.Query(sql, appId, objectType, objectId, seq)
	if err != nil {
		return nil, common.WrapError(err)
	}

	messages := make([]Message, 0, 5)
	for rows.Next() {
		message := new(Message)
		err := rows.Scan(&message.AppId, &message.ObjectType, &message.ObjectId, &message.MessageId, &message.SenderType, &message.SenderId,
			&message.SenderDeviceId, &message.ReceiverType, &message.ReceiverId, &message.ToUserIds, &message.Type, &message.Content, &message.Seq, &message.SendTime)
		if err != nil {
			return nil, common.WrapError(err)
		}
		messages = append(messages, *message)
	}
	return messages, nil
}


func (self *MessageDao) PB2DB(req pb.SendMessageReq,seq int64,sender Sender) error{
	messageType, messageContent := PBToJsonStr(req.MessageBody)
	selfMessage := Message{
		AppId:          sender.AppId,
		ObjectType:     MessageObjectTypeUser,
		ObjectId:       req.ReceiverId,
		MessageId:      req.MessageId,
		SenderType:     int32(sender.SenderType),
		SenderId:       sender.SenderId,
		SenderDeviceId: sender.DeviceId,
		ReceiverType:   int32(req.ReceiverType),
		ReceiverId:     req.ReceiverId,
		ToUserIds:      FormatUserIds(req.ToUserIds),
		Type:           messageType,
		Content:        messageContent,
		Seq:            seq,
		SendTime:       common.UnunixMilliTime(req.SendTime),
		Status:         int32(pb.MessageStatus_MS_NORMAL),
	}

	return self.Add("message", selfMessage)

}